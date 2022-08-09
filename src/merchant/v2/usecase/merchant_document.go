package usecase

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	"gopkg.in/guregu/null.v4"
)

// InsertUpdateDocument function to insert or update merchant document
func (m *MerchantUseCaseImpl) InsertUpdateDocument(ctxReq context.Context, document model.B2CMerchantDocumentInput, merchantDocument model.B2CMerchantDocumentData) (model.B2CMerchantDocumentData, error) {
	var err error
	if document.ID == "" {
		merchantDocument.ID = helper.GenerateDocumentID()
		result := <-m.MerchantDocumentRepo.InsertNewMerchantDocument(ctxReq, &merchantDocument)
		if result.Error != nil {
			err = fmt.Errorf("Failed to save merchant document")
			return merchantDocument, err
		}
	} else {
		result := <-m.MerchantDocumentRepo.UpdateMerchantDocument(ctxReq, document.ID, &merchantDocument)
		if result.Error != nil {
			err = fmt.Errorf("Failed to update merchant document")
			return merchantDocument, err
		}
	}
	return merchantDocument, nil
}

// MerchantDocumentsProcess function to add merchant document
func (m *MerchantUseCaseImpl) MerchantDocumentsProcess(ctxReq context.Context, documents []model.B2CMerchantDocumentInput, merchantID string, userAttribute *model.MerchantUserAttribute) <-chan ResultUseCase {
	ctx := "MerchantUseCase-MerchantDocumentProcess"

	output := make(chan ResultUseCase)

	merchantDocuments := make([]model.B2CMerchantDocumentData, 0)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		var (
			wg sync.WaitGroup
		)
		chResp := make(chan ResultUseCase)

		tags["args"] = documents
		for _, document := range documents {
			document.MerchantID = merchantID
			wg.Add(1)
			go m.insertDoc(ctxReq, merchantID, &wg, document, userAttribute, chResp)
		}

		go func() {
			wg.Wait()
			close(chResp)
		}()

		for resp := range chResp {
			if resp.Error != nil {
				output <- ResultUseCase{Error: resp.Error, HTTPStatus: http.StatusBadRequest}
				return
			}
			doc := resp.Result.(model.B2CMerchantDocumentData)
			merchantDocuments = append(merchantDocuments, doc)
		}
		tags["docs"] = merchantDocuments

		output <- ResultUseCase{Result: merchantDocuments}
	})
	return output
}

func (m *MerchantUseCaseImpl) insertDoc(ctxReq context.Context, merchantID string, wg *sync.WaitGroup, document model.B2CMerchantDocumentInput, userAttribute *model.MerchantUserAttribute, response chan<- ResultUseCase) {
	tr := tracer.StartTrace(ctxReq, "MerchantUseCaseImpl-insertDoc")
	tags := map[string]interface{}{
		"docType":  document.DocumentType,
		"docValue": document.DocumentValue,
	}

	defer func() {
		wg.Done()
		tr.Finish(tags)
	}()

	// validate document
	errValidateDocument := m.ValidateDocument(tr.NewChildContext(), &document)
	if errValidateDocument != nil {
		response <- ResultUseCase{Error: errValidateDocument}
		return
	}
	merchantDocument := model.B2CMerchantDocumentData{
		ID:            document.ID,
		DocumentType:  document.DocumentType,
		DocumentValue: document.DocumentValue,
		MerchantID:    merchantID,
		Version:       1,
		CreatorIP:     userAttribute.UserIP,
		CreatorID:     userAttribute.UserID,
		EditorIP:      userAttribute.UserIP,
		EditorID:      userAttribute.UserID,
		LastModified:  null.TimeFrom(time.Now()),
	}

	if document.ID == "" {
		merchantDocument.Created = null.TimeFrom(time.Now())
	}

	merchantDocument, err := m.InsertUpdateDocument(tr.NewChildContext(), document, merchantDocument)
	if err != nil {
		response <- ResultUseCase{Error: err}
		return
	}
	response <- ResultUseCase{
		Result: merchantDocument,
	}
}

// ValidateDocument function to validate merchant document
func (m *MerchantUseCaseImpl) ValidateDocument(ctxReq context.Context, document *model.B2CMerchantDocumentInput) error {
	query := model.B2CMerchantDocumentQueryInput{
		MerchantID:   document.MerchantID,
		DocumentType: document.DocumentType,
	}
	// validate documentType is required
	if document.DocumentType == "" {
		return fmt.Errorf("document type is required")
	}

	check := <-m.MerchantDocumentRepo.FindMerchantDocumentByParam(ctxReq, &query)
	if check.Result != nil {
		merchantDocument := check.Result.(model.B2CMerchantDocumentData)
		parseUrl, _ := url.Parse(merchantDocument.DocumentValue)
		merchantDocument.DocumentOriginal = strings.TrimLeft(parseUrl.Path, "/")

		document.ID = merchantDocument.ID
		if merchantDocument.DocumentValue != "" && merchantDocument.MerchantID != document.MerchantID {
			return fmt.Errorf("merchant document already used")
		}
	}
	return nil
}
