package repo

import (
	"context"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

// PhoneRepoPostgres data structure
type PhoneRepoPostgres struct {
	*repository.Repository
}

// NewPhoneRepoPostgres function for initializing  repo
func NewPhoneRepoPostgres(repo *repository.Repository) *PhoneRepoPostgres {
	return &PhoneRepoPostgres{repo}
}

// Save function for saving  data
func (mr *PhoneRepoPostgres) Save(ctxReq context.Context, phone sharedModel.B2BPhone) (err error) {
	ctx := "PhoneRepo-create"

	querySave := `INSERT INTO b2b_phone
				(
					id, relation_id, relation_type, type_phone, label,
					number, area, ext, is_primary,
					is_delete, created_at, modified_at, created_by, 
					modified_by
				)
			VALUES
				(
					$1, $2, $3, $4, $5, $6, 
					$7, $8, $9, $10, $11, $12,
					$13, $14
				)
			ON CONFLICT(id)
			DO UPDATE SET
			id=$1, relation_id=$2, relation_type=$3, type_phone=$4, label=$5,
			number=$6, area=$7, ext=$8, is_primary=$9,
			is_delete=$10, created_at=$11, modified_at=$12, created_by=$13, modified_by=$14`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextQuery: querySave,
		helper.TextArgs:  phone,
	}
	defer tr.Finish(tags)

	if err = repository.Exec(
		mr.Repository, querySave,
		phone.ID, phone.RelationID, phone.RelationType, phone.TypePhone,
		phone.Label, phone.Number, phone.Area, phone.Ext,
		phone.IsPrimary, phone.IsDelete,
		phone.CreatedAt, phone.ModifiedAt,
		phone.CreatedBy, phone.ModifiedBy,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, phone)
		return err
	}

	return nil
}

// Update function for update Phone data
func (mr *PhoneRepoPostgres) Update(ctxReq context.Context, phone sharedModel.B2BPhone) (err error) {
	ctx := "PhoneRepo-update"

	queryUpdate := `UPDATE b2b_phone SET relation_id=$2, relation_type=$3, type_phone=$4, label=$5,
	number=$6, area=$7, ext=$8, is_primary=$9,
	is_delete=$10, created_at=$11, modified_at=$12, created_by=$13, modified_by=$14 WHERE id=$1;`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextArgs:  phone,
		helper.TagCtx:    ctx,
		helper.TextQuery: queryUpdate,
	}
	defer tr.Finish(tags)

	if err = repository.Exec(
		mr.Repository, queryUpdate,
		phone.ID, phone.RelationID, phone.RelationType,
		phone.TypePhone,
		phone.Label, phone.Number, phone.Area, phone.Ext,
		phone.IsPrimary, phone.IsDelete, phone.CreatedAt,
		phone.ModifiedAt,
		phone.CreatedBy, phone.ModifiedBy,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, phone)
		return err
	}

	return nil
}

// Delete function for delete Phone data
func (mr *PhoneRepoPostgres) Delete(ctxReq context.Context, phone sharedModel.B2BPhone) (err error) {
	ctx := "PhoneRepo-delete"
	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextArgs: phone,
		helper.TagCtx:   ctx,
	}
	defer tr.Finish(tags)

	if err = repository.DeleteByID(mr.Repository, phone.ID, "b2b_phone"); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, phone)
		return err
	}

	return nil
}
