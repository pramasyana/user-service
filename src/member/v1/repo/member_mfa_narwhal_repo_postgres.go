package repo

import (
	"context"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
)

// EnableNarwhalMFA function for updating status only
func (mr *MemberMFARepoPostgres) EnableNarwhalMFA(ctxReq context.Context, memberID string, mfaKey string) <-chan ResultRepository {
	ctx := "MemberQuery-EnableNarwhalMFA"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)

		lastMfaAdminEnabled := time.Now().Format(time.RFC3339)
		query := `UPDATE member SET "mfaAdminEnabled" = true, "mfaAdminKey"= $1, "lastMfaAdminEnabled" = $2 WHERE id = $3`

		stmt, err := mr.WriteDB.Prepare(query)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, "enable_mfa_admin", err, memberID)
			output <- ResultRepository{Error: err}
			return
		}
		tags[helper.TextParameter] = memberID

		defer stmt.Close()

		if _, err = stmt.Exec(mfaKey, lastMfaAdminEnabled, memberID); err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, memberID)
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Error: nil}
	})

	return output
}

// DisableNarwhalMFA function for updating status only
func (mr *MemberMFARepoPostgres) DisableNarwhalMFA(ctxReq context.Context, memberID string) <-chan ResultRepository {
	ctx := "MemberQuery-DisableNarwhalMFA"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)
		tags[helper.TextParameter] = memberID
		query := `UPDATE member SET "mfaAdminEnabled" = false, "mfaAdminKey" = '', "lastMfaAdminEnabled" = NULL WHERE id = $1`

		stmt, err := mr.WriteDB.Prepare(query)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, memberID)
			output <- ResultRepository{Error: err}
			return
		}

		defer stmt.Close()

		if _, err = stmt.Exec(memberID); err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, memberID)
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Error: nil}
	})

	return output
}
