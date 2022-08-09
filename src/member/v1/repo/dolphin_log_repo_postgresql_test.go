package repo

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/stretchr/testify/assert"
	sqlMock "gopkg.in/DATA-DOG/go-sqlmock.v2"
)

func setupRepoDolphin(t *testing.T) (*DolphinLogRepoPostgres, sqlMock.Sqlmock) {
	db, mock, err := sqlMock.New()
	if err != nil {
		t.Fatal(err)
	}
	repoDolphin := &DolphinLogRepoPostgres{
		db: db,
	}
	return repoDolphin, mock
}

func TestNewDolphinLogRepoPostgres(t *testing.T) {
	t.Run("POSITIVE_NEW_DOLPHIN_LOG_REPO", func(t *testing.T) {
		db, _, _ := sqlMock.New()
		defer db.Close()

		r := NewDolphinLogRepoPostgres(db)
		assert.NotNil(t, r)
	})
}

func TestSaveDolphinLog(t *testing.T) {
	const expectedInsertLogQuery string = `^INSERT INTO "log_dolphin" .*`
	var logData = model.DolphinLog{
		ID:        1,
		UserID:    "123",
		EventType: "INSERT",
		Created:   time.Now(),
	}
	ctx := context.Background()

	t.Run("POSITIVE_SAVE_DOPLHIN_LOG", func(t *testing.T) {
		r, mock := setupRepoDolphin(t)
		defer r.db.Close()

		mock.ExpectPrepare(expectedInsertLogQuery).ExpectExec().WillReturnResult(sqlMock.NewResult(1, 1))

		err := r.Save(ctx, &logData)
		assert.NoError(t, err)
	})

	t.Run("NEGATIVE_SAVE_DOLPHIN_LOG_PREPARE_STATEMENT", func(t *testing.T) {
		r, mock := setupRepoDolphin(t)
		defer r.db.Close()

		mock.ExpectPrepare(expectedInsertLogQuery).WillReturnError(errors.New("error save dolphin log prepare query statement"))

		err := r.Save(ctx, &logData)
		assert.Error(t, err)
	})

	t.Run("NEGATIVE_SAVE_DOLPHIN_LOG_EXEC", func(t *testing.T) {
		r, mock := setupRepoDolphin(t)
		defer r.db.Close()

		mock.ExpectPrepare(expectedInsertLogQuery).ExpectExec().WillReturnError(errors.New("error save dolphin log exec query"))

		err := r.Save(ctx, &logData)
		assert.Error(t, err)
	})
}

func TestDolphinLogRepoPostgresLoad(t *testing.T) {
	t.Run("POSITIVE_DOLPHIN_LOG_LOAD", func(t *testing.T) {
		r, _ := setupRepoDolphin(t)
		defer r.db.Close()
		result := r.Load(0)
		assert.NotNil(t, result)
	})
}
