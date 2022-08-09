package usecase

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"sync"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/gosimple/slug"
)

// GetMerchantByUserID function for get merchant by user ID
func (m *MerchantUseCaseImpl) GetMerchantByUserID(ctxReq context.Context, userID string, token string) <-chan ResultUseCase {
	ctx := "MerchantUseCase-GetMerchantByUserID"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		available := model.ResponseAvailable{}
		merchantResult := m.MerchantRepo.FindMerchantByUser(ctxReq, userID)
		if merchantResult.Result == nil {
			err := fmt.Errorf("Merchant not found")
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusNotFound}
			return
		}

		merchantData := merchantResult.Result.(model.B2CMerchantDataV2)
		mapsResult := <-m.MerchantAddressRepo.FindAddressMaps(ctxReq, merchantData.ID, "b2c_merchant")
		if mapsResult.Error == nil {
			merchantData.Maps, _ = mapsResult.Result.(model.Maps)
			merchantData.IsMapAvailable = helper.ValidateLatLong(merchantData.Maps.Latitude, merchantData.Maps.Longitude)
		}
		isAttachment := "false"
		tags[helper.TextMemberIDCamel] = userID
		tags[helper.TextEmail] = merchantData.MerchantEmail
		available.MerchantServiceAvailable = true
		merchantData = m.adjustMerchantData(ctxReq, merchantData, isAttachment)
		available.MerchantData = merchantData

		tags[helper.TextResponse] = available
		output <- ResultUseCase{Result: available}
	})

	return output
}

// adjustMerchantData function for get merchant by user ID
func (m *MerchantUseCaseImpl) adjustMerchantData(ctxReq context.Context, merchantData model.B2CMerchantDataV2, isAttachment string) model.B2CMerchantDataV2 {
	if merchantData.PicKtpFile.Valid && merchantData.PicKtpFile.String != "" {
		merchantData.PicKtpFile.String = m.getFile(ctxReq, merchantData.PicKtpFile.String, isAttachment)
	}

	if merchantData.NpwpFile.Valid && merchantData.NpwpFile.String != "" {
		merchantData.NpwpFile.String = m.getFile(ctxReq, merchantData.NpwpFile.String, isAttachment)
	}

	if merchantData.MerchantLogo.Valid && merchantData.MerchantLogo.String != "" {
		merchantData.MerchantLogo.String = m.getFile(ctxReq, merchantData.MerchantLogo.String, isAttachment)
	}

	merchantDocuments := make([]model.B2CMerchantDocumentData, 0)
	query := model.B2CMerchantDocumentQueryInput{
		MerchantID: merchantData.ID,
	}

	getMerchantDocuments := <-m.MerchantDocumentRepo.GetListMerchantDocument(ctxReq, &query)
	if getMerchantDocuments.Result != nil {
		resultDocuments := getMerchantDocuments.Result.(model.ListB2CMerchantDocument)
		merchantDocuments = m.getDocuments(ctxReq, resultDocuments, isAttachment)
	}

	merchantData.Documents = merchantDocuments
	if merchantData.StoreClosureDate != nil && merchantData.StoreReopenDate != nil {
		merchantData.IsClosed = merchantData.ChecksIsClosed()
	}
	return merchantData
}

func (m *MerchantUseCaseImpl) getDocuments(ctxReq context.Context, resultDocuments model.ListB2CMerchantDocument, isAttachment string) []model.B2CMerchantDocumentData {
	merchantDocuments := make([]model.B2CMerchantDocumentData, 0)

	maxWorker := 20 // 20 concurrent process for fetch url from upload service
	totalSplit := int(math.Ceil(float64(len(resultDocuments.MerchantDocument)) / float64(maxWorker)))
	documentBuff := make(chan struct {
		index int
		url   string
	})
	var wg sync.WaitGroup
	for i := 0; i < maxWorker; i++ {
		offset := i * totalSplit
		if offset > len(resultDocuments.MerchantDocument) {
			continue
		}
		lastOffset := offset + totalSplit
		if lastOffset > len(resultDocuments.MerchantDocument) {
			lastOffset = len(resultDocuments.MerchantDocument)
		}
		partData := resultDocuments.MerchantDocument[offset:lastOffset]

		wg.Add(1)
		go func(workerOffset int, data []model.B2CMerchantDocumentData) {
			defer wg.Done()
			defer func() { recover() }()
			for i, detail := range data {
				documentBuff <- struct {
					index int
					url   string
				}{
					index: i + workerOffset, url: m.getFile(ctxReq, detail.DocumentValue, isAttachment),
				}
			}
		}(offset, partData)
	}

	go func() { wg.Wait(); close(documentBuff) }()
	for resBuff := range documentBuff {
		resultDocuments.MerchantDocument[resBuff.index].DocumentValue = resBuff.url
	}

	merchantDocuments = append(merchantDocuments, resultDocuments.MerchantDocument...)

	return merchantDocuments
}

// getFile function for get url image
func (m *MerchantUseCaseImpl) getFile(ctxReq context.Context, url string, isAttachment string) string {
	if url == "" {
		return url
	}
	documentURL := <-m.UploadService.GetURLImage(ctxReq, url, isAttachment)
	if documentURL.Result != nil {
		documentURLResult, ok := documentURL.Result.(serviceModel.ResponseUploadService)
		if !ok || documentURLResult.Data.URL == "" {
			return url
		}
		return documentURLResult.Data.URL
	}
	return url
}

// CheckMerchantName usecase function for verify merchant name
func (m *MerchantUseCaseImpl) CheckMerchantName(ctxReq context.Context, merchantName string) <-chan ResultUseCase {
	ctx := "MerchantUseCase-CheckMerchantName"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags["args"] = merchantName
		if merchantName == "" {
			err := fmt.Errorf("merchantName required")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		emailResult := m.MerchantRepo.FindMerchantByName(ctxReq, merchantName)
		if emailResult.Result != nil {
			err := fmt.Errorf("%s already exist", merchantName)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		responseMerchantName := model.ResponseMerchantName{}
		responseMerchantName.MerchantName = merchantName
		responseMerchantName.Slug = slug.MakeLang(merchantName, "en")
		responseMerchantName.URL = "https://bhinneka.com/toko-" + responseMerchantName.Slug

		tags[helper.TextResponse] = responseMerchantName
		output <- ResultUseCase{Result: responseMerchantName}

	})

	return output
}
