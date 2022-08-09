package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
)

// UpdateProfilePicture function for update profile picture
func (mu *MemberUseCaseImpl) UpdateProfilePicture(ctxReq context.Context, data model.ProfilePicture) <-chan ResultUseCase {
	ctx := "MemberUseCase-UpdateProfilePicture"
	var params serviceModel.SendbirdRequestV4

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags[helper.TextParameter] = data
		tags[helper.TextMemberIDCamel] = data.ID
		if !strings.Contains(data.ID, usrFormat) {
			err := fmt.Errorf(helper.ErrorParameterInvalid, msgErrorMemberID)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		if !helper.ValidateDocumentFileURL(data.ProfilePicture) {
			err := errors.New("File not valid")
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		memberResult := <-mu.MemberQueryRead.FindByID(ctxReq, data.ID)
		if memberResult.Error != nil {
			if memberResult.Error == sql.ErrNoRows {
				memberResult.Error = fmt.Errorf(helper.ErrorDataNotFound, labelMember)
				output <- ResultUseCase{Error: memberResult.Error, HTTPStatus: http.StatusBadRequest}
				return
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

		saveResult := <-mu.MemberRepoWrite.UpdateProfilePicture(ctxReq, data)
		if saveResult.Error != nil {
			err := errors.New("failed to update profile picture")
			tracer.SetError(ctxReq, saveResult.Error)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusInternalServerError}
			return
		}

		params.UserID = data.ID
		params.NickName = member.FirstName + " " + member.LastName
		params.ProfileURL = data.ProfilePicture
		mu.UpdateUserSendbirdV4(ctxReq, &params)

		output <- ResultUseCase{Result: data}

	})

	return output
}

// GetProfileComplete function for getting completeness profile information
func (mu *MemberUseCaseImpl) GetProfileComplete(ctxReq context.Context, uid string) <-chan ResultUseCase {
	ctx := "MemberUseCase-GetProfileComplete"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags[helper.TextMemberIDCamel] = uid

		memberResult := <-mu.MemberRepoRead.Load(ctxReq, uid)
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
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		profileComplete := model.ProfileComplete{}
		profileField := []model.ProfileField{}

		fields := map[string]string{"email": "Email", "mobilePhone": "Nomor Ponsel", "shippingAddress": "Alamat Pengiriman", "profile": "Data Diri"}
		if member.Password == "" {
			fields[textPassword] = labelPassword
		}
		totalCompleteness := 100
		percentageField := 100 / len(fields)

		for key, label := range fields {
			data := model.ProfileField{}
			data.Value = true
			data.Key = key
			data.Label = label
			data = mu.GetDataProfileCompleteness(ctxReq, data, member)

			if !data.Value {
				totalCompleteness = totalCompleteness - percentageField
			}
			profileField = append(profileField, data)
		}
		sort.Slice(profileField, func(i, j int) bool { return profileField[i].Step < profileField[j].Step })
		profileComplete.Field = profileField
		profileComplete.Percentage = strconv.Itoa(totalCompleteness)

		output <- ResultUseCase{Result: profileComplete}

	})

	return output
}

// GetDataProfileCompleteness function for getting value
func (mu *MemberUseCaseImpl) GetDataProfileCompleteness(ctxReq context.Context, data model.ProfileField, member model.Member) model.ProfileField {
	switch data.Key {
	case "email":
		data.Step = 1
		if member.Email == "" {
			data.Value = false
		}
	case "mobilePhone":
		data.Step = 2
		if member.Mobile == "" {
			data.Value = false
		}
	case "shippingAddress":
		data.Step = 3
		// count shipping address
		countAddress := <-mu.ShippingAddressRepo.CountShippingAddressByUserID(ctxReq, member.ID)
		if countAddress.Error != nil {
			data.Value = false
		}

		total, _ := countAddress.Result.(int)
		if total < 1 {
			data.Value = false
		}
	case "profile":
		data.Step = 4
		data.Value = validateProfile(data, member)
	case textPassword:
		data.Step = 5
		if member.Password == "" {
			data.Value = false
		}
	}
	return data
}

// UpdateProfileName function for update profile name
func (mu *MemberUseCaseImpl) UpdateProfileName(ctxReq context.Context, data model.ProfileName) <-chan ResultUseCase {
	ctx := "MemberUseCase-UpdateProfileName"
	var params serviceModel.SendbirdRequestV4
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags[helper.TextArgs] = data
		member, err := mu.validateProfileName(ctxReq, data)
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		names := strings.Split(data.ProfileName, " ")
		firstName := names[0]
		lastName := helper.SetLastName(names)
		member.FirstName = firstName
		member.LastName = lastName

		saveResult := <-mu.MemberRepoWrite.Save(ctxReq, member)
		if saveResult.Error != nil {
			err := errors.New("failed to update profile name")
			tracer.SetError(ctxReq, err)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusInternalServerError}
			return
		}

		params.UserID = member.ID
		params.ProfileURL = member.ProfilePicture
		params.NickName = member.FirstName + " " + member.LastName
		params.Token = member.Token
		mu.UpdateUserSendbirdV4(ctxReq, &params)

		output <- ResultUseCase{Result: data}
	})

	return output
}

func (mu *MemberUseCaseImpl) validateProfileName(ctxReq context.Context, data model.ProfileName) (model.Member, error) {
	result := model.Member{}
	memberResult := <-mu.MemberQueryRead.FindByID(ctxReq, data.ID)
	if memberResult.Error != nil {
		if memberResult.Error == sql.ErrNoRows {
			memberResult.Error = fmt.Errorf(helper.ErrorDataNotFound, labelMember)
			return result, memberResult.Error
		}
		return result, memberResult.Error
	}

	if data.ProfileName == "" {
		err := fmt.Errorf("name is required")
		return result, err
	}

	if len(data.ProfileName) > 50 {
		err := fmt.Errorf("name max 50 character")
		return result, err
	}

	if !golib.ValidateAlphabetWithSpace(data.ProfileName) {
		err := fmt.Errorf(helper.ErrorParameterInvalid, "name")
		return result, err
	}

	result = memberResult.Result.(model.Member)
	return result, nil
}

// getProfilePicture function for get url image
func (mu *MemberUseCaseImpl) getProfilePicture(ctxReq context.Context, url string) string {
	ctx := "MemberUseCase-getProfilePicture"
	isAttachment := "false"
	documentURL := <-mu.UploadService.GetURLImage(ctxReq, url, isAttachment)
	if documentURL.Result != nil {
		documentURLResult, ok := documentURL.Result.(serviceModel.ResponseUploadService)
		if !ok {
			err := errors.New("failed get url image")
			helper.SendErrorLog(ctxReq, ctx, "parse_picture", err, url)
			return url
		}
		return documentURLResult.Data.URL
	}
	return url
}

//validateProfile
func validateProfile(data model.ProfileField, member model.Member) bool {
	data.Value = true
	if member.GenderString == "" {
		data.Value = false
	}
	if member.BirthDateString == "" {
		data.Value = false
	}
	if !golib.ValidateAlphabetWithSpace(member.FirstName) {
		data.Value = false
	}
	if len(member.LastName) != 0 {
		if !golib.ValidateAlphabetWithSpace(member.LastName) {
			data.Value = false
		}
	}
	return data.Value
}
