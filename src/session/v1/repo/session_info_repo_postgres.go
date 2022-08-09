package repo

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/session/v1/model"
)

// SessionInfoRepoPostgres data structure
type SessionInfoRepoPostgres struct {
	db *sql.DB
}

// NewSessionInfoRepoPostgres function for initializing session info repo
func NewSessionInfoRepoPostgres(db *sql.DB) *SessionInfoRepoPostgres {
	return &SessionInfoRepoPostgres{db: db}
}

// SaveSessionInfo for save session info user login
func (qp *SessionInfoRepoPostgres) SaveSessionInfo(params *model.SessionInfoRequest) <-chan ResultRepository {
	ctx := "SessionInfoQuery-SaveSessionInfo"
	output := make(chan ResultRepository)
	go func() {
		defer close(output)

		tx, err := qp.db.Begin()
		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, helper.TextDBBegin, err, params)
			output <- ResultRepository{Error: err}
			return
		}

		t := time.Now()

		if len(params.IP) > 20 {
			splitted := helper.SplitStringByN(params.IP, 20)
			params.IP = splitted[0]
		}

		q := `INSERT INTO session_info("userId","userName",ip,"userAgent","deviceId","clientType","grantType",jti,"createdAt") VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9)`

		stmt, err := tx.Prepare(q)
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(context.Background(), ctx, helper.TextPrepareDatabase, err, q)
			output <- ResultRepository{Error: err}
			return
		}
		defer stmt.Close()

		if _, err = stmt.Exec(params.UserID, params.Email, params.IP, params.UserAgent, params.DeviceID, params.DeviceLogin, params.GrantType, params.JTI, t.Format(time.RFC3339)); err != nil {
			tx.Rollback()
			helper.SendErrorLog(context.Background(), ctx, helper.TextQueryDatabase, err, q)
			output <- ResultRepository{Error: errors.New(`unable to logged in`)}
			return
		}

		tx.Commit()

		output <- ResultRepository{Error: nil}
	}()

	return output
}
