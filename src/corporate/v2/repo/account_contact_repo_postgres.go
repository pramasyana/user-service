package repo

import (
	"context"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

// AccountContactRepoPostgres data structure
type AccountContactRepoPostgres struct {
	*repository.Repository
}

// NewAccountContactRepoPostgres function for initializing Account repo
func NewAccountContactRepoPostgres(repo *repository.Repository) *AccountContactRepoPostgres {
	return &AccountContactRepoPostgres{repo}
}

// Save function for saving account data
func (mr *AccountContactRepoPostgres) Save(ctxReq context.Context, accountContact sharedModel.B2BAccountContact) (err error) {
	ctx := "AccountContactRepo-create"

	querySave := `INSERT INTO b2b_account_contact
				(
					id, status, is_delete, is_disabled,
					account_id, contact_id, created_at,
					modified_at, created_by, modified_by, is_admin, department_id
				)
			VALUES
				(
					$1, $2, $3, 
					$4, $5, $6, 
					$7, $8, $9, $10, $11, $12
				)
			ON CONFLICT(id)
			DO UPDATE SET
			id=$1, status=$2, is_delete=$3, is_disabled=$4,
			account_id=$5, contact_id=$6, created_at=$7,
			modified_at=$8, created_by=$9, modified_by=$10, is_admin=$11, department_id=$12`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextQuery: querySave,
		helper.TextArgs:  accountContact,
	}
	defer tr.Finish(tags)

	if err = repository.Exec(
		mr.Repository, querySave,
		accountContact.ID, accountContact.Status, accountContact.IsDelete, accountContact.IsDisabled,
		accountContact.AccountID, accountContact.ContactID, accountContact.CreatedAt,
		accountContact.ModifiedAt, accountContact.CreatedBy, accountContact.ModifiedBy, accountContact.IsAdmin, accountContact.DepartmentID,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextQuery, err, accountContact)
		return err
	}

	return nil
}

// Update function for update account data
func (mr *AccountContactRepoPostgres) Update(ctxReq context.Context, accountContact sharedModel.B2BAccountContact) (err error) {
	ctx := "AccountContactRepo-update"

	query := `UPDATE b2b_account_contact SET status=$2, is_delete=$3, is_disabled=$4,
	account_id=$5, contact_id=$6, created_at=$7,
	modified_at=$8, created_by=$9, modified_by=$10, is_admin=$11, department_id=$12 WHERE id=$1;`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextArgs:  accountContact,
		helper.TextQuery: query,
	}
	defer tr.Finish(tags)

	if err = repository.Exec(
		mr.Repository, query,
		accountContact.ID, accountContact.Status, accountContact.IsDelete,
		accountContact.IsDisabled, accountContact.AccountID, accountContact.ContactID,
		accountContact.CreatedAt, accountContact.ModifiedAt,
		accountContact.CreatedBy,
		accountContact.ModifiedBy, accountContact.IsAdmin, accountContact.DepartmentID,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextQuery, err, accountContact)
		return err
	}

	return nil
}

// Delete function for delete account data
func (mr *AccountContactRepoPostgres) Delete(ctxReq context.Context, accountContact sharedModel.B2BAccountContact) (err error) {
	ctx := "AccountContactRepo-delete"

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextArgs: accountContact,
		helper.TagCtx:   ctx,
	}
	defer tr.Finish(tags)

	if err = repository.DeleteByID(mr.Repository, accountContact.ID, "b2b_account_contact"); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextQuery, err, accountContact)
		return err
	}

	return nil
}
