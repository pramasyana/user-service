package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Bhinneka/golib"
	goString "github.com/Bhinneka/golib/string"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/eapache/go-resiliency/retrier"
)

// ResendActivation usecase function for resend email activation
func (mu *MemberUseCaseImpl) ResendActivation(ctxReq context.Context, email string) <-chan ResultUseCase {
	ctx := "MemberUseCase-ResendActivation"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags[helper.TextEmail] = email
		if email == "" {
			err := fmt.Errorf("email required")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		if err := goString.ValidateEmail(email); err != nil {
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		emailResult := <-mu.MemberQueryRead.FindByEmail(ctxReq, email)
		if emailResult.Result == nil {
			err := fmt.Errorf("Akun anda tidak terdaftar")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		member := emailResult.Result.(model.Member)
		if member.Status == model.Active {
			err := fmt.Errorf("Alamat email Anda sudah aktif")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// check attempt and save it
		attemptResult := <-mu.checkAttemptResendActivation(ctxReq, email, member.LastTokenAttempt)
		if attemptResult.Error != nil {
			tags[helper.TextResponse] = attemptResult.Error
			output <- ResultUseCase{Error: attemptResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}
		mix := member.Email + "-" + member.FirstName
		member.Token = helper.GenerateTokenByString(mix)

		saveResult := <-mu.MemberRepoWrite.Save(ctxReq, member)
		if saveResult.Error != nil {
			err := errors.New(msgErrorSaveMember)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		data := model.SuccessResponse{}
		data.FirstName = member.FirstName
		data.LastName = member.LastName
		data.Token = member.Token
		data.Email = member.Email

		plEmail := model.MemberEmailQueue{
			Member: &member,
			Data:   data,
		}
		go mu.QPublisher.QueueJob(ctxReq, plEmail, data.ID, "SendEmailRegisterMember")

		updateLastAttempt := <-mu.MemberQueryWrite.UpdateLastTokenAttempt(ctxReq, email)
		if updateLastAttempt.Error != nil {
			err := fmt.Errorf(model.ErrorResendActivation)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		res := model.PlainSuccessResponse{Email: member.Email}
		output <- ResultUseCase{Result: res}

	})

	return output
}

// SendEmailRegisterMember usecase function for send email member registration
func (mu *MemberUseCaseImpl) SendEmailRegisterMember(ctxReq context.Context, data model.SuccessResponse, registrationType string) <-chan ResultUseCase {
	ctx := "MemberUseCase-SendEmailRegisterMember"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		// get template email

		templateEmailDetail, errTemplate := mu.GetTemplateEmail(ctxReq, "EMAIL_PERSONAL_REGISTRATION_TEMPLATE_ID")
		if errTemplate != nil {
			output <- ResultUseCase{Error: errTemplate, HTTPStatus: http.StatusBadRequest}
			return
		}

		memberName := data.FirstName + " " + data.LastName

		memberType := "personal"
		if registrationType == "merchant" {
			memberType = registrationType
		}

		// set content to our email
		content := templateEmailDetail.Content
		year := time.Now().Format("2006")
		strURL := fmt.Sprintf("%s/accountactivation/?token=%s&type=%s", mu.SturgeonCFUrl, data.Token, memberType)

		find := findEmail
		replacer := []string{year, memberName, strURL}
		content = golib.StringArrayReplace(content, find, replacer)

		pl := serviceModel.Email{}
		pl.From = serviceModel.NoReply
		pl.FromName = serviceModel.NoReplyName
		pl.To = []string{data.Email}
		pl.ToName = []string{memberName}
		pl.Subject = model.SubjectConfirmEmail
		pl.Content = content

		tags[helper.TextParameter] = pl
		if err := mu.sendEmailMember(ctxReq, pl); err != nil {
			helper.SendErrorLog(ctxReq, ctx, scopeSendEmail, err, data.Email)
			output <- ResultUseCase{Error: errors.New(msgErrorSendEmail), HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: data}
	})

	return output
}

// SendEmailWelcomeMember usecase function for send email success registration
func (mu *MemberUseCaseImpl) SendEmailWelcomeMember(ctxReq context.Context, data model.SuccessResponse) <-chan ResultUseCase {
	ctx := "MemberUseCase-SendEmailWelcomeMember"
	output := make(chan ResultUseCase)
	mu.AccessTokenGenerator.GenerateAnonymous(ctxReq)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		// get template email

		templateEmailDetail, errTemplate := mu.GetTemplateEmail(ctxReq, "EMAIL_SUCCESS_REGISTRATION_TEMPLATE_ID")
		if errTemplate != nil {
			output <- ResultUseCase{Error: errTemplate, HTTPStatus: http.StatusBadRequest}
			return
		}

		memberName := data.FirstName + " " + data.LastName

		// set content to our email
		content := templateEmailDetail.Content
		year := time.Now().Format("2006")
		strURL := mu.B2cCFUrl

		find := findEmail
		replacer := []string{year, memberName, strURL}
		content = golib.StringArrayReplace(content, find, replacer)

		pl := serviceModel.Email{}
		pl.From = serviceModel.NoReply
		pl.FromName = serviceModel.NoReplyName
		pl.To = []string{data.Email}
		pl.ToName = []string{memberName}
		pl.Subject = model.SubjectWelcomeEmail
		pl.Content = content

		tags[helper.TextParameter] = pl
		if err := mu.sendEmailMember(ctxReq, pl); err != nil {
			helper.SendErrorLog(ctxReq, ctx, scopeSendEmail, err, data.Email)
			output <- ResultUseCase{Error: errors.New(msgErrorSendEmail), HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: data}
	})

	return output
}

// SendEmailForgotPassword usecase function for send email forgot password
func (mu *MemberUseCaseImpl) SendEmailForgotPassword(ctxReq context.Context, data model.SuccessResponse) <-chan ResultUseCase {
	ctx := "MemberUseCase-SendEmailForgotPassword"
	output := make(chan ResultUseCase)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		// get template email

		templateEmailDetail, errTemplate := mu.GetTemplateEmail(ctxReq, "EMAIL_FORGOT_PASSWORD_TEMPLATE_ID")
		if errTemplate != nil {
			output <- ResultUseCase{Error: errTemplate, HTTPStatus: http.StatusBadRequest}
			return
		}

		memberName := data.FirstName + " " + data.LastName

		// set content to our email
		content := templateEmailDetail.Content
		year := time.Now().Format("2006")
		strURL := fmt.Sprintf("%s/resetpassword/%s", mu.SturgeonCFUrl, data.Token)

		find := findEmail
		replacer := []string{year, memberName, strURL}
		content = golib.StringArrayReplace(content, find, replacer)

		pl := serviceModel.Email{}
		pl.From = serviceModel.NoReply
		pl.FromName = serviceModel.NoReplyName
		pl.To = []string{data.Email}
		pl.ToName = []string{memberName}
		pl.Subject = model.SubjectForgotPassword
		pl.Content = content

		tags[helper.TextParameter] = pl
		if err := mu.sendEmailMember(ctxReq, pl); err != nil {
			helper.SendErrorLog(ctxReq, ctx, scopeSendEmail, err, data.Email)
			output <- ResultUseCase{Error: errors.New(msgErrorSendEmail), HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: data}
	})

	return output
}

// SendEmailSuccessForgotPassword usecase function for send email forgot password
func (mu *MemberUseCaseImpl) SendEmailSuccessForgotPassword(ctxReq context.Context, data model.SuccessResponse) <-chan ResultUseCase {
	ctx := "MemberUseCase-SendEmailSuccessForgotPassword"
	output := make(chan ResultUseCase)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		// get template email
		templateEmailDetail, errTemplate := mu.GetTemplateEmail(ctxReq, "EMAIL_SUCCESS_FORGOT_PASSWORD_TEMPLATE_ID")
		if errTemplate != nil {
			output <- ResultUseCase{Error: errTemplate, HTTPStatus: http.StatusBadRequest}
			return
		}

		memberName := data.FirstName + " " + data.LastName
		// set content to our email
		content := templateEmailDetail.Content
		year := time.Now().Format("2006")
		find := []string{"##YEAR##", "##EMAIL##", "##FULLNAME##"}
		replacer := []string{year, data.Email, memberName}
		content = golib.StringArrayReplace(content, find, replacer)

		pl := serviceModel.Email{}
		pl.From = serviceModel.NoReply
		pl.FromName = serviceModel.NoReplyName
		pl.To = []string{data.Email}
		pl.ToName = []string{memberName}
		pl.Subject = model.SubjectConfirmPassword
		pl.Content = content

		tags[helper.TextParameter] = pl
		if err := mu.sendEmailMember(ctxReq, pl); err != nil {
			helper.SendErrorLog(ctxReq, ctx, scopeSendEmail, err, data.Email)
			output <- ResultUseCase{Error: errors.New(msgErrorSendEmail), HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: data}
	})

	return output
}

// ValidateEmailDomain function for validating email domain
func (mu *MemberUseCaseImpl) ValidateEmailDomain(ctxReq context.Context, email string) <-chan ResultUseCase {
	ctx := "MemberUseCase-ValidateEmailDomain"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)
		if os.Getenv("ENABLE_VALIDATE_EMAIL_DOMAIN") == "false" {
			output <- ResultUseCase{Result: true}
			return
		}

		ValidateEmailDomain := golib.IsDisabledEmail(email)
		tags[helper.TextParameter] = email
		tags[helper.TextResponse] = ValidateEmailDomain
		if ValidateEmailDomain {
			errMessage := errors.New(model.ErrorInvalidEmailBahasa)
			output <- ResultUseCase{Error: errMessage, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: ValidateEmailDomain}

	})

	return output
}

// SendEmailAddMember usecase function for send email when add member from CMS
func (mu *MemberUseCaseImpl) SendEmailAddMember(ctxReq context.Context, data model.SuccessResponse) <-chan ResultUseCase {
	ctx := "MemberUseCase-SendEmailAddMember"
	output := make(chan ResultUseCase)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		templateDetail, errTemplate := mu.GetTemplateEmail(ctxReq, "EMAIL_ADD_MEMBER_TEMPLATE_ID")
		if errTemplate != nil {
			output <- ResultUseCase{Error: errTemplate, HTTPStatus: http.StatusBadRequest}
			return
		}

		memberFullName := data.FirstName + " " + data.LastName

		// set content to our email
		content := templateDetail.Content
		year := time.Now().Format("2006")
		strURL := fmt.Sprintf("%s/activate-member?token=%s", mu.SturgeonCFUrl, data.Token)

		content = golib.StringArrayReplace(content, findEmail, []string{year, memberFullName, strURL})

		payload := serviceModel.Email{}
		payload.From = serviceModel.NoReply
		payload.FromName = serviceModel.NoReplyName
		payload.To = []string{data.Email}
		payload.ToName = []string{memberFullName}
		payload.Subject = model.SubjectAddMember
		payload.Content = content

		tags[helper.TextParameter] = payload
		if err := mu.sendEmailMember(ctxReq, payload); err != nil {
			helper.SendErrorLog(ctxReq, ctx, scopeSendEmail, err, data.Email)
			output <- ResultUseCase{Error: errors.New(msgErrorSendEmail), HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: data}
	})

	return output
}

// GetTemplateEmail usecase function for get template email
func (mu *MemberUseCaseImpl) GetTemplateEmail(ctxReq context.Context, templateKey string) (template serviceModel.Template, err error) {
	ctx := "MemberUseCase-GetTemplateEmail"
	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		"template.id": templateKey,
	}
	defer tr.Finish(tags)

	templateEmailDetail := serviceModel.Template{}
	emailMerchantRegistrationTemplateID := golib.GetEnvOrFail("GetTemplateEmail", "get_template_email", templateKey)

	templateEmailStatics := <-mu.NotificationService.GetTemplateByID(ctxReq, emailMerchantRegistrationTemplateID, templateKey)
	if templateEmailStatics.Error != nil {
		return templateEmailDetail, templateEmailStatics.Error
	}

	// detail
	templateEmailDetail, ok := templateEmailStatics.Result.(serviceModel.Template)
	if !ok {
		return templateEmailDetail, errors.New(helper.ErrMsgFailedGetTemplate)
	}

	// check template email not null
	if templateEmailDetail.Content == "" {
		return templateEmailDetail, errors.New(helper.ErrMsgEmptyContent)
	}
	tags["content"] = templateEmailDetail.Content

	return templateEmailDetail, nil
}

// / sendEmailMember usecase function for send email member
func (mu *MemberUseCaseImpl) sendEmailMember(ctxReq context.Context, pl serviceModel.Email) error {
	r := retrier.New(retrier.ConstantBackoff(3, 500*time.Millisecond), nil)
	err := r.Run(func() error {
		_, err := mu.NotificationService.SendEmail(ctxReq, pl)

		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		helper.SendErrorLog(ctxReq, "sendEmailMember", "send_email_member", err, pl)
		errMessage := errors.New(helper.ErrMsgFailedSendEmail)
		return errMessage
	}
	return nil
}
