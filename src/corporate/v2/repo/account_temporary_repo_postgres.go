package repo

import (
	"context"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/corporate/v2/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

// AccountTemporaryRepoPostgres data structure
type AccountTemporaryRepoPostgres struct {
	*repository.Repository
}

// NewAccountTemporaryRepoPostgres function for initializing Account repo
func NewAccountTemporaryRepoPostgres(repo *repository.Repository) *AccountTemporaryRepoPostgres {
	return &AccountTemporaryRepoPostgres{repo}
}

// Save function for saving account data
func (mr *AccountTemporaryRepoPostgres) Save(ctxReq context.Context, accountTemp model.B2BAccountTemporary) (err error) {
	ctx := "AccountTemporaryRepo-create"

	querySave := `INSERT INTO b2b_account_temporary
				(
					id, industry_id, number_employee, office_employee,
					business_size, organization_name, legal_entity,
					name, salutation, email, password, token,
					is_delete, is_disabled, phone, parent_id, npwp_number
				)
			VALUES
				(
					$1, $2, $3, $4, $5, $6, $7, $8,
					$9, $10, $11, $12, $13, $14, 
					$15, $16, $17
				)
			ON CONFLICT(id)
			DO UPDATE SET
			id = $1, industry_id = $2, number_employee = $3, office_employee = $4,
			business_size = $5, organization_name = $6, legal_entity = $7,
			name = $8, salutation = $9, email = $10, password = $11, token = $12,
			is_delete = $13, is_disabled = $14, phone = $15, parent_id = $16, npwp_number = $17`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextArgs:  accountTemp,
		helper.TextQuery: querySave,
	}
	defer tr.Finish(tags)
	if err = repository.Exec(
		mr.Repository, querySave,
		accountTemp.ID, accountTemp.IndustryID, accountTemp.NumberEmployee,
		accountTemp.OfficeEmployee, accountTemp.BusinessSize,
		accountTemp.OrganizationName,
		accountTemp.LegalEntity, accountTemp.Name,
		accountTemp.Salutation,
		accountTemp.Email, accountTemp.Password, accountTemp.Token,
		accountTemp.IsDelete, accountTemp.IsDisabled, accountTemp.Phone,
		accountTemp.ParentID, accountTemp.NpwpNumber,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, accountTemp)
		return err
	}

	return nil
}

// Update function for update account data
func (mr *AccountTemporaryRepoPostgres) Update(ctxReq context.Context, accountTemp model.B2BAccountTemporary) (err error) {
	ctx := "AccountTemporaryRepo-update"

	queryUpdate := `UPDATE b2b_account_temporary SET industry_id = $2, number_employee = $3, office_employee = $4,
	business_size = $5, organization_name = $6, legal_entity = $7,
	name = $8, salutation = $9, email = $10, password = $11, token = $12,
	is_delete = $13, is_disabled = $14, phone = $15, parent_id = $16, npwp_number = $17 WHERE id = $1;`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextArgs:  accountTemp,
		helper.TextQuery: queryUpdate,
	}
	defer tr.Finish(tags)

	if err = repository.Exec(
		mr.Repository, queryUpdate,
		accountTemp.ID, accountTemp.IndustryID, accountTemp.NumberEmployee, accountTemp.OfficeEmployee, accountTemp.BusinessSize, accountTemp.OrganizationName,
		accountTemp.LegalEntity, accountTemp.Name, accountTemp.Salutation,
		accountTemp.Email, accountTemp.Password, accountTemp.Token,
		accountTemp.IsDelete, accountTemp.IsDisabled, accountTemp.Phone,
		accountTemp.ParentID, accountTemp.NpwpNumber,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, accountTemp)
		return err
	}

	return nil
}

// Delete function for delete account data
func (mr *AccountTemporaryRepoPostgres) Delete(ctxReq context.Context, accountTemp model.B2BAccountTemporary) (err error) {
	ctx := "AccountTemporaryRepo-delete"

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextArgs: accountTemp,
		helper.TagCtx:   ctx,
	}
	defer tr.Finish(tags)

	if err := repository.DeleteByID(mr.Repository, accountTemp.ID, "b2b_account_temporary"); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, accountTemp)
		return err
	}

	return nil
}
