package repo

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
	"github.com/lib/pq"
	"github.com/spf13/cast"
	"gopkg.in/guregu/null.v4"
)

const (
	DocSKK                = "SKKEMENKUMHAM-file"
	DocSKDP               = "SKDP-file"
	DocSertifikatMerek    = "SertifikatMerek-file"
	DocServiceCenter      = "ServiceCenter-file"
	DocSertifikatKeahlian = "SertifikatKeahlian-file"
	DocSuratIjin          = "SuratIjin-file"
)

// MerchantDocumentRepoPostgres data structure
type MerchantDocumentRepoPostgres struct {
	*repository.Repository
}

// NewMerchantDocumentRepoPostgres function for initializing  repo
func NewMerchantDocumentRepoPostgres(repo *repository.Repository) *MerchantDocumentRepoPostgres {
	return &MerchantDocumentRepoPostgres{repo}
}

// Save function for saving  data
func (mr *MerchantDocumentRepoPostgres) Save(merchantDocument model.B2CMerchantDocument) error {
	ctx := "MerchantDocumentRepo-create"
	query := `INSERT INTO b2c_merchantdocument
				(
					"id", "merchantId", "documentType", "documentValue", "documentExpirationDate",
					"creatorId", "creatorIp", "editorId", "editorIp", "version", "created", "lastModified"
				)
			VALUES
				(
					$1, $2, $3, $4, $5, $6, $7, 
					$8, $9, $10, $11, $12
				)
			ON CONFLICT(id)
			DO UPDATE SET
			"merchantId"=$2, "documentType"=$3, "documentValue"=$4, "documentExpirationDate"=$5,
			"creatorId"=$6, "creatorIp"=$7, "editorId"=$8, "editorIp"=$9, "version"=$10, "created"=$11, "lastModified"=$12`

	stmt, err := mr.WriteDB.Prepare(query)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		merchantDocument.ID, merchantDocument.MerchantID, merchantDocument.DocumentType, merchantDocument.DocumentValue,
		merchantDocument.DocumentExpirationDate, merchantDocument.CreatorID, merchantDocument.CreatorIP,
		merchantDocument.EditorID, merchantDocument.EditorIP, merchantDocument.Version,
		merchantDocument.Created, merchantDocument.LastModified,
	)

	if err != nil {
		helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, merchantDocument)
		return err
	}

	return nil
}

// Delete function for delete merchant data
func (mr *MerchantDocumentRepoPostgres) Delete(id string) error {
	ctx := "MerchantDocumentRepo-delete"

	query := `DELETE FROM b2c_merchantdocument WHERE id=$1;`

	stmt, err := mr.WriteDB.Prepare(query)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(id)

	if err != nil {
		helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, id)
		return err
	}

	return nil
}

// FindMerchantDocumentByParam function for getting detail merchant by slug
func (mr *MerchantDocumentRepoPostgres) FindMerchantDocumentByParam(ctxReq context.Context, param *model.B2CMerchantDocumentQueryInput) <-chan ResultRepository {
	ctx := "MerchantDocumentRepo-FindMerchantDocumentByParam"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)
		tags["documentType"] = param.DocumentType
		var (
			queryParam string
			queryList  []string
		)

		if param.MerchantID != "" {
			queryList = append(queryList, fmt.Sprintf(`"merchantId" = '%s'`, param.MerchantID))
		}

		if param.DocumentType != "" {
			queryList = append(queryList, fmt.Sprintf(`"documentType" = '%s'`, param.DocumentType))
		}

		if len(queryList) > 0 {
			queryParam = fmt.Sprintf(" WHERE %s", strings.Join(queryList, " AND "))
		}

		q := fmt.Sprintf(`SELECT id, "merchantId", "documentType", "documentValue", "documentExpirationDate",
			"creatorId", "creatorIp", "editorId", "editorIp", "version", "created", "lastModified" 
			FROM b2c_merchantdocument %s`, queryParam)

		tags[helper.TextQuery] = q
		stmt, err := mr.ReadDB.Prepare(q)

		if err != nil {
			tags[helper.TextResponse] = err
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, param)
			output <- ResultRepository{Error: err}
			return
		}
		defer stmt.Close()

		var (
			merchantDocument       model.B2CMerchantDocumentData
			documentExpirationDate pq.NullTime
			version                sql.NullInt64
		)

		err = stmt.QueryRow().Scan(
			&merchantDocument.ID, &merchantDocument.MerchantID, &merchantDocument.DocumentType, &merchantDocument.DocumentValue,
			&documentExpirationDate, &merchantDocument.CreatorID, &merchantDocument.CreatorIP,
			&merchantDocument.EditorID, &merchantDocument.EditorIP, &version,
			&merchantDocument.Created, &merchantDocument.LastModified,
		)

		if version.Valid {
			merchantDocument.Version = cast.ToInt(version.Int64)
		}

		if err != nil {
			tags[helper.TextResponse] = err
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, param)
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: merchantDocument}
	})
	return output
}

// InsertNewMerchantDocument function for insert data
func (mr *MerchantDocumentRepoPostgres) InsertNewMerchantDocument(ctxReq context.Context, merchantDocument *model.B2CMerchantDocumentData) <-chan ResultRepository {
	ctx := "MerchantDocumentRepo-create"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)
		var (
			stmt *sql.Stmt
			err  error
		)

		tags["merchantId"] = merchantDocument.MerchantID
		tags["docId"] = merchantDocument.ID

		query := `INSERT INTO b2c_merchantdocument
					(
						"id", "merchantId", "documentType", "documentValue", "documentExpirationDate",
						"creatorId", "creatorIp", "editorId", "editorIp", "version", "created", "lastModified"
					)
				VALUES
					(
						$1, $2, $3, $4, $5, $6, $7, 
						$8, $9, $10, $11, $12

					)
				ON CONFLICT(id)
				DO UPDATE SET
				"id"=$1, "merchantId"=$2, "documentType"=$3, "documentValue"=$4, "documentExpirationDate"=$5,
				"creatorId"=$6, "creatorIp"=$7, "editorId"=$8, "editorIp"=$9, "version"=$10, "created"=$11, "lastModified"=$12`

		if mr.Tx != nil {
			stmt, err = mr.Tx.Prepare(query)
		} else {
			stmt, err = mr.WriteDB.Prepare(query)
		}

		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextPrepareDatabase, err, merchantDocument)
			output <- ResultRepository{Error: err}
			return
		}

		_, err = stmt.Exec(
			merchantDocument.ID, merchantDocument.MerchantID, merchantDocument.DocumentType, merchantDocument.DocumentValue,
			&merchantDocument.DocumentExpirationDate, merchantDocument.CreatorID, merchantDocument.CreatorIP,
			merchantDocument.EditorID, merchantDocument.EditorIP, merchantDocument.Version,
			merchantDocument.Created, merchantDocument.LastModified,
		)

		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, merchantDocument)
			output <- ResultRepository{Error: err}
			return
		}
		output <- ResultRepository{Error: nil}
	})
	return output
}

// UpdateMerchantDocument function for update merchant document data
func (mr *MerchantDocumentRepoPostgres) UpdateMerchantDocument(ctxReq context.Context, id string, merchantDocument *model.B2CMerchantDocumentData) <-chan ResultRepository {
	ctx := "MerchantDocumentRepo-update"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)
		var (
			stmt *sql.Stmt
			err  error
		)
		tags["docId"] = id

		query := `UPDATE b2c_merchantdocument SET "merchantId"=$2, "documentType"=$3, "documentValue"=$4,
		"editorId"=$5, "editorIp"=$6, "lastModified"=$7 WHERE "id"=$1;`

		if mr.Tx != nil {
			stmt, err = mr.Tx.Prepare(query)
		} else {
			stmt, err = mr.WriteDB.Prepare(query)
		}
		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextPrepareDatabase, err, merchantDocument)
			output <- ResultRepository{Error: err}
			return
		}

		_, err = stmt.Exec(
			id, merchantDocument.MerchantID, merchantDocument.DocumentType, merchantDocument.DocumentValue,
			merchantDocument.EditorID, merchantDocument.EditorIP, merchantDocument.LastModified,
		)

		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, merchantDocument)
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Error: nil}
	})
	return output
}

// SaveMerchantDocumentGWS function for saving  data
func (mr *MerchantDocumentRepoPostgres) SaveMerchantDocumentGWS(merchantDocument model.B2CMerchantDocument) error {
	ctx := "MerchantDocumentRepo-SaveMerchantDocumentGWS"

	query := `INSERT INTO b2c_merchantdocument
				(
					"id", "merchantId", "documentType", "documentValue", "documentExpirationDate",
					"creatorId", "creatorIp", "editorId", "editorIp", "version", "created", "lastModified"
				)
			VALUES
				(
					$1, $2, $3, $4, $5, $6, $7, 
					$8, $9, $10, $11, $12

				)
			ON CONFLICT(id)
			DO UPDATE SET
			"id"=$1, "merchantId"=$2, "documentType"=$3, "documentValue"=$4, "documentExpirationDate"=$5,
			"creatorId"=$6, "creatorIp"=$7, "editorId"=$8, "editorIp"=$9, "version"=$10, "created"=$11, "lastModified"=$12`

	stmt, err := mr.WriteDB.Prepare(query)

	if err != nil {
		helper.SendErrorLog(context.Background(), ctx, helper.TextPrepareDatabase, err, merchantDocument)
		return err
	}

	var (
		documentExpirationDate sql.NullString
	)

	if len(*merchantDocument.DocumentExpirationDate) > 0 {
		documentExpirationDate.Valid = true
		documentExpirationDate.String = *merchantDocument.DocumentExpirationDate
	}

	_, err = stmt.Exec(
		merchantDocument.ID, merchantDocument.MerchantID, merchantDocument.DocumentType, merchantDocument.DocumentValue,
		documentExpirationDate, merchantDocument.CreatorID, merchantDocument.CreatorIP,
		merchantDocument.EditorID, merchantDocument.EditorIP, merchantDocument.Version,
		merchantDocument.Created, merchantDocument.LastModified,
	)

	if err != nil {
		helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, merchantDocument)
		return err
	}

	return nil
}

// GetListMerchantDocument function for loading merchant bank data
func (mr *MerchantDocumentRepoPostgres) GetListMerchantDocument(ctxReq context.Context, params *model.B2CMerchantDocumentQueryInput) <-chan ResultRepository {
	ctx := "MerchantDocumentRepo-GetListMerchantDocument"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)

		q := fmt.Sprintf(`SELECT "id", "merchantId", "documentType", "documentValue", "documentExpirationDate",
		"creatorId", "creatorIp", "editorId", "editorIp", "version", "created", "lastModified" 
		FROM b2c_merchantdocument WHERE "merchantId" ='%s'`, params.MerchantID)

		tags[helper.TextQuery] = q
		rows, err := mr.ReadDB.Query(q)
		if err != nil {
			tags[helper.TextResponse] = err
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, params)
			output <- ResultRepository{Error: nil}
			return
		}

		defer rows.Close()

		var listB2CMerchantDocument model.ListB2CMerchantDocument

		for rows.Next() {

			var (
				merchantDocument       model.B2CMerchantDocumentData
				documentExpirationDate null.Time
				version                sql.NullInt64
			)

			err = rows.Scan(
				&merchantDocument.ID, &merchantDocument.MerchantID, &merchantDocument.DocumentType, &merchantDocument.DocumentValue,
				&documentExpirationDate, &merchantDocument.CreatorID, &merchantDocument.CreatorIP,
				&merchantDocument.EditorID, &merchantDocument.EditorIP, &version,
				&merchantDocument.Created, &merchantDocument.LastModified,
			)

			if merchantDocument.DocumentExpirationDate.Valid {
				merchantDocument.DocumentExpirationDate = null.TimeFrom(documentExpirationDate.Time)
			}

			if version.Valid {
				merchantDocument.Version = cast.ToInt(version.Int64)
			}
			parseUrl, _ := url.Parse(merchantDocument.DocumentValue)
			merchantDocument.DocumentOriginal = strings.TrimLeft(parseUrl.Path, "/")

			if err != nil {
				tags[helper.TextResponse] = err
				helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, params)
				output <- ResultRepository{Error: err}
				return
			}

			listB2CMerchantDocument.MerchantDocument = append(listB2CMerchantDocument.MerchantDocument, merchantDocument)
		}
		tags[helper.TextResponse] = listB2CMerchantDocument
		output <- ResultRepository{Result: listB2CMerchantDocument}
	})

	return output

}

func (mr *MerchantDocumentRepoPostgres) ResetRejectedDocument(ctxReq context.Context, param model.B2CMerchantDocumentData) <-chan ResultRepository {
	ctx := "MerchantDocumentRepo-ResetRejectedDocument"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)
		var (
			stmt *sql.Stmt
			err  error
		)
		tags["merchantId"] = param.MerchantID

		query := `UPDATE b2c_merchantdocument SET "documentValue"='',
		"editorId"=$2, "editorIp"=$3, "lastModified"=$10 WHERE "merchantId"=$1 AND "documentType" IN ($4, $5, $6, $7, $8, $9);`

		if mr.Tx != nil {
			stmt, err = mr.Tx.Prepare(query)
		} else {
			stmt, err = mr.WriteDB.Prepare(query)
		}
		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextPrepareDatabase, err, param.MerchantID)
			output <- ResultRepository{Error: err}
			return
		}

		if _, err = stmt.Exec(
			param.MerchantID, param.EditorID, param.EditorIP,
			DocSKK, DocSKDP, DocSertifikatMerek, DocServiceCenter, DocSertifikatKeahlian, DocSuratIjin,
			time.Now(),
		); err != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, param.MerchantID)
			output <- ResultRepository{Error: err}
			return
		}
		output <- ResultRepository{Error: nil}
	})
	return output
}
