package repo

import (
	"context"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/corporate/v2/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
	"github.com/lib/pq"
)

// ContactTempRepoPostgres data structure
type ContactTempRepoPostgres struct {
	*repository.Repository
}

// NewContactTempRepoPostgres function for initializing  repo
func NewContactTempRepoPostgres(repo *repository.Repository) *ContactTempRepoPostgres {
	return &ContactTempRepoPostgres{repo}
}

// Save function for saving  data
func (mr *ContactTempRepoPostgres) Save(ctxReq context.Context, contactTemp model.B2BContactTemp) (err error) {
	ctx := "ContactTempRepo-create"

	querySave := `INSERT INTO b2b_contact_temp
				(
					id, first_name, last_name, salutation,
					job_title, email, nav_contact_id, is_primary,
					birth_date, note, created_at, modified_at, created_by, 
					modified_by, account_id, reference_id
				)
			VALUES
				(
					$1, $2, $3, $4, $5, $6, 
					$7, $8, $9, $10, $11, $12,
					$13, $14, $15, $16
				)
			ON CONFLICT(id)
			DO UPDATE SET
			id=$1, first_name=$2, last_name=$3, salutation=$4,
			job_title=$5, email=$6, nav_contact_id=$7, is_primary=$8,
			birth_date=$9, note=$10, created_at=$11, modified_at=$12, created_by=$13, 
			modified_by=$14, account_id=$15, reference_id=$16`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextQuery: querySave,
		helper.TextArgs:  contactTemp,
	}
	defer tr.Finish(tags)

	var birthDate pq.NullTime

	if !contactTemp.BirthDate.IsZero() {
		birthDate.Time = contactTemp.BirthDate
	}

	if err = repository.Exec(
		mr.Repository, querySave,
		contactTemp.ID, contactTemp.FirstName, contactTemp.LastName,
		contactTemp.Salutation, contactTemp.JobTitle, contactTemp.Email,
		contactTemp.NavContactID, contactTemp.IsPrimary,
		birthDate, contactTemp.Note, contactTemp.CreatedAt,
		contactTemp.ModifiedAt, contactTemp.CreatedBy, contactTemp.ModifiedBy,
		contactTemp.AccountID, contactTemp.ReferenceID,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, contactTemp)
		return err
	}

	return nil
}

// Update function for update contact temp data
func (mr *ContactTempRepoPostgres) Update(ctxReq context.Context, contactTemp model.B2BContactTemp) (err error) {
	ctx := "ContactTempRepo-update"

	queryUpdate := `UPDATE b2b_contact_temp SET first_name=$2, last_name=$3, salutation=$4,
	job_title=$5, email=$6, nav_contact_id=$7, is_primary=$8,
	birth_date=$9, note=$10, created_at=$11, modified_at=$12, created_by=$13, 
	modified_by=$14, account_id=$15, reference_id=$16 WHERE id=$1;`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextQuery: queryUpdate,
		helper.TextArgs:  contactTemp,
	}
	defer tr.Finish(tags)

	var birthDate pq.NullTime

	if !contactTemp.BirthDate.IsZero() {
		birthDate.Time = contactTemp.BirthDate
	}

	if err = repository.Exec(
		mr.Repository, queryUpdate,
		contactTemp.ID, contactTemp.FirstName, contactTemp.LastName,
		contactTemp.Salutation, contactTemp.JobTitle, contactTemp.Email,
		contactTemp.NavContactID, contactTemp.IsPrimary, birthDate,
		contactTemp.Note, contactTemp.CreatedAt,
		contactTemp.ModifiedAt, contactTemp.CreatedBy, contactTemp.ModifiedBy,
		contactTemp.AccountID, contactTemp.ReferenceID,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, contactTemp)
		return err
	}

	return nil
}

// Delete function for delete contact temp data
func (mr *ContactTempRepoPostgres) Delete(ctxReq context.Context, contactTemp model.B2BContactTemp) (err error) {
	ctx := "ContactTempRepo-delete"

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TagCtx:   ctx,
		helper.TextArgs: contactTemp,
	}
	defer tr.Finish(tags)

	if err = repository.DeleteByID(mr.Repository, contactTemp.ID, "b2b_contact_temp"); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, contactTemp)
		return err
	}

	return nil
}
