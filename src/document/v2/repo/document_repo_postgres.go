package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/document/v2/model"
	"github.com/Bhinneka/user-service/src/service"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

// DocumentRepoPostgres data structure
type DocumentRepoPostgres struct {
	*repository.Repository
	UploadService service.UploadService
}

// NewDocumentRepoPostgres function for initializing Document repo
func NewDocumentRepoPostgres(repo *repository.Repository, uploadService *service.UploadService) *DocumentRepoPostgres {
	return &DocumentRepoPostgres{repo, *uploadService}
}

// AddDocument function for saving document data
func (mr *DocumentRepoPostgres) AddDocument(ctxReq context.Context, document model.DocumentData) <-chan ResultRepository {
	ctx := "DocumentRepoPostgres-AddDocument"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {

		var (
			newDoc                model.DocumentData
			createdBy, modifiedBy sql.NullString
		)

		query := `INSERT INTO "document" ("memberId", "documentType",
					"documentFile", "title", "number", "status", "statusText", "description",
					"isDelete", "created", "lastModified", "createdBy",  "modifiedBy")
				VALUES
					($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
				RETURNING
					"id", "memberId", "documentType", "documentFile", "title", "number", "status", "statusText", "description",
					"isDelete", "created", "lastModified", "createdBy",  "modifiedBy"`

		stmt, err := mr.WriteDB.Prepare(query)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		if err = stmt.QueryRow(
			document.MemberID, document.DocumentType, document.DocumentFile, document.Title,
			document.Number, document.Status, document.StatusText, document.Description,
			document.IsDelete, document.Created, document.LastModified, document.CreatedBy, document.ModifiedBy,
		).Scan(
			&newDoc.ID,
			&newDoc.MemberID,
			&newDoc.DocumentType,
			&newDoc.DocumentFile,
			&newDoc.Title,
			&newDoc.Number,
			&newDoc.Status,
			&newDoc.StatusText,
			&newDoc.Description,
			&newDoc.IsDelete,
			&newDoc.Created,
			&newDoc.LastModified,
			&createdBy,
			&modifiedBy,
		); err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		if createdBy.Valid && len(createdBy.String) > 0 {
			newDoc.CreatedBy = createdBy.String
		}

		if modifiedBy.Valid && len(modifiedBy.String) > 0 {
			newDoc.ModifiedBy = modifiedBy.String
		}

		output <- ResultRepository{Result: newDoc}
	})
	return output
}

// UpdateDocument function for update document data
func (mr *DocumentRepoPostgres) UpdateDocument(ctxReq context.Context, document model.DocumentData) <-chan ResultRepository {
	ctx := "DocumentRepoPostgres-UpdateDocument"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		var (
			docUpdate             model.DocumentData
			createdBy, modifiedBy sql.NullString
		)

		document.LastModified = time.Now()

		query := `UPDATE "document" SET  "documentType"  = $2,
						"documentFile" = $3, "title" = $4, "number" = $5, "status" = $6, "statusText" = $7, "description" = $8,
						"isDelete" = $9, "lastModified" = $10, "modifiedBy" = $11
					WHERE 
						"id" = $1
					RETURNING
						"id", "memberId", "documentType", "documentFile", "title", "number", "status", "statusText", "description",
						"isDelete", "created", "lastModified", "createdBy",  "modifiedBy"`
		tags[helper.TextQuery] = query
		stmt, err := mr.WriteDB.Prepare(query)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		if err = stmt.QueryRow(document.ID, document.DocumentType, document.DocumentFile, document.Title,
			document.Number, document.Status, document.StatusText, document.Description,
			document.IsDelete, document.LastModified, document.ModifiedBy,
		).Scan(
			&docUpdate.ID,
			&docUpdate.MemberID,
			&docUpdate.DocumentType,
			&docUpdate.DocumentFile,
			&docUpdate.Title,
			&docUpdate.Number,
			&docUpdate.Status,
			&docUpdate.StatusText,
			&docUpdate.Description,
			&docUpdate.IsDelete,
			&docUpdate.Created,
			&docUpdate.LastModified,
			&createdBy,
			&modifiedBy,
		); err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		if createdBy.Valid && len(createdBy.String) > 0 {
			docUpdate.CreatedBy = createdBy.String
		}

		if modifiedBy.Valid && len(modifiedBy.String) > 0 {
			docUpdate.ModifiedBy = modifiedBy.String
		}

		output <- ResultRepository{Result: docUpdate}
	})
	return output
}

// FindDocumentByParam function for getting detail document by slug
func (mr *DocumentRepoPostgres) FindDocumentByParam(ctxReq context.Context, params *model.DocumentParameters) <-chan ResultRepository {
	ctx := "DocumentRepoPostgres-FindDocumentByParam"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		queryParamFindByParam := mr.generateQuery(params)

		query := fmt.Sprintf(`SELECT "id", "memberId", "documentType", "documentFile", 
						"title", "number", "status", "statusText", "description",
						"isDelete", "created", "lastModified", "createdBy",  "modifiedBy"
						FROM document %s`, queryParamFindByParam)

		tags[helper.TextQuery] = query
		stmt, err := mr.ReadDB.Prepare(query)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, params)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}
		defer stmt.Close()

		var doc model.DocumentData

		err = stmt.QueryRow().Scan(
			&doc.ID,
			&doc.MemberID,
			&doc.DocumentType,
			&doc.DocumentFile,
			&doc.Title,
			&doc.Number,
			&doc.Status,
			&doc.StatusText,
			&doc.Description,
			&doc.IsDelete,
			&doc.Created,
			&doc.LastModified,
			&doc.CreatedBy,
			&doc.ModifiedBy,
		)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, params)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: doc}
	})
	return output
}

// GetListDocument function for loading document
func (mr *DocumentRepoPostgres) GetListDocument(ctxReq context.Context, params *model.DocumentParameters) <-chan ResultRepository {
	ctx := "DocumentRepoPostgres-GetListDocument"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		queryParamGetList := mr.generateQuery(params)

		sq := fmt.Sprintf(`SELECT "id", "memberId", "documentType", "documentFile", "title", "number", "status", "statusText", "description",
		"isDelete", "created", "lastModified", "createdBy",  "modifiedBy" FROM "document" 
		%s
		LIMIT %d OFFSET %d`, queryParamGetList, params.Limit, params.Offset)

		tags[helper.TextQuery] = sq
		rows, err := mr.ReadDB.Query(sq)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, params)
			output <- ResultRepository{Error: nil}
			return
		}

		defer rows.Close()

		var listDocument model.ListDocument

		for rows.Next() {
			var (
				document                            model.DocumentData
				modifiedBy, createdBy, documentFile sql.NullString
			)
			err = rows.Scan(
				&document.ID,
				&document.MemberID,
				&document.DocumentType,
				&documentFile,
				&document.Title,
				&document.Number,
				&document.Status,
				&document.StatusText,
				&document.Description,
				&document.IsDelete,
				&document.Created,
				&document.LastModified,
				&createdBy,
				&modifiedBy,
			)

			document.ModifiedBy = helper.ValidateSQLNullString(modifiedBy)
			document.CreatedBy = helper.ValidateSQLNullString(createdBy)

			if documentFile.Valid && len(documentFile.String) > 0 {
				document.DocumentFile = mr.getURLImage(ctxReq, documentFile.String)
			}

			if err != nil {
				helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, params)
				output <- ResultRepository{Error: err}
				return
			}

			listDocument.Document = append(listDocument.Document, &document)
		}

		tags[helper.TextResponse] = listDocument
		output <- ResultRepository{Result: listDocument}
	})

	return output
}

// GetDetailDocument function for getting detail document by slug
func (mr *DocumentRepoPostgres) GetDetailDocument(ctxReq context.Context, params *model.DocumentParameters) <-chan ResultRepository {
	ctx := "DocumentRepoPostgres-GetDetailDocument"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		queryParamGetDetail := mr.generateQuery(params)

		querygetDetail := fmt.Sprintf(`SELECT "id", "memberId", "documentType", "documentFile", 
						"title", "number", "status", "statusText", "description",
						"isDelete", "created", "lastModified", "createdBy",  "modifiedBy"
						FROM document %s`, queryParamGetDetail)

		tags[helper.TextQuery] = querygetDetail
		stmt, err := mr.ReadDB.Prepare(querygetDetail)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, params)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}
		defer stmt.Close()

		var (
			doc          model.DocumentData
			documentFile sql.NullString
		)

		err = stmt.QueryRow().Scan(
			&doc.ID,
			&doc.MemberID,
			&doc.DocumentType,
			&documentFile,
			&doc.Title,
			&doc.Number,
			&doc.Status,
			&doc.StatusText,
			&doc.Description,
			&doc.IsDelete,
			&doc.Created,
			&doc.LastModified,
			&doc.CreatedBy,
			&doc.ModifiedBy,
		)

		if documentFile.Valid && len(documentFile.String) > 0 {
			doc.DocumentFile = mr.getURLImage(ctxReq, documentFile.String)
		}

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, params)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: doc}
	})
	return output
}

// GetTotalDocument function for getting total of document
func (mr *DocumentRepoPostgres) GetTotalDocument(ctxReq context.Context, params *model.DocumentParameters) <-chan ResultRepository {
	ctx := "DocumentRepoPostgres-GetTotalDocument"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)

		var totalData int

		queryParamGetTotal := mr.generateQuery(params)

		sq := fmt.Sprintf(`SELECT count(id) FROM document %s`, queryParamGetTotal)

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

// DeleteDocumentByID function for delete document data
func (mr *DocumentRepoPostgres) DeleteDocumentByID(ctxReq context.Context, documentID string) <-chan ResultRepository {
	ctx := "DocumentRepoPostgres-DeleteDocumentByID"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		lastModified := time.Now()

		queryDelete := `UPDATE "document" SET   "lastModified" = $2, "isDelete" = true WHERE  "id" = $1`
		tags[helper.TextQuery] = queryDelete
		stmt, err := mr.WriteDB.Prepare(queryDelete)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, queryDelete)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		_, err = stmt.Exec(documentID, lastModified)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, documentID)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: nil}
	})
	return output
}

// generateQuery function for generating query
func (mr *DocumentRepoPostgres) generateQuery(params *model.DocumentParameters) string {
	var (
		queryParam string
		queryList  []string
	)
	queryList = append(queryList, `"isDelete" = false`)

	if params.ID != "" {
		paramID, _ := strconv.Atoi(params.ID)
		queryList = append(queryList, fmt.Sprintf(`"id" = %d`, paramID))
	}

	if params.MemberID != "" {
		queryList = append(queryList, fmt.Sprintf(`"memberId" = '%s'`, params.MemberID))
	}

	if len(queryList) > 0 {
		queryParam = fmt.Sprintf(" WHERE %s", strings.Join(queryList, " AND "))
	}
	return queryParam
}

// getURLImage function for get url image
func (mr *DocumentRepoPostgres) getURLImage(ctxReq context.Context, url string) string {
	ctx := "DocumentRepoPostgres-getUrlImage"
	isAttachment := "false"
	documentURL := <-mr.UploadService.GetURLImage(ctxReq, url, isAttachment)
	if documentURL.Result != nil {
		documentURLResult, ok := documentURL.Result.(serviceModel.ResponseUploadService)
		if !ok {
			err := errors.New("failed get url image")
			helper.SendErrorLog(ctxReq, ctx, "get_url_image", err, documentURLResult)
			return url
		}
		return documentURLResult.Data.URL
	}

	return url
}
