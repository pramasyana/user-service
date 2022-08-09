package model

import (
	"time"

	"github.com/spf13/cast"
	"gopkg.in/guregu/null.v4"
	"gopkg.in/guregu/null.v4/zero"
)

// SetMerchantData function to struct data for create
func (s *B2CMerchantDataV2) SetMerchantData(params *B2CMerchantCreateInput) {

	if params.Source != "" {
		s.Source = zero.StringFrom(params.Source)
	}

	if s.Source.String == "" {
		s.Source = zero.StringFrom("cms")
	}

	if params.UserID != "" {
		s.UserID = params.UserID
	}

	s = s.SetMerchantBasicInfo(params)

	if params.UpgradeStatus != "" {
		s.UpgradeStatus = zero.StringFrom(params.UpgradeStatus)
	}
	if params.ProductType != "" {
		s.ProductType = zero.StringFrom(params.ProductType)
	}

	if params.AccountManager != "" {
		s.AccountManager = zero.StringFrom(params.AccountManager)
	}

	if params.RichContent != "" {
		s.RichContent = zero.StringFrom(params.RichContent)
	}

	if params.NotificationPreferences != 0 {
		notif := cast.ToInt64(params.NotificationPreferences)
		s.NotificationPreferences = zero.Int(null.IntFrom(notif))
	}

	if params.Acquisitor != "" {
		s.Acquisitor = zero.StringFrom(params.Acquisitor)
	}

	if params.LaunchDev != "" {
		s.LaunchDev = zero.StringFrom(params.LaunchDev)
	}

	s = s.SetMerchantStore(params)

	s.IsActive = params.IsActive
	s.Status = params.Status

	if params.SkuLive != "" {
		skuLive, _ := time.Parse(DefaultInputFormat, params.SkuLive)
		s.SkuLive = null.TimeFrom(skuLive)
	}

	if params.MouDate != "" {
		mouDate, _ := time.Parse(DefaultInputFormat, params.MouDate)
		s.MouDate = null.TimeFrom(mouDate)
	}

	if params.Note != "" {
		s.Note.String = params.Note
	}

	if params.AgreementDate != "" {
		agreementDate, _ := time.Parse(DefaultDateFormat, params.AgreementDate)
		s.AgreementDate = null.TimeFrom(agreementDate)
	}
	if params.BusinessType != "" {
		s.BusinessType = zero.StringFrom(params.BusinessType)
	}

	s.setMerchantAddress(params)
	s.setMerchantAddress2(params)

	s.SetMerchantAccount(params)
}

func (s *B2CMerchantDataV2) setMerchantAddress(params *B2CMerchantCreateInput) {
	if params.MerchantVillageID != "" {
		s.MerchantVillageID = zero.StringFrom(params.MerchantVillageID)
	}
	if params.MerchantCityID != "" {
		s.MerchantCityID = zero.StringFrom(params.MerchantCityID)
	}
	if params.StoreVillageID != "" {
		s.StoreVillageID = zero.StringFrom(params.StoreVillageID)
	}
	if params.StoreZipCode != "" {
		s.StoreZipCode = zero.StringFrom(params.StoreZipCode)
	}
	if params.MerchantEmail != "" {
		s.MerchantEmail = zero.StringFrom(params.MerchantEmail)
	}
	if params.BankBranch != "" {
		s.BankBranch = zero.StringFrom(params.BankBranch)
	}
	if params.ZipCode != "" {
		s.ZipCode = zero.StringFrom(params.ZipCode)
	}
	if params.MerchantRank != "" {
		s.MerchantRank = zero.StringFrom(params.MerchantRank)
	}
}

func (s *B2CMerchantDataV2) setMerchantAddress2(params *B2CMerchantCreateInput) {
	if params.MerchantDistrictID != "" {
		s.MerchantDistrictID = zero.StringFrom(params.MerchantDistrictID)
	}
	if params.MerchantProvinceID != "" {
		s.MerchantProvinceID = zero.StringFrom(params.MerchantProvinceID)
	}
	if params.MerchantVillage != "" {
		s.MerchantVillage = zero.StringFrom(params.MerchantVillage)
	}
	if params.MerchantCity != "" {
		s.MerchantCity = zero.StringFrom(params.MerchantCity)
	}
	if params.MerchantProvince != "" {
		s.MerchantProvince = zero.StringFrom(params.MerchantProvince)
	}
	if params.MerchantDistrict != "" {
		s.MerchantDistrict = zero.StringFrom(params.MerchantDistrict)
	}
}

// SetMerchantStore function to struct data for create
func (s *B2CMerchantDataV2) SetMerchantStore(params *B2CMerchantCreateInput) *B2CMerchantDataV2 {
	if (params.StoreClosureDate != "" && params.StoreClosureDate != "0") && (params.StoreReopenDate != "" && params.StoreReopenDate != "0") {
		srd, err := time.Parse(DefaultInputFormat, params.StoreReopenDate)
		if err != nil {
			return nil
		}
		s.StoreReopenDate = &srd

		scd, _ := time.Parse(DefaultInputFormat, params.StoreClosureDate)
		s.StoreClosureDate = &scd

		// if time now is between in closure date
		nowDate, _ := time.Parse(DefaultInputFormat, time.Now().Format(time.RFC3339))

		openDateTime, _ := time.Parse(DefaultInputFormat, params.StoreReopenDate)
		closeDateTime, _ := time.Parse(DefaultInputFormat, params.StoreClosureDate)

		// set isClosed
		s.IsClosed = checkIsClosed(nowDate, closeDateTime, openDateTime)

	} else if params.StoreClosureDate == "0" || params.StoreReopenDate == "0" {
		s.StoreReopenDate = nil
		s.StoreClosureDate = nil

		// set isClosed false
		s.IsClosed = null.BoolFrom(false)
	} else if len(params.StoreClosureDate) < 1 && len(params.StoreReopenDate) < 1 {
		s.StoreClosureDate = nil
		s.StoreReopenDate = nil
		s.IsClosed = null.BoolFrom(false)
	}

	if params.StoreActiveShippingDate != "" && params.StoreActiveShippingDate != "0" {
		sasd, _ := time.Parse(DefaultInputFormat, params.StoreActiveShippingDate)

		s.StoreActiveShippingDate = &sasd
	} else if params.StoreActiveShippingDate == "0" {
		s.StoreActiveShippingDate = nil
	}
	return s
}

func checkIsClosed(nowDate time.Time, closeDateTime time.Time, openDateTime time.Time) null.Bool {
	// set isClosed true
	var result null.Bool
	if nowDate.Equal(closeDateTime) || (nowDate.After(closeDateTime) && nowDate.Before(openDateTime)) {
		result = null.BoolFrom(true)
	}

	// set isClosed false
	if nowDate.Equal(openDateTime) || nowDate.Before(closeDateTime) || nowDate.After(openDateTime) {
		result = null.BoolFrom(false)
	}
	return result
}

// SetMerchantBasicInfo function to struct data for create
func (s *B2CMerchantDataV2) SetMerchantBasicInfo(params *B2CMerchantCreateInput) *B2CMerchantDataV2 {

	if params.MerchantName != "" {
		s.MerchantName = params.MerchantName
	}

	if params.MerchantLogo != "" {
		s.MerchantLogo = zero.StringFrom(params.MerchantLogo)
	}

	if params.VanityURL != "" {
		s.VanityURL = zero.StringFrom(params.VanityURL)
	}

	if params.MerchantCategory != "" {
		s.MerchantCategory = zero.StringFrom(params.MerchantCategory)
	}

	if params.CompanyName != "" && params.BusinessType == PerusahaanType {
		s.CompanyName = zero.StringFrom(params.CompanyName)
	}

	if params.DailyOperationalStaff != "" {
		s.DailyOperationalStaff = zero.StringFrom(params.DailyOperationalStaff)

		// if set to 1 space, set daily operational staff to null
		if s.DailyOperationalStaff.String == " " {
			s.DailyOperationalStaff = zero.StringFrom("")
		}
	}

	if params.MerchantAddress != "" && params.BusinessType == PerusahaanType {
		s.MerchantAddress = zero.StringFrom(params.MerchantAddress)
	}

	if params.StoreAddress != "" {
		s.StoreAddress = zero.StringFrom(params.StoreAddress)
	}

	if params.MerchantDescription != "" {
		s.MerchantDescription = zero.StringFrom(params.MerchantDescription)

		// if set to 1 space, set description to null
		if s.MerchantDescription.String == " " {
			s.MerchantDescription = zero.StringFrom("")
		}
	}
	return s
}

// SetMerchantAccount function to struct data for create
func (s *B2CMerchantDataV2) SetMerchantAccount(params *B2CMerchantCreateInput) {
	if params.Pic != "" {
		s.Pic = zero.StringFrom(params.Pic)
	}

	if params.PicOccupation != "" {
		s.PicOccupation = zero.StringFrom(params.PicOccupation)
	}

	if params.PicKtpFile != "" {
		s.PicKtpFile = zero.StringFrom(params.PicKtpFile)
	}

	if params.PhoneNumber != "" {
		s.PhoneNumber = zero.StringFrom(params.PhoneNumber)

		// if set to 1 space, set phone to null
		if s.PhoneNumber.String == " " {
			s.PhoneNumber = zero.StringFrom("")
		}
	}

	if params.MobilePhoneNumber != "" {
		s.MobilePhoneNumber = zero.StringFrom(params.MobilePhoneNumber)
	}

	if params.AdditionalEmail != "" {
		s.AdditionalEmail = zero.StringFrom(params.AdditionalEmail)
	}

	if params.AccountHolderName != "" {
		s.AccountHolderName = zero.StringFrom(params.AccountHolderName)
	}

	if params.BankID != 0 {
		bankID := cast.ToInt64(params.BankID)
		s.BankID = zero.IntFrom(bankID)
	}

	if params.AccountNumber != "" {
		s.AccountNumber = zero.StringFrom(params.AccountNumber)
	}

	if params.Npwp != "" {
		s.Npwp = zero.StringFrom(params.Npwp)
	}

	if params.NpwpHolderName != "" {
		s.NpwpHolderName = zero.StringFrom(params.NpwpHolderName)
	}

	if params.NpwpFile != "" {
		s.NpwpFile = zero.StringFrom(params.NpwpFile)
	}

	s.IsPKP = params.IsPKP

}

func (s *B2CMerchantDataV2) ChecksIsClosed() null.Bool {
	// isClosed True
	nowDates, _ := time.Parse(DefaultInputFormat, time.Now().Format(time.RFC3339))
	if nowDates.Equal(*s.StoreClosureDate) || (nowDates.After(*s.StoreClosureDate) && nowDates.Before(*s.StoreReopenDate)) {
		s.IsClosed = null.BoolFrom(true)
	}
	// isClosed false
	if nowDates.Equal(*s.StoreReopenDate) || nowDates.Before(*s.StoreClosureDate) || nowDates.After(*s.StoreReopenDate) {
		s.IsClosed = null.BoolFrom(false)
	}
	return s.IsClosed
}
