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
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	memberRepo "github.com/Bhinneka/user-service/src/member/v1/repo"
	modelMerchant "github.com/Bhinneka/user-service/src/merchant/v2/model"
	merchantAddressRepo "github.com/Bhinneka/user-service/src/merchant/v2/repo"
	"github.com/Bhinneka/user-service/src/service"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	sharedRepo "github.com/Bhinneka/user-service/src/shared/repository"
	"github.com/Bhinneka/user-service/src/shipping_address/v2/model"
	"github.com/Bhinneka/user-service/src/shipping_address/v2/repo"
)

const (
	msgErrorSave           = "failed to save address"
	msgErrorFindShipping   = "cannot find shipping address"
	msgErrorDeleteShipping = "failed to delete address"
	msgErrorUpdateShipping = "failed to update address"
	paramsBillingTag       = "params_billing"
	relationNameTable      = "b2c_shippingaddress"
)

// ShippingAddressUseCaseImpl data structure
type ShippingAddressUseCaseImpl struct {
	ShippingAddressRepo      repo.ShippingAddressRepository
	ShippingAddressRedisRepo repo.ShippingAddressRepositoryRedis
	BarracudaService         service.BarracudaServices
	QPublisher               service.QPublisher
	MemberRepo               memberRepo.MemberRepository
	Repository               *sharedRepo.Repository
	ActivityService          service.ActivityServices
	MerchantAddressRepo      merchantAddressRepo.MerchantAddressRepository
}

// NewShippingAddressUseCase function for initialise shipping address use case implementation
func NewShippingAddressUseCase(repository localConfig.ServiceRepository, services localConfig.ServiceShared) ShippingAddressUseCase {
	return &ShippingAddressUseCaseImpl{
		ShippingAddressRepo:      repository.ShippingAddressRepository,
		ShippingAddressRedisRepo: repository.ShippingAddressRedisRepository,
		MemberRepo:               repository.MemberRepository,
		Repository:               repository.Repository,
		BarracudaService:         services.BarracudaService,
		QPublisher:               services.QPublisher,
		ActivityService:          services.ActivityService,
		MerchantAddressRepo:      repository.MerchantAddressRepository,
	}
}

// AddShippingAddress function for add new address
func (s *ShippingAddressUseCaseImpl) AddShippingAddress(ctxReq context.Context, data model.ShippingAddressData) <-chan ResultUseCase {
	ctx := "ShippingAddressUseCase-AddShippingAddress"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags[helper.TextParameter] = data
		data, err := s.validateShippingAddress(ctxReq, data)
		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// Get Member Data from MemberID
		getMemberByID := <-s.MemberRepo.Load(ctxReq, data.MemberID)
		if getMemberByID.Result == nil {
			err := fmt.Errorf("MemberID doesn't exist")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		member, _ := getMemberByID.Result.(memberModel.Member)

		data.Created = time.Now()
		data.LastModified = time.Now()
		data.CreatedBy = data.MemberID

		// parse address structured
		parseResult, err := s.parseShippingAddress(ctxReq, data)
		if err != nil {
			err = errors.New(msgErrorSave)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// add shipping address repository process to database
		saveResult := <-s.ShippingAddressRepo.AddShippingAddress(ctxReq, parseResult)
		if saveResult.Error != nil {
			err := errors.New(msgErrorSave)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		resultData, ok := saveResult.Result.(model.ShippingAddressData)
		if !ok {
			err := errors.New(msgErrorSave)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// set maps data form params
		maps := modelMerchant.Maps{}
		t := time.Now()
		maps.ID = strings.ReplaceAll("MAPS"+t.Format(helper.FormatYmdhisz), ".", "")
		maps.RelationID = resultData.ID
		maps.RelationName = relationNameTable
		maps.Label = data.Label
		maps.Latitude = data.Latitude
		maps.Longitude = data.Longitude

		saveResultMaps := <-s.MerchantAddressRepo.AddUpdateAddressMaps(ctxReq, maps)
		if saveResultMaps.Error != nil {
			output <- ResultUseCase{Error: saveResultMaps.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		resultData.Latitude = data.Latitude
		resultData.Longitude = data.Longitude
		resultData.IsMapAvailable = helper.ValidateLatLong(maps.Latitude, maps.Longitude)

		// set default billing address
		go s.SaveBillingAddress(ctxReq, member, parseResult)

		// send to audit trail activity service
		go s.insertLogShipping(ctxReq, model.ShippingAddressData{}, parseResult, "AddShippingAddress")

		// delete data from redis cache
		err = <-s.ShippingAddressRedisRepo.DeleteMultipleRedis(data.MemberID)
		if err != nil {
			err := errors.New(msgErrorSave)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		tags[helper.TextResponse] = resultData
		output <- ResultUseCase{Result: resultData}
	})

	return output
}

// SaveBillingAddress function for add billing default
func (s *ShippingAddressUseCaseImpl) SaveBillingAddress(ctxReq context.Context, member memberModel.Member, address model.ShippingAddressData) <-chan ResultUseCase {
	ctx := "ShippingAddressUseCase-SaveBillingAddress"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		if member.Address.Address != "" {
			output <- ResultUseCase{Result: nil}
			return
		}
		billAddress := memberModel.Address{
			Street1:       address.Street1,
			Street2:       address.Street2,
			ZipCode:       address.PostalCode,
			SubDistrictID: address.SubDistrictID,
			SubDistrict:   address.SubDistrictName,
			DistrictID:    address.SubDistrictID,
			District:      address.DistrictName,
			CityID:        address.CityID,
			City:          address.CityName,
			ProvinceID:    address.ProvinceID,
			Province:      address.ProvinceName,
		}

		address.Street1 = helper.ClearHTML(address.Street1)

		billAddress.Address = address.Street1
		if len(address.Street2) > 0 {
			address.Street2 = helper.ClearHTML(address.Street2)
			billAddress.Address = fmt.Sprintf("%s\n%s", address.Street1, address.Street2)
		}
		tags[paramsBillingTag] = address
		tags[helper.TextParameter] = billAddress
		member.Address = billAddress

		saveResult := <-s.MemberRepo.Save(ctxReq, member)
		if saveResult.Error != nil {
			err := errors.New("failed to save billing address")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusInternalServerError}
			return
		}

		output <- ResultUseCase{Result: nil}
	})
	return output
}

// DeleteShippingAddressByID function for add new address
func (s *ShippingAddressUseCaseImpl) DeleteShippingAddressByID(ctxReq context.Context, shippingID string, memberID string) <-chan ResultUseCase {
	ctx := "ShippingAddressUseCase-DeleteShippingAddressByID"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		// find shipping by id
		findResult := <-s.ShippingAddressRepo.FindShippingAddressByID(ctxReq, shippingID, memberID)
		if findResult.Error != nil || memberID == "" {
			err := errors.New(msgErrorFindShipping)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusNotFound}
			return
		}

		detailFind, ok := findResult.Result.(model.ShippingAddressData)
		if !ok {
			err := errors.New(msgErrorDeleteShipping)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		if detailFind.IsPrimary {
			err := errors.New("cannot delete primary address")
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// delete shipping address from database
		result := <-s.ShippingAddressRepo.DeleteShippingAddressByID(ctxReq, shippingID)
		if result.Error != nil {
			err := errors.New(msgErrorDeleteShipping)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// send to audit trail activity service
		go s.insertLogShipping(ctxReq, detailFind, model.ShippingAddressData{}, "DeleteShippingAddressByID")

		// delete data from redis cache
		err := <-s.ShippingAddressRedisRepo.DeleteMultipleRedis(memberID)
		if err != nil {
			err := errors.New(msgErrorDeleteShipping)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: nil}
	})
	return output
}

// UpdateShippingAddress function for update address
func (s *ShippingAddressUseCaseImpl) UpdateShippingAddress(ctxReq context.Context, data model.ShippingAddressData) <-chan ResultUseCase {
	ctx := "ShippingAddressUseCase-UpdateShippingAddress"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		member, detailFind, errCode, errFind := s.validateShippingAddressID(ctxReq, data)
		if errFind != nil {
			tags[helper.TextResponse] = errFind
			output <- ResultUseCase{Error: errFind, HTTPStatus: errCode}
			return
		}

		data, err := s.validateShippingAddress(ctxReq, data)
		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		data.Version = detailFind.Version
		data.LastModified = time.Now()
		data.CreatedBy = detailFind.CreatedBy
		data.IsPrimary = detailFind.IsPrimary
		data.ModifiedBy = data.MemberID

		// parse address structured
		parseResult, errorParse := s.parseShippingAddress(ctxReq, data)
		if errorParse != nil {
			err := errors.New(msgErrorUpdateShipping)
			tracer.Log(ctxReq, helper.TextExecUsecase, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// update shipping address repository process to database
		saveResult := <-s.ShippingAddressRepo.UpdateShippingAddress(ctxReq, parseResult)
		if saveResult.Error != nil {
			err := errors.New(msgErrorUpdateShipping)
			tracer.Log(ctxReq, helper.TextExecUsecase, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		resultData, ok := saveResult.Result.(model.ShippingAddressData)
		if !ok {
			err := errors.New(msgErrorUpdateShipping)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// set maps data form params
		maps := modelMerchant.Maps{}
		t := time.Now()
		maps.ID = strings.ReplaceAll("MAPS"+t.Format(helper.FormatYmdhisz), ".", "")
		maps.RelationID = resultData.ID
		maps.RelationName = relationNameTable
		maps.Label = data.Label
		maps.Latitude = data.Latitude
		maps.Longitude = data.Longitude

		saveResultMaps := <-s.MerchantAddressRepo.AddUpdateAddressMaps(ctxReq, maps)
		if saveResultMaps.Error != nil {
			output <- ResultUseCase{Error: saveResultMaps.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		resultData.Latitude = data.Latitude
		resultData.Longitude = data.Longitude
		resultData.IsMapAvailable = helper.ValidateLatLong(maps.Latitude, maps.Longitude)

		// set default billing address
		go s.SaveBillingAddress(ctxReq, member, parseResult)

		// send to audit trail activity service
		go s.insertLogShipping(ctxReq, detailFind, parseResult, "UpdateShippingAddress")

		// delete data from redis cache
		err = <-s.ShippingAddressRedisRepo.DeleteMultipleRedis(data.MemberID)
		if err != nil {
			err := errors.New(msgErrorUpdateShipping)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: resultData}

	})
	return output
}

// insertLogShipping function for send activity log
func (s *ShippingAddressUseCaseImpl) insertLogShipping(ctxReq context.Context, oldData, newData model.ShippingAddressData, action string) error {
	logChanges := model.ShippingAddressLog{
		Before: &oldData,
		After:  &newData,
	}
	return s.QPublisher.QueueJob(ctxReq, logChanges, oldData.ID, action)
}

func (s *ShippingAddressUseCaseImpl) InsertLogShipping(ctxReq context.Context, oldData, newData *model.ShippingAddressData, action string) error {
	targetID := newData.ID
	if targetID == "" {
		targetID = oldData.ID
	}

	if oldData.ID == "" {
		action = helper.TextInsertUpper
	}

	payload := serviceModel.Payload{
		Module:    model.Module,
		Action:    action,
		Target:    targetID,
		CreatorID: newData.CreatedBy,
		EditorID:  newData.ModifiedBy,
	}

	return s.ActivityService.InsertLog(ctxReq, oldData, newData, payload)
}

// validateShippingAddressID function for validate ID shipping address & member ID
func (s *ShippingAddressUseCaseImpl) validateShippingAddressID(ctxReq context.Context, data model.ShippingAddressData) (memberModel.Member, model.ShippingAddressData, int, error) {
	var (
		member     memberModel.Member
		detailFind model.ShippingAddressData
	)
	// Get Member Data from MemberID
	getMemberByID := <-s.MemberRepo.Load(ctxReq, data.MemberID)
	if getMemberByID.Result == nil {
		err := fmt.Errorf("MemberID doesn't exist")
		return member, detailFind, http.StatusBadRequest, err
	}

	member, _ = getMemberByID.Result.(memberModel.Member)

	// find shipping by id
	findResult := <-s.ShippingAddressRepo.FindShippingAddressByID(ctxReq, data.ID, data.MemberID)
	if findResult.Error != nil {
		err := errors.New(msgErrorFindShipping)
		return member, detailFind, http.StatusNotFound, err
	}

	detailFind, ok := findResult.Result.(model.ShippingAddressData)
	if !ok {
		err := errors.New(msgErrorUpdateShipping)
		return member, detailFind, http.StatusBadRequest, err
	}

	return member, detailFind, http.StatusOK, nil
}

// validateShippingAddress function for validating shipping address data
func (s *ShippingAddressUseCaseImpl) validateShippingAddress(ctxReq context.Context, data model.ShippingAddressData) (model.ShippingAddressData, error) {
	// validate maximum number of address, max: 20
	data, err := s.validateMaxAddress(ctxReq, data)
	if err != nil {
		return data, err
	}

	if len(data.Label) <= 0 {
		err := errors.New("label is required")
		return data, err
	}

	if !golib.ValidateAlphanumericWithSpace(data.Label, false) {
		err := errors.New("label only alphabet is allowed")
		return data, err
	}

	if len(data.Name) <= 0 {
		err := errors.New("name is required")
		return data, err
	}

	if !golib.ValidateAlphanumericWithSpace(data.Name, false) {
		err := errors.New("name only alphabet is allowed")
		return data, err
	}

	if len(data.Mobile) <= 0 {
		err := errors.New("mobile is required")
		return data, err
	}

	if helper.ValidateMobileNumberMaxInput(data.Mobile) != nil {
		err := errors.New("mobile phone number is in bad format")
		return data, err
	}

	if len(data.Phone) > 0 && helper.ValidatePhoneNumberMaxInput(data.Phone) != nil {
		err := errors.New("phone number is in bad format")
		return data, err
	}

	data, err = s.validateLocationShippingAddress(data)
	if err != nil {
		return data, err
	}

	return data, nil
}

// validateMaxAddress function for validating maximum address
func (s *ShippingAddressUseCaseImpl) validateMaxAddress(ctxReq context.Context, data model.ShippingAddressData) (model.ShippingAddressData, error) {
	// validate maximum number of address, max: 20
	countAddress := <-s.ShippingAddressRepo.CountShippingAddressByUserID(ctxReq, data.MemberID)
	if countAddress.Error != nil {
		err := errors.New("shipping address reaches maximum limit")
		return data, err
	}

	total, _ := countAddress.Result.(int)
	if data.ID == "" && total >= model.MaximumShippingAddress {
		err := errors.New("shipping address reaches maximum limit")
		return data, err
	}

	if total == 0 || data.IsPrimary {
		data.IsPrimary = true
	} else {
		data.IsPrimary = false
	}

	return data, nil
}

// validateLocationShippingAddress function for validating shipping address data location
func (s *ShippingAddressUseCaseImpl) validateLocationShippingAddress(shipping model.ShippingAddressData) (model.ShippingAddressData, error) {
	if len(shipping.Street1) <= 0 {
		err := errors.New("street is required")
		return shipping, err
	}

	if !stringLib.ValidateLatinOnlyExcepTag(shipping.Street1) {
		err := errors.New("street1 only latin character")
		return shipping, err
	}

	if len(shipping.Street2) > 0 && !stringLib.ValidateLatinOnlyExcepTag(shipping.Street2) {
		err := errors.New("street2 only latin character")
		return shipping, err
	}

	if len(shipping.PostalCode) <= 0 {
		err := errors.New("postal code is required")
		return shipping, err
	}

	shipping, err := s.validateLocationAreaShippingAddress(shipping)
	if err != nil {
		return shipping, err
	}
	return shipping, nil
}

// validateLocationShippingAddress function for validating shipping address data location area
func (s *ShippingAddressUseCaseImpl) validateLocationAreaShippingAddress(shippingData model.ShippingAddressData) (model.ShippingAddressData, error) {
	if !golib.ValidateNumeric(shippingData.PostalCode) {
		err := errors.New("postal code only numeric is allowed")
		return shippingData, err
	}

	if len(shippingData.SubDistrictID) <= 0 {
		err := errors.New("subDistrict ID is required")
		return shippingData, err
	}

	if !golib.ValidateNumeric(shippingData.SubDistrictID) {
		err := errors.New("subdistrict ID only numeric is allowed")
		return shippingData, err
	}

	if len(shippingData.SubDistrictName) <= 0 {
		err := errors.New("subdistrict name is required")
		return shippingData, err
	}

	if len(shippingData.DistrictID) <= 0 {
		err := errors.New("district ID is required")
		return shippingData, err
	}

	if !golib.ValidateNumeric(shippingData.DistrictID) {
		err := errors.New("district ID only numeric is allowed")
		return shippingData, err
	}

	if len(shippingData.DistrictName) <= 0 {
		err := errors.New("district name is required")
		return shippingData, err
	}

	if len(shippingData.CityID) <= 0 {
		err := errors.New("city ID is required")
		return shippingData, err
	}

	if !golib.ValidateNumeric(shippingData.CityID) {
		err := errors.New("city ID only numeric is allowed")
		return shippingData, err
	}

	if len(shippingData.CityName) <= 0 {
		err := errors.New("city name is required")
		return shippingData, err
	}

	if len(shippingData.ProvinceID) <= 0 {
		err := errors.New("province ID is required")
		return shippingData, err
	}

	if !golib.ValidateNumeric(shippingData.ProvinceID) {
		err := errors.New("province ID only numeric is allowed")
		return shippingData, err
	}

	if len(shippingData.ProvinceName) <= 0 {
		err := errors.New("province name is required")
		return shippingData, err
	}

	if len(shippingData.Ext) >= 1 && !golib.ValidateNumeric(shippingData.Ext) {
		err := errors.New("ext only numeric is allowed")
		return shippingData, err
	}

	return shippingData, nil
}

// parseShipping parseShippingAddress function for validating shipping address data
func (s *ShippingAddressUseCaseImpl) parseShippingAddress(ctxReq context.Context, data model.ShippingAddressData) (model.ShippingAddressData, error) {
	ctx := "ShippingAddressUseCase-parseShippingAddress"

	tr := tracer.StartTrace(context.Background(), ctx)
	tags := make(map[string]interface{})
	defer func() {
		tr.Finish(tags)
	}()

	// validate zipcode data and address from barracuda service
	var zipcodeParam serviceModel.ZipCodeQueryParameter
	zipcodeParam.SubDistrictID = data.SubDistrictID
	zipcodeParam.ZipCode = data.PostalCode

	tags["param"] = zipcodeParam
	barracudaService := <-s.BarracudaService.FindZipcode(ctxReq, zipcodeParam)
	if barracudaService.Error != nil {
		return data, barracudaService.Error
	}

	// get detail zipcode
	detailZipCode, ok := barracudaService.Result.(serviceModel.ResponseZipCode)
	if !ok {
		err := errors.New("failed parse zipcode response")
		return data, err
	}

	if len(detailZipCode.Data) == 0 {
		err := errors.New("Data zipcode Not Found")
		return data, err
	}

	var address model.ShippingAddressData
	address.ID = data.ID
	address.MemberID = data.MemberID
	address.Name = data.Name
	address.Mobile = data.Mobile
	address.Phone = data.Phone
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
	address.Street1 = data.Street1
	address.Street2 = data.Street2
	address.Ext = data.Ext
	address.Label = data.Label
	address.IsPrimary = data.IsPrimary
	address.Created = data.Created
	address.LastModified = data.LastModified
	address.Version = data.Version
	address.CreatedBy = data.CreatedBy
	address.ModifiedBy = data.ModifiedBy

	if address.ID == "" {
		t := time.Now()
		address.ID = "ADDR" + t.Format(helper.FormatYmdhisz)
		address.ID = strings.Replace(address.ID, ".", "", -1)
	}

	tags["result"] = address
	return address, nil
}

// GetListShippingAddress function for getting list of shipping address
func (s *ShippingAddressUseCaseImpl) GetListShippingAddress(ctxReq context.Context, paramShipping *model.ParametersShippingAddress) <-chan ResultUseCase {
	ctx := "ShippingAddressUseCase-GetListShippingAddress"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		var err error
		// validate all parameters
		paging, err := helper.ValidatePagination(
			helper.PaginationParameters{
				Page:     1, // default
				StrPage:  paramShipping.StrPage,
				Limit:    20, // default
				StrLimit: paramShipping.StrLimit,
			})

		if err != nil {
			tags[helper.TextResponse] = err.Error()
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		paramShipping.Page = paging.Page
		paramShipping.Limit = paging.Limit
		paramShipping.Offset = paging.Offset
		tags[helper.TextParameter] = paramShipping

		// load first from redis check cache
		shippingAddressRedis := <-s.ShippingAddressRedisRepo.LoadRedisMeta(paramShipping.MemberID, strconv.Itoa(paramShipping.Page), strconv.Itoa(paramShipping.Limit))
		if shippingAddressRedis.Error == nil {
			resp := shippingAddressRedis.Result.(model.ListShippingAddress)
			resp = s.SetAvailableMaps(ctxReq, resp)

			output <- ResultUseCase{Result: resp}
			return
		}

		shippingAddresses := <-s.ShippingAddressRepo.GetListShippingAddress(ctxReq, paramShipping)
		if shippingAddresses.Error != nil {
			httpStatus := http.StatusInternalServerError

			// when data is not found
			if shippingAddresses.Error == sql.ErrNoRows {
				httpStatus = http.StatusNotFound
				shippingAddresses.Error = fmt.Errorf(helper.ErrorDataNotFound, "shipping_address")
			}

			output <- ResultUseCase{Error: shippingAddresses.Error, HTTPStatus: httpStatus}
			return
		}

		shippingAddress := shippingAddresses.Result.(model.ListShippingAddress)
		shippingAddress = s.SetAvailableMaps(ctxReq, shippingAddress)

		totalResult := <-s.ShippingAddressRepo.GetTotalShippingAddress(ctxReq, paramShipping)
		if totalResult.Error != nil {
			output <- ResultUseCase{Error: totalResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		total := totalResult.Result.(int)
		shippingAddress.TotalData = total

		// save data to redis meta cache
		err = <-s.ShippingAddressRedisRepo.SaveRedisMeta(paramShipping.MemberID, strconv.Itoa(paramShipping.Page), strconv.Itoa(paramShipping.Limit), shippingAddress)
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: shippingAddress}
	})

	return output
}

// GetDetailShippingAddress function for getting detail of shipping address
func (s *ShippingAddressUseCaseImpl) GetDetailShippingAddress(ctxReq context.Context, shippingID string, memberID string) <-chan ResultUseCase {
	ctx := "ShippingAddressUseCase-GetDetailShippingAddress"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		// find shipping by id
		findResult := <-s.ShippingAddressRepo.FindShippingAddressByID(ctxReq, shippingID, memberID)
		if findResult.Error != nil {
			err := errors.New(msgErrorFindShipping)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusNotFound}
			return
		}

		result, ok := findResult.Result.(model.ShippingAddressData)
		if !ok {
			err := errors.New(msgErrorUpdateShipping)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		mapsResult := <-s.MerchantAddressRepo.FindAddressMaps(ctxReq, result.ID, "b2c_shippingaddress")
		if mapsResult.Error == nil {
			dataMaps, _ := mapsResult.Result.(modelMerchant.Maps)
			if helper.ValidateLatLong(dataMaps.Latitude, dataMaps.Longitude) {
				result.IsMapAvailable = helper.ValidateLatLong(dataMaps.Latitude, dataMaps.Longitude)

				result.MapsID = dataMaps.ID
				result.RelationID = dataMaps.RelationID
				result.RelationName = dataMaps.RelationName
				result.Label = dataMaps.Label
				result.Latitude = dataMaps.Latitude
				result.Longitude = dataMaps.Longitude
			}
		}

		output <- ResultUseCase{Result: result}

	})
	return output
}

// GetAllListShippingAddress function for getting list of shipping address
func (s *ShippingAddressUseCaseImpl) GetAllListShippingAddress(ctxReq context.Context, paramList *model.ParametersShippingAddress, memberID string) <-chan ResultUseCase {
	ctx := "ShippingAddressUseCase-GetAllListShippingAddress"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		var err error
		// validate all parameters
		paging, err := helper.ValidatePagination(
			helper.PaginationParameters{
				Page:     1, // default
				StrPage:  paramList.StrPage,
				Limit:    20, // default
				StrLimit: paramList.StrLimit,
			})

		if err != nil {
			tags[helper.TextResponse] = err.Error()
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		paramList.Page = paging.Page
		paramList.Limit = paging.Limit
		paramList.Offset = paging.Offset
		tags[helper.TextParameter] = paramList

		shippingAddressesResult := <-s.ShippingAddressRepo.GetListShippingAddress(ctxReq, paramList)
		if shippingAddressesResult.Error != nil {
			httpStatus := http.StatusInternalServerError

			// when data is not found
			if shippingAddressesResult.Error == sql.ErrNoRows {
				httpStatus = http.StatusNotFound
				shippingAddressesResult.Error = fmt.Errorf(helper.ErrorDataNotFound, "shipping_address")
			}

			output <- ResultUseCase{Error: shippingAddressesResult.Error, HTTPStatus: httpStatus}
			return
		}

		shippingAddress := shippingAddressesResult.Result.(model.ListShippingAddress)
		shippingAddress = s.SetAvailableMaps(ctxReq, shippingAddress)

		totalResult := <-s.ShippingAddressRepo.GetTotalShippingAddress(ctxReq, paramList)
		if totalResult.Error != nil {
			output <- ResultUseCase{Error: totalResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		total := totalResult.Result.(int)
		shippingAddress.TotalData = total

		output <- ResultUseCase{Result: shippingAddress}
	})

	return output
}

// UpdatePrimaryShippingAddressByID function for add new address
func (s *ShippingAddressUseCaseImpl) UpdatePrimaryShippingAddressByID(ctxReq context.Context, params model.ParamaterPrimaryShippingAddress) <-chan ResultUseCase {
	ctx := "ShippingAddressUseCase-UpdatePrimaryShippingAddressByID"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		// find shipping by id
		findResult := <-s.ShippingAddressRepo.FindShippingAddressByID(ctxReq, params.ShippingID, params.MemberID)
		if findResult.Error != nil || params.MemberID == "" {
			err := errors.New(msgErrorFindShipping)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusNotFound}
			return
		}

		detailFind, ok := findResult.Result.(model.ShippingAddressData)
		if !ok {
			err := errors.New(msgErrorUpdateShipping)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		oldShipping := detailFind

		s.Repository.StartTransaction()

		// update primary shipping address from database
		result := <-s.ShippingAddressRepo.UpdatePrimaryShippingAddressByID(ctxReq, params.MemberID)
		if result.Error != nil {
			err := errors.New("failed to update primary address")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		detailFind.IsPrimary = true
		if params.UserID != "" {
			detailFind.ModifiedBy = params.UserID
		}

		// update primary shipping address from database
		resultUpdate := <-s.ShippingAddressRepo.UpdateShippingAddress(ctxReq, detailFind)
		if resultUpdate.Error != nil {
			err := errors.New("failed to update primary address")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		s.Repository.Commit()

		// send to audit trail activity service
		go s.insertLogShipping(ctxReq, oldShipping, detailFind, "UpdatePrimaryShippingAddressByID")

		// delete data from redis cache
		err := <-s.ShippingAddressRedisRepo.DeleteMultipleRedis(params.MemberID)
		if err != nil {
			output <- ResultUseCase{Error: errors.New(msgErrorDeleteShipping), HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: nil}
	})
	return output
}

// GetPrimaryShippingAddress function for getting detail of primary shipping address
func (s *ShippingAddressUseCaseImpl) GetPrimaryShippingAddress(ctxReq context.Context, memberID string) <-chan ResultUseCase {
	ctx := "ShippingAddressUseCase-GetPrimaryShippingAddress"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		// find shipping by id
		findResult := <-s.ShippingAddressRepo.FindShippingAddressPrimaryByID(ctxReq, memberID)
		if findResult.Error != nil {
			err := errors.New(msgErrorFindShipping)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusNotFound}
			return
		}

		result, ok := findResult.Result.(model.ShippingAddressData)
		if !ok {
			err := errors.New(msgErrorFindShipping)
			tracer.Log(ctxReq, helper.TextExecUsecase, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: result}

	})
	return output
}

// setAvailableMaps function for validating shipping address data location area
func (s *ShippingAddressUseCaseImpl) SetAvailableMaps(ctxReq context.Context, data model.ListShippingAddress) model.ListShippingAddress {
	for idx, val := range data.ShippingAddress {
		mapsResult := <-s.MerchantAddressRepo.FindAddressMaps(ctxReq, val.ID, relationNameTable)
		if mapsResult.Error == nil {
			detailMaps, _ := mapsResult.Result.(modelMerchant.Maps)
			data.ShippingAddress[idx].IsMapAvailable = helper.ValidateLatLong(detailMaps.Latitude, detailMaps.Longitude)
		}
	}
	return data
}
