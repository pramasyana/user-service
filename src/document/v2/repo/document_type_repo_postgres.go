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
	"github.com/Bhinneka/user-service/src/document/v2/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

// DocumentTypeRepoPostgres data structure
type DocumentTypeRepoPostgres struct {
	*repository.Repository
}

// NewDocumentTypeRepoPostgres function for initializing Document repo
func NewDocumentTypeRepoPostgres(repo *repository.Repository) *DocumentTypeRepoPostgres {
	return &DocumentTypeRepoPostgres{repo}
}

// AddDocumentType function for saving document data
func (mr *DocumentTypeRepoPostgres) AddDocumentType(ctxReq context.Context, document model.DocumentType) <-chan ResultRepository {
	ctx := "DocumentTypeRepo-AddDocumentType"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {

		var (
			docData               model.DocumentType
			createdBy, modifiedBy sql.NullString
		)

		document.Created = time.Now()
		document.LastModified = time.Now()

		query := `INSERT INTO "document_type" ("documentType", "isB2c", "isB2b", 
				"isActive", "created", "lastModified", "createdBy",  "modifiedBy")
				VALUES
					($1, $2, $3, $4, $5, $6, $7, $8)
				RETURNING
					"id", "documentType", "isB2c", "isB2b", 
					"isActive", "created", "lastModified", "createdBy",  "modifiedBy"`
		tags[helper.TextQuery] = query
		stmt, err := mr.WriteDB.Prepare(query)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		if err = stmt.QueryRow(
			document.DocumentType, document.IsB2c, document.IsB2b,
			document.IsActive, document.Created, document.LastModified, document.CreatedBy, document.ModifiedBy,
		).Scan(
			&docData.ID,
			&docData.DocumentType,
			&docData.IsB2c,
			&docData.IsB2b,
			&docData.IsActive,
			&docData.Created,
			&docData.LastModified,
			&createdBy,
			&modifiedBy,
		); err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		if createdBy.Valid && len(createdBy.String) > 0 {
			docData.CreatedBy = createdBy.String
		}

		if modifiedBy.Valid && len(modifiedBy.String) > 0 {
			docData.ModifiedBy = modifiedBy.String
		}

		output <- ResultRepository{Result: docData}
	})
	return output
}

// UpdateDocumentType function for update document data
func (mr *DocumentTypeRepoPostgres) UpdateDocumentType(ctxReq context.Context, document model.DocumentType) <-chan ResultRepository {
	ctx := "DocumentTypeRepo-UpdateDocumentType"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		var (
			doc                   model.DocumentType
			createdBy, modifiedBy sql.NullString
		)

		document.LastModified = time.Now()

		query := `UPDATE "document_type" SET 
					"documentType" = $2, "isB2c" = $3, "isB2b" = $4,  "isActive" = $5,  "lastModified" = $6, "modifiedBy" = $7
				WHERE 
					"id" = $1
				RETURNING
					"id", "documentType", "isB2c", "isB2b", "isActive", "created",
					"lastModified", "createdBy", "modifiedBy"`
		tags[helper.TextQuery] = query
		stmt, err := mr.WriteDB.Prepare(query)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		if err = stmt.QueryRow(
			document.ID, document.DocumentType, document.IsB2c, document.IsB2b,
			document.IsActive, document.LastModified, document.ModifiedBy,
		).Scan(
			&doc.ID,
			&doc.DocumentType,
			&doc.IsB2c,
			&doc.IsB2b,
			&doc.IsActive,
			&doc.Created,
			&doc.LastModified,
			&createdBy,
			&modifiedBy,
		); err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		if createdBy.Valid && len(createdBy.String) > 0 {
			doc.CreatedBy = createdBy.String
		}

		if modifiedBy.Valid && len(modifiedBy.String) > 0 {
			doc.ModifiedBy = modifiedBy.String
		}

		output <- ResultRepository{Result: doc}
	})
	return output
}

// FindDocumentTypeByParam function for getting detail document by slug
func (mr *DocumentTypeRepoPostgres) FindDocumentTypeByParam(ctxReq context.Context, param *model.DocumentTypeParameters) <-chan ResultRepository {
	ctx := "DocumentTypeRepoPostgres-FindDocumentTypeByParam"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		var (
			queryParam string
			queryList  []string
		)

		if param.ID != "" {
			queryList = append(queryList, fmt.Sprintf(`"id" = '%s'`, param.ID))
		}

		if param.DocumentType != "" {
			queryList = append(queryList, fmt.Sprintf(`"documentType" = '%s'`, param.DocumentType))
		}

		if param.IsActive != "" {
			isActive, _ := strconv.ParseBool(param.IsActive)
			queryList = append(queryList, fmt.Sprintf(`"isActive" = %t`, isActive))
		}

		if len(queryList) > 0 {
			queryParam = fmt.Sprintf(" WHERE %s", strings.Join(queryList, " AND "))
		}

		q := fmt.Sprintf(`SELECT id, "documentType", "isB2c", "isB2b", 
			"isActive", "created", "lastModified", "createdBy",  "modifiedBy" 
			FROM document_type %s`, queryParam)

		tags[helper.TextQuery] = q
		stmt, err := mr.ReadDB.Prepare(q)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, param)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}
		defer stmt.Close()

		var documentType model.DocumentType

		err = stmt.QueryRow().Scan(
			&documentType.ID, &documentType.DocumentType, &documentType.IsB2c,
			&documentType.IsB2b, &documentType.IsActive, &documentType.Created,
			&documentType.LastModified, &documentType.CreatedBy, &documentType.ModifiedBy,
		)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, param)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: documentType}
	})
	return output
}

// GetListDocumentType function for loading document
func (mr *DocumentTypeRepoPostgres) GetListDocumentType(ctxReq context.Context, params *model.DocumentTypeParameters) <-chan ResultRepository {
	ctx := "DocumentTypeRepoPostgres-GetListDocumentType"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)

		queryParam := mr.generateQuery(params)

		sq := fmt.Sprintf(`SELECT "id", "documentType", "isB2c", "isB2b", "isActive", 
						"created", "lastModified", "createdBy",  "modifiedBy" FROM "document_type" %s
						LIMIT %d OFFSET %d`, queryParam, params.Limit, params.Offset)
		tags[helper.TextQuery] = sq
		rows, err := mr.ReadDB.Query(sq)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, params)
			output <- ResultRepository{Error: nil}
			return
		}

		defer rows.Close()

		var listDocument model.ListDocumentType

		for rows.Next() {
			var (
				document              model.DocumentType
				modifiedBy, createdBy sql.NullString
			)
			err = rows.Scan(
				&document.ID,
				&document.DocumentType,
				&document.IsB2c,
				&document.IsB2b,
				&document.IsActive,
				&document.Created,
				&document.LastModified,
				&createdBy,
				&modifiedBy,
			)

			if modifiedBy.Valid && len(modifiedBy.String) > 0 {
				document.ModifiedBy = modifiedBy.String
			}

			if createdBy.Valid && len(createdBy.String) > 0 {
				document.CreatedBy = createdBy.String
			}

			if err != nil {
				helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, params)
				output <- ResultRepository{Error: err}
				return
			}
			listDocument.DocumentType = append(listDocument.DocumentType, &document)
		}

		tags[helper.TextResponse] = listDocument
		output <- ResultRepository{Result: listDocument}
	})

	return output
}

// GetTotalDocumentType function for getting total of document
func (mr *DocumentTypeRepoPostgres) GetTotalDocumentType(ctxReq context.Context, params *model.DocumentTypeParameters) <-chan ResultRepository {
	ctx := "DocumentTypeRepoPostgres-GetTotalDocumentType"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)

		var totalData int

		queryParam := mr.generateQuery(params)

		sq := fmt.Sprintf(`SELECT count(id) FROM document_type %s`, queryParam)

		tags[helper.TextQuery] = sq
		err := mr.ReadDB.QueryRow(sq).Scan(&totalData)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextQuery, err, params)
			output <- ResultRepository{Error: err}
			return
		}

		tags[helper.TextResponse] = totalData
		output <- ResultRepository{Result: totalData}
	})

	return output
}

// generateQuery function for generating query
func (mr *DocumentTypeRepoPostgres) generateQuery(params *model.DocumentTypeParameters) string {
	var (
		queryParam string
		queryList  []string
	)
	if params.IsB2b != "" {
		isB2b, _ := strconv.ParseBool(params.IsB2b)
		queryList = append(queryList, fmt.Sprintf(`"isB2b" = %t`, isB2b))
	}

	if params.IsB2c != "" {
		isB2c, _ := strconv.ParseBool(params.IsB2c)
		queryList = append(queryList, fmt.Sprintf(`"isB2c" = %t`, isB2c))
	}

	if len(queryList) > 0 {
		queryParam = fmt.Sprintf(" WHERE %s", strings.Join(queryList, " AND "))
	}
	return queryParam
}
