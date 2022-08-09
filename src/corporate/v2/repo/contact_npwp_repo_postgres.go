package repo

import (
	"context"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/corporate/v2/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

// ContactNpwpRepoPostgres data structure
type ContactNpwpRepoPostgres struct {
	*repository.Repository
}

// NewContactNpwpRepoPostgres function for initializing ContactNpwp repo
func NewContactNpwpRepoPostgres(repo *repository.Repository) *ContactNpwpRepoPostgres {
	return &ContactNpwpRepoPostgres{repo}
}

// Save function for saving ContactNpwp data
func (mr *ContactNpwpRepoPostgres) Save(ctxReq context.Context, contactNpwp model.B2BContactNpwp) (err error) {
	ctx := "ContactNpwpRepo-create"

	querySave := `INSERT INTO b2b_contact_npwp
				(
					id, number, name, address, address2, 
					contact_id, file, created_at, modified_at, created_by,
					modified_by, is_delete

				)
			VALUES
				(
					$1, $2, $3, $4, $5, $6, $7, $8,
					$9, $10, $11, $12
				)
				ON CONFLICT(id)
				DO UPDATE SET
					id=$1, number=$2, name=$3, address=$4, address2=$5, 
					contact_id=$6, file=$7, created_at=$8, modified_at=$9, created_by=$10,
					modified_by=$11, is_delete=$12`

	tr := tracer.StartTrace(context.Background(), ctx)
	tags := map[string]interface{}{
		helper.TextQuery: querySave,
		helper.TextArgs:  contactNpwp,
	}
	defer tr.Finish(tags)

	if err = repository.Exec(
		mr.Repository, querySave,
		contactNpwp.ID, contactNpwp.Number, contactNpwp.Name, contactNpwp.Address,
		contactNpwp.Address2, contactNpwp.ContactID, contactNpwp.File,
		contactNpwp.CreatedAt,
		contactNpwp.ModifiedAt, contactNpwp.CreatedBy, contactNpwp.ModifiedBy, contactNpwp.IsDelete,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, contactNpwp)
		return err
	}

	return nil
}

// Update function for update Contact Npwp data
func (mr *ContactNpwpRepoPostgres) Update(ctxReq context.Context, contactNpwp model.B2BContactNpwp) (err error) {
	ctx := "ContactNpwpRepo-update"

	queryUpdate := `UPDATE b2b_contact_npwp SET number=$2, name=$3, address=$4, address2=$5, 
	contact_id=$6, file=$7, created_at=$8, modified_at=$9, created_by=$10,
	modified_by=$11, is_delete=$12 WHERE id=$1;`

	tr := tracer.StartTrace(context.Background(), ctx)
	tags := map[string]interface{}{
		helper.TextQuery: queryUpdate,
		helper.TextArgs:  contactNpwp,
	}
	defer tr.Finish(tags)

	if err = repository.Exec(
		mr.Repository, queryUpdate,
		contactNpwp.ID, contactNpwp.Number, contactNpwp.Name, contactNpwp.Address,
		contactNpwp.Address2, contactNpwp.ContactID, contactNpwp.File, contactNpwp.CreatedAt,
		contactNpwp.ModifiedAt, contactNpwp.CreatedBy,
		contactNpwp.ModifiedBy, contactNpwp.IsDelete,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, contactNpwp)
		return err
	}

	return nil
}

// Delete function for delete ContactNpwp data
func (mr *ContactNpwpRepoPostgres) Delete(ctxReq context.Context, contactNpwp model.B2BContactNpwp) (err error) {
	ctx := "ContactNpwpRepo-delete"

	tr := tracer.StartTrace(context.Background(), ctx)
	tags := map[string]interface{}{
		helper.TextArgs: contactNpwp,
		helper.TagCtx:   ctx,
	}
	defer tr.Finish(tags)

	if err = repository.DeleteByID(mr.Repository, contactNpwp.ID, "b2b_contact_npwp"); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, contactNpwp)
		return err
	}

	return nil
}
