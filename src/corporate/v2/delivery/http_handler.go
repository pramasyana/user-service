package delivery

import (
	"encoding/base64"
	"math"
	"net/http"
	"strings"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/corporate/v2/model"
	"github.com/Bhinneka/user-service/src/corporate/v2/usecase"
	"github.com/Bhinneka/user-service/src/shared"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/labstack/echo"
)

const (
	messageContactListSuccess = "Contact list Success"
	scopeParseContact         = "parse_contact"
	messageErrorResult        = "result is not proper response"
)

// HTTPCorporateHandler model
type HTTPCorporateHandler struct {
	CorporateUseCase usecase.CorporateUseCase
}

// NewHTTPHandler function for initialise *HTTPAuthHandler
func NewHTTPHandler(CorporateUseCase usecase.CorporateUseCase) *HTTPCorporateHandler {
	return &HTTPCorporateHandler{CorporateUseCase: CorporateUseCase}
}

// MountCorporate function for mounting endpoints
func (s *HTTPCorporateHandler) MountCorporate(group *echo.Group) {
	group.GET("/contact", s.GetContactList)
	group.GET("/contact/:contactID", s.GetContactDetail)
	group.POST("/contact/import", s.ImportContact)
}

// GetContactList function for getting list of contact
func (s *HTTPCorporateHandler) GetContactList(c echo.Context) error {
	params := model.ParametersContact{
		StrPage:  c.QueryParam("page"),
		StrLimit: c.QueryParam("limit"),
		Query:    c.QueryParam("query"),
		Status:   strings.ToUpper(c.QueryParam("status")),
	}

	contactResult := <-s.CorporateUseCase.GetAllListContact(c.Request().Context(), &params)
	if contactResult.Error != nil {
		return shared.NewHTTPResponse(contactResult.HTTPStatus, contactResult.Error.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	contact, ok := contactResult.Result.(sharedModel.ListContact)
	if !ok {
		return shared.NewHTTPResponse(http.StatusBadRequest, "result is not list of contact", make(helper.EmptySlice, 0)).JSON(c)
	}

	totalPage := math.Ceil(float64(contact.TotalData) / float64(params.Limit))

	if len(contact.Contact) <= 0 {
		response := shared.NewHTTPResponse(http.StatusOK, messageContactListSuccess, make(helper.EmptySlice, 0))
		response.SetSuccess(false)
		return response.JSON(c)
	}

	meta := shared.Meta{
		Page:         params.Page,
		Limit:        params.Limit,
		TotalRecords: contact.TotalData,
		TotalPages:   int(totalPage),
	}
	return shared.NewHTTPResponse(http.StatusOK, messageContactListSuccess, contact.Contact, meta).JSON(c)
}

// GetContactDetail function for getting detail of contact
func (s *HTTPCorporateHandler) GetContactDetail(c echo.Context) error {
	contactResult := <-s.CorporateUseCase.GetDetailContact(c.Request().Context(), c.Param("contactID"))
	if contactResult.Error != nil {
		return shared.NewHTTPResponse(contactResult.HTTPStatus, contactResult.Error.Error()).JSON(c)
	}

	result, ok := contactResult.Result.(sharedModel.B2BContactData)
	if !ok {
		return shared.NewHTTPResponse(http.StatusBadRequest, messageErrorResult).JSON(c)
	}
	return shared.NewHTTPResponse(http.StatusOK, messageContactListSuccess, result).JSON(c)
}

// ImportContact function for import contact using excel file
func (s *HTTPCorporateHandler) ImportContact(c echo.Context) error {
	source := model.ImportFile{}
	if err := c.Bind(&source); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	fileDecode, err := base64.StdEncoding.DecodeString(source.File)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	if _, err = s.CorporateUseCase.ImportContact(c.Request().Context(), fileDecode); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	return shared.NewHTTPResponse(http.StatusOK, "Success Import Contact").JSON(c)
}
