package repo

import (
	"context"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

// LeadsRepoPostgres data structure
type LeadsRepoPostgres struct {
	*repository.Repository
}

// NewLeadsRepoPostgres function for initializing  repo
func NewLeadsRepoPostgres(repo *repository.Repository) *LeadsRepoPostgres {
	return &LeadsRepoPostgres{repo}
}

// Save function for saving  data
func (mr *LeadsRepoPostgres) Save(ctxReq context.Context, leadsTemp sharedModel.B2BLeads) (err error) {
	ctx := "LeadsRepo-create"

	querySave := `INSERT INTO b2b_leads
				(
					id, name, email, phone,
					source_id, source, source_status, status,
					notes, created_at, modified_at, created_by, 
					modified_by
				)
			VALUES
				(
					$1, $2, $3, $4, $5, $6, 
					$7, $8, $9, $10, $11, $12, $13
				)
			ON CONFLICT(id)
			DO UPDATE SET
			id=$1, name=$2, email=$3, phone=$4,
			source_id=$5, source=$6, source_status=$7, status=$8,
			notes=$9, created_at=$10, modified_at=$11, created_by=$12, 
			modified_by=$13`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextQuery: querySave,
		helper.TextArgs:  leadsTemp,
	}
	defer tr.Finish(tags)

	if err = repository.Exec(
		mr.Repository, querySave,
		leadsTemp.ID, leadsTemp.Name, leadsTemp.Email,
		leadsTemp.Phone, leadsTemp.SourceID, leadsTemp.Source,
		leadsTemp.SourceStatus,
		leadsTemp.Status, leadsTemp.Notes, leadsTemp.CreatedAt, leadsTemp.ModifiedAt,
		leadsTemp.CreatedBy, leadsTemp.ModifiedBy,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, leadsTemp)
		return err
	}

	return nil
}

// Update function for update leads data
func (mr *LeadsRepoPostgres) Update(ctxReq context.Context, leadsTemp sharedModel.B2BLeads) (err error) {
	ctx := "LeadsRepo-update"

	queryUpdate := `UPDATE b2b_leads SET name=$2, email=$3, phone=$4,
	source_id=$5, source=$6, source_status=$7, status=$8,
	notes=$9, created_at=$10, modified_at=$11, created_by=$12, 
	modified_by=$13 WHERE id=$1;`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextQuery: queryUpdate,
		helper.TextArgs:  leadsTemp,
	}
	defer tr.Finish(tags)

	if err = repository.Exec(
		mr.Repository, queryUpdate,
		leadsTemp.ID, leadsTemp.Name, leadsTemp.Email,
		leadsTemp.Phone, leadsTemp.SourceID, leadsTemp.Source, leadsTemp.SourceStatus,
		leadsTemp.Status, leadsTemp.Notes, leadsTemp.CreatedAt, leadsTemp.ModifiedAt,
		leadsTemp.CreatedBy, leadsTemp.ModifiedBy,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, leadsTemp)
		return err
	}

	return nil
}

// Delete function for delete leads data
func (mr *LeadsRepoPostgres) Delete(ctxReq context.Context, leadsTemp sharedModel.B2BLeads) (err error) {
	ctx := "LeadsRepo-delete"

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextArgs: leadsTemp,
		helper.TagCtx:   ctx,
	}
	defer tr.Finish(tags)

	if err = repository.DeleteByID(mr.Repository, leadsTemp.ID, "b2b_leads"); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, leadsTemp)
		return err
	}

	return nil
}
