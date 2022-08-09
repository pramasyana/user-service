package usecase

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/jsonschema"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/tealeg/xlsx"
)

func (mu *MemberUseCaseImpl) processRows(row *xlsx.Row) (*model.Member, error) {
	passwordPlain := row.Cells[8].String()
	if len(passwordPlain) > 0 && !helper.IsValidPass(passwordPlain) {
		return nil, fmt.Errorf(helper.ErrorParameterInvalid, "format password")
	}

	member := model.Member{}

	firstName := row.Cells[0].String()
	lastName := row.Cells[1].String()
	email := row.Cells[2].String()

	genderStr := row.Cells[3].String()
	mobile := row.Cells[4].String()
	if len(mobile) <= 0 {
		return nil, fmt.Errorf("row with email %s has no mobile number", email)
	}
	phone := row.Cells[5].String()
	ext := row.Cells[6].String()
	birthDateStr := row.Cells[7].String()

	if len(passwordPlain) > 0 {
		member.Salt = mu.Hash.GenerateSalt()
		err := mu.Hash.ParseSalt(member.Salt)
		if err != nil {
			return nil, err
		}
		member.Password = base64.StdEncoding.EncodeToString(mu.Hash.Hash([]byte(passwordPlain)))
		if len(member.Salt) > 255 {
			return nil, fmt.Errorf("password for %s is invalid", email)
		}
	}
	address := row.Cells[9].String()
	province := row.Cells[10].String()
	provinceIDStr := row.Cells[11].String()
	provinceID := helper.FloatToString(provinceIDStr)

	city := row.Cells[12].String()
	cityIDStr := row.Cells[13].String()
	cityID := helper.FloatToString(cityIDStr)

	district := row.Cells[14].String()
	districtIDStr := row.Cells[15].String()
	districtID := helper.FloatToString(districtIDStr)

	subDistrict := row.Cells[16].String()
	subDistrictIDStr := row.Cells[17].String()
	subDistrictID := helper.FloatToString(subDistrictIDStr)

	zipCodeStr := row.Cells[18].String()
	zipCode := helper.FloatToString(zipCodeStr)
	jobTitle := row.Cells[19].String()
	department := row.Cells[20].String()
	isAdmin := false
	isStaff := false
	statusStr := row.Cells[21].String()
	signUpFrom := row.Cells[22].String()

	member.FirstName = firstName
	member.LastName = lastName
	member.Email = strings.ToLower(email)
	member.Gender = model.StringToGender(genderStr)
	member.GenderString = genderStr
	member.Mobile = mobile
	member.Phone = phone
	member.Ext = ext
	birthDate, _ := time.Parse("02/01/2006", birthDateStr)
	member.NewPassword = passwordPlain
	member.BirthDate = birthDate
	member.BirthDateString = birthDate.Format("02/01/2006")
	member.Address.ZipCode = zipCode
	member.Address.Street1 = address
	member.Address.CityID = cityID
	member.Address.City = city
	member.Address.ProvinceID = provinceID
	member.Address.Province = province
	member.Address.DistrictID = districtID
	member.Address.District = district
	member.Address.SubDistrictID = subDistrictID
	member.Address.SubDistrict = subDistrict
	member.JobTitle = jobTitle
	member.Department = department
	member.IsAdmin = isAdmin
	member.IsStaff = isStaff
	member.Status = model.StringToStatus(statusStr)
	member.StatusString = statusStr
	member.SignUpFrom = signUpFrom
	member.Type = "import"
	return &member, nil

}

// ParseMemberData parse member data from xls
func (mu *MemberUseCaseImpl) ParseMemberData(ctxReq context.Context, input []byte) ([]*model.Member, error) {
	ctx := "MemberUseCase-ParseMemberData"
	defer func() {
		if r := recover(); r != nil {
			helper.SendErrorLog(ctxReq, ctx, "recover_import_data", fmt.Errorf("%v", r), nil)
		}
	}()

	xlsBinary, err := xlsx.OpenBinary(input)
	if err != nil {
		return nil, err
	}

	if xlsBinary.Sheets[0].MaxCol < 22 {
		return nil, fmt.Errorf("excel format template not valid")
	}
	members := make([]*model.Member, 0)
	emails := make([]string, 0)
	for _, row := range xlsBinary.Sheets[0].Rows {
		if len(row.Cells) == 0 {
			continue
		}

		if row.Cells[0].String() == model.FieldFirstName ||
			(row.Cells[0].String() == "" && row.Cells[2].String() == "") {
			continue
		}

		email := row.Cells[2].String()
		if golib.StringInSlice(email, emails, false) {
			emails = nil
			return nil, fmt.Errorf("%s found as a duplicate email address in file source", email)
		}
		emails = append(emails, strings.ToLower(email))

		member, err := mu.processRows(row)
		if err != nil {
			return nil, err
		}

		members = append(members, member)
	}
	emails = nil
	return members, nil
}

// MigrateMember function for migrating member from squid
func (mu *MemberUseCaseImpl) MigrateMember(ctxReq context.Context, members *model.Members) <-chan ResultUseCase {
	ctx := "MemberUseCase-MigrateMember"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		var (
			errs []model.MemberError
		)

		tags[helper.TextArgs] = members
		for _, data := range members.Data {
			if mErr, err := mu.parseMigrationData(ctxReq, data, ctx); err != nil {
				tracer.SetError(ctxReq, err)
				errs = append(errs, mErr)
				continue
			}
		}

		if len(errs) > 0 {
			output <- ResultUseCase{Error: errors.New("failed to save member(s)"), ErrorData: errs, HTTPStatus: http.StatusBadRequest}
			return
		}
	})

	return output
}

func (mu *MemberUseCaseImpl) parseMigrationData(ctxReq context.Context, data model.Member, ctx string) (model.MemberError, error) {
	var errModel model.MemberError
	// check data existence first
	memberResult := <-mu.MemberRepoRead.Load(ctxReq, data.ID)
	if memberResult.Error == nil {
		err := fmt.Errorf("failed to save member %s, data exists", data.ID)
		return errModel, err
	}

	// check member when the data does not exist
	if strings.TrimSpace(memberResult.Error.Error()) == fmt.Errorf(helper.ErrorDataNotFound, labelMember).Error() {

		if err := mu.validate(&data); err != nil {
			return model.MemberError{ID: data.ID, Message: err.Error()}, err
		}

		saveResult := <-mu.MemberRepoWrite.Save(ctxReq, data)
		if saveResult.Error != nil {
			err := fmt.Errorf("failed to save member %s", data.ID)
			return model.MemberError{ID: data.ID, Message: err.Error()}, err
		}
	}
	return errModel, nil
}

// ImportMember function for import member from excel
func (mu *MemberUseCaseImpl) ImportMember(ctxReq context.Context, data *model.Member) <-chan ResultUseCase {
	ctx := "MemberUseCase-ImportMember"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags[helper.TextArgs] = data.Email
		response, err := mu.importMember(ctxReq, ctx, data)
		if err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		output <- ResultUseCase{Result: response}
	})

	return output
}

func (mu *MemberUseCaseImpl) importMember(ctxReq context.Context, ctx string, data *model.Member) (interface{}, error) {
	ctxP := "MemberUseCaseImpl-ImportMember"
	var hasPassword bool

	// validate schema
	mErr := jsonschema.ValidateTemp("import_member_params", data)
	if mErr != nil {
		helper.SendErrorLog(ctxReq, ctxP, "validate_params", mErr, data)
		return nil, mErr
	}

	data.Email = strings.ToLower(data.Email)

	// validate data first
	if err := mu.validateMemberData(ctxReq, data); err != nil {
		helper.SendErrorLog(ctxReq, ctxP, "validate_member_date", err, data)
		return nil, err
	}

	// generate member id
	if len(data.ID) <= 0 {
		data.ID = helper.GenerateMemberIDv2()
	}

	hasPassword = false
	// optional when password exists only
	if len(data.NewPassword) > 0 {
		// encode the new password then replace the old password and salt
		data.Salt = mu.Hash.GenerateSalt()
		err := mu.Hash.ParseSalt(data.Salt)
		if err != nil {
			tracer.SetError(ctxReq, err)
			return nil, err
		}

		data.Password = base64.StdEncoding.EncodeToString(mu.Hash.Hash([]byte(data.NewPassword)))
		hasPassword = true
	}

	// generate random string
	mix := data.Email + "-" + data.FirstName
	data.Token = helper.GenerateTokenByString(mix)

	saveResult := <-mu.MemberRepoWrite.Save(ctxReq, *data)
	if saveResult.Error != nil {
		err := errors.New(msgErrorSaveMember)
		tracer.SetError(ctxReq, err)
		return nil, err
	}

	// publish to kafka
	data.Created = time.Now()
	go mu.PublishToKafkaUser(ctxReq, data, textRegister)

	response := model.SuccessResponse{
		ID:          data.ID,
		Message:     helper.SuccessMessage,
		Token:       data.Token,
		HasPassword: hasPassword,
		FirstName:   data.FirstName,
		LastName:    data.LastName,
		Email:       strings.ToLower(data.Email),
	}
	return response, nil
}

// BulkImportMember bulk import
func (mu *MemberUseCaseImpl) BulkImportMember(ctxReq context.Context, data []*model.Member) <-chan ResultUseCase {
	ctx := "MemberUseCase-BulkImportMember"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, _ map[string]interface{}) {
		defer close(output)
		responses := make([]model.SuccessResponse, 0)

		var wg sync.WaitGroup
		chResp := make(chan ResultUseCase)

		for _, m := range data {
			wg.Add(1)
			go mu.populateMemberData(ctxReq, m, &wg, chResp)
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
			temp := resp.Result.(model.SuccessResponse)
			responses = append(responses, temp)
		}

		saveResult := <-mu.MemberRepoWrite.BulkImportSave(ctxReq, data)
		if saveResult.Error != nil {
			helper.SendErrorLog(ctxReq, ctx, "bulk_import", saveResult.Error, nil)
			err := errors.New(msgErrorSaveMember)
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusInternalServerError}
			responses = nil
			return
		}
		memberAfterSave := saveResult.Result.([]*model.Member)

		if err := mu.pushToThirdParty(ctxReq, memberAfterSave); err != nil {
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}

		output <- ResultUseCase{Result: responses}
	})

	return output
}

func (mu *MemberUseCaseImpl) populateMemberData(ctxReq context.Context, m *model.Member, wg *sync.WaitGroup, response chan<- ResultUseCase) {
	defer wg.Done()
	hasPassword, err := mu.validateBulkImportData(ctxReq, m)
	if err != nil {
		response <- ResultUseCase{Error: err}
	}
	response <- ResultUseCase{
		Result: model.SuccessResponse{
			ID:          m.ID,
			Message:     helper.SuccessMessage,
			Token:       m.Token,
			HasPassword: hasPassword,
			FirstName:   m.FirstName,
			LastName:    m.LastName,
			Email:       strings.ToLower(m.Email),
		},
	}
}

func (mu *MemberUseCaseImpl) validate(data *model.Member) error {
	var (
		err error
		ok  bool
	)

	data.Gender, ok = model.ValidateGender(data.GenderString)
	if !ok {
		err = fmt.Errorf(helper.ErrorParameterInvalid, textGender)
		return err
	}

	data.Status, ok = model.ValidateStatus(data.StatusString)
	if !ok {
		err = fmt.Errorf(helper.ErrorParameterInvalid, scopeStatus)
		return err
	}

	// parse and replace dob value
	if len(data.BirthDateString) > 0 {
		data.BirthDate, err = time.Parse(helper.FormatDOB, data.BirthDateString)
		if err != nil {
			err = fmt.Errorf(helper.ErrorParameterInvalid, msgErrorBirthdate)
			return err
		}
	}

	if len(data.LastLoginString) > 0 {
		data.LastLogin, err = time.Parse(time.RFC3339, data.LastLoginString)
		if err != nil {
			err = fmt.Errorf(helper.ErrorParameterInvalid, "last login date")
			return err
		}
	}

	if len(data.CreatedString) > 0 {
		data.Created, err = time.Parse(time.RFC3339, data.CreatedString)
		if err != nil {
			err = fmt.Errorf(helper.ErrorParameterInvalid, "created date")
			return err
		}
	}

	if len(data.LastModifiedString) > 0 {
		data.LastModified, err = time.Parse(time.RFC3339, data.LastModifiedString)
		if err != nil {
			err = fmt.Errorf(helper.ErrorParameterInvalid, "last modified date")
			return err
		}
	}

	return nil
}

func (mu *MemberUseCaseImpl) validateBulkImportData(ctxReq context.Context, m *model.Member) (hasPassword bool, err error) {
	tr := tracer.StartTrace(ctxReq, "MemberUseCase-validateBulkImportData")
	defer tr.Finish(map[string]interface{}{"member": m})

	if mErr := jsonschema.ValidateTemp("import_member_params", m); mErr != nil {
		return hasPassword, mErr
	}

	m.Email = strings.ToLower(m.Email)

	// validate data first
	if err := mu.validateMemberData(tr.NewChildContext(), m); err != nil {
		return hasPassword, err
	}

	// generate member id
	if len(m.ID) <= 0 {
		m.ID = helper.GenerateMemberIDv2()
	}

	// generate random string
	mix := m.Email + "-" + m.FirstName
	m.Token = helper.GenerateTokenByString(mix)
	return hasPassword, nil
}

func (mu *MemberUseCaseImpl) pushToThirdParty(ctxReq context.Context, data []*model.Member) error {
	tc := tracer.StartTrace(ctxReq, "MemberUseCaseImpl-pushToThirdParty")
	defer tc.Finish(nil)
	ch := make(chan error)
	var wg sync.WaitGroup
	for _, m := range data {
		wg.Add(1)
		// publish to kafka async
		go mu.pushAsync(tc.NewChildContext(), m, textRegister, ch, &wg)
	}
	go func() {
		wg.Wait()
		close(ch)
	}()

	for err := range ch {
		if err != nil {
			return err
		}
	}

	return nil
}

func (mu *MemberUseCaseImpl) pushAsync(ctxReq context.Context, member *model.Member, eventType string, ch chan<- error, wg *sync.WaitGroup) {
	defer wg.Done()
	member.Created = time.Now()
	if err := mu.PublishToKafkaUser(ctxReq, member, textRegister); err != nil {
		ch <- err
	}
}
