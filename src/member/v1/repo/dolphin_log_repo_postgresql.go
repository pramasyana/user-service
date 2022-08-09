package repo

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/Bhinneka/user-service/src/member/v1/model"
)

// DolphinLogRepoPostgres data structure
type DolphinLogRepoPostgres struct {
	db *sql.DB
}

// NewDolphinLogRepoPostgres function for initializing dolphin log repo
func NewDolphinLogRepoPostgres(db *sql.DB) *DolphinLogRepoPostgres {
	return &DolphinLogRepoPostgres{db: db}
}

// Save function function for save dolphin log
func (r *DolphinLogRepoPostgres) Save(ctxReq context.Context, data *model.DolphinLog) error {
	query := `INSERT INTO "log_dolphin" ("user_id", "event_type", "log_data", "created") VALUES($1, $2, $3::jsonb, $4)`

	stmt, err := r.db.Prepare(query)

	if err != nil {
		return err
	}

	logDataJSON, _ := json.Marshal(data.LogData)

	_, err = stmt.Exec(data.UserID, data.EventType, logDataJSON, data.Created)

	if err != nil {
		return err
	}

	return nil
}

// Load function by id
func (r *DolphinLogRepoPostgres) Load(id int) ResultRepository {
	return ResultRepository{}
}
