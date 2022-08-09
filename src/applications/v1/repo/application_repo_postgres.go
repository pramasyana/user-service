package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/applications/v1/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

// ApplicationRepoPostgres data structure
type ApplicationRepoPostgres struct {
	*repository.Repository
}

// NewApplicationRepoPostgres function for initializing Application repo
func NewApplicationRepoPostgres(repo *repository.Repository) *ApplicationRepoPostgres {
	return &ApplicationRepoPostgres{repo}
}

// Save function for saving application data
func (mr *ApplicationRepoPostgres) Save(ctxReq context.Context, application model.Application) <-chan ResultRepository {
	ctx := "ApplicationRepo-Save"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {

		var appSave model.Application

		query := `INSERT INTO "applications" ("name",
					"url", "logo", "created", "lastModified")
				VALUES
					($1, $2, $3, $4, $5)
				RETURNING
					"id", "name", "url", "logo", "created", "lastModified"`

		tags[helper.TextQuery] = query
		stmt, err := mr.WriteDB.Prepare(query)
		if err != nil {
			tracer.Log(ctxReq, helper.TextExecQuery, err)
			output <- ResultRepository{Error: err}
			return
		}

		err = stmt.QueryRow(application.Name, application.URL, application.Logo,
			application.Created, application.LastModified,
		).Scan(
			&appSave.ID,
			&appSave.Name,
			&appSave.URL,
			&appSave.Logo,
			&appSave.Created,
			&appSave.LastModified,
		)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextStmtError, err, application)
			output <- ResultRepository{Error: err}
			return
		}
		tags[helper.TextResponse] = appSave

		output <- ResultRepository{Result: appSave}
	})
	return output
}

// Update function for update application data
func (mr *ApplicationRepoPostgres) Update(ctxReq context.Context, application model.Application) <-chan ResultRepository {
	ctx := "ApplicationRepo-Update"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {

		var appUpdate model.Application

		query := `UPDATE "applications" SET 
					"name" = $2, "url" = $3, "logo" = $4, "lastModified" = $5
				WHERE
					"id" = $1
				RETURNING
					"id", "name", "url", "logo", "created", "lastModified"`

		tags[helper.TextQuery] = query
		stmt, err := mr.WriteDB.Prepare(query)
		if err != nil {
			tracer.Log(ctxReq, helper.TextExecQuery, err)
			output <- ResultRepository{Error: err}
			return
		}

		application.LastModified = time.Now()

		if err = stmt.QueryRow(
			application.ID, application.Name, application.URL, application.Logo, application.LastModified,
		).Scan(
			&appUpdate.ID,
			&appUpdate.Name,
			&appUpdate.URL,
			&appUpdate.Logo,
			&appUpdate.Created,
			&appUpdate.LastModified,
		); err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextStmtError, err, application)
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: appUpdate}
	})
	return output
}

// FindApplicationByID function for getting detail shipping address by id
func (mr *ApplicationRepoPostgres) FindApplicationByID(ctxReq context.Context, id string) <-chan ResultRepository {
	ctx := "ApplicationRepo-FindApplicationByID"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {

		query := `SELECT "id", "name", "url", "logo", "created",
		"lastModified" FROM "applications" WHERE "id" = $1`

		stmt, err := mr.ReadDB.Prepare(query)
		tags[helper.TextQuery] = query
		if err != nil {
			if err != sql.ErrNoRows {
				helper.SendErrorLog(ctxReq, ctx, helper.TextStmtError, err, id)
			}
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}
		defer stmt.Close()

		var appGet model.Application

		err = stmt.QueryRow(id).Scan(
			&appGet.ID,
			&appGet.Name,
			&appGet.URL,
			&appGet.Logo,
			&appGet.Created,
			&appGet.LastModified,
		)

		if err != nil {
			if err != sql.ErrNoRows {
				helper.SendErrorLog(ctxReq, ctx, helper.TextStmtError, err, id)
			}
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		tags["args"] = appGet
		output <- ResultRepository{Result: appGet}
	})
	return output
}

// Delete function for delete Application data
func (mr *ApplicationRepoPostgres) Delete(ctxReq context.Context, id string) <-chan ResultRepository {
	ctx := "ApplicationRepo-Delete"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		query := `DELETE FROM applications WHERE id=$1;`
		tags[helper.TextQuery] = query
		stmt, err := mr.WriteDB.Prepare(query)

		if err != nil {
			tracer.Log(ctxReq, helper.TextExecQuery, err)
			output <- ResultRepository{Error: err}
			return
		}

		_, err = stmt.Exec(id)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextStmtError, err, id)
			output <- ResultRepository{Error: err}
			return
		}
		output <- ResultRepository{Result: nil}
	})
	return output
}

// GetListApplication function for loading shipping address
func (mr *ApplicationRepoPostgres) GetListApplication(ctxReq context.Context, params *model.ParametersApplication) <-chan ResultRepository {
	ctx := "ApplicationRepo-GetListApplication"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		query := fmt.Sprintf(`SELECT "id", "name", "url",  "logo", "created",
		"lastModified" FROM "applications"
		LIMIT %d OFFSET %d`, params.Limit, params.Offset)

		tags[helper.TextQuery] = query
		rows, err := mr.ReadDB.Query(query)
		if err != nil {
			if err != sql.ErrNoRows {
				helper.SendErrorLog(ctxReq, ctx, helper.TextStmtError, err, params)
			}
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
		}

		defer rows.Close()

		var listApplication model.ListApplication

		for rows.Next() {
			var appList model.Application
			err = rows.Scan(
				&appList.ID,
				&appList.Name,
				&appList.URL,
				&appList.Logo,
				&appList.Created,
				&appList.LastModified,
			)

			if err != nil {
				if err != sql.ErrNoRows {
					helper.SendErrorLog(ctxReq, ctx, helper.TextStmtError, err, params)
				}
				tags[helper.TextResponse] = err
				output <- ResultRepository{Error: err}
			}

			listApplication.Application = append(listApplication.Application, &appList)
		}

		tags[helper.TextResponse] = listApplication
		output <- ResultRepository{Result: listApplication}
	})

	return output

}

// GetTotalApplication function for getting total of m
func (mr *ApplicationRepoPostgres) GetTotalApplication(ctxReq context.Context, params *model.ParametersApplication) <-chan ResultRepository {
	ctx := "ApplicationRepo-GetTotalApplication"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		var totalData int
		sq := `SELECT count(id) FROM applications`

		tags[helper.TextQuery] = sq
		err := mr.ReadDB.QueryRow(sq).Scan(&totalData)
		if err != nil {
			if err != sql.ErrNoRows {
				helper.SendErrorLog(ctxReq, ctx, helper.TextStmtError, err, params)
			}
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
		}

		tags[helper.TextResponse] = totalData
		output <- ResultRepository{Result: totalData}
	})

	return output
}
