package model

import (
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
)

// B2BAccountTemporary data structure
type B2BAccountTemporary struct {
	ID               int     `json:"id"`
	IndustryID       *string `json:"industryId"`
	NumberEmployee   *string `json:"numberEmployee"`
	OfficeEmployee   *int64  `json:"officeEmployee"`
	BusinessSize     *string `json:"businessSize"`
	OrganizationName *string `json:"organizationName"`
	LegalEntity      *string `json:"legalEntity"`
	Name             *string `json:"name"`
	Salutation       *string `json:"salutation"`
	Email            *string `json:"email"`
	Password         *string `json:"password"`
	Token            *string `json:"token"`
	IsDelete         bool    `json:"isDelete"`
	IsDisabled       bool    `json:"isDisabled"`
	Phone            *string `json:"phone"`
	ParentID         *string `json:"parentId"`
	NpwpNumber       *string `json:"npwpNumber"`
}

// RestructCorporateAccountTemporary function for restruct from cdc
func RestructCorporateAccountTemporary(pl serviceModel.AccountTemporaryPayloadCDC) B2BAccountTemporary {
	var accountModel B2BAccountTemporary
	accountModel.ID = pl.Payload.After.ID
	accountModel.IndustryID = pl.Payload.After.IndustryID
	accountModel.NumberEmployee = pl.Payload.After.NumberEmployee
	accountModel.OfficeEmployee = pl.Payload.After.OfficeEmployee
	accountModel.BusinessSize = pl.Payload.After.BusinessSize
	accountModel.OrganizationName = pl.Payload.After.OrganizationName
	accountModel.LegalEntity = pl.Payload.After.LegalEntity
	accountModel.Name = pl.Payload.After.Name
	accountModel.Salutation = pl.Payload.After.Salutation
	accountModel.Email = pl.Payload.After.Email
	accountModel.Password = pl.Payload.After.Password
	accountModel.Token = pl.Payload.After.Token
	accountModel.IsDelete = pl.Payload.After.IsDelete
	accountModel.IsDisabled = pl.Payload.After.IsDisabled
	accountModel.Phone = pl.Payload.After.Phone
	accountModel.ParentID = pl.Payload.After.ParentID
	accountModel.NpwpNumber = pl.Payload.After.NpwpNumber
	return accountModel
}
