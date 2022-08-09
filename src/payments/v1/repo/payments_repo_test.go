package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Bhinneka/user-service/src/payments/v1/model"
	sharedRepository "github.com/Bhinneka/user-service/src/shared/repository"
	"github.com/stretchr/testify/assert"
	sqlMock "gopkg.in/DATA-DOG/go-sqlmock.v2"
)

const (
	testEmail                  = "sansa@mailinator.com"
	testChannel                = "b2c"
	testMethod                 = "kredivo"
	expectedQueryInsertPayment = `^INSERT INTO member_payment.*`
	expectedQueryLoad          = `^SELECT .*`
	expectedQueryAddUpdate     = "SELECT email FROM member_payment WHERE id = ?"
	expectedQueryEmail         = `^SELECT "email".*`
)

var (
	paymentField = []string{"Id", "email", "channel", "method", "token", "expiredAt"}

	paymentVal = []model.Payments{
		{
			Email:     testEmail,
			Token:     "token_test",
			Channel:   testChannel,
			Method:    testMethod,
			ExpiredAt: time.Now(),
		},
	}
)

func setupRepoPostgres(t *testing.T) (*PaymentsRepoPostgres, sqlMock.Sqlmock) {
	db, mock, err := sqlMock.New()
	if err != nil {
		t.Fatalf("error new mock %v", err)
	}
	repoMember := &PaymentsRepoPostgres{
		Repository: &sharedRepository.Repository{
			ReadDB:  db,
			WriteDB: db,
		},
	}
	return repoMember, mock
}

func closeRepoPostgres(r *PaymentsRepoPostgres) {
	if r.ReadDB != nil {
		r.ReadDB.Close()
	}
	if r.WriteDB != nil {
		r.WriteDB.Close()
	}
}
func TestNewPaymentRepoPostgres(t *testing.T) {
	t.Run("POSITIVE_NEW_Payment_REPO", func(t *testing.T) {
		repo := NewPaymentsRepoPostgres(&sharedRepository.Repository{})
		assert.NotNil(t, repo)
	})
}
func TestPaymentsRepoPostgres_AddUpdatePayment(t *testing.T) {
	t.Run("Case Failed to Connect DB", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		mock.ExpectBegin().WillReturnError(errors.New("error connnect to DB"))
		res := <-r.AddUpdatePayment(context.Background(), model.Payments{})
		assert.Error(t, res.Error)
	})

	t.Run("Case Error Prepare", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		mock.ExpectBegin()
		mock.ExpectPrepare(expectedQueryAddUpdate).WillReturnError(errors.New("error exec query version"))
		mock.ExpectRollback()

		res := <-r.AddUpdatePayment(context.Background(), model.Payments{})
		assert.Error(t, res.Error)
	})

	t.Run("Case Error Query Row", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		mock.ExpectBegin()
		mock.ExpectPrepare(expectedQueryAddUpdate).ExpectQuery().WillReturnError(errors.New("error query"))
		mock.ExpectRollback()

		defer mock.ExpectClose()

		res := <-r.AddUpdatePayment(context.Background(), model.Payments{})
		assert.Error(t, res.Error)
	})

	t.Run("Case Error Add Update Payment", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)
		rows := sqlMock.NewRows([]string{"email"}).AddRow(paymentVal[0].Email)

		mock.ExpectBegin()
		mock.ExpectPrepare(expectedQueryAddUpdate).ExpectQuery().WillReturnRows(rows)
		mock.ExpectPrepare(expectedQueryInsertPayment).ExpectExec().WillReturnError(errors.New("Error Add Update Payment"))
		mock.ExpectCommit()

		res := <-r.AddUpdatePayment(context.Background(), paymentVal[0])
		assert.Error(t, res.Error)
	})

	t.Run("Case Error Find Payment", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)
		rows := sqlMock.NewRows([]string{"email"}).AddRow(paymentVal[0].Email)

		mock.ExpectBegin()
		mock.ExpectPrepare(expectedQueryAddUpdate).ExpectQuery().WillReturnRows(rows)
		mock.ExpectPrepare(expectedQueryInsertPayment).ExpectExec().WillReturnError(errors.New("Error Add Update Payment"))
		mock.ExpectCommit()

		res := <-r.AddUpdatePayment(context.Background(), model.Payments{})
		assert.Error(t, res.Error)
	})

	t.Run("Case Success", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)
		rows := sqlMock.NewRows([]string{"email"}).AddRow(paymentVal[0].Email)

		mock.ExpectBegin()
		mock.ExpectPrepare(expectedQueryAddUpdate).ExpectQuery().WillReturnRows(rows)
		mock.ExpectPrepare(expectedQueryInsertPayment).ExpectExec().WillReturnResult(sqlMock.NewResult(1, 1))
		mock.ExpectCommit()

		res := <-r.AddUpdatePayment(context.Background(), paymentVal[0])
		assert.NoError(t, res.Error)
	})

}

func TestPaymentsRepoPostgres_FindPaymentByEmailChannelMethod(t *testing.T) {
	t.Run("Test Payment Repo FindPaymentByEmailChannelMethod Success", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		rows := sqlMock.NewRows(paymentField).
			AddRow("TKNXXXXXX", "tes@getnada.com", "b2c", "kredivo", "random token", time.Now())

		query := `SELECT member_payment."id", "email","channel","method","token","expiredAt"`

		mock.ExpectPrepare(query).ExpectQuery().WithArgs("tes@getnada.com", "b2c", "kredivo").WillReturnRows(rows)

		r.FindPaymentByEmailChannelMethod(context.Background(), "tes@getnada.com", "b2c", "kredivo")
	})

	t.Run("Test Payment Repo FindPaymentByEmailChannelMethod Failed Prepared", func(t *testing.T) {
		r, mock := setupRepoPostgres(t)
		defer closeRepoPostgres(r)

		query := `SELECT member_payment."id", "email","channel","method","token","expiredAt"`

		mock.ExpectPrepare(query).ExpectQuery().WillReturnError(errors.New("error prepare"))

		r.FindPaymentByEmailChannelMethod(context.Background(), "tes@getnada.com", "b2c", "kredivo")
	})
}
