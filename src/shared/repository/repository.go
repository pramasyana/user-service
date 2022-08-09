package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/Bhinneka/user-service/helper"
)

// Repository parent for all repository
type Repository struct {
	ReadDB, WriteDB *sql.DB
	Tx              *sql.Tx
	mut             sync.Mutex
}

// NewRepository create new Repository domain
func NewRepository(readDB, writeDB *sql.DB) *Repository {
	return &Repository{ReadDB: readDB, WriteDB: writeDB}
}

// StartTransaction for starting transaction on each child repository
func (r *Repository) StartTransaction() {
	defer func() {
		if r := recover(); r != nil {
			helper.SendErrorLog(context.Background(), "repo-transaction", helper.TextExecQuery, errors.New("error transaction query"), r)
		}
	}()

	tx, err := r.WriteDB.Begin()
	if err != nil {
		helper.SendErrorLog(context.Background(), "transaction_repo", helper.TextExecQuery, err, r)
	}
	r.mut.Lock()
	r.Tx = tx
}

// Rollback saved data if error happened
func (r *Repository) Rollback() {
	if r.Tx != nil {
		if err := r.Tx.Rollback(); err != nil {
			helper.SendErrorLog(context.Background(), "Repo-Rollback", "rollback_repo", err, nil)
		}
	}
	r.Tx = nil
	r.mut.Unlock()
}

// Commit for final save data
func (r *Repository) Commit() {
	if r.Tx != nil {
		if err := r.Tx.Commit(); err != nil {
			helper.SendErrorLog(context.Background(), "Repo-Commit", "commit_repo", err, nil)
		}
	}
	r.Tx = nil
	r.mut.Unlock()
}

// WithTransaction run db transaction in repository with clean closure txFunc
func (r *Repository) WithTransaction(txFunc func() error) (err error) {
	r.StartTransaction()
	defer func() {
		if rec := recover(); rec != nil {
			r.Rollback()
			err = fmt.Errorf("panic: %v", rec)
		} else if err != nil {
			r.Rollback()
		} else {
			r.Commit()
		}
	}()

	err = txFunc()
	return
}

// Prepare wrapper
func (r *Repository) Prepare(ctxReq context.Context, q string) (*sql.Stmt, error) {
	ctx := "Prepare-Repository"
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("error prepare query %s", q)
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, r)
		}
	}()
	if strings.Contains(q, strings.ToLower("select")) {
		return r.ReadDB.Prepare(q)
	}
	return r.WriteDB.Prepare(q)
}

// Query wrapper
func (r *Repository) Query(query string, args ...interface{}) (*sql.Rows, error) {
	ctx := "Query-Repository"
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("error execute query %s", query)
			helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, r)
		}
	}()
	return r.ReadDB.Query(query, args...)
}
