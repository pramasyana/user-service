package repo

import (
	"context"
	"strconv"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

// AccountRepoPostgres data structure
type AccountRepoPostgres struct {
	*repository.Repository
}

// NewAccountRepoPostgres function for initializing Account repo
func NewAccountRepoPostgres(repo *repository.Repository) *AccountRepoPostgres {
	return &AccountRepoPostgres{repo}
}

// Save function for saving account data
func (mr *AccountRepoPostgres) Save(ctxReq context.Context, account sharedModel.B2BAccount) (err error) {
	ctx := "AccountRepo-create"

	querySave := `INSERT INTO b2b_account
				(
					id, parent_id, nav_id, legal_entity, name, industry_id, number_employee, office_employee,
					business_size, business_group, established_year,
					is_delete, term_of_payment, customer_category, user_id, created_at,
					modified_at, created_by, modified_by, accountgroup_id,
					is_disabled, status, payment_method_id, payment_method_type,
					sub_payment_method_name, is_cf, logo, is_parent, is_microsite, member_type, erp_id
				)
			VALUES
				(
					$1, $2, $3, $4, $5, $6, $7, $8,
					$9, $10, $11,
					$12, $13, $14, $15, $16, $17,
					$18, $19, $20, $21,
					$22, $23,
					$24, $25, $26,
					$27, $28, $29, $30, $31
				)
				ON CONFLICT(id)
				DO UPDATE SET
					id=$1, parent_id=$2, nav_id=$3, legal_entity=$4, name=$5, industry_id=$6, number_employee=$7, office_employee=$8,
					business_size=$9, business_group=$10, established_year=$11,
					is_delete=$12, term_of_payment=$13, customer_category=$14, user_id=$15, created_at=$16,
					modified_at=$17, created_by=$18, modified_by=$19, accountgroup_id=$20,
					is_disabled=$21, status=$22, payment_method_id=$23, payment_method_type=$24,
					sub_payment_method_name=$25, is_cf=$26, logo=$27, is_parent=$28, is_microsite=$29, member_type=$30, erp_id=$31`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextQuery: querySave,
		helper.TextArgs:  account,
	}
	defer tr.Finish(tags)

	if err = repository.Exec(
		mr.Repository, querySave,
		account.ID, account.ParentID, account.NavID, account.LegalEntity,
		account.Name, account.IndustryID, account.NumberEmployee,
		account.OfficeEmployee, account.BusinessSize, account.BusinessGroup,
		account.EstablishedYear, account.IsDelete, account.TermOfPayment,
		account.CustomerCategory, account.UserID, account.CreatedAt,
		account.ModifiedAt, account.CreatedBy, account.ModifiedBy, account.AccountGroupID,
		account.IsDisabled, account.Status, account.PaymentMethodID, account.PaymentMethodType,
		account.SubPaymentMethodName, account.IsCf, account.Logo, account.IsParent, account.IsMicrosite, account.MemberType, account.ErpID,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, account)
		return err
	}
	return nil
}

// Update function for update account data
func (mr *AccountRepoPostgres) Update(ctxReq context.Context, account sharedModel.B2BAccount) (err error) {
	ctx := "AccountRepo-update"

	queryUpdate := `UPDATE b2b_account SET parent_id=$2, nav_id=$3, legal_entity=$4, name=$5, industry_id=$6, number_employee=$7, office_employee=$8,
	business_size=$9, business_group=$10, established_year=$11,
	is_delete=$12, term_of_payment=$13, customer_category=$14, user_id=$15, created_at=$16,
	modified_at=$17, created_by=$18, modified_by=$19, accountgroup_id=$20,
	is_disabled=$21, status=$22, payment_method_id=$23, payment_method_type=$24,
	sub_payment_method_name=$25, is_cf=$26, logo=$27, is_parent=$28, is_microsite=$29, member_type=$30, erp_id=$31 WHERE id=$1;`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextQuery: queryUpdate,
		helper.TextArgs:  account,
	}
	defer tr.Finish(tags)

	if err = repository.Exec(
		mr.Repository, queryUpdate,
		account.ID, account.ParentID, account.NavID, account.LegalEntity,
		account.Name, account.IndustryID, account.NumberEmployee,
		account.OfficeEmployee, account.BusinessSize, account.BusinessGroup,
		account.EstablishedYear, account.IsDelete, account.TermOfPayment,
		account.CustomerCategory, account.UserID, account.CreatedAt,
		account.ModifiedAt, account.CreatedBy, account.ModifiedBy,
		account.AccountGroupID,
		account.IsDisabled, account.Status, account.PaymentMethodID,
		account.PaymentMethodType,
		account.SubPaymentMethodName, account.IsCf, account.Logo, account.IsParent, account.IsMicrosite, account.MemberType, account.ErpID,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, account)
		return err
	}

	return nil
}

// Delete function for delete account data
func (mr *AccountRepoPostgres) Delete(ctxReq context.Context, account sharedModel.B2BAccount) (err error) {
	ctx := "AccountRepo-delete"

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TagCtx:   ctx,
		helper.TextArgs: account,
	}
	defer tr.Finish(tags)

	num, err := strconv.Atoi(account.ID)
	if err != nil {
		return err
	}

	if err = repository.DeleteByID(mr.Repository, num, "b2b_account"); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, account)
		return err
	}

	return nil
}
