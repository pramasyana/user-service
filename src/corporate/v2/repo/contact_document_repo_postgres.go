package repo

import (
	"context"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

// ContactDocumentRepoPostgres data structure
type ContactDocumentRepoPostgres struct {
	*repository.Repository
}

// NewContactDocumentRepoPostgres function for initializing Contact Document repo
func NewContactDocumentRepoPostgres(repo *repository.Repository) *ContactDocumentRepoPostgres {
	return &ContactDocumentRepoPostgres{repo}
}

// Save function for saving Contact Document data
func (mr *ContactDocumentRepoPostgres) Save(ctxReq context.Context, cd sharedModel.B2BContactDocument) (err error) {
	ctx := "ContactDocumentRepo-Save"

	querySaveCD := `INSERT INTO b2b_contact_document
				(
					id, document_file, document_type, document_title, document_description, 
					npwp_number, npwp_name, npwp_address, npwp_address2, siup_number, siup_company_name,
					siup_type, created_at, modified_at, created_by, modified_by,
					is_delete, contact_id
				)
			VALUES
				(
					$1, $2, $3, $4, $5, $6, $7, $8,
					$9, $10, $11, $12, $13, $14, 
					$15, $16, $17, $18
				)
				ON CONFLICT(id)
				DO UPDATE SET
					id=$1, document_file=$2, document_type=$3, document_title=$4, document_description=$5, 
					npwp_number=$6, npwp_name=$7, npwp_address=$8, npwp_address2=$9, siup_number=$10, siup_company_name=$11,
					siup_type=$12, created_at=$13, modified_at=$14, created_by=$15, modified_by=$16, 
					is_delete=$17, contact_id=$18`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextQuery: querySaveCD,
		helper.TextArgs:  cd,
	}
	defer tr.Finish(tags)

	if err = repository.Exec(
		mr.Repository, querySaveCD,
		cd.ID, cd.DocumentFile, cd.DocumentType, cd.DocumentTitle, cd.DocumentDescription,
		cd.NpwpNumber, cd.NpwpName, cd.NpwpAddress, cd.NpwpAddress2, cd.SiupNumber, cd.SiupCompanyName,
		cd.SiupType, cd.CreatedAt, cd.ModifiedAt, cd.CreatedBy, cd.ModifiedBy,
		cd.IsDelete, cd.ContactID,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, cd)
		return err
	}

	return nil
}

// Delete function for delete Document data
func (mr *ContactDocumentRepoPostgres) Delete(ctxReq context.Context, cd sharedModel.B2BContactDocument) (err error) {
	ctx := "ContactDocumentRepo-delete"

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextArgs: cd,
		helper.TagCtx:   ctx,
	}
	defer tr.Finish(tags)

	if err = repository.DeleteByID(mr.Repository, cd.ID, "b2b_contact_document"); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, cd)
		return err
	}

	return nil
}
