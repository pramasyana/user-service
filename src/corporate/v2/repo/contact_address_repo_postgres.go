package repo

import (
	"context"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

// ContactAddressRepoPostgres data structure
type ContactAddressRepoPostgres struct {
	*repository.Repository
}

// NewContactAddressRepoPostgres function for initializing ContactAddress repo
func NewContactAddressRepoPostgres(repo *repository.Repository) *ContactAddressRepoPostgres {
	return &ContactAddressRepoPostgres{repo}
}

// Save function for saving ContactAddress data
func (mr *ContactAddressRepoPostgres) Save(ctxReq context.Context, contactAddress sharedModel.B2BContactAddress) (err error) {
	ctx := "ContactAddressRepo-create"

	querySave := `INSERT INTO b2b_contact_address
				(
					id, name, pic_name, phone, street, province_id, province_name, 
					city_id, city_name, district_id, district_name, subdistrict_id, 
					subdistrict_name, is_billing, is_shipping, postal_code, created_at, modified_at, 
					created_by, modified_by, is_delete, contact_id
				)
			VALUES
				(
					$1, $2, $3, $4, $5, $6, $7, $8,
					$9, $10, $11,
					$12, $13, $14, $15, $16, $17,
					$18, $19, $20, $21, $22
				)
				ON CONFLICT(id)
				DO UPDATE SET
					id=$1, name=$2, pic_name=$3, phone=$4, street=$5, province_id=$6, province_name=$7, 
					city_id=$8, city_name=$9, district_id=$10, district_name=$11, subdistrict_id=$12, 
					subdistrict_name=$13, is_billing=$14, is_shipping=$15, postal_code=$16, created_at=$17, modified_at=$18, 
					created_by=$19, modified_by=$20, is_delete=$21, contact_id=$22`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextQuery: querySave,
		helper.TextArgs:  contactAddress,
	}
	defer tr.Finish(tags)

	if err = repository.Exec(
		mr.Repository, querySave,
		contactAddress.ID, contactAddress.Name, contactAddress.PicName, contactAddress.Phone,
		contactAddress.Street, contactAddress.ProvinceID, contactAddress.ProvinceName, contactAddress.CityID,
		contactAddress.CityName, contactAddress.DistrictID,
		contactAddress.DistrictName, contactAddress.SubDistrictID,
		contactAddress.SubDistrictName, contactAddress.IsBilling,
		contactAddress.IsShipping, contactAddress.PostalCode,
		contactAddress.CreatedAt, contactAddress.ModifiedAt, contactAddress.CreatedBy, contactAddress.ModifiedBy,
		contactAddress.IsDelete, contactAddress.ContactID,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, contactAddress)
		return err
	}

	return nil
}

// Update function for update Contact Address data
func (mr *ContactAddressRepoPostgres) Update(ctxReq context.Context, contactAddress sharedModel.B2BContactAddress) (err error) {
	ctx := "ContactAddressRepo-update"

	queryUpdate := `UPDATE b2b_contact_address SET name=$2, pic_name=$3, phone=$4, street=$5, province_id=$6, province_name=$7, 
	city_id=$8, city_name=$9, district_id=$10, district_name=$11, subdistrict_id=$12, 
	subdistrict_name=$13, is_billing=$14, is_shipping=$15, postal_code=$16, created_at=$17, modified_at=$18, 
	created_by=$19, modified_by=$20, is_delete=$21, contact_id=$22 WHERE id=$1;`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextQuery: queryUpdate,
		helper.TextArgs:  contactAddress,
	}
	defer tr.Finish(tags)

	if err = repository.Exec(
		mr.Repository, queryUpdate,
		contactAddress.ID, contactAddress.Name,
		contactAddress.PicName, contactAddress.Phone,
		contactAddress.Street, contactAddress.ProvinceID, contactAddress.ProvinceName, contactAddress.CityID,
		contactAddress.CityName, contactAddress.DistrictID, contactAddress.DistrictName, contactAddress.SubDistrictID,
		contactAddress.SubDistrictName, contactAddress.IsBilling, contactAddress.IsShipping, contactAddress.PostalCode,
		contactAddress.CreatedAt, contactAddress.ModifiedAt,
		contactAddress.CreatedBy, contactAddress.ModifiedBy,
		contactAddress.IsDelete, contactAddress.ContactID,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, contactAddress)
		return err
	}

	return nil
}

// Delete function for delete ContactAddress data
func (mr *ContactAddressRepoPostgres) Delete(ctxReq context.Context, contactAddress sharedModel.B2BContactAddress) (err error) {
	ctx := "ContactAddressRepo-delete"

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextArgs: contactAddress,
		helper.TagCtx:   ctx,
	}
	defer tr.Finish(tags)

	if err = repository.DeleteByID(mr.Repository, contactAddress.ID, "b2b_contact_address"); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, contactAddress)
		return err
	}

	return nil
}
