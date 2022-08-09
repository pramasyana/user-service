package delivery

import (
	"errors"
	"math"
	"net/http"
	"strings"

	"github.com/Bhinneka/golib/jsonschema"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/middleware"
	"github.com/Bhinneka/user-service/src/document/v2/model"
	"github.com/Bhinneka/user-service/src/document/v2/usecase"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/labstack/echo"
)

const (
	tempDocument            = "add_document_params_v2"
	scopeExtractMemberID    = "extract_member_id"
	scopeDeleteDocument     = "delete_document"
	scopeUpdateDocument     = "update_document"
	scopeSaveDocument       = "save_document"
	scopeSaveDocumentType   = "save_document_type"
	scopeUpdateDocumentType = "update_document_type"
	scopeParseDocument      = "parse_document"
	messageErrorResult      = "result is not proper response"
	messageTypeSuccess      = "Document Type Success"
	messageSuccess          = "Document Success"
	messageDeleteSuccess    = "Document Deleted"
	urlParamDocumentID      = "/me/:documentID"
	fieldDocumentType       = "documentType"
	paramDocumentID         = "documentID"
)

// HTTPDocumentHandler model
type HTTPDocumentHandler struct {
	DocumentUseCase usecase.DocumentUseCase
}

// NewHTTPHandler function for initialise *HTTPAuthHandler
func NewHTTPHandler(DocumentUseCase usecase.DocumentUseCase) *HTTPDocumentHandler {
	return &HTTPDocumentHandler{DocumentUseCase: DocumentUseCase}
}

// MountMe function for mounting routes
func (s *HTTPDocumentHandler) MountMe(group *echo.Group) {
	group.POST("/me", s.AddDocumentMe)
	group.PUT(urlParamDocumentID, s.UpdateDocumentMe)
	group.DELETE(urlParamDocumentID, s.DeleteDocumentMe)
	group.GET("/me", s.GetDocumentMe)
	group.GET(urlParamDocumentID, s.GetDocumentDetailMe)
	group.GET("/list", s.GetRequiredDocument)
}

// MountDocumentType function for mounting endpoints
func (s *HTTPDocumentHandler) MountDocumentType(group *echo.Group) {
	group.POST("", s.AddDocumentType)
	group.PUT("/:documentTypeID", s.UpdateDocumentType)
	group.GET("", s.GetDocumentType)
}

// AddDocumentMe function for add new document
func (s *HTTPDocumentHandler) AddDocumentMe(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	payload := model.DocumentData{}
	if err := c.Bind(&payload); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, helper.ErrorPayload).JSON(c)
	}

	payload.MemberID = memberID
	payload.DocumentType = strings.ToUpper(payload.DocumentType)
	payload.CreatedBy = memberID
	// validate schema fields
	errJsonSchema := jsonschema.ValidateTemp(tempDocument, payload)
	if errJsonSchema != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, errJsonSchema.Error()).JSON(c)
	}

	// Add document usecase process
	updateResult := <-s.DocumentUseCase.AddUpdateDocument(c.Request().Context(), payload)
	if updateResult.Error != nil {
		return shared.NewHTTPResponse(updateResult.HTTPStatus, updateResult.Error.Error()).JSON(c)
	}

	result, ok := updateResult.Result.(model.DocumentData)
	if !ok {
		err := errors.New(messageErrorResult)
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusCreated, messageSuccess, result).JSON(c)
}

// UpdateDocumentMe function for add new document
func (s *HTTPDocumentHandler) UpdateDocumentMe(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	payload := model.DocumentData{}
	if err := c.Bind(&payload); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, helper.ErrorPayload).JSON(c)
	}

	payload.ID = c.Param(paramDocumentID)
	payload.MemberID = memberID
	payload.DocumentType = strings.ToUpper(payload.DocumentType)
	payload.ModifiedBy = memberID

	// validate schema fields
	errValidateUpdate := jsonschema.ValidateTemp(tempDocument, payload)
	if errValidateUpdate != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, errValidateUpdate.Error()).JSON(c)
	}

	// Add  document usecase process
	updateResult := <-s.DocumentUseCase.AddUpdateDocument(c.Request().Context(), payload)
	if updateResult.Error != nil {
		return shared.NewHTTPResponse(updateResult.HTTPStatus, updateResult.Error.Error()).JSON(c)
	}

	result, ok := updateResult.Result.(model.DocumentData)
	if !ok {
		err := errors.New(messageErrorResult)
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusCreated, messageSuccess, result).JSON(c)
}

// DeleteDocumentMe function for add new document
func (s *HTTPDocumentHandler) DeleteDocumentMe(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	// Add  document usecase process
	deleteResult := <-s.DocumentUseCase.DeleteDocument(c.Request().Context(), c.Param(paramDocumentID), memberID)
	if deleteResult.Error != nil {
		return shared.NewHTTPResponse(deleteResult.HTTPStatus, deleteResult.Error.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, messageDeleteSuccess).JSON(c)
}

// GetDocumentMe function for getting list of document
func (s *HTTPDocumentHandler) GetDocumentMe(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	params := model.DocumentParameters{
		StrPage:  c.QueryParam("page"),
		StrLimit: c.QueryParam("limit"),
		MemberID: memberID,
	}

	ctxReq := c.Request().Context()
	ctxReq = shared.SetDataToContext(ctxReq, shared.ContextKey(helper.TextToken), c.Request().Header.Get(echo.HeaderAuthorization))
	documentResult := <-s.DocumentUseCase.GetListDocument(ctxReq, &params)
	if documentResult.Error != nil {
		return shared.NewHTTPResponse(documentResult.HTTPStatus, documentResult.Error.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	document, ok := documentResult.Result.(model.ListDocument)
	if !ok {
		err := errors.New("result is not list of document")
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	totalPage := math.Ceil(float64(document.TotalData) / float64(params.Limit))

	if len(document.Document) <= 0 {
		response := shared.NewHTTPResponse(http.StatusOK, messageSuccess, make(helper.EmptySlice, 0))
		response.SetSuccess(false)
		return response.JSON(c)
	}

	meta := shared.Meta{
		Page:         params.Page,
		Limit:        params.Limit,
		TotalRecords: document.TotalData,
		TotalPages:   int(totalPage),
	}
	return shared.NewHTTPResponse(http.StatusOK, messageSuccess, document.Document, meta).JSON(c)
}

// GetDocumentDetailMe function for getting detail of document
func (s *HTTPDocumentHandler) GetDocumentDetailMe(c echo.Context) error {
	memberID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	ctxReq := c.Request().Context()
	ctxReq = shared.SetDataToContext(ctxReq, shared.ContextKey(helper.TextToken), c.Request().Header.Get(echo.HeaderAuthorization))
	documentResult := <-s.DocumentUseCase.GetDetailDocument(ctxReq, c.Param(paramDocumentID), memberID)
	if documentResult.Error != nil {
		return shared.NewHTTPResponse(documentResult.HTTPStatus, documentResult.Error.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	document, ok := documentResult.Result.(model.DocumentData)
	if !ok {
		err := errors.New("result is not detail of document")
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, messageSuccess, document).JSON(c)
}

// AddDocumentType function for add new document
func (s *HTTPDocumentHandler) AddDocumentType(c echo.Context) error {
	err := middleware.ExtractClaimsIsAdmin(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, err.Error()).JSON(c)
	}

	userID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	payloadCreate := model.DocumentTypePayload{}
	if err := c.Bind(&payloadCreate); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, helper.ErrorPayload).JSON(c)
	}

	data := model.DocumentType{
		ID:             c.Param("documentTypeID"),
		DocumentType:   strings.ToUpper(payloadCreate.DocumentType),
		IsB2bString:    payloadCreate.IsB2b,
		IsB2cString:    payloadCreate.IsB2c,
		IsActiveString: payloadCreate.IsActive,
		CreatedBy:      userID,
	}

	// Add  document type usecase process
	upsertResult := <-s.DocumentUseCase.AddUpdateDocumentType(c.Request().Context(), data)
	if upsertResult.Error != nil {
		return shared.NewHTTPResponse(upsertResult.HTTPStatus, upsertResult.Error.Error()).JSON(c)
	}

	result, ok := upsertResult.Result.(model.DocumentType)
	if !ok {
		err := errors.New(messageErrorResult)
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusCreated, messageTypeSuccess, result).JSON(c)
}

// UpdateDocumentType function for add new document
func (s *HTTPDocumentHandler) UpdateDocumentType(c echo.Context) error {
	err := middleware.ExtractClaimsIsAdmin(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusUnauthorized, err.Error()).JSON(c)
	}

	userID, err := middleware.ExtractMemberIDFromToken(c)
	if err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	payloadUpdate := model.DocumentTypePayload{}
	if err := c.Bind(&payloadUpdate); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, helper.ErrorPayload).JSON(c)
	}

	data := model.DocumentType{
		ID:             c.Param("documentTypeID"),
		DocumentType:   strings.ToUpper(payloadUpdate.DocumentType),
		IsB2bString:    payloadUpdate.IsB2b,
		IsB2cString:    payloadUpdate.IsB2c,
		IsActiveString: payloadUpdate.IsActive,
		ModifiedBy:     userID,
	}

	// Update  document type usecase process
	saveResult := <-s.DocumentUseCase.AddUpdateDocumentType(c.Request().Context(), data)
	if saveResult.Error != nil {
		return shared.NewHTTPResponse(saveResult.HTTPStatus, saveResult.Error.Error()).JSON(c)
	}

	resultDoc, ok := saveResult.Result.(model.DocumentType)
	if !ok {
		err := errors.New(messageErrorResult)
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, messageTypeSuccess, resultDoc).JSON(c)
}

// GetDocumentType function for getting list of document
func (s *HTTPDocumentHandler) GetDocumentType(c echo.Context) error {
	params := model.DocumentTypeParameters{
		StrPage:  c.QueryParam("page"),
		StrLimit: c.QueryParam("limit"),
		IsB2b:    c.QueryParam("isB2b"),
		IsB2c:    c.QueryParam("isB2c"),
	}

	documentResult := <-s.DocumentUseCase.GetListDocumentType(c.Request().Context(), &params)
	if documentResult.Error != nil {
		return shared.NewHTTPResponse(documentResult.HTTPStatus, documentResult.Error.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	document, ok := documentResult.Result.(model.ListDocumentType)
	if !ok {
		err := errors.New("result is not list of document type")
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	totalPage := math.Ceil(float64(document.TotalData) / float64(params.Limit))

	if len(document.DocumentType) <= 0 {
		response := shared.NewHTTPResponse(http.StatusOK, messageTypeSuccess, make(helper.EmptySlice, 0))
		response.SetSuccess(false)
		return response.JSON(c)
	}

	meta := shared.Meta{
		Page:         params.Page,
		Limit:        params.Limit,
		TotalRecords: document.TotalData,
		TotalPages:   int(totalPage),
	}
	return shared.NewHTTPResponse(http.StatusOK, messageTypeSuccess, document.DocumentType, meta).JSON(c)
}

// GetRequiredDocument function for getting list of document
func (s *HTTPDocumentHandler) GetRequiredDocument(c echo.Context) error {
	resultList := <-s.DocumentUseCase.GetRequiredDocument(c.Request().Context())

	if resultList.Error != nil {
		res := shared.NewHTTPResponse(resultList.HTTPStatus, resultList.Error.Error())
		return res.JSON(c)
	}

	response, _ := resultList.Result.(model.RequiredDocuments)

	res := shared.NewHTTPResponse(http.StatusOK, "Success Get Required Documents", response.Merchant)
	return res.JSON(c)
}
