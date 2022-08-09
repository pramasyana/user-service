package repo

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"github.com/Bhinneka/user-service/src/member/v1/model"
	sharedRepository "github.com/Bhinneka/user-service/src/shared/repository"
	"github.com/stretchr/testify/assert"
	sqlMock "gopkg.in/DATA-DOG/go-sqlmock.v2"
)

const (
	userID   = "USR1827715"
	authType = "LDAP"
)

var logData = model.MemberAdditionalInfo{
	ID:       1,
	MemberID: "123",
	AuthType: authType,
	Data:     "",
}

func setupRepoAdditional(t *testing.T) (*MemberAdditionalInfoRepoPostgres, sqlMock.Sqlmock) {
	db, mock, err := sqlMock.New()
	if err != nil {
		t.Fatal(err)
	}
	repoAdditional := &MemberAdditionalInfoRepoPostgres{
		&sharedRepository.Repository{
			ReadDB:  db,
			WriteDB: db,
		},
	}
	return repoAdditional, mock
}

func TestSaveAdditionalInfo(t *testing.T) {
	expectedInsertAdditionalQuery := `^INSERT INTO "member_additional_info" .*`

	t.Run("POSITIVE_SAVE_ADDITIONAL_INFO", func(t *testing.T) {
		r, mock := setupRepoAdditional(t)
		defer func() {
			r.WriteDB.Close()
			r.ReadDB.Close()
		}()
		mock.ExpectPrepare(expectedInsertAdditionalQuery).ExpectExec().WillReturnResult(sqlMock.NewResult(1, 1))
		saveResult := <-r.Save(context.Background(), &logData)
		assert.NoError(t, saveResult.Error)
	})

	t.Run("NEGATIVE_SAVE_ADDITIONAL_INFO_EXEC_QUERY", func(t *testing.T) {
		r, mock := setupRepoAdditional(t)
		defer func() {
			r.WriteDB.Close()
			r.ReadDB.Close()
		}()
		mock.ExpectPrepare(expectedInsertAdditionalQuery).ExpectExec().WillReturnError(errors.New("error insert additional info exec query"))
		saveResult := <-r.Save(context.Background(), &logData)
		assert.Error(t, saveResult.Error)
	})

	t.Run("NEGATIVE_SAVE_ADDITIONAL_INFO", func(t *testing.T) {
		r, _ := setupRepoAdditional(t)
		r.WriteDB.Close()
		r.ReadDB.Close()
		saveResult := <-r.Save(context.Background(), &logData)
		assert.Error(t, saveResult.Error)
	})

}

func TestUpdateAdditionalInfo(t *testing.T) {
	expectedQueryUpdate := `^UPDATE "member_additional_info" .*`

	t.Run("POSITIVE_UPDATE_ADDITIONAL_INFO", func(t *testing.T) {
		r, mock := setupRepoAdditional(t)
		defer func() {
			r.WriteDB.Close()
			r.ReadDB.Close()
		}()
		mock.ExpectPrepare(expectedQueryUpdate).ExpectExec().WillReturnResult(sqlMock.NewResult(1, 1))
		saveResult := <-r.Update(context.Background(), &logData)
		assert.NoError(t, saveResult.Error)
	})

	t.Run("NEGATIVE_UPDATE_ADDITIONAL_INFO_PREPARE_STATEMENT", func(t *testing.T) {
		r, mock := setupRepoAdditional(t)
		defer func() {
			r.WriteDB.Close()
			r.ReadDB.Close()
		}()
		mock.ExpectPrepare(expectedQueryUpdate).WillReturnError(errors.New("error update additional info prepare statement"))
		saveResult := <-r.Update(context.Background(), &logData)
		assert.Error(t, saveResult.Error)
	})

	t.Run("NEGATIVE_UPDATE_ADDITIONAL_INFO_EXEC_QUERY", func(t *testing.T) {
		r, mock := setupRepoAdditional(t)
		defer func() {
			r.WriteDB.Close()
			r.ReadDB.Close()
		}()
		mock.ExpectPrepare(expectedQueryUpdate).ExpectExec().WillReturnError(errors.New("error update additional info exec query"))
		saveResult := <-r.Update(context.Background(), &logData)
		assert.Error(t, saveResult.Error)
	})

	t.Run("NEGATIVE_UPDATE_ADDITIONAL_INFO", func(t *testing.T) {
		r, _ := setupRepoAdditional(t)
		r.WriteDB.Close()
		r.ReadDB.Close()
		saveResult := <-r.Update(context.Background(), &logData)
		assert.Error(t, saveResult.Error)
	})
}

func TestLoadAdditionalInfo(t *testing.T) {

	memberField := []string{"id", "memberId", "authType",
		"data", "created", "lastModified",
	}
	data := "[{\"key\": \"objectClass\"}]"
	expectedQuery := `SELECT * FROM member_additional_info`

	t.Run("POSITIVE_LOAD_ADDITIONAL_INFO", func(t *testing.T) {
		r, mock := setupRepoAdditional(t)
		defer func() {
			r.WriteDB.Close()
			r.ReadDB.Close()
		}()

		query := regexp.QuoteMeta(expectedQuery)

		rows := sqlMock.NewRows(memberField).
			AddRow("1", "USR123", authType, data, time.Now(), time.Now())

		mock.ExpectPrepare(query).ExpectQuery().WithArgs(userID, authType).WillReturnRows(rows)

		memberResult := <-r.Load(context.Background(), userID, authType)

		assert.NoError(t, memberResult.Error)
		assert.IsType(t, model.MemberAdditionalInfo{}, memberResult.Result)
	})

	t.Run("NEGATIVE_LOAD_ADDITIONAL_INFO_INVALID_AUTH", func(t *testing.T) {
		r, mock := setupRepoAdditional(t)
		defer func() {
			r.WriteDB.Close()
			r.ReadDB.Close()
		}()

		query := regexp.QuoteMeta(expectedQuery)

		rows := sqlMock.NewRows(memberField).
			AddRow("2", "USR456", authType, data, time.Now(), time.Now())

		mock.ExpectPrepare(query).ExpectQuery().WithArgs(userID, authType).WillReturnRows(rows)

		memberResult := <-r.Load(context.Background(), "USR182771s", "a")

		assert.Error(t, memberResult.Error)
	})

	t.Run("NEGATIVE_LOAD_ADDITIONAL_INFO_PREPARE_STATEMENT", func(t *testing.T) {
		r, mock := setupRepoAdditional(t)
		defer func() {
			r.WriteDB.Close()
			r.ReadDB.Close()
		}()

		query := regexp.QuoteMeta(expectedQuery)

		mock.ExpectPrepare(query).WillReturnError(errors.New("error load additional info prepare query statement"))

		memberResult := <-r.Load(context.Background(), userID, authType)

		assert.Error(t, memberResult.Error)
	})
}

func TestNewMemberAdditionalInfoRepoPostgres(t *testing.T) {
	t.Run("POSITIVE_NEW_ADDITIONAL_INFO_REPO", func(t *testing.T) {
		repo := NewMemberAdditionalInfoRepoPostgres(&sharedRepository.Repository{})
		assert.NotNil(t, repo)
	})
}
