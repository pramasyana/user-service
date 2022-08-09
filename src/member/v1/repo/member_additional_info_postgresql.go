package repo

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

// MemberAdditionalInfoRepoPostgres data structure
type MemberAdditionalInfoRepoPostgres struct {
	*repository.Repository
}

// NewMemberAdditionalInfoRepoPostgres function for initializing additional info repo
func NewMemberAdditionalInfoRepoPostgres(repo *repository.Repository) *MemberAdditionalInfoRepoPostgres {
	return &MemberAdditionalInfoRepoPostgres{repo}
}

// Save function function for save additional info
func (r *MemberAdditionalInfoRepoPostgres) Save(ctxReq context.Context, data *model.MemberAdditionalInfo) <-chan ResultRepository {
	ctx := "MemberAdditionalInfoRepo-Save"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)

		query := `INSERT INTO "member_additional_info" 
				("memberId", "authType", "data", "created", "lastModified") 
				VALUES($1, $2, $3::jsonb, $4, $5)`

		tags[helper.TextQuery] = query
		stmt, err := r.WriteDB.Prepare(query)

		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		dataJSON, _ := json.Marshal(data.Data)

		_, err = stmt.Exec(data.MemberID, data.AuthType, dataJSON, time.Now(), time.Now())

		if err != nil {
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: nil}
	})
	return output
}

// Update function function for update additional info
func (r *MemberAdditionalInfoRepoPostgres) Update(ctxReq context.Context, data *model.MemberAdditionalInfo) <-chan ResultRepository {
	ctx := "MemberAdditionalInfoRepo-Update"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		query := `UPDATE "member_additional_info"  
				SET "memberId"=$2, "authType"=$3, "data"=$4, "lastModified"=$5 WHERE "id"=$1`

		tags[helper.TextQuery] = query
		stmt, err := r.WriteDB.Prepare(query)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, data)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		dataJSON, _ := json.Marshal(data.Data)

		_, err = stmt.Exec(data.ID, data.MemberID, data.AuthType, dataJSON, time.Now())

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, data)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: nil}
	})
	return output
}

// Load function for loading member additional info based on member id
func (r *MemberAdditionalInfoRepoPostgres) Load(ctxReq context.Context, uid string, authType string) <-chan ResultRepository {
	ctx := "MemberAdditionalInfoRepo-Load"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		q := `SELECT * FROM member_additional_info WHERE "memberId" = $1 AND "authType"=$2`

		tags[helper.TextQuery] = q

		stmt, err := r.ReadDB.Prepare(q)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, uid)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}
		defer stmt.Close()

		var member model.MemberAdditionalInfo

		err = stmt.QueryRow(uid, authType).Scan(
			&member.ID,
			&member.MemberID,
			&member.AuthType,
			&member.Data,
			&member.Created,
			&member.LastModified,
		)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, uid)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		tags["args"] = member
		output <- ResultRepository{Result: member}
	})

	return output
}
