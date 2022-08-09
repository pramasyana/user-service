package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/payments/v1/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

type PaymentsRepoPostgres struct {
	*repository.Repository
}

func NewPaymentsRepoPostgres(repo *repository.Repository) *PaymentsRepoPostgres {
	return &PaymentsRepoPostgres{repo}
}

func (pr *PaymentsRepoPostgres) AddUpdatePayment(ctxReq context.Context, payments model.Payments) <-chan ResultRepository {
	ctx := "PaymentsRepo-AddUpdatePayments"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		tx, err := pr.WriteDB.Begin()
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, payments)
			output <- ResultRepository{Error: err}
			return
		}
		readStmt, err := tx.Prepare(`SELECT email FROM member_payment WHERE id = $1`)

		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, payments)
			output <- ResultRepository{Error: err}
			return
		}
		defer readStmt.Close()

		var (
			email string
		)
		err = readStmt.QueryRow(payments.ID).Scan(&email)
		if err != nil && err != sql.ErrNoRows {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, payments)
			tx.Rollback()
			output <- ResultRepository{Error: err}
			return
		}

		if err != sql.ErrNoRows && email != payments.Email {
			tx.Rollback()
			err := fmt.Errorf("there is conflict during save, unable to save model. You can discard the changes or just try again. Persistence Email=%s, Entity Email=%s, ID=%s",
				email, payments.Email, payments.ID)
			helper.SendErrorLog(ctxReq, ctx, "payment", err, payments)
			output <- ResultRepository{Error: err}
			return
		}
		tags[helper.TextEmail] = payments.Email
		err = pr.insertPayment(ctxReq, tx, payments)
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, payments)
			output <- ResultRepository{Error: err}
			return
		}
		// commit statement
		tx.Commit()

		output <- ResultRepository{Error: nil}
	})
	return output
}

func (pr *PaymentsRepoPostgres) insertPayment(ctxReq context.Context, tx *sql.Tx, payments model.Payments) error {
	ctx := "PaymentRepo-insertPayment"

	queryInsert := `INSERT INTO member_payment
					(
						"id", "email", "channel", "method", "token", "expiredAt"
					)
				VALUES
					(
						$1, $2, $3, $4, $5, $6	
					)
				ON CONFLICT("email","channel","method")
				DO UPDATE SET
				"token"=$5, "expiredAt"=$6`

	stmt, err := tx.Prepare(queryInsert)
	tr := tracer.StartTrace(ctxReq, ctx)
	tags := make(map[string]interface{})
	defer func() {
		tr.Finish(tags)
	}()

	tags[helper.TextQuery] = queryInsert
	tags[helper.TextArgs] = payments

	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "insert_payments", err, payments)
		return err
	}

	defer stmt.Close()

	// set null-able variables
	var (
		channel, token, method sql.NullString
	)
	method = helper.ValidateStringToSQLNullString(payments.Method)
	channel = helper.ValidateStringToSQLNullString(payments.Channel)
	token = helper.ValidateStringToSQLNullString(payments.Token)

	_, err = stmt.Exec(
		payments.ID, payments.Email, channel, method, token, payments.ExpiredAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (pr *PaymentsRepoPostgres) FindPaymentByEmailChannelMethod(ctxReq context.Context, email, channel, method string) ResultRepository {
	ctx := "PaymentRepo-FindPaymentByEmailChannelMethod"

	var (
		filter      string
		queryValues []interface{}
	)

	queryValues = append(queryValues, email, channel, method)
	filter = `WHERE "email" = $1 AND "channel" = $2 AND "method" = $3`

	merchant, err := pr.findPayment(ctxReq, ctx, filter, queryValues)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, email)
		return ResultRepository{Error: err}
	}

	return ResultRepository{Result: merchant}
}
func (pr *PaymentsRepoPostgres) findPayment(ctxReq context.Context, ctx, filter string, queryValues []interface{}) (model.Payments, error) {
	var (
		payment model.Payments
	)
	query := fmt.Sprintf(`%s FROM member_payment %s`, pr.getSelect(), filter)

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := make(map[string]interface{})
	defer func() {
		tr.Finish(tags)
	}()

	tags[helper.TextQuery] = query
	stmt, err := pr.ReadDB.Prepare(query)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
		tags[helper.TextResponse] = err
		return payment, err
	}

	defer stmt.Close()

	err = stmt.QueryRow(queryValues...).Scan(
		&payment.ID, &payment.Email, &payment.Channel, &payment.Method, &payment.Token, &payment.ExpiredAt,
	)

	if _, err := json.Marshal(&payment); err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, payment)
		return payment, err
	}

	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
		tags[helper.TextResponse] = err
		return payment, err
	}

	tags[helper.TextResponse] = payment
	return payment, nil
}

func (pr *PaymentsRepoPostgres) getSelect() string {
	return `
	SELECT member_payment."id", "email","channel","method","token","expiredAt"`

}
