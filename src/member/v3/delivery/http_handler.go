package delivery

import (
	"context"
	"errors"
	"net/http"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	authModel "github.com/Bhinneka/user-service/src/auth/v1/model"
	authUc "github.com/Bhinneka/user-service/src/auth/v1/usecase"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/member/v1/usecase"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/labstack/echo"
)

const (
	badResponse  = "result is not proper response"
	invalidInput = "invalid input"
)

// MemberHandlerV3 model receiver
type MemberHandlerV3 struct {
	MemberUseCase usecase.MemberUseCase
	AuthUseCase   authUc.AuthUseCase
}

// NewHTTPHandlerV3 function for initialise *MemberHandlerV3
func NewHTTPHandlerV3(memberUseCase usecase.MemberUseCase, authUsecase authUc.AuthUseCase) *MemberHandlerV3 {
	return &MemberHandlerV3{
		MemberUseCase: memberUseCase,
		AuthUseCase:   authUsecase,
	}
}

// Mount function for mounting routes
func (h *MemberHandlerV3) Mount(group *echo.Group) {
	group.POST("/forgot-password", h.ForgotPassword)
	group.POST("/register", h.RegisterMember)
	group.POST("/member", h.AddMember)
}

func (h *MemberHandlerV3) MountMeV3(group *echo.Group) {
	group.POST("/mfa/activation", h.ActivateMFASettingsV3)
}

// ForgotPassword function for requesting email for forgot password
func (h *MemberHandlerV3) ForgotPassword(c echo.Context) error {
	input := model.ForgotPasswordInput{}
	if err := c.Bind(&input); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "please check your input")
	}
	response := model.GenericResponse{
		Message: "we have sent email your password reset link",
	}

	passResultV3 := <-h.MemberUseCase.ForgotPassword(c.Request().Context(), input.Email)
	if passResultV3.Error != nil {
		return c.JSON(http.StatusOK, response)
	}

	return c.JSON(http.StatusOK, response)
}

// RegisterMember implement register, supprt json body
func (h *MemberHandlerV3) RegisterMember(c echo.Context) error {
	memberV3 := model.Member{}
	if err := c.Bind(&memberV3); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	memberV3.Type = "register"
	memberV3.NewPassword = memberV3.Password

	checkResult := <-h.MemberUseCase.CheckEmailAndMobileExistence(c.Request().Context(), &memberV3)
	if checkResult.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, checkResult.Error.Error()).JSON(c)
	}
	memberV3.APIVersion = helper.Version3

	saveResultV3 := <-h.MemberUseCase.RegisterMember(c.Request().Context(), &memberV3)
	if saveResultV3.Error != nil {
		return shared.NewHTTPResponse(saveResultV3.HTTPStatus, saveResultV3.Error.Error()).JSON(c)
	}

	result, ok := saveResultV3.Result.(model.SuccessResponse)
	if !ok {
		return shared.NewHTTPResponse(http.StatusInternalServerError, "bad response").JSON(c)
	}

	// if member doesn't have social media, mean regular register
	// send verification email
	if !memberV3.IsSocialMediaExist() {
		return shared.NewHTTPResponse(http.StatusCreated, "Member Register Response V3", result).JSON(c)
	}

	// return with token
	rt := authModel.RequestToken{
		GrantType:   authModel.AuthTypePassword,
		Email:       memberV3.Email,
		Password:    memberV3.NewPassword,
		DeviceLogin: authModel.DefaultDeviceLogin,
		DeviceID:    authModel.DefaultDeviceID,
	}
	tokenResultV3 := <-h.AuthUseCase.GenerateToken(c.Request().Context(), "", rt)
	if tokenResultV3.Error != nil {
		return shared.NewHTTPResponse(tokenResultV3.HTTPStatus, tokenResultV3.Error.Error()).JSON(c)
	}
	tokenV3, ok := tokenResultV3.Result.(authModel.RequestToken)
	if !ok {
		return shared.NewHTTPResponse(http.StatusInternalServerError, "result is not token").JSON(c)
	}
	res := model.PlainSuccessResponse{
		Email:        memberV3.Email,
		Token:        tokenV3.Token,
		RefreshToken: tokenV3.RefreshToken,
	}
	// member register using social media flow
	// send welcome email
	// moved to queue on RegisterMember usecase
	// <-h.MemberUseCase.SendEmailWelcomeMember(c.Request().Context(), result)

	return shared.NewHTTPResponse(http.StatusCreated, "Member Register Response V3", res).JSON(c)
}

// AddMember create member from internal cms
func (h *MemberHandlerV3) AddMember(c echo.Context) error {
	memberV3 := model.Member{}
	if err := c.Bind(&memberV3); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	checkResult := <-h.MemberUseCase.CheckEmailAndMobileExistence(c.Request().Context(), &memberV3)
	if checkResult.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, checkResult.Error.Error()).JSON(c)
	}
	memberV3.Type = "add"
	memberV3.APIVersion = helper.Version3

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	saveResult := <-h.MemberUseCase.RegisterMember(newCtx, &memberV3)
	if saveResult.Error != nil {
		return shared.NewHTTPResponse(saveResult.HTTPStatus, saveResult.Error.Error()).JSON(c)
	}

	resultV3, ok := saveResult.Result.(model.SuccessResponse)
	if !ok {
		err := errors.New(badResponse)
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Add Member Response", resultV3).JSON(c)
}

func (h *MemberHandlerV3) ActivateMFASettingsV3(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	activateData := model.MFAActivateSettings{}
	activateData.Otp = c.FormValue("otp")
	activateData.SharedKeyText = c.FormValue("sharedKeyText")
	activateData.Password = c.FormValue("password")
	activateData.MemberID = memberID
	activateData.RequestFrom = helper.TextAccount

	activateResult := <-h.MemberUseCase.ActivateMFASettingV3(c.Request().Context(), activateData)
	if activateResult.Error != nil {
		return shared.NewHTTPResponse(activateResult.HTTPStatus, activateResult.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, model.SuccessMFAActivation).JSON(c)
}
