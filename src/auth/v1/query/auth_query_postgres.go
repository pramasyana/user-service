package query

import (
	"context"
	"database/sql"

	"github.com/Bhinneka/user-service/helper"
)

// AuthQueryPostgres data structure
type AuthQueryPostgres struct {
	db *sql.DB
}

// NewAuthQueryPostgres function for initializing auth query
func NewAuthQueryPostgres(db *sql.DB) *AuthQueryPostgres {
	return &AuthQueryPostgres{db: db}
}

//UpdateLastLogin for update last login date by USER ID
func (aqp *AuthQueryPostgres) UpdateLastLogin(uid string) <-chan ResultQuery {
	ctx := "AuthRepo-UpdateLastLogin"

	output := make(chan ResultQuery)
	go func() {
		defer close(output)

		tx, err := aqp.db.Begin()
		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextDBBegin, err, uid)
			output <- ResultQuery{Error: err}
			return
		}

		q := `UPDATE member SET "lastLogin" = now() WHERE id = $1`

		stmt, err := tx.Prepare(q)
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(context.Background(), ctx, helper.TextPrepareDatabase, err, q)
			output <- ResultQuery{Error: err}
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(uid)
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(context.Background(), ctx, helper.TextPrepareDatabase, err, uid)
			output <- ResultQuery{Error: err}
			return
		}

		// commit statement
		tx.Commit()

		output <- ResultQuery{Error: nil}
	}()

	return output
}

func (aqp *AuthQueryPostgres) GetAccountId(ctxReq context.Context, contactId int) <-chan ResultQuery {
	ctx:="AuthRepo-GetAccountId"
	output := make(chan ResultQuery)
	go func ()  {
		defer close(output)

		tx, err := aqp.db.Begin()
		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextDBBegin, err, contactId)
			output <- ResultQuery{Error: err}
			return
		}

		q := `SELECT bbac.id, bbac.status, bbac.is_disabled, bbac.is_delete, bbac.account_id FROM b2b_account_contact bbac LEFT JOIN b2b_account bba ON bba.id = bbac.account_id WHERE bbac.contact_id = $1 AND bba.member_type = 'corporate';`

		stmt, err := tx.Prepare(q)
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(context.Background(), ctx, helper.TextPrepareDatabase, err, q)
			output <- ResultQuery{Error: err}
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(contactId)
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(context.Background(), ctx, helper.TextPrepareDatabase, err, contactId)
			output <- ResultQuery{Error: err}
			return
		}

		// commit statement
		tx.Commit()

		output <- ResultQuery{Error: nil}
	}()

	return output
}
