package usecase

import (
	"context"
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Bhinneka/golib"
	goString "github.com/Bhinneka/golib/string"
	"github.com/Bhinneka/golib/tracer"
	localConfig "github.com/Bhinneka/user-service/config"
	"github.com/Bhinneka/user-service/helper"
	authModel "github.com/Bhinneka/user-service/src/auth/v1/model"
	authRepo "github.com/Bhinneka/user-service/src/auth/v1/repo"
	"github.com/Bhinneka/user-service/src/auth/v1/token"
	authUsecase "github.com/Bhinneka/user-service/src/auth/v1/usecase"
	corporateQuery "github.com/Bhinneka/user-service/src/corporate/v2/query"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/member/v1/query"
	"github.com/Bhinneka/user-service/src/member/v1/repo"
	merchantModel "github.com/Bhinneka/user-service/src/merchant/v2/model"
	merchantRepoRead "github.com/Bhinneka/user-service/src/merchant/v2/repo"
	service "github.com/Bhinneka/user-service/src/service"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	sessionModel "github.com/Bhinneka/user-service/src/session/v1/model"
	sessionQuery "github.com/Bhinneka/user-service/src/session/v1/query"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	shippingRepo "github.com/Bhinneka/user-service/src/shipping_address/v2/repo"
	"github.com/golang-jwt/jwt"
)

// MemberUseCaseImpl data structure
type MemberUseCaseImpl struct {
	MemberRepoRead                    repo.MemberRepository
	MemberRepoWrite                   repo.MemberRepository
	MemberMFARepoWrite                repo.MemberMFARepository
	MemberRepoRedis                   repo.MemberRepositoryRedis
	TokenActivationRepo               repo.TokenActivationRepository
	ShippingAddressRepo               shippingRepo.ShippingAddressRepository
	LoginAttemptRepo                  authRepo.AttemptRepository
	LoginSessionRedis                 authRepo.LoginSessionRepository
	MemberQueryRead                   query.MemberQuery
	MemberQueryWrite                  query.MemberQuery
	MemberMFAQueryRead                query.MemberMFAQuery
	SessionQueryRead                  sessionQuery.SessionInfoQuery
	StaticService                     service.StaticServices
	UploadService                     service.UploadServices
	ActivityService                   service.ActivityServices
	QPublisher                        service.QPublisher
	Hash                              model.PasswordHasher
	TokenActivationExpiration         time.Duration
	ResendActivationAttemptAge        string
	ResendActivationAttemptAgeRequest string
	Topic                             string
	IsProductionStage                 bool
	SturgeonCFUrl                     string
	B2cCFUrl                          string
	AccessTokenGenerator              token.AccessTokenGenerator
	NotificationService               service.NotificationServices
	SendbirdService                   service.SendbirdServices
	AuthUseCase                       authUsecase.AuthUseCase
	CorporateContactQueryRead         corporateQuery.ContactQuery
	CorporateAccContactQueryRead      corporateQuery.AccountContactQuery
	MerchantRepoRead                  merchantRepoRead.MerchantRepository
	MerchantEmployeeRead              merchantRepoRead.MerchantEmployeeRepository
}

// NewMemberUseCase function for initialise member use case implementation
func NewMemberUseCase(
	repository localConfig.ServiceRepository,
	query localConfig.ServiceQuery,
	services localConfig.ServiceShared,
	params localConfig.MembershipParameters,
	authUsecase authUsecase.AuthUseCase,
) MemberUseCase {

	return &MemberUseCaseImpl{
		MemberRepoRead:                    repository.MemberRepository,
		MemberRepoWrite:                   repository.MemberRepository,
		MemberMFARepoWrite:                repository.MemberMFARepository,
		MemberRepoRedis:                   repository.MemberRedisRepository,
		TokenActivationRepo:               repository.TokenActivationRepoRedis,
		LoginAttemptRepo:                  repository.AttemptRepositoryRedis,
		LoginSessionRedis:                 repository.LoginSessionRepositoryRedis,
		ShippingAddressRepo:               repository.ShippingAddressRepository,
		MemberQueryRead:                   query.MemberQueryRead,
		MemberQueryWrite:                  query.MemberQueryWrite,
		MemberMFAQueryRead:                query.MemberMFAQueryRead,
		SessionQueryRead:                  query.SessionInfoQuery,
		CorporateContactQueryRead:         query.CorporateContactQueryRead,
		CorporateAccContactQueryRead:      query.CorporateAccContactQueryRead,
		QPublisher:                        services.QPublisher,
		UploadService:                     services.UploadService,
		ActivityService:                   services.ActivityService,
		Hash:                              params.Hash,
		TokenActivationExpiration:         params.TokenActivationExpiration,
		Topic:                             params.Topic,
		IsProductionStage:                 params.IsProductionStage,
		SturgeonCFUrl:                     params.SturgeonCFUrl,
		B2cCFUrl:                          params.B2cCFUrl,
		ResendActivationAttemptAge:        params.ResendActivationAttemptAge,
		ResendActivationAttemptAgeRequest: params.ResendActivationAttemptAgeRequest,
		AccessTokenGenerator:              params.AccessTokenGenerator,
		NotificationService:               services.NotificationService,
		SendbirdService:                   services.SendbirdService,
		AuthUseCase:                       authUsecase,
		MerchantRepoRead:                  repository.MerchantRepository,
		MerchantEmployeeRead:              repository.MerchantEmployeeRepository,
	}
}

func cleanTags(_ context.Context, data *model.Member) (maskedData *model.Member) {
	dataMask := new(model.Member)
	helper.CloneStruct(data, dataMask)

	dataMask.Password = ""
	dataMask.RePassword = ""
	return dataMask
}

// CheckEmailAndMobileExistence function for checking member existence by email and mobile
func (mu *MemberUseCaseImpl) CheckEmailAndMobileExistence(ctxReq context.Context, data *model.Member) <-chan ResultUseCase {
	ctx := "MemberUseCase-CheckEmailAndMobileExistence"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags[helper.TextArgs] = cleanTags(ctxReq, data)
		if err := mu.ValidateEmailAndPhone(ctxReq, data); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
	})

	return output
}

// BulkValidateEmailAndPhone return existing member
func (mu *MemberUseCaseImpl) BulkValidateEmailAndPhone(ctxReq context.Context, members []*model.Member) []string {
	tc := tracer.StartTrace(ctxReq, "MemberUseCase-BulkValidateEmailAndPhone")
	tags := map[string]interface{}{}
	defer tc.Finish(tags)
	var emails []string
	var responses []string
	for _, e := range members {
		if err := goString.ValidateEmail(e.Email); err != nil {
			responses = append(responses, fmt.Sprintf("email %s is invalid", e.Email))
			return responses
		}
		emails = append(emails, strings.ToLower(e.Email))
	}
	tags["emails"] = emails
	emailResults := <-mu.MemberQueryRead.BulkFindByEmail(ctxReq, emails)
	if emailResults.Error != nil {
		responses = append(responses, "please try again")
		return responses
	}
	if emailResults.Result != nil {
		members := emailResults.Result.([]model.Member)
		for _, member := range members {
			responses = append(responses, fmt.Sprintf("email %s already exists", member.Email))
		}
	}
	tags["response"] = responses
	return responses
}

// ValidateEmailAndPhone when import
func (mu *MemberUseCaseImpl) ValidateEmailAndPhone(ctxReq context.Context, data *model.Member) error {
	if err := goString.ValidateEmail(data.Email); err != nil {
		return err
	}

	emailResult := <-mu.MemberQueryRead.FindByEmail(ctxReq, data.Email)
	// error nil means email already exists
	if emailResult.Result != nil {
		member := emailResult.Result.(model.Member)
		if !member.IsBhinnekaEmail() || data.Type == "import" {
			return fmt.Errorf("email %s already exists", member.Email)
		}
		return fmt.Errorf("Alamat email sudah terdaftar, silakan masuk dengan pilih Google")
	}

	if data.Mobile == "" {
		return fmt.Errorf("mobile required for %s ", data.Email)
	}

	if os.Getenv("ENABLE_VALIDATE_MOBILE") == "true" {
		mobileResult := <-mu.MemberQueryRead.FindByMobile(ctxReq, data.Mobile)
		// error nil means mobile already exists
		if mobileResult.Result != nil {
			return fmt.Errorf("mobile number %s already exists", data.Mobile)
		}
	}
	return nil
}

// GetDetailMemberByEmail function for getting detail member based on email
func (mu *MemberUseCaseImpl) GetDetailMemberByEmail(email string) <-chan ResultUseCase {
	output := make(chan ResultUseCase)
	go func() {
		defer close(output)

		if err := goString.ValidateEmail(email); err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		memberResult := <-mu.MemberQueryRead.FindByEmail(context.Background(), email)
		if memberResult.Error != nil {
			// when data not found
			if memberResult.Error == sql.ErrNoRows {
				memberResult.Error = fmt.Errorf(helper.ErrorDataNotFound, msgErrorMember+email)
			}

			output <- ResultUseCase{Error: memberResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		member, ok := memberResult.Result.(model.Member)
		if !ok {
			err := errors.New(msgErrorResultMember)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// adjust data member then return
		output <- ResultUseCase{Result: mu.adjustMemberData(context.Background(), member)}
	}()

	return output
}

// GetDetailMemberByID function for getting detail member based on member id
func (mu *MemberUseCaseImpl) GetDetailMemberByID(ctxReq context.Context, uid string) <-chan ResultUseCase {
	ctx := "MemberUseCase-GetDetailMemberByID"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags[helper.TextMemberIDCamel] = uid
		if !strings.Contains(uid, usrFormat) {
			err := fmt.Errorf(helper.ErrorParameterInvalid, msgErrorMemberID)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		memberResult := <-mu.MemberRepoRead.Load(ctxReq, uid)
		if memberResult.Error != nil {
			if memberResult.Error == sql.ErrNoRows {
				memberResult.Error = fmt.Errorf(helper.ErrorDataNotFound, labelMember)
			}
			tags[helper.TextResponse] = memberResult.Error
			output <- ResultUseCase{Error: memberResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		member, ok := memberResult.Result.(model.Member)
		if !ok {
			err := errors.New(msgErrorResultMember)
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// adjust data member then return
		output <- ResultUseCase{Result: mu.adjustMemberData(ctxReq, member)}
	})

	return output
}

// UpdateDetailMemberByID function for validating and updating member detail
func (mu *MemberUseCaseImpl) UpdateDetailMemberByID(ctxReq context.Context, data model.Member) <-chan ResultUseCase {
	ctx := "MemberUseCase-UpdateDetailMemberByID"
	var params serviceModel.SendbirdRequestV4
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags[helper.TextMemberIDCamel] = data.ID
		tags[helper.TextParameter] = data
		// validate data first
		if err := mu.validateMemberData(ctxReq, &data); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		memberResult := <-mu.MemberRepoRead.Load(ctxReq, data.ID)
		if memberResult.Error != nil {
			if memberResult.Error == sql.ErrNoRows {
				memberResult.Error = fmt.Errorf(helper.ErrorDataNotFound, labelMember)
			}

			tags[helper.TextResponse] = memberResult.Error
			output <- ResultUseCase{Error: memberResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		member, ok := memberResult.Result.(model.Member)
		if !ok {
			err := errors.New(msgErrorResultMember)
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		tags[helper.TextEmail] = member.Email
		member.StatusString = member.Status.String()
		oldMember := member
		data.IsActive, _ = strconv.ParseBool(data.IsActiveString)
		data = checkIsActiveStatus(data)

		mu.deleteSessionByStatus(ctxReq, data, member)

		member, httpStatus, err := mu.adjustUpdateMember(ctxReq, member, data)
		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: httpStatus}
			return
		}

		tags[helper.TextArgs] = member
		saveResult := <-mu.MemberRepoWrite.Save(ctxReq, member)
		if saveResult.Error != nil {
			err := errors.New(msgErrorSaveMember)
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusInternalServerError}
			return
		}

		// publish to kafka
		member.StatusString = strings.ToUpper(data.StatusString)
		go mu.PublishToKafkaUser(ctxReq, &member, textUpdate)
		mu.deleteSessionMember(ctxReq, member)

		// send to audit trail activity service
		member.ModifiedBy = data.ModifiedBy
		plLog := model.MemberLog{
			Before: &oldMember,
			After:  &member,
		}

		params.UserID = member.ID
		params.ProfileURL = member.ProfilePicture
		params.NickName = member.FirstName + " " + member.LastName
		params.Token = member.Token
		mu.UpdateUserSendbirdV4(ctxReq, &params)

		go mu.QPublisher.QueueJob(ctxReq, plLog, data.ID, "InsertLogUpdateMember")

		output <- ResultUseCase{Result: mu.adjustMemberData(ctxReq, member)}
	})

	return output
}

func (mu *MemberUseCaseImpl) adjustUpdateMember(ctxReq context.Context, member, data model.Member) (model.Member, int, error) {
	// check the existence mobile phone data
	if member.Mobile != data.Mobile && os.Getenv("ENABLE_VALIDATE_MOBILE") == "true" {
		mobileResult := <-mu.MemberQueryRead.FindByMobile(ctxReq, data.Mobile)
		// error nil means email already exists
		if mobileResult.Error == nil {
			err := errors.New("mobile number already exists")
			return member, http.StatusBadRequest, err
		}
	}
	oldmember := member
	// replace old value
	member.FirstName = data.FirstName
	member.LastName = data.LastName
	member.Gender = data.Gender
	member.BirthDate = data.BirthDate
	member.BirthDateString = data.BirthDateString
	member.Phone = data.Phone
	member.Ext = data.Ext
	member.Mobile = data.Mobile
	member.Status = data.Status
	member.StatusString = data.StatusString
	member.IsActive = data.IsActive
	member.IsActiveString = data.IsActiveString
	member = checkIsActiveStatusV2(member, oldmember)
	member.Status, _ = model.ValidateStatus(member.StatusString)

	if data.RequestFrom != model.Sturgeon {
		member.Address = data.Address
	}

	// when new password exists
	if len(data.NewPassword) > 0 {
		// encode the new password then replace the old password and salt
		member.Salt = mu.Hash.GenerateSalt()
		err := mu.Hash.ParseSalt(member.Salt)
		if err != nil {
			return member, http.StatusInternalServerError, err
		}

		member.Password = base64.StdEncoding.EncodeToString(mu.Hash.Hash([]byte(data.NewPassword)))
		member.LastPasswordModified = time.Now()
	}

	if len(data.IsStaffString) > 0 {
		member.IsStaff = data.IsStaff
	}
	if len(data.IsAdminString) > 0 {
		member.IsAdmin = data.IsAdmin
	}
	if len(member.Password) > 0 {
		member.HasPassword = true
	}

	return member, http.StatusOK, nil
}

func (mu *MemberUseCaseImpl) deleteSessionByStatus(ctxReq context.Context, data model.Member, member model.Member) error {
	// compare old data with the new one and process the following steps
	if member.StatusString == model.BlockedString && (data.StatusString == model.ActiveString || data.StatusString == model.InactiveString || data.StatusString == model.NewString) {
		// delete the redis data of login attempt
		key := fmt.Sprintf("ATTEMPT:%s", member.Email)

		delResult := <-mu.LoginAttemptRepo.Delete(ctxReq, key)
		if delResult.Error != nil {
			return delResult.Error
		}
	}
	return nil
}
func checkIsActiveStatus(input model.Member) model.Member {
	if input.StatusString == "" {
		switch input.IsActive {
		case true:
			input.StatusString = model.ActiveString
		default:
			input.StatusString = model.NewString
		}
	} else {
		switch input.StatusString {
		case model.ActiveString:
			input.IsActive = true
			input.IsActiveString = "true"
		default:
			input.IsActive = false
			input.IsActiveString = "false"
		}
	}
	return input
}
func checkIsActiveStatusV2(member model.Member, oldData model.Member) model.Member {
	if member.StatusString == model.NewString {
		switch member.IsActive {
		case true:
			member.StatusString = model.ActiveString
		default:
			member.StatusString = oldData.StatusString
			if member.StatusString == model.ActiveString {
				member.StatusString = model.InactiveString
			}
		}
	} else {
		switch member.StatusString {
		case model.ActiveString:
			member.IsActive = true
		default:
			member.IsActive = false
		}
	}
	return member
}

func (mu *MemberUseCaseImpl) deleteSessionMember(ctxReq context.Context, member model.Member) error {
	if member.StatusString == model.BlockedString || member.StatusString == model.InactiveString || member.StatusString == model.NewString {
		session := <-mu.SessionQueryRead.GetListSessionInfo(ctxReq, &sessionModel.ParamList{Page: 1, Limit: 100, Sort: "desc", OrderBy: "createdAt", Email: member.Email})
		if session.Error == nil {
			if res, ok := session.Result.(sessionModel.SessionInfoList); ok {
				for _, s := range res.Data {
					key := strings.Join([]string{"STG", *s.UserID, *s.DeviceID, *s.ClientType}, "-")
					mu.LoginSessionRedis.Delete(ctxReq, key)
				}
			}
		}
	}
	return nil
}

// RegisterMember function for registering new member
func (mu *MemberUseCaseImpl) RegisterMember(ctxReq context.Context, data *model.Member) <-chan ResultUseCase {
	ctx := "MemberUseCase-RegisterMember"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		data.Email = strings.ToLower(data.Email)
		tags[helper.TextEmail] = data.Email

		// validate data first
		if err := mu.validateMemberData(ctxReq, data); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		data, httpStatus, err := mu.adjustRegistrationData(ctxReq, data)
		if err != nil {
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: err, HTTPStatus: httpStatus}
			return
		}
		tags[helper.TextArgs] = data
		saveResult := <-mu.MemberRepoWrite.Save(ctxReq, *data)
		if saveResult.Error != nil {
			err := errors.New(msgErrorSaveMember)
			tracer.SetError(ctxReq, saveResult.Error)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusInternalServerError}
			return
		}

		// check sign up from process
		eventType := textRegister
		if data.SignUpFrom == model.Dolphin {
			eventType = textUpdate
		}

		data.StatusString = strings.ToUpper(model.NewString)
		data.IsActive = false
		data.Created = time.Now()
		go mu.PublishToKafkaUser(ctxReq, data, eventType)

		response := model.SuccessResponse{
			ID:          data.ID,
			Message:     helper.SuccessMessage,
			Token:       data.Token,
			HasPassword: data.HasPassword,
			FirstName:   data.FirstName,
			LastName:    data.LastName,
			Email:       strings.ToLower(data.Email),
		}
		plLog := model.MemberLog{
			Before: &model.Member{},
			After:  data,
		}
		plEmail := model.MemberEmailQueue{
			Member: data,
			Data:   response,
		}

		mu.conditionalEmail(ctxReq, data, plEmail)

		// send to audit trail activity service
		go mu.QPublisher.QueueJob(ctxReq, plLog, data.ID, "InsertLogRegisterMember")

		// return the token to be sent to email notification service
		output <- ResultUseCase{Result: response}
	})

	return output
}

func (mu *MemberUseCaseImpl) adjustRegistrationV3(ctxReq context.Context, data *model.Member) error {

	if data.SignUpFrom != model.Dolphin {
		data.Status = model.StringToStatus(model.NewString)
	}
	// check key redis for temporary
	socialMediaID := ""
	if data.SocialMedia.GoogleID != "" {
		socialMediaID = data.SocialMedia.GoogleID
		data.SocialMedia.GoogleConnect = time.Now()
	} else if data.SocialMedia.FacebookID != "" {
		socialMediaID = data.SocialMedia.FacebookID
		data.SocialMedia.FacebookConnect = time.Now()
	} else if data.SocialMedia.AppleID != "" {
		socialMediaID = data.SocialMedia.AppleID
		data.SocialMedia.AppleConnect = time.Now()
	}
	// to accomodate regular rgister, need to check social media ID existence
	if socialMediaID != "" {
		data.Status = model.StringToStatus(model.ActiveString)
		key := fmt.Sprintf("STG:%s:%s", data.Email, socialMediaID)
		tempKey := <-mu.LoginSessionRedis.Load(ctxReq, key)
		if tempKey.Error != nil {
			return errors.New("invalid session ID")
		}
	}
	return nil
}

func (mu *MemberUseCaseImpl) adjustRegistrationData(ctxReq context.Context, data *model.Member) (*model.Member, int, error) {
	data.Status = model.StringToStatus(model.NewString)
	if data.APIVersion == helper.Version3 {
		if err := mu.adjustRegistrationV3(ctxReq, data); err != nil {
			return data, http.StatusBadRequest, err
		}
	}

	// generate member id
	if len(data.ID) <= 0 {
		data.ID = helper.GenerateMemberIDv2()
	}

	data.HasPassword = false
	// optional when password exists only
	if len(data.NewPassword) > 0 {
		// encode the new password then replace the old password and salt
		passwordHasher := model.NewPBKDF2Hasher(model.SaltSize, model.SaltSize, model.IterationsCount, sha1.New)
		data.Salt = passwordHasher.GenerateSalt()
		err := passwordHasher.ParseSalt(data.Salt)
		if err != nil {
			return data, http.StatusInternalServerError, err
		}

		data.Password = base64.StdEncoding.EncodeToString(passwordHasher.Hash([]byte(data.NewPassword)))
		// set only if validation is success
		data.RePassword = ""
		data.HasPassword = true
	}

	// generate random string
	mix := data.Email + "-" + data.FirstName
	data.Token = helper.GenerateTokenByString(mix)

	return data, http.StatusOK, nil
}

// ActivateMember function for activating inactive member
func (mu *MemberUseCaseImpl) ActivateMember(ctxReq context.Context, token, requestFrom string) <-chan ResultUseCase {
	ctx := "MemberUseCase-ActivateMember"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags[helper.TextParameter] = token
		if len(token) == 0 {
			err := fmt.Errorf(helper.ErrorParameterRequired, "token")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		memberResult := <-mu.MemberQueryRead.FindByToken(ctxReq, token)
		if memberResult.Error != nil {
			// when data not found
			if memberResult.Error == sql.ErrNoRows {
				memberResult.Error = fmt.Errorf(helper.ErrorDataNotFound, "member with token "+token)
			}

			tags[helper.TextResponse] = memberResult.Error
			output <- ResultUseCase{Error: memberResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		member, ok := memberResult.Result.(model.Member)
		if !ok {
			err := errors.New(msgErrorResultMember)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		oldMember := member

		// append the new data
		member.Status = model.StringToStatus(model.ActiveString)
		member.IsActive = true
		member.Token = ""

		saveResult := <-mu.MemberRepoWrite.Save(ctxReq, member)
		if saveResult.Error != nil {
			tracer.SetError(ctxReq, saveResult.Error)
			err := errors.New(msgErrorSaveMember)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// for has password field
		hasPassword := len(member.Password) > 0

		// check sign up from process
		if member.SignUpFrom != model.Dolphin {
			// publish to kafka
			member.IsActiveString = "true"
			member.StatusString = strings.ToUpper(model.ActiveString)
			member.Created = time.Now()
			go mu.PublishToKafkaUser(ctxReq, &member, "activation")
		}

		response := model.SuccessResponse{
			ID:          golib.RandomString(8),
			Message:     helper.SuccessMessage,
			HasPassword: hasPassword,
			Email:       strings.ToLower(member.Email),
			FirstName:   member.FirstName,
			LastName:    member.LastName,
		}

		tags[helper.TextResponse] = response

		// send to audit trail activity service
		plLog := model.MemberLog{
			Before: &oldMember,
			After:  &member,
		}

		if requestFrom == model.Sturgeon {
			plEmail := model.MemberEmailQueue{
				Member: &member,
				Data:   response,
			}
			go mu.QPublisher.QueueJob(ctxReq, plEmail, member.ID, "SendEmailWelcomeMember")
		}

		go mu.QPublisher.QueueJob(ctxReq, plLog, member.ID, "InsertLogUpdateMember")
		output <- ResultUseCase{Result: response}
	})

	return output
}

// GetListMembers function for getting list of members
func (mu *MemberUseCaseImpl) GetListMembers(ctxReq context.Context, params *model.Parameters) <-chan ResultUseCase {
	output := make(chan ResultUseCase)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				err := fmt.Errorf("%v", r)
				output <- ResultUseCase{HTTPStatus: http.StatusInternalServerError, Error: err}
			}
			close(output)
		}()

		var err error

		// validate all parameters
		paging, err := helper.ValidatePagination(
			helper.PaginationParameters{
				Page:     1, // default
				StrPage:  params.StrPage,
				Limit:    10, // default
				StrLimit: params.StrLimit,
			})

		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		params.Page = paging.Page
		params.Limit = paging.Limit
		params.Offset = paging.Offset

		params, err = mu.validateFilterParamsList(params)
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		memberResult := <-mu.MemberQueryRead.GetListMembers(ctxReq, params)
		if memberResult.Error != nil {
			httpStatus := http.StatusInternalServerError

			// when data is not found
			if memberResult.Error == sql.ErrNoRows {
				httpStatus = http.StatusNotFound
				memberResult.Error = fmt.Errorf(helper.ErrorDataNotFound, labelMember)
			}

			output <- ResultUseCase{Error: memberResult.Error, HTTPStatus: httpStatus}
			return
		}

		member := memberResult.Result.(model.ListMembers)

		totalResult := <-mu.MemberQueryRead.GetTotalMembers(params)
		if totalResult.Error != nil {
			output <- ResultUseCase{Error: totalResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		member.TotalData = totalResult.Result.(int)

		output <- ResultUseCase{Result: member}

	}()

	return output
}

func (mu *MemberUseCaseImpl) validateFilterParamsList(params *model.Parameters) (*model.Parameters, error) {
	var err error
	if len(params.Status) > 0 {
		params.Status = strings.ToUpper(params.Status)
		if !helper.StringInSlice(params.Status, []string{model.ActiveString, model.InactiveString, model.NewString, model.BlockedString}) {
			err = fmt.Errorf(helper.ErrorParameterInvalid, scopeStatus)
			return params, err
		}
	}

	if len(params.OrderBy) > 0 && !helper.StringInSlice(params.OrderBy, model.AllowedSortFields) {
		err = fmt.Errorf(helper.ErrorParameterInvalid, "order by")
		return params, err
	}

	if len(params.Sort) > 0 && !helper.StringInSlice(params.Sort, []string{"asc", "desc"}) {
		err = fmt.Errorf(helper.ErrorParameterInvalid, "sort")
		return params, err
	}

	if params.OrderBy == "" {
		params.OrderBy = "id"
	}
	if params.Sort == "" {
		params.Sort = "asc"
	}
	if goString.IsValidEmail(params.Query) {
		params.Email = params.Query
		params.Query = ""
	}
	if strings.HasPrefix(params.Query, "USR") {
		params.UserID = params.Query
		params.Query = ""
	}

	return params, err
}

// GetHistorySession  function for getting detail session by param
func (mu *MemberUseCaseImpl) GetHistorySession(ctxReq context.Context, activeID []string, params *model.ParametersLoginActivity) <-chan ResultUseCase {
	ctx := "MemberUseCase-GetHistorySession"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		excludeID := strings.Join(activeID, ",")
		// validate all parameters
		paging, err := helper.ValidatePagination(
			helper.PaginationParameters{
				Page:     1, // default
				StrPage:  params.StrPage,
				Limit:    5, // default
				StrLimit: params.StrLimit,
			})

		if err != nil {
			tags[helper.TextResponse] = err.Error()
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		params.Page = paging.Page
		params.Limit = paging.Limit
		params.Offset = paging.Offset
		params.ExcludeID = excludeID
		tags["params"] = params
		session := <-mu.SessionQueryRead.GetHistorySessionInfo(ctxReq, params)
		if session.Error != nil {
			err := errors.New(model.ErrorGetLoginActivity)
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		resultSession := session.Result.(sessionModel.SessionInfoList)

		historySession := []model.SessionInfoDetail{}
		for _, s := range resultSession.Data {
			browser, isMobile, isApp := helper.ParseUserAgent(*s.UserAgent)
			data := model.SessionInfoDetail{
				ID:         *s.ID,
				DeviceType: *s.ClientType,
				IP:         *s.IP,
				UserAgent:  browser,
				ActiveNow:  false,
				LastLogin:  s.CreatedAt,
				IsMobile:   isMobile,
				IsApp:      isApp,
			}
			historySession = append(historySession, data)
		}

		totalResult := <-mu.SessionQueryRead.GetTotalHistorySessionInfo(ctxReq, params)
		if totalResult.Error != nil {
			output <- ResultUseCase{Error: totalResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		total := totalResult.Result.(int)

		resultSessionHistory := model.SessionHistoryInfoList{}
		resultSessionHistory.Data = historySession
		resultSessionHistory.TotalData = total

		output <- ResultUseCase{Result: resultSessionHistory}
	})
	return output
}

func (mu *MemberUseCaseImpl) conditionalEmail(ctxReq context.Context, data *model.Member, plEmail model.MemberEmailQueue) error {
	if data.IsSocialMediaExist() {
		go mu.QPublisher.QueueJob(ctxReq, plEmail, data.ID, "SendEmailWelcomeMember")
	} else {
		if data.SignUpFrom != model.Dolphin {
			go mu.QPublisher.QueueJob(ctxReq, plEmail, data.ID, "SendEmailRegisterMember")
		}
	}

	if data.APIVersion == helper.Version3 && data.SignUpFrom == model.Dolphin {
		go mu.QPublisher.QueueJob(ctxReq, plEmail, data.ID, "SendEmailAddMember")
	}

	return nil
}

// Clients  function for getting detail session by param
func (mu *MemberUseCaseImpl) Clients(ctxReq context.Context, token string) <-chan ResultUseCase {
	ctx := "MemberUseCase-Clients"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		var result model.User

		claims := jwt.MapClaims{}
		token = strings.Replace(token, "Bearer ", "", -1)
		jwtResult, err := jwt.ParseWithClaims(token, claims, func(tkn *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if (jwtResult == nil && err != nil) || len(claims) == 0 {
			output <- ResultUseCase{Error: errors.New("invalid token")}
			return
		}

		result.Email = claims["email"].(string)

		// get member
		result = mu.GetUserMember(ctxReq, claims, result)

		// get b2c_merchant
		result = mu.GetUserMerchant(ctxReq, claims, result)

		// get b2b_contact
		result = mu.GetUserCorporate(ctxReq, claims, result)

		output <- ResultUseCase{Result: result}
	})
	return output
}

func (mu *MemberUseCaseImpl) GetUserMember(ctxReq context.Context, claims jwt.MapClaims, result model.User) model.User {
	memberResult := <-mu.MemberQueryRead.FindByEmail(ctxReq, claims["email"].(string))
	if memberResult.Error == nil {
		member, ok := memberResult.Result.(model.Member)
		if ok {
			result.Users = append(result.Users, model.ListClient{
				UserType:  authModel.UserTypePersonal,
				FirstName: &member.FirstName,
				LastName:  &member.LastName,
				Logo:      member.ProfilePicture,
				IsSync:    member.IsSync,
			})

			if member.Password != "" {
				result.HasPassword = true
			}
		}

		// get b2c_merchant_employee
		result = mu.GetUserMerchantEmployee(ctxReq, member.ID, result)
	}

	return result
}

func (mu *MemberUseCaseImpl) GetUserMerchant(ctxReq context.Context, claims jwt.MapClaims, result model.User) model.User {
	merchantResult := mu.MerchantRepoRead.FindMerchantByEmail(ctxReq, claims["email"].(string))
	if merchantResult.Result != nil {
		merchant := merchantResult.Result.(merchantModel.B2CMerchantDataV2)
		result.Users = append(result.Users, model.ListClientMerchant{
			UserType:                 authModel.UserTypeMerchant,
			SellerId:                 merchant.ID,
			MerchantName:             merchant.MerchantName,
			Logo:                     merchant.MerchantLogo.String,
			MerchantServiceAvailable: true,
			VanityURL:                merchant.VanityURL.String,
			IsActive:                 merchant.IsActive,
			UpgradeStatus:            merchant.UpgradeStatus.String,
			MerchantType:             merchant.MerchantTypeString.String,
		})
	}

	return result
}

func (mu *MemberUseCaseImpl) GetUserMerchantEmployee(ctxReq context.Context, memberId string, result model.User) model.User {

	filter := &merchantModel.QueryMerchantEmployeeParameters{}
	filter.MemberID = memberId
	merchantEmployeeResult := <-mu.MerchantEmployeeRead.GetAllMerchantEmployees(ctxReq, filter)
	if merchantEmployeeResult.Result != nil {
		merchantEmployee := merchantEmployeeResult.Result.([]merchantModel.B2CMerchantEmployeeData)
		for _, val := range merchantEmployee {
			result.Sellers = append(result.Sellers, model.ListClientMerchant{
				UserType:     val.MerchantID,
				FirstName:    val.FirstName,
				LastName:     val.LastName.String,
				MerchantName: val.MerchantName.String,
				Logo:         val.MerchantLogo.String,
				MerchantType: val.MerchantType.String,
				VanityURL:    val.VanityURL.String,
				IsActive:     val.IsActive,
				IsPKP:        val.IsPKP,
			})
		}
	}

	return result
}

func (mu *MemberUseCaseImpl) GetUserCorporate(ctxReq context.Context, claims jwt.MapClaims, result model.User) model.User {
	corporateResult := <-mu.CorporateContactQueryRead.FindContactByEmail(ctxReq, claims["email"].(string))
	if corporateResult.Error == nil {
		corporate, ok := corporateResult.Result.(sharedModel.B2BContactData)
		if ok {
			for _, val := range corporate.TransactionType {
				if val.Microsite == authModel.UserTypeCorporate {
					result.Users = append(result.Users, model.ListClient{
						UserType:        val.Microsite,
						FirstName:       &corporate.FirstName,
						LastName:        &corporate.LastName,
						Logo:            corporate.Avatar,
						IsSync:          corporate.IsSync,
						TransactionType: val.Type,
					})
				} else {
					result.Microsites = append(result.Microsites, model.ListClient{
						UserType:        val.Microsite,
						FirstName:       &corporate.FirstName,
						LastName:        &corporate.LastName,
						Logo:            corporate.Avatar,
						IsSync:          corporate.IsSync,
						TransactionType: val.Type,
					})
				}
			}
		}
	}

	return result
}

// ActivateMerchantEmployee function for activating inactive member
func (mu *MemberUseCaseImpl) ActivateMerchantEmployee(ctxReq context.Context, token string) <-chan ResultUseCase {
	ctx := "MemberUseCase-ActivateMerchantEmployee"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags[helper.TextParameter] = token
		if len(token) == 0 {
			err := fmt.Errorf(helper.ErrorParameterRequired, "token")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// decode token
		data, err := hex.DecodeString(token) // convert to byte
		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		err, decode := golib.Decrypt(data, "SECRET")
		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// merchantId-memberId
		extractToken := strings.Split(decode, "-")

		// find merchant employee
		params := &merchantModel.QueryMerchantEmployeeParameters{
			MerchantID: extractToken[0],
			MemberID:   extractToken[1],
		}
		merchantEmployeeResult := <-mu.MerchantEmployeeRead.GetMerchantEmployees(ctxReq, params)
		if merchantEmployeeResult.Error != nil {
			tags[helper.TextResponse] = merchantEmployeeResult.Error
			output <- ResultUseCase{Error: merchantEmployeeResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		merchantEmployee, ok := merchantEmployeeResult.Result.(merchantModel.B2CMerchantEmployeeData)
		if !ok {
			err := errors.New(msgErrorResultMerchantEmployee)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// check status employee must invited
		if merchantEmployee.Status.String != helper.TextInvited {
			err := errors.New(msgErrorActivationMerchantEmployee)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// change status
		payload := merchantModel.B2CMerchantEmployee{}
		payload.MerchantID = extractToken[0]
		payload.MemberID = extractToken[1]
		payload.Status = helper.TextActive
		payload.ModifiedAt = time.Now()
		payload.ModifiedBy = &extractToken[1]
		err = mu.MerchantEmployeeRead.ChangeStatus(ctxReq, payload)
		if err != nil {
			err := errors.New("failed change status employee")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: merchantEmployee}
	})

	return output
}
