package usecase

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/document/v2/model"
	"github.com/Bhinneka/user-service/src/document/v2/repo"
	memberRepo "github.com/Bhinneka/user-service/src/member/v1/repo"
	sharedRepo "github.com/Bhinneka/user-service/src/shared/repository"
)

const (
	msgErrorSave        = "failed to save document"
	msgErrorUpdate      = "failed to update document"
	msgErrorSaveType    = "failed to save document type"
	msgErrorUpdateType  = "failed to update document type"
	msgErrorDoesntExist = "document doesn't exists"
)

// DocumentUseCaseImpl data structure
type DocumentUseCaseImpl struct {
	DocumentRepo     repo.DocumentRepository
	DocumentTypeRepo repo.DocumentTypeRepository
	MemberRepoRead   memberRepo.MemberRepository
	Repository       *sharedRepo.Repository
	ENVKey           string
}

// NewDocumentUseCase function for initialise document use case implementation
func NewDocumentUseCase(documentRepo repo.DocumentRepository,
	documentTypeRepo repo.DocumentTypeRepository,
	memberRepoRead memberRepo.MemberRepository,
	repository *sharedRepo.Repository,
	envKey string) DocumentUseCase {
	return &DocumentUseCaseImpl{
		DocumentRepo:     documentRepo,
		DocumentTypeRepo: documentTypeRepo,
		MemberRepoRead:   memberRepoRead,
		Repository:       repository,
		ENVKey:           envKey,
	}
}

// AddUpdateDocument function for add new address
func (s *DocumentUseCaseImpl) AddUpdateDocument(ctxReq context.Context, data model.DocumentData) <-chan ResultUseCase {
	ctx := "DocumentUseCase-AddUpdateDocument"
	output := make(chan ResultUseCase)
	var (
		err            error
		documentResult model.DocumentData
	)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		data, err = s.ValidateDocument(ctxReq, data)
		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// Get Member Data from MemberID
		getMemberByID := <-s.MemberRepoRead.Load(ctxReq, data.MemberID)
		if getMemberByID.Result == nil {
			err := fmt.Errorf("MemberID doesn't exist")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		if data.ID == "" {
			documentResult, err = s.ProcessAddDocument(ctxReq, data)
		} else {
			documentResult, err = s.ProcessUpdateDocument(ctxReq, data)
		}

		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: documentResult}
	})

	return output
}

// ValidateDocument function for validating document data
func (s *DocumentUseCaseImpl) ValidateDocument(ctxReq context.Context, data model.DocumentData) (model.DocumentData, error) {

	if len(data.Title) == 0 {
		err := errors.New("title is required")
		return data, err
	}

	if !golib.ValidateLatinOnly(data.Title) {
		err := errors.New("title only latin character")
		return data, err
	}

	if len(data.DocumentFile) == 0 {
		err := errors.New("document file is required")
		return data, err
	}

	if !helper.ValidateDocumentFileURL(data.DocumentFile) {
		err := errors.New("document file not valid")
		return data, err
	}

	if len(data.DocumentType) == 0 {
		err := errors.New("document type is required")
		return data, err
	}

	if len(data.Description) != 0 && !golib.ValidateLatinOnly(data.Description) {
		err := errors.New("description only latin character")
		return data, err
	}

	query := model.DocumentTypeParameters{
		DocumentType: data.DocumentType,
		IsActive:     "true",
	}

	documentTypeResult := <-s.DocumentTypeRepo.FindDocumentTypeByParam(ctxReq, &query)
	if documentTypeResult.Result == nil {
		err := errors.New("document type not valid")
		return data, err
	}

	if len(data.Number) <= 0 {
		err := errors.New("number required")
		return data, err
	}

	if !golib.ValidateNumeric(data.Number) {
		err := errors.New("number only alphanumeric")
		return data, err
	}

	data.DocumentFile = s.SplitURLFile(data.DocumentFile)

	return data, nil
}

// SplitURLFile function for split url
func (s *DocumentUseCaseImpl) SplitURLFile(url string) string {
	AWSMerchantDocumentURL, ok := os.LookupEnv("AWS_MERCHANT_DOCUMENT_URL")
	if ok {
		contains := strings.Replace(AWSMerchantDocumentURL, "https://", "", -1)
		if strings.Contains(url, contains) {
			url = strings.Replace(url, AWSMerchantDocumentURL, "", -1)
		}
	}

	AWSMerchantDocumentURLSalmon, ok := os.LookupEnv("AWS_MERCHANT_DOCUMENT_URL_SALMON")
	if ok {
		containsSalmon := strings.Replace(AWSMerchantDocumentURLSalmon, "https://", "", -1)
		if strings.Contains(url, containsSalmon) {
			url = strings.Replace(url, AWSMerchantDocumentURLSalmon, "", -1)
		}
	}
	return url
}

// ProcessAddDocument function process document
func (s *DocumentUseCaseImpl) ProcessAddDocument(ctxReq context.Context, data model.DocumentData) (model.DocumentData, error) {
	data.Created = time.Now()
	data.LastModified = time.Now()

	data.StatusText = model.WaitingString
	data.Status = model.StringToStatus(data.StatusText)

	// add document repository process to database
	saveResult := <-s.DocumentRepo.AddDocument(ctxReq, data)
	if saveResult.Error != nil {
		err := errors.New(msgErrorSave)
		return data, err
	}

	resultData, ok := saveResult.Result.(model.DocumentData)
	if !ok {
		err := errors.New(msgErrorSave)
		return data, err
	}
	return resultData, nil
}

// ProcessUpdateDocument function process document
func (s *DocumentUseCaseImpl) ProcessUpdateDocument(ctxReq context.Context, data model.DocumentData) (model.DocumentData, error) {
	query := model.DocumentParameters{
		ID: data.ID,
	}

	documentResult := <-s.DocumentRepo.FindDocumentByParam(ctxReq, &query)
	if documentResult.Result == nil {
		err := errors.New(msgErrorDoesntExist)
		return data, err
	}

	data.LastModified = time.Now()
	data.StatusText = model.WaitingString
	data.Status = model.StringToStatus(data.StatusText)

	// add document repository process to database
	saveResult := <-s.DocumentRepo.UpdateDocument(ctxReq, data)
	if saveResult.Error != nil {
		err := errors.New(msgErrorUpdate)
		return data, err
	}

	resultData, ok := saveResult.Result.(model.DocumentData)
	if !ok {
		err := errors.New(msgErrorUpdate)
		return data, err
	}
	return resultData, nil
}

// DeleteDocument function process delete document
func (s *DocumentUseCaseImpl) DeleteDocument(ctxReq context.Context, documentID string, memberID string) <-chan ResultUseCase {
	ctx := "DocumentUseCase-DeleteDocument"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		query := model.DocumentParameters{
			ID:       documentID,
			MemberID: memberID,
		}

		documentResult := <-s.DocumentRepo.FindDocumentByParam(ctxReq, &query)
		if documentResult.Result == nil {
			err := errors.New(msgErrorDoesntExist)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusNotFound}
			return
		}

		// delete document
		result := <-s.DocumentRepo.DeleteDocumentByID(ctxReq, documentID)
		if result.Error != nil {
			err := errors.New("Document delete failed")
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: nil}
	})
	return output
}

// GetListDocument function for getting list of document
func (s *DocumentUseCaseImpl) GetListDocument(ctxReq context.Context, params *model.DocumentParameters) <-chan ResultUseCase {
	ctx := "DocumentUseCase-GetListDocument"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		var err error

		paging, err := helper.ValidatePagination(
			helper.PaginationParameters{
				Page:     1, // default
				StrPage:  params.StrPage,
				Limit:    20, // default
				StrLimit: params.StrLimit,
			})

		if err != nil {
			tags[helper.TextResponse] = err.Error()
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		params.Page = paging.Page
		params.Limit = paging.Limit
		params.Offset = paging.Offset
		tags[helper.TextParameter] = params

		documentResult := <-s.DocumentRepo.GetListDocument(ctxReq, params)
		if documentResult.Error != nil {
			httpStatus := http.StatusInternalServerError

			// when data is not found
			if documentResult.Error == sql.ErrNoRows {
				httpStatus = http.StatusNotFound
				documentResult.Error = fmt.Errorf(helper.ErrorDataNotFound, "document")
			}

			output <- ResultUseCase{Error: documentResult.Error, HTTPStatus: httpStatus}
			return
		}

		document := documentResult.Result.(model.ListDocument)

		totalResult := <-s.DocumentRepo.GetTotalDocument(ctxReq, params)
		if totalResult.Error != nil {
			output <- ResultUseCase{Error: totalResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		total := totalResult.Result.(int)
		document.TotalData = total

		output <- ResultUseCase{Result: document}
	})

	return output
}

// GetDetailDocument function for getting detail of document
func (s *DocumentUseCaseImpl) GetDetailDocument(ctxReq context.Context, documentID string, memberID string) <-chan ResultUseCase {
	ctx := "DocumentUseCase-GetDetailDocument"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {

		params := model.DocumentParameters{
			ID:       documentID,
			MemberID: memberID,
		}
		documentResult := <-s.DocumentRepo.GetDetailDocument(ctxReq, &params)
		if documentResult.Error != nil {
			err := errors.New(msgErrorDoesntExist)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusNotFound}
			return
		}

		document, ok := documentResult.Result.(model.DocumentData)
		if !ok {
			err := errors.New(msgErrorDoesntExist)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: document}
	})

	return output
}

// AddUpdateDocumentType function for add new address
func (s *DocumentUseCaseImpl) AddUpdateDocumentType(ctxReq context.Context, data model.DocumentType) <-chan ResultUseCase {
	ctx := "DocumentUseCase-AddUpdateDocumentType"
	output := make(chan ResultUseCase)
	var (
		err                error
		documentTypeResult model.DocumentType
	)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		data, err = s.ValidateDocumentType(data)
		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		if data.ID == "" {
			documentTypeResult, err = s.ProcessAddDocumentType(ctxReq, data)
		} else {
			documentTypeResult, err = s.ProcessUpdateDocumentType(ctxReq, data)
		}

		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		tags[helper.TextResponse] = documentTypeResult

		output <- ResultUseCase{Result: documentTypeResult}
	})

	return output
}

// ValidateDocumentType function validate data document type
func (s *DocumentUseCaseImpl) ValidateDocumentType(data model.DocumentType) (model.DocumentType, error) {
	var err error
	if len(data.DocumentType) == 0 {
		err := errors.New("document type required")
		return data, err
	}

	if len(data.IsB2cString) > 0 {
		data.IsB2c, err = strconv.ParseBool(data.IsB2cString)
		if err != nil {
			err = fmt.Errorf(helper.ErrorParameterInvalid, "isB2c")
			return data, err
		}
	}

	if len(data.IsB2bString) > 0 {
		data.IsB2b, err = strconv.ParseBool(data.IsB2bString)
		if err != nil {
			err = fmt.Errorf(helper.ErrorParameterInvalid, "isB2b")
			return data, err
		}
	}

	if len(data.IsActiveString) > 0 {
		data.IsActive, err = strconv.ParseBool(data.IsActiveString)
		if err != nil {
			err = fmt.Errorf(helper.ErrorParameterInvalid, "isB2c")
			return data, err
		}
	}

	return data, nil
}

// ProcessAddDocumentType function process document type
func (s *DocumentUseCaseImpl) ProcessAddDocumentType(ctxReq context.Context, data model.DocumentType) (model.DocumentType, error) {

	query := model.DocumentTypeParameters{
		DocumentType: data.DocumentType,
	}
	documentTypeResult := <-s.DocumentTypeRepo.FindDocumentTypeByParam(ctxReq, &query)
	if documentTypeResult.Result != nil {
		err := errors.New("document type already exists")
		return data, err
	}

	// add document repository process to database
	saveResult := <-s.DocumentTypeRepo.AddDocumentType(ctxReq, data)
	if saveResult.Error != nil {
		err := errors.New(msgErrorSaveType)
		return data, err
	}

	resultData, ok := saveResult.Result.(model.DocumentType)
	if !ok {
		err := errors.New(msgErrorSaveType)
		return data, err
	}

	return resultData, nil
}

// ProcessUpdateDocumentType function process document type
func (s *DocumentUseCaseImpl) ProcessUpdateDocumentType(ctxReq context.Context, data model.DocumentType) (model.DocumentType, error) {

	query := model.DocumentTypeParameters{
		ID: data.ID,
	}

	documentTypeResult := <-s.DocumentTypeRepo.FindDocumentTypeByParam(ctxReq, &query)
	if documentTypeResult.Result == nil {
		err := errors.New("document type doesn't exists")
		return data, err
	}

	documentType, ok := documentTypeResult.Result.(model.DocumentType)
	if !ok {
		err := errors.New("document type doesn't exists")
		return data, err
	}

	if documentType.DocumentType != data.DocumentType {
		query := model.DocumentTypeParameters{
			DocumentType: data.DocumentType,
		}

		documentTypeResult := <-s.DocumentTypeRepo.FindDocumentTypeByParam(ctxReq, &query)
		if documentTypeResult.Result != nil {
			err := errors.New("document type already exists")
			return data, err
		}
	}

	// update document repository process to database
	saveResult := <-s.DocumentTypeRepo.UpdateDocumentType(ctxReq, data)
	if saveResult.Error != nil {
		err := errors.New(msgErrorUpdateType)
		return data, err
	}

	resultData, ok := saveResult.Result.(model.DocumentType)
	if !ok {
		err := errors.New(msgErrorUpdateType)
		return data, err
	}

	return resultData, nil
}

// GetListDocumentType function for getting list of document
func (s *DocumentUseCaseImpl) GetListDocumentType(ctxReq context.Context, params *model.DocumentTypeParameters) <-chan ResultUseCase {
	ctx := "DocumentUseCase-GetListDocumentType"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		var err error

		paging, err := helper.ValidatePagination(
			helper.PaginationParameters{
				Page:     1, // default
				StrPage:  params.StrPage,
				Limit:    20, // default
				StrLimit: params.StrLimit,
			})

		if err != nil {
			tags[helper.TextResponse] = err.Error()
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		params.Page = paging.Page
		params.Limit = paging.Limit
		params.Offset = paging.Offset
		tags[helper.TextParameter] = params

		documentResult := <-s.DocumentTypeRepo.GetListDocumentType(ctxReq, params)
		if documentResult.Error != nil {
			httpStatus := http.StatusInternalServerError

			// when data is not found
			if documentResult.Error == sql.ErrNoRows {
				httpStatus = http.StatusNotFound
				documentResult.Error = fmt.Errorf(helper.ErrorDataNotFound, "document type")
			}
			tags[helper.TextResponse] = documentResult.Error
			output <- ResultUseCase{Error: documentResult.Error, HTTPStatus: httpStatus}
			return
		}

		document := documentResult.Result.(model.ListDocumentType)

		totalResult := <-s.DocumentTypeRepo.GetTotalDocumentType(ctxReq, params)
		if totalResult.Error != nil {
			tags[helper.TextResponse] = documentResult.Error
			output <- ResultUseCase{Error: totalResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		total := totalResult.Result.(int)
		document.TotalData = total
		tags[helper.TextResponse] = document
		output <- ResultUseCase{Result: document}
	})

	return output
}

// GetRequiredDocument function for getting list of document
func (s *DocumentUseCaseImpl) GetRequiredDocument(ctxReq context.Context) <-chan ResultUseCase {
	ctx := "DocumentUseCase-GetRequiredDocument"
	output := make(chan ResultUseCase)

	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		var err error

		resp, err := http.Get(os.Getenv("DOCUMENTS_JSON"))
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, "get_document_json", err, nil)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		defer resp.Body.Close()
		buf := new(bytes.Buffer)
		buf.ReadFrom(resp.Body)
		respByte := buf.Bytes()

		var data model.RequiredDocuments
		errorUnmarshal := json.Unmarshal([]byte(respByte), &data)
		if errorUnmarshal != nil {
			helper.SendErrorLog(ctxReq, ctx, "unmarshal_response", errorUnmarshal, string(respByte))
			output <- ResultUseCase{Error: errorUnmarshal, HTTPStatus: http.StatusBadRequest}
			return
		}
		tags[helper.TextResponse] = data
		output <- ResultUseCase{Result: data}
	})

	return output
}
