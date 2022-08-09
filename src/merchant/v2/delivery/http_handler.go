package delivery

import (
	"context"
	"encoding/base64"
	"errors"
	"math"
	"net/http"
	"strings"

	"github.com/Bhinneka/golib/jsonschema"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/Bhinneka/user-service/src/merchant/v2/usecase"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
	"github.com/tealeg/xlsx"
)

const (
	// TextExtractMemberID text
	TextExtractMemberID = "merchant_extract_member_id"
	// TextAddressID text
	TextAddressID       = "addressId"
	msgErrorResult      = "result is not proper response"
	scopeParseWarehouse = "parse_warehouse"
	messageSuccessGet   = "Success Get Data"
	pathWarehouseID     = "/warehouse/:addressId"
	privacy             = "private"
)

// HTTPMerchantHandler model
type HTTPMerchantHandler struct {
	MerchantUseCase        usecase.MerchantUseCase
	MerchantAddressUseCase usecase.MerchantAddressUseCase
}

// NewHTTPHandler function for initialise *HTTPAuthHandler
func NewHTTPHandler(MerchantUseCase usecase.MerchantUseCase,
	MerchantAddressUseCase usecase.MerchantAddressUseCase) *HTTPMerchantHandler {
	return &HTTPMerchantHandler{
		MerchantUseCase:        MerchantUseCase,
		MerchantAddressUseCase: MerchantAddressUseCase,
	}
}

// MountMe function for mounting routes
func (m *HTTPMerchantHandler) MountMe(group *echo.Group) {
	group.POST("", m.AddMerchant)
	group.POST("/upgrade", m.UpgradeMerchant)
	group.GET("", m.GetMerchantByUserID)
	group.POST("/warehouse", m.AddWarehouse)
	group.PUT(pathWarehouseID, m.UpdateWarehouse)
	group.PUT("/warehouse/:addressId/set-primary", m.UpdateWarehousePrimary)
	group.GET("/warehouse", m.GetWarehouse)
	group.GET(pathWarehouseID, m.GetWarehouseDetail)
	group.DELETE(pathWarehouseID, m.DeleteWarehouse)
	group.PUT("", m.UpdateMerchant)
	group.PUT("/change-name", m.ChangeMerchantName)
	group.PATCH("", m.UpdateMerchantPartial)

	group.POST("/employees", m.AddEmployee)
	group.GET("/employees", m.ListEmployee)
	group.GET("/employees/:memberId", m.GetEmployee)
	group.PUT("/employees/:memberId", m.UpdateEmployee)
	group.POST("/employees/resend-email", m.ResendEmailEmployee)

	group.POST("/clear-upgrade", m.clearRejectMerchantUpgrade)
}

// MountMerchant function for mounting routes
func (m *HTTPMerchantHandler) MountMerchant(group *echo.Group) {
	group.GET("/merchantbank", m.GetMerchantBank)
	group.POST("/bulk-merchant-send", m.BulkMerchantSend)
	group.POST("/merchant-send", m.MerchantSend)
	group.POST("/merchant/validate", m.CheckMerchantName)
}

// AddMerchant function for registering new merchant
func (m *HTTPMerchantHandler) AddMerchant(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	payload := &model.B2CMerchantCreateInput{}
	if err := c.Bind(&payload); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	userAttribute := &model.MerchantUserAttribute{}
	userAttribute.UserID = memberID
	userAttribute.UserIP = c.RealIP()

	payload.UserID = memberID

	mErr := helper.ValidateTemp("add_merchant_params_v2", payload)
	if mErr != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, mErr.Error()).JSON(c)
	}

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	saveResult := <-m.MerchantUseCase.AddMerchant(newCtx, payload, userAttribute)
	if saveResult.Error != nil {
		return shared.NewHTTPResponse(saveResult.HTTPStatus, saveResult.Error.Error()).JSON(c)
	}

	result, ok := saveResult.Result.(model.B2CMerchantDataV2)
	if !ok {
		err := errors.New(msgErrorResult)
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusCreated, "Merchant Register Response", result).JSON(c)
}

// UpgradeMerchant function for upgrade merchant
func (m *HTTPMerchantHandler) UpgradeMerchant(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	payload := model.B2CMerchantCreateInput{}
	userAttribute := &model.MerchantUserAttribute{}
	userAttribute.UserID = memberID
	userAttribute.UserIP = c.RealIP()

	if err := c.Bind(&payload); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	payload.UserID = memberID
	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	saveResult := <-m.MerchantUseCase.UpgradeMerchant(newCtx, &payload, userAttribute)
	if saveResult.Error != nil {
		return shared.NewHTTPResponse(saveResult.HTTPStatus, saveResult.Error.Error()).JSON(c)
	}

	result, ok := saveResult.Result.(model.B2CMerchantDataV2)
	if !ok {
		err := errors.New(msgErrorResult)
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Merchant Upgrade Response", result).JSON(c)
}

// GetMerchantBank function for getting list of merchant bank
func (m *HTTPMerchantHandler) GetMerchantBank(c echo.Context) error {
	params := model.ParametersMerchantBank{
		StrPage:  c.QueryParam("page"),
		StrLimit: c.QueryParam("limit"),
		Sort:     c.QueryParam("sort"),
		OrderBy:  c.QueryParam("orderBy"),
		Status:   c.QueryParam("status"),
	}

	merchantBankResult := <-m.MerchantUseCase.GetListMerchantBank(&params)
	if merchantBankResult.Error != nil {
		return shared.NewHTTPResponse(merchantBankResult.HTTPStatus, merchantBankResult.Error.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	merchantBank, ok := merchantBankResult.Result.(model.ListMerchantBank)
	if !ok {
		err := errors.New("result is not list of merchant bank")
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	totalPage := math.Ceil(float64(merchantBank.TotalData) / float64(params.Limit))

	if len(merchantBank.MerchantBank) <= 0 {
		response := shared.NewHTTPResponse(http.StatusOK, "Get Merchant Bank Response", make(helper.EmptySlice, 0))
		response.SetSuccess(false)
		return response.JSON(c)
	}

	meta := shared.Meta{
		Page:         params.Page,
		Limit:        params.Limit,
		TotalRecords: merchantBank.TotalData,
		TotalPages:   int(totalPage),
	}
	return shared.NewHTTPResponse(http.StatusOK, "Get Merchant Bank Response", merchantBank.MerchantBank, meta).JSON(c)
}

// CheckMerchantName for check merchant name availability
func (m *HTTPMerchantHandler) CheckMerchantName(c echo.Context) error {
	var payload struct {
		MerchantName string `json:"merchantName"`
	}

	if err := c.Bind(&payload); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	res := <-m.MerchantUseCase.CheckMerchantName(c.Request().Context(), payload.MerchantName)
	if res.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, res.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "success merchant name available", res.Result).JSON(c)
}

// GetMerchantByUserID for get merchant detail by user
func (m *HTTPMerchantHandler) GetMerchantByUserID(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	token, _ := c.Get("token").(*jwt.Token)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	ctxReq := context.WithValue(c.Request().Context(), shared.ContextKey(helper.TextToken), c.Request().Header.Get(echo.HeaderAuthorization))
	saveResult := <-m.MerchantUseCase.GetMerchantByUserID(ctxReq, memberID, token.Raw)
	if saveResult.Error != nil {
		return shared.NewHTTPResponse(saveResult.HTTPStatus, saveResult.Error.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "success get merchant detail", saveResult.Result).JSON(c)
}

// MerchantSend function for send to kafka for updated nav
func (m *HTTPMerchantHandler) MerchantSend(c echo.Context) error {
	merchantID := c.FormValue(helper.TextMerchantID)
	eventType := c.FormValue("eventType")
	isAttachment := "false"
	err := middleware.ExtractClaimsIsAdmin(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, err.Error()).JSON(c)
	}

	// get merchant data
	merchantResult := <-m.MerchantUseCase.GetMerchantByID(c.Request().Context(), merchantID, privacy, isAttachment)
	if merchantResult.Error != nil {
		return shared.NewHTTPResponse(merchantResult.HTTPStatus, merchantResult.Error.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	merchant, ok := merchantResult.Result.(model.B2CMerchantDataV2)
	if !ok {
		err := errors.New("failed")
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	// publish to kafka
	m.MerchantUseCase.PublishToKafkaMerchant(c.Request().Context(), merchant, eventType)

	return shared.NewHTTPResponse(http.StatusOK, "Merchant Send Response", merchant).JSON(c)
}

// BulkMerchantSend function for getting list of merchant bank
func (m *HTTPMerchantHandler) BulkMerchantSend(c echo.Context) error {
	base64File := c.FormValue("file")
	eventType := c.FormValue("eventType")
	isAttachment := "false"

	err := middleware.ExtractClaimsIsAdmin(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, err.Error()).JSON(c)
	}

	fileDecode, err := base64.StdEncoding.DecodeString(base64File)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	xlsBinary, err := xlsx.OpenBinary(fileDecode)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)

	}

	successMerchant := make([]string, 0)
	for _, row := range xlsBinary.Sheets[0].Rows {
		if row.Cells[0].String() == helper.TextMerchantID || row.Cells[0].String() == "" {
			continue
		}

		merchantID := row.Cells[0].String()

		// get merchant data
		merchantResult := <-m.MerchantUseCase.GetMerchantByID(c.Request().Context(), merchantID, privacy, isAttachment)
		if merchantResult.Error != nil {
			//Error, skip row
			continue
		}
		merchant, ok := merchantResult.Result.(model.B2CMerchantDataV2)
		if !ok {
			//Error, skip row
			continue
		}

		successMerchant = append(successMerchant, merchantID)

		// publish to kafka
		m.MerchantUseCase.PublishToKafkaMerchant(c.Request().Context(), merchant, eventType)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Bulk Merchant Send Response", successMerchant).JSON(c)
}

// AddWarehouse function for add new merchant address
func (m *HTTPMerchantHandler) AddWarehouse(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	// create warehouse process with parameter `add`
	result, errCode, err := m.SaveUpdateWarehouse(c, helper.TextAdd, helper.TextMe, memberID)
	if err != nil {
		return shared.NewHTTPResponse(errCode, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusCreated, "Warehouse location Response", result).JSON(c)
}

// UpdateWarehouse function for update merchant address
func (m *HTTPMerchantHandler) UpdateWarehouse(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	// create warehouse process with parameter `update`
	result, errCode, err := m.SaveUpdateWarehouse(c, helper.TextUpdate, helper.TextMe, memberID)
	if err != nil {
		return shared.NewHTTPResponse(errCode, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Warehouse location Response", result).JSON(c)
}

// SaveUpdateWarehouse function for save update warehouse address
func (m *HTTPMerchantHandler) SaveUpdateWarehouse(c echo.Context, action, source, memberID string) (model.WarehouseData, int, error) {
	result := model.WarehouseData{}
	warehouse := &model.WarehouseData{}
	if err := c.Bind(&warehouse); err != nil {
		return result, http.StatusBadRequest, err
	}

	warehouse.Address = helper.ValidateHTML(warehouse.Address)

	// validate schema fields
	mErr := jsonschema.ValidateTemp("add_warehouse_params_v2", warehouse)
	if mErr != nil {
		return result, http.StatusBadRequest, mErr
	}

	// set id for update
	if action == helper.TextUpdate {
		warehouse.ID = c.Param(TextAddressID)
	}
	if action == helper.TextAdd {
		warehouse.Status = helper.TextActive
	}

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	saveResult := <-m.MerchantAddressUseCase.AddUpdateWarehouseAddress(newCtx, *warehouse, memberID, action)

	if saveResult.Error != nil {
		return result, saveResult.HTTPStatus, saveResult.Error
	}

	result, ok := saveResult.Result.(model.WarehouseData)
	if !ok {
		err := errors.New(msgErrorResult)
		return result, http.StatusBadRequest, err
	}
	return result, 0, nil
}

// UpdateWarehousePrimary function for update primary merchant address
func (m *HTTPMerchantHandler) UpdateWarehousePrimary(c echo.Context) error {
	// extract user ID
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	params := model.ParameterPrimaryWarehouse{}
	params.MemberID = memberID
	params.AddressID = c.Param(TextAddressID)

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	// Update primary address usecase process
	result := <-m.MerchantAddressUseCase.UpdatePrimaryWarehouseAddress(newCtx, params)
	if result.Error != nil {
		return shared.NewHTTPResponse(result.HTTPStatus, result.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Warehouse location Updated").JSON(c)
}

// GetWarehouse function for getting list of warehouse address
func (m *HTTPMerchantHandler) GetWarehouse(c echo.Context) error {
	// extract user ID
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	params := model.ParameterWarehouse{
		StrPage:  c.QueryParam("page"),
		StrLimit: c.QueryParam("limit"),
		MemberID: memberID,
		ShowAll:  c.QueryParam("showAll"),
	}

	warehouseResult := <-m.MerchantAddressUseCase.GetWarehouseAddresses(c.Request().Context(), &params)
	if warehouseResult.Error != nil {
		return shared.NewHTTPResponse(warehouseResult.HTTPStatus, warehouseResult.Error.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	warehouse, ok := warehouseResult.Result.(model.ListWarehouse)
	if !ok {
		err := errors.New("result is not list of warehouse")
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	totalPage := math.Ceil(float64(warehouse.TotalData) / float64(params.Limit))

	if len(warehouse.WarehouseData) <= 0 {
		response := shared.NewHTTPResponse(http.StatusOK, messageSuccessGet, make(helper.EmptySlice, 0))
		return response.JSON(c)
	}

	meta := shared.Meta{
		Page:         params.Page,
		Limit:        params.Limit,
		TotalRecords: warehouse.TotalData,
		TotalPages:   int(totalPage),
	}
	return shared.NewHTTPResponse(http.StatusOK, messageSuccessGet, warehouse.WarehouseData, meta).JSON(c)
}

// GetWarehouseDetail function for getting list of merchant address
func (m *HTTPMerchantHandler) GetWarehouseDetail(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	result := model.WarehouseData{}
	warehouse := <-m.MerchantAddressUseCase.GetDetailWarehouseAddress(c.Request().Context(), c.Param(TextAddressID), memberID)
	if warehouse.Error != nil {
		return shared.NewHTTPResponse(warehouse.HTTPStatus, warehouse.Error.Error()).JSON(c)
	}

	result, ok := warehouse.Result.(model.WarehouseData)
	if !ok {
		err := errors.New(msgErrorResult)
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, messageSuccessGet, result).JSON(c)
}

// DeleteWarehouse function for removing merchant address
func (m *HTTPMerchantHandler) DeleteWarehouse(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	// Delete warehouse address usecase process
	result := <-m.MerchantAddressUseCase.DeleteWarehouseAddress(newCtx, c.Param(TextAddressID), memberID)
	if result.Error != nil {
		return shared.NewHTTPResponse(result.HTTPStatus, result.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Merchant Address Deleted").JSON(c)
}

func (m *HTTPMerchantHandler) UpdateMerchant(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	updatePayload := &model.B2CMerchantCreateInput{}
	if err := c.Bind(&updatePayload); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	userInfo := &model.MerchantUserAttribute{}
	userInfo.UserID = memberID
	userInfo.UserIP = c.RealIP()

	updatePayload.UserID = memberID

	mErr := jsonschema.ValidateTemp("self_update_merchant_params_v2", updatePayload)
	if mErr != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, mErr.Error()).JSON(c)
	}

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	saveResult := <-m.MerchantUseCase.SelfUpdateMerchant(newCtx, updatePayload, userInfo)
	if saveResult.Error != nil {
		return shared.NewHTTPResponse(saveResult.HTTPStatus, saveResult.Error.Error()).JSON(c)
	}

	resultResponse, ok := saveResult.Result.(model.B2CMerchantDataV2)
	if !ok {
		err := errors.New(msgErrorResult)
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Merchant Update Response", resultResponse).JSON(c)
}

func (m *HTTPMerchantHandler) UpdateMerchantPartial(c echo.Context) error {
	merchantMemberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	updatePayload := &model.B2CMerchantCreateInput{}
	if err := c.Bind(&updatePayload); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	userInfo := &model.MerchantUserAttribute{}
	userInfo.UserID = merchantMemberID
	userInfo.UserIP = c.RealIP()

	updatePayload.UserID = merchantMemberID

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	saveResult := <-m.MerchantUseCase.SelfUpdateMerchantPartial(newCtx, updatePayload, userInfo)
	if saveResult.Error != nil {
		return shared.NewHTTPResponse(saveResult.HTTPStatus, saveResult.Error.Error()).JSON(c)
	}

	resultResponse, ok := saveResult.Result.(model.B2CMerchantDataV2)
	if !ok {
		err := errors.New(msgErrorResult)
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	newResultResponse := model.B2CMerchantDataPatial(resultResponse)
	return shared.NewHTTPResponse(http.StatusOK, "Merchant Update Response", newResultResponse).JSON(c)
}

func (m *HTTPMerchantHandler) ChangeMerchantName(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	updatePayload := &model.B2CMerchantCreateInput{}
	if err := c.Bind(&updatePayload); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	userInfo := &model.MerchantUserAttribute{}
	userInfo.UserID = memberID
	userInfo.UserIP = c.RealIP()

	updatePayload.UserID = memberID

	mErr := jsonschema.ValidateTemp("self_update_merchant_partial_params_v2", updatePayload)
	if mErr != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, mErr.Error()).JSON(c)
	}

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	saveResult := <-m.MerchantUseCase.ChangeMerchantName(newCtx, updatePayload, userInfo)
	if saveResult.Error != nil {
		return shared.NewHTTPResponse(saveResult.HTTPStatus, saveResult.Error.Error()).JSON(c)
	}

	resultResponse, ok := saveResult.Result.(model.B2CMerchantDataV2)
	if !ok {
		err := errors.New(msgErrorResult)
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Merchant Change Name Response", resultResponse).JSON(c)
}

// AddEmployee ...
func (h *HTTPMerchantHandler) AddEmployee(c echo.Context) error {
	var payload struct {
		Email     string `json:"email" form:"email"`
		FirstName string `json:"firstName" form:"firstName"`
		Token     string `json:"token" form:"token"`
	}
	if err := c.Bind(&payload); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	if payload.Email == "" || payload.FirstName == "" {
		return shared.NewHTTPResponse(http.StatusBadRequest, "Email and Firstname can't be empty").JSON(c)
	}

	authorization := c.Request().Header.Get(echo.HeaderAuthorization)
	var token string
	if split := strings.Split(authorization, " "); len(split) > 1 {
		token = split[1]
	}

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	result := <-h.MerchantUseCase.AddEmployee(newCtx, token, payload.Email, payload.FirstName)
	if result.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, result.Error.Error()).JSON(c)
	}

	payload.Token = result.Result.(string)
	return shared.NewHTTPResponse(http.StatusOK, "Success Merchant employee invited", payload).JSON(c)
}

// ListEmployee ...
func (m *HTTPMerchantHandler) ListEmployee(c echo.Context) error {
	params := model.QueryMerchantEmployeeParameters{}
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

	result := <-m.MerchantUseCase.GetAllMerchantEmployee(c.Request().Context(), token, &params)
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

// GetEmployee ...
func (m *HTTPMerchantHandler) GetEmployee(c echo.Context) error {
	params := model.QueryMerchantEmployeeParameters{}

	memberId := c.Param("memberId")
	if memberId == "" {
		return shared.NewHTTPResponse(http.StatusBadRequest, "required memberId").JSON(c)
	}
	params.MemberID = memberId

	authorization := c.Request().Header.Get(echo.HeaderAuthorization)
	var token string
	if split := strings.Split(authorization, " "); len(split) > 1 {
		token = split[1]
	}

	result := <-m.MerchantUseCase.GetMerchantEmployee(c.Request().Context(), token, &params)
	if result.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, result.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success get merchant employee", result.Result).JSON(c)
}

// UpdateEmployee ...
func (m *HTTPMerchantHandler) UpdateEmployee(c echo.Context) error {
	params := model.QueryMerchantEmployeeParameters{}
	if err := c.Bind(&params); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	memberId := c.Param("memberId")
	if memberId == "" {
		return shared.NewHTTPResponse(http.StatusBadRequest, "required memberId").JSON(c)
	}
	params.MemberID = memberId

	authorization := c.Request().Header.Get(echo.HeaderAuthorization)
	var token string
	if split := strings.Split(authorization, " "); len(split) > 1 {
		token = split[1]
	}

	result := <-m.MerchantUseCase.UpdateMerchantEmployee(c.Request().Context(), token, &params)
	if result.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, result.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success update merchant employee", result.Result).JSON(c)
}

// ResendEmailEmployee ...
func (m *HTTPMerchantHandler) ResendEmailEmployee(c echo.Context) error {
	var payload struct {
		Email string `json:"email" form:"email"`
		Token string `json:"token" form:"token"`
	}
	if err := c.Bind(&payload); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	authorization := c.Request().Header.Get(echo.HeaderAuthorization)
	var token string
	if split := strings.Split(authorization, " "); len(split) > 1 {
		token = split[1]
	}

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	result := <-m.MerchantUseCase.AddEmployee(newCtx, token, payload.Email, "{resend-email}")
	if result.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, result.Error.Error()).JSON(c)
	}

	payload.Token = result.Result.(string)
	return shared.NewHTTPResponse(http.StatusOK, "Success resend email merchant employee", payload).JSON(c)
}

func (m *HTTPMerchantHandler) clearRejectMerchantUpgrade(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	uaMerchantUpgrade := &model.MerchantUserAttribute{}
	uaMerchantUpgrade.UserID = memberID
	uaMerchantUpgrade.UserIP = c.RealIP()

	ctxReq := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	ctxReq = context.WithValue(ctxReq, middleware.ContextKeyClientIP, c.RealIP())

	cleanRejectUpgrade := <-m.MerchantUseCase.ClearRejectUpgrade(ctxReq, memberID, uaMerchantUpgrade)
	if cleanRejectUpgrade.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, cleanRejectUpgrade.Error.Error()).JSON(c)
	}

	merchant, ok := cleanRejectUpgrade.Result.(model.B2CMerchantDataV2)
	if !ok {
		return shared.NewHTTPResponse(http.StatusBadRequest, msgErrorResult).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success cleared reject upgrade status", merchant).JSON(c)
}
