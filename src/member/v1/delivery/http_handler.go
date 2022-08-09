package delivery

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"os"
	"strings"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/jsonschema"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/member/v1/usecase"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/google/jsonapi"
	"github.com/labstack/echo"
	"golang.org/x/net/context"
)

const (
	requestFrom = "requestFrom"
)

// HTTPMemberHandler model
type HTTPMemberHandler struct {
	MemberUseCase usecase.MemberUseCase
}

// NewHTTPHandler function for initialise *HTTPAuthHandler
func NewHTTPHandler(memberUseCase usecase.MemberUseCase) *HTTPMemberHandler {
	return &HTTPMemberHandler{MemberUseCase: memberUseCase}
}

// Mount function for mounting routes
func (h *HTTPMemberHandler) Mount(group *echo.Group) {
	group.POST("/register", h.RegisterMember)
	group.POST("/forgot-password", h.ForgotPassword)
	group.POST("/activation", h.ActivationMember)
	group.POST("/validate-token", h.ValidateToken)
	group.POST("/change-password", h.ChangeForgotPassword)
	group.POST("/activate-new-password", h.ActivateNewPassword)
}

// MountMe function for mounting routes
func (h *HTTPMemberHandler) MountMe(group *echo.Group) {
	group.GET("", h.FindDetailProfile)
	group.PUT("", h.EditProfile)
	group.POST("/change-password", h.UpdatePassword)
	group.POST("/new-password", h.AddPassword)
}

// MountAdmin function for mounting admin endpoints
func (h *HTTPMemberHandler) MountAdmin(group *echo.Group) {
	group.POST("/member/import", h.ImportMember)
	group.PUT("/member/regenerate-token/:memberID", h.RegenerateToken)
	group.POST("/member", h.AddNewMember)
	group.GET("/member", h.GetMembers)

	group.GET("/member/:memberID", h.GetDetailMember)
	group.POST("/member/migrate", h.MigrateData)
	// group.GET("/member/migrate_legacy/:memberID", h.MigrateLegacyData)
}

// MountMember function for mounting member endpoints using for internal cms login
func (h *HTTPMemberHandler) MountMember(group *echo.Group) {
	group.PUT("/:memberID", h.UpdateMember)
}

// GetMember function for get data member
func (h *HTTPMemberHandler) GetMember(c echo.Context, memberID string) model.GetMemberResult {
	ctxReq := c.Request().Context()
	ctxReq = context.WithValue(ctxReq, shared.ContextKey(helper.TextToken), c.Request().Header.Get(echo.HeaderAuthorization))
	memberResult := <-h.MemberUseCase.GetDetailMemberByID(ctxReq, memberID)
	if memberResult.Error != nil {
		if memberResult.Error == fmt.Errorf(helper.ErrorDataNotFound, "member") {
			memberResult.Error = errors.New(helper.ErrorUnauthorized)
			return model.GetMemberResult{Error: memberResult.Error, HTTPStatus: http.StatusUnauthorized, Scope: "get_detail_member_nofound"}
		}

		if memberResult.HTTPStatus > 0 {
			return model.GetMemberResult{Error: memberResult.Error, HTTPStatus: memberResult.HTTPStatus, Scope: "get_detail_member"}
		}
		return model.GetMemberResult{Error: memberResult.Error, HTTPStatus: http.StatusUnauthorized, Scope: "get_detail_member"}
	}

	member, ok := memberResult.Result.(model.Member)
	if !ok {
		err := errors.New("result is not member")
		return model.GetMemberResult{Error: err, HTTPStatus: http.StatusInternalServerError, Scope: helper.ScopeParseResponse}
	}

	return model.GetMemberResult{Result: member}
}

// FindDetailProfile function for getting detail member based on member id
func (h *HTTPMemberHandler) FindDetailProfile(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// get member data
	memberResult := h.GetMember(c, memberID)
	if memberResult.Error != nil {
		return echo.NewHTTPError(memberResult.HTTPStatus, memberResult.Error.Error())
	}
	member := memberResult.Result

	// adjust value of object
	if member.HasPassword {
		member.Password = ""
		member.Salt = ""
	}

	payload, err := helper.MarshalConvertOnePayload(&member)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrorOccured)
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(payload)
}

// EditProfile function for updating profile data
func (h *HTTPMemberHandler) EditProfile(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// adjust data from form value
	member := h.generateFormValueMember(c)
	member.ID = memberID
	member.StatusString = model.ActiveString
	member.Type = helper.TextUpdate
	member.ModifiedBy = memberID

	validationError := jsonschema.ValidateTemp("update_member_params_v1", member)
	if validationError != nil {
		return echo.NewHTTPError(http.StatusBadRequest, validationError.Error())
	}
	headerAuth := c.Request().Header.Get(helper.TextAuthorization)

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, headerAuth)
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	newCtx = shared.SetDataToContext(newCtx, shared.ContextKey(helper.TextToken), headerAuth)
	saveResultV1 := <-h.MemberUseCase.UpdateDetailMemberByID(newCtx, member)
	if saveResultV1.Error != nil {
		return echo.NewHTTPError(saveResultV1.HTTPStatus, saveResultV1.Error.Error())
	}

	// adjust value of object
	if member.HasPassword {
		member.Password = ""
		member.Salt = ""
	}

	payload, err := helper.MarshalConvertOnePayload(&member)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrorOccured)
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(payload)
}

// UpdatePassword function for updating password
func (h *HTTPMemberHandler) UpdatePassword(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	oldPassword := c.FormValue(model.FieldOldPassword)
	newPassword := c.FormValue(model.FieldNewPassword)

	authorization := c.Request().Header.Get(echo.HeaderAuthorization)
	var token string
	if split := strings.Split(authorization, " "); len(split) > 1 {
		token = split[1]
	}

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	passResultV1 := <-h.MemberUseCase.UpdatePassword(newCtx, token, memberID, oldPassword, newPassword)

	if passResultV1.Error != nil {
		return echo.NewHTTPError(passResultV1.HTTPStatus, passResultV1.Error.Error())
	}

	result, ok := passResultV1.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(helper.ErrorResultNotProper)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	response, err := helper.MarshalConvertOnePayload(&result)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrorOccured)
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(response)
}

// AddPassword function for adding new password for logged in member
func (h *HTTPMemberHandler) AddPassword(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// get member data
	memberResult := h.GetMember(c, memberID)
	if memberResult.Error != nil {
		return echo.NewHTTPError(memberResult.HTTPStatus, memberResult.Error.Error())
	}

	member := memberResult.Result

	// check whether password exists or not
	if len(member.Password) > 0 && len(member.Salt) > 0 {
		err := errors.New("you are not allowed to add new password")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// append the member object
	member.NewPassword = c.FormValue(model.FieldPassword)
	member.RePassword = c.FormValue(model.FieldRePassword)

	passResult := <-h.MemberUseCase.AddNewPassword(c.Request().Context(), member)

	if passResult.Error != nil {
		return echo.NewHTTPError(passResult.HTTPStatus, passResult.Error.Error())
	}

	res := model.SuccessResponse{
		ID:          golib.RandomString(8),
		Message:     helper.SuccessMessage,
		HasPassword: true,
		Email:       strings.ToLower(member.Email),
		FirstName:   member.FirstName,
		LastName:    member.LastName,
	}
	payload, err := helper.MarshalConvertOnePayload(&res)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrorOccured)
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(payload)
}

// RegisterMember function for registering new member
func (h *HTTPMemberHandler) RegisterMember(c echo.Context) error {
	member := &model.Member{}
	member.FirstName = c.FormValue(model.FieldFirstName)
	member.LastName = c.FormValue(model.FieldLastName)
	member.Email = c.FormValue(model.FieldEmail)
	member.NewPassword = c.FormValue(model.FieldPassword)
	member.RePassword = c.FormValue(model.FieldRePassword)
	member.GenderString = c.FormValue(model.FieldGender)
	member.BirthDateString = c.FormValue(model.FieldDOB)
	member.Mobile = c.FormValue(model.FieldMobile)
	member.Type = "register"

	// check member existence first when error not happens
	// or the error is not equal data not found
	checkResult := <-h.MemberUseCase.CheckEmailAndMobileExistence(c.Request().Context(), member)
	if checkResult.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, checkResult.Error.Error())
	}

	registerResult := <-h.MemberUseCase.RegisterMember(c.Request().Context(), member)
	if registerResult.Error != nil {
		return echo.NewHTTPError(registerResult.HTTPStatus, registerResult.Error.Error())
	}

	result, ok := registerResult.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(helper.ErrorResultNotProper)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	responseRegister, err := helper.MarshalConvertOnePayload(&result)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrorOccured)
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(responseRegister)
}

// ActivationMember function for activating inactivate member
func (h *HTTPMemberHandler) ActivationMember(c echo.Context) error {
	token := c.FormValue("token")

	activateResult := <-h.MemberUseCase.ActivateMember(c.Request().Context(), token, "")
	if activateResult.Error != nil {
		return echo.NewHTTPError(activateResult.HTTPStatus, activateResult.Error.Error())
	}

	response, ok := activateResult.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(helper.ErrorResultNotProper)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	activationResponse, err := helper.MarshalConvertOnePayload(&response)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrorOccured)
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(activationResponse)
}

// ForgotPassword function for requesting token of forgot password
func (h *HTTPMemberHandler) ForgotPassword(c echo.Context) error {
	email := c.FormValue(model.FieldEmail)

	passResult := <-h.MemberUseCase.ForgotPassword(c.Request().Context(), email)
	if passResult.Error != nil {
		return echo.NewHTTPError(passResult.HTTPStatus, passResult.Error.Error())
	}

	res, ok := passResult.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(helper.ErrorResultNotProper)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	env, ok := os.LookupEnv("ENV")
	if !ok {
		env = "PROD"
	}
	if env != "DEV" {
		res.Token = "" //Hide Token
	}

	payload, err := helper.MarshalConvertOnePayload(&res)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrorOccured)
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(payload)
}

// ValidateToken function for validating token of forgot password
func (h *HTTPMemberHandler) ValidateToken(c echo.Context) error {
	token := c.FormValue("token")

	validateResult := <-h.MemberUseCase.ValidateToken(c.Request().Context(), token)
	if validateResult.Error != nil {
		return echo.NewHTTPError(validateResult.HTTPStatus, validateResult.Error.Error())
	}

	responseValidate, ok := validateResult.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(helper.ErrorResultNotProper)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	responseToken, err := helper.MarshalConvertOnePayload(&responseValidate)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrorOccured)
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(responseToken)
}

// ChangeForgotPassword function for creating new password after requesting token
func (h *HTTPMemberHandler) ChangeForgotPassword(c echo.Context) error {
	token := c.FormValue("token")
	newPassword := c.FormValue(model.FieldNewPassword)
	rePassword := c.FormValue(model.FieldRePassword)
	requestFrom := c.FormValue(requestFrom)

	changePassResult := <-h.MemberUseCase.ChangeForgotPassword(c.Request().Context(), token, newPassword, rePassword, requestFrom)
	if changePassResult.Error != nil {
		return echo.NewHTTPError(changePassResult.HTTPStatus, changePassResult.Error.Error())
	}

	changePassResponse, ok := changePassResult.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(helper.ErrorResultNotProper)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	finalResponse, err := helper.MarshalConvertOnePayload(&changePassResponse)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrorOccured)
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(finalResponse)
}

// ActivateNewPassword function for activating new password for member who is registered from dolphin
func (h *HTTPMemberHandler) ActivateNewPassword(c echo.Context) error {
	token := c.FormValue("token")
	newPassword := c.FormValue(model.FieldNewPassword)
	rePassword := c.FormValue(model.FieldRePassword)

	activateResult := <-h.MemberUseCase.ActivateNewPassword(c.Request().Context(), token, newPassword, rePassword)
	if activateResult.Error != nil {
		return echo.NewHTTPError(activateResult.HTTPStatus, activateResult.Error.Error())
	}

	responseNewPass, ok := activateResult.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(helper.ErrorResultNotProper)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	finalResponseActivate, err := helper.MarshalConvertOnePayload(&responseNewPass)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrorOccured)
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(finalResponseActivate)
}

// RegenerateToken function for regenerating token for activation
func (h *HTTPMemberHandler) RegenerateToken(c echo.Context) error {
	// get member data
	memberResult := h.GetMember(c, c.Param(helper.TextMemberID))
	if memberResult.Error != nil {
		return echo.NewHTTPError(memberResult.HTTPStatus, memberResult.Error.Error())
	}

	generateResult := <-h.MemberUseCase.RegenerateToken(c.Request().Context(), memberResult.Result)
	if generateResult.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, generateResult.Error.Error())
	}

	responseGenerate, ok := generateResult.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(helper.ErrorResultNotProper)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	finalResponse, err := helper.MarshalConvertOnePayload(&responseGenerate)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrorOccured)
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(finalResponse)
}

// AddNewMember function for adding new member from other pack
func (h *HTTPMemberHandler) AddNewMember(c echo.Context) error {
	// adjust data from form value
	member := h.generateFormValueMember(c)
	member.Email = c.FormValue(model.FieldEmail)
	member.Password = c.FormValue(model.FieldPassword)
	member.RePassword = c.FormValue(model.FieldRePassword)
	member.SignUpFrom = c.FormValue("signUpFrom")
	member.Type = "add"

	// check member existence first when error not happens
	// or the error is not equal data not found
	checkResult := <-h.MemberUseCase.CheckEmailAndMobileExistence(c.Request().Context(), &member)
	if checkResult.Error != nil {
		return echo.NewHTTPError(http.StatusBadRequest, checkResult.Error.Error())
	}

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	addMemberResult := <-h.MemberUseCase.RegisterMember(newCtx, &member)
	if addMemberResult.Error != nil {
		return echo.NewHTTPError(addMemberResult.HTTPStatus, addMemberResult.Error.Error())
	}

	result, ok := addMemberResult.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(helper.ErrorResultNotProper)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	finalResponseAddMember, err := helper.MarshalConvertOnePayload(&result)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrorOccured)
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(finalResponseAddMember)
}

// UpdateMember function for updating member from other pack
func (h *HTTPMemberHandler) UpdateMember(c echo.Context) error {
	isAdmin := middleware.ExtractClaimsIsAdmin(c)
	paramMemberID := c.Param(helper.TextMemberID)
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, err.Error())
	}

	if paramMemberID != memberID && isAdmin != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, isAdmin.Error())
	}

	// adjust data from form value
	member := h.generateFormValueMember(c)
	member.ID = c.Param(helper.TextMemberID)
	// optional
	member.IsStaffString = c.FormValue("isStaff")
	member.IsAdminString = c.FormValue("isAdmin")
	member.StatusString = c.FormValue(model.FieldStatus)
	member.IsActiveString = c.FormValue("isActive")
	member.Type = helper.TextUpdate
	member.UpdateFrom = model.Dolphin

	mErr := jsonschema.ValidateTemp("update_member_params_v1", member)
	if mErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest, mErr.Error())
	}
	headerAuth := c.Request().Header.Get(helper.TextAuthorization)

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, headerAuth)
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	newCtx = shared.SetDataToContext(newCtx, shared.ContextKey(helper.TextToken), headerAuth)
	saveResult := <-h.MemberUseCase.UpdateDetailMemberByID(newCtx, member)
	if saveResult.Error != nil {
		return echo.NewHTTPError(saveResult.HTTPStatus, saveResult.Error.Error())
	}

	member, ok := saveResult.Result.(model.Member)
	if !ok {
		err := errors.New("result is not member")
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	// adjust value of object
	if !member.HasPassword {
		member.Password = ""
		member.Salt = ""
	}

	payload, err := helper.MarshalConvertOnePayload(&member)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrorOccured)
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(payload)
}

// GetDetailMember function for getting detail member based on id
func (h *HTTPMemberHandler) GetDetailMember(c echo.Context) error {
	// get member data
	memberResult := h.GetMember(c, c.Param(helper.TextMemberID))
	if memberResult.Error != nil {
		return echo.NewHTTPError(memberResult.HTTPStatus, memberResult.Error.Error())
	}

	payload, err := helper.MarshalConvertOnePayload(&memberResult.Result)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrorOccured)
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(payload)
}

// GetMembers function for getting list of members
func (h *HTTPMemberHandler) GetMembers(c echo.Context) error {
	params := model.Parameters{
		Query:    c.QueryParam("query"),
		StrPage:  c.QueryParam("page"),
		StrLimit: c.QueryParam("limit"),
		Sort:     c.QueryParam("sort"),
		OrderBy:  c.QueryParam("orderBy"),
		Status:   c.QueryParam(model.FieldStatus),
		Email:    c.QueryParam(model.FieldEmail),
	}

	memberResult := <-h.MemberUseCase.GetListMembers(c.Request().Context(), &params)
	if memberResult.Error != nil {
		return echo.NewHTTPError(memberResult.HTTPStatus, memberResult.Error.Error())
	}

	member, ok := memberResult.Result.(model.ListMembers)
	if !ok {
		err := errors.New(helper.ErrorResultNotProper)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	member.ID = golib.RandomString(8)
	member.Name = "list of members"

	payload, err := helper.MarshalConvertOnePayload(&member)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrorOccured)
	}

	totalPage := math.Ceil(float64(member.TotalData) / float64(params.Limit))

	payload.Meta = &jsonapi.Meta{
		"page":      params.Page,
		"limit":     params.Limit,
		"totalData": member.TotalData,
		"totalPage": int(totalPage),
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(payload)
}

// MigrateData function for migrating data from squid
func (h *HTTPMemberHandler) MigrateData(c echo.Context) error {
	// bind and get the body
	members := &model.Members{}
	if err := c.Bind(members); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	migrateResult := <-h.MemberUseCase.MigrateMember(c.Request().Context(), members)
	if migrateResult.Error != nil {

		// if error data exists then return its array
		if len(migrateResult.ErrorData) > 0 {
			data := struct {
				Message string              `json:"message"`
				Data    []model.MemberError `json:"data"`
			}{migrateResult.Error.Error(), migrateResult.ErrorData}

			return c.JSON(http.StatusInternalServerError, data)
		}
		return echo.NewHTTPError(http.StatusInternalServerError, migrateResult.Error.Error())
	}

	res := model.SuccessResponse{
		ID:      golib.RandomString(8),
		Message: helper.SuccessMessage,
	}
	payload, err := helper.MarshalConvertOnePayload(&res)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, helper.ErrorOccured)
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(payload)
}

// generateFormValueMember function for generate data Form Value
func (h *HTTPMemberHandler) generateFormValueMember(c echo.Context) model.Member {
	memberAddressV1 := model.Address{}
	memberValue := model.Member{}

	memberValue.FirstName = c.FormValue(model.FieldFirstName)
	memberValue.LastName = c.FormValue(model.FieldLastName)
	memberValue.Mobile = c.FormValue(model.FieldMobile)
	memberValue.Phone = c.FormValue(model.FieldPhone)
	memberValue.Ext = c.FormValue(model.FieldExt)
	memberValue.GenderString = c.FormValue(model.FieldGender)
	memberValue.BirthDateString = c.FormValue(model.FieldDOB)

	memberAddressV1.Street1 = c.FormValue(model.FieldStreet1)
	memberAddressV1.Street2 = c.FormValue(model.FieldStreet2)
	memberAddressV1.ZipCode = c.FormValue(model.FieldPostalCode)
	memberAddressV1.SubDistrictID = c.FormValue(model.FieldSubDistrictID)
	memberAddressV1.SubDistrict = c.FormValue(model.FieldSubDistrictName)
	memberAddressV1.DistrictID = c.FormValue(model.FieldDistrictID)
	memberAddressV1.District = c.FormValue(model.FieldDistrictName)
	memberAddressV1.CityID = c.FormValue(model.FieldCityID)
	memberAddressV1.City = c.FormValue(model.FieldCityName)
	memberAddressV1.ProvinceID = c.FormValue(model.FieldProvinceID)
	memberAddressV1.Province = c.FormValue(model.FieldProvinceName)

	memberValue.Address = memberAddressV1
	return memberValue
}
