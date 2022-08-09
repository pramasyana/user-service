package query

import (
	"context"
	"database/sql"

	"github.com/Bhinneka/user-service/helper"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
)

// AccountContactQueryPostgres data structure
type AccountContactQueryPostgres struct {
	db *sql.DB
}

// NewAccountContactQueryPostgres function for initializing Accountcontact query
func NewAccountContactQueryPostgres(db *sql.DB) *AccountContactQueryPostgres {
	return &AccountContactQueryPostgres{db: db}
}

// FindByAccountContactID function for getting detail Accountcontact by contact id
func (mq *AccountContactQueryPostgres) FindByAccountContactID(id int) <-chan ResultQuery {
	ctx := "AccountContactQuery-FindByAccountContactID"

	output := make(chan ResultQuery)
	go func() {
		defer close(output)

		querySelect := `SELECT id, status, is_disabled, is_delete, account_id
			FROM b2b_account_contact WHERE contact_id = $1`

		stmt, err := mq.db.Prepare(querySelect)
		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, querySelect)
			output <- ResultQuery{Error: err}
			return
		}
		defer stmt.Close()

		var accContact sharedModel.B2BAccountContact

		if err := stmt.QueryRow(id).Scan(&accContact.ID, &accContact.Status, &accContact.IsDisabled, &accContact.IsDelete, &accContact.AccountID); err != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, querySelect)
			output <- ResultQuery{Error: err}
			return
		}

		output <- ResultQuery{Result: accContact}
	}()

	return output
}

// FindAccountMicrositeByContactID function for getting detail Accountcontact by contact id
func (mq *AccountContactQueryPostgres) FindAccountMicrositeByContactID(id int) <-chan ResultQuery {
	ctx := "AccountContactQuery-FindAccountMicrositeByContactID"

	output := make(chan ResultQuery)
	go func() {
		defer close(output)

		q := `SELECT b2b_contact.id, b2b_contact.first_name, b2b_contact.last_name, b2b_contact.phone_number, ac.account_id 
		FROM b2b_account_contact ac
		LEFT JOIN b2b_account ON b2b_account.id = ac.account_id
		LEFT JOIN b2b_contact ON b2b_contact.id = ac.contact_id 
		WHERE b2b_account.is_microsite = true and contact_id = $1`

		stmt, err := mq.db.Prepare(q)
		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextPrepareDatabase, err, q)
			output <- ResultQuery{Error: err}
			return
		}
		defer stmt.Close()

		// initialize needed variables
		var (
			accContact                       sharedModel.B2BContactData
			LastName, PhoneNumber, AccountID sql.NullString
		)

		err = stmt.QueryRow(id).Scan(
			&accContact.ID, &accContact.FirstName, &LastName, &PhoneNumber, &AccountID,
		)

		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextQueryDatabase, err, id)
			output <- ResultQuery{Error: err}
			return
		}

		accContact.LastName = helper.ValidateSQLNullString(LastName)
		accContact.PhoneNumber = helper.ValidateSQLNullString(PhoneNumber)
		accContact.AccountID = helper.ValidateSQLNullString(AccountID)

		output <- ResultQuery{Result: accContact}

	}()

	return output
}

// FindByAccountMicrositeContactID function for getting detail Accountcontact by contact id
func (mq *AccountContactQueryPostgres) FindByAccountMicrositeContactID(id int) <-chan ResultQuery {
	ctx := "AccountContactQuery-FindByAccountMicrositeContactID"

	output := make(chan ResultQuery)
	go func() {
		defer close(output)

		q := `SELECT ac.id, ac.status, ac.is_disabled, ac.is_delete 
			FROM b2b_account_contact ac
			LEFT JOIN b2b_account ON b2b_account.id = ac.account_id
			LEFT JOIN b2b_contact ON b2b_contact.id = ac.contact_id 
			WHERE b2b_account.is_microsite = true and contact_id = $1`

		stmt, err := mq.db.Prepare(q)
		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextPrepareDatabase, err, id)
			output <- ResultQuery{Error: err}
			return
		}
		defer stmt.Close()
		// initialize needed variables
		var accContact sharedModel.B2BAccountContact

		err = stmt.QueryRow(id).Scan(
			&accContact.ID, &accContact.Status, &accContact.IsDisabled, &accContact.IsDelete,
		)

		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextQueryDatabase, err, id)
			output <- ResultQuery{Error: err}
			return
		}

		output <- ResultQuery{Result: accContact}

	}()

	return output
}
