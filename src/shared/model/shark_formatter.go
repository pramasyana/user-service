package model

import (
	"encoding/json"
	"time"

	"github.com/Bhinneka/user-service/helper"
)

// RestructCorporateLeads function for restruct from cdc
func RestructCorporateLeads(pl B2BLeadsCDC) B2BLeads {
	var leadsModel B2BLeads
	leadsModel.ID = pl.ID
	leadsModel.Name = pl.Name
	leadsModel.Email = pl.Email
	leadsModel.Phone = pl.Phone
	leadsModel.SourceID = pl.SourceID
	leadsModel.Source = pl.Source
	leadsModel.SourceStatus = pl.SourceStatus
	leadsModel.Status = pl.Status
	leadsModel.Notes = pl.Notes
	leadsModel.CreatedAt = pl.CreatedAt
	leadsModel.ModifiedAt = pl.ModifiedAt
	leadsModel.CreatedBy = pl.CreatedBy
	leadsModel.ModifiedBy = pl.ModifiedBy
	return leadsModel
}

// RestructCorporateContactAddress function for restruct from cdc
func RestructCorporateContactAddress(pl B2BContactAddressCDC) B2BContactAddress {
	var contactAddressModel B2BContactAddress
	contactAddressModel.ID = pl.ID
	contactAddressModel.Name = pl.Name
	contactAddressModel.PicName = pl.PicName
	contactAddressModel.Phone = pl.Phone
	contactAddressModel.Street = pl.Street
	contactAddressModel.ProvinceID = pl.ProvinceID
	contactAddressModel.ProvinceName = pl.ProvinceName
	contactAddressModel.CityID = pl.CityID
	contactAddressModel.CityName = pl.CityName
	contactAddressModel.DistrictID = pl.DistrictID
	contactAddressModel.DistrictName = pl.DistrictName
	contactAddressModel.SubDistrictID = pl.SubDistrictID
	contactAddressModel.SubDistrictName = pl.SubDistrictName
	contactAddressModel.IsBilling = pl.IsBilling
	contactAddressModel.IsShipping = pl.IsShipping
	contactAddressModel.PostalCode = pl.PostalCode
	contactAddressModel.CreatedAt = pl.CreatedAt
	contactAddressModel.ModifiedAt = pl.ModifiedAt
	contactAddressModel.CreatedBy = pl.CreatedBy
	contactAddressModel.ModifiedBy = pl.ModifiedBy
	contactAddressModel.IsDelete = pl.IsDelete
	contactAddressModel.ContactID = pl.ContactID
	return contactAddressModel
}

// RestructCorporateDocument function for restruct from cdc
func RestructCorporateDocument(pl B2BDocumentCDC) B2BDocument {
	var documentModule B2BDocument
	documentModule.ID = pl.ID
	documentModule.DocumentType = pl.DocumentType
	documentModule.DocumentFile = pl.DocumentFile
	documentModule.DocumentTitle = pl.DocumentTitle
	documentModule.DocumentDescription = pl.DocumentDescription
	documentModule.NpwpNumber = pl.NpwpNumber
	documentModule.NpwpName = pl.NpwpName
	documentModule.NpwpAddress = pl.NpwpAddress
	documentModule.SiupNumber = pl.SiupNumber
	documentModule.SiupCompanyName = pl.SiupCompanyName
	documentModule.SiupType = pl.SiupType
	documentModule.IsDelete = pl.IsDelete
	documentModule.CreatedAt = pl.CreatedAt
	documentModule.ModifiedAt = pl.ModifiedAt
	documentModule.CreatedBy = pl.CreatedBy
	documentModule.ModifiedBy = pl.ModifiedBy
	documentModule.AccountID = pl.AccountID
	documentModule.IsDisabled = pl.IsDisabled
	documentModule.ProvinceID = pl.ProvinceID
	documentModule.ProvinceName = pl.ProvinceName
	documentModule.CityID = pl.CityID
	documentModule.CityName = pl.CityName
	documentModule.DistrictID = pl.DistrictID
	documentModule.DistrictName = pl.DistrictName
	documentModule.SubdistrictID = pl.SubdistrictID
	documentModule.SubdistrictName = pl.SubdistrictName
	documentModule.PostalCode = pl.PostalCode
	documentModule.Street = pl.Street
	return documentModule
}

// RestructCorporatePhone function for restruct from cdc
func RestructCorporatePhone(pl B2BPhoneCDC) B2BPhone {
	var phoneModel B2BPhone
	phoneModel.ID = pl.ID
	phoneModel.RelationID = pl.RelationID
	phoneModel.RelationType = pl.RelationType
	phoneModel.TypePhone = pl.TypePhone
	phoneModel.Label = pl.Label
	phoneModel.Number = pl.Number
	phoneModel.Area = pl.Area
	phoneModel.Ext = pl.Ext
	phoneModel.IsPrimary = pl.IsPrimary
	phoneModel.IsDelete = pl.IsDelete
	phoneModel.CreatedAt = pl.CreatedAt
	phoneModel.ModifiedAt = pl.ModifiedAt
	phoneModel.CreatedBy = pl.CreatedBy
	phoneModel.ModifiedBy = pl.ModifiedBy
	return phoneModel
}

// RestructCorporateContact function for restruct from cdc
func RestructCorporateContact(pl B2BContactCDC) B2BContact {
	var contactModel B2BContact
	contactModel.ID = pl.ID
	contactModel.ReferenceID = pl.ReferenceID
	contactModel.FirstName = pl.FirstName
	contactModel.LastName = pl.LastName
	contactModel.Salutation = pl.Salutation
	contactModel.JobTitle = pl.JobTitle
	contactModel.Email = pl.Email
	contactModel.NavContactID = pl.NavContactID
	contactModel.IsPrimary = pl.IsPrimary
	contactModel.Note = pl.Note
	contactModel.CreatedAt = pl.CreatedAt
	contactModel.ModifiedAt = pl.ModifiedAt
	contactModel.CreatedBy = pl.CreatedBy
	contactModel.ModifiedBy = pl.ModifiedBy
	contactModel.AccountID = pl.AccountID
	contactModel.IsDisabled = pl.IsDisabled
	contactModel.Password = pl.Password
	contactModel.Avatar = pl.Avatar
	contactModel.Status = pl.Status
	contactModel.Token = pl.Token
	contactModel.IsNew = pl.IsNew
	contactModel.PhoneNumber = pl.PhoneNumber
	contactModel.OtherPhoneNumber = pl.OtherPhoneNumber
	contactModel.Gender = pl.Gender

	if birthDateString, ok := pl.BirthDate.(string); ok {
		if birthdate, err := time.Parse(helper.FormatDateDB, birthDateString); err == nil {
			contactModel.BirthDate = birthdate
		}
	} else if birthDateInteger, okInteger := pl.BirthDate.(int32); okInteger {
		birthdate := helper.DateSinceEpoch(birthDateInteger)
		contactModel.BirthDate = birthdate
	} else if birthDateFloat, okInteger := pl.BirthDate.(float64); okInteger {
		birthDateInteger = int32(birthDateFloat)
		birthdate := helper.DateSinceEpoch(birthDateInteger)
		contactModel.BirthDate = birthdate
	}
	if transactionType, ok := pl.TransactionType.(string); ok {
		m := []TransactionType{}
		if err := json.Unmarshal([]byte(transactionType), &m); err == nil {
			contactModel.TransactionType = m
		}
	}
	contactModel.ErpID = pl.ErpID
	contactModel.IsSync = pl.IsSync
	contactModel.Salt = pl.Salt
	contactModel.LastPasswordModified = pl.LastPasswordModified

	return contactModel
}

// RestructCorporateAddress function for restruct from cdc
func RestructCorporateAddress(pl B2BAddressCDC) B2BAddress {
	var addressModel B2BAddress
	addressModel.ID = pl.ID
	addressModel.TypeAddress = pl.TypeAddress
	addressModel.Name = pl.Name
	addressModel.Street = pl.Street
	addressModel.ProvinceID = pl.ProvinceID
	addressModel.ProvinceName = pl.ProvinceName
	addressModel.CityID = pl.CityID
	addressModel.CityName = pl.CityName
	addressModel.DistrictID = pl.DistrictID
	addressModel.DistrictName = pl.DistrictName
	addressModel.SubDistrictID = pl.SubDistrictID
	addressModel.SubDistrictName = pl.SubDistrictName
	addressModel.IsPrimary = pl.IsPrimary
	addressModel.PostalCode = pl.PostalCode
	addressModel.CreatedAt = pl.CreatedAt
	addressModel.ModifiedAt = pl.ModifiedAt
	addressModel.CreatedBy = pl.CreatedBy
	addressModel.ModifiedBy = pl.ModifiedBy
	addressModel.AccountID = pl.AccountID
	addressModel.IsDisabled = pl.IsDisabled
	addressModel.IsDelete = pl.IsDelete
	return addressModel
}

// RestructCorporateAccount function for restruct from cdc
func RestructCorporateAccount(pl B2BAccountCDC) B2BAccount {
	var accountModel B2BAccount
	accountModel.ID = pl.ID
	accountModel.ParentID = pl.ParentID
	accountModel.NavID = pl.NavID
	accountModel.LegalEntity = pl.LegalEntity
	accountModel.Name = pl.Name
	accountModel.IndustryID = pl.IndustryID
	accountModel.NumberEmployee = pl.NumberEmployee
	accountModel.OfficeEmployee = pl.OfficeEmployee
	accountModel.BusinessSize = pl.BusinessSize
	accountModel.BusinessGroup = pl.BusinessGroup
	accountModel.EstablishedYear = pl.EstablishedYear
	accountModel.IsDelete = pl.IsDelete
	accountModel.TermOfPayment = pl.TermOfPayment
	accountModel.CustomerCategory = pl.CustomerCategory
	accountModel.UserID = pl.UserID
	accountModel.CreatedAt = pl.CreatedAt

	accountModel.CreatedBy = pl.CreatedBy
	accountModel.ModifiedBy = pl.ModifiedBy
	accountModel.AccountGroupID = pl.AccountGroupID
	accountModel.IsDisabled = pl.IsDisabled
	accountModel.Status = pl.Status
	accountModel.PaymentMethodID = pl.PaymentMethodID
	accountModel.PaymentMethodType = pl.PaymentMethodType
	accountModel.SubPaymentMethodName = pl.SubPaymentMethodName
	accountModel.IsCf = pl.IsCf
	accountModel.Logo = pl.Logo
	accountModel.IsParent = pl.IsParent
	accountModel.IsMicrosite = pl.IsMicrosite
	accountModel.MemberType = pl.Membertype
	accountModel.ModifiedAt = pl.ModifiedAt
	if pl.ModifiedAt == nil {
		mod := time.Now().Format(time.RFC3339)
		accountModel.ModifiedAt = &mod
	}
	accountModel.ErpID = pl.ErpID

	return accountModel
}

// RestructCorporateAccountContact function for restruct from cdc
func RestructCorporateAccountContact(pl B2BAccountContactCDC) B2BAccountContact {
	var accountModel B2BAccountContact
	accountModel.ID = pl.ID
	accountModel.Status = pl.Status
	accountModel.IsDelete = pl.IsDelete
	accountModel.IsDisabled = pl.IsDisabled
	accountModel.AccountID = pl.AccountID
	accountModel.ContactID = pl.ContactID
	accountModel.CreatedAt = pl.CreatedAt
	accountModel.ModifiedAt = pl.ModifiedAt
	accountModel.CreatedBy = pl.CreatedBy
	accountModel.ModifiedBy = pl.ModifiedBy
	accountModel.IsAdmin = pl.IsAdmin
	accountModel.DepartmentID = pl.DepartmentID
	return accountModel
}
