package query

import (
	"context"
	"database/sql"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/lib/pq"
)

// MemberMFAQueryPostgres data structure
type MemberMFAQueryPostgres struct {
	db *sql.DB
}

// NewMemberMFAQueryPostgres function for initializing member query
func NewMemberMFAQueryPostgres(db *sql.DB) *MemberMFAQueryPostgres {
	return &MemberMFAQueryPostgres{db: db}
}

// FindMFASettings function for loading member data based on user id
func (mr *MemberMFAQueryPostgres) FindMFASettings(ctxReq context.Context, uid string) <-chan ResultQuery {
	ctx := "MemberMFAQuery-FindMFASettings"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		q := `SELECT "mfaEnabled", "lastMfaEnabled" FROM member WHERE id = $1`

		tags[helper.TextQuery] = q

		stmt, err := mr.db.Prepare(q)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, q)
			tags[helper.TextResponse] = err
			output <- ResultQuery{Error: err}
			return
		}
		defer stmt.Close()

		var (
			member         model.MFASettings
			lastMfaEnabled pq.NullTime
		)
		err = stmt.QueryRow(uid).Scan(
			&member.MfaEnabled,
			&lastMfaEnabled,
		)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, q)
			tags[helper.TextResponse] = err
			output <- ResultQuery{Error: err}
			return
		}

		if lastMfaEnabled.Valid {
			member.LastMfaEnabled = lastMfaEnabled.Time
			member.LastMfaEnabledString = member.LastMfaEnabled.Format(time.RFC3339)
		}

		tags["args"] = member
		output <- ResultQuery{Result: member}
	})
	return output
}

// FindNarwhalMFASettings function for loading member data based on user id
func (mr *MemberMFAQueryPostgres) FindNarwhalMFASettings(ctxReq context.Context, uid string) <-chan ResultQuery {
	ctx := "MemberMFAQuery-FindNarwhalMFASettings"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		q := `SELECT "mfaAdminEnabled", "lastMfaAdminEnabled" FROM member WHERE id = $1`

		stmt, err := mr.db.Prepare(q)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, q)
			output <- ResultQuery{Error: err}
			return
		}
		defer stmt.Close()

		var (
			admin               model.MFAAdminSettings
			lastMfaAdminEnabled pq.NullTime
		)
		err = stmt.QueryRow(uid).Scan(
			&admin.MfaAdminEnabled,
			&lastMfaAdminEnabled,
		)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, q)
			output <- ResultQuery{Error: err}
			return
		}

		if lastMfaAdminEnabled.Valid {
			admin.LastMfaAdminEnabled = lastMfaAdminEnabled.Time
			admin.LastMfaAdminEnabledString = admin.LastMfaAdminEnabled.Format(time.RFC3339)
		}

		tags["args"] = admin
		output <- ResultQuery{Result: admin}
	})
	return output
}
