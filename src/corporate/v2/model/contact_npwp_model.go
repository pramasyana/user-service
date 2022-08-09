package model

import serviceModel "github.com/Bhinneka/user-service/src/service/model"

// B2BContactNpwp data structure
type B2BContactNpwp struct {
	ID         int     `json:"id"`
	Number     *string `json:"number"`
	Name       *string `json:"name"`
	Address    *string `json:"address"`
	Address2   *string `json:"address2"`
	ContactID  *int    `json:"contactId"`
	File       *string `json:"file"`
	CreatedAt  *string `json:"createdAt"`
	ModifiedAt *string `json:"modifiedAt"`
	CreatedBy  *int64  `json:"createdBy"`
	ModifiedBy *int64  `json:"modifiedBy"`
	IsDelete   bool    `json:"isDelete"`
}

// RestructCorporateContactNpwp function for restruct from cdc
func RestructCorporateContactNpwp(pl serviceModel.ContactNpwpPayloadCDC) B2BContactNpwp {
	var contactNpwpModule B2BContactNpwp
	contactNpwpModule.ID = pl.Payload.After.ID
	contactNpwpModule.Number = pl.Payload.After.Number
	contactNpwpModule.Name = pl.Payload.After.Name
	contactNpwpModule.Address = pl.Payload.After.Address
	contactNpwpModule.Address2 = pl.Payload.After.Address2
	contactNpwpModule.ContactID = pl.Payload.After.ContactID
	contactNpwpModule.File = pl.Payload.After.File
	contactNpwpModule.CreatedAt = pl.Payload.After.CreatedAt
	contactNpwpModule.ModifiedAt = pl.Payload.After.ModifiedAt
	contactNpwpModule.CreatedBy = pl.Payload.After.CreatedBy
	contactNpwpModule.ModifiedBy = pl.Payload.After.ModifiedBy
	contactNpwpModule.IsDelete = pl.Payload.After.IsDelete
	return contactNpwpModule
}
