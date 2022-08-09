package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/gosimple/slug"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"
)

const (
	producerSturgeon        = "sturgeon"
	errUpgradeAlreadyActive = "merchant upgrade status already active"
	errUpgradePending       = "merchant upgrade status already pending"
)

// GetCorporateData function for validate merchant bank
func (m *MerchantUseCaseImpl) GetCorporateData(merchant *model.B2CMerchantDataV2, data *model.B2CMerchantCreateInput) {
	if data.CompanyName != "" {
		merchant.CompanyName = zero.StringFrom(data.CompanyName)
	}
	if data.MerchantAddress != "" {
		merchant.MerchantAddress = zero.StringFrom(data.MerchantAddress)
	}
}

// AddMerchant function for self-register merchant
func (m *MerchantUseCaseImpl) AddMerchant(ctxReq context.Context, data *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase {
	ctx := "MerchantUseCase-AddMerchant"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		merchant := model.B2CMerchantDataV2{}
		data.Status = checkStatus(data.Status)

		tags[helper.TextEmail] = data.MerchantEmail
		if err := m.validateMerchantField(ctxReq, data); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// Validate merchant bank id
		if err := m.ValidateMerchantBank(ctxReq, data, &merchant); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// Generate ID
		t := time.Now()
		data.ID = "MCH" + t.Format(helper.FormatYmdhis)
		data.Source = "cf"

		// Make slugify
		data.VanityURL = slug.MakeLang(data.MerchantName, "en")

		// Get Member Data from userID
		getMemberByID := <-m.MemberRepoRead.Load(ctxReq, data.UserID)
		if getMemberByID.Result == nil {
			err := fmt.Errorf("userID doesn't exist")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		memberData := getMemberByID.Result.(memberModel.Member)
		data.MerchantEmail = memberData.Email
		// Validate merchant data
		if err := m.ValidateMerchantData(ctxReq, data); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		// set merchant data from params
		merchant.SetMerchantData(data)
		merchant.ID = data.ID
		merchant.CreatorID = null.StringFrom(userAttribute.UserID)
		merchant.CreatorIP = null.StringFrom(userAttribute.UserIP)
		merchant.EditorID = null.StringFrom(userAttribute.UserID)
		merchant.EditorIP = null.StringFrom(userAttribute.UserIP)
		merchant.Created = zero.TimeFrom(time.Now())
		merchant.LastModified = null.TimeFrom(time.Now())
		merchant.Version = zero.IntFrom(merchant.Version.Int64 + 1)
		merchant.MerchantTypeString = zero.StringFrom(model.RegularString)
		merchant.MerchantType = model.StringToMerchantType(model.RegularString)
		merchant.GenderPic = data.GenderPic
		merchant.MerchantGroup = zero.StringFrom(data.MerchantGroup)
		merchant.GenderPicString = zero.StringFrom(data.GenderPic.String())
		merchant.MerchantEmail = zero.StringFrom(memberData.Email)
		merchant.LegalEntity = zero.IntFrom(int64(data.LegalEntity))
		merchant.NumberOfEmployee = zero.IntFrom(int64(data.NumberOfEmployee))
		merchant.Status = data.Status
		merchant.CountUpdateNameAvailable = 1

		m.Repository.StartTransaction()
		saveResult := <-m.MerchantRepo.AddUpdateMerchant(ctxReq, merchant)
		if saveResult.Error != nil {
			err := errors.New("failed to save merchant")
			tracer.Log(ctxReq, helper.TextExecUsecase, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}

		// set maps data form params
		maps := model.Maps{}
		maps.ID = strings.ReplaceAll("MAPS"+t.Format(helper.FormatYmdhisz), ".", "")
		maps.RelationID = data.ID
		maps.RelationName = "b2c_merchant"
		maps.Label = data.Maps.Label
		maps.Latitude = data.Maps.Latitude
		maps.Longitude = data.Maps.Longitude

		saveResultMaps := <-m.MerchantAddressRepo.AddUpdateAddressMaps(ctxReq, maps)
		if saveResultMaps.Error != nil {
			output <- ResultUseCase{Error: saveResultMaps.Error, HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}
		merchant.Maps = maps
		merchant.IsMapAvailable = helper.ValidateLatLong(maps.Latitude, maps.Longitude)

		documentProccess := <-m.MerchantDocumentsProcess(ctxReq, data.Documents, data.ID, userAttribute)
		if documentProccess.Error != nil {
			output <- ResultUseCase{Error: documentProccess.Error, HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}

		merchant.Documents = documentProccess.Result.([]model.B2CMerchantDocumentData)
		memberName := memberData.FirstName + " " + memberData.LastName

		plQueue := model.MerchantPayloadEmail{
			MemberName: memberName,
			Data:       merchant,
		}
		plLog := model.MerchantLog{
			Before: model.B2CMerchantDataV2{},
			After:  merchant,
		}

		go m.QueuePublisher.QueueJob(ctxReq, plQueue, merchant.ID, "SendEmailMerchantAdd")
		go m.QueuePublisher.QueueJob(ctxReq, plLog, merchant.ID, "InsertLogMerchantCreate")
		go m.MerchantService.PublishToKafkaUserMerchant(ctxReq, &merchant, helper.EventProduceCreateMerchant, producerSturgeon)

		m.Repository.Commit()

		tags[helper.TextResponse] = merchant
		output <- ResultUseCase{Result: merchant}
	})

	return output
}

// UpgradeMerchant function for upgrade merchant
func (m *MerchantUseCaseImpl) UpgradeMerchant(ctxReq context.Context, data *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase {
	ctx := "MerchantUseCase-UpgradeMerchant"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		data, err := m.validateMerchantFieldUpgrade(data)
		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// Get Member Data from userID
		getMemberByID := <-m.MemberRepoRead.Load(ctxReq, data.UserID)
		if getMemberByID.Result == nil {
			err := fmt.Errorf("UserID doesn't exist")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		memberData := getMemberByID.Result.(memberModel.Member)
		tags[helper.TextEmail] = memberData.Email

		// Load merchant data by id and user id
		merchantData := m.MerchantRepo.FindMerchantByID(ctxReq, data.ID, memberData.ID)
		if merchantData.Result == nil {
			tags[helper.TextResponse] = merchantData.Error
			output <- ResultUseCase{Error: merchantData.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		merchant := merchantData.Result.(model.B2CMerchantDataV2)
		if !merchant.IsActive {
			output <- ResultUseCase{Error: errors.New("your merchant is not active"), HTTPStatus: http.StatusBadRequest}
			return
		}
		if err := CheckUpgradeStatusBeforeSelfUpgrade(merchant); err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		oldMerchant := merchant
		m.GetCorporateData(&merchant, data)
		merchant.EditorID = null.StringFrom(userAttribute.UserID)
		merchant.EditorIP = null.StringFrom(userAttribute.UserIP)
		merchant.LastModified = null.TimeFrom(time.Now())
		merchant.Version = zero.IntFrom(merchant.Version.Int64 + 1)
		merchant.GenderPic = data.GenderPic
		merchant.MerchantGroup = zero.StringFrom(data.MerchantGroup)
		merchant.UpgradeStatus = zero.StringFrom(data.UpgradeStatus)
		merchant.GenderPicString = zero.StringFrom(data.GenderPic.String())
		merchant.ProductType = zero.StringFrom(data.ProductType)
		merchant.MerchantTypeString = zero.StringFrom(data.MerchantTypeString)
		merchant.LegalEntity = zero.IntFrom(int64(data.LegalEntity))
		merchant.NumberOfEmployee = zero.IntFrom(int64(data.NumberOfEmployee))

		m.Repository.StartTransaction()
		saveResult := <-m.MerchantRepo.AddUpdateMerchant(ctxReq, merchant)
		if saveResult.Error != nil {
			tracer.Log(ctxReq, helper.TextExecUsecase, saveResult.Error)
			output <- ResultUseCase{Error: errors.New("failed to upgrade merchant"), HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}

		documentProccess := <-m.MerchantDocumentsProcess(ctxReq, data.Documents, data.ID, userAttribute)
		if documentProccess.Error != nil {
			output <- ResultUseCase{Error: documentProccess.Error, HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}

		merchant.Documents = documentProccess.Result.([]model.B2CMerchantDocumentData)

		plLog := model.MerchantLog{
			Before: oldMerchant,
			After:  merchant,
		}
		plQueue := model.MerchantPayloadEmail{
			MemberName: memberData.FirstName + " " + memberData.LastName,
			Data:       merchant,
		}

		go m.QueuePublisher.QueueJob(ctxReq, plQueue, merchant.ID, "SendEmailMerchantUpgrade")
		go m.QueuePublisher.QueueJob(ctxReq, plLog, merchant.ID, "InsertLogMerchantUpdate")
		go func() {
			m.PublishToKafkaMerchant(ctxReq, merchant, helper.EventProduceUpdateMerchant)
		}()

		m.Repository.Commit()

		tags[helper.TextResponse] = merchant
		output <- ResultUseCase{Result: merchant}
	})

	return output
}
func CheckUpgradeStatusBeforeSelfUpgrade(merchant model.B2CMerchantDataV2) error {
	if merchant.UpgradeStatus.String == model.ActiveString {
		return errors.New(errUpgradeAlreadyActive)
	} else if merchant.UpgradeStatus.String == model.PendingAssociateString || merchant.UpgradeStatus.String == model.PendingManageString {
		return errors.New(errUpgradePending)
	}
	return nil
}
func (m *MerchantUseCaseImpl) InsertLogMerchant(ctxReq context.Context, old model.B2CMerchantDataV2, new model.B2CMerchantDataV2, action string) error {
	return m.MerchantService.InsertLogMerchant(ctxReq, old, new, action, model.Module)
}

// SelfUpdateMerchant update merchant from CF
func (m *MerchantUseCaseImpl) SelfUpdateMerchant(ctxReq context.Context, input *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase {
	ctx := "MerchantUseCase-SelfUpdateMerchant"
	output := make(chan ResultUseCase)
	var params serviceModel.SendbirdRequestV4

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		input := checkIsActiveStatus(input)

		merchantResult := m.MerchantRepo.FindMerchantByUser(ctxReq, userAttribute.UserID)
		if merchantResult.Error != nil {
			tags[helper.TextResponse] = merchantResult.Error.Error()
			output <- ResultUseCase{Error: merchantResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		currentData := merchantResult.Result.(model.B2CMerchantDataV2)
		oldData := currentData
		if err := m.ValidateMerchantBank(ctxReq, input, &currentData); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		if err := m.validateRequiredDataSelf(ctxReq, input, oldData); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		input.MerchantEmail = currentData.MerchantEmail.String
		m.GetCorporateData(&currentData, input)

		currentData.SetMerchantData(input)
		currentData.DailyOperationalStaff = zero.StringFrom(input.DailyOperationalStaff)
		currentData.EditorID = null.StringFrom(userAttribute.UserID)
		currentData.EditorIP = null.StringFrom(userAttribute.UserIP)
		currentData.LastModified = null.TimeFrom(time.Now())
		currentData.Version = zero.IntFrom(currentData.Version.ValueOrZero() + 1)
		currentData.GenderPic = input.GenderPic
		currentData.MerchantGroup = zero.StringFrom(input.MerchantGroup)
		currentData.GenderPicString = zero.StringFrom(input.GenderPic.String())
		currentData.ProductType = zero.StringFrom(input.ProductType)
		currentData.LegalEntity = zero.IntFrom(int64(input.LegalEntity))
		currentData.NumberOfEmployee = zero.IntFrom(int64(input.NumberOfEmployee))
		currentData.MerchantTypeString = zero.StringFrom(input.MerchantTypeString)
		currentData.MerchantType = model.StringToMerchantType(input.MerchantTypeString)

		currentData = checkIsActiveStatusV2(currentData, oldData)

		//Store Address
		currentData.StoreAddress = zero.StringFrom(input.StoreAddress)
		currentData.StoreVillageID = zero.StringFrom(input.StoreVillageID)
		currentData.StoreVillageID = zero.StringFrom(input.StoreVillageID)
		currentData.StoreDistrictID = zero.StringFrom(input.StoreDistrictID)
		currentData.StoreCityID = zero.StringFrom(input.StoreCityID)
		currentData.StoreProvinceID = zero.StringFrom(input.StoreProvinceID)
		currentData.StoreVillage = zero.StringFrom(input.StoreVillage)
		currentData.StoreDistrict = zero.StringFrom(input.StoreDistrict)
		currentData.StoreCity = zero.StringFrom(input.StoreCity)
		currentData.StoreProvince = zero.StringFrom(input.StoreProvince)
		currentData.StoreZipCode = zero.StringFrom(input.StoreZipCode)

		// set disallow field
		currentData.VanityURL = oldData.VanityURL
		currentData.MerchantName = oldData.MerchantName
		currentData.MerchantEmail = oldData.MerchantEmail

		m.Repository.StartTransaction()
		updateResult := <-m.MerchantRepo.AddUpdateMerchant(ctxReq, currentData)
		if updateResult.Error != nil {
			output <- ResultUseCase{Error: errors.New(model.MerchantFailedUpdateError), HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}

		// set maps data form params
		maps := model.Maps{}
		t := time.Now()
		maps.ID = strings.ReplaceAll("MAPS"+t.Format(helper.FormatYmdhisz), ".", "")
		maps.RelationID = currentData.ID
		maps.RelationName = "b2c_merchant"
		maps.Label = input.Maps.Label
		maps.Latitude = input.Maps.Latitude
		maps.Longitude = input.Maps.Longitude

		saveResultMaps := <-m.MerchantAddressRepo.AddUpdateAddressMaps(ctxReq, maps)
		if saveResultMaps.Error != nil {
			output <- ResultUseCase{Error: saveResultMaps.Error, HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}
		currentData.Maps = maps
		currentData.IsMapAvailable = helper.ValidateLatLong(maps.Latitude, maps.Longitude)
		currentData.Maps.ID = ""

		if len(input.Documents) > 0 {
			documentProccess := <-m.MerchantDocumentsProcess(ctxReq, input.Documents, currentData.ID, userAttribute)
			if documentProccess.Error != nil {
				output <- ResultUseCase{Error: documentProccess.Error, HTTPStatus: http.StatusBadRequest}
				m.Repository.Rollback()
				return
			}
			currentData.Documents = documentProccess.Result.([]model.B2CMerchantDocumentData)
		}

		params.UserID = currentData.ID
		params.NickName = currentData.MerchantName
		params.ProfileURL = currentData.MerchantLogo.String
		m.SendbirdService.UpdateUserSendbirdV4(ctxReq, &params)

		plLog := model.MerchantLog{
			Before: oldData,
			After:  currentData,
		}

		go m.QueuePublisher.QueueJob(ctxReq, plLog, oldData.ID, "InsertLogMerchantUpdate")
		go func() {
			m.PublishToKafkaMerchant(ctxReq, currentData, helper.EventProduceUpdateMerchant)
		}()

		m.Repository.Commit()

		output <- ResultUseCase{Result: currentData}
	})
	return output
}

func checkStatus(data string) string {
	status := data
	if data == "" {
		status = model.NewString
	}
	return status
}

func checkIsActiveStatus(input *model.B2CMerchantCreateInput) *model.B2CMerchantCreateInput {
	if input.Status == "" {
		switch input.IsActive {
		case true:
			input.Status = model.ActiveString
		default:
			input.Status = model.NewString
		}
	} else {
		switch input.Status {
		case model.ActiveString:
			input.IsActive = true
		default:
			input.IsActive = false
		}
	}
	return input
}

func checkIsActiveStatusV2(currentData model.B2CMerchantDataV2, oldData model.B2CMerchantDataV2) model.B2CMerchantDataV2 {
	if currentData.Status == model.NewString {
		switch currentData.IsActive {
		case true:
			currentData.Status = model.ActiveString
		default:
			currentData.Status = oldData.Status
			if currentData.Status == model.ActiveString {
				currentData.Status = model.InactiveString
			}
		}
	} else {
		switch currentData.Status {
		case model.ActiveString:
			currentData.IsActive = true
		default:
			currentData.IsActive = false
		}
	}
	return currentData
}

func (m *MerchantUseCaseImpl) SelfUpdateMerchantPartial(ctxReq context.Context, input *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase {
	ctx := "MerchantUseCase-SelfUpdateMerchantPartial"
	output := make(chan ResultUseCase)
	var params serviceModel.SendbirdRequestV4

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		input := checkIsActiveStatus(input)
		if err := m.validateRequiredData(ctxReq, input); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		merchantResult := m.MerchantRepo.FindMerchantByUser(ctxReq, userAttribute.UserID)
		if merchantResult.Error != nil {
			tags[helper.TextResponse] = merchantResult.Error.Error()
			output <- ResultUseCase{Error: merchantResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}
		isAttachment := "false"

		currentData := merchantResult.Result.(model.B2CMerchantDataV2)
		currentData = m.adjustMerchantData(ctxReq, currentData, isAttachment)
		currentData = m.CheckMapsBefore(ctxReq, currentData)
		oldData := currentData

		if err := m.ValidateMerchantBank(ctxReq, input, &currentData); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		input = CheckEmptyMerchantInput(oldData, input)
		input.MerchantEmail = currentData.MerchantEmail.String
		m.GetCorporateData(&currentData, input)
		currentData.SetMerchantData(input)
		currentData.EditorID = null.StringFrom(userAttribute.UserID)
		currentData.EditorIP = null.StringFrom(userAttribute.UserIP)
		currentData.LastModified = null.TimeFrom(time.Now())
		currentData.Version = zero.IntFrom(currentData.Version.ValueOrZero() + 1)

		currentData.MerchantGroup = zero.StringFrom(input.MerchantGroup)
		currentData.GenderPicString = zero.StringFrom(input.GenderPic.String())
		currentData.GenderPic = input.GenderPic
		currentData.ProductType = zero.StringFrom(input.ProductType)
		currentData.LegalEntity = zero.IntFrom(int64(input.LegalEntity))
		currentData.NumberOfEmployee = zero.IntFrom(int64(input.NumberOfEmployee))
		currentData.MerchantTypeString = zero.StringFrom(input.MerchantTypeString)
		currentData.MerchantType = input.MerchantType

		//Store Address
		currentData.StoreAddress = zero.StringFrom(input.StoreAddress)
		currentData.StoreVillageID = zero.StringFrom(input.StoreVillageID)
		currentData.StoreVillageID = zero.StringFrom(input.StoreVillageID)
		currentData.StoreDistrictID = zero.StringFrom(input.StoreDistrictID)
		currentData.StoreCityID = zero.StringFrom(input.StoreCityID)
		currentData.StoreProvinceID = zero.StringFrom(input.StoreProvinceID)
		currentData.StoreVillage = zero.StringFrom(input.StoreVillage)
		currentData.StoreDistrict = zero.StringFrom(input.StoreDistrict)
		currentData.StoreCity = zero.StringFrom(input.StoreCity)
		currentData.StoreProvince = zero.StringFrom(input.StoreProvince)
		currentData.StoreZipCode = zero.StringFrom(input.StoreZipCode)
		currentData.Documents = oldData.Documents

		currentData = checkIsActiveStatusV2(currentData, oldData)
		currentData = CheckEmptyCurrentData(currentData, oldData)

		// set disallow field
		currentData.MerchantName = input.MerchantName
		newVanity := slug.MakeLang(input.MerchantName, "en")
		currentData.VanityURL = zero.StringFrom(newVanity)

		m.Repository.StartTransaction()
		updateResult := <-m.MerchantRepo.AddUpdateMerchant(ctxReq, currentData)
		if updateResult.Error != nil {
			output <- ResultUseCase{Error: errors.New(model.MerchantFailedUpdateError), HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}

		// set maps data form params
		maps := model.Maps{}
		t := time.Now()
		maps.ID = strings.ReplaceAll("MAPS"+t.Format(helper.FormatYmdhisz), ".", "")
		maps.RelationID = currentData.ID
		maps.RelationName = "b2c_merchant"
		maps.Label = input.Maps.Label
		maps.Latitude = input.Maps.Latitude
		maps.Longitude = input.Maps.Longitude

		saveResultMaps := <-m.MerchantAddressRepo.AddUpdateAddressMaps(ctxReq, maps)
		if saveResultMaps.Error != nil {
			output <- ResultUseCase{Error: saveResultMaps.Error, HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}
		currentData.Maps = maps
		currentData.IsMapAvailable = helper.ValidateLatLong(maps.Latitude, maps.Longitude)
		currentData.Maps.ID = ""

		if len(input.Documents) > 0 {
			documentProccess := <-m.MerchantDocumentsProcess(ctxReq, input.Documents, currentData.ID, userAttribute)
			if documentProccess.Error != nil {
				output <- ResultUseCase{Error: documentProccess.Error, HTTPStatus: http.StatusBadRequest}
				m.Repository.Rollback()
				return
			}
			currentData.Documents = documentProccess.Result.([]model.B2CMerchantDocumentData)
		}

		params.UserID = currentData.ID
		params.NickName = currentData.MerchantName
		params.ProfileURL = currentData.MerchantLogo.String
		m.SendbirdService.UpdateUserSendbirdV4(ctxReq, &params)

		plLog := model.MerchantLog{
			Before: oldData,
			After:  currentData,
		}

		go m.QueuePublisher.QueueJob(ctxReq, plLog, oldData.ID, "InsertLogMerchantUpdate")
		go func() {
			m.PublishToKafkaMerchant(ctxReq, currentData, helper.EventProduceUpdateMerchant)
		}()

		m.Repository.Commit()

		output <- ResultUseCase{Result: currentData}
	})
	return output
}

func (m *MerchantUseCaseImpl) CheckMapsBefore(ctxReq context.Context, currentData model.B2CMerchantDataV2) model.B2CMerchantDataV2 {
	mapsResult := <-m.MerchantAddressRepo.FindAddressMaps(ctxReq, currentData.ID, "b2c_merchant")
	if mapsResult.Error == nil {
		currentData.Maps, _ = mapsResult.Result.(model.Maps)
		currentData.IsMapAvailable = helper.ValidateLatLong(currentData.Maps.Latitude, currentData.Maps.Longitude)
	}
	return currentData
}

func CheckEmptyMerchantInput(oldData model.B2CMerchantDataV2, input *model.B2CMerchantCreateInput) *model.B2CMerchantCreateInput {
	if input.GenderPicString == "" {
		input.GenderPicString = oldData.GenderPicString.String
	}
	if input.MerchantGroup == "" {
		input.MerchantGroup = oldData.MerchantGroup.String
	}
	if input.MerchantTypeString == "" {
		input.MerchantTypeString = oldData.MerchantTypeString.String
	}
	CheckStoreAddressInputEmpty(oldData, input)
	CheckMapsInputAndProductEmpty(oldData, input)

	//Maps
	if input.NumberOfEmployee == 0 {
		input.NumberOfEmployee = int(oldData.NumberOfEmployee.Int64)
	}
	if input.VanityURL == "" {
		input.VanityURL = oldData.VanityURL.String
	}
	if input.MerchantName == "" {
		input.MerchantName = oldData.MerchantName
	}
	if input.MerchantEmail == "" {
		input.MerchantEmail = oldData.MerchantEmail.String
	}

	input.GenderPic = model.StringToGenderPic(input.GenderPicString)
	input.MerchantType = model.StringToMerchantType(input.MerchantTypeString)
	return input
}
func CheckMapsInputAndProductEmpty(oldData model.B2CMerchantDataV2, input *model.B2CMerchantCreateInput) *model.B2CMerchantCreateInput {
	if input.Maps.Label == "" {
		input.Maps.Label = oldData.Maps.Label
	}
	if input.Maps.Latitude == 0 {
		input.Maps.Latitude = oldData.Maps.Latitude
	}

	if input.Maps.Longitude == 0 {
		input.Maps.Longitude = oldData.Maps.Longitude
	}

	if input.ProductType == "" {
		input.ProductType = oldData.ProductType.String
	}
	if input.LegalEntity == 0 {
		input.LegalEntity = int(oldData.LegalEntity.Int64)
	}
	return input
}
func CheckStoreAddressInputEmpty(oldData model.B2CMerchantDataV2, input *model.B2CMerchantCreateInput) *model.B2CMerchantCreateInput {
	if input.StoreAddress == "" {
		input.StoreAddress = oldData.StoreAddress.String
	}
	if input.StoreVillage == "" {
		input.StoreVillage = oldData.StoreVillage.String
	}
	if input.StoreVillageID == "" {
		input.StoreVillageID = oldData.StoreVillageID.String
	}
	if input.StoreCity == "" {
		input.StoreCity = oldData.StoreCity.String
	}
	if input.StoreCityID == "" {
		input.StoreCityID = oldData.StoreCity.String
	}
	CheckStoreAddressInputEmpty2(oldData, input)
	return input
}
func CheckStoreAddressInputEmpty2(oldData model.B2CMerchantDataV2, input *model.B2CMerchantCreateInput) *model.B2CMerchantCreateInput {
	if input.StoreDistrict == "" {
		input.StoreDistrict = oldData.StoreDistrict.String
	}
	if input.StoreDistrictID == "" {
		input.StoreDistrictID = oldData.StoreDistrictID.String
	}
	if input.StoreProvince == "" {
		input.StoreProvince = oldData.StoreProvince.String
	}
	if input.StoreProvinceID == "" {
		input.StoreProvinceID = oldData.StoreProvinceID.String
	}
	if input.StoreZipCode == "" {
		input.StoreZipCode = oldData.StoreZipCode.String
	}
	return input
}

func CheckEmptyCurrentData(currentData, oldData model.B2CMerchantDataV2) model.B2CMerchantDataV2 {
	if currentData.StoreClosureDate == nil {
		currentData.StoreClosureDate = oldData.StoreClosureDate
	}
	if currentData.StoreReopenDate == nil {
		currentData.StoreReopenDate = oldData.StoreReopenDate
	}
	return currentData
}

func (m *MerchantUseCaseImpl) ChangeMerchantName(ctxReq context.Context, input *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase {
	ctx := "MerchantUseCase-ChangeMerchantName"
	output := make(chan ResultUseCase)
	var params serviceModel.SendbirdRequestV4
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		input := checkIsActiveStatus(input)

		if err := m.validateMerchantName(input.MerchantName); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		checkName := m.MerchantRepo.FindMerchantByName(ctxReq, input.MerchantName)
		if checkName.Result != nil {
			err := fmt.Errorf("Nama merchant sudah ada")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		merchantResult := m.MerchantRepo.FindMerchantByUser(ctxReq, userAttribute.UserID)
		if merchantResult.Error != nil {
			tags[helper.TextResponse] = merchantResult.Error.Error()
			output <- ResultUseCase{Error: merchantResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		currentData := merchantResult.Result.(model.B2CMerchantDataV2)
		currentData.GenderPic = model.StringToGenderPic(currentData.GenderPicString.String)

		if err := m.ValidateMerchantBank(ctxReq, input, &currentData); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		if !currentData.IsActive {
			err := fmt.Errorf("Merchant belum aktif, tidak bisa mengubah nama")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		// set disallow field
		if err := m.validateMerchantName(input.MerchantName); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		currentData.MerchantName = input.MerchantName
		newVanity := slug.MakeLang(input.MerchantName, "en")
		currentData.VanityURL = zero.StringFrom(newVanity)
		err, status, modelB2CMerchantDatav2 := m.CheckCount(ctxReq, currentData, input, userAttribute)

		params.UserID = currentData.ID
		params.NickName = currentData.MerchantName
		params.ProfileURL = currentData.MerchantLogo.String
		m.SendbirdService.UpdateUserSendbirdV4(ctxReq, &params)

		output <- ResultUseCase{Error: err, HTTPStatus: status, Result: modelB2CMerchantDatav2}
	})
	return output
}

func (m *MerchantUseCaseImpl) CheckCount(ctxReq context.Context, currentData model.B2CMerchantDataV2, input *model.B2CMerchantCreateInput, userAttribute *model.MerchantUserAttribute) (error, int, model.B2CMerchantDataV2) {
	oldData := currentData

	switch currentData.CountUpdateNameAvailable {
	case 1:
		currentData.CountUpdateNameAvailable = 0
		m.Repository.StartTransaction()
		updateResult := <-m.MerchantRepo.AddUpdateMerchant(ctxReq, currentData)
		if updateResult.Error != nil {
			m.Repository.Rollback()
			return errors.New(model.MerchantFailedUpdateError), http.StatusBadRequest, model.B2CMerchantDataV2{}
		}

		// set maps data form params
		maps := model.Maps{}
		t := time.Now()
		maps.ID = strings.ReplaceAll("MAPS"+t.Format(helper.FormatYmdhisz), ".", "")
		maps.RelationID = currentData.ID
		maps.RelationName = "b2c_merchant"
		maps.Label = input.Maps.Label
		maps.Latitude = input.Maps.Latitude
		maps.Longitude = input.Maps.Longitude

		saveResultMaps := <-m.MerchantAddressRepo.AddUpdateAddressMaps(ctxReq, maps)
		if saveResultMaps.Error != nil {
			m.Repository.Rollback()
			return saveResultMaps.Error, http.StatusBadRequest, model.B2CMerchantDataV2{}
		}
		currentData.Maps = maps
		currentData.IsMapAvailable = helper.ValidateLatLong(maps.Latitude, maps.Longitude)
		currentData.Maps.ID = ""

		if len(input.Documents) > 0 {
			documentProccess := <-m.MerchantDocumentsProcess(ctxReq, input.Documents, currentData.ID, userAttribute)
			if documentProccess.Error != nil {
				m.Repository.Rollback()
				return documentProccess.Error, http.StatusBadRequest, model.B2CMerchantDataV2{}
			}
			currentData.Documents = documentProccess.Result.([]model.B2CMerchantDocumentData)
		}

		plLog := model.MerchantLog{
			Before: oldData,
			After:  currentData,
		}
		if err := m.UpdateMerchantSendbirdV4(ctxReq, oldData, input); err != nil {
			m.Repository.Rollback()
			return err, http.StatusBadRequest, model.B2CMerchantDataV2{}
		}

		go m.QueuePublisher.QueueJob(ctxReq, plLog, oldData.ID, "InsertLogMerchantUpdate")
		go func() {
			m.PublishToKafkaMerchant(ctxReq, currentData, helper.EventProduceUpdateMerchant)
		}()

		m.Repository.Commit()

	case 0:
		err := fmt.Errorf("User sudah mencapai limit untuk update nama")

		return err, http.StatusBadRequest, model.B2CMerchantDataV2{}
	}
	return nil, 200, currentData
}

func (m *MerchantUseCaseImpl) ClearRejectUpgrade(ctxReq context.Context, memberId string, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase {
	ctx := "MerchantUseCase-ClearRejectUpgrade"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags[helper.TextMerchantIDCamel] = memberId
		merchantResult := m.MerchantRepo.FindMerchantByUser(ctxReq, userAttribute.UserID)
		if merchantResult.Error != nil {
			tags[helper.TextResponse] = merchantResult.Error.Error()
			output <- ResultUseCase{Error: merchantResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		currentData := merchantResult.Result.(model.B2CMerchantDataV2)
		oldData := currentData
		if currentData.UpgradeStatus.String == "" {
			output <- ResultUseCase{Error: errMerchantNotValidForClearUpgrade, HTTPStatus: http.StatusBadRequest}
			return
		}
		if currentData.UpgradeStatus.String == model.ActiveString {
			output <- ResultUseCase{Error: errUnableToRejectUpgrade, HTTPStatus: http.StatusBadRequest}
			return
		}
		if currentData.UpgradeStatus.String == model.PendingAssociateString || currentData.UpgradeStatus.String == model.PendingManageString {
			output <- ResultUseCase{Error: errClearPendingUpgrade, HTTPStatus: http.StatusBadRequest}
			return
		}

		currentData.EditorID = null.StringFrom(userAttribute.UserID)
		currentData.EditorIP = null.StringFrom(userAttribute.UserIP)
		currentData.LastModified = null.TimeFrom(time.Now())
		currentData.UpgradeStatus = zero.StringFrom("")
		currentData.Reason = zero.StringFrom("")

		m.Repository.StartTransaction()

		clearReject := <-m.MerchantRepo.ClearRejectUpgrade(ctxReq, currentData)
		if clearReject.Error != nil {
			m.Repository.Rollback()
			output <- ResultUseCase{Error: clearReject.Error}
			return
		}

		m.Repository.Commit()

		plLog := model.MerchantLog{
			Before: oldData,
			After:  currentData,
		}
		fmt.Println(userAttribute.UserID)

		go m.QueuePublisher.QueueJob(ctxReq, plLog, currentData.ID, "InsertLogMerchantUpdate")
		go func() {
			m.PublishToKafkaMerchant(ctxReq, currentData, helper.EventProduceUpdateMerchant)
		}()

		output <- ResultUseCase{Result: currentData}

	})
	return output
}
