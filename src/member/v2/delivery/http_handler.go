package delivery

import (
	"context"
	"errors"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/jsonschema"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/member/v1/usecase"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo"
)

const (
	extractMemberID           = "extract_member_id"
	resultIsNotProperResponse = "result is not proper response"
	parseMember               = "parse_members"
	requestFrom               = "requestFrom"
	updateType                = "update"
	memberID                  = "memberID"
	errInvalidInput           = "invalid input"
	URL_CHANGE_PASSWORD       = "/change-password"
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
	group.POST(URL_CHANGE_PASSWORD, h.ChangeForgotPassword)
	group.POST("/activate-new-password", h.ActivateNewPassword)
	group.POST("/validate-email", h.ValidateEmailDomain)
	group.POST("/resend-activation", h.ResendActivation)
	group.POST("/validate-token", h.ValidateToken)
	group.GET("/employee/activation", h.ActivationMerchantEmployee)
}

// MountMe function for mounting me routes
func (h *HTTPMemberHandler) MountMe(group *echo.Group) {
	group.GET("", h.FindDetailProfile)
	group.PUT("", h.EditProfile)
	group.PUT("/change-name", h.ChangeProfileName)
	group.POST(URL_CHANGE_PASSWORD, h.UpdatePassword)
	group.POST("/new-password", h.AddPassword)
	group.POST("/change-profile-picture", h.ChangeProfilePicture)
	group.GET("/mfa", h.GetMFASettings)
	group.GET("/mfa/generate", h.GenerateMFASettings)
	group.POST("/mfa/activation", h.ActivateMFASettings)
	group.DELETE("/mfa", h.DisabledMFASetting)
	group.DELETE("/revoke-all", h.RevokeAllAccess)
	group.GET("/login-activity", h.GetLoginActivity)
	group.GET("/profile-complete", h.GetProfileComplete)
	group.GET("/revoke", h.RevokeAccess)

	// specific for narwhal
	group.GET("/mfa-narwhal", h.GetNarwhalMFASettings)
	group.GET("/mfa-narwhal/generate", h.GenerateNarwhalMFASettings)
	group.POST("/mfa-narwhal/activation", h.ActivateNarwhalMFASettings)
	group.DELETE("/mfa-narwhal", h.DisabledNarwhalMFASetting)

	// sync password
	group.POST("/sync-password", h.SyncPassword)
	group.POST(URL_CHANGE_PASSWORD, h.ChangePassword)
	group.GET("/clients", h.Clients)

}

// MountAdmin function for mounting anonymous membership endpoints
func (h *HTTPMemberHandler) MountAdmin(group *echo.Group) {
	group.POST("/import", h.ImportMember)
	group.GET("/member", h.GetMembers)
	group.PUT("/member/regenerate-token/:memberID", h.RegenerateToken)
	group.GET("/member/:memberID", h.GetDetailMember)
	group.POST("/member/migrate", h.MigrateData)
	group.POST("/bulk-member-send", h.BulkMemberSend)
	group.POST("/member-send", h.MemberSend)
	group.POST("/member", h.AddNewMember) // add new member
}

// MountMember function for mounting member endpoints
// special endpoints for dolphin
func (h *HTTPMemberHandler) MountMember(group *echo.Group) {
	group.PUT("/:memberID", h.UpdateMember)
	group.DELETE("/:memberID/mfa", h.DisabledMFAMember)
	group.DELETE("/:memberID/mfa-narwhal", h.DisabledMFAMemberNarwhal)
}

// GetMember function for getting list of members
func (h *HTTPMemberHandler) GetMember(ctxReq context.Context, memberID string) model.GetMemberResult {

	memberResult := <-h.MemberUseCase.GetDetailMemberByID(ctxReq, memberID)
	if memberResult.Error != nil {
		if memberResult.HTTPStatus > 0 {
			return model.GetMemberResult{Error: memberResult.Error, HTTPStatus: memberResult.HTTPStatus, Scope: "get_detail_member"}
		}

		return model.GetMemberResult{Error: memberResult.Error, HTTPStatus: http.StatusUnauthorized, Scope: "get_detail_member"}
	}

	member, ok := memberResult.Result.(model.Member)
	if !ok {
		err := errors.New("result is not member")
		return model.GetMemberResult{Error: err, HTTPStatus: http.StatusInternalServerError, Scope: "parse_detail_member"}
	}
	return model.GetMemberResult{Result: member}

}

// FindDetailProfile function for getting detail member based on member id
func (h *HTTPMemberHandler) FindDetailProfile(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ctxReq := c.Request().Context()
	ctxReq = context.WithValue(ctxReq, shared.ContextKey(helper.TextToken), c.Request().Header.Get(echo.HeaderAuthorization))
	// get member data
	memberResult := h.GetMember(ctxReq, memberID)
	if memberResult.Error != nil {
		return shared.NewHTTPResponse(memberResult.HTTPStatus, memberResult.Error.Error()).JSON(c)
	}

	member := memberResult.Result

	// adjust value of object
	if member.HasPassword {
		member.Password = ""
		member.Salt = ""
	}

	return shared.NewHTTPResponse(http.StatusOK, "Member Detail Response", member).JSON(c)
}

// EditProfile function for updating profile data
func (h *HTTPMemberHandler) EditProfile(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	// adjust data from form value
	memberAddressEdit := model.Address{}
	member := model.Member{}
	member.ID = memberID
	member.FirstName = c.FormValue(model.FieldFirstName)
	member.LastName = c.FormValue(model.FieldLastName)
	member.Mobile = c.FormValue(model.FieldMobile)
	member.Phone = c.FormValue(model.FieldPhone)
	member.Ext = c.FormValue(model.FieldExt)
	member.GenderString = c.FormValue(model.FieldGender)
	member.BirthDateString = c.FormValue(model.FieldDOB)
	member.StatusString = model.ActiveString

	// flagging for sturgeon update profile
	member.RequestFrom = c.FormValue(requestFrom)
	member.ModifiedBy = memberID

	memberAddressEdit.Street1 = c.FormValue(model.FieldStreet1)
	memberAddressEdit.Street2 = c.FormValue(model.FieldStreet2)
	memberAddressEdit.ZipCode = c.FormValue(model.FieldPostalCode)
	memberAddressEdit.SubDistrictID = c.FormValue(model.FieldSubDistrictID)
	memberAddressEdit.SubDistrict = c.FormValue(model.FieldSubDistrictName)
	memberAddressEdit.DistrictID = c.FormValue(model.FieldDistrictID)
	memberAddressEdit.District = c.FormValue(model.FieldDistrictName)
	memberAddressEdit.CityID = c.FormValue(model.FieldCityID)
	memberAddressEdit.City = c.FormValue(model.FieldCityName)
	memberAddressEdit.ProvinceID = c.FormValue(model.FieldProvinceID)
	memberAddressEdit.Province = c.FormValue(model.FieldProvinceName)

	member.Address = memberAddressEdit
	member.Type = updateType

	var mErr error
	schemaJSON := "update_member_params_v2"
	if member.RequestFrom == model.Sturgeon {
		schemaJSON = "update_member_params_v2_non_address"
	}

	mErr = jsonschema.ValidateTemp(schemaJSON, member)

	if mErr != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, mErr.Error()).JSON(c)
	}

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	newCtx = shared.SetDataToContext(newCtx, shared.ContextKey(helper.TextToken), c.Request().Header.Get(helper.TextAuthorization))
	saveResult := <-h.MemberUseCase.UpdateDetailMemberByID(newCtx, member)
	if saveResult.Error != nil {
		return shared.NewHTTPResponse(saveResult.HTTPStatus, saveResult.Error.Error()).JSON(c)
	}

	// adjust value of object
	if member.HasPassword {
		member.Password = ""
		member.Salt = ""
	}

	if member.RequestFrom == model.Sturgeon {
		result, _ := saveResult.Result.(model.Member)
		member.Address = result.Address
	}

	return shared.NewHTTPResponse(http.StatusOK, "Member Edit Response", member).JSON(c)
}

// UpdatePassword function for updating password
func (h *HTTPMemberHandler) UpdatePassword(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	oldPasswordV2 := c.FormValue(model.FieldOldPassword)
	newPasswordV2 := c.FormValue(model.FieldNewPassword)

	authorization := c.Request().Header.Get(echo.HeaderAuthorization)
	var token string
	if split := strings.Split(authorization, " "); len(split) > 1 {
		token = split[1]
	}

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	passResultV2 := <-h.MemberUseCase.UpdatePassword(newCtx, token, memberID, oldPasswordV2, newPasswordV2)

	if passResultV2.Error != nil {
		return shared.NewHTTPResponse(passResultV2.HTTPStatus, passResultV2.Error.Error()).JSON(c)
	}

	res, ok := passResultV2.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(resultIsNotProperResponse)
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Update Password Response", res).JSON(c)
}

// AddPassword function for adding new password for logged in member
func (h *HTTPMemberHandler) AddPassword(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// get member data
	memberResult := h.GetMember(c.Request().Context(), memberID)
	if memberResult.Error != nil {
		return shared.NewHTTPResponse(memberResult.HTTPStatus, memberResult.Error.Error()).JSON(c)
	}

	member := memberResult.Result

	// check whether password exists or not
	if len(member.Password) > 0 && len(member.Salt) > 0 {
		err := errors.New("you are not allowed to add new password")
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	// append the member object
	member.NewPassword = c.FormValue(model.FieldPassword)
	member.RePassword = c.FormValue(model.FieldRePassword)

	if member.RePassword != member.NewPassword {
		err := errors.New("rePassword doesn't match password")
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	newCtx := c.Request().Context()
	newCtx = shared.SetDataToContext(newCtx, shared.ContextKey(helper.TextToken), c.Request().Header.Get(helper.TextAuthorization))

	passResult := <-h.MemberUseCase.AddNewPassword(newCtx, member)

	if passResult.Error != nil {
		err := errors.New("failed to add password")
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	res := model.SuccessResponse{
		ID:          golib.RandomString(8),
		Message:     helper.SuccessMessage,
		HasPassword: true,
		Email:       member.Email,
		FirstName:   member.FirstName,
		LastName:    member.LastName,
	}

	return shared.NewHTTPResponse(http.StatusOK, "Add Password Response", res).JSON(c)
}

// RegisterMember function for registering new member
func (h *HTTPMemberHandler) RegisterMember(c echo.Context) error {
	memberV2 := &model.Member{}
	memberV2.FirstName = c.FormValue(model.FieldFirstName)
	memberV2.LastName = c.FormValue(model.FieldLastName)
	memberV2.Email = c.FormValue(model.FieldEmail)
	memberV2.NewPassword = c.FormValue(model.FieldPassword)
	memberV2.RePassword = c.FormValue(model.FieldRePassword)
	memberV2.GenderString = c.FormValue(model.FieldGender)
	memberV2.BirthDateString = c.FormValue(model.FieldDOB)
	memberV2.Mobile = c.FormValue(model.FieldMobile)
	memberV2.Type = "register"

	registerType := c.FormValue("registerType")
	signUpFrom := c.FormValue("signUpFrom")
	memberV2.SignUpFrom = signUpFrom
	memberV2.RegisterType = registerType
	// check member existence first when error not happens
	// or the error is not equal data not found
	checkResultV2 := <-h.MemberUseCase.CheckEmailAndMobileExistence(c.Request().Context(), memberV2)
	if checkResultV2.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, checkResultV2.Error.Error()).JSON(c)
	}

	saveResult := <-h.MemberUseCase.RegisterMember(c.Request().Context(), memberV2)
	if saveResult.Error != nil {
		return shared.NewHTTPResponse(saveResult.HTTPStatus, saveResult.Error.Error()).JSON(c)
	}

	result, ok := saveResult.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(resultIsNotProperResponse)
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	if signUpFrom == model.Sturgeon {
		res := model.PlainSuccessResponse{Email: memberV2.Email}
		return shared.NewHTTPResponse(http.StatusCreated, "Member Register Response", res).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusCreated, "Member Register Response", result).JSON(c)
}

// ActivationMember function for activating inactivate member
func (h *HTTPMemberHandler) ActivationMember(c echo.Context) error {
	token := c.FormValue("token")
	requestFrom := c.FormValue(requestFrom)

	activateResult := <-h.MemberUseCase.ActivateMember(c.Request().Context(), token, requestFrom)
	if activateResult.Error != nil {
		return shared.NewHTTPResponse(activateResult.HTTPStatus, activateResult.Error.Error()).JSON(c)
	}

	res, ok := activateResult.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(resultIsNotProperResponse)
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	if requestFrom == model.Sturgeon {
		res := model.PlainSuccessResponse{Email: res.Email}
		return shared.NewHTTPResponse(http.StatusOK, "Member Activation Response", res).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Member Activation Response", res).JSON(c)
}

// ForgotPassword function for requesting token of forgot password
func (h *HTTPMemberHandler) ForgotPassword(c echo.Context) error {
	email := c.FormValue(model.FieldEmail)

	passResult := <-h.MemberUseCase.ForgotPassword(c.Request().Context(), email)
	if passResult.Error != nil {
		return shared.NewHTTPResponse(passResult.HTTPStatus, passResult.Error.Error()).JSON(c)
	}

	res, ok := passResult.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(resultIsNotProperResponse)
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}
	env, ok := os.LookupEnv("ENV")
	if !ok {
		env = "PROD"
	}
	if env != "DEV" {
		return shared.NewHTTPResponse(http.StatusOK, "Forgot Password Response").JSON(c)
	}
	return shared.NewHTTPResponse(http.StatusOK, "Forgot Password Response", res).JSON(c)

}

// ChangeForgotPassword function for creating new password after requesting token
func (h *HTTPMemberHandler) ChangeForgotPassword(c echo.Context) error {
	token := c.FormValue("token")
	newPassword := c.FormValue(model.FieldNewPassword)
	rePassword := c.FormValue(model.FieldRePassword)
	requestFrom := c.FormValue(requestFrom)

	passResult := <-h.MemberUseCase.ChangeForgotPassword(c.Request().Context(), token, newPassword, rePassword, requestFrom)
	if passResult.Error != nil {
		return shared.NewHTTPResponse(passResult.HTTPStatus, passResult.Error.Error()).JSON(c)
	}

	res, ok := passResult.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(resultIsNotProperResponse)
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Change Password Response", res).JSON(c)
}

// ActivateNewPassword function for activating new password for member who is registered from dolphin
func (h *HTTPMemberHandler) ActivateNewPassword(c echo.Context) error {
	token := c.FormValue("token")
	newPassword := c.FormValue(model.FieldNewPassword)
	rePassword := c.FormValue(model.FieldRePassword)

	activateResult := <-h.MemberUseCase.ActivateNewPassword(c.Request().Context(), token, newPassword, rePassword)
	if activateResult.Error != nil {
		return shared.NewHTTPResponse(activateResult.HTTPStatus, activateResult.Error.Error()).JSON(c)
	}

	res, ok := activateResult.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(resultIsNotProperResponse)
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Activate Password Response", res).JSON(c)
}

// RegenerateToken function for regenerating token for activation
func (h *HTTPMemberHandler) RegenerateToken(c echo.Context) error {
	// get member data
	memberResult := h.GetMember(c.Request().Context(), c.Param(memberID))
	if memberResult.Error != nil {
		return shared.NewHTTPResponse(memberResult.HTTPStatus, memberResult.Error.Error()).JSON(c)
	}

	member := memberResult.Result

	saveResult := <-h.MemberUseCase.RegenerateToken(c.Request().Context(), member)
	if saveResult.Error != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, saveResult.Error.Error()).JSON(c)
	}

	result, ok := saveResult.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(resultIsNotProperResponse)
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Regenerate Token Response", result).JSON(c)
}

// AddNewMember function for adding new member from other pack
func (h *HTTPMemberHandler) AddNewMember(c echo.Context) error {
	// adjust data from form value
	memberAddress := model.Address{}
	member := &model.Member{}

	member.Email = c.FormValue(model.FieldEmail)
	member.FirstName = c.FormValue(model.FieldFirstName)
	member.LastName = c.FormValue(model.FieldLastName)
	member.Mobile = c.FormValue(model.FieldMobile)
	member.Phone = c.FormValue(model.FieldPhone)
	member.Ext = c.FormValue(model.FieldExt)
	member.GenderString = c.FormValue(model.FieldGender)
	member.BirthDateString = c.FormValue(model.FieldDOB)
	member.Password = c.FormValue(model.FieldPassword)
	member.RePassword = c.FormValue(model.FieldRePassword)
	member.SignUpFrom = c.FormValue("signUpFrom")
	member.Type = "add"

	memberAddress.Street1 = c.FormValue(model.FieldStreet1)
	memberAddress.Street2 = c.FormValue(model.FieldStreet2)
	memberAddress.ZipCode = c.FormValue(model.FieldPostalCode)
	memberAddress.SubDistrictID = c.FormValue(model.FieldSubDistrictID)
	memberAddress.SubDistrict = c.FormValue(model.FieldSubDistrictName)
	memberAddress.DistrictID = c.FormValue(model.FieldDistrictID)
	memberAddress.District = c.FormValue(model.FieldDistrictName)
	memberAddress.CityID = c.FormValue(model.FieldCityID)
	memberAddress.City = c.FormValue(model.FieldCityName)
	memberAddress.ProvinceID = c.FormValue(model.FieldProvinceID)
	memberAddress.Province = c.FormValue(model.FieldProvinceName)
	member.Address = memberAddress

	// check member existence first when error not happens
	// or the error is not equal data not found
	checkResult := <-h.MemberUseCase.CheckEmailAndMobileExistence(c.Request().Context(), member)
	if checkResult.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, checkResult.Error.Error()).JSON(c)
	}

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	saveResult := <-h.MemberUseCase.RegisterMember(newCtx, member)
	if saveResult.Error != nil {
		return shared.NewHTTPResponse(saveResult.HTTPStatus, saveResult.Error.Error()).JSON(c)
	}

	result, ok := saveResult.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(resultIsNotProperResponse)
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusCreated, "Add Member Response", result).JSON(c)
}

// GetDetailMember function for getting detail member based on id
func (h *HTTPMemberHandler) GetDetailMember(c echo.Context) error {
	// get member data
	ctxReq := c.Request().Context()
	ctxReq = context.WithValue(ctxReq, shared.ContextKey(helper.TextToken), c.Request().Header.Get(echo.HeaderAuthorization))
	memberResult := h.GetMember(ctxReq, c.Param(memberID))
	if memberResult.Error != nil {
		return shared.NewHTTPResponse(memberResult.HTTPStatus, memberResult.Error.Error()).JSON(c)
	}

	member := memberResult.Result

	return shared.NewHTTPResponse(http.StatusOK, "Get Member Response", member).JSON(c)
}

// GetMembers function for getting list of members
func (h *HTTPMemberHandler) GetMembers(c echo.Context) error {
	params := model.Parameters{
		Query:    c.QueryParam("query"),
		StrPage:  c.QueryParam("page"),
		StrLimit: c.QueryParam("limit"),
		Sort:     c.QueryParam("sort"),
		OrderBy:  c.QueryParam("orderBy"),
		Status:   c.QueryParam("status"),
		Email:    c.QueryParam("email"),
		IsStaff:  c.QueryParam("isStaff"),
		IsAdmin:  c.QueryParam("isAdmin"),
	}

	memberResult := <-h.MemberUseCase.GetListMembers(c.Request().Context(), &params)
	if memberResult.Error != nil {
		return shared.NewHTTPResponse(memberResult.HTTPStatus, memberResult.Error.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	member, ok := memberResult.Result.(model.ListMembers)
	if !ok {
		err := errors.New(helper.ErrorResultNotProper)
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	member.ID = golib.RandomString(8)
	member.Name = "list of members"

	totalPage := math.Ceil(float64(member.TotalData) / float64(params.Limit))

	if len(member.Members) <= 0 {
		response := shared.NewHTTPResponse(http.StatusOK, "Get Members Response", make(helper.EmptySlice, 0))
		response.SetSuccess(false)
		return response.JSON(c)
	}

	meta := shared.Meta{
		Page:         params.Page,
		Limit:        params.Limit,
		TotalRecords: member.TotalData,
		TotalPages:   int(totalPage),
	}
	return shared.NewHTTPResponse(http.StatusOK, "Get Members Response", member.Members, meta).JSON(c)
}

// MemberSend function for send to kafka for updated nav
func (h *HTTPMemberHandler) MemberSend(c echo.Context) error {
	memberID := c.FormValue(memberID)

	// get member data
	memberResult := <-h.MemberUseCase.GetDetailMemberByID(c.Request().Context(), memberID)
	if memberResult.Error != nil {
		return shared.NewHTTPResponse(memberResult.HTTPStatus, memberResult.Error.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	member, ok := memberResult.Result.(model.Member)
	if !ok {
		err := errors.New(helper.ErrorResultNotProper)
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	// publish to kafka
	h.MemberUseCase.PublishToKafkaUser(c.Request().Context(), &member, updateType)

	return shared.NewHTTPResponse(http.StatusOK, "Member Send Response", member).JSON(c)
}

// ResendActivation for resend activation email
func (h *HTTPMemberHandler) ResendActivation(c echo.Context) error {
	email := c.FormValue("email")

	res := <-h.MemberUseCase.ResendActivation(c.Request().Context(), strings.ToLower(email))
	if res.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, res.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, model.SuccessResendActivation, res.Result).JSON(c)
}

// RevokeAllAccess function for revoke all access login by user ID
func (h *HTTPMemberHandler) RevokeAllAccess(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	token, _ := c.Get("token").(*jwt.Token)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	memberResult := <-h.MemberUseCase.RevokeAllAccess(c.Request().Context(), memberID, token.Raw)
	if memberResult.Error != nil {
		return shared.NewHTTPResponse(memberResult.HTTPStatus, memberResult.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success revoke all access").JSON(c)
}

// GetLoginActivity function for getting status mfa settings
func (h *HTTPMemberHandler) GetLoginActivity(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	token, _ := c.Get("token").(*jwt.Token)
	strPage := c.QueryParam("page")
	strLimit := c.QueryParam("limit")

	params := model.ParametersLoginActivity{
		StrPage:  strPage,
		StrLimit: strLimit,
		MemberID: memberID,
		Token:    token.Raw,
	}
	if _, err := helper.ValidatePagination(helper.PaginationParameters{StrPage: strPage, StrLimit: strLimit}); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	sessionInfoResult := <-h.MemberUseCase.GetLoginActivity(c.Request().Context(), &params)
	if sessionInfoResult.Error != nil {
		return shared.NewHTTPResponse(sessionInfoResult.HTTPStatus, sessionInfoResult.Error.Error()).JSON(c)
	}

	sessionInfo, ok := sessionInfoResult.Result.(model.SessionInfoList)
	if !ok {
		err := errors.New("result is not list of session history")
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	totalPage := math.Ceil(float64(sessionInfo.TotalData) / float64(params.Limit))
	meta := shared.Meta{
		Page:         params.Page,
		Limit:        params.Limit,
		TotalRecords: sessionInfo.TotalData,
		TotalPages:   int(totalPage),
	}
	return shared.NewHTTPResponse(http.StatusOK, "Get login activity", sessionInfo.Data, meta).JSON(c)
}

// RevokeAccess function for revoke specific session
func (h *HTTPMemberHandler) RevokeAccess(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	sessionID := c.QueryParam("sid")
	if _, err := strconv.ParseInt(sessionID, 10, 32); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, errInvalidInput).JSON(c)
	}

	if err != nil || sessionID == "" {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	revokeAction := <-h.MemberUseCase.RevokeAccess(c.Request().Context(), memberID, sessionID)
	if revokeAction.Error != nil {
		return shared.NewHTTPResponse(revokeAction.HTTPStatus, revokeAction.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success revoke access").JSON(c)
}

// SyncPassword function update flag on member and b2b_contact
func (h *HTTPMemberHandler) SyncPassword(c echo.Context) error {
	var payload struct {
		OldPassword string `json:"oldPassword" form:"oldPassword"`
		NewPassword string `json:"newPassword" form:"newPassword"`
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
	syncPassword := <-h.MemberUseCase.SyncPassword(newCtx, token, payload.OldPassword, payload.NewPassword)
	if syncPassword.Error != nil {
		return shared.NewHTTPResponse(syncPassword.HTTPStatus, syncPassword.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success sync password").JSON(c)
}

// ChangePassword ...
func (h *HTTPMemberHandler) ChangePassword(c echo.Context) error {
	var payload struct {
		OldPassword string `json:"oldPassword" form:"oldPassword"`
		NewPassword string `json:"newPassword" form:"newPassword"`
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
	result := <-h.MemberUseCase.ChangePassword(newCtx, token, payload.OldPassword, payload.NewPassword)
	if result.Error != nil {
		return shared.NewHTTPResponse(result.HTTPStatus, result.Error.Error()).JSON(c)
	}

	res, ok := result.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(resultIsNotProperResponse)
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Update Password Response", res).JSON(c)
}

// Clients ...
func (h *HTTPMemberHandler) Clients(c echo.Context) error {
	result := <-h.MemberUseCase.Clients(c.Request().Context(), c.Request().Header.Get(helper.TextAuthorization))
	if result.Error != nil {
		return shared.NewHTTPResponse(result.HTTPStatus, result.Error.Error()).JSON(c)
	}

	data, ok := result.Result.(model.User)
	if !ok {
		err := errors.New("result is not list")
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success list client", data).JSON(c)
}

// ActivationMerchantEmployee function for activating merchant employee
func (h *HTTPMemberHandler) ActivationMerchantEmployee(c echo.Context) error {
	token := c.QueryParam("token")
	if token == "" {
		err := errors.New("token not found")
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	activateResult := <-h.MemberUseCase.ActivateMerchantEmployee(c.Request().Context(), token)
	if activateResult.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadGateway, activateResult.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success activation merchant employee").JSON(c)
}
