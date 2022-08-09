package delivery

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/payments/v1/model"
	"github.com/Bhinneka/user-service/src/payments/v1/usecase"
	"github.com/Bhinneka/user-service/src/shared"
	"github.com/labstack/echo"
)

type HTTPPaymentsHandler struct {
	PaymentsUseCase usecase.PaymentsUseCase
}

func NewHTTPHandler(paymentsUseCase usecase.PaymentsUseCase) *HTTPPaymentsHandler {
	return &HTTPPaymentsHandler{PaymentsUseCase: paymentsUseCase}
}

func (p *HTTPPaymentsHandler) MountInfo(group *echo.Group) {
	group.POST("", p.AddUpdatePayment)
	group.GET("", p.GetDetailPayment)
}

func (p *HTTPPaymentsHandler) AddUpdatePayment(c echo.Context) error {
	paymentsInput := model.PaymentsInput{}
	headerAuth := c.Request().Header.Get(helper.TextAuthorization)

	if err := c.Bind(&paymentsInput); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	var auth string
	if split := strings.Split(headerAuth, " "); len(split) > 1 {
		auth = split[1]
	}
	payments := model.Payments{}
	payments.Channel = paymentsInput.Channel
	payments.Email = strings.ToLower(paymentsInput.Email)
	payments.Method = paymentsInput.Method
	payments.Token = paymentsInput.Token
	payments.ExpiredAt, _ = time.Parse(time.RFC3339, paymentsInput.ExpiredAt)

	checkResult := <-p.PaymentsUseCase.CompareHeaderAndBody(c.Request().Context(), &payments, auth)
	if checkResult.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, checkResult.Error.Error(), checkResult.ErrorData).JSON(c)
	}

	saveResult := <-p.PaymentsUseCase.AddUpdatePayments(c.Request().Context(), &payments)
	if saveResult.Error != nil {
		return shared.NewHTTPResponse(saveResult.HTTPStatus, saveResult.Error.Error(), saveResult.ErrorData).JSON(c)
	}

	result, ok := saveResult.Result.(model.SuccessResponse)
	if !ok {
		return shared.NewHTTPResponse(http.StatusInternalServerError, "save failed").JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusCreated, model.SaveSuccess, result).JSON(c)
}

func (p *HTTPPaymentsHandler) GetDetailPayment(c echo.Context) error {
	payments := model.Payments{}

	params := model.QueryPaymentParameters{}
	if err := c.Bind(&params); err != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, err.Error()).JSON(c)
	}
	payments.Channel = params.Channel
	payments.Email = params.Email
	payments.Method = params.Method

	headerAuth := c.Request().Header.Get(helper.TextAuthorization)
	var auth string
	if split := strings.Split(headerAuth, " "); len(split) > 1 {
		auth = split[1]
	}
	checkResult := <-p.PaymentsUseCase.CompareHeaderAndBody(c.Request().Context(), &payments, auth)
	if checkResult.Error != nil {
		return shared.NewHTTPResponse(http.StatusBadRequest, checkResult.Error.Error(), checkResult.ErrorData).JSON(c)
	}

	ctxReq := context.WithValue(c.Request().Context(), shared.ContextKey(helper.TextToken), c.Request().Header.Get(echo.HeaderAuthorization))
	getResult := <-p.PaymentsUseCase.GetPaymentDetail(ctxReq, &payments)
	if getResult.Error != nil {
		return shared.NewHTTPResponse(getResult.HTTPStatus, getResult.Error.Error(), getResult.ErrorData).JSON(c)
	}

	return shared.NewHTTPResponse(http.StatusOK, model.MessageGetData, getResult.Result).JSON(c)
}
