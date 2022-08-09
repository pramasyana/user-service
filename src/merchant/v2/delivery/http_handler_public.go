package delivery

import (
	"math"
	"net/http"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/labstack/echo"
)

const (
	merchantPublicPath = "/merchant/public"
)

func (m *HTTPMerchantHandler) MountMerchantPublic(group *echo.Group) {
	group.GET("/merchant/:merchantId/public", m.GetPublicDetailMerchant)
	group.GET("/merchant/:merchantId/warehouses/public", m.getPublicMerchantWarehouse)
	group.GET(merchantPublicPath+"/:vanityUrl", m.GetMerchantByVanity)
	group.GET(merchantPublicPath, m.GetListMerchantPublic)
}

func (m *HTTPMerchantHandler) getPublicMerchantWarehouse(c echo.Context) error {
	whPublicParams := model.ParameterWarehouse{
		StrPage:    c.QueryParam("page"),
		StrLimit:   c.QueryParam("limit"),
		MerchantID: c.Param(merchantIDParam),
		Query:      c.QueryParam("query"),
		Sort:       c.QueryParam("sort"),
		OrderBy:    c.QueryParam("order"),
		ShowAll:    c.QueryParam("showAll"),
	}

	warehousesResults := <-m.MerchantAddressUseCase.GetWarehouseAddresses(c.Request().Context(), &whPublicParams)
	if warehousesResults.Error != nil {
		return shared.NewHTTPResponse(warehousesResults.HTTPStatus, warehousesResults.Error.Error(), make(helper.EmptySlice, 0)).JSON(c)
	}

	warehouses, ok := warehousesResults.Result.(model.ListWarehouse)
	if !ok {
		return shared.NewHTTPResponse(http.StatusBadRequest, "result is not list of warehouse", make(helper.EmptySlice, 0)).JSON(c)
	}

	totalPage := math.Ceil(float64(warehouses.TotalData) / float64(whPublicParams.Limit))

	if len(warehouses.WarehouseData) <= 0 {
		return shared.NewHTTPResponse(http.StatusOK, messageSuccessGet, make(helper.EmptySlice, 0)).JSON(c)
	}

	meta := shared.Meta{
		Page:         whPublicParams.Page,
		Limit:        whPublicParams.Limit,
		TotalRecords: warehouses.TotalData,
		TotalPages:   int(totalPage),
	}
	warehouses.RestructToPublic()
	return shared.NewHTTPResponse(http.StatusOK, messageSuccessGet, meta, warehouses.WarehouseData).JSON(c)
}

func (m *HTTPMerchantHandler) GetPublicDetailMerchant(c echo.Context) error {
	merchantID := c.Param(merchantIDParam)
	isAttachment := "false"
	privacy := "public"
	// return merchant data
	merchantResult := <-m.MerchantUseCase.GetMerchantByID(c.Request().Context(), merchantID, privacy, isAttachment)
	if merchantResult.Error != nil {
		return shared.NewHTTPResponse(merchantResult.HTTPStatus, merchantResult.Error.Error(), helper.EmptySlice{}).JSON(c)
	}

	merchant, ok := merchantResult.Result.(model.B2CMerchantDataV2)
	if !ok {
		return shared.NewHTTPResponse(http.StatusBadRequest, "failed get merchant", make(helper.EmptySlice, 0)).JSON(c)
	}
	res := merchant.RestructForPublic()
	return shared.NewHTTPResponse(http.StatusOK, "Success get merchant", res).JSON(c)
}

func (m *HTTPMerchantHandler) GetMerchantByVanity(c echo.Context) error {
	vanityURL := c.Param(merchantVanityURL)
	// return merchant data
	merchantResult := <-m.MerchantUseCase.GetMerchantByVanityURL(c.Request().Context(), vanityURL)
	if merchantResult.Error != nil {
		return shared.NewHTTPResponse(merchantResult.HTTPStatus, merchantResult.Error.Error(), helper.EmptySlice{}).JSON(c)
	}

	merchant, ok := merchantResult.Result.(model.B2CMerchantDataV2)
	if !ok {
		return shared.NewHTTPResponse(http.StatusBadRequest, "failed get merchant", make(helper.EmptySlice, 0)).JSON(c)
	}
	res := merchant.RestructForPublic()
	return shared.NewHTTPResponse(http.StatusOK, "Success get merchant", res).JSON(c)
}

func (m *HTTPMerchantHandler) GetListMerchantPublic(c echo.Context) error {
	params := model.QueryParametersPublic{}
	if err := c.Bind(&params); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	if err := params.ValidatePublic(); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	result := <-m.MerchantUseCase.GetMerchantsPublic(c.Request().Context(), &params)
	if result.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, result.Error.Error()).JSON(c)
	}
	totalData, ok := result.TotalData.(int)
	if !ok {
		return shared.NewHTTPResponse(http.StatusBadRequest, "invalid result").JSON(c)
	}

	totalPage := math.Ceil(float64(totalData) / float64(params.Limit))

	meta := shared.Meta{
		Page:         params.Page,
		Limit:        params.Limit,
		TotalRecords: totalData,
		TotalPages:   int(totalPage),
	}

	return shared.NewHTTPResponse(http.StatusOK, "Success get merchants", result.Result, meta).JSON(c)
}
