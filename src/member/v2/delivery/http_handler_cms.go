package delivery

import (
	"context"
	"errors"
	"net/http"

	"github.com/Bhinneka/golib/jsonschema"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/labstack/echo"
)

// UpdateMember function for updating member from other pack
func (h *HTTPMemberHandler) UpdateMember(c echo.Context) error {
	paramMemberID := c.Param(memberID)
	// adjust data from form value
	ma := model.Address{}
	member := model.Member{}
	member.ID = paramMemberID
	member.FirstName = c.FormValue(model.FieldFirstName)
	member.LastName = c.FormValue(model.FieldLastName)
	member.Mobile = c.FormValue(model.FieldMobile)
	member.Phone = c.FormValue(model.FieldPhone)
	member.Ext = c.FormValue(model.FieldExt)
	member.GenderString = c.FormValue(model.FieldGender)
	member.BirthDateString = c.FormValue(model.FieldDOB)
	member.StatusString = c.FormValue(model.FieldStatus)

	// optional
	member.IsStaffString = c.FormValue("isStaff")
	member.IsAdminString = c.FormValue("isAdmin")
	member.IsActiveString = c.FormValue("isActive")

	ma.Street1 = c.FormValue(model.FieldStreet1)
	ma.Street2 = c.FormValue(model.FieldStreet2)
	ma.ZipCode = c.FormValue(model.FieldPostalCode)
	ma.SubDistrictID = c.FormValue(model.FieldSubDistrictID)
	ma.SubDistrict = c.FormValue(model.FieldSubDistrictName)
	ma.DistrictID = c.FormValue(model.FieldDistrictID)
	ma.District = c.FormValue(model.FieldDistrictName)
	ma.CityID = c.FormValue(model.FieldCityID)
	ma.City = c.FormValue(model.FieldCityName)
	ma.ProvinceID = c.FormValue(model.FieldProvinceID)
	ma.Province = c.FormValue(model.FieldProvinceName)

	member.Address = ma
	member.Type = updateType
	member.UpdateFrom = model.Dolphin

	mErr := jsonschema.ValidateTemp("update_member_params_v2", member)
	if mErr != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, mErr.Error()).JSON(c)
	}

	newCtx := context.WithValue(c.Request().Context(), helper.TextAuthorization, c.Request().Header.Get(helper.TextAuthorization))
	newCtx = context.WithValue(newCtx, middleware.ContextKeyClientIP, c.RealIP())
	saveResult := <-h.MemberUseCase.UpdateDetailMemberByID(newCtx, member)
	if saveResult.Error != nil {
		return shared.NewHTTPResponse(saveResult.HTTPStatus, saveResult.Error.Error()).JSON(c)
	}

	member, ok := saveResult.Result.(model.Member)
	if !ok {
		err := errors.New("result is not member")
		return shared.NewHTTPResponse(http.StatusInternalServerError, err.Error()).JSON(c)
	}

	// adjust value of object
	if !member.HasPassword {
		member.Password = ""
		member.Salt = ""
	}

	return shared.NewHTTPResponse(http.StatusOK, "Update Member Response", member).JSON(c)
}

// DisabledMFAMember function for disabled mfa member
func (h *HTTPMemberHandler) DisabledMFAMember(c echo.Context) error {
	memberID := c.Param(memberID)

	memberResult := <-h.MemberUseCase.DisabledMFASetting(c.Request().Context(), memberID, helper.TextAccount)
	if memberResult.Error != nil {
		return shared.NewHTTPResponse(memberResult.HTTPStatus, memberResult.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success disable MFA").JSON(c)
}

// DisabledMFAMemberNarwhal function to disable MFA Narwhal from CMS
func (h *HTTPMemberHandler) DisabledMFAMemberNarwhal(c echo.Context) error {
	memberID := c.Param(memberID)

	actionResult := <-h.MemberUseCase.DisabledMFASetting(c.Request().Context(), memberID, helper.TextNarwhal)
	if actionResult.Error != nil {
		return shared.NewHTTPResponse(actionResult.HTTPStatus, actionResult.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success disable MFA Narwhal").JSON(c)
}
