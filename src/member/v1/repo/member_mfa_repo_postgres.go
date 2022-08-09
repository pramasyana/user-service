package repo

import (
	"context"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

// MemberMFARepoPostgres data structure
type MemberMFARepoPostgres struct {
	*repository.Repository
}

// NewMemberMFARepoPostgres function for initializing member repo
func NewMemberMFARepoPostgres(repo *repository.Repository) *MemberMFARepoPostgres {
	return &MemberMFARepoPostgres{repo}
}

// MFAEnabled function for updating status only
func (mr *MemberMFARepoPostgres) MFAEnabled(ctxReq context.Context, memberID string, mfaKey string) <-chan ResultRepository {
	ctx := "MemberQuery-MFAEnabled"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)

		lastMfaEnabled := time.Now()
		query := `UPDATE member SET "mfaEnabled" = true, "mfaKey"= $1, "lastMfaEnabled" = $2 WHERE id = $3`
		tags[helper.TextQuery] = query
		stmt, err := mr.WriteDB.Prepare(query)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, "update_mfa_enable_member", err, memberID)
			output <- ResultRepository{Error: err}
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(mfaKey, lastMfaEnabled, memberID)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, memberID)
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Error: nil}
	})

	return output
}

// MFADisabled function for updating status only
func (mr *MemberMFARepoPostgres) MFADisabled(ctxReq context.Context, memberID string) <-chan ResultRepository {
	ctx := "MemberQuery-MFADisabled"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)
		query := `UPDATE member SET "mfaEnabled" = false, "mfaKey" = '', "lastMfaEnabled" = NULL WHERE id = $1`
		tags[helper.TextQuery] = query
		stmt, err := mr.WriteDB.Prepare(query)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, memberID)
			output <- ResultRepository{Error: err}
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(memberID)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, memberID)
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Error: nil}
	})

	return output
}
