package delivery

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/google/jsonapi"
	"github.com/labstack/echo"
	"golang.org/x/net/context"
)

// ImportMember function for import new member using excel file
func (h *HTTPMemberHandler) ImportMember(c echo.Context) error {
	ctx := "MemberPresenter-ImportMember"
	base64File := c.FormValue("file")

	fileDecode, err := base64.StdEncoding.DecodeString(base64File)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	members, err := h.MemberUseCase.ParseMemberData(c.Request().Context(), fileDecode)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	var memberExistLines []string

	invalidRows := h.validateRows(c.Request().Context(), members)
	if len(invalidRows) > 0 {
		memberExistLines = append(memberExistLines, invalidRows...)
		c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(strings.Join(memberExistLines, ","))
	}
	invalidResult := h.importMemberData(c.Request().Context(), ctx, members)
	if len(invalidResult) > 0 {
		memberExistLines = append(memberExistLines, invalidResult...)
		c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
		c.Response().WriteHeader(http.StatusBadRequest)
		return json.NewEncoder(c.Response()).Encode(memberExistLines)
	}

	c.Response().Header().Set(echo.HeaderContentType, jsonapi.MediaType)
	c.Response().WriteHeader(http.StatusOK)
	return json.NewEncoder(c.Response()).Encode(memberExistLines)
}

func (h *HTTPMemberHandler) validateRows(ctxReq context.Context, members []*model.Member) []string {
	var memberExistLines []string
	for _, member := range members {
		checkResult := <-h.MemberUseCase.CheckEmailAndMobileExistence(ctxReq, member)
		if checkResult.Error != nil {
			lineExist := fmt.Sprintf("%s already exist", member.Email)
			memberExistLines = append(memberExistLines, lineExist)
		}
	}
	return memberExistLines
}

func (h *HTTPMemberHandler) importMemberData(ctxReq context.Context, ctx string, members []*model.Member) []string {
	var memberImportLines []string
	for _, member := range members {
		saveResult := <-h.MemberUseCase.ImportMember(ctxReq, member)
		if saveResult.Error != nil {
			lineExist := fmt.Sprintf("%s %s", member.Email, saveResult.Error)
			memberImportLines = append(memberImportLines, lineExist)
			continue
		}

		_, ok := saveResult.Result.(model.SuccessResponse)
		if !ok {
			err := errors.New(helper.ErrorResultNotProper)
			memberImportLines = append(memberImportLines, err.Error())
			continue
		}
	}
	return memberImportLines
}
