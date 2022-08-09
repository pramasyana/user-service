package query

import (
	"context"
	"database/sql"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/phone_area/v1/model"
)

// PhoneAreaQueryPostgres data structure
type PhoneAreaQueryPostgres struct {
	db *sql.DB
}

// NewPhoneAreaQueryPostgres function for initializing member query
func NewPhoneAreaQueryPostgres(db *sql.DB) *PhoneAreaQueryPostgres {
	return &PhoneAreaQueryPostgres{db: db}
}

// FindAll function for getting list of phone area
func (q *PhoneAreaQueryPostgres) FindAll(ctxReq context.Context) <-chan ResultQuery {
	ctx := "PhoneAreaQuery-FindAll"

	output := make(chan ResultQuery)

	query := `SELECT "codeArea", "areaName", "provinceName" FROM phone_area`

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tags[helper.TextQuery] = query

		rows, err := q.db.Query(query)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
			output <- ResultQuery{Error: err}
			return
		}
		defer rows.Close()

		var phoneArea []model.PhoneArea

		for rows.Next() {
			var pa model.PhoneArea

			err = rows.Scan(&pa.CodeArea, &pa.AreaName, &pa.ProvinceName)
			if err != nil {
				helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
				output <- ResultQuery{Error: err}
				return
			}

			phoneArea = append(phoneArea, pa)
		}

		tags["args"] = phoneArea
		output <- ResultQuery{Result: phoneArea}
	})

	return output
}

// Count function for getting total phone area
func (q *PhoneAreaQueryPostgres) Count(ctxReq context.Context) <-chan ResultQuery {
	ctx := "PhoneAreaQuery-Count"

	output := make(chan ResultQuery)

	query := `SELECT COUNT(*) FROM phone_area`

	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)

		tags[helper.TextQuery] = query

		var res model.TotalPhoneArea
		err := q.db.QueryRow(query).Scan(&res.TotalData)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
			output <- ResultQuery{Error: err}
			return
		}

		tags["args"] = res
		output <- ResultQuery{Result: res}
	})

	return output
}
