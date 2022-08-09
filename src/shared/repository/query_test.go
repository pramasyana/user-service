package repository

import (
	"testing"

	"github.com/Bhinneka/user-service/src/shared/model"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestQueryParameter(t *testing.T) {
	param := make(map[model.DBOperator]interface{})
	param[model.DBOperator{Field: "id", Operator: "="}] = 1
	param[model.DBOperator{Field: "name", Operator: "="}] = "someName"
	param[model.DBOperator{Field: "isDeleted", Operator: "="}] = false

	query, variadic, bindVar := getClause(param)
	assert.Contains(t, query, "isDeleted=$isDeleted")
	assert.Contains(t, query, "id=$id")
	assert.Contains(t, query, "name=$name")
	assert.Equal(t, 1, bindVar["$id"])
	assert.Equal(t, "someName", bindVar["$name"])
	assert.Equal(t, false, bindVar["$isDeleted"])
	assert.NotNil(t, variadic)
	assert.Equal(t, 3, len(variadic))
}

func TestDeleteByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("DELETE FROM b2b_contact WHERE id=$1").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repository := &Repository{
		WriteDB: db,
	}

	// now we execute our method
	DeleteByID(repository, 1, "b2b_contact")
}

func TestExec(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("SELECT * FROM b2b_contact").WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	repository := &Repository{
		WriteDB: db,
	}

	// now we execute our method
	Exec(repository, "SELECT * FROM b2b_contact", nil)
}
