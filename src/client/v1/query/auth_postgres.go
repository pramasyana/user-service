package query

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	clientModel "github.com/Bhinneka/user-service/src/client/v1/model"
)

func NewClientAppQuery(repo *sql.DB) *ClientAppQueryPostgres {
	return &ClientAppQueryPostgres{
		db: repo,
	}
}

type ClientAppQueryPostgres struct {
	db *sql.DB
}

func (q *ClientAppQueryPostgres) Validate(ctxReq context.Context, clientID, clientSecret string) (bool, error) {
	ctx := "ClientAppQuery-Validate"
	tr := tracer.StartTrace(ctxReq, ctx)
	tags := map[string]interface{}{
		"clientId":     clientID,
		"clientSecret": clientSecret,
	}

	defer tr.Finish(tags)

	var (
		client                                         clientModel.ClientApp
		clientIDString, clientSecretString, nameString sql.NullString
		clientStatus                                   string
		valid                                          bool
	)
	query := `SELECT id, "clientId", "clientSecret", status FROM client_app WHERE "clientId" = $1`

	stmt, err := q.db.Prepare(query)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "validate_client", err, tags)
		return valid, errors.New("oops something bad on us")
	}

	defer stmt.Close()

	err = stmt.QueryRow(clientID).Scan(&client.ID, &clientIDString, &clientSecretString, &clientStatus)
	if err == sql.ErrNoRows {
		return valid, errors.New("invalid credential")
	} else if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "query_row", err, nil)
		return valid, errors.New("oops")
	}

	if clientIDString.Valid {
		client.ClientID = clientIDString.String
	}

	if clientSecretString.Valid {
		client.ClientSecret = clientSecretString.String
	}

	if nameString.Valid {
		client.Name = nameString.String
	}

	client.Status = clientModel.StringToStatus(clientStatus)

	// validate client app secret
	check := client.Authenticate(clientSecret)
	if !check {
		return valid, errors.New("client credential does not match")
	}
	if !client.IsActive() {
		return valid, errors.New("client status is not active")
	}
	return true, nil
}
