package repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

// MerchantBankRepoPostgres data structure
type MerchantBankRepoPostgres struct {
	*repository.Repository
}

// NewMerchantBankRepoPostgres function for initializing  repo
func NewMerchantBankRepoPostgres(repo *repository.Repository) *MerchantBankRepoPostgres {
	return &MerchantBankRepoPostgres{repo}
}

// Save function for saving  data
func (mr *MerchantBankRepoPostgres) Save(merchantBank model.B2CMerchantBank) error {
	ctx := "MerchantBankRepo-create"

	query := `INSERT INTO b2c_merchantbank
				(
					"id", "bankCode", "bankName", "status", "creatorId",
					"creatorIp", "editorId", "editorIp", "lastModified", "created"
				)
			VALUES
				(
					$1, $2, $3, $4, $5, 
					$6, $7, $8, $9, $10

				)
			ON CONFLICT(id)
			DO UPDATE SET
			"id"=$1, "bankCode"=$2, "bankName"=$3, "status"=$4, "creatorId"=$5,
			"creatorIp"=$6, "editorId"=$7, "editorIp"=$8, "lastModified"=$9, "created"=$10`

	stmt, err := mr.WriteDB.Prepare(query)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		merchantBank.ID, merchantBank.BankCode, merchantBank.BankName, merchantBank.Status,
		merchantBank.CreatorID, merchantBank.CreatorIP, merchantBank.EditorID,
		merchantBank.EditorIP, merchantBank.LastModified, merchantBank.Created,
	)

	if err != nil {
		helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, merchantBank)
		return err
	}

	return nil
}

// Delete function for delete merchant bank data
func (mr *MerchantBankRepoPostgres) Delete(merchantBank model.B2CMerchantBank) error {
	ctx := "MerchantBankRepo-delete"

	query := `DELETE FROM b2c_merchantbank WHERE id=$1;`

	stmt, err := mr.WriteDB.Prepare(query)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(merchantBank.ID)

	if err != nil {
		helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, merchantBank)
		return err
	}

	return nil
}

// Load function for loading merchant bank data based on id
func (mr *MerchantBankRepoPostgres) Load(uid string) ResultRepository {
	ctx := "MerchantBankRepo-Load"

	q := `SELECT "id", "bankCode","bankName","status","creatorId","creatorIp",
	"editorId","editorIp","created","lastModified"  
	FROM b2c_merchantbank WHERE id = $1 AND "deletedAt" IS NULL`

	stmt, err := mr.ReadDB.Prepare(q)

	if err != nil {
		helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, uid)
		return ResultRepository{Error: err}
	}
	defer stmt.Close()

	var merchant model.B2CMerchantBank

	err = stmt.QueryRow(uid).Scan(
		&merchant.ID, &merchant.BankCode, &merchant.BankName, &merchant.Status,
		&merchant.CreatorID, &merchant.CreatorIP, &merchant.Created, &merchant.EditorID,
		&merchant.EditorIP, &merchant.LastModified,
	)

	if err != nil {
		helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, uid)
		return ResultRepository{Error: err}
	}

	return ResultRepository{Result: merchant}

}

// FindActiveMerchantBankByID function for loading merchant bank data based on id and status
func (mr *MerchantBankRepoPostgres) FindActiveMerchantBankByID(ctxReq context.Context, bankID int) <-chan ResultRepository {
	ctx := "MerchantBankRepo-Load"

	output := make(chan ResultRepository)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		q := `SELECT "id", "bankCode","bankName","status","creatorId","creatorIp",
		"editorId","editorIp","created","lastModified" 
		FROM b2c_merchantbank WHERE id = $1 AND status=$2 AND "deletedAt" IS NULL`

		stmt, err := mr.ReadDB.Prepare(q)

		if err != nil {
			tags[helper.TextResponse] = err
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, q)
			output <- ResultRepository{Error: err}
			return
		}
		defer stmt.Close()
		tags[helper.TextQuery] = q

		var merchant model.B2CMerchantBankData

		err = stmt.QueryRow(bankID, true).Scan(
			&merchant.ID, &merchant.BankCode, &merchant.BankName, &merchant.Status,
			&merchant.CreatorID, &merchant.CreatorIP, &merchant.Created, &merchant.EditorID,
			&merchant.EditorIP, &merchant.LastModified,
		)

		if err != nil {
			tags[helper.TextResponse] = err
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, q)
			output <- ResultRepository{Error: err}
			return
		}

		tags["args"] = merchant
		output <- ResultRepository{Result: merchant}
	})

	return output

}

// GetListMerchantBank function for loading merchant bank data
func (mr *MerchantBankRepoPostgres) GetListMerchantBank(params *model.ParametersMerchantBank) <-chan ResultRepository {
	ctx := "MerchantBankRepo-GetListMerchantBank"
	output := make(chan ResultRepository)
	go func() {
		defer close(output)

		if len(params.OrderBy) > 0 {
			params.OrderBy = fmt.Sprintf(`"%s"`, params.OrderBy)
		}

		queryParam := ""
		if params.Status != "" {
			queryParam = ` AND "status" = ` + params.Status
		}
		sq := fmt.Sprintf(`SELECT "id", "bankCode","bankName","status","creatorId","creatorIp",
		"editorId","editorIp","created","lastModified" 
		FROM b2c_merchantbank WHERE "deletedAt" IS NULL %s ORDER BY %s %s
		LIMIT %d OFFSET %d`, queryParam, params.OrderBy, params.Sort, params.Limit, params.Offset)

		rows, err := mr.ReadDB.Query(sq)
		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, params)
			output <- ResultRepository{Error: nil}
			return
		}

		defer rows.Close()

		var listMerchantBank model.ListMerchantBank

		for rows.Next() {
			var merchant model.B2CMerchantBankData
			err = rows.Scan(
				&merchant.ID, &merchant.BankCode, &merchant.BankName, &merchant.Status,
				&merchant.CreatorID, &merchant.CreatorIP, &merchant.Created, &merchant.EditorID,
				&merchant.EditorIP, &merchant.LastModified,
			)

			if err != nil {
				helper.SendErrorLog(context.Background(), ctx, helper.TextQueryDatabase, err, params)
				output <- ResultRepository{Error: err}
				return
			}

			listMerchantBank.MerchantBank = append(listMerchantBank.MerchantBank, &merchant)
		}
		output <- ResultRepository{Result: listMerchantBank}
	}()

	return output

}

// GetTotalMerchantBank function for getting total of merchant bank
func (mr *MerchantBankRepoPostgres) GetTotalMerchantBank(params *model.ParametersMerchantBank) <-chan ResultRepository {
	ctx := "MemberQuery-GetTotalMerchantBank"

	output := make(chan ResultRepository)
	go func() {
		defer close(output)

		queryParam := ""
		if params.Status != "" {
			queryParam = `AND "status" = ` + params.Status
		}
		var totalData int
		sq := fmt.Sprintf(`SELECT count(id) FROM b2c_merchantbank WHERE "deletedAt" IS NULL %s`, queryParam)
		err := mr.ReadDB.QueryRow(sq).Scan(&totalData)
		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, params)
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: totalData}
	}()

	return output
}

// SaveMasterBankGWS function for saving  data
func (mr *MerchantBankRepoPostgres) SaveMasterBankGWS(ctxReq context.Context, merchantBank model.B2CMerchantBank) error {
	ctx := "MerchantBankRepo-SaveMasterBankGWS"
	tr := tracer.StartTrace(ctxReq, ctx)
	tags := make(map[string]interface{})
	defer func() {
		tr.Finish(tags)
	}()

	query := `INSERT INTO b2c_merchantbank
				(
					"id", "bankCode", "bankName", "status", "creatorId",
					"creatorIp", "editorId", "editorIp", "lastModified", "created", "deletedAt"
				)
			VALUES
				(
					$1, $2, $3, $4, $5, 
					$6, $7, $8, $9, $10, $11
				)
			ON CONFLICT(id)
			DO UPDATE SET
			"id"=$1, "bankCode"=$2, "bankName"=$3, "status"=$4, "creatorId"=$5,
			"creatorIp"=$6, "editorId"=$7, "editorIp"=$8, "lastModified"=$9, "created"=$10, "deletedAt"=$11`

	tags[helper.TextQuery] = query
	stmt, err := mr.WriteDB.Prepare(query)
	if err != nil {
		tags[helper.TextResponse] = err
		return err
	}

	var deletedAt sql.NullString

	if len(*merchantBank.DeletedAt) > 0 {
		deletedAt.Valid = true
		deletedAt.String = *merchantBank.DeletedAt
	}

	tags[helper.TextParameter] = merchantBank
	_, err = stmt.Exec(
		merchantBank.ID, merchantBank.BankCode, merchantBank.BankName, merchantBank.Status,
		merchantBank.CreatorID, merchantBank.CreatorIP, merchantBank.EditorID,
		merchantBank.EditorIP, merchantBank.LastModified, merchantBank.Created, deletedAt,
	)

	if err != nil {
		tags[helper.TextResponse] = err
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, merchantBank)
		return err
	}

	tags[helper.TextResponse] = "success"

	return nil
}
