package usecase

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/Bhinneka/golib"
	stringLib "github.com/Bhinneka/golib/string"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
)

var (
	errMerchantNotFound                = errors.New("merchant not found")
	errMerchantNotValidForUpgrade      = errors.New("merchant does not have upgrade request")
	errMerchantNotValidForClearUpgrade = errors.New("merchant upgrade status already cleared")
	errClearPendingUpgrade             = errors.New("merchant upgrade status still in pending")
	errUnableToRejectUpgrade           = errors.New("cannot reject request, merchant upgrade already approved")
	errUnableToRejectRegistration      = errors.New("merchant is already active")
	errBankNotExist                    = errors.New("bank doesn't exists")
	timeFormat                         = "2006-01-02T15:04:05Z07:00"
)

// ValidateMerchantBank function for validate merchant bank
func (m *MerchantUseCaseImpl) ValidateMerchantBank(ctxReq context.Context, data *model.B2CMerchantCreateInput, merchant *model.B2CMerchantDataV2) error {
	bankID := int(data.BankID)
	if bankID != 0 {
		bankResult := <-m.MerchantBankRepo.FindActiveMerchantBankByID(ctxReq, bankID)
		if bankResult.Error != nil {
			return errBankNotExist
		}
		if bankResult.Result == nil {
			return errBankNotExist
		}

		bankDetail, ok := bankResult.Result.(model.B2CMerchantBankData)
		if !ok {
			return errBankNotExist
		}
		merchant.BankName = bankDetail.BankName
	}
	return nil
}

// ValidateMerchantData function to validate merchant data
func (m *MerchantUseCaseImpl) ValidateMerchantData(ctxReq context.Context, data *model.B2CMerchantCreateInput) error {
	checkID := m.MerchantRepo.LoadMerchant(ctxReq, data.ID, private)
	if checkID.Result != nil {
		return fmt.Errorf("merchant ID already exist")
	}

	// validate email unique
	checkEmail := m.MerchantRepo.FindMerchantByEmail(ctxReq, data.MerchantEmail)
	if checkEmail.Result != nil {
		return fmt.Errorf("email merchant %s already exist", data.MerchantEmail)
	}

	// validate user id unique
	checkUser := m.MerchantRepo.FindMerchantByUser(ctxReq, data.UserID)
	if checkUser.Result != nil {
		return fmt.Errorf("user %s already exist", data.UserID)
	}

	// validate name unique
	checkName := m.MerchantRepo.FindMerchantByName(ctxReq, data.MerchantName)
	if checkName.Result != nil {
		return fmt.Errorf("name %s already exist", data.MerchantName)
	}
	// validate close date
	if data.StoreClosureDate != "" || data.StoreReopenDate != "" {
		err := ValidateMerchantCloseDate(data.StoreClosureDate, data.StoreReopenDate)
		if err != nil {
			return err
		}
	}

	// validate daily
	if data.DailyOperationalStaff != "" {
		err := ValidateDailyOperationalStaff(data.DailyOperationalStaff)
		if err != nil {
			return err
		}
	}

	if data.VanityURL != "" {
		// get merchant data
		merchantDataResultBySlug := m.MerchantRepo.FindMerchantBySlug(ctxReq, data.VanityURL)
		if merchantDataResultBySlug.Result != nil {
			return fmt.Errorf("slug already exists")
		}
	}

	return nil
}

// ValidateMerchantCloseDate validate function for merchant only would be here
func ValidateMerchantCloseDate(storeClosureDate, storeReopenDate string) error {
	// set now date
	nowDate, _ := helper.ConvertTimeToDate(time.Now())

	status := helper.XNOR(storeClosureDate == "0", storeReopenDate == "0")

	if !status {
		return fmt.Errorf("tanggal buka dan tanggal tutup tidak boleh kosong")
	}

	openDateTime, _ := time.Parse(model.DefaultInputFormat, storeReopenDate)
	closeDateTime, _ := time.Parse(model.DefaultInputFormat, storeClosureDate)

	if closeDateTime.Before(nowDate) {
		return fmt.Errorf("tanggal tutup sudah lewat")
	}

	if openDateTime.Before(nowDate) {
		return fmt.Errorf("tanggal buka sudah lewat")
	}

	if closeDateTime.Equal(openDateTime) {
		return fmt.Errorf("tanggal tutup tidak boleh sama dengan tanggal buka")
	}

	if closeDateTime.After(openDateTime) {
		return fmt.Errorf("tanggal tutup harus sebelum tanggal buka")
	}

	if openDateTime.Before(closeDateTime) {
		return fmt.Errorf("tanggal buka harus setelah tanggal buka")
	}

	return nil
}

// ValidateDailyOperationalStaff function for validate daily operational staff
func ValidateDailyOperationalStaff(dailyOperationalStaff string) error {
	var (
		regexRule = `^ *([a-zA-Z0-9,.] ?)+ *$`
	)

	if dailyOperationalStaff != "" && dailyOperationalStaff != " " {
		if len(dailyOperationalStaff) < 3 {
			return fmt.Errorf("panjang minimal nama operational staff adalah 3 karakter")
		}

		if len(dailyOperationalStaff) > 200 {
			return fmt.Errorf("panjang minimal nama operational staff adalah 200 karakter")
		}

		checkName := regexp.MustCompile(regexRule).MatchString
		if !checkName(dailyOperationalStaff) {
			return fmt.Errorf("nama operational staff hanya berupa karakter alphanumeric, spasi ( ), koma (,) dan titik (.)")
		}
	}

	return nil
}

// validateMerchantFieldUpgrade function for validating merchant data
func (m *MerchantUseCaseImpl) validateMerchantFieldUpgrade(data *model.B2CMerchantCreateInput) (*model.B2CMerchantCreateInput, error) {
	var ok bool

	_, ok = model.ValidateMerchantType(data.MerchantTypeString)
	if !ok {
		err := fmt.Errorf(helper.ErrorParameterInvalid, "merchantType")
		return data, err
	}

	switch strings.ToUpper(data.MerchantTypeString) {
	case model.ManageString:
		data.UpgradeStatus = model.PendingManageString
	case model.AssociateString:
		data.UpgradeStatus = model.PendingAssociateString
	}

	data.GenderPic, ok = model.ValidateGenderPic(data.GenderPicString)
	if !ok {
		err := fmt.Errorf(helper.ErrorParameterInvalid, "genderPic")
		return data, err
	}

	data.MerchantGroup, ok = model.ValidateMerchantGroup(data.MerchantGroup)
	if !ok {
		err := fmt.Errorf(helper.ErrorParameterInvalid, "merchantGroup")
		return data, err
	}

	data.UpgradeStatus, ok = model.ValidateUpgradeStatus(data.UpgradeStatus)
	if !ok {
		err := fmt.Errorf(helper.ErrorParameterInvalid, "upgradeStatus")
		return data, err
	}
	if !model.ValidateProductType(data.ProductType) {
		return data, fmt.Errorf(helper.ErrorParameterInvalid, "productType")
	}
	data.ProductType = strings.ToUpper(data.ProductType)

	if data.LegalEntity > 0 {
		legalEntity := m.MerchantRepo.LoadLegalEntity(context.Background(), data.LegalEntity)
		if legalEntity.Error != nil {
			return data, fmt.Errorf(helper.ErrorParameterInvalid, "legalEntity")
		}
	}

	return data, nil
}

// validateMerchantField function for validating merchant data
func (m *MerchantUseCaseImpl) validateMerchantField(ctxReq context.Context, data *model.B2CMerchantCreateInput) error {
	errorValidateBasicMerchant := m.validateMandatoryField(data)
	if errorValidateBasicMerchant != nil {
		return errorValidateBasicMerchant
	}

	// validate mobile phone value existence
	if len(data.MobilePhoneNumber) <= 0 {
		return fmt.Errorf("mobile phone number required")
	}

	errValidate := m.validateAdditionalField(data)
	if errValidate != nil {
		return errValidate
	}

	// validate npwp file
	errorMerchantFile := m.validateMerchantFile(data)
	if errorMerchantFile != nil {
		return errorMerchantFile
	}
	if err := m.validateFieldContent(data); err != nil {
		return err
	}
	if err := m.validateKTPAndNPWPFolder(data); err != nil {
		return err
	}

	return nil
}

// validateAdditionalField function for validating merchant additional information
func (m *MerchantUseCaseImpl) validateAdditionalField(data *model.B2CMerchantCreateInput) error {
	var ok bool
	if data.GenderPicString != "" {
		data.GenderPic, ok = model.ValidateGenderPic(data.GenderPicString)
		if !ok {
			err := fmt.Errorf(helper.ErrorParameterInvalid, "genderPic")
			return err
		}
	}

	if data.MerchantGroup != "" {
		data.MerchantGroup, ok = model.ValidateMerchantGroup(data.MerchantGroup)
		if !ok {
			err := fmt.Errorf(helper.ErrorParameterInvalid, "merchantGroup")
			return err
		}
	}
	if !model.ValidateMerchantStatus(data.Status) {
		return fmt.Errorf(helper.ErrorParameterInvalid, "status")
	}
	return nil
}

func isBusinessTypeValid(input string) error {
	if strings.ToLower(input) != model.PeroranganType && strings.ToLower(input) != model.PerusahaanType {
		return fmt.Errorf(helper.ErrorParameterInvalid, "business type")
	}
	return nil
}

func (m *MerchantUseCaseImpl) validateMerchantName(merchantName string) error {
	// validate merchant name value existence
	if len(merchantName) <= 0 {
		err := fmt.Errorf("merchant name required")
		return err
	}
	if !helper.ValidationMerchantName(merchantName) {
		err := fmt.Errorf(helper.ErrorParameterInvalid, "merchantName")
		return err
	}
	if !golib.ValidateLatinOnly(merchantName) {
		err := fmt.Errorf("merchant name only latin character")
		return err
	}
	return nil
}

// validateMandatoryField function for validating merchant data basic information
func (m *MerchantUseCaseImpl) validateMandatoryField(data *model.B2CMerchantCreateInput) error {
	if err := isBusinessTypeValid(data.BusinessType); err != nil {
		return err
	}
	if err := m.validateMerchantName(data.MerchantName); err != nil {
		return err
	}

	// validate phone number value existence
	if len(data.PhoneNumber) <= 0 {
		return fmt.Errorf("phone number required")
	}

	// validate merchant description value existence
	if len(data.MerchantDescription) <= 0 {
		return fmt.Errorf("merchant description required")
	}

	// validate business type value existence
	errorMerchantBusiness := m.validateMerchantCompany(data)
	if errorMerchantBusiness != nil {
		return errorMerchantBusiness
	}

	// validate pic value existence
	if len(data.Pic) <= 0 {
		return fmt.Errorf("Pic required")
	}

	// validate picktp
	if len(data.PicKtpFile) <= 0 {
		return fmt.Errorf("Pic KTP required")
	}

	// validate pic occupation value existence
	if len(data.PicOccupation) <= 0 {
		return fmt.Errorf("pic occupation required")
	}
	return nil
}

// validateMerchantCompany function for validating merchant data for business types
func (m *MerchantUseCaseImpl) validateMerchantCompany(data *model.B2CMerchantCreateInput) error {
	if strings.ToLower(data.BusinessType) != model.PerusahaanType {
		return nil
	}
	if len(data.CompanyName) <= 0 {
		return fmt.Errorf("company name required")
	}

	if len(data.CompanyName) < 3 || len(data.CompanyName) > 200 {
		return fmt.Errorf("company name min 3 and max 200 character")
	}
	if !golib.ValidateLatinOnly(data.CompanyName) {
		return fmt.Errorf("company name only latin character")
	}

	if err := m.validateAddress(data); err != nil {
		return err
	}

	if data.LegalEntity > 0 {
		legalEntity := m.MerchantRepo.LoadLegalEntity(context.Background(), data.LegalEntity)
		if legalEntity.Error != nil {
			return fmt.Errorf(helper.ErrorParameterInvalid, "legalEntity")
		}
	}
	if data.NumberOfEmployee > 0 {
		companySize := m.MerchantRepo.LoadCompanySize(context.Background(), data.NumberOfEmployee)
		if companySize.Error != nil {
			return fmt.Errorf(helper.ErrorParameterInvalid, "numberOfEmployee")
		}
	}
	return nil
}

func (m *MerchantUseCaseImpl) validateAddress(data *model.B2CMerchantCreateInput) error {
	if len(data.MerchantAddress) <= 0 {
		return fmt.Errorf("company address required")
	}

	if len(data.MerchantAddress) < 5 || len(data.MerchantAddress) > 255 {
		return fmt.Errorf("company address min 5 and max 255 character")
	}

	if !stringLib.ValidateLatinOnlyExcepTagCurly(data.MerchantAddress) {
		return fmt.Errorf("company address only latin character except '<' '>' '{' '}' ")
	}
	return nil
}

// validateMerchantFile function for validating merchant data file
func (m *MerchantUseCaseImpl) validateMerchantFile(data *model.B2CMerchantCreateInput) error {
	if len(data.NpwpFile) <= 0 {
		return fmt.Errorf("npwp file required")
	}

	if len(data.NpwpHolderName) <= 0 {
		err := fmt.Errorf("npwp holder name required")
		return err
	}

	// validate npwp value existence
	if len(data.Npwp) <= 0 {
		return fmt.Errorf("npwp required")
	}

	// validate account number value existence
	if len(data.AccountNumber) <= 0 {
		err := fmt.Errorf("account number required")
		return err
	}

	// validate account holder name value existence
	if len(data.AccountHolderName) <= 0 {
		return fmt.Errorf("account holder name required")
	}

	// validate bank branch value existence
	if len(data.BankBranch) <= 0 {
		return fmt.Errorf("bank branch required")
	}

	return nil
}

func (m *MerchantUseCaseImpl) validateUpdateField(ctxReq context.Context, payload *model.B2CMerchantCreateInput, oldData model.B2CMerchantDataV2) error {
	if len(payload.MerchantDescription) > 0 && !stringLib.ValidateLatinOnlyExcepTag(payload.MerchantDescription) {
		return fmt.Errorf("merchant description only accept latin character except '<' '>' ")
	}

	if strings.ToLower(payload.BusinessType) == model.PerusahaanType && len(payload.MerchantAddress) > 0 && !stringLib.ValidateLatinOnlyExcepTagCurly(payload.MerchantAddress) {
		return fmt.Errorf("company address only accept latin character except '<' '>' '{' '}' ")
	}

	if err := m.validateAdditionalField(payload); err != nil {
		return err
	}
	// validate close date
	if err := m.checkDateCloseAndReopen(payload, oldData); err != nil {
		return err
	}
	if payload.UpgradeStatus != "" {
		if err := m.validateStatus(payload.UpgradeStatus); err != nil {
			return err
		}
	}

	return nil
}

func (m *MerchantUseCaseImpl) checkDateCloseAndReopen(payload *model.B2CMerchantCreateInput, oldData model.B2CMerchantDataV2) error {
	var oldDataClose string
	var oldDataOpen string
	closeDateTime, _ := time.Parse(model.DefaultInputFormat, payload.StoreClosureDate)
	reopenDateTime, _ := time.Parse(model.DefaultInputFormat, payload.StoreReopenDate)
	closeDate := closeDateTime.Format(timeFormat)
	reopenDate := reopenDateTime.Format(timeFormat)
	if oldData.StoreClosureDate != nil || oldData.StoreReopenDate != nil {
		oldDataClose = oldData.StoreClosureDate.Local().Format(timeFormat)
		oldDataOpen = oldData.StoreReopenDate.Local().Format(timeFormat)
	}
	if closeDate != oldDataClose || reopenDate != oldDataOpen {
		if payload.StoreClosureDate != "" || payload.StoreReopenDate != "" {
			if err := ValidateMerchantCloseDate(payload.StoreClosureDate, payload.StoreReopenDate); err != nil {
				return err
			}
		}
	}
	return nil
}

func (m *MerchantUseCaseImpl) validateStatus(status string) error {
	if !golib.StringInSlice(strings.ToUpper(status), []string{model.PendingManageString, model.AssociateString, model.RegularString, model.PendingAssociateString, model.ActiveString, model.RejectManageString, model.RejectAssociateString}, false) {
		return fmt.Errorf(helper.ErrorParameterInvalid, "upgradeStatus")
	}
	return nil
}

// refer to gws required data when create merchant from cms
func (m *MerchantUseCaseImpl) validateBasicData(data *model.B2CMerchantCreateInput) error {
	if err := m.validateMerchantName(data.MerchantName); err != nil {
		return err
	}

	if err := m.validateAdditionalField(data); err != nil {
		return err
	}
	if err := m.validateKTPAndNPWPFolder(data); err != nil {
		return err
	}
	return m.validateFieldContent(data)
}

func (m *MerchantUseCaseImpl) validateFieldContent(data *model.B2CMerchantCreateInput) error {
	if data.MerchantDescription != "" && !stringLib.ValidateLatinOnlyExcepTag(data.MerchantDescription) {
		return fmt.Errorf("merchant description only latin character except '<' '>'")
	}
	if len(data.PhoneNumber) > 0 && !golib.ValidateNumeric(data.PhoneNumber) {
		return fmt.Errorf("phone number only numeric")
	}
	// validate pic value existence
	if len(data.Pic) > 0 && stringLib.ValidateAlphaNumericInput(data.Pic) != nil {
		return fmt.Errorf("pic only alphanumeric")
	}
	// validate pic occupation value existence
	if len(data.PicOccupation) > 0 && stringLib.ValidateAlphaNumericInput(data.PicOccupation) != nil {
		return fmt.Errorf("pic occupation only alphanumeric")
	}

	if len(data.AccountNumber) > 0 && !golib.ValidateNumeric(data.AccountNumber) {
		return fmt.Errorf("account number only alphanumeric")
	}
	if err := m.validateMoreFieldContent(data); err != nil {
		return err
	}

	return nil
}

func (m *MerchantUseCaseImpl) validateMoreFieldContent(data *model.B2CMerchantCreateInput) error {
	err := validateNpwpHolder(data)
	if err != nil {
		return err
	}
	if len(data.Npwp) > 0 && (!golib.ValidateNumeric(data.Npwp) || len(data.Npwp) < 15 || len(data.Npwp) > 15) {
		return fmt.Errorf("npwp number length must be equal to 15")
	}

	// validate account holder name value existence
	if len(data.AccountHolderName) > 0 && (!golib.ValidateLatinOnly(data.AccountHolderName) || (len(data.AccountHolderName) < 3 || len(data.AccountHolderName) > 200)) {
		return fmt.Errorf("account holder name min 3 and max 200 valid character")
	}
	// validate daily staff value existence
	if len(data.DailyOperationalStaff) > 0 && (!golib.ValidateLatinOnly(data.DailyOperationalStaff) || (len(data.DailyOperationalStaff) < 3 || len(data.DailyOperationalStaff) > 250)) {
		return fmt.Errorf("opperational staff min 3 and max 250 latin character")
	}
	if len(data.MobilePhoneNumber) > 0 && helper.ValidateMobileNumberMaxInput(data.MobilePhoneNumber) != nil {
		return fmt.Errorf("mobile phone number is in bad format")
	}
	return nil
}

func (m *MerchantUseCaseImpl) validateRequiredData(ctxReq context.Context, input *model.B2CMerchantCreateInput) error {
	oldData := model.B2CMerchantDataV2{}
	if err := m.validateUpdateField(ctxReq, input, oldData); err != nil {
		return err
	}
	if err := m.validateMerchantCompany(input); err != nil {
		return err
	}
	if err := m.validateKTPAndNPWPFolder(input); err != nil {
		return err
	}
	return nil
}
func (m *MerchantUseCaseImpl) validateRequiredDataSelf(ctxReq context.Context, input *model.B2CMerchantCreateInput, oldData model.B2CMerchantDataV2) error {
	if err := m.validateUpdateField(ctxReq, input, oldData); err != nil {
		return err
	}
	if err := m.validateMerchantCompany(input); err != nil {
		return err
	}
	if err := m.validateKTPAndNPWPFolder(input); err != nil {
		return err
	}
	return nil
}

func validateNpwpHolder(data *model.B2CMerchantCreateInput) error {
	if len(data.NpwpHolderName) > 0 && (!golib.ValidateLatinOnly(data.NpwpHolderName) || len(data.NpwpHolderName) < 3 || len(data.NpwpHolderName) > 200) {
		return fmt.Errorf("npwp holder name min 3 and max 200 character")
	}
	return nil
}

func (m *MerchantUseCaseImpl) validateKTPAndNPWPFolder(data *model.B2CMerchantCreateInput) error {
	picKtpFile := strings.Contains(data.PicKtpFile, "static.bmdstatic.com")
	npwpFile := strings.Contains(data.NpwpFile, "static.bmdstatic.com")

	if picKtpFile {
		err := fmt.Errorf("KTP File should not be in folder at static.bmdstatic.com")
		return err
	}
	if npwpFile {
		err := fmt.Errorf("NPWP File should not be in folder at static.bmdstatic.com")
		return err
	}

	return nil
}
