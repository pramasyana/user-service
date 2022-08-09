package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	memberQuery "github.com/Bhinneka/user-service/src/member/v1/query"
	memberRepo "github.com/Bhinneka/user-service/src/member/v1/repo"
	"gopkg.in/ldap.v2"
)

const (
	defErrorMessage = "Invalid username or password"
)

//LDAPServiceImpl LDAPService implementation
type LDAPServiceImpl struct {
	ldapServer           string
	ldapBind             string
	ldapPasswordBind     string
	baseDN               string
	filterDN             string
	MemberAdditionalRepo memberRepo.MemberAdditionalInfoRepository
	MemberQueryRead      memberQuery.MemberQuery
}

//NewLDAPService LDAPService's constructor
func NewLDAPService(ldapServer,
	ldapBind,
	ldapPassword,
	baseDN,
	filterDN string,
	memberAdditional memberRepo.MemberAdditionalInfoRepository,
	memberQueryRead memberQuery.MemberQuery) (*LDAPServiceImpl, error) {
	return &LDAPServiceImpl{
		ldapServer:           ldapServer,
		ldapBind:             ldapBind,
		ldapPasswordBind:     ldapPassword,
		baseDN:               baseDN,
		filterDN:             filterDN,
		MemberAdditionalRepo: memberAdditional,
		MemberQueryRead:      memberQueryRead}, nil
}

//Auth function for validate LDAP user authentication
func (l *LDAPServiceImpl) Auth(ctxReq context.Context, loginUsername, loginPassword string) (*model.LDAPProfile, error) {
	ctx := "LDAPService"
	message := defErrorMessage

	trace := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{"username": loginUsername}
	defer trace.Finish(tags)

	conn, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", l.ldapServer, 389))

	if err != nil {
		return nil, err
	}

	if err := conn.Bind(l.ldapBind, l.ldapPasswordBind); err != nil {
		return nil, err
	}

	result, err := conn.Search(ldap.NewSearchRequest(
		l.baseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		0,
		false,
		l.filter(loginUsername),
		[]string{},
		nil,
	))

	if err != nil {
		tags["error_ldap"] = err.Error()
		return nil, fmt.Errorf(message)
	}

	if len(result.Entries) != 1 {
		return nil, fmt.Errorf(message)
	}

	if err := conn.Bind(result.Entries[0].DN, loginPassword); err != nil {
		if strings.Contains(err.Error(), "775") {
			// https://ldapwiki.com/wiki/Common%20Active%20Directory%20Bind%20Errors
			message = "Maximum failed login attempt reached, account is locked. Please contact IT Helpdesk"
		}
		return nil, fmt.Errorf(message)
	}

	response := &model.LDAPProfile{}
	attrResponse := result.Entries[0].Attributes

	for _, s := range attrResponse {
		switch s.Name {
		case "title":
			response.JobTitle = s.Values[0]
		case "givenName":
			response.FirstName = s.Values[0]
		case "sn":
			response.LastName = s.Values[0]
		case "displayName":
			response.DisplayName = s.Values[0]
		case "department":
			response.Department = s.Values[0]
		case "objectGUID":
			response.ObjectID = base64.StdEncoding.EncodeToString([]byte(s.Values[0]))
		case "mail":
			response.Email = s.Values[0]
		default:
		}
	}

	// save member additional
	go l.SaveMemberAdditionalInfo(ctxReq, result, response.Email)

	return response, nil
}

//filter private function for filtering DN by username
func (l *LDAPServiceImpl) filter(needle string) string {
	res := strings.Replace(
		l.filterDN,
		"{username}",
		needle,
		-1,
	)

	return res
}

// SaveMemberAdditionalInfo function saves member information
func (l *LDAPServiceImpl) SaveMemberAdditionalInfo(ctxReq context.Context, result *ldap.SearchResult, email string) {
	ctx := "LDAPService-SaveMemberAdditionalInfo"

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		authType := "LDAP"
		memberResult := <-l.MemberQueryRead.FindByEmail(ctxReq, email)
		member, _ := memberResult.Result.(memberModel.Member)

		data := []memberModel.DataMemberAdditionalInfo{}
		for _, s := range result.Entries[0].Attributes {
			if s.Name == "objectGUID" || s.Name == "objectSid" {
				s.Values[0] = base64.StdEncoding.EncodeToString([]byte(s.Values[0]))
			}
			attributes := memberModel.DataMemberAdditionalInfo{
				Key:   s.Name,
				Value: s.Values[0],
			}
			data = append(data, attributes)
		}

		savedData := memberModel.MemberAdditionalInfo{
			MemberID: member.ID,
			AuthType: authType,
			Data:     data,
		}

		findMemberAdditional := <-l.MemberAdditionalRepo.Load(ctxReq, member.ID, authType)
		if findMemberAdditional.Error != nil {
			save := <-l.MemberAdditionalRepo.Save(ctxReq, &savedData)
			if save.Error != nil {
				tags[helper.TextResponse] = save.Error
				return
			}
			tags[helper.TextResponse] = savedData
			return
		}

		resultMemberAdditional, _ := findMemberAdditional.Result.(memberModel.MemberAdditionalInfo)
		savedData.ID = resultMemberAdditional.ID
		update := <-l.MemberAdditionalRepo.Update(ctxReq, &savedData)
		if update.Error != nil {
			tags[helper.TextResponse] = update.Error
			return
		}

		tags[helper.TextResponse] = savedData
	})
}
