package usecase

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/golang-jwt/jwt"
)

const StringBearer = "Bearer "
const ErrorInvalidToken = "invalid token"
const ErrorMerchantNotFound = "merchant not found"
const ErrorMerchantNotRegister = "you are not registered as a merchant"
const ParamsResendEmail = "{resend-email}"

func (m *MerchantUseCaseImpl) AddEmployee(ctxReq context.Context, token, email, firstName string) <-chan ResultUseCase {
	ctx := "MerchantUseCase-AddEmployee"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {

		claims := jwt.MapClaims{}
		token = strings.Replace(token, StringBearer, "", -1)
		jwtResult, err := jwt.ParseWithClaims(token, claims, func(tkn *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if (jwtResult == nil && err != nil) || len(claims) == 0 {
			output <- ResultUseCase{Error: errors.New(ErrorInvalidToken)}
			return
		}

		// validate email merchant
		if claims["email"].(string) == email {
			output <- ResultUseCase{Error: errors.New("your email is already registered as a merchant")}
			return
		}

		// get merchant
		merchantResult := m.MerchantRepo.FindMerchantByEmail(ctxReq, claims["email"].(string))
		if merchantResult.Error != nil {
			output <- ResultUseCase{Error: errors.New(ErrorMerchantNotRegister)}
			return
		}
		merchant := merchantResult.Result.(model.B2CMerchantDataV2)

		token, err = m.CreateMerchantEmployee(ctxReq, merchant, email, firstName, claims)
		if err != nil {
			output <- ResultUseCase{Error: err}
			return
		}

		output <- ResultUseCase{Result: token}
	})
	return output

}

func (m *MerchantUseCaseImpl) CreateMerchantEmployee(ctxReq context.Context, merchant model.B2CMerchantDataV2, email, firstName string, claims jwt.MapClaims) (string, error) {
	// get member
	var err error
	paramsMember := memberModel.Member{}
	template_key := "SendEmailMerchantEmployeeLogin"
	memberNotFound := false
	memberResult := <-m.MemberQueryRead.FindByEmail(ctxReq, email)
	if memberResult.Result == nil {
		template_key = "SendEmailMerchantEmployeeRegister"
		memberNotFound = true

		// check firstname == resend-email and resul == nil then error
		if firstName == ParamsResendEmail {
			err := errors.New("email is not registered as a member")
			return "", err
		}

		// created member
		paramsMember, err = m.CreateMember(ctxReq, merchant, email, firstName)
		if err != nil {
			return "", err
		}
	}

	member := memberModel.Member{}
	if memberNotFound {
		member = memberModel.Member(paramsMember)
	} else {
		member = memberResult.Result.(memberModel.Member)

		// generate token
		mix := merchant.ID + "-" + member.ID
		_, v := golib.Encrypt([]byte(mix), "SECRET")
		member.Token = hex.EncodeToString(v)

		saveResult := <-m.MemberRepoRead.Save(ctxReq, member)
		if saveResult.Error != nil {
			return "", saveResult.Error
		}
	}

	// checking merchant employee
	err = m.CheckingMerchantEmployee(ctxReq, merchant, &member, claims, firstName)
	if err != nil {
		return "", err
	}

	// send email
	plQueue := memberModel.MemberPayloadEmail{
		Merchant: &merchant,
		Member:   &member,
	}

	go m.QueuePublisher.QueueJob(ctxReq, plQueue, merchant.ID, template_key)

	return member.Token, nil
}

func (m *MerchantUseCaseImpl) CheckingMerchantEmployee(ctxReq context.Context, merchant model.B2CMerchantDataV2, member *memberModel.Member, claims jwt.MapClaims, firstName string) error {
	// checking merchant employee
	// get data employee
	params := &model.QueryMerchantEmployeeParameters{
		MerchantID: merchant.ID,
		MemberID:   member.ID,
	}

	mr := <-m.MerchantEmployeeRepo.GetMerchantEmployees(ctxReq, params)
	if mr.Result == nil {
		// check firstname == resend-email and resul == nil then error
		if firstName == ParamsResendEmail {
			err := errors.New("email is not registered as a employee")
			return err
		}

		// save merchant employee
		param := model.B2CMerchantEmployee{}
		param.ID = "EMP" + time.Now().Format(helper.FormatYmdhis)
		param.MerchantID = merchant.ID
		param.MemberID = member.ID
		param.CreatedAt = time.Now()
		param.CreatedBy = claims["sub"].(string)

		err := m.MerchantEmployeeRepo.Save(param)
		if err != nil {
			return err
		}
	} else {
		if firstName == ParamsResendEmail {
			merchantEmployee := mr.Result.(model.B2CMerchantEmployeeData)
			if merchantEmployee.Status.String != helper.TextInvited {
				return errors.New("employee status must be invited")
			}
		} else {
			merchantEmployee := mr.Result.(model.B2CMerchantEmployeeData)
			err := m.CheckEmployeeRevoked(ctxReq, merchantEmployee)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
func (m *MerchantUseCaseImpl) CheckEmployeeRevoked(ctxReq context.Context, merchantEmployee model.B2CMerchantEmployeeData) error {
	if merchantEmployee.Status.String == helper.TextRevoked {
		param := model.B2CMerchantEmployee{}
		param.ID = merchantEmployee.ID
		param.MerchantID = merchantEmployee.MerchantID
		param.MemberID = merchantEmployee.MemberID
		param.ModifiedAt = time.Now()
		param.ModifiedBy = merchantEmployee.ModifiedBy
		param.Status = helper.TextInvited
		err := m.MerchantEmployeeRepo.ChangeStatus(ctxReq, param)
		if err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("email already register")
	}
}

func (m *MerchantUseCaseImpl) CreateMember(ctxReq context.Context, merchant model.B2CMerchantDataV2, email, firstName string) (memberModel.Member, error) {
	paramsMember := memberModel.Member{}
	paramsMember.ID = helper.GenerateMemberIDv2()
	paramsMember.Email = email
	paramsMember.FirstName = firstName
	paramsMember.Status = memberModel.New
	paramsMember.SignUpFrom = "seller"

	// generate token
	mix := merchant.ID + "-" + paramsMember.ID
	_, v := golib.Encrypt([]byte(mix), "SECRET")
	paramsMember.Token = hex.EncodeToString(v)

	saveResult := <-m.MemberRepoRead.Save(ctxReq, paramsMember)
	if saveResult.Error != nil {
		return memberModel.Member{}, saveResult.Error
	}

	// check sign up from process
	eventType := "register"

	paramsMember.StatusString = strings.ToUpper(model.NewString)
	paramsMember.Created = time.Now()
	go m.PublishToKafkaUser(ctxReq, &paramsMember, eventType)

	plLog := memberModel.MemberLog{
		Before: &memberModel.Member{},
		After:  &paramsMember,
	}

	// send to audit trail activity service
	go m.QueuePublisher.QueueJob(ctxReq, plLog, paramsMember.ID, "InsertLogRegisterMember")

	return paramsMember, nil
}

// PublishToKafkaUser function for publish to kafka user
func (m *MerchantUseCaseImpl) PublishToKafkaUser(ctxReq context.Context, dataMember *memberModel.Member, eventType string) error {
	ctx := "MerchantUseCase-PublishToKafkaUser"

	trace := tracer.StartTrace(ctxReq, ctx)
	tags := make(map[string]interface{})

	dataMember.GenderString = dataMember.Gender.String()
	dataMember.Gender = memberModel.StringToGender(dataMember.GenderString)
	memberDolpin := serviceModel.MemberDolphin{
		ID:              dataMember.ID,
		Email:           strings.ToLower(dataMember.Email),
		FirstName:       dataMember.FirstName,
		LastName:        dataMember.LastName,
		Gender:          dataMember.Gender.GetDolpinGender(),
		DOB:             dataMember.BirthDateString,
		Phone:           dataMember.Phone,
		Ext:             dataMember.Ext,
		Mobile:          dataMember.Mobile,
		Street1:         dataMember.Address.Street1,
		Street2:         dataMember.Address.Street2,
		PostalCode:      dataMember.Address.ZipCode,
		SubDistrictID:   dataMember.Address.SubDistrictID,
		SubDistrictName: dataMember.Address.SubDistrict,
		DistrictID:      dataMember.Address.DistrictID,
		DistrictName:    dataMember.Address.District,
		CityID:          dataMember.Address.CityID,
		CityName:        dataMember.Address.City,
		ProvinceID:      dataMember.Address.ProvinceID,
		ProvinceName:    dataMember.Address.Province,
		Status:          strings.ToUpper(dataMember.StatusString),
		Created:         dataMember.Created.Format(time.RFC3339),
		LastModified:    time.Now().Format(time.RFC3339),
		FacebookID:      dataMember.SocialMedia.FacebookID,
		GoogleID:        dataMember.SocialMedia.GoogleID,
		LDAPID:          dataMember.SocialMedia.LDAPID,
		AppleID:         dataMember.SocialMedia.AppleID,
	}

	payload := serviceModel.DolphinPayloadNSQ{
		EventOrchestration:     "UpdateMember",
		TimestampOrchestration: time.Now().Format(time.RFC3339),
		EventType:              eventType,
		Counter:                0,
		Payload:                memberDolpin,
	}

	// prepare to send to nsq
	payloadJSON, e := json.Marshal(payload)
	if e != nil {
		helper.SendErrorLog(ctxReq, ctx, "publish_payload", e, payload)
		return e
	}

	messageKey := memberDolpin.ID
	tags[helper.TextArgs] = payload
	trace.Finish(tags)

	// excluded from parent context for kafka producer
	return m.QueuePublisher.PublishKafka(trace.NewChildContext(), os.Getenv("KAFKA_USER_SERVICE_TOPIC"), messageKey, payloadJSON)
}

// GetAllMerchantEmployee return all merchants by given parameters
func (m *MerchantUseCaseImpl) GetAllMerchantEmployee(ctxReq context.Context, token string, params *model.QueryMerchantEmployeeParameters) <-chan ResultUseCase {
	ctx := "MerchantUseCase-GetAllMerchantEmployee"
	output := make(chan ResultUseCase)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		claims := jwt.MapClaims{}
		token = strings.Replace(token, StringBearer, "", -1)
		jwtResult, err := jwt.ParseWithClaims(token, claims, func(tkn *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if (jwtResult == nil && err != nil) || len(claims) == 0 {
			output <- ResultUseCase{Error: errors.New(ErrorInvalidToken)}
			return
		}

		paging, err := helper.ValidatePagination(
			helper.PaginationParameters{
				Page:     1,
				StrPage:  params.StrPage,
				Limit:    10,
				StrLimit: params.StrLimit,
			},
		)
		if err != nil {
			tags[helper.TextResponse] = err.Error()
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		params.Offset = paging.Offset
		params.Page = paging.Page
		params.Limit = paging.Limit
		if params.Status == "" {
			params.Status = "INVITED,ACTIVE,INACTIVE"
		}

		// get merchant
		merchantResult := m.MerchantRepo.FindMerchantByUser(ctxReq, claims["sub"].(string))
		if merchantResult.Error != nil {
			err := errors.New(ErrorMerchantNotRegister)
			tags[helper.TextResponse] = merchantResult.Error
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		merchant := merchantResult.Result.(model.B2CMerchantDataV2)
		params.MerchantID = merchant.ID

		mr := <-m.MerchantEmployeeRepo.GetAllMerchantEmployees(ctxReq, params)
		if mr.Error != nil {
			tags[helper.TextResponse] = mr.Error.Error()
			output <- ResultUseCase{Error: mr.Error, HTTPStatus: http.StatusBadRequest}
			return
		}
		merchantEmployee := mr.Result.([]model.B2CMerchantEmployeeData)

		merchantEmployeeQuery := <-m.MerchantEmployeeRepo.GetTotalMerchantEmployees(ctxReq, params)
		if merchantEmployeeQuery.Error != nil {
			output <- ResultUseCase{Error: merchantEmployeeQuery.Error, HTTPStatus: http.StatusBadRequest}
			return
		}
		totalData := merchantEmployeeQuery.Result.(int)

		output <- ResultUseCase{Result: merchantEmployee, TotalData: totalData}
	})

	return output
}

// GetMerchantEmployee return all merchants by given parameters
func (m *MerchantUseCaseImpl) GetMerchantEmployee(ctxReq context.Context, token string, params *model.QueryMerchantEmployeeParameters) <-chan ResultUseCase {
	ctx := "MerchantUseCase-GetMerchantEmployee"
	output := make(chan ResultUseCase)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		claims := jwt.MapClaims{}
		token = strings.Replace(token, StringBearer, "", -1)
		jwtResult, err := jwt.ParseWithClaims(token, claims, func(tkn *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if (jwtResult == nil && err != nil) || len(claims) == 0 {
			output <- ResultUseCase{Error: errors.New(ErrorInvalidToken)}
			return
		}

		// get merchant
		merchantResult := m.MerchantRepo.FindMerchantByUser(ctxReq, claims["sub"].(string))
		if merchantResult.Error != nil {
			err := errors.New(ErrorMerchantNotRegister)
			tags[helper.TextResponse] = merchantResult.Error
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		merchant := merchantResult.Result.(model.B2CMerchantDataV2)
		params.MerchantID = merchant.ID

		// get data employee
		mr := <-m.MerchantEmployeeRepo.GetMerchantEmployees(ctxReq, params)
		if mr.Error != nil {
			err := errors.New("employee not found")
			tags[helper.TextResponse] = mr.Error
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		merchantEmployee := mr.Result.(model.B2CMerchantEmployeeData)

		output <- ResultUseCase{Result: merchantEmployee}
	})

	return output
}

// UpdateMerchantEmployee return all merchants by given parameters
func (m *MerchantUseCaseImpl) UpdateMerchantEmployee(ctxReq context.Context, token string, params *model.QueryMerchantEmployeeParameters) <-chan ResultUseCase {
	ctx := "MerchantUseCase-UpdateMerchantEmployee"
	output := make(chan ResultUseCase)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		claims := jwt.MapClaims{}
		token = strings.Replace(token, StringBearer, "", -1)
		jwtResult, err := jwt.ParseWithClaims(token, claims, func(tkn *jwt.Token) (interface{}, error) {
			return []byte("secret"), nil
		})

		if (jwtResult == nil && err != nil) || len(claims) == 0 {
			output <- ResultUseCase{Error: errors.New(ErrorInvalidToken)}
			return
		}

		// get merchant
		merchantResult := m.MerchantRepo.FindMerchantByUser(ctxReq, claims["sub"].(string))
		if merchantResult.Error != nil {
			err := errors.New(ErrorMerchantNotRegister)
			tags[helper.TextResponse] = merchantResult.Error
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		merchant := merchantResult.Result.(model.B2CMerchantDataV2)
		params.MerchantID = merchant.ID

		// get data employee
		filter := model.QueryMerchantEmployeeParameters{}
		filter.MerchantID = params.MerchantID
		filter.MemberID = params.MemberID

		mr := <-m.MerchantEmployeeRepo.GetMerchantEmployees(ctxReq, &filter)
		if mr.Error != nil {
			err := errors.New("employee not found")
			tags[helper.TextResponse] = mr.Error
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		merchantEmployee := mr.Result.(model.B2CMerchantEmployeeData)

		// validate status
		if err := params.ValidateStatus(merchantEmployee.Status.String, params.Status); err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// updateStatus
		modifiedBy := claims["sub"].(string)

		payload := model.B2CMerchantEmployee{}
		payload.Status = params.Status
		payload.MerchantID = params.MerchantID
		payload.MemberID = params.MemberID
		payload.ModifiedAt = time.Now()
		payload.ModifiedBy = &modifiedBy
		err = m.MerchantEmployeeRepo.ChangeStatus(ctxReq, payload)
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		merchantEmployee.Status.String = params.Status

		output <- ResultUseCase{Result: merchantEmployee}
	})

	return output
}
