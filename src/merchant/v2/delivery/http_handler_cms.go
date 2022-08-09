package delivery

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"

	"github.com/Bhinneka/golib/jsonschema"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/labstack/echo"
)

const (
	merchantIDPath                 = "/:merchantId"
	merchantIDParam                = "merchantId"
	rejectMerchantRegistrationPath = "/:merchantId/reject"
	rejectMerchantUpgradePath      = "/:merchantId/reject-upgrade"
	merchantWarehousePath          = "/:merchantId/warehouses"
	merchantWarehouseAddressPath   = "/:merchantId/warehouses/:addressId"
	merchantVanityURL              = "vanityUrl"
)

// MountCMS specific for CMS usage
func (m *HTTPMerchantHandler) MountCMS(group *echo.Group) {
	group.POST("", m.createMerchant)                        // STG-67
	group.PUT(merchantIDPath, m.updateMerchant)             // update merchant [STG-562]
	group.DELETE(merchantIDPath, m.deleteMerchant)          // delete merchant [STG-563]
	group.GET("/list", m.getList)                           // list and filter [STG-564]
	group.GET(merchantIDPath, m.getMerchant)                // get single merchant [STG-565]
	group.POST(merchantIDPath+"/officer", m.addMerchantPIC) //set merchant pic [STG-939]

	// reject registration
	group.POST(rejectMerchantRegistrationPath, m.rejectMerchantRegistration) // reject merchant registration [STG-778]
	group.POST(rejectMerchantUpgradePath, m.rejectMerchantUpgrade)           // reject merchant upgrade [STG-779]

	// merchant warehouse
	group.GET(merchantWarehousePath, m.getMerchantWarehouse)       // get warehouse list per merchant [STG-823]
	group.GET(merchantWarehouseAddressPath, m.getWarehouseAddress) // get single warehouse address detail [STG-824]

	group.GET("/employees", m.CmsListEmployee)
}

func (m *HTTPMerchantHandler) createMerchant(c echo.Context) error {
	creatorID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	creatorAttr := &model.MerchantUserAttribute{}
	creatorAttr.UserID = creatorID
	creatorAttr.UserIP = c.RealIP()

	payload := &model.B2CMerchantCreateInput{}
	if err := c.Bind(&payload); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	validationError := helper.ValidateTemp("create_merchant_params", payload)
	if validationError != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, validationError.Error()).JSON(c)
	}

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	createResult := <-m.MerchantUseCase.CreateMerchant(newCtx, payload, creatorAttr)
	if createResult.Error != nil {
		return shared.NewHTTPResponse(createResult.HTTPStatus, createResult.Error.Error()).JSON(c)
	}

	result, ok := createResult.Result.(model.B2CMerchantDataV2)
	if !ok {
		err := errors.New(msgErrorResult)
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	// publish to kafka
	m.MerchantUseCase.PublishToKafkaMerchant(c.Request().Context(), result, helper.EventProduceCreateMerchant)

	return shared.NewHTTPResponse(http.StatusCreated, "Create Merchant Response", result).JSON(c)
}

func (m *HTTPMerchantHandler) updateMerchant(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	merchantID := c.Param(merchantIDParam)

	userAttribute := &model.MerchantUserAttribute{}
	userAttribute.UserID = memberID
	userAttribute.UserIP = c.RealIP()

	payload := &model.B2CMerchantCreateInput{}
	payload.ID = merchantID

	if err := c.Bind(&payload); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	mErr := jsonschema.ValidateTemp("update_merchant_params_v2", payload)
	if mErr != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, mErr.Error()).JSON(c)
	}
	headerAuth := c.Request().Header.Get(helper.TextAuthorization)
	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, headerAuth)
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	newCtx = shared.SetDataToContext(newCtx, shared.ContextKey(helper.TextToken), headerAuth)

	mr := <-m.MerchantUseCase.UpdateMerchant(newCtx, payload, userAttribute)
	if mr.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, mr.Error.Error()).JSON(c)
	}

	result, ok := mr.Result.(model.B2CMerchantDataV2)
	if !ok {
		err := errors.New(msgErrorResult)
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	m.MerchantUseCase.PublishToKafkaMerchant(c.Request().Context(), result, helper.EventProduceUpdateMerchant)

	return shared.NewHTTPResponse(http.StatusOK, "Success update merchant", mr.Result).JSON(c)
}

func (m *HTTPMerchantHandler) deleteMerchant(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	merchantID := c.Param(merchantIDParam)

	uaDelete := &model.MerchantUserAttribute{}
	uaDelete.UserID = memberID
	uaDelete.UserIP = c.RealIP()

	ctxReq := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	ctxReq = context.WithValue(ctxReq, middleware.ContextKeyClientIP, c.RealIP())

	mr := <-m.MerchantUseCase.DeleteMerchant(ctxReq, merchantID, uaDelete)
	if mr.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, mr.Error.Error()).JSON(c)
	}

	result, ok := mr.Result.(model.B2CMerchantDataV2)
	if !ok {
		err := errors.New(msgErrorResult)
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	m.MerchantUseCase.PublishToKafkaMerchant(c.Request().Context(), result, helper.EventProduceDeleteMerchant)

	return shared.NewHTTPResponse(http.StatusOK, "Success delete merchant", mr.Result).JSON(c)
}

func (m *HTTPMerchantHandler) getList(c echo.Context) error {
	params := model.QueryParameters{}
	if err := c.Bind(&params); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	if err := params.Validate(); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	result := <-m.MerchantUseCase.GetMerchants(c.Request().Context(), &params)
	if result.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, result.Error.Error()).JSON(c)
	}
	totalData, ok := result.TotalData.(int)
	if !ok {
		return shared.NewHTTPResponse(http.StatusBadRequest, "invalid result").JSON(c)
	}

	totalPage := math.Ceil(float64(totalData) / float64(params.Limit))

	meta := shared.Meta{
		Page:         params.Page,
		Limit:        params.Limit,
		TotalRecords: totalData,
		TotalPages:   int(totalPage),
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success get merchants", result.Result, meta).JSON(c)
}

func (m *HTTPMerchantHandler) getMerchant(c echo.Context) error {
	merchantID := c.Param(merchantIDParam)
	param := model.Param{}
	if err := c.Bind(&param); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	mErr := jsonschema.ValidateTemp("get_merchant_param_v2", param)
	if mErr != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, mErr.Error()).JSON(c)
	}
	if param.IsAttachment == "" {
		param.IsAttachment = "false"
	}

	isAttachment := param.IsAttachment
	fmt.Println(isAttachment)
	// return merchant data
	ctxReq := context.WithValue(c.Request().Context(), shared.ContextKey(helper.TextToken), c.Request().Header.Get(echo.HeaderAuthorization))
	merchantResult := <-m.MerchantUseCase.GetMerchantByID(ctxReq, merchantID, privacy, isAttachment)
	if merchantResult.Error != nil {
		return shared.NewHTTPResponse(merchantResult.HTTPStatus, merchantResult.Error.Error(), helper.EmptySlice{}).JSON(c)
	}

	merchant, ok := merchantResult.Result.(model.B2CMerchantDataV2)
	if !ok {
		return shared.NewHTTPResponse(http.StatusBadRequest, "failed get merchant", make(helper.EmptySlice, 0)).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success get merchant", merchant).JSON(c)
}

func (m *HTTPMerchantHandler) rejectMerchantRegistration(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	merchantID := c.Param(merchantIDParam)

	userAttribute := &model.MerchantUserAttribute{}
	userAttribute.UserID = memberID
	userAttribute.UserIP = c.RealIP()

	ctxReq := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	ctxReq = context.WithValue(ctxReq, middleware.ContextKeyClientIP, c.RealIP())

	mr := <-m.MerchantUseCase.RejectMerchantRegistration(ctxReq, merchantID, userAttribute)
	if mr.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, mr.Error.Error()).JSON(c)
	}

	merchant, ok := mr.Result.(model.B2CMerchantDataV2)
	if !ok {
		return shared.NewHTTPResponse(http.StatusBadRequest, msgErrorResult).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success reject merchant registration", merchant).JSON(c)
}

func (m *HTTPMerchantHandler) rejectMerchantUpgrade(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	var payload struct {
		Reason string `json:"reason,omitempty"`
	}

	c.Bind(&payload)

	merchantID := c.Param(merchantIDParam)

	uaMerchantUpgrade := &model.MerchantUserAttribute{}
	uaMerchantUpgrade.UserID = memberID
	uaMerchantUpgrade.UserIP = c.RealIP()

	ctxReq := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	ctxReq = context.WithValue(ctxReq, middleware.ContextKeyClientIP, c.RealIP())

	rejectRequest := <-m.MerchantUseCase.RejectMerchantUpgrade(ctxReq, merchantID, uaMerchantUpgrade, payload.Reason)
	if rejectRequest.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, rejectRequest.Error.Error()).JSON(c)
	}

	merchant, ok := rejectRequest.Result.(model.B2CMerchantDataV2)
	if !ok {
		return shared.NewHTTPResponse(http.StatusBadRequest, msgErrorResult).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success reject merchant upgrade", merchant).JSON(c)
}

func (m *HTTPMerchantHandler) getMerchantWarehouse(c echo.Context) error {
	whParams := model.ParameterWarehouse{
		StrPage:    c.QueryParam("page"),
		StrLimit:   c.QueryParam("limit"),
		MerchantID: c.Param(merchantIDParam),
		Query:      c.QueryParam("query"),
		Sort:       c.QueryParam("sort"),
		OrderBy:    c.QueryParam("order"),
		ShowAll:    c.QueryParam("showAll"),
	}

	warehouseResults := <-m.MerchantAddressUseCase.GetWarehouseAddresses(c.Request().Context(), &whParams)
	if warehouseResults.Error != nil {
		return shared.NewHTTPResponse(warehouseResults.HTTPStatus, warehouseResults.Error.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	warehouses, ok := warehouseResults.Result.(model.ListWarehouse)
	if !ok {
		err := errors.New("result is not list of warehouse")
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	totalPage := math.Ceil(float64(warehouses.TotalData) / float64(whParams.Limit))

	if len(warehouses.WarehouseData) <= 0 {
		return shared.NewHTTPResponse(http.StatusOK, messageSuccessGet, make(helper.EmptySlice, 0)).JSON(c)
	}

	meta := shared.Meta{
		Page:         whParams.Page,
		Limit:        whParams.Limit,
		TotalRecords: warehouses.TotalData,
		TotalPages:   int(totalPage),
	}
	return shared.NewHTTPResponse(http.StatusOK, messageSuccessGet, meta, warehouses.WarehouseData).JSON(c)
}

func (m *HTTPMerchantHandler) getWarehouseAddress(c echo.Context) error {
	ctxReq := context.WithValue(c.Request().Context(), shared.ContextKey(helper.TextToken), c.Request().Header.Get(echo.HeaderAuthorization))
	result := model.WarehouseData{}
	warehouse := <-m.MerchantAddressUseCase.GetWarehouseAddressByID(ctxReq, c.Param(merchantIDParam), c.Param(TextAddressID))
	if warehouse.Error != nil {
		return shared.NewHTTPResponse(warehouse.HTTPStatus, warehouse.Error.Error()).JSON(c)
	}

	result, ok := warehouse.Result.(model.WarehouseData)
	if !ok {
		return shared.NewHTTPResponse(http.StatusBadRequest, msgErrorResult).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, messageSuccessGet, result).JSON(c)
}

// CmsListEmployee ...
func (m *HTTPMerchantHandler) CmsListEmployee(c echo.Context) error {
	params := model.QueryCmsMerchantEmployeeParameters{}
	if err := c.Bind(&params); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	if err := params.Validate(); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	authorization := c.Request().Header.Get(echo.HeaderAuthorization)
	var token string
	if split := strings.Split(authorization, " "); len(split) > 1 {
		token = split[1]
	}

	result := <-m.MerchantUseCase.CmsGetAllMerchantEmployee(c.Request().Context(), token, &params)
	if result.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, result.Error.Error()).JSON(c)
	}

	totalData, ok := result.TotalData.(int)
	if !ok {
		return shared.NewHTTPResponse(http.StatusBadRequest, "invalid result").JSON(c)
	}

	totalPage := math.Ceil(float64(totalData) / float64(params.Limit))

	meta := shared.Meta{
		Page:         params.Page,
		Limit:        params.Limit,
		TotalRecords: totalData,
		TotalPages:   int(totalPage),
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success get all merchant employee", result.Result, meta).JSON(c)
}

func (m *HTTPMerchantHandler) addMerchantPIC(c echo.Context) error {
	memberId, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	userAttribute := &model.MerchantUserAttribute{}
	userAttribute.UserID = memberId
	userAttribute.UserIP = c.RealIP()

	merchantID := c.Param(merchantIDParam)
	// return merchant data
	payload := &model.B2CMerchantCreateInput{}
	payload.ID = merchantID

	if err := c.Bind(&payload); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	mErr := jsonschema.ValidateTemp("add_merchant_pic_params_v2", payload)
	if mErr != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, mErr.Error()).JSON(c)
	}
	headerAuth := c.Request().Header.Get(helper.TextAuthorization)
	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, headerAuth)
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	newCtx = shared.SetDataToContext(newCtx, shared.ContextKey(helper.TextToken), headerAuth)

	sellerOfficerResult := <-m.MerchantUseCase.AddMerchantPIC(newCtx, payload, userAttribute)
	if sellerOfficerResult.Error != nil {
		return shared.NewHTTPResponse(sellerOfficerResult.HTTPStatus, sellerOfficerResult.Error.Error(), helper.EmptySlice{}).JSON(c)
	}

	sellerOfficer, ok := sellerOfficerResult.Result.(model.B2CMerchantDataV2)
	if !ok {
		return shared.NewHTTPResponse(http.StatusBadRequest, "failed add seller officer", make(helper.EmptySlice, 0)).JSON(c)
	}
	m.MerchantUseCase.PublishToKafkaMerchant(c.Request().Context(), sellerOfficer, helper.EventProduceUpdateMerchant)

	return shared.NewHTTPResponse(http.StatusOK, "Success add seller officer", sellerOfficer).JSON(c)
}
