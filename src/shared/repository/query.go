package repository

import (
	"fmt"
	"strings"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/shared/model"
	"github.com/pkg/errors"
)

// DeleteByID delete row by ID
func DeleteByID(repo *Repository, rowID int, tableName string) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1;`, tableName)

	stmt, err := repo.WriteDB.Prepare(query)
	if err != nil {
		return errors.Wrap(err, helper.TextPrepareDatabase)
	}
	if _, err = stmt.Exec(rowID); err != nil {
		return errors.Wrap(err, helper.TextExecQuery)
	}
	return nil
}

// Exec function to execute specific query with parameters
func Exec(repo *Repository, query string, param ...interface{}) (err error) {
	stmt, err := repo.WriteDB.Prepare(query)

	if err != nil {
		return errors.Wrap(err, helper.TextPrepareDatabase)
	}

	if _, err = stmt.Exec(param...); err != nil {
		return errors.Wrap(err, helper.TextExecQuery)
	}

	return nil
}

// Delete delete row by specific criteria
func Delete(repo *Repository, param map[model.DBOperator]interface{}, tableName string) error {
	clause, variadic, _ := getClause(param)
	query := fmt.Sprintf(`DELETE FROM %s WHERE %s;`, tableName, clause)

	stmt, err := repo.WriteDB.Prepare(query)
	if err != nil {
		return errors.Wrap(err, helper.TextPrepareDatabase)
	}
	if _, err = stmt.Exec(variadic); err != nil {
		return errors.Wrap(err, helper.TextExecQuery)
	}
	return nil
}

func getClause(param map[model.DBOperator]interface{}) (clause string, variadic []interface{}, bindVariadic map[string]interface{}) {
	if len(param) == 0 {
		return "", nil, nil
	}
	where := []string{}
	bindVar := make(map[string]interface{})
	for a, b := range param {
		key := `$` + a.Field
		where = append(where, a.Field+a.Operator+`$`+a.Field)
		bindVar[key] = b
		variadic = append(variadic, b)
	}

	return strings.Join(where, " AND "), variadic, bindVar
}
