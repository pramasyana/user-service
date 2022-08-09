package delivery

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/labstack/echo"
	"github.com/tealeg/xlsx"
)

// MigrateData function for migrating data from squid
func (h *HTTPMemberHandler) MigrateData(c echo.Context) error {
	// bind and get the body
	members := &model.Members{}
	if err := c.Bind(members); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	migrateResult := <-h.MemberUseCase.MigrateMember(c.Request().Context(), members)
	if migrateResult.Error != nil {
		// if error data exists then return its array
		if len(migrateResult.ErrorData) > 0 {
			data := struct {
				Message string              `json:"message"`
				Data    []model.MemberError `json:"data"`
			}{migrateResult.Error.Error(), migrateResult.ErrorData}
			return shared.NewHTTPResponse(http.StatusInternalServerError, "Error Migrate Members", data).JSON(c)
		}

		return shared.NewHTTPResponse(http.StatusInternalServerError, migrateResult.Error.Error()).JSON(c)
	}

	res := model.SuccessResponse{
		ID:      golib.RandomString(8),
		Message: helper.SuccessMessage,
	}

	return shared.NewHTTPResponse(http.StatusOK, "Migrate Member Response", res).JSON(c)
}

// BulkMemberSend function for bulk send member using excel file
func (h *HTTPMemberHandler) BulkMemberSend(c echo.Context) error {
	base64File := c.FormValue("file")

	fileDecode, err := base64.StdEncoding.DecodeString(base64File)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}
	xlsBinary, err := xlsx.OpenBinary(fileDecode)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)

	}

	successMember := make([]string, 0)

	for _, row := range xlsBinary.Sheets[0].Rows {
		if len(row.Cells) == 0 {
			continue
		}

		if row.Cells[0].String() == memberID || row.Cells[0].String() == "" {
			continue
		}

		memberID := row.Cells[0].String()

		// get member data
		memberResult := <-h.MemberUseCase.GetDetailMemberByID(c.Request().Context(), memberID)
		if memberResult.Error != nil {
			//Error, skip row
			continue
		}
		member, ok := memberResult.Result.(model.Member)
		if !ok {
			//Error, skip row
			continue
		}

		successMember = append(successMember, memberID)

		// publish to kafka
		h.MemberUseCase.PublishToKafkaUser(c.Request().Context(), &member, updateType)
	}

	return shared.NewHTTPResponse(http.StatusCreated, "Bulk Member Send Response", successMember).JSON(c)
}

// ImportMember function for import new member using excel file
func (h *HTTPMemberHandler) ImportMember(c echo.Context) error {
	ctx := "MemberPresenter-ImportMember"
	base64File := c.FormValue("file")

	fileDecode, err := base64.StdEncoding.DecodeString(base64File)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	var memberExistLines []string
	members, err := h.MemberUseCase.ParseMemberData(c.Request().Context(), fileDecode)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	invalidRows := h.MemberUseCase.BulkValidateEmailAndPhone(c.Request().Context(), members)
	if len(invalidRows) > 0 {
		memberExistLines = append(memberExistLines, invalidRows...)
		response := shared.NewHTTPResponse(http.StatusBadRequest, strings.Join(memberExistLines, ","))
		return response.JSON(c)
	}

	invalidResult := h.importMemberData(c.Request().Context(), ctx, members)
	if len(invalidResult) > 0 {
		memberExistLines = append(memberExistLines, invalidResult...)
		response := shared.NewHTTPResponse(http.StatusBadRequest, strings.Join(memberExistLines, ","))
		return response.JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusCreated, "Import Member Response").JSON(c)
}

func (h *HTTPMemberHandler) importMemberData(ctxReq context.Context, ctx string, members []*model.Member) (memberImportLines []string) {
	tc := tracer.StartTrace(ctxReq, "HTTPMemberHandler-importMemberData")
	saveResult := <-h.MemberUseCase.BulkImportMember(tc.NewChildContext(), members)
	if saveResult.Error != nil {
		lineExist := fmt.Sprintf("%s", saveResult.Error)
		memberImportLines = append(memberImportLines, lineExist)
		return
	}
	defer tc.Finish(nil)

	_, ok := saveResult.Result.([]model.SuccessResponse)
	if !ok {
		err := errors.New(resultIsNotProperResponse)
		memberImportLines = append(memberImportLines, err.Error())
		return
	}

	return memberImportLines
}
