package usecase

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/eapache/go-resiliency/retrier"
)

const (
	merchantPlaceholder      = "##MERCHANTNAME##"
	fullnamePlaceholder      = "##FULLNAME##"
	urlPlaceholder           = "##URL##"
	yearPlaceholder          = "##YEAR##"
	upgradeStatusPlaceholder = "##UPGRADE_STATUS##"
	reasonRejectText         = "##REASON_REJECT##"
	adminName                = "##ADMIN_NAME##"

	textErrorSturgeonCFURL = "you need to specify %s in the environment variable"
)

// SendEmailMerchantAdd usecase function for send email merchant registration
func (m *MerchantUseCaseImpl) SendEmailMerchantAdd(ctxReq context.Context, data model.B2CMerchantDataV2, fullName string) <-chan ResultUseCase {
	ctx := "MerchantUseCase-SendEmailMerchantAdd"
	output := make(chan ResultUseCase)
	go tracer.WithTrace(ctxReq, ctx, nil, func(ctxReq context.Context) {
		defer close(output)

		// get template email
		templateKeyID := "EMAIL_MERCHANT_REGISTRATION_TEMPLATE_ID"
		templateEmailDetail, errTemplate := m.GetTemplateEmail(ctxReq, templateKeyID)
		if errTemplate != nil {
			output <- ResultUseCase{Error: errTemplate, HTTPStatus: http.StatusBadRequest}
			return
		}

		// set emailContent to our email
		emailContents := templateEmailDetail.Content
		emailContents = strings.Replace(emailContents, merchantPlaceholder, data.MerchantName, -1)

		bCCEmail, err := m.getBCC()
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		emailPayload := serviceModel.Email{}
		emailPayload.From = serviceModel.EmailCare
		emailPayload.FromName = serviceModel.NoReplyName
		emailPayload.To = []string{data.MerchantEmail.String}
		emailPayload.ToName = []string{fullName}
		if bCCEmail != "" {
			emailPayload.BCC = []string{bCCEmail}
			emailPayload.BCCName = []string{serviceModel.NoReplyName}
		}
		emailPayload.Subject = model.MerchantRegistrationSubject
		emailPayload.Content = emailContents

		attachments := []serviceModel.Attachment{}
		attachmentFile, _ := os.LookupEnv("EMAIL_ATTACHMENT_SLA_MERCHANT_REGULAR")
		slaFile := m.getAttachmentFile(ctxReq, attachmentFile, "SERVICE LEVEL AGREEMENT (SLA) REGULAR MERCHANT.PDF")
		attachments = append(attachments, slaFile)

		attachmentFile, _ = os.LookupEnv("EMAIL_ATTACHMENT_TERMS_MERCHANT")
		slaFile = m.getAttachmentFile(ctxReq, attachmentFile, "KETENTUAN STANDAR KERJASAMA PENYEDIA BARANG JASA.PDF")
		attachments = append(attachments, slaFile)

		emailPayload.Attachments = attachments

		if err := m.sendEmailMerchant(ctxReq, emailPayload); err != nil {
			errMessage := errors.New(errMsgFailedSendEmail)
			output <- ResultUseCase{Error: errMessage, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: data}
	})

	return output
}

// SendEmailMerchantUpgrade usecase function for send email merchant upgrade
func (m *MerchantUseCaseImpl) SendEmailMerchantUpgrade(ctxReq context.Context, data model.B2CMerchantDataV2, memberName string) <-chan ResultUseCase {
	ctx := "MerchantUseCase-SendEmailMerchantUpgrade"
	output := make(chan ResultUseCase)
	go tracer.WithTrace(ctxReq, ctx, nil, func(ctxReq context.Context) {
		defer close(output)

		// get template email
		templateEmailDetail, errTemplate := m.GetTemplateEmail(ctxReq, "EMAIL_MERCHANT_UPGRADE_TEMPLATE_ID")
		if errTemplate != nil {
			output <- ResultUseCase{Error: errTemplate, HTTPStatus: http.StatusBadRequest}
			return
		}

		// set emailContent to our email
		emailContent := templateEmailDetail.Content
		emailContent = strings.Replace(emailContent, merchantPlaceholder, data.MerchantName, -1)

		bCCEmail, err := m.getBCC()
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		pl := serviceModel.Email{}
		pl.From = serviceModel.EmailCare
		pl.FromName = serviceModel.NoReplyName
		pl.To = []string{data.MerchantEmail.String}
		pl.ToName = []string{memberName}
		if bCCEmail != "" {
			pl.BCC = []string{bCCEmail}
			pl.BCCName = []string{serviceModel.NoReplyName}
		}
		pl.Subject = model.MerchantUpgradeSubject
		pl.Content = emailContent

		attachment, err := m.getAttachmentFileUpgrade(ctxReq, data.UpgradeStatus.String)
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		pl.Attachments = attachment

		err = m.sendEmailMerchant(ctxReq, pl)
		if err != nil {
			output <- ResultUseCase{Error: errors.New(errMsgFailedSendEmail), HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: data}
	})

	return output
}

func (m *MerchantUseCaseImpl) getAttachmentFileUpgrade(ctxReq context.Context, upgradeStatus string) ([]serviceModel.Attachment, error) {
	attachments := []serviceModel.Attachment{}
	if upgradeStatus == model.PendingManageString {
		attachmentFile, ok := os.LookupEnv("EMAIL_ATTACHMENT_SLA_MERCHANT_MANAGED")
		if !ok {
			return attachments, errors.New("unable to find EMAIL_ATTACHMENT_SLA_MERCHANT_MANAGED")
		}
		slaFile := m.getAttachmentFile(ctxReq, attachmentFile, "SERVICE LEVEL AGREEMENT (SLA) MANAGED MERCHANT.PDF")
		attachments = append(attachments, slaFile)
	} else if upgradeStatus == model.PendingAssociateString {
		attachmentFile, ok := os.LookupEnv("EMAIL_ATTACHMENT_SLA_MERCHANT_ASSOCIATE")
		if !ok {
			return attachments, errors.New("unable to find EMAIL_ATTACHMENT_SLA_MERCHANT_ASSOCIATE")
		}
		slaFile := m.getAttachmentFile(ctxReq, attachmentFile, "SERVICE LEVEL AGREEMENT (SLA) ASSOCIATE MERCHANT.PDF")
		attachments = append(attachments, slaFile)
	}
	return attachments, nil
}

// sendEmailMerchant usecase function for send email merchant
func (m *MerchantUseCaseImpl) sendEmailMerchant(ctxReq context.Context, pl serviceModel.Email) error {
	tr := tracer.StartTrace(ctxReq, "MerchantUseCaseImpl-sendEmailMerchant")
	tags := map[string]interface{}{"email": pl.To[0]}
	defer tr.Finish(tags)

	r := retrier.New(retrier.ConstantBackoff(3, 100*time.Millisecond), nil)
	err := r.RunCtx(ctxReq, func(_ context.Context) error {
		_, err := m.NotificationService.SendEmail(ctxReq, pl)
		return err
	})

	if err != nil {
		helper.SendErrorLog(ctxReq, "sendEmailMerchant", "send_email_merchant", err, pl)
		errMessage := errors.New(errMsgFailedSendEmail)
		return errMessage
	}
	return nil
}

// sendEmailMerchant usecase function for send email merchant
func (m *MerchantUseCaseImpl) sendEmailAdmin(ctxReq context.Context, pl serviceModel.Email) error {
	tr := tracer.StartTrace(ctxReq, "MerchantUseCaseImpl-sendEmailAdmin")
	tags := map[string]interface{}{"email": pl.To[0]}
	defer tr.Finish(tags)

	r := retrier.New(retrier.ConstantBackoff(3, 100*time.Millisecond), nil)
	err := r.RunCtx(ctxReq, func(_ context.Context) error {
		_, err := m.NotificationService.SendEmail(ctxReq, pl)
		return err
	})

	if err != nil {
		helper.SendErrorLog(ctxReq, "sendEmailAdmin", "send_email_admin", err, pl)
		errMessage := errors.New(errMsgFailedSendEmail)
		return errMessage
	}
	return nil
}

func (m *MerchantUseCaseImpl) getAttachmentFile(ctxReq context.Context, attachmentFile, fileName string) serviceModel.Attachment {
	ctx := "MerchantUseCase-getAttachmentFile"
	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{"fileName": fileName}
	defer tr.Finish(tags)

	attach := serviceModel.Attachment{}
	resp, err := http.Get(attachmentFile)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "get_attachment_file", err, attachmentFile)
		return attach
	}

	defer resp.Body.Close()
	buf := new(bytes.Buffer)
	buf.ReadFrom(resp.Body)
	base64File := base64.StdEncoding.EncodeToString([]byte(buf.Bytes()))

	//Attachment file sla
	attach = serviceModel.Attachment{
		Type:    "application/pdf",
		Name:    fileName,
		Content: base64File,
	}

	return attach
}

// GetTemplateEmail usecase function for get template email
func (m *MerchantUseCaseImpl) GetTemplateEmail(ctxReq context.Context, templateKey string) (serviceModel.Template, error) {
	templateEmailDetail := serviceModel.Template{}
	EmailMerchantRegistrationTemplateID, ok := os.LookupEnv(templateKey)
	if !ok {
		return templateEmailDetail, fmt.Errorf(textErrorSturgeonCFURL, templateKey)
	}

	templateEmailStatics := <-m.NotificationService.GetTemplateByID(ctxReq, EmailMerchantRegistrationTemplateID, templateKey)
	if templateEmailStatics.Error != nil {
		return templateEmailDetail, templateEmailStatics.Error
	}

	// detail
	templateEmailDetail, ok = templateEmailStatics.Result.(serviceModel.Template)
	if !ok {
		return templateEmailDetail, errors.New(helper.ErrMsgFailedGetTemplate)
	}

	// check template email not null
	if templateEmailDetail.Content == "" {
		return templateEmailDetail, errors.New(helper.ErrMsgEmptyContent)
	}

	return templateEmailDetail, nil
}

// SendEmailMerchantAdd usecase function for send email merchant registration
func (m *MerchantUseCaseImpl) SendEmailApproval(ctxReq context.Context, old model.B2CMerchantDataV2) <-chan ResultUseCase {
	ctx := "MerchantUseCase-SendEmailApproval"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		// get template email
		templateID := "EMAIL_MERCHANT_UPGRADE_APPROVAL"
		templateEmailDetail, errTemplate := m.GetTemplateEmail(ctxReq, templateID)
		if errTemplate != nil {
			output <- ResultUseCase{Error: errTemplate, HTTPStatus: http.StatusBadRequest}
			return
		}
		tags["email"] = old.MerchantEmail.String

		merchantName := old.MerchantName
		emailContent := templateEmailDetail.Content
		emailContent = strings.Replace(emailContent, merchantPlaceholder, old.MerchantName, -1)

		bCCEmail, err := m.getBCC()
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		pl := serviceModel.Email{}
		pl.From = serviceModel.EmailCare
		pl.FromName = serviceModel.NoReplyName
		pl.To = []string{old.MerchantEmail.String}
		pl.ToName = []string{merchantName}
		if bCCEmail != "" {
			pl.BCC = []string{bCCEmail}
			pl.BCCName = []string{serviceModel.NoReplyName}
		}
		pl.Subject = model.MerchantUpgradeApproved
		pl.Content = emailContent

		if err = m.sendEmailMerchant(ctxReq, pl); err != nil {
			errMessage := errors.New(errMsgFailedSendEmail)
			output <- ResultUseCase{Error: errMessage, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: old}
	})

	return output
}

func (m *MerchantUseCaseImpl) getBCC() (string, error) {
	bCCEmail, ok := os.LookupEnv("EMAIL_BCC_HUNTER")
	if !ok {
		return "", errors.New("you need to specify EMAIL_BCC_HUNTER in the environment variable")
	}
	if golib.StringInSlice(os.Getenv("ENV"), []string{helper.EnvDev, helper.EnvStaging}, false) {
		return "", nil
	}
	return bCCEmail, nil
}

// sendEmailActivation usecase function for send email merchant registration
func (m *MerchantUseCaseImpl) SendEmailActivation(ctxReq context.Context, merchant model.B2CMerchantDataV2) <-chan ResultUseCase {
	ctx := "MerchantUseCase-sendEmailActivation"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		// get template email
		templateID := "EMAIL_MERCHANT_ACTIVATION"
		templateEmailDetail, errTemplate := m.GetTemplateEmail(ctxReq, templateID)
		if errTemplate != nil {
			output <- ResultUseCase{Error: errTemplate, HTTPStatus: http.StatusBadRequest}
			return
		}
		tags["templateId"] = templateID
		tags["merchantId"] = merchant.ID

		merchantName := merchant.MerchantName
		emailContent := templateEmailDetail.Content
		emailContent = strings.Replace(emailContent, merchantPlaceholder, merchant.MerchantName, -1)

		bCCEmail, err := m.getBCC()
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		payloadActivation := serviceModel.Email{}
		payloadActivation.From = serviceModel.EmailCare
		payloadActivation.FromName = serviceModel.NoReplyName
		payloadActivation.To = []string{merchant.MerchantEmail.String}
		payloadActivation.ToName = []string{merchantName}
		if bCCEmail != "" {
			payloadActivation.BCC = []string{bCCEmail}
			payloadActivation.BCCName = []string{serviceModel.NoReplyName}
		}
		payloadActivation.Subject = model.MerchantActivated
		payloadActivation.Content = emailContent

		if err = m.sendEmailMerchant(ctxReq, payloadActivation); err != nil {
			errMessage := errors.New(errMsgFailedSendEmail)
			output <- ResultUseCase{Error: errMessage, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: merchant}
	})

	return output
}

func (m *MerchantUseCaseImpl) sendEmailActivationOrApproval(ctxReq context.Context, oldData model.B2CMerchantDataV2, payload *model.B2CMerchantCreateInput) error {
	if !oldData.IsActive && payload.IsActive {
		plQueue := model.MerchantPayloadEmail{
			MemberName: oldData.MerchantName,
			Data:       oldData,
		}
		go m.QueuePublisher.QueueJob(ctxReq, plQueue, oldData.ID, "SendEmailActivation")
	}

	// send email only if current upgradeStatus is not ACTIVE and new payload is ACTIVE
	if oldData.UpgradeStatus.String != model.ActiveString && payload.UpgradeStatus == model.ActiveString {
		// prevent repetitive sending email if merchant being updated
		if oldData.UpgradeStatus.String == model.ActiveString && payload.UpgradeStatus == model.ActiveString {
			// skip send email, current status is already active
			return nil
		}

		plQueue := model.MerchantPayloadEmail{
			MemberName: oldData.MerchantName,
			Data:       oldData,
		}
		go m.QueuePublisher.QueueJob(ctxReq, plQueue, oldData.ID, "SendEmailApproval")
	}
	return nil
}

// SendEmailMerchantRejectRegistration usecase function for send email merchant reject registration
func (m *MerchantUseCaseImpl) SendEmailMerchantRejectRegistration(ctxReq context.Context, data model.B2CMerchantDataV2, fullName string) <-chan ResultUseCase {
	ctx := "MerchantUseCase-SendEmailRejectRegistration"
	output := make(chan ResultUseCase)
	go tracer.WithTrace(ctxReq, ctx, nil, func(ctxReq context.Context) {
		defer close(output)

		// get template email
		templateKey := "EMAIL_MERCHANT_REJECT"
		templateEmailDetail, errTemplate := m.GetTemplateEmail(ctxReq, templateKey)
		if errTemplate != nil {
			output <- ResultUseCase{Error: errTemplate, HTTPStatus: http.StatusBadRequest}
			return
		}

		// set emailContent to our email
		emailContent := templateEmailDetail.Content
		emailContent = strings.Replace(emailContent, merchantPlaceholder, data.MerchantName, -1)

		bCCEmail, err := m.getBCC()
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		payloadReject := serviceModel.Email{}
		payloadReject.From = serviceModel.EmailCare
		payloadReject.FromName = serviceModel.NoReplyName
		payloadReject.To = []string{data.MerchantEmail.String}
		payloadReject.ToName = []string{fullName}
		if bCCEmail != "" {
			payloadReject.BCC = []string{bCCEmail}
			payloadReject.BCCName = []string{serviceModel.NoReplyName}
		}
		payloadReject.Subject = model.MerchantRegistrationRejectSubject
		payloadReject.Content = emailContent

		if err := m.sendEmailMerchant(ctxReq, payloadReject); err != nil {
			errMessage := errors.New(errMsgFailedSendEmail)
			output <- ResultUseCase{Error: errMessage, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: data}
	})

	return output
}

// SendEmailMerchantRejectUpgrade usecase function for send email merchant reject upgrade request
func (m *MerchantUseCaseImpl) SendEmailMerchantRejectUpgrade(ctxReq context.Context, data model.B2CMerchantDataV2, fullName string, reasonReject string) <-chan ResultUseCase {
	ctx := "MerchantUseCase-SendEmailMerchantRejectUpgrade"
	output := make(chan ResultUseCase)
	go tracer.WithTrace(ctxReq, ctx, nil, func(ctxReq context.Context) {
		defer close(output)

		// get template email

		templateKeyRejectRequest := "EMAIL_MERCHANT_UPGRADE_REJECT"
		if reasonReject != "" {
			templateKeyRejectRequest = "EMAIL_MERCHANT_UPGRADE_REJECT_WITH_REASON"
		}

		templateEmailDetail, errTemplate := m.GetTemplateEmail(ctxReq, templateKeyRejectRequest)
		if errTemplate != nil {
			output <- ResultUseCase{Error: errTemplate, HTTPStatus: http.StatusBadRequest}
			return
		}
		var upgradeStatus string
		if data.UpgradeStatus.String == model.PendingAssociateString {
			upgradeStatus = model.AssociateString
		} else if data.UpgradeStatus.String == model.PendingManageString {
			upgradeStatus = model.ManageString
		}

		emailContentRejectRequest := templateEmailDetail.Content
		emailContentRejectRequest = strings.Replace(emailContentRejectRequest, merchantPlaceholder, data.MerchantName, -1)
		emailContentRejectRequest = strings.Replace(emailContentRejectRequest, upgradeStatusPlaceholder, upgradeStatus, 1)
		if reasonReject != "" {
			emailContentRejectRequest = strings.Replace(emailContentRejectRequest, reasonRejectText, reasonReject, 1)
		}
		payloadRejectRequest := serviceModel.Email{}
		payloadRejectRequest.From = serviceModel.EmailCare
		payloadRejectRequest.FromName = serviceModel.NoReplyName
		payloadRejectRequest.To = []string{data.MerchantEmail.String}
		payloadRejectRequest.ToName = []string{fullName}

		payloadRejectRequest.Subject = model.MerchantUpgradeRejectSubject
		payloadRejectRequest.Content = emailContentRejectRequest

		if err := m.sendEmailMerchant(ctxReq, payloadRejectRequest); err != nil {
			errMessage := errors.New(errMsgFailedSendEmail)
			output <- ResultUseCase{Error: errMessage, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: data}
	})

	return output
}

// SendEmailMerchantEmployeeLogin usecase function for send email merchant upgrade
func (m *MerchantUseCaseImpl) SendEmailMerchantEmployeeLogin(ctxReq context.Context, dataMerchant model.B2CMerchantDataV2, dataMember memberModel.Member) <-chan ResultUseCase {
	ctx := "MerchantUseCase-SendEmailMerchantEmployeeLogin"
	output := make(chan ResultUseCase)
	go tracer.WithTrace(ctxReq, ctx, nil, func(ctxReq context.Context) {
		defer close(output)

		// get template email
		templateEmailDetail, errTemplate := m.GetTemplateEmail(ctxReq, "EMAIL_MERCHANT_EMPLOYEE_LOGIN")
		if errTemplate != nil {
			output <- ResultUseCase{Error: errTemplate, HTTPStatus: http.StatusBadRequest}
			return
		}

		// set emailContent to our email
		emailContent := templateEmailDetail.Content

		year := time.Now().Format("2006")
		url, ok := os.LookupEnv("STURGEON_CF_URL")
		if !ok {
			output <- ResultUseCase{Error: fmt.Errorf(textErrorSturgeonCFURL, "STURGEON_CF_URL"), HTTPStatus: http.StatusBadRequest}
			return
		}

		strURL := fmt.Sprintf("%s/employeeactivation?token=%s", url, dataMember.Token)

		find := []string{merchantPlaceholder, fullnamePlaceholder, urlPlaceholder, yearPlaceholder}
		replacer := []string{dataMerchant.MerchantName, dataMember.FirstName + " " + dataMember.LastName, strURL, year}
		emailContent = golib.StringArrayReplace(emailContent, find, replacer)

		pl := serviceModel.Email{}
		pl.From = serviceModel.NoReply
		pl.FromName = serviceModel.NoReplyName
		pl.To = []string{dataMember.Email}
		pl.ToName = []string{dataMember.FirstName + " " + dataMember.LastName}
		pl.Subject = fmt.Sprintf("Seller %s Mengundang Anda!", dataMerchant.MerchantName)
		pl.Content = emailContent

		err := m.sendEmailMerchant(ctxReq, pl)
		if err != nil {
			output <- ResultUseCase{Error: errors.New(errMsgFailedSendEmail), HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: nil}
	})

	return output
}

// SendEmailMerchantEmployeeRegister usecase function for send email merchant upgrade
func (m *MerchantUseCaseImpl) SendEmailMerchantEmployeeRegister(ctxReq context.Context, dataMerchant model.B2CMerchantDataV2, dataMember memberModel.Member) <-chan ResultUseCase {
	ctx := "MerchantUseCase-SendEmailMerchantEmployeeRegister"
	output := make(chan ResultUseCase)
	go tracer.WithTrace(ctxReq, ctx, nil, func(ctxReq context.Context) {
		defer close(output)

		// get template email
		templateEmailDetail, errTemplate := m.GetTemplateEmail(ctxReq, "EMAIL_MERCHANT_EMPLOYEE_REGISTER")
		if errTemplate != nil {
			output <- ResultUseCase{Error: errTemplate, HTTPStatus: http.StatusBadRequest}
			return
		}

		// set emailContent to our email
		emailContent := templateEmailDetail.Content

		year := time.Now().Format("2006")
		url, ok := os.LookupEnv("STURGEON_CF_URL")
		if !ok {
			output <- ResultUseCase{Error: fmt.Errorf(textErrorSturgeonCFURL, "STURGEON_CF_URL"), HTTPStatus: http.StatusBadRequest}
			return
		}

		strURL := fmt.Sprintf("%s/activate-member?token=%s", url, dataMember.Token)

		find := []string{merchantPlaceholder, fullnamePlaceholder, urlPlaceholder, yearPlaceholder}
		replacer := []string{dataMerchant.MerchantName, dataMember.FirstName + " " + dataMember.LastName, strURL, year}
		emailContent = golib.StringArrayReplace(emailContent, find, replacer)

		pl := serviceModel.Email{}
		pl.From = serviceModel.NoReply
		pl.FromName = serviceModel.NoReplyName
		pl.To = []string{dataMember.Email}
		pl.ToName = []string{dataMember.FirstName + " " + dataMember.LastName}
		pl.Subject = fmt.Sprintf("Seller %s Mengundang Anda!", dataMerchant.MerchantName)
		pl.Content = emailContent

		err := m.sendEmailMerchant(ctxReq, pl)
		if err != nil {
			output <- ResultUseCase{Error: errors.New(errMsgFailedSendEmail), HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: nil}
	})

	return output
}

func (m *MerchantUseCaseImpl) SendEmailAdmin(ctxReq context.Context, data model.B2CMerchantDataV2, fullName string, reasonReject string, adminCMS string) <-chan ResultUseCase {
	ctx := "MerchantUseCase-SendEmailAdmin"
	output := make(chan ResultUseCase)
	go tracer.WithTrace(ctxReq, ctx, nil, func(ctxReq context.Context) {
		defer close(output)

		// get template email
		templateKeyAdminRejectMerchant := "EMAIL_ADMIN_REJECT_MERCHANT_UPGRADE"
		if reasonReject != "" {
			templateKeyAdminRejectMerchant = "EMAIL_ADMIN_REJECT_MERCHANT_UPGRADE_WITH_REASON"
		}
		templateEmailDetailAdmin, errAdminTemplate := m.GetTemplateEmail(ctxReq, templateKeyAdminRejectMerchant)
		if errAdminTemplate != nil {
			output <- ResultUseCase{Error: errAdminTemplate, HTTPStatus: http.StatusBadRequest}
			return
		}

		var upgradeStatus string
		if data.UpgradeStatus.String == model.PendingAssociateString {
			upgradeStatus = model.AssociateString
		} else if data.UpgradeStatus.String == model.PendingManageString {
			upgradeStatus = model.ManageString
		}

		emailContentRejectRequestAdmin := templateEmailDetailAdmin.Content
		emailContentRejectRequestAdmin = strings.Replace(emailContentRejectRequestAdmin, merchantPlaceholder, data.MerchantName, -1)
		emailContentRejectRequestAdmin = strings.Replace(emailContentRejectRequestAdmin, upgradeStatusPlaceholder, upgradeStatus, 1)
		emailContentRejectRequestAdmin = strings.Replace(emailContentRejectRequestAdmin, adminName, adminCMS, 1)
		if reasonReject != "" {
			emailContentRejectRequestAdmin = strings.Replace(emailContentRejectRequestAdmin, reasonRejectText, reasonReject, 1)
		}
		payloadRejectRequestForAdmin := serviceModel.Email{}
		payloadRejectRequestForAdmin.From = serviceModel.EmailCare
		payloadRejectRequestForAdmin.FromName = serviceModel.NoReplyName
		emailToAdmin, ok := os.LookupEnv("EMAIL_INTERNAL_TEAM")
		if !ok {
			errs := errors.New("erorr")
			output <- ResultUseCase{Error: errs, HTTPStatus: http.StatusBadRequest}
			return
		}

		payloadRejectRequestForAdmin.To = []string{emailToAdmin}
		payloadRejectRequestForAdmin.ToName = []string{"Admin Team"}
		payloadRejectRequestForAdmin.Subject = model.MerchantUpgradeRejectSubject
		payloadRejectRequestForAdmin.Content = emailContentRejectRequestAdmin

		if err := m.sendEmailAdmin(ctxReq, payloadRejectRequestForAdmin); err != nil {
			errMessage := errors.New(errMsgFailedSendEmail)
			output <- ResultUseCase{Error: errMessage, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: data}
	})

	return output
}
