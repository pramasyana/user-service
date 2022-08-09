package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
)

const (
	txtScanDB   = "scan_database"
	txtClientID = "client id is not found"
)

// ClientAppRepoPostgres data structure
type ClientAppRepoPostgres struct {
	db *sql.DB
}

// NewClientAppRepoPostgres function for initializing auth repo
func NewClientAppRepoPostgres(db *sql.DB) *ClientAppRepoPostgres {
	return &ClientAppRepoPostgres{db: db}
}

// Save function
func (repo *ClientAppRepoPostgres) Save(m *model.ClientApp) <-chan ResultRepository {
	ctx := "ClientAppRepo-Save"
	ctxReq := context.Background()

	output := make(chan ResultRepository)
	go func() {
		defer close(output)

		var version int
		tx, err := repo.db.Begin()
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextDBBegin, err, m)
			output <- ResultRepository{Error: err}
			return
		}

		// check version client app
		version, err = repo.CheckClientApp(tx, m)
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, m)
			output <- ResultRepository{Error: err}
			return
		}

		// MVCC https://en.wikipedia.org/wiki/Multiversion_concurrency_control
		m.Version++

		if version > m.Version {
			tx.Rollback()
			err := fmt.Errorf("there is conflict during save, unable to save model. You can discard the changes or just try again. Persistence Version=%d, Entity Version=%d, ID=%s",
				version, m.Version, m.ID)
			helper.SendErrorLog(ctxReq, ctx, helper.ScopeSaveMember, err, m)
			output <- ResultRepository{Error: err}
			return
		}

		q := `INSERT INTO client_app
				(
					"clientId", "clientSecret", name, status,
					created, "lastModified", version
				)
			VALUES
				(
					$1, $2, $3, $4, $5, $6, $7
				)
			ON CONFLICT(id)
			DO UPDATE SET
				"clientId" = $1, "clientSecret" = $2, name = $3, status = $4,
				 "lastModified" = $6, version = $7`

		stmt, err := tx.Prepare(q)
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(ctxReq, ctx, "insert_client_app", err, q)
			output <- ResultRepository{Error: err}
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(
			m.ClientID, m.ClientSecret, m.Name, m.Status.String(),
			time.Now(), time.Now(), m.Version,
		)
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, m)
			output <- ResultRepository{Error: err}
			return
		}

		// commit statement
		tx.Commit()

		output <- ResultRepository{Error: nil}
	}()

	return output
}

// check version client app
func (repo *ClientAppRepoPostgres) CheckClientApp(tx *sql.Tx, m *model.ClientApp) (int, error) {
	var version int
	if m.ID != "" {
		readStmt, err := tx.Prepare(`SELECT "version" FROM client_app WHERE id = $1`)
		if err != nil {
			return version, err
		}
		defer readStmt.Close()

		err = readStmt.QueryRow(m.ID).Scan(&version)
		if err != nil && err != sql.ErrNoRows {
			return version, err
		}
	}

	return version, nil
}

// Load function for loading basic auth based on username
func (repo *ClientAppRepoPostgres) Load(id int) <-chan ResultRepository {
	ctx := "ClientAppRepo-LoadByID"

	output := make(chan ResultRepository)

	go func() {
		defer close(output)
		resp := repo.findBy(ctx, `id`, id)
		output <- resp
	}()

	return output
}

// FindByClientID function for loading client app based on client id
func (repo *ClientAppRepoPostgres) FindByClientID(clientID string) <-chan ResultRepository {
	ctx := "ClientAppRepo-LoadByClientID"

	output := make(chan ResultRepository)

	go func() {
		defer close(output)
		resp := repo.findBy(ctx, `"clientId"`, clientID)
		output <- resp
	}()

	return output
}

func (repo *ClientAppRepoPostgres) findBy(ctx, field string, value interface{}) ResultRepository {
	var (
		app                model.ClientApp
		clientIDString     sql.NullString
		clientSecretString sql.NullString
		nameString         sql.NullString
		status             string
	)
	query := `
		SELECT id, "clientId", "clientSecret", name, status, created, 
		"lastModified", "version" FROM client_app`

	query += ` WHERE ` + field + ` = $1`
	stmt, err := repo.db.Prepare(query)

	if err != nil {
		helper.SendErrorLog(context.Background(), ctx, helper.TextPrepareDatabase, err, query)
		return ResultRepository{Error: err}
	}
	defer stmt.Close()

	err = stmt.QueryRow(value).Scan(&app.ID,
		&clientIDString, &clientSecretString, &nameString,
		&status, &app.Created, &app.LastModified, &app.Version)

	if err == sql.ErrNoRows {
		return ResultRepository{nil, errors.New(txtClientID)}
	} else if err != nil {
		helper.SendErrorLog(context.Background(), ctx, txtScanDB, err, query)
		return ResultRepository{Error: err}
	}

	if clientIDString.Valid {
		app.ClientID = clientIDString.String
	}

	if clientSecretString.Valid {
		app.ClientSecret = clientSecretString.String
	}

	if nameString.Valid {
		app.Name = nameString.String
	}

	app.Status = model.StringToStatus(status)
	return ResultRepository{Result: app}
}
