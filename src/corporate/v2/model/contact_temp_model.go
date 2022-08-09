package model

import (
	"time"

	"github.com/Bhinneka/user-service/helper"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
)

// B2BContactTemp data structure
type B2BContactTemp struct {
	ID           int       `json:"id"`
	FirstName    *string   `json:"firstName"`
	LastName     *string   `json:"lastName"`
	Salutation   *string   `json:"salutation"`
	JobTitle     *string   `json:"jobTitle"`
	Email        *string   `json:"email"`
	NavContactID *string   `json:"navContactId"`
	IsPrimary    bool      `json:"isPrimary"`
	BirthDate    time.Time `json:"birthDate"`
	Note         *string   `json:"note"`
	CreatedAt    *string   `json:"createdAt"`
	ModifiedAt   *string   `json:"modifiedAt"`
	CreatedBy    *int64    `json:"createdBy"`
	ModifiedBy   *int64    `json:"modifiedBy"`
	AccountID    *string   `json:"accountId"`
	ReferenceID  *string   `json:"referenceId"`
}

// RestructCorporateContactTemp function for restruct from cdc
func RestructCorporateContactTemp(pl serviceModel.ContactTempPayloadCDC) B2BContactTemp {
	var contactTempModel B2BContactTemp
	contactTempModel.ID = pl.Payload.After.ID
	contactTempModel.FirstName = pl.Payload.After.FirstName
	contactTempModel.LastName = pl.Payload.After.LastName
	contactTempModel.Salutation = pl.Payload.After.Salutation
	contactTempModel.JobTitle = pl.Payload.After.JobTitle
	contactTempModel.Email = pl.Payload.After.Email
	contactTempModel.NavContactID = pl.Payload.After.NavContactID
	contactTempModel.IsPrimary = pl.Payload.After.IsPrimary

	contactTempModel.Note = pl.Payload.After.Note
	contactTempModel.CreatedAt = pl.Payload.After.CreatedAt
	contactTempModel.ModifiedAt = pl.Payload.After.ModifiedAt
	contactTempModel.CreatedBy = pl.Payload.After.CreatedBy
	contactTempModel.ModifiedBy = pl.Payload.After.ModifiedBy
	contactTempModel.AccountID = pl.Payload.After.AccountID
	contactTempModel.ReferenceID = pl.Payload.After.ReferenceID
	if birthDateString, ok := pl.Payload.After.BirthDate.(string); ok {
		if birthdate, err := time.Parse(helper.FormatDateDB, birthDateString); err == nil {
			contactTempModel.BirthDate = birthdate
		}
	} else if birthDateInteger, okInteger := pl.Payload.After.BirthDate.(int32); okInteger {
		birthdate := helper.DateSinceEpoch(birthDateInteger)
		contactTempModel.BirthDate = birthdate
	} else if birthDateFloat, okInteger := pl.Payload.After.BirthDate.(float64); okInteger {
		birthDateInteger = int32(birthDateFloat)
		birthdate := helper.DateSinceEpoch(birthDateInteger)
		contactTempModel.BirthDate = birthdate
	}

	return contactTempModel
}
