package consumer

import (
	"context"
	"encoding/json"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	memberUC "github.com/Bhinneka/user-service/src/member/v1/usecase"
	merchantModel "github.com/Bhinneka/user-service/src/merchant/v2/model"
	merchantUC "github.com/Bhinneka/user-service/src/merchant/v2/usecase"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	shippingModel "github.com/Bhinneka/user-service/src/shipping_address/v2/model"
	shippingUC "github.com/Bhinneka/user-service/src/shipping_address/v2/usecase"
)

func Dispatch(ctxReq context.Context, payload []byte, merchantUsecase merchantUC.MerchantUseCase, memberUsecase memberUC.MemberUseCase, shippingUsecase shippingUC.ShippingAddressUseCase) (err error) {
	ctx := "SturgeonDispatcher"
	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{}
	defer tr.Finish(tags)

	data := serviceModel.QueuePayload{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return err
	}
	tags["eventType"] = data.EventType
	newCtx := context.WithValue(tr.NewChildContext(), helper.TextAuthorization, data.Auth)
	switch data.EventType {
	case "SendEmailMerchantAdd", "SendEmailMerchantRejectRegistration", "SendEmailMerchantRejectUpgrade", "SendEmailMerchantUpgrade", "SendEmailAdmin":
		err = sendEmailMerchant(newCtx, data.Payload, merchantUsecase, data.EventType)
	case "SendEmailMerchantEmployeeLogin", "SendEmailMerchantEmployeeRegister":
		err = sendEmailMerchantEmployee(newCtx, data.Payload, merchantUsecase, data.EventType)
	case "InsertLogMerchantCreate":
		err = insertLogMerchant(newCtx, data.Payload, merchantUsecase, helper.TextInsertUpper)
	case "InsertLogMerchantDelete":
		err = insertLogMerchant(newCtx, data.Payload, merchantUsecase, helper.TextDeleteUpper)
	case "InsertLogMerchantUpdate":
		err = insertLogMerchant(newCtx, data.Payload, merchantUsecase, helper.TextUpdateUpper)
	case "InsertLogRegisterMember", "InsertLogAuth":
		err = insertLogMember(newCtx, data.Payload, memberUsecase, helper.TextInsertUpper)
	case "InsertLogUpdateMember":
		err = insertLogMember(newCtx, data.Payload, memberUsecase, helper.TextUpdateUpper)
	case "SendEmailRegisterMember", "SendEmailWelcomeMember", "SendEmailSuccessForgotPassword", "SendEmailForgotPassword", "SendEmailAddMember":
		err = sendEmailMember(newCtx, data.Payload, memberUsecase, data.EventType)
	case "AddShippingAddress":
		err = insertLogShipping(newCtx, data.Payload, shippingUsecase, helper.TextInsertUpper)
	case "UpdatePrimaryShippingAddressByID", "UpdateShippingAddress":
		err = insertLogShipping(newCtx, data.Payload, shippingUsecase, helper.TextUpdateUpper)
	case "DeleteShippingAddressByID":
		err = insertLogShipping(newCtx, data.Payload, shippingUsecase, helper.TextDeleteUpper)
	}
	if err != nil {
		helper.SendErrorLog(tr.Context(), "SturgeonWorker", "exec_dispatcher", err, data.Payload)
		helper.SendNotification(data.EventType, string(payload), ctx, err)
		return err
	}
	return nil
}

func sendEmailMerchant(ctxReq context.Context, payload interface{}, merchantUsecase merchantUC.MerchantUseCase, eventType string) error {
	merchantPl := merchantModel.MerchantPayloadEmail{}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &merchantPl); err != nil {
		return err
	}

	switch eventType {
	case "SendEmailMerchantAdd":
		result := <-merchantUsecase.SendEmailMerchantAdd(ctxReq, merchantPl.Data, merchantPl.MemberName)
		return result.Error
	case "SendEmailMerchantRejectRegistration":
		result := <-merchantUsecase.SendEmailMerchantRejectRegistration(ctxReq, merchantPl.Data, merchantPl.MemberName)
		return result.Error
	case "SendEmailMerchantRejectUpgrade":
		result := <-merchantUsecase.SendEmailMerchantRejectUpgrade(ctxReq, merchantPl.Data, merchantPl.MemberName, merchantPl.ReasonReject)
		return result.Error
	case "SendEmailActivation":
		result := <-merchantUsecase.SendEmailActivation(ctxReq, merchantPl.Data)
		return result.Error
	case "SendEmailApproval":
		result := <-merchantUsecase.SendEmailApproval(ctxReq, merchantPl.Data)
		return result.Error
	case "SendEmailMerchantUpgrade":
		result := <-merchantUsecase.SendEmailMerchantUpgrade(ctxReq, merchantPl.Data, merchantPl.MemberName)
		return result.Error
	case "SendEmailAdmin":
		result := <-merchantUsecase.SendEmailAdmin(ctxReq, merchantPl.Data, merchantPl.MemberName, merchantPl.ReasonReject, merchantPl.AdminCMS)
		return result.Error
	}

	return nil
}

func sendEmailMerchantEmployee(ctxReq context.Context, payload interface{}, merchantUsecase merchantUC.MerchantUseCase, eventType string) error {

	merchantPl := memberModel.MemberPayloadEmail{}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &merchantPl); err != nil {
		return err
	}

	switch eventType {
	case "SendEmailMerchantEmployeeLogin":
		result := <-merchantUsecase.SendEmailMerchantEmployeeLogin(ctxReq, *merchantPl.Merchant, *merchantPl.Member)
		return result.Error
	case "SendEmailMerchantEmployeeRegister":
		result := <-merchantUsecase.SendEmailMerchantEmployeeRegister(ctxReq, *merchantPl.Merchant, *merchantPl.Member)
		return result.Error
	}

	return nil
}

func sendEmailMember(ctxReq context.Context, payload interface{}, memberUsecase memberUC.MemberUseCase, eventType string) error {
	memberEmailPayload := memberModel.MemberEmailQueue{}
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(b, &memberEmailPayload); err != nil {
		return err
	}
	switch eventType {
	case "SendEmailRegisterMember":
		result := <-memberUsecase.SendEmailRegisterMember(ctxReq, memberEmailPayload.Data, memberEmailPayload.Member.RegisterType)
		return result.Error
	case "SendEmailWelcomeMember":
		result := <-memberUsecase.SendEmailWelcomeMember(ctxReq, memberEmailPayload.Data)
		return result.Error
	case "SendEmailSuccessForgotPassword":
		result := <-memberUsecase.SendEmailSuccessForgotPassword(ctxReq, memberEmailPayload.Data)
		return result.Error
	case "SendEmailForgotPassword":
		result := <-memberUsecase.SendEmailForgotPassword(ctxReq, memberEmailPayload.Data)
		return result.Error
	case "SendEmailAddMember":
		result := <-memberUsecase.SendEmailAddMember(ctxReq, memberEmailPayload.Data)
		return result.Error
	}
	return nil
}

// payload is type of MerchantLog
// type MerchantLog struct {
// Before B2CMerchantDataV2 `json:"before"`
// After  B2CMerchantDataV2 `json:"after"`
// }
func insertLogMerchant(ctxReq context.Context, payload interface{}, merchantUsecase merchantUC.MerchantUseCase, eventType string) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	merchantData := merchantModel.MerchantLog{}
	if err := json.Unmarshal(b, &merchantData); err != nil {
		return err
	}
	return merchantUsecase.InsertLogMerchant(ctxReq, merchantData.Before, merchantData.After, eventType)
}

func insertLogMember(ctxReq context.Context, payload interface{}, memberUsecase memberUC.MemberUseCase, eventType string) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	memberData := memberModel.MemberLog{}
	if err := json.Unmarshal(b, &memberData); err != nil {
		return err
	}
	return memberUsecase.InsertLogMember(ctxReq, memberData.Before, memberData.After, eventType)
}

func insertLogShipping(ctxReq context.Context, payload interface{}, shippingUsecase shippingUC.ShippingAddressUseCase, eventType string) error {
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	shippinglog := shippingModel.ShippingAddressLog{}
	if err := json.Unmarshal(b, &shippinglog); err != nil {
		return err
	}

	return shippingUsecase.InsertLogShipping(ctxReq, shippinglog.Before, shippinglog.After, eventType)
}
