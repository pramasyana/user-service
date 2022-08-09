package repo

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
	"github.com/spf13/cast"
)

// MerchantAddressRepoPostgres data structure
type MerchantAddressRepoPostgres struct {
	*repository.Repository
}

// NewMerchantAddressRepoPostgres function for initializing  repo
func NewMerchantAddressRepoPostgres(repo *repository.Repository) MerchantAddressRepository {
	return &MerchantAddressRepoPostgres{repo}
}

// CountAddress function for count merchant address by merchant ID
func (mr *MerchantAddressRepoPostgres) CountAddress(ctxReq context.Context, relationID, relationName string, params *model.ParameterWarehouse) <-chan ResultRepository {
	ctx := "MerchantAddressRepo-CountAddress"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {

		defer close(output)
		var totalData int
		bindVars := []interface{}{relationID, relationName}

		query := `SELECT COUNT(id) FROM "address" WHERE "relationId" = $1 AND "relationName" = $2`
		if params.Query != "" {
			query = query + ` AND lower("label") like $3`
			bindVars = append(bindVars, "%"+strings.ToLower(params.Query)+"%")
		}

		tags[query] = query
		stmt, err := mr.ReadDB.Prepare(query)

		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}
		defer stmt.Close()

		if err := stmt.QueryRow(bindVars...).Scan(&totalData); err != nil {
			tags[helper.TextResponse] = err
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, []string{relationID, relationName})
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: totalData}
	})
	return output
}

// AddUpdateAddress function for saving address data
func (mr *MerchantAddressRepoPostgres) AddUpdateAddress(ctxReq context.Context, address model.AddressData) <-chan ResultRepository {
	ctx := "MerchantAddressRepo-AddAddress"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		query := `INSERT INTO "address" 
					("id", "relationId", "relationName", "type", "label",
					"address", "provinceId", "provinceName",  "cityId", "cityName", 
					"districtId", "districtName", "subdistrictId", "subdistrictName", "postalCode",  
					"created", "lastModified", "createdBy", "modifiedBy", "version", "isPrimary", "status")
				VALUES
					($1, $2, $3, $4, $5, 
					$6, $7, $8, $9, $10, 
					$11, $12, $13, $14, $15, 
					$16, $17, $18, $19, $20, $21, $22)
				ON CONFLICT("id")
				DO UPDATE SET
					"relationId" = $2, "relationName" = $3, "type" = $4, "label" = $5,
					"address" = $6, "provinceId" = $7, "provinceName" = $8,  "cityId" = $9, "cityName" = $10, 
					"districtId" = $11, "districtName" = $12, "subdistrictId" = $13, "subdistrictName" = $14, "postalCode" = $15,  
					"created" = $16, "lastModified" = $17, "createdBy" = $18, "modifiedBy" = $19, "version" = $20, "isPrimary" = $21,
					"status" = $22
				RETURNING
					"id", "relationId", "relationName", "type", "label",
					"address", "provinceId", "provinceName",  "cityId", "cityName", 
					"districtId", "districtName", "subdistrictId", "subdistrictName", "postalCode",  
					"created", "lastModified", "createdBy", "modifiedBy", "version", "isPrimary", "status"`

		var (
			addressData           model.AddressData
			createdBy, modifiedBy sql.NullString
			version               sql.NullInt64
			isPrimary             sql.NullBool
			stmt                  *sql.Stmt
			err                   error
		)

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

		err = stmt.QueryRow(
			address.ID, address.RelationID, address.RelationName, address.Type, address.Label,
			address.Address, address.ProvinceID, address.ProvinceName, address.CityID, address.CityName,
			address.DistrictID, address.DistrictName, address.SubDistrictID, address.SubDistrictName, address.PostalCode,
			address.Created, address.LastModified, address.CreatedBy, address.ModifiedBy, address.Version, address.IsPrimary,
			address.Status,
		).Scan(
			&addressData.ID, &addressData.RelationID, &addressData.RelationName, &addressData.Type, &addressData.Label,
			&addressData.Address, &addressData.ProvinceID, &addressData.ProvinceName, &addressData.CityID, &addressData.CityName,
			&addressData.DistrictID, &addressData.DistrictName, &addressData.SubDistrictID, &addressData.SubDistrictName, &addressData.PostalCode,
			&addressData.Created, &addressData.LastModified, &createdBy, &modifiedBy, &version, &isPrimary, &addressData.Status,
		)
		if isPrimary.Valid {
			addressData.IsPrimary = isPrimary.Bool
		}

		addressData.CreatedString = addressData.Created.Format(time.RFC3339)
		addressData.LastModifiedString = addressData.LastModified.Format(time.RFC3339)
		addressData.CreatedBy = helper.ValidateSQLNullString(createdBy)
		addressData.ModifiedBy = helper.ValidateSQLNullString(modifiedBy)
		addressData.Version = cast.ToInt(helper.ValidateSQLNullInt64(version))

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: addressData}
	})
	return output
}

// UpdatePhoneAddress function for saving address data
func (mr *MerchantAddressRepoPostgres) UpdatePhoneAddress(ctxReq context.Context, address model.PhoneData) <-chan ResultRepository {
	ctx := "MerchantAddressRepo-UpdatePhoneAddress"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		updateQuery := `UPDATE "phone" SET 
					"label" = $3, "number" = $4, "isPrimary" = $6, "created" = $7, 
					"lastModified" = $8, "createdBy" = $9, "modifiedBy" = $10, "version" = $11
				WHERE  
					"relationId" = $1 AND "relationName" = $2 AND "type" = $5
				RETURNING
					"relationId", "relationName", "label", "number", "type", 
					"isPrimary", "created",  "lastModified", "createdBy", "modifiedBy", 
					"version"`

		updatePhoneResult, err := mr.AddUpdatePhoneAddress(ctxReq, updateQuery, address)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, address)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: updatePhoneResult}
	})
	return output
}

// AddPhoneAddress function for saving address data
func (mr *MerchantAddressRepoPostgres) AddPhoneAddress(ctxReq context.Context, address model.PhoneData) <-chan ResultRepository {
	ctx := "MerchantAddressRepo-AddPhoneAddress"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		addQuery := `INSERT INTO "phone" 
					("relationId", "relationName", "label", "number", "type", 
					"isPrimary", "created", "lastModified", "createdBy", "modifiedBy", 
					"version")
				VALUES
					($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
				RETURNING
					"relationId", "relationName", "label", "number", "type", 
					"isPrimary", "created", "lastModified", "createdBy", "modifiedBy", 
					"version"`

		addPhoneResult, err := mr.AddUpdatePhoneAddress(ctxReq, addQuery, address)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, address)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: addPhoneResult}
	})
	return output
}

// AddUpdatePhoneAddress function for saving merchant address data
func (mr *MerchantAddressRepoPostgres) AddUpdatePhoneAddress(ctxReq context.Context, query string, phone model.PhoneData) (model.PhoneData, error) {

	var (
		phoneData             model.PhoneData
		isPrimary             sql.NullBool
		createdBy, modifiedBy sql.NullString
		version               sql.NullInt64
		stmt                  *sql.Stmt
		err                   error
	)

	if mr.Tx != nil {
		stmt, err = mr.Tx.Prepare(query)
	} else {
		stmt, err = mr.WriteDB.Prepare(query)
	}
	if err != nil {
		return phoneData, err
	}

	err = stmt.QueryRow(
		phone.RelationID, phone.RelationName, phone.Label, phone.Number, phone.Type,
		phone.IsPrimary, phone.Created, phone.LastModified, phone.CreatedBy, phone.ModifiedBy, phone.Version,
	).Scan(
		&phoneData.RelationID, &phoneData.RelationName, &phoneData.Label, &phoneData.Number, &phoneData.Type,
		&isPrimary, &phoneData.Created, &phoneData.LastModified, &createdBy, &modifiedBy, &version,
	)

	phoneData.CreatedString = phoneData.Created.Format(time.RFC3339)
	phoneData.LastModifiedString = phoneData.LastModified.Format(time.RFC3339)
	phoneData.CreatedBy = helper.ValidateSQLNullString(createdBy)
	phoneData.ModifiedBy = helper.ValidateSQLNullString(modifiedBy)
	phoneData.Version = cast.ToInt(helper.ValidateSQLNullInt64(version))

	if isPrimary.Valid {
		phoneData.IsPrimary = isPrimary.Bool
	}

	if err != nil {
		return phoneData, err
	}
	return phoneData, nil
}

// FindMerchantAddress function for get address by ID
func (mr *MerchantAddressRepoPostgres) FindMerchantAddress(ctxReq context.Context, id string) <-chan ResultRepository {
	ctx := "MerchantAddressRepo-FindMerchantAddress"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {

		defer close(output)
		var (
			addressData           model.WarehouseData
			createdBy, modifiedBy sql.NullString
			version               sql.NullInt64
			isPrimary             sql.NullBool
		)

		query := `SELECT 
		"id", "relationId", "label", "provinceId", "provinceName",  
		"cityId", "cityName", "districtId", "districtName", "subdistrictId", 
		"subdistrictName", "postalCode", "address", "created", "lastModified", 
		"createdBy", "modifiedBy", "version", "type", "isPrimary", "status"
		FROM "address" WHERE "id" = $1`

		tags[query] = query
		stmt, err := mr.ReadDB.Prepare(query)

		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		err = stmt.QueryRow(id).Scan(
			&addressData.ID, &addressData.MerchantID, &addressData.Label, &addressData.ProvinceID, &addressData.ProvinceName,
			&addressData.CityID, &addressData.CityName, &addressData.DistrictID, &addressData.DistrictName, &addressData.SubDistrictID,
			&addressData.SubDistrictName, &addressData.PostalCode, &addressData.Address, &addressData.Created, &addressData.LastModified,
			&createdBy, &modifiedBy, &version, &addressData.Type, &isPrimary, &addressData.Status,
		)

		if isPrimary.Valid {
			addressData.IsPrimary = isPrimary.Bool
		}

		addressData.CreatedBy = helper.ValidateSQLNullString(createdBy)
		addressData.ModifiedBy = helper.ValidateSQLNullString(modifiedBy)
		addressData.CreatedString = addressData.Created.Format(time.RFC3339)
		addressData.LastModifiedString = addressData.LastModified.Format(time.RFC3339)
		addressData.Version = cast.ToInt(helper.ValidateSQLNullInt64(version))

		if err != nil {
			tags[helper.TextResponse] = err
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, id)
			output <- ResultRepository{Error: err}
			return
		}

		resultPhone := <-mr.GetPhoneWarehouse(ctxReq, addressData)
		addressData, _ = resultPhone.Result.(model.WarehouseData)

		output <- ResultRepository{Result: addressData}
	})
	return output
}

// GetPhoneWarehouse function for loading phone data
func (mr *MerchantAddressRepoPostgres) GetPhoneWarehouse(ctxReq context.Context, addressData model.WarehouseData) <-chan ResultRepository {
	ctx := "MerchantAddressRepo-GetPhone"
	output := make(chan ResultRepository)
	go func() {
		defer close(output)

		query := `SELECT 
		"relationId", "relationName", "label", "number", "type", 
		"isPrimary", "created", "lastModified", "createdBy", "modifiedBy", "version"
		FROM "phone" 
		WHERE "relationId" = $1 AND "relationName" = $2`

		rows, err := mr.ReadDB.Query(query, addressData.ID, model.AddressString)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, addressData)
			output <- ResultRepository{Error: nil}
			return
		}

		defer rows.Close()

		for rows.Next() {
			var (
				phone                 model.PhoneData
				isPrimary             sql.NullBool
				createdBy, modifiedBy sql.NullString
				version               sql.NullInt64
			)
			if err := rows.Scan(
				&phone.RelationID, &phone.RelationName, &phone.Label, &phone.Number, &phone.Type,
				&isPrimary, &phone.Created, &phone.LastModified, &createdBy, &modifiedBy, &version,
			); err != nil {
				helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, addressData)
				output <- ResultRepository{Error: nil}
				return
			}

			addressData.Name = phone.Label
			switch phone.Type {
			case model.MobileString:
				addressData.Mobile = phone.Number
			case model.PhoneString:
				addressData.Phone = phone.Number
			}
		}

		output <- ResultRepository{Result: addressData}
	}()

	return output

}

// GetListAddress function for loading merchant address
func (mr *MerchantAddressRepoPostgres) GetListAddress(ctxReq context.Context, params *model.ParameterWarehouse) <-chan ResultRepository {
	ctx := "MerchantAddressRepo-GetListAddress"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)

		bindVars := []interface{}{params.MerchantID, model.MerchantString, params.Limit, params.Offset}
		orderBy := "id"
		sort := "DESC"

		sq := `SELECT
		"id", "relationId", "label", "provinceId", "provinceName",
		"cityId", "cityName", "districtId", "districtName", "subdistrictId",
		"subdistrictName", "postalCode", "address", "created", "lastModified",
		"createdBy", "modifiedBy", "version", "type", "isPrimary", "status"
		FROM "address" WHERE "relationId"=$1 AND "relationName"=$2
		`
		if params.Query != "" {
			sq = sq + `AND LOWER("label") like $5`
			bindVars = append(bindVars, "%"+strings.ToLower(params.Query)+"%")
		}
		if params.OrderBy != "" {
			orderBy = params.OrderBy
		}
		if params.Sort != "" {
			sort = params.Sort
		}

		sq = sq + ` ORDER BY "` + orderBy + `" ` + sort
		sq = sq + ` LIMIT $3 OFFSET $4`

		tags[helper.TextQuery] = sq
		rows, err := mr.ReadDB.Query(sq, bindVars...)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, params)
			output <- ResultRepository{Error: nil}
			return
		}

		defer rows.Close()
		var listWarehouse model.ListWarehouse

		for rows.Next() {
			var (
				address                         model.WarehouseData
				isPrimary                       sql.NullBool
				labelStr, modifiedBy, createdBy sql.NullString
				version                         sql.NullInt64
			)

			err = rows.Scan(
				&address.ID, &address.MerchantID, &labelStr, &address.ProvinceID, &address.ProvinceName,
				&address.CityID, &address.CityName, &address.DistrictID, &address.DistrictName, &address.SubDistrictID,
				&address.SubDistrictName, &address.PostalCode, &address.Address, &address.Created, &address.LastModified,
				&createdBy, &modifiedBy, &version, &address.Type, &isPrimary, &address.Status,
			)
			if err != nil {
				helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, params)
				output <- ResultRepository{Error: err}
				return
			}

			mr.adjustData(&address, labelStr, modifiedBy, createdBy, version, isPrimary)

			resultPhone := <-mr.GetPhoneWarehouse(ctxReq, address)
			address, _ = resultPhone.Result.(model.WarehouseData)

			listWarehouse.WarehouseData = append(listWarehouse.WarehouseData, &address)
		}

		tags[helper.TextResponse] = listWarehouse
		output <- ResultRepository{Result: listWarehouse}
	})

	return output
}

func (mr *MerchantAddressRepoPostgres) adjustData(address *model.WarehouseData, labelStr, modifiedBy, createdBy sql.NullString, version sql.NullInt64, isPrimary sql.NullBool) {
	address.Label = helper.ValidateSQLNullString(labelStr)
	address.ModifiedBy = helper.ValidateSQLNullString(modifiedBy)
	address.CreatedBy = helper.ValidateSQLNullString(createdBy)
	address.Version = cast.ToInt(helper.ValidateSQLNullInt64(version))
	address.CreatedString = address.Created.Format(time.RFC3339)
	address.LastModifiedString = address.LastModified.Format(time.RFC3339)

	if isPrimary.Valid {
		address.IsPrimary = isPrimary.Bool
	}
}

// UpdatePrimaryAddressByRelationID function for update primary as false all address data
func (mr *MerchantAddressRepoPostgres) UpdatePrimaryAddressByRelationID(ctxReq context.Context, relationID, relationName string) <-chan ResultRepository {
	ctx := "MerchantAddressRepo-UpdatePrimaryAddressByRelationID"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		queryUpdatePrimary := `UPDATE address SET "isPrimary"=false WHERE "relationId"=$1 AND "relationName"=$2`
		tags[helper.TextQuery] = queryUpdatePrimary

		var queryValuesUpdatePrimary []interface{}
		queryValuesUpdatePrimary = append(queryValuesUpdatePrimary, relationID)
		queryValuesUpdatePrimary = append(queryValuesUpdatePrimary, relationName)
		err := mr.ExecQuery(queryUpdatePrimary, queryValuesUpdatePrimary)
		if err != nil {
			tags[helper.TextResponse] = err
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, []string{relationID, relationName})
			output <- ResultRepository{Error: err}
			return
		}
		output <- ResultRepository{Result: nil}
	})
	return output
}

// DeleteWarehouseAddress function for delete address data
func (mr *MerchantAddressRepoPostgres) DeleteWarehouseAddress(ctxReq context.Context, id string) <-chan ResultRepository {
	ctx := "MerchantAddressRepo-DeleteWarehouseAddress"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		queryDeleteWarehouse := `DELETE FROM address WHERE id=$1;`
		tags[helper.TextQuery] = queryDeleteWarehouse

		var queryValuesDeleteWarehouse []interface{}
		queryValuesDeleteWarehouse = append(queryValuesDeleteWarehouse, id)
		err := mr.ExecQuery(queryDeleteWarehouse, queryValuesDeleteWarehouse)
		if err != nil {
			tags[helper.TextResponse] = err
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, id)
			output <- ResultRepository{Error: err}
			return
		}
		output <- ResultRepository{Result: nil}
	})
	return output
}

// DeletePhoneAddress function for delete address data
func (mr *MerchantAddressRepoPostgres) DeletePhoneAddress(ctxReq context.Context, relationID, relationName string) <-chan ResultRepository {
	ctx := "MerchantAddressRepo-DeletePhoneAddress"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		queryDeletePhone := `DELETE FROM phone WHERE "relationId"=$1 AND "relationName"=$2`
		tags[helper.TextQuery] = queryDeletePhone

		var queryValuesDeletePhone []interface{}
		queryValuesDeletePhone = append(queryValuesDeletePhone, relationID)
		queryValuesDeletePhone = append(queryValuesDeletePhone, relationName)
		err := mr.ExecQuery(queryDeletePhone, queryValuesDeletePhone)
		if err != nil {
			tags[helper.TextResponse] = err
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, []string{relationID, relationName})
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: nil}
	})
	return output
}

// ExecQuery function for execution query
func (mr *MerchantAddressRepoPostgres) ExecQuery(query string, queryValues []interface{}) error {
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
		return err
	}

	_, err = stmt.Exec(queryValues...)
	if err != nil {
		return err
	}

	return nil
}

// AddUpdateAddressMaps function for saving  data
func (mr *MerchantAddressRepoPostgres) AddUpdateAddressMaps(ctxReq context.Context, maps model.Maps) <-chan ResultRepository {
	ctx := "MerchantAddressRepo-AddUpdateAddressMaps"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		var (
			stmt *sql.Stmt
			err  error
		)

		query := `INSERT INTO maps
					(
						"id", "relationId", "relationName", "label", "latitude", "longitude"
					)
				VALUES
					(
						$1, $2, $3, $4, $5, $6
					)
				ON CONFLICT("relationId", "relationName")
				DO UPDATE SET
				"relationId"=$2, "relationName"=$3, "label"=$4, "latitude"=$5, "longitude"=$6;`

		tags[helper.TextQuery] = query
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
		defer stmt.Close()

		_, err = stmt.Exec(
			maps.ID, maps.RelationID, maps.RelationName, maps.Label, maps.Latitude, maps.Longitude,
		)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		tags["args"] = maps

		output <- ResultRepository{Error: nil}
	})
	return output
}

// FindAddressMaps function for get address by ID
func (mr *MerchantAddressRepoPostgres) FindAddressMaps(ctxReq context.Context, relationID, relationName string) <-chan ResultRepository {
	ctx := "MerchantAddressRepo-FindAddressMaps"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {

		defer close(output)
		var (
			mapsData model.Maps
		)

		query := `SELECT 
		"id", "relationId", "relationName", "label", "latitude", "longitude"
		FROM "maps" WHERE "relationId" = $1 AND "relationName" = $2`

		tags[query] = query
		stmt, err := mr.ReadDB.Prepare(query)

		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		err = stmt.QueryRow(relationID, relationName).Scan(
			&mapsData.ID, &mapsData.RelationID, &mapsData.RelationName, &mapsData.Label, &mapsData.Latitude, &mapsData.Longitude,
		)

		output <- ResultRepository{Result: mapsData}
	})
	return output
}
