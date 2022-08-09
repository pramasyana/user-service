package repo

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/shared/repository"
	"github.com/Bhinneka/user-service/src/shipping_address/v2/model"
	"github.com/spf13/cast"
)

// ShippingAddressRepoPostgres data structure
type ShippingAddressRepoPostgres struct {
	*repository.Repository
}

// NewShippingAddressRepoPostgres function for initializing ShippingAddress repo
func NewShippingAddressRepoPostgres(repo *repository.Repository) *ShippingAddressRepoPostgres {
	return &ShippingAddressRepoPostgres{repo}
}

// Save function for saving ShippingAddress data
func (mr *ShippingAddressRepoPostgres) Save(shippingAddress model.ShippingAddress) error {
	ctx := "ShippingAddressRepo-create"

	queryInsert := `INSERT INTO b2c_shippingaddress
				(
					"id", "memberId", "name", "mobile", "phone", "provinceId", "provinceName", 
					"cityId", "cityName", "districtId", "districtName", "subdistrictId", 
					"subdistrictName", "postalCode", "street1", "street2", "version",
					"created", "lastModified", "ext", "label", "isPrimary"
				)
			VALUES
				(
					$1, $2, $3, $4, $5, $6, $7, $8,
					$9, $10, $11, $12, $13, $14, 
					$15, $16, $17, $18, $19, $20, 
					$21, $22
				)
				ON CONFLICT(id)
				DO UPDATE SET
				"id"=$1, "memberId"=$2, "name"=$3, "mobile"=$4, "phone"=$5, "provinceId"=$6, "provinceName"=$7, 
				"cityId"=$8, "cityName"=$9, "districtId"=$10, "districtName"=$11, "subdistrictId"=$12, 
				"subdistrictName"=$13, "postalCode"=$14, "street1"=$15, "street2"=$16, "version"=$17,
				"created"=$18, "lastModified"=$19, "ext"=$20, "label"=$21, "isPrimary"=$22`

	if err := repository.Exec(
		mr.Repository, queryInsert,
		shippingAddress.ID, shippingAddress.MemberID, shippingAddress.Name, shippingAddress.Mobile, shippingAddress.Phone,
		shippingAddress.ProvinceID, shippingAddress.ProvinceName, shippingAddress.CityID, shippingAddress.CityName,
		shippingAddress.DistrictID, shippingAddress.DistrictName, shippingAddress.SubDistrictID, shippingAddress.SubDistrictName,
		shippingAddress.PostalCode, shippingAddress.Street1, shippingAddress.Street2, shippingAddress.Version,
		shippingAddress.Created, shippingAddress.LastModified, shippingAddress.Ext, shippingAddress.Label, shippingAddress.IsPrimary,
	); err != nil {
		helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, shippingAddress)
		return err
	}

	return nil
}

// Update function for update Shipping Address data
func (mr *ShippingAddressRepoPostgres) Update(shippingAddress model.ShippingAddress) error {
	ctx := "ShippingAddressRepo-update"

	queryUpdate := `UPDATE b2c_shippingaddress SET "memberId"=$2, "name"=$3, "mobile"=$4, "phone"=$5, "provinceId"=$6, "provinceName"=$7, 
	"cityId"=$8, "cityName"=$9, "districtId"=$10, "districtName"=$11, "subdistrictId"=$12, 
	"subdistrictName"=$13, "postalCode"=$14, "street1"=$15, "street2"=$16, "version"=$17,
	"created"=$18, "lastModified"=$19, "ext"=$20, "label"=$21, "isPrimary"=$22 WHERE "id"=$1;`

	if err := repository.Exec(
		mr.Repository, queryUpdate,
		shippingAddress.ID, shippingAddress.MemberID, shippingAddress.Name, shippingAddress.Mobile, shippingAddress.Phone,
		shippingAddress.ProvinceID, shippingAddress.ProvinceName, shippingAddress.CityID, shippingAddress.CityName,
		shippingAddress.DistrictID, shippingAddress.DistrictName, shippingAddress.SubDistrictID, shippingAddress.SubDistrictName,
		shippingAddress.PostalCode, shippingAddress.Street1, shippingAddress.Street2, shippingAddress.Version,
		shippingAddress.Created, shippingAddress.LastModified, shippingAddress.Ext, shippingAddress.Label, shippingAddress.IsPrimary,
	); err != nil {
		helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, shippingAddress)
		return err
	}

	return nil
}

// Delete function for delete ShippingAddress data
func (mr *ShippingAddressRepoPostgres) Delete(shippingAddress model.ShippingAddress) error {
	ctx := "ShippingAddressRepo-delete"

	query := `DELETE FROM b2c_shippingaddress WHERE id=$1;`

	stmt, err := mr.WriteDB.Prepare(query)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(shippingAddress.ID)

	if err != nil {
		helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, shippingAddress)
		return err
	}

	return nil
}

// AddShippingAddress function for saving ShippingAddress data
func (mr *ShippingAddressRepoPostgres) AddShippingAddress(ctxReq context.Context, shippingAddress model.ShippingAddressData) <-chan ResultRepository {
	ctx := "ShippingAddressRepo-AddShipingAddress"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {

		var (
			shipAddrs                   model.ShippingAddressData
			isPrimary                   sql.NullBool
			extStr, labelStr, createdBy sql.NullString
		)

		query := `INSERT INTO "b2c_shippingaddress" ("id", "memberId", "name",
					"mobile", "phone", "provinceId", "provinceName", "cityId", "cityName",
					"districtId", "districtName", "subdistrictId", "subdistrictName", "postalCode", 
					"street1", "street2", "version", "created", "lastModified", "ext", "label", "isPrimary","createdBy")
				VALUES
					($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23)
				RETURNING
					"id", "memberId", "name", "mobile", "phone", "provinceId", "provinceName", "cityId", "cityName",
					"districtId", "districtName", "subdistrictId", "subdistrictName", "postalCode", "street1", "street2", "version", "created",
					"lastModified", "ext", "label","isPrimary","createdBy"`

		stmt, err := mr.WriteDB.Prepare(query)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		err = stmt.QueryRow(
			shippingAddress.ID, shippingAddress.MemberID, shippingAddress.Name, shippingAddress.Mobile, shippingAddress.Phone,
			shippingAddress.ProvinceID, shippingAddress.ProvinceName, shippingAddress.CityID, shippingAddress.CityName,
			shippingAddress.DistrictID, shippingAddress.DistrictName, shippingAddress.SubDistrictID, shippingAddress.SubDistrictName,
			shippingAddress.PostalCode, shippingAddress.Street1, shippingAddress.Street2, shippingAddress.Version,
			shippingAddress.Created, shippingAddress.LastModified, shippingAddress.Ext, shippingAddress.Label, shippingAddress.IsPrimary, shippingAddress.CreatedBy,
		).Scan(
			&shipAddrs.ID,
			&shipAddrs.MemberID,
			&shipAddrs.Name,
			&shipAddrs.Mobile,
			&shipAddrs.Phone,
			&shipAddrs.ProvinceID,
			&shipAddrs.ProvinceName,
			&shipAddrs.CityID,
			&shipAddrs.CityName,
			&shipAddrs.DistrictID,
			&shipAddrs.DistrictName,
			&shipAddrs.SubDistrictID,
			&shipAddrs.SubDistrictName,
			&shipAddrs.PostalCode,
			&shipAddrs.Street1,
			&shipAddrs.Street2,
			&shipAddrs.Version,
			&shipAddrs.Created,
			&shipAddrs.LastModified,
			&extStr,
			&labelStr,
			&isPrimary,
			&createdBy,
		)

		if isPrimary.Valid {
			shipAddrs.IsPrimary = isPrimary.Bool
		}

		if extStr.Valid && len(extStr.String) > 0 {
			shipAddrs.Ext = extStr.String
		}

		if labelStr.Valid && len(labelStr.String) > 0 {
			shipAddrs.Label = labelStr.String
		}

		if createdBy.Valid && len(createdBy.String) > 0 {
			shipAddrs.CreatedBy = createdBy.String
		}

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, shippingAddress)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: shipAddrs}
	})
	return output
}

// UpdateShippingAddress function for update ShippingAddress data
func (mr *ShippingAddressRepoPostgres) UpdateShippingAddress(ctxReq context.Context, shippingAddress model.ShippingAddressData) <-chan ResultRepository {
	ctx := "ShippingAddressRepo-UpdateShippingAddress"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {

		var (
			stmt                                    *sql.Stmt
			err                                     error
			shipAddrs                               model.ShippingAddressData
			extStr, labelStr, modifiedBy, createdBy sql.NullString
			version                                 sql.NullInt64
		)

		query := `UPDATE "b2c_shippingaddress" SET 
					"memberId" = $2, "name" = $3, "mobile" = $4, "phone" = $5, "provinceId" = $6, 
					"provinceName" = $7, "cityId" = $8, "cityName" = $9, "districtId" = $10, 
					"districtName" = $11, "subdistrictId" = $12, "subdistrictName" = $13,
					"postalCode" = $14, "street1" = $15, "street2" = $16, "version" = $17, "lastModified" = $18,
					"ext" = $19, "label" = $20, "modifiedBy" = $21, "isPrimary" = $22
				WHERE
					"id" = $1
				RETURNING
					"id", "memberId", "name", "mobile", "phone", "provinceId", "provinceName", "cityId", "cityName",
					"districtId", "districtName", "subdistrictId", "subdistrictName", "postalCode", "street1", "street2", "version", "created",
					"lastModified", "ext", "label", "isPrimary", "createdBy", "modifiedBy"`

		if mr.Tx != nil {
			stmt, err = mr.Tx.Prepare(query)
		} else {
			stmt, err = mr.WriteDB.Prepare(query)
		}

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		shippingAddress.Version++
		shippingAddress.LastModified = time.Now()

		err = stmt.QueryRow(
			shippingAddress.ID, shippingAddress.MemberID, shippingAddress.Name, shippingAddress.Mobile, shippingAddress.Phone,
			shippingAddress.ProvinceID, shippingAddress.ProvinceName, shippingAddress.CityID, shippingAddress.CityName,
			shippingAddress.DistrictID, shippingAddress.DistrictName, shippingAddress.SubDistrictID, shippingAddress.SubDistrictName,
			shippingAddress.PostalCode, shippingAddress.Street1, shippingAddress.Street2, shippingAddress.Version,
			shippingAddress.LastModified, shippingAddress.Ext, shippingAddress.Label, shippingAddress.ModifiedBy, shippingAddress.IsPrimary,
		).Scan(
			&shipAddrs.ID,
			&shipAddrs.MemberID,
			&shipAddrs.Name,
			&shipAddrs.Mobile,
			&shipAddrs.Phone,
			&shipAddrs.ProvinceID,
			&shipAddrs.ProvinceName,
			&shipAddrs.CityID,
			&shipAddrs.CityName,
			&shipAddrs.DistrictID,
			&shipAddrs.DistrictName,
			&shipAddrs.SubDistrictID,
			&shipAddrs.SubDistrictName,
			&shipAddrs.PostalCode,
			&shipAddrs.Street1,
			&shipAddrs.Street2,
			&version,
			&shipAddrs.Created,
			&shipAddrs.LastModified,
			&extStr,
			&labelStr,
			&shipAddrs.IsPrimary,
			&createdBy,
			&modifiedBy,
		)

		shipAddrs.Ext = helper.ValidateSQLNullString(extStr)
		shipAddrs.Label = helper.ValidateSQLNullString(labelStr)
		shipAddrs.ModifiedBy = helper.ValidateSQLNullString(modifiedBy)
		shipAddrs.CreatedBy = helper.ValidateSQLNullString(createdBy)
		shipAddrs.Version = cast.ToInt(helper.ValidateSQLNullInt64(version))

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, shippingAddress)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: shipAddrs}
	})
	return output
}

// CountShippingAddressByUserID function for count shipping address by user ID
func (mr *ShippingAddressRepoPostgres) CountShippingAddressByUserID(ctxReq context.Context, id string) <-chan ResultRepository {
	ctx := "ShippingAddressRepo-CountShippingAddressByUserID"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {

		defer close(output)
		var totalData int
		query := `SELECT COUNT(id) FROM "b2c_shippingaddress" WHERE "memberId" = $1`

		tags[query] = query
		stmt, err := mr.ReadDB.Prepare(query)

		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		err = stmt.QueryRow(id).Scan(&totalData)

		if err != nil {
			tags[helper.TextResponse] = err
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, id)
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: totalData}
	})
	return output
}

// FindShippingAddressPrimaryByID function for getting detail primary shipping address by id
func (mr *ShippingAddressRepoPostgres) FindShippingAddressPrimaryByID(ctxReq context.Context, memberID string) <-chan ResultRepository {
	ctx := "MerchantRepo-FindShippingAddressPrimaryByID"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {

		q := `SELECT "id", "memberId", "name", "mobile", "phone", "provinceId", "provinceName", "cityId", "cityName",
		"districtId", "districtName", "subdistrictId", "subdistrictName", "postalCode", "street1", "street2", "version", "created",
		"lastModified", "ext", "label", "isPrimary", "createdBy", "modifiedBy" FROM "b2c_shippingaddress" WHERE   "memberId" = $1 AND "isPrimary" = true`

		stmt, err := mr.ReadDB.Prepare(q)
		tags[helper.TextQuery] = q

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, q)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}
		defer stmt.Close()

		shipAddrs, err := mr.GetFieldShippingAddress(stmt, memberID)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, memberID)
			tags[helper.TextResponse] = err.Error()
			output <- ResultRepository{Error: err}
			return
		}

		tags["args"] = shipAddrs
		output <- ResultRepository{Result: shipAddrs}
	})
	return output
}

// DeleteShippingAddressByID function for delete ShippingAddress data
func (mr *ShippingAddressRepoPostgres) DeleteShippingAddressByID(ctxReq context.Context, id string) <-chan ResultRepository {
	ctx := "ShippingAddressRepo-DeleteShippingAddressByID"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		query := `DELETE FROM b2c_shippingaddress WHERE id=$1;`
		tags[helper.TextQuery] = query
		stmt, err := mr.WriteDB.Prepare(query)

		if err != nil {
			output <- ResultRepository{Error: err}
			return
		}

		_, err = stmt.Exec(id)

		if err != nil {
			tags["err"] = err
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, id)
			output <- ResultRepository{Error: err}
			return
		}
		output <- ResultRepository{Result: nil}
	})
	return output
}

// GetListShippingAddress function for loading shipping address
func (mr *ShippingAddressRepoPostgres) GetListShippingAddress(ctxReq context.Context, params *model.ParametersShippingAddress) <-chan ResultRepository {
	ctx := "ShippingAddressRepo-GetListShippingAddress"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)

		strQuery, queryValues := mr.generateQuery(params)

		sq := fmt.Sprintf(`SELECT "id", "memberId", "name", "mobile", "phone", "provinceId", "provinceName", "cityId", "cityName",
		"districtId", "districtName", "subdistrictId", "subdistrictName", "postalCode", "street1", "street2", "version", "created",
		"lastModified", "ext", "label", "isPrimary", "createdBy", "modifiedBy" FROM "b2c_shippingaddress" %s 
		ORDER BY "isPrimary" desc, "lastModified" desc
		LIMIT %d OFFSET %d`, strQuery, params.Limit, params.Offset)

		tags[helper.TextQuery] = sq
		rows, err := mr.ReadDB.Query(sq, queryValues...)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, sq)
			output <- ResultRepository{Error: nil}
			return
		}

		defer rows.Close()

		var listShippingAddress model.ListShippingAddress

		for rows.Next() {
			var (
				shipAddress                             model.ShippingAddressData
				isPrimary                               sql.NullBool
				extStr, labelStr, modifiedBy, createdBy sql.NullString
				version                                 sql.NullInt64
			)

			err = rows.Scan(
				&shipAddress.ID,
				&shipAddress.MemberID,
				&shipAddress.Name,
				&shipAddress.Mobile,
				&shipAddress.Phone,
				&shipAddress.ProvinceID,
				&shipAddress.ProvinceName,
				&shipAddress.CityID,
				&shipAddress.CityName,
				&shipAddress.DistrictID,
				&shipAddress.DistrictName,
				&shipAddress.SubDistrictID,
				&shipAddress.SubDistrictName,
				&shipAddress.PostalCode,
				&shipAddress.Street1,
				&shipAddress.Street2,
				&version,
				&shipAddress.Created,
				&shipAddress.LastModified,
				&extStr,
				&labelStr,
				&isPrimary,
				&createdBy,
				&modifiedBy,
			)

			if isPrimary.Valid {
				shipAddress.IsPrimary = isPrimary.Bool
			}

			shipAddress.Ext = helper.ValidateSQLNullString(extStr)
			shipAddress.Label = helper.ValidateSQLNullString(labelStr)
			shipAddress.ModifiedBy = helper.ValidateSQLNullString(modifiedBy)
			shipAddress.CreatedBy = helper.ValidateSQLNullString(createdBy)
			shipAddress.Version = cast.ToInt(helper.ValidateSQLNullInt64(version))

			if err != nil {
				helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, params)
				output <- ResultRepository{Error: err}
				return
			}

			listShippingAddress.ShippingAddress = append(listShippingAddress.ShippingAddress, &shipAddress)
		}

		tags[helper.TextResponse] = listShippingAddress
		output <- ResultRepository{Result: listShippingAddress}
	})

	return output

}

// GetTotalShippingAddress function for getting total of shipping address
func (mr *ShippingAddressRepoPostgres) GetTotalShippingAddress(ctxReq context.Context, params *model.ParametersShippingAddress) <-chan ResultRepository {
	ctx := "ShippingAddressRepo-GetTotalShippingAddress"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)

		var totalData int

		filter := ""
		if params.MemberID != "" {
			filter = fmt.Sprintf(`WHERE "memberId" = '%s'`, params.MemberID)
		}
		sq := fmt.Sprintf(`SELECT count(id) FROM b2c_shippingaddress %s`, filter)

		tags[helper.TextQuery] = sq
		err := mr.ReadDB.QueryRow(sq).Scan(&totalData)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, params)
			output <- ResultRepository{Error: err}
			return
		}

		tags[helper.TextResponse] = totalData
		output <- ResultRepository{Result: totalData}
	})

	return output
}

// generateQuery function for generating query
func (mr *ShippingAddressRepoPostgres) generateQuery(params *model.ParametersShippingAddress) (string, []interface{}) {
	var (
		queryString, idx string
		queryStringOR    []string
		queryStringAND   []string
		queryValues      []interface{}
		lq               int
	)

	if len(params.Query) > 0 {
		queries := strings.Split(params.Query, " ")
		lq = len(queries)
		if lq > 1 {
			queryValues = append(queryValues, "%"+params.Query+"%")
			queryStringOR = append(queryStringOR, `("name" || ' ' || "mobile"  || ' ' || "phone"   || ' ' || "id"   || ' ' || "label"   ilike $`+strconv.Itoa(len(queryStringOR)+1)+`)`)

		} else {
			queryStringOR = append(queryStringOR, `"name" ilike $`+strconv.Itoa(len(queryStringOR)+1))
			queryValues = append(queryValues, "%"+params.Query+"%")
			queryStringOR = append(queryStringOR, `"mobile" ilike $`+strconv.Itoa(len(queryStringOR)+1))
			queryValues = append(queryValues, "%"+params.Query+"%")
			queryStringOR = append(queryStringOR, `"phone" ilike $`+strconv.Itoa(len(queryStringOR)+1))
			queryValues = append(queryValues, "%"+params.Query+"%")
			queryStringOR = append(queryStringOR, `"id" ilike $`+strconv.Itoa(len(queryStringOR)+1))
			queryValues = append(queryValues, "%"+params.Query+"%")
			queryStringOR = append(queryStringOR, `"label" ilike $`+strconv.Itoa(len(queryStringOR)+1))
			queryValues = append(queryValues, "%"+params.Query+"%")
		}
	}

	if len(params.MemberID) > 0 {
		intLentOR := 0
		if len(queryStringOR) > 0 {
			intLentOR = len(queryStringOR)
		}
		idx = strconv.Itoa(intLentOR + 1)
		queryStringAND = append(queryStringAND, `"memberId" = $`+idx)
		queryValues = append(queryValues, params.MemberID)
	}

	if len(queryStringOR) > 0 || len(queryStringAND) > 0 {
		if len(queryStringOR) > 0 {
			queryString = fmt.Sprintf(`(%s)`, strings.Join(queryStringOR, " OR "))
			queryStringAND = append(queryStringAND, queryString)
		}

		if len(queryStringAND) > 0 {
			queryString = strings.Join(queryStringAND, " AND ")
		}
	}

	if len(queryString) > 0 {
		queryString = fmt.Sprintf(" WHERE %s", queryString)
	}

	return queryString, queryValues
}

// UpdatePrimaryShippingAddressByID function for update primary ShippingAddress data
func (mr *ShippingAddressRepoPostgres) UpdatePrimaryShippingAddressByID(ctxReq context.Context, memberID string) <-chan ResultRepository {
	ctx := "ShippingAddressRepo-UpdatePrimaryShippingAddressByID"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		query := `UPDATE b2c_shippingaddress SET "isPrimary"=false WHERE "memberId"=$1;`
		tags[helper.TextQuery] = query
		var (
			stmt *sql.Stmt
			err  error
		)

		if mr.Tx != nil {
			stmt, err = mr.Tx.Prepare(query)
		} else {
			stmt, err = mr.WriteDB.Prepare(query)
		}

		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		_, err = stmt.Exec(memberID)

		if err != nil {
			tags[helper.TextResponse] = err
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, memberID)
			output <- ResultRepository{Error: err}
			return
		}
		output <- ResultRepository{Result: nil}
	})
	return output
}

// FindShippingAddressByID function for getting detail shipping address by id
func (mr *ShippingAddressRepoPostgres) FindShippingAddressByID(ctxReq context.Context, id string, memberID string) <-chan ResultRepository {
	ctx := "MerchantRepo-FindShippingAddressById"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {

		filter := ""
		if memberID != "" {
			filter = fmt.Sprintf(`AND "memberId" = '%s'`, memberID)
		}
		q := fmt.Sprintf(`SELECT "id", "memberId", "name", "mobile", "phone", "provinceId", "provinceName", "cityId", "cityName",
		"districtId", "districtName", "subdistrictId", "subdistrictName", "postalCode", "street1", "street2", "version", "created",
		"lastModified", "ext", "label", "isPrimary", "createdBy", "modifiedBy" FROM "b2c_shippingaddress" WHERE "id" = $1 %s`, filter)

		stmt, err := mr.ReadDB.Prepare(q)
		tags[helper.TextQuery] = q

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, q)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}
		defer stmt.Close()

		shipAddrs, err := mr.GetFieldShippingAddress(stmt, id)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, memberID)
			tags[helper.TextResponse] = err.Error()
			output <- ResultRepository{Error: err}
			return
		}

		tags["args"] = shipAddrs
		output <- ResultRepository{Result: shipAddrs}
	})
	return output
}

// GetFieldShippingAddress function for getting field struct shipping address
func (mr *ShippingAddressRepoPostgres) GetFieldShippingAddress(stmt *sql.Stmt, id string) (model.ShippingAddressData, error) {
	var (
		shipAddrs                               model.ShippingAddressData
		extStr, labelStr, modifiedBy, createdBy sql.NullString
		isPrimary                               sql.NullBool
		version                                 sql.NullInt64
	)

	err := stmt.QueryRow(id).Scan(
		&shipAddrs.ID,
		&shipAddrs.MemberID,
		&shipAddrs.Name,
		&shipAddrs.Mobile,
		&shipAddrs.Phone,
		&shipAddrs.ProvinceID,
		&shipAddrs.ProvinceName,
		&shipAddrs.CityID,
		&shipAddrs.CityName,
		&shipAddrs.DistrictID,
		&shipAddrs.DistrictName,
		&shipAddrs.SubDistrictID,
		&shipAddrs.SubDistrictName,
		&shipAddrs.PostalCode,
		&shipAddrs.Street1,
		&shipAddrs.Street2,
		&version,
		&shipAddrs.Created,
		&shipAddrs.LastModified,
		&extStr,
		&labelStr,
		&isPrimary,
		&createdBy,
		&modifiedBy,
	)

	if isPrimary.Valid {
		shipAddrs.IsPrimary = isPrimary.Bool
	}

	shipAddrs.Ext = helper.ValidateSQLNullString(extStr)
	shipAddrs.Label = helper.ValidateSQLNullString(labelStr)
	shipAddrs.ModifiedBy = helper.ValidateSQLNullString(modifiedBy)
	shipAddrs.CreatedBy = helper.ValidateSQLNullString(createdBy)
	shipAddrs.Version = cast.ToInt(helper.ValidateSQLNullInt64(version))

	if err != nil {
		return shipAddrs, err
	}
	return shipAddrs, nil
}
