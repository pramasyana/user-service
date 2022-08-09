package delivery

import (
	"net/http"

	"github.com/Bhinneka/user-service/middleware"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/labstack/echo"
)

// ChangeProfilePicture function for updating member profile picture
func (h *HTTPMemberHandler) ChangeProfilePicture(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	mp := model.ProfilePicture{}
	mp.ID = memberID
	mp.ProfilePicture = c.FormValue("file")
	saveResult := <-h.MemberUseCase.UpdateProfilePicture(c.Request().Context(), mp)
	if saveResult.Error != nil {
		return shared.NewHTTPResponse(saveResult.HTTPStatus, saveResult.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Profile picture updated", mp).JSON(c)
}

// GetProfileComplete function for getting completeness profile information
func (h *HTTPMemberHandler) GetProfileComplete(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	memberResult := <-h.MemberUseCase.GetProfileComplete(c.Request().Context(), memberID)
	if memberResult.Error != nil {
		return shared.NewHTTPResponse(memberResult.HTTPStatus, memberResult.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Profile Completeness", memberResult.Result).JSON(c)
}

// ChangeProfileName function for updating member name
func (h *HTTPMemberHandler) ChangeProfileName(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	mp := model.ProfileName{}
	mp.ID = memberID
	mp.ProfileName = c.FormValue("name")
	saveResult := <-h.MemberUseCase.UpdateProfileName(c.Request().Context(), mp)
	if saveResult.Error != nil {
		return shared.NewHTTPResponse(saveResult.HTTPStatus, saveResult.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Profile name updated", mp).JSON(c)
}
