package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Bhinneka/golib"
	goString "github.com/Bhinneka/golib/string"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	contactModel "github.com/Bhinneka/user-service/src/client/v1/model"
	corporateModel "github.com/Bhinneka/user-service/src/corporate/v2/model"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
)

// PublishToKafkaUser function for publish to kafka user
func (mu *MemberUseCaseImpl) PublishToKafkaUser(ctxReq context.Context, data *model.Member, eventType string) error {
	ctx := "MemberUseCase-PublishToKafkaUser"

	trace := tracer.StartTrace(ctxReq, ctx)
	tags := make(map[string]interface{})

	data.GenderString = data.Gender.String()
	data.Gender = model.StringToGender(data.GenderString)
	member := serviceModel.MemberDolphin{
		ID:              data.ID,
		Email:           strings.ToLower(data.Email),
		FirstName:       data.FirstName,
		LastName:        data.LastName,
		Gender:          data.Gender.GetDolpinGender(),
		DOB:             data.BirthDateString,
		Phone:           data.Phone,
		Ext:             data.Ext,
		Mobile:          data.Mobile,
		Street1:         data.Address.Street1,
		Street2:         data.Address.Street2,
		PostalCode:      data.Address.ZipCode,
		SubDistrictID:   data.Address.SubDistrictID,
		SubDistrictName: data.Address.SubDistrict,
		DistrictID:      data.Address.DistrictID,
		DistrictName:    data.Address.District,
		CityID:          data.Address.CityID,
		CityName:        data.Address.City,
		ProvinceID:      data.Address.ProvinceID,
		ProvinceName:    data.Address.Province,
		Status:          strings.ToUpper(data.StatusString),
		Created:         data.Created.Format(time.RFC3339),
		LastModified:    time.Now().Format(time.RFC3339),
		FacebookID:      data.SocialMedia.FacebookID,
		GoogleID:        data.SocialMedia.GoogleID,
		LDAPID:          data.SocialMedia.LDAPID,
		AppleID:         data.SocialMedia.AppleID,
	}

	payload := serviceModel.DolphinPayloadNSQ{
		EventOrchestration:     "UpdateMember",
		TimestampOrchestration: time.Now().Format(time.RFC3339),
		EventType:              eventType,
		Counter:                0,
		Payload:                member,
	}

	// prepare to send to nsq
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "publish_payload", err, payload)
		return err
	}

	messageKey := member.ID
	tags[helper.TextArgs] = payload
	trace.Finish(tags)

	// excluded from parent context for kafka producer
	return mu.QPublisher.PublishKafka(trace.NewChildContext(), mu.Topic, messageKey, payloadJSON)
}

// InsertLogMember function to write log activity service for merchant warehouse
func (mu *MemberUseCaseImpl) InsertLogMember(ctxReq context.Context, oldData, newData *model.Member, action string) error {
	targetID := newData.ID
	if targetID == "" {
		targetID = oldData.ID
	}

	if newData.CreatedBy == "" {
		newData.CreatedBy = newData.ID
	}

	payload := serviceModel.Payload{
		Module:    model.Module,
		Action:    action,
		Target:    targetID,
		CreatorID: newData.CreatedBy,
		EditorID:  newData.ModifiedBy,
	}

	// insert log to activity service
	tokenResult := <-mu.AccessTokenGenerator.GenerateAnonymous(ctxReq)
	newCtx := context.WithValue(ctxReq, helper.TextAuthorization, tokenResult.AccessToken.AccessToken)
	mu.ActivityService.InsertLog(newCtx, oldData, newData, payload)
	return nil
}

// adjustMemberData function for adjusting member data
// all about date, status and gender
func (mu *MemberUseCaseImpl) adjustMemberData(ctxReq context.Context, member model.Member) model.Member {
	if member.Gender != 0 {
		member.GenderString = member.Gender.String()
	}
	member.StatusString = member.Status.String()
	if len(member.Password) > 0 {
		member.HasPassword = true

	}

	member.BirthDateString = member.BirthDate.Format(helper.FormatDOB)
	if member.BirthDateString == helper.DefaultDOB {
		member.BirthDateString = ""
	}

	member.LastLoginString = member.LastLogin.Format(time.RFC3339)
	if member.LastLoginString == model.DefaultDateTime {
		member.LastLoginString = ""
	}
	member.LastModifiedString = member.LastModified.Format(time.RFC3339)
	if member.LastModifiedString == model.DefaultDateTime {
		member.LastModifiedString = ""
	}
	member.CreatedString = member.Created.Format(time.RFC3339)
	if member.ProfilePicture != "" {
		member.ProfilePicture = mu.getProfilePicture(ctxReq, member.ProfilePicture)
	}
	member.Password = ""
	member.Salt = ""

	return member
}

func (mu *MemberUseCaseImpl) validateMemberDataUpdate(data *model.Member) error {
	var (
		err error
		ok  bool
	)
	if !strings.Contains(data.ID, usrFormat) {
		err = fmt.Errorf(helper.ErrorParameterInvalid, "user id")
		return err
	}

	// validate status when exists
	if len(data.StatusString) > 0 {
		data.Status, ok = model.ValidateStatus(data.StatusString)
		if !ok {
			err = fmt.Errorf(helper.ErrorParameterInvalid, scopeStatus)
			return err
		}
	}

	// validate new password when exists
	if len(data.NewPassword) > 0 {
		if data.NewPassword != data.RePassword {
			err := errors.New(msgErrorPasswordMatch)
			return err
		}

		// validate password format
		if err := helper.ValidatePassword(data.NewPassword); err != nil {
			return err
		}
	}

	if err = mu.checkMemberAuthorization(data); err != nil {
		return err
	}

	return nil
}

func (mu *MemberUseCaseImpl) checkMemberAuthorization(data *model.Member) (err error) {
	// optional staff and admin
	if len(data.IsStaffString) > 0 {
		data.IsStaff, err = strconv.ParseBool(data.IsStaffString)
		if err != nil {
			err = fmt.Errorf(helper.ErrorParameterInvalid, "is staff")
			return err
		}
	}
	if len(data.IsAdminString) > 0 {
		data.IsAdmin, err = strconv.ParseBool(data.IsAdminString)
		if err != nil {
			err = fmt.Errorf(helper.ErrorParameterInvalid, "is admin")
			return err
		}
	}
	return nil
}

func (mu *MemberUseCaseImpl) validateMemberDataRegister(ctxReq context.Context, data *model.Member) error {
	var err error
	if len(data.NewPassword) == 0 {
		err = fmt.Errorf(helper.ErrorParameterRequired, textPassword)
		return err
	}

	// validate password format
	if err := helper.ValidatePassword(data.NewPassword); err != nil {
		return err
	}

	if data.NewPassword != data.RePassword {
		err := errors.New(msgErrorPasswordMatch)
		return err
	}

	// validate email domain
	if err := <-mu.ValidateEmailDomain(ctxReq, data.Email); err.Error != nil {
		return err.Error
	}
	return nil
}

func (mu *MemberUseCaseImpl) validateMemberDataAddNew(data *model.Member) error {
	ctx := "MemberUseCase-validateMemberDataAddNew"
	var err error
	// validate email domain
	if err := <-mu.ValidateEmailDomain(context.Background(), data.Email); err.Error != nil {
		return err.Error
	}

	// validate address value existence
	if err = mu.validateUpdateParam(ctx, data); err != nil {
		return err
	}

	if len(data.Address.Street1) > 255 {
		err := errors.New("street 1 length cannot greater than 255 characters")
		return err
	}

	if len(data.Address.Street2) > 95 {
		err := errors.New("street 2 length cannot greater than 95 characters")
		return err
	}

	data.Address.Address = fmt.Sprintf("%s\n%s", data.Address.Street1, data.Address.Street2)

	if len(data.Password) > 0 {
		if data.Password != data.RePassword {
			err := errors.New(msgErrorPasswordMatch)
			return err
		}

		// validate password format
		if err := helper.ValidatePassword(data.Password); err != nil {
			return err
		}
	}

	if !helper.StringInSlice(data.SignUpFrom, model.SignUpFrom) {
		err = fmt.Errorf(helper.ErrorParameterInvalid, "sign up from")
		return err
	}
	return nil
}

func (mu *MemberUseCaseImpl) validateMemberDataImport(ctxReq context.Context, data *model.Member) error {
	tr := tracer.StartTrace(ctxReq, "MemberUseCase-validateMemberDataImport")
	defer tr.Finish(map[string]interface{}{"payload": data})

	if err := <-mu.ValidateEmailDomain(tr.NewChildContext(), data.Email); err.Error != nil {
		return err.Error
	}
	if err := mu.validateMemberImportContent(tr.NewChildContext(), data); err != nil {
		return err
	}
	return nil
}

func (mu *MemberUseCaseImpl) validateMemberImportContent(ctxReq context.Context, data *model.Member) error {
	tr := tracer.StartTrace(ctxReq, "MemberUseCase-validateMemberImportContent")
	defer tr.Finish(map[string]interface{}{"payload": data})

	var (
		ok bool
	)
	// validate address value existence
	if !golib.ValidateLatinOnly(data.Address.Street1) {
		return errors.New(msgErrorStreetName)
	}

	data.Address.Street1 = helper.ClearHTML(data.Address.Street1)

	data.Address.Address = data.Address.Street1

	// validate status when exists
	if len(data.StatusString) > 0 {
		data.Status, ok = model.ValidateStatus(strings.ToUpper(data.StatusString))
		if !ok {
			return fmt.Errorf(helper.ErrorParameterInvalid, scopeStatus)
		}
	}
	return nil
}

func (mu *MemberUseCaseImpl) validateMemberName(data *model.Member) (err error) {
	data.FirstName = strings.Trim(data.FirstName, " ")
	if len(data.FirstName) <= 0 {
		err = fmt.Errorf(helper.ErrorParameterRequired, textFirstName)
		return err
	}

	if !golib.ValidateAlphabetWithSpace(data.FirstName) {
		return fmt.Errorf("entry with first name %s is invalid", data.FirstName)
	}

	err = golib.ValidateMaxInput(data.FirstName, 25)
	if err != nil {
		return fmt.Errorf(helper.ErrorParameterLength, fmt.Sprintf("%s %s", textFirstName, data.FirstName), 25)
	}

	if len(data.LastName) != 0 {
		if !golib.ValidateAlphabetWithSpace(data.LastName) {
			return fmt.Errorf("entry with last name %s is invalid", data.LastName)
		}

		err = golib.ValidateMaxInput(data.LastName, 25)
		if err != nil {
			return fmt.Errorf(helper.ErrorParameterLength, fmt.Sprintf("%s %s", textLastName, data.LastName), 25)
		}
	}
	return nil
}

// validateMemberData function for validating member data
func (mu *MemberUseCaseImpl) validateMemberData(ctxReq context.Context, data *model.Member) error {
	tr := tracer.StartTrace(ctxReq, "MemberUseCase-validateMemberData")

	defer func() {
		tr.Finish(map[string]interface{}{"member": cleanTags(ctxReq, data)})
	}()
	ctx := "MemberUseCase-validateMemberData"

	var (
		err error
		ok  bool
	)
	switch data.Type {
	case textUpdate:
		err = mu.validateMemberDataUpdate(data)
	case textRegister:
		err = mu.validateMemberDataRegister(ctxReq, data)
	case textAdd:
		err = mu.validateMemberDataAddNew(data)
	case textImport:
		err = mu.validateMemberDataImport(tr.NewChildContext(), data)
	}

	if err != nil {
		return err
	}
	if err = mu.validateMemberName(data); err != nil {
		return err
	}
	if err = mu.validateBirthdate(data); err != nil {
		return err
	}

	if data.GenderString != "" {
		data.Gender, ok = model.ValidateGender(data.GenderString)
		if !ok {
			err = fmt.Errorf(helper.ErrorParameterInvalid, textGender)
			return err
		}

		data.GenderString = data.Gender.String()
	}

	if err = mu.validatePhone(ctx, data); err != nil {
		return err
	}

	// clear data from html tag
	data.FirstName = helper.ClearHTML(data.FirstName)
	data.LastName = helper.ClearHTML(data.LastName)
	return nil
}

func (mu *MemberUseCaseImpl) validateBirthdate(data *model.Member) (err error) {
	if os.Getenv("ENABLE_VALIDATE_DOB") == "true" && len(data.BirthDateString) <= 0 {
		err = fmt.Errorf(helper.ErrorParameterRequired, "birth date")
		return err
	}

	if len(data.BirthDateString) > 0 {
		// parse and replace dob value
		data.BirthDate, err = time.Parse(helper.FormatDOB, data.BirthDateString)
		if err != nil {
			err := fmt.Errorf(helper.ErrorParameterInvalid, msgErrorBirthdate)
			return err
		}

		if !goString.IsValidBirthDate(data.BirthDateString) {
			err := fmt.Errorf(helper.ErrorParameterInvalid, msgErrorBirthdate)
			return err
		}
	}
	return nil
}

func (mu *MemberUseCaseImpl) validateUpdateParam(ctx string, data *model.Member) (err error) {
	if len(data.Address.ProvinceID) == 0 {
		err = fmt.Errorf(helper.ErrorParameterRequired, "province id")
		return err
	}
	if len(data.Address.CityID) == 0 {
		err = fmt.Errorf(helper.ErrorParameterRequired, "city id")
		return err
	}
	if len(data.Address.DistrictID) == 0 {
		err = fmt.Errorf(helper.ErrorParameterRequired, "district id")
		return err
	}
	if len(data.Address.SubDistrictID) == 0 {
		err = fmt.Errorf(helper.ErrorParameterRequired, "sub district id")
		return err
	}
	if len(data.Address.ZipCode) == 0 {
		err = fmt.Errorf(helper.ErrorParameterRequired, "zip code")
		return err
	}

	if len(data.Address.Street1) == 0 {
		err = fmt.Errorf(helper.ErrorParameterRequired, "address")
		return err
	}

	if !golib.ValidateLatinOnly(data.Address.Street1) {
		err = errors.New(msgErrorStreetName)
		return fmt.Errorf("%s", err.Error())
	}

	data.Address.Street1 = helper.ClearHTML(data.Address.Street1)

	data.Address.Address = data.Address.Street1
	if len(data.Address.Street2) > 0 {
		if !golib.ValidateLatinOnly(data.Address.Street2) {
			err = errors.New(msgErrorStreetName)
			return fmt.Errorf("%s", err.Error())
		}

		data.Address.Street2 = helper.ClearHTML(data.Address.Street2)

		data.Address.Address = fmt.Sprintf("%s\n%s", data.Address.Street1, data.Address.Street2)
	}
	return nil
}

func (mu *MemberUseCaseImpl) validatePhone(ctx string, data *model.Member) (err error) {
	if len(data.Mobile) <= 0 {
		return fmt.Errorf(helper.ErrorParameterRequired, textMobileNumber)
	}

	if len(data.Mobile) > 0 && len(data.Mobile) < 8 {
		return fmt.Errorf(helper.ErrorParameterInvalid, fmt.Sprintf("%s %s", textMobileNumber, data.Mobile))
	}

	if len(data.Mobile) > 0 {
		if err = helper.ValidateMobileNumberMaxInput(data.Mobile); err != nil {
			return fmt.Errorf(helper.ErrorParameterInvalid, fmt.Sprintf("%s %s", textMobileNumber, data.Mobile))
		}
	}

	if len(data.Phone) > 0 {
		if err = helper.ValidatePhoneNumberMaxInput(data.Phone); err != nil {
			return fmt.Errorf(helper.ErrorParameterInvalid, "phone number")
		}
	}

	if len(data.Ext) > 0 {
		if err = helper.ValidateExtPhoneNumberMaxInput(data.Ext); err != nil {
			return fmt.Errorf(helper.ErrorParameterInvalid, "ext")
		}
	}

	if len(data.Phone) == 0 && len(data.Ext) > 0 {
		return fmt.Errorf(helper.ErrorParameterRequired, "phone number")
	}

	// when phone is empty extension should be empty
	if len(data.Phone) == 0 {
		data.Ext = ""
	}
	return nil
}

// PublishToKafkaDolphin function for publish new member to kafka dolphin
func (au *MemberUseCaseImpl) PublishToKafkaContact(ctxReq context.Context, payload corporateModel.ContactPayload, eventType string) error {
	contactPayload := contactModel.ContactKafka{
		EventType: eventType,
		Payload:   payload,
	}

	payloadJSON, err := json.Marshal(contactPayload)
	if err != nil {
		helper.SendErrorLog(ctxReq, "PublishToKafkaContact", "unmarshal_payload", err, contactPayload)
		return err
	}

	if err := au.QPublisher.PublishKafka(ctxReq, os.Getenv("KAFKA_SHARK_IMPORT"), strconv.Itoa(payload.ID), payloadJSON); err != nil {
		helper.SendErrorLog(ctxReq, "PublishKafka", "publish_to_contact", err, payload.Email)
		return err
	}
	return nil
}
