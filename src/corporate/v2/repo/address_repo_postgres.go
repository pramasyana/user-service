package repo

import (
	"context"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

// AddressRepoPostgres data structure
type AddressRepoPostgres struct {
	*repository.Repository
}

// NewAddressRepoPostgres function for initializing address repo
func NewAddressRepoPostgres(repo *repository.Repository) *AddressRepoPostgres {
	return &AddressRepoPostgres{repo}
}

// Save function for saving address data
func (mr *AddressRepoPostgres) Save(ctxReq context.Context, address sharedModel.B2BAddress) (err error) {
	ctx := "AddressRepo-create"

	querySave := `INSERT INTO b2b_address
				(
					id, type_address, name, street, province_id, province_name, 
					city_id, city_name, district_id, district_name, subdistrict_id, 
					subdistrict_name, is_primary, postal_code, created_at, modified_at, 
					created_by, modified_by, account_id, is_disabled, is_delete
				)
			VALUES
				(
					$1, $2, $3, $4, $5, $6, $7, $8,
					$9, $10, $11,
					$12, $13, $14, $15, $16, $17,
					$18, $19, $20, $21
				)
				ON CONFLICT(id)
				DO UPDATE SET
					id=$1, type_address=$2, name=$3, street=$4, province_id=$5, province_name=$6, 
					city_id=$7, city_name=$8, district_id=$9, district_name=$10, subdistrict_id=$11, 
					subdistrict_name=$12, is_primary=$13, postal_code=$14, created_at=$15, modified_at=$16, 
					created_by=$17, modified_by=$18, account_id=$19, is_disabled=$20, is_delete=$21`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextQuery: querySave,
		helper.TextArgs:  address,
	}
	defer tr.Finish(tags)

	if err = repository.Exec(
		mr.Repository, querySave,
		address.ID, address.TypeAddress, address.Name, address.Street,
		address.ProvinceID, address.ProvinceName,
		address.CityID, address.CityName,
		address.DistrictID, address.DistrictName, address.SubDistrictID, address.SubDistrictName,
		address.IsPrimary, address.PostalCode, address.CreatedAt, address.ModifiedAt,
		address.CreatedBy, address.ModifiedBy,
		address.AccountID, address.IsDisabled,
		address.IsDelete,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, address)
		return err
	}

	return nil
}

// Update function for update address data
func (mr *AddressRepoPostgres) Update(ctxReq context.Context, address sharedModel.B2BAddress) (err error) {
	ctx := "AddressRepo-update"

	queryUpdate := `UPDATE b2b_address SET type_address=$2, name=$3, street=$4, province_id=$5, province_name=$6, 
	city_id=$7, city_name=$8, district_id=$9, district_name=$10, subdistrict_id=$11, 
	subdistrict_name=$12, is_primary=$13, postal_code=$14, created_at=$15, modified_at=$16, 
	created_by=$17, modified_by=$18, account_id=$19, is_disabled=$20, is_delete=$21 WHERE id=$1;`

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextQuery: queryUpdate,
		helper.TextArgs:  address,
	}
	defer tr.Finish(tags)

	if err = repository.Exec(
		mr.Repository, queryUpdate,
		address.ID, address.TypeAddress, address.Name, address.Street,
		address.ProvinceID, address.ProvinceName, address.CityID, address.CityName,
		address.DistrictID, address.DistrictName,
		address.SubDistrictID, address.SubDistrictName,
		address.IsPrimary, address.PostalCode, address.CreatedAt, address.ModifiedAt,
		address.CreatedBy, address.ModifiedBy, address.AccountID, address.IsDisabled,
		address.IsDelete,
	); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, address)
		return err
	}

	return nil
}

// Delete function for delete address data
func (mr *AddressRepoPostgres) Delete(ctxReq context.Context, address sharedModel.B2BAddress) (err error) {
	ctx := "AddressRepo-delete"

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		helper.TextArgs: address,
		helper.TagCtx:   ctx,
	}
	defer tr.Finish(tags)

	if err = repository.DeleteByID(mr.Repository, address.ID, "b2b_address"); err != nil {
		tags[helper.TagError] = err
		helper.SendErrorLog(tr.Context(), ctx, helper.TextExecQuery, err, address)
		return err
	}

	return nil
}
