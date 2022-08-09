package repo

import (
	"context"
	"encoding/json"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
	"github.com/lib/pq"
)

// ContactRepoPostgres data structure
type ContactRepoPostgres struct {
	*repository.Repository
}

// NewContactRepoPostgres function for initializing  repo
func NewContactRepoPostgres(repo *repository.Repository) *ContactRepoPostgres {
	return &ContactRepoPostgres{repo}
}

// Save function for saving  data
func (mr *ContactRepoPostgres) Save(ctxReq context.Context, contact sharedModel.B2BContact) (err error) {
	ctx := "ContactRepo-create"

	querySave := `INSERT INTO b2b_contact
				(
					id, reference_id, first_name, last_name, salutation,
					job_title, email, nav_contact_id, is_primary,
					birth_date, note, created_at, modified_at, created_by, 
					modified_by, account_id, is_disabled, password, avatar, status,
					token, is_new, phone_number, other_phone_number, gender, transaction_type, erp_id, 
					is_sync, salt, last_password_modified
				)
			VALUES
				(
					$1, $2, $3, $4, $5, $6, 
					$7, $8, $9, $10, $11, $12,
					$13, $14, $15, $16, $17, $18,
					$19, $20, $21, $22, $23, $24, $25, $26, $27,
					$28, $29, $30
				)
			ON CONFLICT(id)
			DO UPDATE SET
			id = $1, reference_id = $2, first_name = $3, last_name = $4, salutation = $5,
			job_title = $6, email = $7, nav_contact_id = $8, is_primary = $9,
			birth_date = $10, note = $11, created_at = $12, modified_at = $13, created_by = $14, 
			modified_by = $15, account_id = $16, is_disabled = $17, password = $18, avatar = $19, status = $20,
			token = $21, is_new = $22, phone_number = $23, other_phone_number = $24, gender=$25, transaction_type=$26::jsonb, erp_id=$27, 
			is_sync=$28, salt=$29, last_password_modified=$30`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextQuery: querySave,
		helper.TextArgs:  contact,
	}
	defer tr.Finish(tags)

	var birthDate pq.NullTime

	if !contact.BirthDate.IsZero() {
		birthDate.Time = contact.BirthDate
		birthDate.Valid = true
	}
	jsonTransactionType, _ := json.Marshal(contact.TransactionType)

	if err = repository.Exec(
		mr.Repository, querySave,
		contact.ID, contact.ReferenceID, contact.FirstName, contact.LastName,
		contact.Salutation, contact.JobTitle, contact.Email, contact.NavContactID,
		contact.IsPrimary, birthDate,
		contact.Note, contact.CreatedAt,
		contact.ModifiedAt, contact.CreatedBy, contact.ModifiedBy,
		contact.AccountID, contact.IsDisabled, contact.Password, contact.Avatar,
		contact.Status, contact.Token, contact.IsNew, contact.PhoneNumber, contact.OtherPhoneNumber, contact.Gender,
		jsonTransactionType, contact.ErpID, contact.IsSync, contact.Salt, contact.LastPasswordModified,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, contact)
		return err
	}

	return nil
}

// Update function for update contact data
func (mr *ContactRepoPostgres) Update(ctxReq context.Context, contact sharedModel.B2BContact) (err error) {
	ctx := "ContactRepo-update"

	queryUpdate := `
	UPDATE b2b_contact SET reference_id = $2, first_name = $3, last_name = $4, salutation = $5,
	job_title = $6, email = $7, nav_contact_id = $8, is_primary = $9,
	birth_date = $10, note = $11, created_at = $12, modified_at = $13, created_by = $14, 
	modified_by = $15, account_id = $16, is_disabled = $17, password = $18, avatar = $19, status = $20,
	token = $21, is_new = $22, phone_number = $23, other_phone_number = $24, gender = $25, transaction_type = $26::jsonb, erp_id=$27
	WHERE id = $1;`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextQuery: queryUpdate,
		helper.TextArgs:  contact,
	}
	defer tr.Finish(tags)

	var birthDate pq.NullTime

	if !contact.BirthDate.IsZero() {
		birthDate.Time = contact.BirthDate
		birthDate.Valid = true
	}
	jsonTransactionType, _ := json.Marshal(contact.TransactionType)

	if err = repository.Exec(
		mr.Repository, queryUpdate,
		contact.ID, contact.ReferenceID, contact.FirstName, contact.LastName,
		contact.Salutation, contact.JobTitle, contact.Email,
		contact.NavContactID,
		contact.IsPrimary, birthDate, contact.Note, contact.CreatedAt,
		contact.ModifiedAt, contact.CreatedBy, contact.ModifiedBy,
		contact.AccountID, contact.IsDisabled,
		contact.Password, contact.Avatar,
		contact.Status, contact.Token, contact.IsNew, contact.PhoneNumber,
		contact.OtherPhoneNumber, contact.Gender, jsonTransactionType, contact.ErpID,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, contact)
		return err
	}

	return nil
}

// Delete function for delete contact data
func (mr *ContactRepoPostgres) Delete(ctxReq context.Context, contact sharedModel.B2BContact) (err error) {
	ctx := "ContactRepo-delete"

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextArgs: contact,
		helper.TagCtx:   ctx,
	}
	defer tr.Finish(tags)

	if err = repository.DeleteByID(mr.Repository, contact.ID, "b2b_contact"); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, contact)
		return err
	}

	return nil
}
