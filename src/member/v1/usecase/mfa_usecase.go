package usecase

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/sha1"
	"database/sql"
	"encoding/base32"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	dgoogauth "github.com/dgryski/dgoogauth"
	qr "rsc.io/qr"
)

// GetMFASettings function for getting detail member based on member id
func (mu *MemberUseCaseImpl) GetMFASettings(ctxReq context.Context, uid string) <-chan ResultUseCase {
	ctx := "MemberUseCase-GetMFASettings"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags[helper.TextArgs] = uid
		if !strings.Contains(uid, usrFormat) {
			err := fmt.Errorf(helper.ErrorParameterInvalid, msgErrorMemberID)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		memberResult := <-mu.MemberMFAQueryRead.FindMFASettings(ctxReq, uid)
		if memberResult.Error != nil {
			tracer.SetError(ctxReq, memberResult.Error)
			output <- ResultUseCase{Error: memberResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		result := memberResult.Result.(model.MFASettings)
		output <- ResultUseCase{Result: result}

	})

	return output
}

// GenerateMFASettings function for getting detail member based on member id
func (mu *MemberUseCaseImpl) GenerateMFASettings(ctxReq context.Context, userID, requestFrom string) <-chan ResultUseCase {
	ctx := "MemberUseCase-GenerateMFASettings"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		if !strings.Contains(userID, usrFormat) {
			err := fmt.Errorf(helper.ErrorParameterInvalid, msgErrorMemberID)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		memberResult := <-mu.MemberRepoRead.Load(ctxReq, userID)
		if memberResult.Error != nil {
			if memberResult.Error == sql.ErrNoRows {
				memberResult.Error = fmt.Errorf(helper.ErrorDataNotFound, labelMember)
			}

			output <- ResultUseCase{Error: memberResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		member, ok := memberResult.Result.(model.Member)
		if !ok {
			err := errors.New(msgErrorResultMember)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		generatedResult := <-mu.GenerateMFASharedKey(ctxReq, member, requestFrom)

		result, ok := generatedResult.Result.(model.MFAGenerateSettings)
		if !ok {
			err := errors.New(model.ErrorGenerateMFA)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: result}

	})

	return output
}

// GenerateMFASharedKey function for generate shared key
func (mu *MemberUseCaseImpl) GenerateMFASharedKey(ctxReq context.Context, member model.Member, requestFrom string) <-chan ResultUseCase {
	ctx := "MemberUseCase-GenerateMFASharedKey"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags[helper.TextArgs] = member
		account := member.Email
		var issuer string
		if requestFrom == "account" {
			issuer = "Bhinneka"
		} else {
			issuer = "Narwhal"
		}

		qrFilename := "/tmp/qr.png"

		// Generate random secret
		secret := make([]byte, 10)
		_, err := rand.Read(secret)
		if err != nil {
			err := errors.New(model.ErrorGenerateMFA)
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// encode secret
		secretBase32 := base32.StdEncoding.EncodeToString(secret)

		URL, err := url.Parse("otpauth://totp")
		if err != nil {
			err := errors.New(model.ErrorGenerateMFA)
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// set url path
		URL.Path += "/" + url.PathEscape(issuer) + ":" + url.PathEscape(account)

		params := url.Values{}
		params.Add(textSecret, secretBase32)
		params.Add("issuer", issuer)
		URL.RawQuery = params.Encode()

		code, err := qr.Encode(URL.String(), qr.Q)
		if err != nil {
			err := errors.New(model.ErrorGenerateMFA)
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		b := code.PNG()

		// write file as qr code image
		err = ioutil.WriteFile(qrFilename, b, 0600)
		if err != nil {
			err := errors.New(model.ErrorGenerateMFA)
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		imgFile, err := os.Open(qrFilename) // a QR code image
		if err != nil {
			err := errors.New(model.ErrorGenerateMFA)
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		defer imgFile.Close()

		// create a new buffer base on file size
		fInfo, _ := imgFile.Stat()
		var size int64 = fInfo.Size()
		buf := make([]byte, size)
		// read file content into buffer
		fReader := bufio.NewReader(imgFile)
		fReader.Read(buf)

		// convert the buffer bytes to base64 string - use buf.Bytes() for new image
		imgBase64Str := base64.StdEncoding.EncodeToString(buf)

		resultData := model.MFAGenerateSettings{}
		resultData.SharedKeyQRCode = "data:image/png;base64," + imgBase64Str
		resultData.SharedKeyText = secretBase32

		output <- ResultUseCase{Result: resultData}
	})
	return output
}

// ActivateMFASettings function for activate mfa otp
func (mu *MemberUseCaseImpl) ActivateMFASettings(ctxReq context.Context, activateData model.MFAActivateSettings) <-chan ResultUseCase {
	ctx := "MemberUseCase-ActivateMFASettings"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags[helper.TextArgs] = activateData.MemberID
		memberResult := <-mu.MemberRepoRead.Load(ctxReq, activateData.MemberID)
		if memberResult.Error != nil {
			if memberResult.Error == sql.ErrNoRows {
				memberResult.Error = fmt.Errorf(helper.ErrorDataNotFound, labelMember)
			}

			output <- ResultUseCase{Error: memberResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		errorValidateMFA := mu.validateMFA(activateData)
		if errorValidateMFA != nil {
			tracer.SetError(ctxReq, errorValidateMFA)
			output <- ResultUseCase{Error: errorValidateMFA, HTTPStatus: http.StatusBadRequest}
			return
		}

		// encode shared key
		mfaKey := base64.URLEncoding.EncodeToString([]byte(activateData.SharedKeyText))
		// update field mfaEnabled
		if err := mu.activateMFA(ctxReq, activateData.RequestFrom, activateData.MemberID, mfaKey); err != nil {
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: activateData}
	})

	return output
}
func (mu *MemberUseCaseImpl) ActivateMFASettingV3(ctxReq context.Context, activateDatas model.MFAActivateSettings) <-chan ResultUseCase {
	ctx := "MemberUseCase-ActivateMFASettingV3"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags[helper.TextArgs] = activateDatas.MemberID
		membersResult := <-mu.MemberRepoRead.Load(ctxReq, activateDatas.MemberID)
		if membersResult.Error != nil {
			if membersResult.Error == sql.ErrNoRows {
				membersResult.Error = fmt.Errorf(helper.ErrorDataNotFound, labelMember)
			}

			output <- ResultUseCase{Error: membersResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}
		members, ok := membersResult.Result.(model.Member)
		if !ok {
			err := errors.New(msgErrorResultMember)
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		passwordSalt := strings.TrimSpace(members.Salt)
		passwordHasher := model.NewPBKDF2Hasher(model.SaltSize, model.SaltSize, model.IterationsCount, sha1.New)
		passwordHasher.ParseSalt(passwordSalt)
		base64Data := base64.StdEncoding.EncodeToString(passwordHasher.Hash([]byte(activateDatas.Password)))
		if base64Data != members.Password {
			err := fmt.Errorf(model.ErrorMFAPassword)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		errorValidateMFA := mu.validateMFAV3(activateDatas, base64Data, members)
		if errorValidateMFA != nil {
			tracer.SetError(ctxReq, errorValidateMFA)
			output <- ResultUseCase{Error: errorValidateMFA, HTTPStatus: http.StatusBadRequest}
			return
		}

		// encode shared key
		mfaKeyV3 := base64.URLEncoding.EncodeToString([]byte(activateDatas.SharedKeyText))
		// update field mfaEnabled
		if err := mu.activateMFA(ctxReq, activateDatas.RequestFrom, activateDatas.MemberID, mfaKeyV3); err != nil {
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: activateDatas}
	})

	return output
}

func (mu *MemberUseCaseImpl) activateMFA(ctxReq context.Context, requestFrom, memberID, mfaKey string) error {
	if requestFrom == helper.TextAccount {
		updateResult := <-mu.MemberMFARepoWrite.MFAEnabled(ctxReq, memberID, mfaKey)
		if updateResult.Error != nil {
			return errors.New(model.ErrorEnabledMFA)
		}
	} else if requestFrom == helper.TextNarwhal {
		updateResult := <-mu.MemberMFARepoWrite.EnableNarwhalMFA(ctxReq, memberID, mfaKey)
		if updateResult.Error != nil {
			return errors.New(model.ErrorEnabledMFA)
		}
	}
	return nil
}

func (mu *MemberUseCaseImpl) validateMFA(activateData model.MFAActivateSettings) error {
	byPassStatic := false
	// only for exclude prod & key static
	if os.Getenv("ENV") != "PROD" && activateData.SharedKeyText == model.StaticSharedMfaKeyForDev && activateData.Otp == model.StaticOTPMfaForDev {
		byPassStatic = true
	}

	if !byPassStatic {
		// The OTPConfig gets modified by otpc.Authenticate() to prevent passcode replay, etc.,
		// so allocate it once and reuse it for multiple calls.
		otpc := &dgoogauth.OTPConfig{
			Secret:      activateData.SharedKeyText,
			WindowSize:  3,
			HotpCounter: 0,
			UTC:         true,
		}

		// authentication OTPConfig
		val, err := otpc.Authenticate(activateData.Otp)
		if err != nil {
			return err
		}

		if !val {
			err := errors.New(model.ErrorMFAOTP)
			return err
		}
	}
	return nil
}

func (mu *MemberUseCaseImpl) validateMFAV3(activateData model.MFAActivateSettings, base64Data string, member model.Member) error {
	byPassStatic := false
	// only for exclude prod & key static
	if os.Getenv("ENV") != "PROD" && activateData.SharedKeyText == model.StaticSharedMfaKeyForDev && activateData.Otp == model.StaticOTPMfaForDev && base64Data == member.Password {
		byPassStatic = true
	}

	if !byPassStatic {
		// The OTPConfig gets modified by otpc.Authenticate() to prevent passcode replay, etc.,
		// so allocate it once and reuse it for multiple calls.
		otpc := &dgoogauth.OTPConfig{
			Secret:      activateData.SharedKeyText,
			WindowSize:  3,
			HotpCounter: 0,
			UTC:         true,
		}

		// authentication OTPConfig
		val, err := otpc.Authenticate(activateData.Otp)
		if err != nil {
			return err
		}

		if !val {
			err := errors.New(model.ErrorMFAOTP)
			return err
		}
	}
	return nil
}

// DisabledMFASetting function for getting detail member based on member id
func (mu *MemberUseCaseImpl) DisabledMFASetting(ctxReq context.Context, userID, requestFrom string) <-chan ResultUseCase {
	ctx := "MemberUseCase-DisabledMFASetting"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags[helper.TextMemberIDCamel] = userID
		if !strings.Contains(userID, usrFormat) {
			err := fmt.Errorf(helper.ErrorParameterInvalid, msgErrorMemberID)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		memberResult := <-mu.MemberRepoRead.Load(ctxReq, userID)
		if memberResult.Error != nil {
			if memberResult.Error == sql.ErrNoRows {
				memberResult.Error = fmt.Errorf(helper.ErrorDataNotFound, labelMember)
			}

			output <- ResultUseCase{Error: memberResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		member, ok := memberResult.Result.(model.Member)
		if !ok {
			err := errors.New(msgErrorResultMember)
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		// update field mfaEnabled to false
		if err := mu.disabledMFA(ctxReq, requestFrom, member.ID); err != nil {
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: member}

	})

	return output
}

func (mu *MemberUseCaseImpl) disabledMFA(ctxReq context.Context, requestFrom string, memberID string) error {
	if requestFrom == helper.TextAccount {
		updateResult := <-mu.MemberMFARepoWrite.MFADisabled(ctxReq, memberID)
		if updateResult.Error != nil {
			return errors.New(model.ErrorDisabledMFA)
		}
	} else if requestFrom == helper.TextNarwhal {
		adminResult := <-mu.MemberMFARepoWrite.DisableNarwhalMFA(ctxReq, memberID)
		if adminResult.Error != nil {
			return errors.New(model.ErrorDisabledMFA)
		}
	}
	return nil
}
