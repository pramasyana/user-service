package usecase

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"

	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/gosimple/slug"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"
)

const private = "private"

// GetDataUserMerchant func
func (m *MerchantUseCaseImpl) GetDataUserMerchant(ctxReq context.Context, merchantID string) *serviceModel.SendbirdRequest {

	var response serviceModel.SendbirdRequest

	merchantResult := m.MerchantRepo.LoadMerchant(ctxReq, merchantID, private)
	merchantData := merchantResult.Result.(model.B2CMerchantDataV2)
	if merchantResult.Result != nil {
		response.Metadata.Merchant.IsMerchant = "true"
		response.Metadata.Merchant.MerchantID = merchantData.ID
		response.Metadata.Merchant.MerchantName = merchantData.MerchantName
		response.Metadata.MerchantLogo = merchantData.MerchantLogo
	}

	emailResult := <-m.MemberQueryRead.FindByID(ctxReq, merchantData.UserID)
	if emailResult.Result != nil {
		member := emailResult.Result.(memberModel.Member)
		response.UserID = merchantData.UserID
		response.NickName = member.FirstName
		response.ProfileURL = member.ProfilePicture
	}

	return &response
}

// GetDataUserMerchantV4 func
func (m *MerchantUseCaseImpl) GetDataUserMerchantV4(ctxReq context.Context, merchantID string) *serviceModel.SendbirdRequestV4 {

	var response serviceModel.SendbirdRequestV4

	merchantResult := m.MerchantRepo.LoadMerchant(ctxReq, merchantID, private)
	merchantData := merchantResult.Result.(model.B2CMerchantDataV2)
	if merchantResult.Result != nil {
		response.UserID = merchantData.ID
		response.NickName = merchantData.MerchantName
		response.ProfileURL = merchantData.MerchantLogo.String
		response.MetadataV4.Reference = merchantData.UserID
	}

	return &response
}

// CreateMerchantSendbirdV4  func
func (m *MerchantUseCaseImpl) CreateMerchantSendbirdV4(ctxReq context.Context, oldData model.B2CMerchantDataV2, payload *model.B2CMerchantCreateInput) error {

	// send email only if current upgradeStatus is not ACTIVE and new payload is ACTIVE
	if !oldData.IsActive && payload.IsActive {

		inputData := m.GetDataUserMerchantV4(ctxReq, payload.ID)

		// check userID sendbird
		checkUserSendbird := m.SendbirdService.CheckUserSenbirdV4(ctxReq, inputData)

		if checkUserSendbird.Error != nil {

			checkStatusCode := checkUserSendbird.Result.(serviceModel.SendbirdStringResponseV4)
			// check if user sendbird is not found
			if checkStatusCode.Code == 400201 {

				// create user sendbird
				createUserSendbird := m.SendbirdService.CreateUserSendbirdV4(ctxReq, inputData)
				if createUserSendbird.Error != nil {
					return createUserSendbird.Error
				}

				return nil

			}

			return checkUserSendbird.Error
		}
		return nil

	}
	return nil
}

func (m *MerchantUseCaseImpl) UpdateMerchantSendbirdV4(ctxReq context.Context, oldData model.B2CMerchantDataV2, payload *model.B2CMerchantCreateInput) error {
	// send email only if current upgradeStatus is not ACTIVE and new payload is ACTIVE
	if !oldData.IsActive && payload.IsActive {

		inputData := m.GetDataUserMerchantV4(ctxReq, payload.ID)

		// check userID sendbird
		checkUserSendbird := m.SendbirdService.CheckUserSenbirdV4(ctxReq, inputData)

		if checkUserSendbird.Error != nil {

			checkStatusCode := checkUserSendbird.Result.(serviceModel.SendbirdStringResponseV4)
			// check if user sendbird is not found
			if checkStatusCode.Code == 400201 {

				// create user sendbird
				createUserSendbird := m.SendbirdService.UpdateUserSendbirdV4(ctxReq, inputData)
				if createUserSendbird.Error != nil {
					return createUserSendbird.Error
				}

				return nil

			}

			return checkUserSendbird.Error
		}
		return nil

	}
	return nil
}

// UpdateMerchant update merchant
func (m *MerchantUseCaseImpl) UpdateMerchant(ctxReq context.Context, payload *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase {
	ctx := "MerchantUseCase-UpdateMerchant"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		payload.Status, payload.IsActive = validateStatusIsActive(payload)
		if err := m.validateRequiredData(ctxReq, payload); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		merchant := m.MerchantRepo.LoadMerchant(ctxReq, payload.ID, private)
		if merchant.Error != nil && merchant.Error == sql.ErrNoRows {
			output <- ResultUseCase{Error: errMerchantNotFound, HTTPStatus: http.StatusNotFound}
			return
		}

		existingMerchant := merchant.Result.(model.B2CMerchantDataV2)
		oldData := existingMerchant
		if err := m.ValidateMerchantBank(ctxReq, payload, &existingMerchant); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		payload.MerchantEmail = existingMerchant.MerchantEmail.String
		m.GetCorporateData(&existingMerchant, payload)
		existingMerchant.SetMerchantData(payload)
		existingMerchant.EditorID = null.StringFrom(userAttribute.UserID)
		existingMerchant.EditorIP = null.StringFrom(userAttribute.UserIP)
		existingMerchant.LastModified = null.TimeFrom(time.Now())
		existingMerchant.Version = zero.IntFrom(existingMerchant.Version.ValueOrZero() + 1)
		existingMerchant.GenderPic = payload.GenderPic
		existingMerchant.MerchantGroup = zero.StringFrom(payload.MerchantGroup)
		existingMerchant.GenderPicString = zero.StringFrom(payload.GenderPic.String())
		existingMerchant.ProductType = zero.StringFrom(payload.ProductType)
		existingMerchant.LegalEntity = zero.IntFrom(int64(payload.LegalEntity))
		existingMerchant.NumberOfEmployee = zero.IntFrom(int64(payload.NumberOfEmployee))
		existingMerchant.MerchantTypeString = zero.StringFrom(payload.MerchantTypeString)
		existingMerchant.MerchantType = model.StringToMerchantType(payload.MerchantTypeString)

		existingMerchant.Status, existingMerchant.IsActive = validateExistingStatusIsActive(oldData, existingMerchant)
		// set disallow field
		existingMerchant.VanityURL = oldData.VanityURL
		existingMerchant.MerchantName = oldData.MerchantName
		existingMerchant.MerchantEmail = oldData.MerchantEmail

		m.Repository.StartTransaction()
		updateResult := <-m.MerchantRepo.AddUpdateMerchant(ctxReq, existingMerchant)
		if updateResult.Error != nil {
			output <- ResultUseCase{Error: errors.New("failed to update merchant"), HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}
		documentProccess := <-m.MerchantDocumentsProcess(ctxReq, payload.Documents, existingMerchant.ID, userAttribute)
		if documentProccess.Error != nil {
			output <- ResultUseCase{Error: documentProccess.Error, HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}
		existingMerchant.Documents = documentProccess.Result.([]model.B2CMerchantDocumentData)

		// create user sendbird
		if err := m.CreateMerchantSendbirdV4(ctxReq, oldData, payload); err != nil {
			m.Repository.Rollback()
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		if err := m.sendEmailActivationOrApproval(ctxReq, oldData, payload); err != nil {
			m.Repository.Rollback()
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		m.Repository.Commit()
		plLog := model.MerchantLog{
			Before: oldData,
			After:  existingMerchant,
		}

		go m.QueuePublisher.QueueJob(ctxReq, plLog, existingMerchant.ID, "InsertLogMerchantUpdate")

		output <- ResultUseCase{Result: existingMerchant}
	})
	return output
}

// DeleteMerchant soft delete by flagging on DB
func (m *MerchantUseCaseImpl) DeleteMerchant(ctxReq context.Context, merchantID string, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase {
	ctx := "MerchantUseCase-DeleteMerchant"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags[helper.TextMerchantIDCamel] = merchantID

		merchant := m.MerchantRepo.LoadMerchant(ctxReq, merchantID, private)
		if merchant.Error != nil && merchant.Error == sql.ErrNoRows {
			output <- ResultUseCase{Error: errMerchantNotFound, HTTPStatus: http.StatusNotFound}
			return
		}

		existingMerchant := merchant.Result.(model.B2CMerchantDataV2)
		existingMerchant.Status = model.DeletedString
		oldData := existingMerchant
		del := <-m.MerchantRepo.SoftDelete(ctxReq, merchantID)
		if del.Error != nil {
			output <- ResultUseCase{Error: del.Error}
			return
		}

		existingMerchant.DeletedAt = null.TimeFrom(time.Now())
		plLog := model.MerchantLog{
			Before: oldData,
			After:  existingMerchant,
		}

		go m.QueuePublisher.QueueJob(ctxReq, plLog, merchantID, "InsertLogMerchantDelete")

		output <- ResultUseCase{Result: existingMerchant}

	})
	return output
}

// GetMerchants return all merchants by given parameters
func (m *MerchantUseCaseImpl) GetMerchants(ctxReq context.Context, params *model.QueryParameters) <-chan ResultUseCase {
	ctx := "MerchantUseCase-GetMerchants"
	output := make(chan ResultUseCase)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

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

		mr := <-m.MerchantRepo.GetMerchants(ctxReq, params)
		if mr.Error != nil {
			tags[helper.TextResponse] = mr.Error.Error()
			output <- ResultUseCase{Error: mr.Error, HTTPStatus: http.StatusBadRequest}
			return
		}
		merchants := mr.Result.([]model.B2CMerchantDataV2)
		for idx, val := range merchants {
			mapsResult := <-m.MerchantAddressRepo.FindAddressMaps(ctxReq, val.ID, "b2c_merchant")
			if mapsResult.Error == nil {
				detailMaps, _ := mapsResult.Result.(model.Maps)
				merchants[idx].IsMapAvailable = helper.ValidateLatLong(detailMaps.Latitude, detailMaps.Longitude)
			}
		}

		merchantQuery := <-m.MerchantRepo.GetTotalMerchant(ctxReq, params)
		if merchantQuery.Error != nil {
			output <- ResultUseCase{Error: merchantQuery.Error, HTTPStatus: http.StatusBadRequest}
			return
		}
		totalData := merchantQuery.Result.(int)

		output <- ResultUseCase{Result: merchants, TotalData: totalData}
	})

	return output
}

// GetMerchantByID function for get merchant by merchant ID
func (m *MerchantUseCaseImpl) GetMerchantByID(ctxReq context.Context, id string, privacy string, isAttachment string) <-chan ResultUseCase {
	ctx := "MerchantUseCase-GetMerchantByID"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		merchantResult := m.MerchantRepo.LoadMerchant(ctxReq, id, privacy)
		if merchantResult.Result == nil {
			output <- ResultUseCase{Error: errMerchantNotFound, HTTPStatus: http.StatusNotFound}
			return
		}

		merchantData := merchantResult.Result.(model.B2CMerchantDataV2)
		if privacy == "private" {
			mapsResult := <-m.MerchantAddressRepo.FindAddressMaps(ctxReq, merchantData.ID, "b2c_merchant")
			if mapsResult.Error == nil {
				merchantData.Maps, _ = mapsResult.Result.(model.Maps)
				merchantData.IsMapAvailable = helper.ValidateLatLong(merchantData.Maps.Latitude, merchantData.Maps.Longitude)
			}

			merchantData = m.adjustMerchantData(ctxReq, merchantData, isAttachment)
			tags[helper.TextMerchantIDCamel] = id
			tags[helper.TextEmail] = merchantData.MerchantEmail
		}
		output <- ResultUseCase{Result: merchantData}

	})

	return output
}
func (m *MerchantUseCaseImpl) GetMerchantByVanityURL(ctxReq context.Context, vanityURL string) <-chan ResultUseCase {
	ctx := "MerchantUseCase-GetMerchantByVanityURL"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		merchantResult := m.MerchantRepo.LoadMerchantByVanityURL(ctxReq, vanityURL)
		if merchantResult.Result == nil {
			output <- ResultUseCase{Error: errMerchantNotFound, HTTPStatus: http.StatusNotFound}
			return
		}

		merchantData := merchantResult.Result.(model.B2CMerchantDataV2)
		mapsResult := <-m.MerchantAddressRepo.FindAddressMaps(ctxReq, merchantData.ID, "b2c_merchant")
		if mapsResult.Error == nil {
			merchantData.Maps, _ = mapsResult.Result.(model.Maps)
			merchantData.IsMapAvailable = helper.ValidateLatLong(merchantData.Maps.Latitude, merchantData.Maps.Longitude)
		}
		isAttachment := "false"
		tags[helper.TextMerchantIDCamel] = vanityURL
		tags[helper.TextEmail] = merchantData.MerchantEmail

		merchantData = m.adjustMerchantData(ctxReq, merchantData, isAttachment)

		output <- ResultUseCase{Result: merchantData}

	})

	return output
}

func (m *MerchantUseCaseImpl) findMemberByMerchantEmail(ctxReq context.Context, email string) (*memberModel.Member, error) {
	memberQuery := <-m.MemberQueryRead.FindByEmail(ctxReq, email)
	if memberQuery.Error != nil {
		err := memberQuery.Error
		if memberQuery.Error == sql.ErrNoRows {
			err = errors.New("email does not exist")
		}
		return nil, err
	}
	member, ok := memberQuery.Result.(memberModel.Member)
	if !ok {
		return nil, errors.New("malformed member data")
	}
	return &member, nil
}

// CreateMerchant function for create new merchant
func (m *MerchantUseCaseImpl) CreateMerchant(ctxReq context.Context, params *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase {
	ctx := "MerchantUseCase-CreateMerchant"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		merchantInput := model.B2CMerchantDataV2{}
		tags[helper.TextEmail] = params.MerchantEmail
		member, err := m.findMemberByMerchantEmail(ctxReq, params.MerchantEmail)
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		merchantInput.UserID = member.ID
		if params.Status == "" {
			params.Status = model.NewString
		}

		if err := m.validateBasicData(params); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// Validate merchant bank id
		if err := m.ValidateMerchantBank(ctxReq, params, &merchantInput); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// Generate ID
		params.ID = "MCH" + time.Now().Format(helper.FormatYmdhis)
		params.Source = "cms"
		params.UserID = member.ID

		// Make slugify
		params.VanityURL = slug.MakeLang(params.MerchantName, "en")

		// Validate merchant data
		if err := m.ValidateMerchantData(ctxReq, params); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		// set merchant data from params
		merchantInput.SetMerchantData(params)
		merchantInput.ID = params.ID
		merchantInput.CreatorID = null.StringFrom(userAttribute.UserID)
		merchantInput.CreatorIP = null.StringFrom(userAttribute.UserIP)
		merchantInput.Created = zero.TimeFrom(time.Now())
		merchantInput.LastModified = null.TimeFrom(time.Now())
		merchantInput.Version = zero.IntFrom(merchantInput.Version.ValueOrZero() + 1)
		merchantInput.MerchantTypeString = zero.StringFrom(model.RegularString)
		merchantInput.MerchantType = model.StringToMerchantType(model.RegularString)
		merchantInput.GenderPic = params.GenderPic
		merchantInput.MerchantGroup = zero.StringFrom(params.MerchantGroup)
		merchantInput.GenderPicString = zero.StringFrom(params.GenderPic.String())
		merchantInput.LegalEntity = zero.IntFrom(int64(params.LegalEntity))
		merchantInput.NumberOfEmployee = zero.IntFrom(int64(params.NumberOfEmployee))
		merchantInput.Status = model.NewString
		merchantInput.CountUpdateNameAvailable = 1

		m.Repository.StartTransaction()
		saveResult := <-m.MerchantRepo.AddUpdateMerchant(ctxReq, merchantInput)
		if saveResult.Error != nil {
			err := errors.New("failed to save merchant")
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}

		documentProccess := <-m.MerchantDocumentsProcess(ctxReq, params.Documents, params.ID, userAttribute)
		if documentProccess.Error != nil {
			output <- ResultUseCase{Error: documentProccess.Error, HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}

		merchantInput.Documents = documentProccess.Result.([]model.B2CMerchantDocumentData)

		plLog := model.MerchantLog{
			Before: model.B2CMerchantDataV2{},
			After:  merchantInput,
		}

		go m.QueuePublisher.QueueJob(ctxReq, plLog, merchantInput.ID, "InsertLogMerchantCreate")

		m.Repository.Commit()

		tags[helper.TextResponse] = merchantInput
		output <- ResultUseCase{Result: merchantInput}
	})

	return output
}

// RejectMerchantRegistration soft delete by flagging on DB
func (m *MerchantUseCaseImpl) RejectMerchantRegistration(ctxReq context.Context, merchantID string, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase {
	ctx := "MerchantUseCase-RejectMerchantRegistration"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags[helper.TextMerchantIDCamel] = merchantID

		merchantDB := m.MerchantRepo.LoadMerchant(ctxReq, merchantID, private)
		if merchantDB.Error != nil && merchantDB.Error == sql.ErrNoRows {
			output <- ResultUseCase{Error: errMerchantNotFound, HTTPStatus: http.StatusNotFound}
			return
		}

		existingMerchant := merchantDB.Result.(model.B2CMerchantDataV2)
		if existingMerchant.IsActive && existingMerchant.Status == model.ActiveString {
			output <- ResultUseCase{Error: errUnableToRejectRegistration, HTTPStatus: http.StatusBadRequest}
			return
		}

		oldData := existingMerchant
		del := <-m.MerchantRepo.SoftDelete(ctxReq, merchantID)
		if del.Error != nil {
			output <- ResultUseCase{Error: del.Error}
			return
		}

		existingMerchant.DeletedAt = null.TimeFrom(time.Now())
		existingMerchant.Status = model.DeletedString
		memberName := existingMerchant.MerchantName

		plQueue := model.MerchantPayloadEmail{
			MemberName: memberName,
			Data:       existingMerchant,
		}
		plLog := model.MerchantLog{
			Before: existingMerchant,
			After:  model.B2CMerchantDataV2{},
		}

		go m.QueuePublisher.QueueJob(ctxReq, plQueue, existingMerchant.ID, "SendEmailMerchantRejectRegistration")
		go m.QueuePublisher.QueueJob(ctxReq, plLog, existingMerchant.ID, "InsertLogMerchantDelete")
		go func() {
			m.PublishToKafkaMerchant(ctxReq, oldData, helper.EventProduceDeleteMerchant)
		}()

		output <- ResultUseCase{Result: existingMerchant}

	})
	return output
}

// RejectMerchantUpgrade revert previous upgradeStatus
func (m *MerchantUseCaseImpl) RejectMerchantUpgrade(ctxReq context.Context, merchantID string, userAttribute *model.MerchantUserAttribute, reasonReject string) <-chan ResultUseCase {
	ctx := "MerchantUseCase-RejectMerchantUpgrade"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags[helper.TextMerchantIDCamel] = merchantID

		merchantDB := m.MerchantRepo.LoadMerchant(ctxReq, merchantID, private)
		if merchantDB.Error != nil && merchantDB.Error == sql.ErrNoRows {
			output <- ResultUseCase{Error: errMerchantNotFound, HTTPStatus: http.StatusNotFound}
			return
		}

		merchantOnDB := merchantDB.Result.(model.B2CMerchantDataV2)
		// only check if upgrade status is not empty,
		// upgrade status should be PENDING_ASSOCIATE or PENDING_MANAGE
		// empty value on upgradeStatus mean not yet request for upgrade
		// avtive mean merchant already approved
		if merchantOnDB.UpgradeStatus.String == "" {
			output <- ResultUseCase{Error: errMerchantNotValidForUpgrade, HTTPStatus: http.StatusBadRequest}
			return
		}
		if merchantOnDB.UpgradeStatus.String == model.ActiveString {
			output <- ResultUseCase{Error: errUnableToRejectUpgrade, HTTPStatus: http.StatusBadRequest}
			return
		}

		oldData := merchantOnDB

		merchantOnDB.EditorID = null.StringFrom(userAttribute.UserID)
		merchantOnDB.EditorIP = null.StringFrom(userAttribute.UserIP)
		merchantOnDB.LastModified = null.TimeFrom(time.Now())
		merchantOnDB.MerchantType = model.Regular
		merchantOnDB.MerchantTypeString = zero.StringFrom(model.RegularString)
		merchantOnDB.Reason = zero.StringFrom(reasonReject)
		merchantOnDB = checkUpgradeStatusBeforeReject(merchantOnDB)

		m.Repository.StartTransaction()

		reject := <-m.MerchantRepo.RejectUpgrade(ctxReq, merchantOnDB, reasonReject)
		if reject.Error != nil {
			m.Repository.Rollback()
			output <- ResultUseCase{Error: reject.Error}
			return
		}
		resetedDocument := model.B2CMerchantDocumentData{
			MerchantID:   merchantID,
			EditorIP:     userAttribute.UserIP,
			EditorID:     userAttribute.UserID,
			LastModified: null.TimeFrom(time.Now()),
		}
		resetAction := <-m.MerchantDocumentRepo.ResetRejectedDocument(ctxReq, resetedDocument)
		if resetAction.Error != nil {
			m.Repository.Rollback()
			output <- ResultUseCase{Error: resetAction.Error}
			return
		}
		m.Repository.Commit()

		memberName := merchantOnDB.MerchantName
		emailResult := <-m.MemberQueryRead.FindByID(ctxReq, userAttribute.UserID)
		member := memberModel.Member{}
		if emailResult.Result != nil {
			member = emailResult.Result.(memberModel.Member)

		}
		plQueue := model.MerchantPayloadEmail{
			MemberName:   memberName,
			Data:         oldData, // to get upgradeStatus, pass old data
			ReasonReject: reasonReject,
			AdminCMS:     member.FirstName + " " + member.LastName,
		}
		plLog := model.MerchantLog{
			Before: oldData,
			After:  merchantOnDB,
		}

		go m.QueuePublisher.QueueJob(ctxReq, plQueue, merchantOnDB.ID, "SendEmailMerchantRejectUpgrade")
		go m.QueuePublisher.QueueJob(ctxReq, plQueue, merchantOnDB.ID, "SendEmailAdmin")
		go m.QueuePublisher.QueueJob(ctxReq, plLog, merchantOnDB.ID, "InsertLogMerchantUpdate")
		go func() {
			m.PublishToKafkaMerchant(ctxReq, merchantOnDB, helper.EventProduceUpdateMerchant)
		}()

		output <- ResultUseCase{Result: merchantOnDB}

	})
	return output
}
func checkUpgradeStatusBeforeReject(merchantOnDB model.B2CMerchantDataV2) model.B2CMerchantDataV2 {
	if merchantOnDB.UpgradeStatus.String == "PENDING_MANAGE" {
		merchantOnDB.UpgradeStatus = zero.StringFrom("REJECT_MANAGE")
	} else if merchantOnDB.UpgradeStatus.String == "PENDING_ASSOCIATE" {
		merchantOnDB.UpgradeStatus = zero.StringFrom("REJECT_ASSOCIATE")
	}
	return merchantOnDB
}

func validateStatusIsActive(payload *model.B2CMerchantCreateInput) (string, bool) {
	if payload.Status == "" {
		switch payload.IsActive {
		case true:
			payload.Status = model.ActiveString
		default:
			payload.Status = model.NewString
		}
	} else {
		switch payload.Status {
		case model.ActiveString:
			payload.IsActive = true
		default:
			payload.IsActive = false
		}
	}
	return payload.Status, payload.IsActive
}

func validateExistingStatusIsActive(oldData model.B2CMerchantDataV2, existingMerchant model.B2CMerchantDataV2) (string, bool) {
	if existingMerchant.Status == model.NewString {
		switch existingMerchant.IsActive {
		case true:
			existingMerchant.Status = model.ActiveString
		default:
			existingMerchant.Status = oldData.Status
			if existingMerchant.Status == model.ActiveString {
				existingMerchant.Status = model.InactiveString
			}
		}
	} else {
		switch existingMerchant.Status {
		case model.ActiveString:
			existingMerchant.IsActive = true
		default:
			existingMerchant.IsActive = false
		}
	}
	return existingMerchant.Status, existingMerchant.IsActive
}

func (m *MerchantUseCaseImpl) AddMerchantPIC(ctxReq context.Context, payload *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase {
	ctx := "MerchantUseCase-AddMerchantPIC"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		merchant := m.MerchantRepo.LoadMerchant(ctxReq, payload.ID, private)
		if merchant.Error != nil && merchant.Error == sql.ErrNoRows {
			output <- ResultUseCase{Error: errMerchantNotFound, HTTPStatus: http.StatusNotFound}
			return
		}

		existingMerchant := merchant.Result.(model.B2CMerchantDataV2)
		oldData := existingMerchant
		payload.MerchantEmail = existingMerchant.MerchantEmail.String
		if payload.SellerOfficerEmail != "" {
			checkEmailPayload := payload.IsBhinnekaEmail()
			if !checkEmailPayload {
				output <- ResultUseCase{Error: errors.New("Seller Officer Email (PIC) is Not Bhinneka"), HTTPStatus: http.StatusBadRequest}
				return
			}
		}
		m.GetCorporateData(&existingMerchant, payload)
		existingMerchant.EditorID = null.StringFrom(userAttribute.UserID)
		existingMerchant.EditorIP = null.StringFrom(userAttribute.UserIP)
		existingMerchant.LastModified = null.TimeFrom(time.Now())
		existingMerchant.SellerOfficerName = zero.StringFrom(payload.SellerOfficerName)
		existingMerchant.SellerOfficerEmail = zero.StringFrom(payload.SellerOfficerEmail)

		m.Repository.StartTransaction()
		addOwnerResult := <-m.MerchantRepo.AddUpdateMerchant(ctxReq, existingMerchant)
		if addOwnerResult.Error != nil {
			output <- ResultUseCase{Error: errors.New("failed to add merchant owner"), HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}
		emailResult := <-m.MemberQueryRead.FindByID(ctxReq, userAttribute.UserID)
		member := memberModel.Member{}
		if emailResult.Result != nil {
			member = emailResult.Result.(memberModel.Member)

		}

		m.Repository.Commit()
		plLog := model.MerchantLog{
			Before: oldData,
			After:  existingMerchant,
		}

		go m.QueuePublisher.QueueJob(ctxReq, plLog, payload.ID, "InsertLogMerchantPICUpdate")

		m.MerchantService.InsertLogMerchantPIC(ctxReq, oldData, existingMerchant, "UPDATE", "merchantPIC", member)
		output <- ResultUseCase{Result: existingMerchant}
	})
	return output
}
