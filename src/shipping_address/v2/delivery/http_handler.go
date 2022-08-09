package delivery

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/Bhinneka/golib/jsonschema"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/Bhinneka/user-service/src/shipping_address/v2/model"
	"github.com/Bhinneka/user-service/src/shipping_address/v2/usecase"
	"github.com/labstack/echo"
)

const (
	scopeSaveShippingAddress     = "save_shipping_address"
	scopeExtractMemberID         = "extract_member_id"
	scopeParseShippingAddress    = "parse_shipping_address"
	tempShippingAddress          = "add_shipping_address_params_v2"
	scopeValidateShippingAddress = "validate_add_shipping_address_params_v2"
	messageErrorResult           = "result is not proper response"
	messageSuccess               = "Shipping Address Success"
	mobileText                   = "mobile"
	shippingIDParams             = "/:shippingId"
	meShippingIDParams           = "/me/:shippingId"
	street1                      = "street1"
	street2                      = "street2"
	postalCode                   = "postalCode"
	subDistrictID                = "subDistrictId"
	subDistrictName              = "subDistrictName"
	districtID                   = "districtId"
	districtName                 = "districtName"
	cityID                       = "cityId"
	cityName                     = "cityName"
	provinceID                   = "provinceId"
	provinceName                 = "provinceName"
	shippingID                   = "shippingId"
	scopeClaimToken              = "get_claims_from_token"
	errLatlong                   = "Latitude/longitude data type must be float"
)

// HTTPShippingAddressHandler model
type HTTPShippingAddressHandler struct {
	ShippingAddressUseCase usecase.ShippingAddressUseCase
}

// NewHTTPHandler function for initialise *HTTPAuthHandler
func NewHTTPHandler(ShippingAddressUseCase usecase.ShippingAddressUseCase) *HTTPShippingAddressHandler {
	return &HTTPShippingAddressHandler{ShippingAddressUseCase: ShippingAddressUseCase}
}

// MountMe function for mounting routes
func (s *HTTPShippingAddressHandler) MountMe(group *echo.Group) {
	group.POST("/me", s.AddShippingAddressMe)
	group.PUT(meShippingIDParams, s.UpdateShippingAddressMe)
	group.DELETE(meShippingIDParams, s.DeleteShippingAddressMe)
	group.GET("/me", s.GetShippingAddressMe)
	group.GET(meShippingIDParams, s.GetShippingAddressDetailMe)
	group.GET("/me/primary", s.GetShippingAddressPrimaryMe)
	group.PUT("/me/:shippingId/set-primary", s.UpdateIsPrimaryMe)
}

// MountShippingAddress function for mounting endpoints
func (s *HTTPShippingAddressHandler) MountShippingAddress(group *echo.Group) {
	group.POST("", s.AddShippingAddress)
	group.PUT(shippingIDParams, s.UpdateShippingAddress)
	group.DELETE(shippingIDParams, s.DeleteShippingAddress)
	group.GET("", s.GetShippingAddress)
	group.GET(shippingIDParams, s.GetShippingAddressDetail)
	group.PUT("/:shippingId/set-primary", s.UpdateIsPrimary)
}

// AddShippingAddressMe function for add new shipping address
func (s *HTTPShippingAddressHandler) AddShippingAddressMe(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	result, errCode, err := s.SaveUpdateShippingAddress(c, "add", helper.TextMe, memberID)
	if err != nil {
		return shared.NewHTTPResponse(errCode, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusCreated, messageSuccess, result).JSON(c)
}

// DeleteShippingAddressMe function for removing shipping address
func (s *HTTPShippingAddressHandler) DeleteShippingAddressMe(c echo.Context) error {
	shippingID := c.Param(shippingID)
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	// Delete shipping address usecase process
	result := <-s.ShippingAddressUseCase.DeleteShippingAddressByID(newCtx, shippingID, memberID)
	if result.Error != nil {
		return shared.NewHTTPResponse(result.HTTPStatus, result.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Shipping Address Deleted").JSON(c)
}

// UpdateShippingAddressMe function for update shipping address
func (s *HTTPShippingAddressHandler) UpdateShippingAddressMe(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	result, errCode, err := s.SaveUpdateShippingAddress(c, helper.TextUpdate, helper.TextMe, memberID)
	if err != nil {
		return shared.NewHTTPResponse(errCode, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, messageSuccess, result).JSON(c)
}

// GetShippingAddressMe function for getting list of shipping address
func (s *HTTPShippingAddressHandler) GetShippingAddressMe(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	params := model.ParametersShippingAddress{
		StrPage:  c.QueryParam("page"),
		StrLimit: c.QueryParam("limit"),
		MemberID: memberID,
	}

	shippingAddressResults := <-s.ShippingAddressUseCase.GetListShippingAddress(c.Request().Context(), &params)
	if shippingAddressResults.Error != nil {
		return shared.NewHTTPResponse(shippingAddressResults.HTTPStatus, shippingAddressResults.Error.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	shippingAddress, ok := shippingAddressResults.Result.(model.ListShippingAddress)
	if !ok {
		err := errors.New("result is not list of shipping address")
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	totalPage := math.Ceil(float64(shippingAddress.TotalData) / float64(params.Limit))

	if len(shippingAddress.ShippingAddress) <= 0 {
		response := shared.NewHTTPResponse(http.StatusOK, messageSuccess, make(helper.EmptySlice, 0))
		return response.JSON(c)
	}

	meta := shared.Meta{
		Page:         params.Page,
		Limit:        params.Limit,
		TotalRecords: shippingAddress.TotalData,
		TotalPages:   int(totalPage),
	}
	return shared.NewHTTPResponse(http.StatusOK, messageSuccess, shippingAddress.ShippingAddress, meta).JSON(c)
}

// GetShippingAddressDetailMe function for getting list of shipping address
func (s *HTTPShippingAddressHandler) GetShippingAddressDetailMe(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	result, errCode, err := s.GetDetail(c, memberID)
	if err != nil {
		return shared.NewHTTPResponse(errCode, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, messageSuccess, result).JSON(c)
}

// UpdateIsPrimaryMe function for update primary shipping address
func (s *HTTPShippingAddressHandler) UpdateIsPrimaryMe(c echo.Context) error {
	// extract user ID
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	params := model.ParamaterPrimaryShippingAddress{}
	params.MemberID = memberID
	params.ShippingID = c.Param(shippingID)

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	// Delete shipping address usecase process
	result := <-s.ShippingAddressUseCase.UpdatePrimaryShippingAddressByID(newCtx, params)
	if result.Error != nil {
		return shared.NewHTTPResponse(result.HTTPStatus, result.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Shipping Address Updated").JSON(c)
}

// AddShippingAddress function for add new shipping address
func (s *HTTPShippingAddressHandler) AddShippingAddress(c echo.Context) error {
	err := middleware.ExtractClaimsIsAdmin(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, err.Error()).JSON(c)
	}

	userID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	result, errCode, err := s.SaveUpdateShippingAddress(c, "add", "", userID)
	if err != nil {
		return shared.NewHTTPResponse(errCode, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusCreated, messageSuccess, result).JSON(c)
}

// GetShippingAddress function for getting list of shipping address
func (s *HTTPShippingAddressHandler) GetShippingAddress(c echo.Context) error {
	err := middleware.ExtractClaimsIsAdmin(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, err.Error()).JSON(c)
	}

	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	params := model.ParametersShippingAddress{
		StrPage:  c.QueryParam("page"),
		StrLimit: c.QueryParam("limit"),
		MemberID: c.QueryParam(helper.TextMemberIDCamel),
		Query:    c.QueryParam("query"),
	}

	if !helper.ValidateMamberID(params.MemberID) {
		err = errors.New("memberId not valid")
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	shippingAddressResult := <-s.ShippingAddressUseCase.GetAllListShippingAddress(c.Request().Context(), &params, memberID)
	if shippingAddressResult.Error != nil {
		return shared.NewHTTPResponse(shippingAddressResult.HTTPStatus, shippingAddressResult.Error.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	shippingAddress, ok := shippingAddressResult.Result.(model.ListShippingAddress)
	if !ok {
		err := errors.New("result is not list of shipping address")
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	totalPage := math.Ceil(float64(shippingAddress.TotalData) / float64(params.Limit))

	if len(shippingAddress.ShippingAddress) <= 0 {
		response := shared.NewHTTPResponse(http.StatusOK, messageSuccess, make(helper.EmptySlice, 0))
		return response.JSON(c)
	}

	meta := shared.Meta{
		Page:         params.Page,
		Limit:        params.Limit,
		TotalRecords: shippingAddress.TotalData,
		TotalPages:   int(totalPage),
	}
	return shared.NewHTTPResponse(http.StatusOK, messageSuccess, shippingAddress.ShippingAddress, meta).JSON(c)
}

// GetShippingAddressDetail function for getting list of shipping address
func (s *HTTPShippingAddressHandler) GetShippingAddressDetail(c echo.Context) error {
	err := middleware.ExtractClaimsIsAdmin(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, err.Error()).JSON(c)
	}
	result, errCode, err := s.GetDetail(c, "")
	if err != nil {
		return shared.NewHTTPResponse(errCode, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, messageSuccess, result).JSON(c)
}

// GetDetail function for getting list of shipping address
func (s *HTTPShippingAddressHandler) GetDetail(c echo.Context, memberID string) (model.ShippingAddressData, int, error) {
	result := model.ShippingAddressData{}
	shippingAddressResult := <-s.ShippingAddressUseCase.GetDetailShippingAddress(c.Request().Context(), c.Param(shippingID), memberID)
	if shippingAddressResult.Error != nil {
		return result, shippingAddressResult.HTTPStatus, shippingAddressResult.Error
	}

	result, ok := shippingAddressResult.Result.(model.ShippingAddressData)
	if !ok {
		err := errors.New(messageErrorResult)
		return result, http.StatusBadRequest, err
	}
	return result, 0, nil
}

// UpdateShippingAddress function for update shipping address
func (s *HTTPShippingAddressHandler) UpdateShippingAddress(c echo.Context) error {
	err := middleware.ExtractClaimsIsAdmin(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, err.Error()).JSON(c)
	}

	userID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	result, errCode, err := s.SaveUpdateShippingAddress(c, helper.TextUpdate, "", userID)
	if err != nil {
		return shared.NewHTTPResponse(errCode, err.Error()).JSON(c)
	}
	return shared.NewHTTPResponse(http.StatusOK, messageSuccess, result).JSON(c)
}

// SaveUpdateShippingAddress function for save update shipping address
func (s *HTTPShippingAddressHandler) SaveUpdateShippingAddress(c echo.Context, action, source, memberID string) (model.ShippingAddressData, int, error) {
	var err error
	result := model.ShippingAddressData{}
	address := model.ShippingAddressData{}
	address.Name = c.FormValue("name")
	address.Mobile = c.FormValue(mobileText)
	address.Phone = c.FormValue("phone")
	address.SubDistrictID = c.FormValue(subDistrictID)
	address.SubDistrictName = c.FormValue(subDistrictName)
	address.DistrictID = c.FormValue(districtID)
	address.DistrictName = c.FormValue(districtName)
	address.CityID = c.FormValue(cityID)
	address.CityName = c.FormValue(cityName)
	address.ProvinceID = c.FormValue(provinceID)
	address.ProvinceName = c.FormValue(provinceName)
	address.PostalCode = c.FormValue(postalCode)
	address.Street1 = helper.ValidateHTML(c.FormValue(street1))
	address.Street2 = helper.ValidateHTML(c.FormValue(street2))
	address.Ext = c.FormValue("ext")
	address.Label = c.FormValue("label")
	latitude := c.FormValue("latitude")
	longitude := c.FormValue("longitude")

	if latitude != "" {
		address.Latitude, err = strconv.ParseFloat(latitude, 8)
		if err != nil {
			return result, http.StatusBadRequest, errors.New(errLatlong)
		}
	}

	if longitude != "" {
		address.Longitude, err = strconv.ParseFloat(longitude, 8)
		if err != nil {
			return result, http.StatusBadRequest, errors.New(errLatlong)
		}
	}

	if action == helper.TextUpdate {
		address.ID = c.Param(shippingID)
	}

	if source == helper.TextMe {
		address.MemberID = memberID
	} else {
		address.MemberID = c.FormValue(helper.TextMemberIDCamel)
		address.ModifiedBy = memberID
	}

	// validate schema fields
	mErr := jsonschema.ValidateTemp(tempShippingAddress, address)
	if mErr != nil {
		return result, http.StatusBadRequest, mErr
	}
	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())

	var saveResult usecase.ResultUseCase
	if action == helper.TextUpdate {
		// update shipping address usecase process
		saveResult = <-s.ShippingAddressUseCase.UpdateShippingAddress(newCtx, address)
	} else {
		// Add shipping address usecase process
		saveResult = <-s.ShippingAddressUseCase.AddShippingAddress(newCtx, address)
	}

	if saveResult.Error != nil {
		return result, saveResult.HTTPStatus, saveResult.Error
	}

	result, ok := saveResult.Result.(model.ShippingAddressData)
	if !ok {
		err := errors.New(messageErrorResult)
		return result, http.StatusBadRequest, err
	}
	return result, 0, nil
}

// DeleteShippingAddress function for removing shipping address
func (s *HTTPShippingAddressHandler) DeleteShippingAddress(c echo.Context) error {
	err := middleware.ExtractClaimsIsAdmin(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, err.Error()).JSON(c)
	}

	shippingID := c.Param(shippingID)
	memberID := c.QueryParam(helper.TextMemberIDCamel)

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	// Delete shipping address usecase process
	result := <-s.ShippingAddressUseCase.DeleteShippingAddressByID(newCtx, shippingID, memberID)
	if result.Error != nil {
		return shared.NewHTTPResponse(result.HTTPStatus, result.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Shipping Address Deleted").JSON(c)
}

// UpdateIsPrimary function for update primary shipping address
func (s *HTTPShippingAddressHandler) UpdateIsPrimary(c echo.Context) error {
	err := middleware.ExtractClaimsIsAdmin(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, err.Error()).JSON(c)
	}

	// extract user ID
	userID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	params := model.ParamaterPrimaryShippingAddress{}
	params.MemberID = c.FormValue(helper.TextMemberIDCamel)
	params.ShippingID = c.Param(shippingID)
	params.UserID = userID

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	// Delete shipping address usecase process
	result := <-s.ShippingAddressUseCase.UpdatePrimaryShippingAddressByID(newCtx, params)
	if result.Error != nil {
		return shared.NewHTTPResponse(result.HTTPStatus, result.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Shipping Address Updated").JSON(c)
}

// GetShippingAddressPrimaryMe function for get primary shipping address
func (s *HTTPShippingAddressHandler) GetShippingAddressPrimaryMe(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	shippingAddressResult := <-s.ShippingAddressUseCase.GetPrimaryShippingAddress(c.Request().Context(), memberID)
	if shippingAddressResult.Error != nil {
		return shared.NewHTTPResponse(shippingAddressResult.HTTPStatus, shippingAddressResult.Error.Error()).JSON(c)
	}

	result, ok := shippingAddressResult.Result.(model.ShippingAddressData)
	if !ok {
		err := errors.New(messageErrorResult)
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, messageSuccess, result).JSON(c)
}
