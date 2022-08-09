package usecase

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
)

// createMemberFromSocMed function for create member
func (au *AuthUseCaseImpl) createMemberFromSocMed(ctxReq context.Context, Socmed interface{}, grantType string) <-chan ResultUseCase {
	ctx := "AuthUseCase-createMemberFromSocMed"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {

		memberID := helper.GenerateMemberIDv2()

		// set default gender as male
		gender := memberModel.MaleString

		// set the member status directly to active
		dataMember := memberModel.Member{
			ID:     memberID,
			Gender: memberModel.StringToGender(gender),
			Status: memberModel.StringToStatus(memberModel.ActiveString),
		}

		// generate data member from media social data
		mediaData := <-au.generateSocmedData(ctxReq, dataMember, Socmed, grantType, true)

		if mediaData.Error != nil {
			output <- ResultUseCase{HTTPStatus: mediaData.HTTPStatus, Error: mediaData.Error}
			return
		}

		dataMember, ok := mediaData.Result.(memberModel.Member)
		if !ok {
			err := errors.New("result is not member data")
			output <- ResultUseCase{HTTPStatus: http.StatusInternalServerError, Error: err}
			return
		}

		saveResult := <-au.MemberRepoWrite.Save(ctxReq, dataMember)
		if saveResult.Error != nil {
			err := errors.New("failed to save member")
			helper.SendErrorLog(ctxReq, ctx, "save_member", saveResult.Error, dataMember)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusInternalServerError}
			return
		}

		tags["args"] = dataMember
		output <- ResultUseCase{Result: dataMember}
	})
	return output
}

func (au *AuthUseCaseImpl) publishMemberData(ctxReq context.Context, dataMember memberModel.Member, eventType string) (payload serviceModel.MemberDolphin, err error) {
	ctx := "AuthUseCaseImpl-publishMemberData"
	defer func(err error) {
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, "publish_member_kafka", err, dataMember)
		}
	}(err)

	dolphinData := serviceModel.MemberDolphin{
		ID:         dataMember.ID,
		Email:      strings.ToLower(dataMember.Email),
		FirstName:  dataMember.FirstName,
		LastName:   dataMember.LastName,
		Gender:     dataMember.Gender.GetDolpinGender(),
		DOB:        dataMember.BirthDateString,
		Mobile:     dataMember.Mobile,
		Status:     strings.ToUpper(memberModel.ActiveString),
		AzureID:    dataMember.SocialMedia.AzureID,
		AppleID:    dataMember.SocialMedia.AppleID,
		FacebookID: dataMember.SocialMedia.FacebookID,
		GoogleID:   dataMember.SocialMedia.GoogleID,
		LDAPID:     dataMember.SocialMedia.LDAPID,
		Created:    time.Now().Format(time.RFC3339),
	}

	if err = au.PublishToKafkaDolphin(ctxReq, dolphinData, eventType); err != nil {
		return dolphinData, err
	}
	plLog := memberModel.MemberLog{
		Before: &memberModel.Member{},
		After:  &dataMember,
	}
	go au.QPublisher.QueueJob(ctxReq, plLog, dataMember.ID, "InsertLogAuth")
	return dolphinData, nil
}

// CheckMemberSocmedType function for validate request token with grantType is azure
func (au *AuthUseCaseImpl) CheckMemberSocmedType(ctxReq context.Context, data *model.RequestToken, socialMedia interface{}, Email string) *model.ValidateSocmedRequest {
	ctx := "AuthUseCase-CheckMemberSocmedType"
	newMember := false
	eventType := "update"
	trace := tracer.StartTrace(ctxReq, ctx)
	defer trace.Finish(nil)

	memberResult := <-au.MemberQueryRead.FindByEmail(ctxReq, Email)

	// when email does not exist then save the detail email
	if memberResult.Error != nil && memberResult.Error == sql.ErrNoRows {
		// if request coming from version 3, only return availability
		if data.Version == helper.Version3 && data.GrantType != model.AuthTypeApple {
			return &model.ValidateSocmedRequest{HTTPStatus: http.StatusForbidden}
		}
		eventType = "register"
		// create member
		createMember := <-au.createMemberFromSocMed(ctxReq, socialMedia, data.GrantType)
		if createMember.Error != nil {
			return &model.ValidateSocmedRequest{HTTPStatus: createMember.HTTPStatus, Error: createMember.Error}
		}

		member, ok := createMember.Result.(memberModel.Member)
		if !ok {
			return &model.ValidateSocmedRequest{HTTPStatus: http.StatusInternalServerError, Error: errors.New(msgResultNotMember)}
		}
		data.NewMember = true

		au.publishMemberData(ctxReq, member, eventType)

		return &model.ValidateSocmedRequest{Data: &member, NewMember: true, HTTPStatus: 200, Error: nil}
	}

	member, ok := memberResult.Result.(memberModel.Member)
	if !ok {
		return &model.ValidateSocmedRequest{HTTPStatus: http.StatusInternalServerError, Error: errors.New(msgResultNotMember)}
	}

	// check member status
	checkMember := <-au.checkMemberStatus(ctxReq, member, socialMedia, data.GrantType)
	if checkMember.Error != nil {
		return &model.ValidateSocmedRequest{HTTPStatus: checkMember.HTTPStatus, Error: checkMember.Error}
	}

	memberData, ok := checkMember.Result.(memberModel.Member)
	if !ok {
		return &model.ValidateSocmedRequest{HTTPStatus: http.StatusInternalServerError, Error: errors.New(msgResultNotMember)}
	}

	memberData.MFAEnabled = member.MFAEnabled
	memberData.AdminMFAEnabled = member.AdminMFAEnabled

	if len(memberData.FirstName) == 0 {
		newMember = true
	}
	if err := au.publishUpdateMemberData(ctxReq, data.GrantType, member, memberData); err != nil {
		return &model.ValidateSocmedRequest{HTTPStatus: http.StatusInternalServerError, Error: err}
	}

	return &model.ValidateSocmedRequest{Data: &memberData, NewMember: newMember, HTTPStatus: 200, Error: nil}
}

func (au *AuthUseCaseImpl) publishUpdateMemberData(ctxReq context.Context, grandType string, oldMember memberModel.Member, newMemberData memberModel.Member) error {
	if (golib.StringInSlice(grandType, []string{model.AuthTypeGoogle, model.AuthTypeGoogleBackend}) && oldMember.SocialMedia.GoogleID == "") ||
		(grandType == model.AuthTypeFacebook && oldMember.SocialMedia.FacebookID == "") ||
		(grandType == model.AuthTypeApple && oldMember.SocialMedia.AppleID == "") {
		if _, err := au.publishMemberData(ctxReq, newMemberData, "update"); err != nil {
			return err
		}
	}
	return nil
}

// generateSocmedData function for generate data social media
func (au *AuthUseCaseImpl) generateSocmedData(ctxReq context.Context, member memberModel.Member, Socmed interface{}, grantType string, new bool) <-chan ResultUseCase {
	ctx := "AuthUseCase-generateSocmedData"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		if err := au.processSocialMediaData(grantType, Socmed, &member, new); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "generate_sosmed", err, member)
			output <- ResultUseCase{HTTPStatus: http.StatusInternalServerError, Error: err}
			return
		}

		if (member.Email == "") && grantType != model.AuthTypeApple {
			err := errors.New("failed to fetch social media data")
			tags["response"] = member
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		member.Email = strings.ToLower(member.Email)
		tags["args"] = member
		output <- ResultUseCase{Result: member}
	})

	return output
}

func (au *AuthUseCaseImpl) processSocialMediaData(grantType string, socialMedia interface{}, existingMember *memberModel.Member, new bool) error {
	ctx := "processSocialMedia-AuthUseCaseImpl"
	ctxReq := context.Background()
	switch grantType {
	case model.AuthTypeAzure:
		if err := au.parseAzureData(socialMedia, existingMember, new); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "parse_azure_data", err, existingMember)
			return err
		}
	case model.AuthTypeFacebook:
		if err := au.parseFacebookData(socialMedia, existingMember, new); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "parse_facebook_data", err, existingMember)
			return err
		}

	case model.AuthTypeGoogle, model.AuthTypeGoogleBackend:
		if err := au.parseGoogleData(socialMedia, existingMember, new); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "parse_google", err, existingMember)
			return err
		}

	case model.AuthTypeGoogleOAauth:
		if err := au.parseGoogleOAuthData(socialMedia, existingMember, new); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "parse_google_oauth", err, existingMember)
			return err
		}

	case model.AuthTypeLDAP:
		if err := au.parseLDAPData(socialMedia, existingMember, new); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "parse_ldap", err, existingMember)
			return err
		}

	case model.AuthTypeApple:
		if err := au.parseAppleData(socialMedia, existingMember, new); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "parse_apple_data", err, existingMember)
			return err
		}
	case model.AuthTypePassword:
		if err := au.parsePasswordMicrosite(socialMedia, existingMember, new); err != nil {
			helper.SendErrorLog(ctxReq, ctx, "parse_microsite_data", err, existingMember)
			return err
		}

	default:
		return errors.New("invalid grant type")
	}
	return nil
}
func (au *AuthUseCaseImpl) parsePasswordMicrosite(input interface{}, member *memberModel.Member, new bool) error {
	media, ok := input.(model.MicrositeClient)
	if !ok {
		return errors.New("result is not microsite data")
	}
	if member.FirstName == "" {
		member.FirstName = media.Firstname
	}

	member.Email = media.Email
	member.Status = memberModel.StringToStatus(memberModel.ActiveString)
	member.SignUpFrom = signUpFromLKPP
	member.NewMember = true
	return nil
}

func (au *AuthUseCaseImpl) parseAppleData(input interface{}, member *memberModel.Member, new bool) error {
	media, ok := input.(model.AppleProfile)
	if !ok {
		return errors.New("result is not apple data")
	}
	member.SocialMedia.AppleID = media.Sub
	if new || member.SocialMedia.AppleConnect.IsZero() {
		member.SocialMedia.AppleConnect = time.Now()
	}
	member.Status = memberModel.StringToStatus(memberModel.ActiveString)
	isStaff := media.IsBhinnekaEmail()
	member.IsStaff = isStaff
	member.Email = media.Email
	if new && len(media.FirstName) != 0 {
		member.FirstName = media.FirstName
		member.LastName = media.LastName
	}
	return nil
}
func (au *AuthUseCaseImpl) parseLDAPData(input interface{}, member *memberModel.Member, new bool) error {
	media, ok := input.(*model.LDAPProfile)
	if !ok {
		return errors.New("result is not ldap data")
	}
	member.SocialMedia.LDAPID = media.ObjectID
	member.Status = memberModel.StringToStatus(memberModel.ActiveString)
	member.IsStaff = true
	member.Email = media.Email
	member.JobTitle = media.JobTitle
	member.Department = media.Department

	// set isAdmin = true by default login ldap
	member.IsAdmin = true

	names := strings.Split(media.DisplayName, " ")

	member.FirstName = names[0]
	member.LastName = helper.SetLastName(names)
	return nil
}

func (au *AuthUseCaseImpl) parseGoogleData(input interface{}, existingMember *memberModel.Member, new bool) error {
	media, ok := input.(model.GoogleOAuth2Response)
	if !ok {
		return errors.New("result is not google data")
	}
	existingMember.SocialMedia.GoogleID = media.ID

	if new || existingMember.SocialMedia.GoogleConnect.IsZero() {
		existingMember.SocialMedia.GoogleConnect = time.Now()
		if existingMember.FirstName == "" || existingMember.LastName == "" {
			names := strings.Split(media.Name, " ")
			existingMember.FirstName = names[0]
			existingMember.LastName = helper.SetLastName(names)
		}
	}
	existingMember.Status = memberModel.StringToStatus(memberModel.ActiveString)
	existingMember.IsStaff = media.IsBhinnekaEmail()
	existingMember.Email = media.Email
	return nil
}

func (au *AuthUseCaseImpl) parseGoogleOAuthData(input interface{}, existingMember *memberModel.Member, new bool) error {
	media, ok := input.(model.GoogleOAuthToken)
	if !ok {
		return errors.New("result is not google data")
	}
	existingMember.SocialMedia.GoogleID = media.Sub

	if new || existingMember.SocialMedia.GoogleConnect.IsZero() {
		existingMember.SocialMedia.GoogleConnect = time.Now()
		if existingMember.FirstName == "" || existingMember.LastName == "" {
			names := strings.Split(media.Name, " ")
			existingMember.FirstName = names[0]
			existingMember.LastName = helper.SetLastName(names)
		}
	}
	existingMember.Status = memberModel.StringToStatus(memberModel.ActiveString)
	existingMember.IsStaff = media.IsBhinnekaEmail()
	existingMember.Email = media.Email
	return nil
}

func (au *AuthUseCaseImpl) parseFacebookData(input interface{}, existingMember *memberModel.Member, new bool) error {
	socialMedia, ok := input.(model.FacebookResponse)
	if !ok {
		return errors.New("result is not facebook data")
	}
	existingMember.SocialMedia.FacebookID = socialMedia.ID

	if new || existingMember.SocialMedia.FacebookConnect.IsZero() {
		existingMember.SocialMedia.FacebookConnect = time.Now()
		if existingMember.FirstName == "" || existingMember.LastName == "" {
			names := strings.Split(socialMedia.Name, " ")
			existingMember.FirstName = names[0]
			existingMember.LastName = helper.SetLastName(names)
		}
	}

	existingMember.Status = memberModel.StringToStatus(memberModel.ActiveString)
	existingMember.IsStaff = false
	existingMember.Email = socialMedia.Email
	return nil
}

func (au *AuthUseCaseImpl) parseAzureData(input interface{}, member *memberModel.Member, new bool) error {
	media, ok := input.(model.AzureResponse)
	if !ok {
		return errors.New("result is not azure data")
	}
	member.SocialMedia.AzureID = media.ObjectID
	isStaff := media.IsBhinnekaEmail()
	member.IsStaff = isStaff
	member.Email = media.Email
	member.JobTitle = media.JobTitle
	member.Department = media.Department
	names := strings.Split(media.DisplayName, " ")

	member.FirstName = names[0]
	member.LastName = helper.SetLastName(names)
	return nil
}
