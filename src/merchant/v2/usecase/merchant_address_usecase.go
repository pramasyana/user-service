package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Bhinneka/golib"
	stringLib "github.com/Bhinneka/golib/string"
	"github.com/Bhinneka/golib/tracer"
	localConfig "github.com/Bhinneka/user-service/config"
	"github.com/Bhinneka/user-service/helper"
	memberRepo "github.com/Bhinneka/user-service/src/member/v1/repo"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/Bhinneka/user-service/src/merchant/v2/repo"
	"github.com/Bhinneka/user-service/src/service"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	sharedRepo "github.com/Bhinneka/user-service/src/shared/repository"
)

// MerchantAddressUseCaseImpl data structure
type MerchantAddressUseCaseImpl struct {
	MerchantRepo        repo.MerchantRepository
	MerchantAddressRepo repo.MerchantAddressRepository
	MemberRepoRead      memberRepo.MemberRepository
	Repository          *sharedRepo.Repository
	BarracudaService    service.BarracudaServices
	ActivityService     service.ActivityServices
}

// NewMerchantAddressUseCase function for initialise merchant use case implementation mo el
func NewMerchantAddressUseCase(repository localConfig.ServiceRepository, services localConfig.ServiceShared) MerchantAddressUseCase {
	return &MerchantAddressUseCaseImpl{
		Repository:          repository.Repository,
		MerchantRepo:        repository.MerchantRepository,
		MerchantAddressRepo: repository.MerchantAddressRepository,
		MemberRepoRead:      repository.MemberRepository,
		BarracudaService:    services.BarracudaService,
		ActivityService:     services.ActivityService,
	}
}

// AddUpdateWarehouseAddress function for add or update address
func (m *MerchantAddressUseCaseImpl) AddUpdateWarehouseAddress(ctxReq context.Context, data model.WarehouseData, memberID, action string) <-chan ResultUseCase {
	ctx := "MerchantAddressUseCase-AddUpdateWarehouseAddress"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags[helper.TextParameter] = data

		// validate address data and return back with parsed data
		oldAddress, data, err := m.validateAddress(ctxReq, data, memberID)
		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// start transaction process
		m.Repository.StartTransaction()

		// process function for save warehouse address data
		data, err = m.saveAddressProccess(ctxReq, data, action, memberID)
		if err != nil {
			err := errors.New(msgErrorSave)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}

		// process function for save phone address data
		data, err = m.saveAddressPhoneProccess(ctxReq, data, action)
		if err != nil {
			err := errors.New(msgErrorSave)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}

		// commit process
		m.Repository.Commit()

		//insert log data
		go m.InsertLogAddressWarehouse(ctxReq, oldAddress, data, helper.TextUpdateUpper)

		tags[helper.TextResponse] = data
		output <- ResultUseCase{Result: data}
	})

	return output
}

// saveAddressProccess function for process save address
func (m *MerchantAddressUseCaseImpl) saveAddressProccess(ctxReq context.Context, data model.WarehouseData, action, memberID string) (model.WarehouseData, error) {
	ctx := "MerchantAddressUseCase-saveAddressProccess"

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := make(map[string]interface{})
	defer func() {
		tr.Finish(tags)
	}()

	var saveResult repo.ResultRepository

	// parse address structured
	address, err := m.parseAddress(ctxReq, data)
	if err != nil {
		return data, err
	}

	if action == helper.TextUpdate {
		address.ModifiedBy = memberID
	} else {
		address.CreatedBy = memberID
	}

	tags[helper.TextParameter+"_address"] = address
	// update address repository process to database
	saveResult = <-m.MerchantAddressRepo.AddUpdateAddress(ctxReq, address)
	if saveResult.Error != nil {
		err := errors.New(msgErrorSave)
		return data, err
	}

	resultData, ok := saveResult.Result.(model.AddressData)
	if !ok {
		err := errors.New(msgErrorSave)
		return data, err
	}
	data.ID = resultData.ID
	data.ModifiedBy = resultData.ModifiedBy
	data.Version = resultData.Version
	data.CreatedBy = resultData.CreatedBy
	data.IsPrimary = resultData.IsPrimary
	data.Type = resultData.Type

	data.ProvinceID = resultData.ProvinceID
	data.ProvinceName = resultData.ProvinceName
	data.CityID = resultData.CityID
	data.CityName = resultData.CityName
	data.DistrictID = resultData.DistrictID
	data.DistrictName = resultData.DistrictName
	data.SubDistrictID = resultData.SubDistrictID
	data.SubDistrictName = resultData.SubDistrictName
	data.PostalCode = resultData.PostalCode

	// set maps data form params
	maps := model.Maps{}
	t := time.Now()
	maps.ID = strings.ReplaceAll("MAPS"+t.Format(helper.FormatYmdhisz), ".", "")
	maps.RelationID = data.ID
	maps.RelationName = "address"
	maps.Label = data.Maps.Label
	maps.Latitude = data.Maps.Latitude
	maps.Longitude = data.Maps.Longitude

	data.Maps = maps
	data.IsMapAvailable = helper.ValidateLatLong(maps.Latitude, maps.Longitude)
	tags[helper.TextParameter+"_maps"] = maps
	saveResultMaps := <-m.MerchantAddressRepo.AddUpdateAddressMaps(tr.Context(), maps)
	if saveResultMaps.Error != nil {
		return data, saveResultMaps.Error
	}

	tags[helper.TextResponse] = data
	return data, nil

}

// saveAddressPhoneProccess function for process save phone address
func (m *MerchantAddressUseCaseImpl) saveAddressPhoneProccess(ctxReq context.Context, data model.WarehouseData, action string) (model.WarehouseData, error) {
	var saveResult repo.ResultRepository

	var phoneData = []model.PhoneData{
		{
			Number:    data.Mobile,
			Type:      model.MobileString,
			IsPrimary: true,
		},
		{
			Number: data.Phone,
			Type:   model.PhoneString,
		},
	}

	for _, phone := range phoneData {
		phone.RelationID = data.ID
		phone.RelationName = model.AddressString
		phone.Label = data.Name
		phone.Created = data.Created
		phone.LastModified = data.LastModified
		phone.CreatedBy = data.CreatedBy
		phone.ModifiedBy = data.ModifiedBy
		phone.Version = data.Version
		if action == helper.TextUpdate {
			// update address repository process to database
			saveResult = <-m.MerchantAddressRepo.UpdatePhoneAddress(ctxReq, phone)
		} else {
			// add address repository process to database
			saveResult = <-m.MerchantAddressRepo.AddPhoneAddress(ctxReq, phone)
		}
	}

	if saveResult.Error != nil {
		err := errors.New(msgErrorSave)
		return data, err
	}

	return data, nil
}

func (m *MerchantAddressUseCaseImpl) validateAddress(ctxReq context.Context, data model.WarehouseData, memberID string) (model.WarehouseData, model.WarehouseData, error) {
	oldData := model.WarehouseData{}

	// validate address information
	_, err := m.validateLocationAddress(data)
	if err != nil {
		return oldData, data, err
	}

	// get merchant data by member id
	getMerchantByUser := m.MerchantRepo.FindMerchantByUser(ctxReq, memberID)
	if getMerchantByUser.Result == nil {
		err := fmt.Errorf("Merchant doesn't exist")
		return oldData, data, err
	}

	merchant, _ := getMerchantByUser.Result.(model.B2CMerchantDataV2)

	data.MerchantID = merchant.ID
	data.Created = time.Now()

	var phoneData = []model.PhoneData{
		{
			Number: data.Mobile,
			Type:   model.MobileString,
		},
		{
			Number: data.Phone,
			Type:   model.PhoneString,
		},
	}

	for _, phone := range phoneData {
		phone.Label = data.Name
		// validate phone information
		_, err := m.validatePhone(phone)
		if err != nil {
			return oldData, data, err
		}
	}

	if data.ID != "" {
		// get merchant address from merchant id and address id
		getMerchantAddressByID := <-m.MerchantAddressRepo.FindMerchantAddress(ctxReq, data.ID)
		if getMerchantAddressByID.Result == nil {
			err := fmt.Errorf("Address doesn't exist")
			return oldData, data, err
		}

		detailFind, ok := getMerchantAddressByID.Result.(model.WarehouseData)
		if !ok {
			err := errors.New(msgErrorSave)
			return oldData, data, err
		}
		oldData = detailFind

		data.Version = detailFind.Version + 1
		data.Created = detailFind.Created
		data.LastModified = time.Now()
		data.LastModifiedString = data.LastModified.Format(time.RFC3339)
		data.CreatedBy = detailFind.CreatedBy
		data.IsPrimary = detailFind.IsPrimary
	}
	data.CreatedString = data.Created.Format(time.RFC3339)
	if oldData.IsPrimary && data.Status == helper.TextInactive {
		return oldData, data, errors.New("Cannot set primary address to inactive")
	}

	return oldData, data, nil
}

// validatePhone function for validating address data
func (m *MerchantAddressUseCaseImpl) validatePhone(data model.PhoneData) (model.PhoneData, error) {
	if len(data.Label) <= 0 {
		err := errors.New("name is required")
		return data, err
	}

	if !golib.ValidateAlphanumericWithSpace(data.Label, false) {
		err := errors.New("name only alphabet is allowed")
		return data, err
	}

	if data.Type == model.PhoneString && len(data.Number) > 0 && helper.ValidatePhoneNumberMaxInput(data.Number) != nil {
		err := errors.New("phone number is in bad format")
		return data, err
	}

	return data, nil
}

// validateLocationAddress function for validating address data location
func (m *MerchantAddressUseCaseImpl) validateLocationAddress(data model.WarehouseData) (model.WarehouseData, error) {
	if len(data.Address) <= 0 {
		err := errors.New("address is required")
		return data, err
	}

	if !stringLib.ValidateLatinOnlyExcepTag(data.Address) {
		err := errors.New("address only latin character")
		return data, err
	}

	if len(data.PostalCode) <= 0 {
		err := errors.New("postal code is required")
		return data, err
	}

	// validate area
	data, err := m.validateLocationAreaAddress(data)
	if err != nil {
		return data, err
	}
	return data, nil
}

// validateLocationAddress function for validating address data location area
func (m *MerchantAddressUseCaseImpl) validateLocationAreaAddress(data model.WarehouseData) (model.WarehouseData, error) {
	if !golib.ValidateNumeric(data.PostalCode) {
		err := errors.New("postal code only numeric is allowed")
		return data, err
	}

	if len(data.SubDistrictID) <= 0 {
		err := errors.New("subDistrict ID is required")
		return data, err
	}

	if !golib.ValidateNumeric(data.SubDistrictID) {
		err := errors.New("subdistrict ID only numeric is allowed")
		return data, err
	}

	if len(data.SubDistrictName) <= 0 {
		err := errors.New("subdistrict name is required")
		return data, err
	}

	if len(data.DistrictID) <= 0 {
		err := errors.New("district ID is required")
		return data, err
	}

	if !golib.ValidateNumeric(data.DistrictID) {
		err := errors.New("district ID only numeric is allowed")
		return data, err
	}

	if len(data.DistrictName) <= 0 {
		err := errors.New("district name is required")
		return data, err
	}

	if len(data.CityID) <= 0 {
		err := errors.New("city ID is required")
		return data, err
	}

	if !golib.ValidateNumeric(data.CityID) {
		err := errors.New("city ID only numeric is allowed")
		return data, err
	}

	if len(data.CityName) <= 0 {
		err := errors.New("city name is required")
		return data, err
	}

	if len(data.ProvinceID) <= 0 {
		err := errors.New("province ID is required")
		return data, err
	}

	if !golib.ValidateNumeric(data.ProvinceID) {
		err := errors.New("province ID only numeric is allowed")
		return data, err
	}

	if len(data.ProvinceName) <= 0 {
		err := errors.New("province name is required")
		return data, err
	}
	if !golib.StringInSlice(data.Status, []string{helper.TextActive, helper.TextInactive}) {
		return data, errors.New("status should be `ACTIVE` of `INACTIVE`")
	}

	return data, nil
}

// parseAddress function for validating merchant address data
func (m *MerchantAddressUseCaseImpl) parseAddress(ctxReq context.Context, data model.WarehouseData) (model.AddressData, error) {
	ctx := "MerchantAddressUseCase-parseAddress"

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := make(map[string]interface{})
	defer func() {
		tr.Finish(tags)
	}()

	address := model.AddressData{
		ID:                 data.ID,
		RelationID:         data.MerchantID,
		RelationName:       model.MerchantString,
		Type:               model.WarehouseString,
		Label:              data.Label,
		ProvinceID:         data.ProvinceID,
		ProvinceName:       data.ProvinceName,
		CityID:             data.CityID,
		CityName:           data.CityName,
		DistrictID:         data.DistrictID,
		DistrictName:       data.DistrictName,
		SubDistrictID:      data.SubDistrictID,
		SubDistrictName:    data.SubDistrictName,
		PostalCode:         data.PostalCode,
		Address:            data.Address,
		Version:            data.Version,
		Created:            data.Created,
		CreatedString:      data.CreatedString,
		LastModified:       data.LastModified,
		LastModifiedString: data.LastModifiedString,
		CreatedBy:          data.CreatedBy,
		IsPrimary:          data.IsPrimary,
		Status:             data.Status,
	}

	zipcodeParam := serviceModel.ZipCodeQueryParameter{
		SubDistrictID: data.SubDistrictID,
		ZipCode:       data.PostalCode,
	}
	tags[helper.TextParameter] = zipcodeParam
	// validate & find address information from barracuda service by zipcode
	barracudaService := <-m.BarracudaService.FindZipcode(ctxReq, zipcodeParam)
	if barracudaService.Error != nil {
		tags[helper.TextResponse] = barracudaService.Error
		return address, barracudaService.Error
	}

	// get detail zipcode
	detailZipCode, ok := barracudaService.Result.(serviceModel.ResponseZipCode)
	if !ok {
		err := errors.New("failed parse zipcode response")
		tags[helper.TextResponse] = err
		return address, err
	}

	if len(detailZipCode.Data) == 0 {
		err := errors.New("Data zipcode Not Found")
		tags[helper.TextResponse] = err
		return address, err
	}

	// restruct the data with valid information
	for _, c := range detailZipCode.Data {
		address.ProvinceID = c.Province.ProvinceID
		address.ProvinceName = c.Province.ProvinceName
		address.CityID = c.City.CityID
		address.CityName = c.City.CityName
		address.DistrictID = c.District.DistrictID
		address.DistrictName = c.District.DistrictName
		address.SubDistrictID = c.SubDistrict.SubDistrictID
		address.SubDistrictName = c.SubDistrict.SubDistrictName
		address.PostalCode = strconv.Itoa(c.ZipCode)
	}

	if address.ID == "" {
		// generate new address id format
		t := time.Now()
		address.ID = "ADDRS" + t.Format(helper.FormatYmdhisz)
		address.ID = strings.Replace(address.ID, ".", "", -1)
	}

	countAddress := <-m.MerchantAddressRepo.CountAddress(ctxReq, address.RelationID, address.RelationName, &model.ParameterWarehouse{})
	if countAddress.Error != nil {
		err := errors.New("count address failed")
		return address, err
	}

	total, _ := countAddress.Result.(int)
	if total == 0 {
		// set default address as primary
		address.IsPrimary = true
	}

	tags[helper.TextResponse] = address
	return address, nil
}

// UpdatePrimaryWarehouseAddress function for add new address
func (m *MerchantAddressUseCaseImpl) UpdatePrimaryWarehouseAddress(ctxReq context.Context, params model.ParameterPrimaryWarehouse) <-chan ResultUseCase {
	ctx := "MerchantAddressUseCase-UpdatePrimaryWarehouseAddress"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		// find address by id
		findResult := <-m.FindWarehouseAddress(ctxReq, params.AddressID, params.MemberID)
		if findResult.Error != nil {
			tags[helper.TextResponse] = findResult.Error
			output <- ResultUseCase{Error: findResult.Error, HTTPStatus: findResult.HTTPStatus}
			return
		}

		addressDetail, ok := findResult.Result.(model.WarehouseData)
		if !ok {
			err := errors.New(msgErrorFindAddress)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// set old address data
		oldAddress := addressDetail

		m.Repository.StartTransaction()

		// update primary address
		result := <-m.MerchantAddressRepo.UpdatePrimaryAddressByRelationID(ctxReq, addressDetail.MerchantID, model.MerchantString)
		if result.Error != nil {
			err := errors.New(msgErrorUpdatePrimary)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}

		addressDetail.IsPrimary = true
		addressDetail.ModifiedBy = params.MemberID
		addressData := model.RestructToAddress(addressDetail)

		// update address
		resultUpdate := <-m.MerchantAddressRepo.AddUpdateAddress(ctxReq, addressData)
		if resultUpdate.Error != nil {
			err := errors.New(msgErrorUpdatePrimary)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}

		m.Repository.Commit()

		//insert log data
		go m.InsertLogAddressWarehouse(ctxReq, oldAddress, addressDetail, helper.TextUpdateUpper)

		output <- ResultUseCase{Result: nil}
	})
	return output
}

// InsertLogAddressWarehouse function to write log activity service for merchant warehouse
func (m *MerchantAddressUseCaseImpl) InsertLogAddressWarehouse(ctxReq context.Context, oldData, newData model.WarehouseData, action string) error {

	targetID := newData.ID
	if targetID == "" {
		targetID = oldData.ID
	}

	if oldData.ID == "" {
		action = helper.TextInsertUpper
		oldData = model.WarehouseData{}
	}

	payload := serviceModel.Payload{
		Module:    model.ModuleWarehouse,
		Action:    action,
		Target:    targetID,
		CreatorID: newData.CreatedBy,
		EditorID:  newData.ModifiedBy,
	}

	m.ActivityService.InsertLog(ctxReq, oldData, newData, payload)
	return nil
}

func (m *MerchantAddressUseCaseImpl) validateQueryParams(params *model.ParameterWarehouse) error {
	if params.OrderBy != "" && !golib.StringInSlice(params.OrderBy, []string{"id", "created", "lastModified"}, false) {
		return errors.New("invalid field ordering")
	}
	if params.Sort != "" && !golib.StringInSlice(params.Sort, []string{"desc", "asc"}, false) {
		return errors.New("invalid sort order")
	}
	if params.ShowAll != "" && !golib.StringInSlice(params.ShowAll, []string{"true", "false"}, false) {
		return errors.New("invalid show all parameter")
	}

	paging, err := helper.ValidatePagination(
		helper.PaginationParameters{
			Page:     1, // default
			StrPage:  params.StrPage,
			Limit:    5, // default
			StrLimit: params.StrLimit,
		})

	if err != nil {
		return err
	}
	if paging.Limit > 10 {
		return errors.New("max 10 result per page")
	}
	params.Page = paging.Page
	if params.ShowAll != "" && params.ShowAll == "true" {
		params.Limit = 9999
		params.StrLimit = ""
	} else {
		params.Limit = paging.Limit
	}

	params.Offset = paging.Offset
	return nil
}

// GetWarehouseAddresses function for getting list of warehouse address
func (m *MerchantAddressUseCaseImpl) GetWarehouseAddresses(ctxReq context.Context, params *model.ParameterWarehouse) <-chan ResultUseCase {
	ctx := "MerchantAddressUseCase-GetWarehouseAddresses"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		var (
			merchant model.B2CMerchantDataV2
		)
		if err := m.validateQueryParams(params); err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// get merchant data by member id
		if err := m.loadMerchant(ctxReq, params, &merchant); err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		params.MerchantID = merchant.ID
		tags[helper.TextParameter] = params

		warehouseResult := <-m.MerchantAddressRepo.GetListAddress(ctxReq, params)
		if warehouseResult.Error != nil {
			httpStatus := http.StatusInternalServerError

			// when data is not found
			if warehouseResult.Error == sql.ErrNoRows {
				httpStatus = http.StatusNotFound
				warehouseResult.Error = fmt.Errorf(helper.ErrorDataNotFound, "merchant_address")
			}

			output <- ResultUseCase{Error: warehouseResult.Error, HTTPStatus: httpStatus}
			return
		}

		warehouse, ok := warehouseResult.Result.(model.ListWarehouse)
		if !ok {
			output <- ResultUseCase{Error: errors.New("invalid result set"), HTTPStatus: http.StatusBadRequest}
		}

		totalResult := <-m.MerchantAddressRepo.CountAddress(ctxReq, params.MerchantID, model.MerchantString, params)
		if totalResult.Error != nil {
			output <- ResultUseCase{Error: totalResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		warehouse.TotalData = totalResult.Result.(int)

		output <- ResultUseCase{Result: warehouse}
	})

	return output
}

func (m *MerchantAddressUseCaseImpl) loadMerchant(ctxReq context.Context, params *model.ParameterWarehouse, mch *model.B2CMerchantDataV2) error {
	// get merchant data by member id
	var merchant model.B2CMerchantDataV2
	privacy := "private"
	if params.MemberID != "" {
		getMerchantByUser := m.MerchantRepo.FindMerchantByUser(ctxReq, params.MemberID)
		if getMerchantByUser.Result == nil {
			return fmt.Errorf(msgErrorFindAddress)
		}
		merchant, _ = getMerchantByUser.Result.(model.B2CMerchantDataV2)
		mch.ID = merchant.ID
	} else if params.MerchantID != "" {
		getMerchantByID := m.MerchantRepo.LoadMerchant(ctxReq, params.MerchantID, privacy)
		if getMerchantByID.Result == nil {
			return fmt.Errorf(msgErrorFindAddress)
		}
		merchant, _ = getMerchantByID.Result.(model.B2CMerchantDataV2)
		mch.ID = merchant.ID
	}

	return nil
}

// GetDetailWarehouseAddress function for getting detail of merchant address
func (m *MerchantAddressUseCaseImpl) GetDetailWarehouseAddress(ctxReq context.Context, addressID string, memberID string) <-chan ResultUseCase {
	ctx := "MerchantAddressUseCase-GetDetailWarehouseAddress"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		// find address by id
		findResult := <-m.FindWarehouseAddress(ctxReq, addressID, memberID)
		if findResult.Error != nil {
			tags[helper.TextResponse] = findResult.Error
			output <- ResultUseCase{Error: findResult.Error, HTTPStatus: findResult.HTTPStatus}
			return
		}

		addressDetail, ok := findResult.Result.(model.WarehouseData)
		if !ok {
			err := errors.New(msgErrorFindAddress)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: addressDetail}

	})
	return output
}

// DeleteWarehouseAddress function for add new address
func (m *MerchantAddressUseCaseImpl) DeleteWarehouseAddress(ctxReq context.Context, addressID, memberID string) <-chan ResultUseCase {
	ctx := "MerchantAddressUseCase-DeleteWarehouseAddress"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		// find address by id
		findResult := <-m.FindWarehouseAddress(ctxReq, addressID, memberID)
		if findResult.Error != nil {
			tags[helper.TextResponse] = findResult.Error
			output <- ResultUseCase{Error: findResult.Error, HTTPStatus: findResult.HTTPStatus}
			return
		}

		addressDetail, ok := findResult.Result.(model.WarehouseData)
		if !ok {
			err := errors.New(msgErrorFindAddress)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// check if primary address
		if addressDetail.IsPrimary {
			err := errors.New("cannot delete primary address")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		m.Repository.StartTransaction()
		// delete warehouse address from database
		result := <-m.MerchantAddressRepo.DeleteWarehouseAddress(ctxReq, addressID)
		if result.Error != nil {
			err := errors.New(msgErrorDeleteAddress)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}

		// delete phone address from database
		resultPhone := <-m.MerchantAddressRepo.DeletePhoneAddress(ctxReq, addressID, model.AddressString)
		if resultPhone.Error != nil {
			err := errors.New(msgErrorDeleteAddress)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			m.Repository.Rollback()
			return
		}
		m.Repository.Commit()

		//insert log data
		go m.InsertLogAddressWarehouse(ctxReq, addressDetail, model.WarehouseData{}, helper.TextDeleteUpper)

		output <- ResultUseCase{Result: nil}
	})
	return output
}

// FindWarehouseAddress function for find address
func (m *MerchantAddressUseCaseImpl) FindWarehouseAddress(ctxReq context.Context, addressID string, memberID string) <-chan ResultUseCase {
	ctx := "MerchantAddressUseCase-FindWarehouseAddress"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		// find address by id
		findResult := <-m.MerchantAddressRepo.FindMerchantAddress(ctxReq, addressID)
		if findResult.Error != nil {
			err := errors.New(msgErrorFindAddress)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusNotFound}
			return
		}

		addressDetail, ok := findResult.Result.(model.WarehouseData)
		if !ok {
			err := errors.New(msgErrorFindAddress)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// get merchant data by member id & merchant id
		getMerchantByUser := m.MerchantRepo.FindMerchantByID(ctxReq, addressDetail.MerchantID, memberID)
		if getMerchantByUser.Result == nil {
			err := fmt.Errorf(msgErrorFindAddress)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusNotFound}
			return
		}

		tags[helper.TextResponse] = addressDetail
		output <- ResultUseCase{Result: addressDetail}
	})

	return output
}

// GetWarehouseAddressByID function for getting detail of merchant address
func (m *MerchantAddressUseCaseImpl) GetWarehouseAddressByID(ctxReq context.Context, merchantID, addressID string) <-chan ResultUseCase {
	ctx := "MerchantAddressUseCase-GetDetailWarehouseAddress"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		// find address by addressID
		findResult := <-m.MerchantAddressRepo.FindMerchantAddress(ctxReq, addressID)
		if findResult.Error != nil {
			err := errors.New(msgErrorFindAddress)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusNotFound}
			return
		}

		addressDetail, ok := findResult.Result.(model.WarehouseData)
		if !ok {
			err := errors.New(msgErrorFindAddress)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		if addressDetail.MerchantID != merchantID {
			output <- ResultUseCase{Error: errors.New("warehouse doesn't belong to current merchant"), HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: addressDetail}

	})
	return output
}
