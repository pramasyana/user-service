package repo

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

const StringWhere = " WHERE "

// MerchantEmployeeRepoPostgres data structure
type MerchantEmployeeRepoPostgres struct {
	*repository.Repository
}

// NewMerchantEmployeeRepoPostgres function for initializing  repo
func NewMerchantEmployeeRepoPostgres(repo *repository.Repository) *MerchantEmployeeRepoPostgres {
	return &MerchantEmployeeRepoPostgres{repo}
}

// Save function for saving  data
func (mr *MerchantEmployeeRepoPostgres) Save(merchantBank model.B2CMerchantEmployee) error {
	ctx := "MerchantEmployeeRepo-create"

	query := `INSERT INTO b2c_merchant_employees
				(
					"id", "merchantId", "memberId", "createdAt", "createdBy"
				)
			VALUES
				(
					$1, $2, $3, $4, $5

				)
			ON CONFLICT("merchantId", "memberId") 
			DO UPDATE SET 
			"merchantId"=$2, "memberId"=$3, "modifiedAt"=$4, "modifiedBy"=$5`

	stmt, err := mr.WriteDB.Prepare(query)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		merchantBank.ID, merchantBank.MerchantID, merchantBank.MemberID, merchantBank.CreatedAt, merchantBank.CreatedBy,
	)

	if err != nil {
		helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, merchantBank)
		return err
	}

	return nil
}

// ChangeStatus function for flagging delete
func (mr *MerchantEmployeeRepoPostgres) ChangeStatus(ctxReq context.Context, params model.B2CMerchantEmployee) error {
	ctx := "MerchantEmployeeRepo-ChangeStatus"

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := make(map[string]interface{})
	defer func() {
		tr.Finish(tags)
	}()

	query := `UPDATE "b2c_merchant_employees" SET "status"=$1, "modifiedAt"=$2, "modifiedBy"=$3 WHERE "merchantId"=$4 AND "memberId"=$5;`

	tags[helper.TextQuery] = query
	tags[helper.TextMerchantIDCamel] = params.MerchantID
	tags[helper.TextMemberIDCamel] = params.MemberID
	tags[helper.TextStatus] = params.Status

	stmt, err := mr.WriteDB.Prepare(query)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, map[string]string{
			helper.TextMerchantIDCamel: params.MerchantID,
			params.MemberID:            params.MemberID,
			helper.TextStatus:          params.Status,
		})
		tags[helper.TextResponse] = err
		return err
	}

	if _, err = stmt.Exec(params.Status, params.ModifiedAt, params.ModifiedBy, params.MerchantID, params.MemberID); err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, map[string]string{
			helper.TextMerchantIDCamel: params.MerchantID,
			params.MemberID:            params.MemberID,
			helper.TextStatus:          params.Status,
		})
		tags[helper.TextResponse] = err
		return err
	}

	return nil
}

// GetAllMerchantEmployees retrieve merchant data based on given parameters
func (mr *MerchantEmployeeRepoPostgres) GetAllMerchantEmployees(ctxReq context.Context, params *model.QueryMerchantEmployeeParameters) <-chan ResultRepository {
	ctx := "MerchantRepoPostgres-GetAllMerchantEmployees"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		query := `SELECT 
					me."id", me."merchantId", me."memberId", 
					m."firstName", m."lastName", m."email", m."gender", m."mobile", m."phone", m."birthDate", 
					me."createdAt", me."modifiedAt", me."status", m."profilePicture", 
					bcm."merchantLogo", bcm."merchantName", bcm."merchantType", 
					bcm."vanityURL", bcm."isActive", bcm."isPKP"
				FROM b2c_merchant_employees me
				JOIN member m on me."memberId" = m."id"
				JOIN b2c_merchant bcm on me."merchantId" = bcm."id"`

		queries, queryValues := params.Build()
		if len(queries) > 0 {
			query += StringWhere + fmt.Sprintf("%s", strings.Join(queries, ` AND `))
		}

		if params.OrderBy != "" {
			params.OrderBy = `me."` + params.OrderBy + `"`
		}
		query = helper.RawQuery(query, params.Offset, params.Limit, params.OrderBy, params.SortBy)
		tags[helper.TextQuery] = query

		rows, err := mr.ReadDB.Query(query, queryValues...)
		if err != nil {
			output <- ResultRepository{Error: err}
			return
		}
		results := []model.B2CMerchantEmployeeData{}
		for rows.Next() {
			var (
				row model.B2CMerchantEmployeeData
			)
			err := rows.Scan(
				&row.ID, &row.MerchantID, &row.MemberID,
				&row.FirstName, &row.LastName, &row.Email, &row.Gender, &row.Mobile, &row.Phone, &row.BirthDate,
				&row.CreatedAt, &row.ModifiedAt, &row.Status, &row.ProfilePicture,
				&row.MerchantLogo, &row.MerchantName, &row.MerchantType,
				&row.VanityURL, &row.IsActive, &row.IsPKP)

			row.BirthDateString.String = row.BirthDate.Format("02/01/2006")

			if err != nil {
				helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, params)
				output <- ResultRepository{Error: err}
				return
			}
			if _, err := json.Marshal(&row); err != nil {
				helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, params)
				output <- ResultRepository{Error: err}
				return
			}

			results = append(results, row)
		}

		output <- ResultRepository{Result: results}
	})

	return output
}

// GetTotalMerchantEmployees return total merchants
func (mr *MerchantEmployeeRepoPostgres) GetTotalMerchantEmployees(ctxReq context.Context, params *model.QueryMerchantEmployeeParameters) <-chan ResultRepository {
	ctx := "MerchantRepoPostgres-GetTotalMerchantEmployees"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)
		var totalData int

		query := `SELECT 
					COUNT(me.id) 
				FROM b2c_merchant_employees me
				JOIN member m on me."memberId" = m."id"
				JOIN b2c_merchant bcm on me."merchantId" = bcm."id"`
		queries, queryValues := params.Build()
		if len(queries) > 0 {
			query += StringWhere + fmt.Sprintf("%s", strings.Join(queries, ` AND `))
		}

		tags[helper.TextQuery] = query

		if err := mr.ReadDB.QueryRow(query, queryValues...).Scan(&totalData); err != nil {
			output <- ResultRepository{Error: err}
			return
		}
		output <- ResultRepository{Result: totalData}
	})
	return output
}

// GetMerchantEmployees retrieve merchant data based on given parameters
func (mr *MerchantEmployeeRepoPostgres) GetMerchantEmployees(ctxReq context.Context, params *model.QueryMerchantEmployeeParameters) <-chan ResultRepository {
	ctx := "MerchantRepoPostgres-GetMerchantEmployees"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		query := `SELECT 
					me."id", me."merchantId", me."memberId", 
					m."firstName", m."lastName", m."email", m."gender", m."mobile", m."phone", m."birthDate", 
					me."createdAt", me."modifiedAt", me."status", m."profilePicture",
					bcm."merchantLogo", bcm."merchantName", bcm."merchantType", 
					bcm."vanityURL", bcm."isActive", bcm."isPKP" 
				FROM b2c_merchant_employees me
				JOIN member m on me."memberId" = m."id"
				JOIN b2c_merchant bcm on me."merchantId" = bcm."id"`

		queries, queryValues := params.Build()
		if len(queries) > 0 {
			query += StringWhere + fmt.Sprintf("%s", strings.Join(queries, ` AND `))
		}

		tags[helper.TextQuery] = query

		stmt, err := mr.ReadDB.Prepare(query)

		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, params)
			output <- ResultRepository{Error: err}
			return
		}
		defer stmt.Close()

		var (
			row model.B2CMerchantEmployeeData
		)

		err = stmt.QueryRow(queryValues...).Scan(
			&row.ID, &row.MerchantID, &row.MemberID,
			&row.FirstName, &row.LastName, &row.Email, &row.Gender, &row.Mobile, &row.Phone, &row.BirthDate,
			&row.CreatedAt, &row.ModifiedAt, &row.Status, &row.ProfilePicture,
			&row.MerchantLogo, &row.MerchantName, &row.MerchantType,
			&row.VanityURL, &row.IsActive, &row.IsPKP)

		if err != nil {
			output <- ResultRepository{Error: err}
			return
		}

		row.BirthDateString.String = row.BirthDate.Format("02/01/2006")

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, params)
			output <- ResultRepository{Error: err}
			return
		}

		if _, err := json.Marshal(&row); err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, params)
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: row}
	})

	return output
}
