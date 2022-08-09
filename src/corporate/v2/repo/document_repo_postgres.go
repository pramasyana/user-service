package repo

import (
	"context"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

// DocumentRepoPostgres data structure
type DocumentRepoPostgres struct {
	*repository.Repository
}

// NewDocumentRepoPostgres function for initializing Document repo
func NewDocumentRepoPostgres(repo *repository.Repository) *DocumentRepoPostgres {
	return &DocumentRepoPostgres{repo}
}

// Save function for saving Document data
func (mr *DocumentRepoPostgres) Save(ctxReq context.Context, document sharedModel.B2BDocument) (err error) {
	ctx := "DocumentRepo-create"

	querySave := `INSERT INTO b2b_document
				(
					id, document_type, document_file, document_title, document_description, 
					npwp_number, npwp_name, npwp_address, siup_number, siup_company_name,
					siup_type, is_delete, created_at, modified_at, created_by, modified_by, 
					account_id, is_disabled, province_id, province_name, city_id, city_name,
					district_id, district_name, subdistrict_id, subdistrict_name, 
					postal_code, street

				)
			VALUES
				(
					$1, $2, $3, $4, $5, $6, $7, $8,
					$9, $10, $11, $12, $13, $14, 
					$15, $16, $17, $18, $19, $20, 
					$21, $22, $23, $24, $25, $26,
					$27, $28
				)
				ON CONFLICT(id)
				DO UPDATE SET
					id=$1, document_type=$2, document_file=$3, document_title=$4, document_description=$5, 
					npwp_number=$6, npwp_name=$7, npwp_address=$8, siup_number=$9, siup_company_name=$10,
					siup_type=$11, is_delete=$12, created_at=$13, modified_at=$14, created_by=$15, modified_by=$16, 
					account_id=$17, is_disabled=$18, province_id=$19, province_name=$20, city_id=$21, city_name=$22,
					district_id=$23, district_name=$24, subdistrict_id=$25, subdistrict_name=$26, 
					postal_code=$27, street=$28`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextQuery: querySave,
		helper.TextArgs:  document,
	}
	defer tr.Finish(tags)

	if err = repository.Exec(
		mr.Repository, querySave,
		document.ID, document.DocumentType, document.DocumentFile,
		document.DocumentTitle,
		document.DocumentDescription, document.NpwpNumber,
		document.NpwpName,
		document.NpwpAddress, document.SiupNumber, document.SiupCompanyName,
		document.SiupType, document.IsDelete, document.CreatedAt, document.ModifiedAt,
		document.CreatedBy, document.ModifiedBy,
		document.AccountID, document.IsDisabled, document.ProvinceID, document.ProvinceName,
		document.CityID, document.CityName,
		document.DistrictID,
		document.DistrictName, document.SubdistrictID, document.SubdistrictName,
		document.PostalCode, document.Street,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, document)
		return err
	}

	return nil
}

// Update function for update Document data
func (mr *DocumentRepoPostgres) Update(ctxReq context.Context, document sharedModel.B2BDocument) (err error) {
	ctx := "DocumentRepo-update"

	queryUpdate := `UPDATE b2b_document SET document_type=$2, document_file=$3, document_title=$4, document_description=$5, 
	npwp_number=$6, npwp_name=$7, npwp_address=$8, siup_number=$9, siup_company_name=$10,
	siup_type=$11, is_delete=$12, created_at=$13, modified_at=$14, created_by=$15, modified_by=$16, 
	account_id=$17, is_disabled=$18, province_id=$19, province_name=$20, city_id=$21, city_name=$22,
	district_id=$23, district_name=$24, subdistrict_id=$25, subdistrict_name=$26, 
	postal_code=$27, street=$28 WHERE id=$1;`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextQuery: queryUpdate,
		helper.TextArgs:  document,
	}
	defer tr.Finish(tags)

	if err = repository.Exec(
		mr.Repository, queryUpdate,
		document.ID, document.DocumentType, document.DocumentFile, document.DocumentTitle,
		document.DocumentDescription, document.NpwpNumber, document.NpwpName,
		document.NpwpAddress, document.SiupNumber, document.SiupCompanyName,
		document.SiupType, document.IsDelete, document.CreatedAt, document.ModifiedAt,
		document.CreatedBy, document.ModifiedBy,
		document.AccountID, document.IsDisabled, document.ProvinceID, document.ProvinceName,
		document.CityID, document.CityName, document.DistrictID,
		document.DistrictName, document.SubdistrictID, document.SubdistrictName,
		document.PostalCode, document.Street,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, document)
		return err
	}

	return nil
}

// Delete function for delete Document data
func (mr *DocumentRepoPostgres) Delete(ctxReq context.Context, document sharedModel.B2BDocument) (err error) {
	ctx := "DocumentRepo-delete"

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextArgs: document,
		helper.TagCtx:   ctx,
	}
	defer tr.Finish(tags)

	if err = repository.DeleteByID(mr.Repository, document.ID, "b2b_document"); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, document)
		return err
	}

	return nil
}
