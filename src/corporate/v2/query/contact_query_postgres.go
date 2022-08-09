package query

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/corporate/v2/model"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
	"github.com/lib/pq"
)

// ContactQueryPostgres data structure
type ContactQueryPostgres struct {
	db *sql.DB
}

// NewContactQueryPostgres function for initializing contact query
func NewContactQueryPostgres(db *sql.DB) *ContactQueryPostgres {
	return &ContactQueryPostgres{db: db}
}

// FindByEmail function for getting detail contact by email
func (mq *ContactQueryPostgres) FindByEmail(ctxReq context.Context, email string) <-chan ResultQuery {
	ctx := "ContactQuery-Corporate-FindByEmail"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		var queryValues []interface{}

		email = strings.ToLower(email)
		queryValues = append(queryValues, email)
		filter := `WHERE "is_disabled" = false AND LOWER("email") = $1`
		tags[helper.TextParameter] = filter
		member := <-mq.FindContact(ctxReq, filter, queryValues)
		if member.Error != nil {
			output <- ResultQuery{Error: member.Error}
			return
		}

		contact := member.Result.(sharedModel.B2BContactData)
		output <- ResultQuery{Result: contact}

	})

	return output
}

// FindByID function for getting detail contact by ID
func (mq *ContactQueryPostgres) FindByID(ctxReq context.Context, uid string) <-chan ResultQuery {
	ctx := "ContactQuery-Corporate-FindByID"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		var queryValues []interface{}

		queryValues = append(queryValues, uid)
		filter := `WHERE "is_disabled" = false AND "id" = $1`
		tags[helper.TextParameter] = filter
		member := <-mq.FindContact(ctxReq, filter, queryValues)
		if member.Error != nil {
			err := member.Error
			if member.Error != sql.ErrNoRows {
				err = fmt.Errorf(helper.ErrorDataNotFound, "user")
			}

			output <- ResultQuery{Error: err}
			return

		}

		contact := member.Result.(sharedModel.B2BContactData)
		output <- ResultQuery{Result: contact}
	})

	return output
}

// FindContact function for getting detail contact by ID
func (mq *ContactQueryPostgres) FindContact(ctxReq context.Context, filter string, queryValues []interface{}) <-chan ResultQuery {
	ctx := "ContactQuery-Corporate-FindContact"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		querySelect := fmt.Sprintf(`SELECT 
			id, "reference_id", "first_name", "last_name", "salutation", 
			"job_title", "email", "nav_contact_id", "is_primary", "birth_date", 
			"note", "created_at", "modified_at","created_by", "modified_by", 
			"account_id",  "is_disabled", "password","token", "avatar", 
			"status", "is_new", "phone_number"
			FROM b2b_contact %s`, filter)

		tags[helper.TextQuery] = querySelect
		stmt, err := mq.db.Prepare(querySelect)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, querySelect)
			tags[helper.TextResponse] = err
			output <- ResultQuery{Error: err}
			return
		}

		defer stmt.Close()

		// initialize needed variables
		var (
			contact                                                                            sharedModel.B2BContactData
			isNew                                                                              sql.NullBool
			birthDate                                                                          pq.NullTime
			referenceID, navContactID, accountID, phoneNumber, avatar, status, password, token sql.NullString
			createdBy, modifiedBy                                                              sql.NullInt64
		)

		err = stmt.QueryRow(queryValues...).Scan(
			&contact.ID, &referenceID, &contact.FirstName, &contact.LastName, &contact.Salutation,
			&contact.JobTitle, &contact.Email, &navContactID, &contact.IsPrimary, &birthDate,
			&contact.Note, &contact.CreatedAt, &contact.ModifiedAt,
			&createdBy, &modifiedBy,
			&accountID, &contact.IsDisabled, &password, &token, &avatar,
			&status, &isNew, &phoneNumber,
		)

		if isNew.Valid {
			contact.IsNew = isNew.Bool
		}

		if birthDate.Valid {
			contact.BirthDate = birthDate.Time
		}

		contact.ReferenceID = helper.ValidateSQLNullString(referenceID)
		contact.NavContactID = helper.ValidateSQLNullString(navContactID)
		contact.AccountID = helper.ValidateSQLNullString(accountID)
		contact.PhoneNumber = helper.ValidateSQLNullString(phoneNumber)
		contact.Password = helper.ValidateSQLNullString(password)
		contact.Token = helper.ValidateSQLNullString(token)
		contact.Avatar = helper.ValidateSQLNullString(avatar)
		contact.Status = helper.ValidateSQLNullString(status)
		contact.CreatedBy = helper.ValidateSQLNullInt64(createdBy)
		contact.ModifiedBy = helper.ValidateSQLNullInt64(modifiedBy)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, nil)
			tags[helper.TextResponse] = err
			output <- ResultQuery{Error: err}
			return
		}

		tags["args"] = contact
		output <- ResultQuery{Result: contact}

	})

	return output
}

// GetListContact function for loading contact
func (mq *ContactQueryPostgres) GetListContact(ctxReq context.Context, params *model.ParametersContact) <-chan ResultQuery {
	ctx := "ContactQuery-GetListContact"
	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		strQuery, queryValues := mq.generateQuery(params)

		sq := fmt.Sprintf(`SELECT 
		id, "reference_id", "first_name", "last_name", "salutation", 
		"job_title", "email", "nav_contact_id", "is_primary", "birth_date", 
		"note", "created_at", "modified_at","created_by", "modified_by", 
		"account_id",  "is_disabled","avatar", "status", "is_new", 
		"phone_number"
		FROM b2b_contact %s
		ORDER BY id DESC
		LIMIT %d OFFSET %d`, strQuery, params.Limit, params.Offset)

		tags[helper.TextQuery] = sq
		rows, err := mq.db.Query(sq, queryValues...)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextQuery, err, params)
			output <- ResultQuery{Error: err}
			return
		}

		defer rows.Close()

		var listContact sharedModel.ListContact

		for rows.Next() {
			var (
				contact                                                           sharedModel.B2BContactData
				isNew                                                             sql.NullBool
				birthDate                                                         pq.NullTime
				referenceID, navContactID, accountID, phoneNumber, avatar, status sql.NullString
				createdBy, modifiedBy                                             sql.NullInt64
			)

			err = rows.Scan(
				&contact.ID, &referenceID, &contact.FirstName, &contact.LastName, &contact.Salutation,
				&contact.JobTitle, &contact.Email, &navContactID, &contact.IsPrimary, &birthDate,
				&contact.Note, &contact.CreatedAt, &contact.ModifiedAt, &createdBy, &modifiedBy,
				&accountID, &contact.IsDisabled, &avatar, &status, &isNew,
				&phoneNumber,
			)
			if err != nil {
				helper.SendErrorLog(ctxReq, ctx, helper.TextQuery, err, params)
				output <- ResultQuery{Error: err}
				return
			}

			if isNew.Valid {
				contact.IsNew = isNew.Bool
			}

			if birthDate.Valid {
				contact.BirthDate = birthDate.Time
			}

			contact.ReferenceID = helper.ValidateSQLNullString(referenceID)
			contact.NavContactID = helper.ValidateSQLNullString(navContactID)
			contact.AccountID = helper.ValidateSQLNullString(accountID)
			contact.PhoneNumber = helper.ValidateSQLNullString(phoneNumber)
			contact.Avatar = helper.ValidateSQLNullString(avatar)
			contact.Status = helper.ValidateSQLNullString(status)
			contact.CreatedBy = helper.ValidateSQLNullInt64(createdBy)
			contact.ModifiedBy = helper.ValidateSQLNullInt64(modifiedBy)

			listContact.Contact = append(listContact.Contact, &contact)
		}

		tags[helper.TextResponse] = listContact
		output <- ResultQuery{Result: listContact}
	})

	return output

}

// generateQuery function for generating query
func (mq *ContactQueryPostgres) generateQuery(params *model.ParametersContact) (string, []interface{}) {
	var (
		strQuery, status, isNew string
		queryStrOR              []string
		queryStrAND             []string
		queryValues             []interface{}
	)

	if len(params.Query) > 0 {
		queries := strings.Split(params.Query, " ")
		if len(queries) > 1 {
			queryValues = append(queryValues, "%"+params.Query+"%")
			queryStrOR = append(queryStrOR, `("first_name" || ' ' || "last_name"  || ' ' || "email" || "phone_number"  ilike $`+strconv.Itoa(len(queryStrOR)+1)+`)`)
		} else {
			queryStrOR = append(queryStrOR, `"first_name" ilike $`+strconv.Itoa(len(queryStrOR)+1))
			queryValues = append(queryValues, "%"+params.Query+"%")
			queryStrOR = append(queryStrOR, `"last_name" ilike $`+strconv.Itoa(len(queryStrOR)+1))
			queryValues = append(queryValues, "%"+params.Query+"%")
			queryStrOR = append(queryStrOR, `"email" ilike $`+strconv.Itoa(len(queryStrOR)+1))
			queryValues = append(queryValues, "%"+params.Query+"%")
			queryStrOR = append(queryStrOR, `"phone_number" ilike $`+strconv.Itoa(len(queryStrOR)+1))
			queryValues = append(queryValues, "%"+params.Query+"%")
		}
	}

	if len(params.IsNew) > 0 {
		intLentOR := 0
		if len(queryStrOR) > 0 {
			intLentOR = len(queryStrOR)
		}

		isNew = strconv.Itoa(intLentOR + 1)
		queryStrAND = append(queryStrAND, `"is_new" = $`+isNew)

		isNewBool, _ := strconv.ParseBool(params.IsNew)
		queryValues = append(queryValues, isNewBool)
	}

	if len(params.Status) > 0 {
		intLentAnd := 0
		if len(queryStrAND) > 0 {
			intLentAnd = len(queryStrAND) + len(queryStrOR)
		}
		status = strconv.Itoa(intLentAnd + 1)
		queryStrAND = append(queryStrAND, `"status" = $`+status)
		queryValues = append(queryValues, params.Status)
	}

	strQuery = mq.generateFilterQuery(queryStrOR, queryStrAND)

	return strQuery, queryValues
}

func (mq *ContactQueryPostgres) generateFilterQuery(queryStrOR, queryStrAND []string) string {
	var strQuery string
	if len(queryStrOR) > 0 || len(queryStrAND) > 0 {
		if len(queryStrOR) > 0 {
			strQuery = fmt.Sprintf(`(%s)`, strings.Join(queryStrOR, " OR "))
			queryStrAND = append(queryStrAND, strQuery)
		}

		if len(queryStrAND) > 0 {
			strQuery = strings.Join(queryStrAND, " AND ")
		}
	}

	if len(strQuery) > 0 {
		strQuery = fmt.Sprintf(" WHERE %s", strQuery)
	}
	return strQuery
}

// GetTotalContact function for getting total of contact
func (mq *ContactQueryPostgres) GetTotalContact(ctxReq context.Context, params *model.ParametersContact) <-chan ResultQuery {
	ctx := "ContactQuery-GetTotalContact"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		var totalData int

		strQuery, queryValues := mq.generateQuery(params)
		sq := fmt.Sprintf(`SELECT count(id) FROM b2b_contact %s`, strQuery)

		tags[helper.TextQuery] = sq
		err := mq.db.QueryRow(sq, queryValues...).Scan(&totalData)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextQuery, err, params)
			output <- ResultQuery{Error: err}
			return
		}

		tags[helper.TextResponse] = totalData
		output <- ResultQuery{Result: totalData}
	})

	return output
}

// FindContactMicrositeByEmail function for getting detail contact by email
func (mq *ContactQueryPostgres) FindContactMicrositeByEmail(ctxReq context.Context, email, transactionType, memberType string) <-chan ResultQuery {
	ctx := "ContactQuery-Microsite-FindContactMicrositeByEmail"

	outputFindByEmail := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(outputFindByEmail)
		querySelect := `
			SELECT b2b_contact.id, b2b_contact."reference_id", b2b_contact."first_name", b2b_contact."last_name", 
				b2b_contact."salutation", b2b_contact."job_title", b2b_contact."email", b2b_contact."nav_contact_id", 
				b2b_contact."is_primary", b2b_contact."birth_date", b2b_contact."note", b2b_contact."created_at", 
				b2b_contact."modified_at",b2b_contact."created_by", b2b_contact."modified_by", 
				b2b_contact."account_id",  b2b_contact."is_disabled", b2b_contact."password",b2b_contact."token", 
				b2b_contact."avatar", b2b_contact."status", b2b_contact."is_new", b2b_contact."phone_number",
				b2b_contact.transaction_type, elems->>'microsite' as microsite, b2b_contact."is_sync", 
				b2b_contact."salt", b2b_contact."last_password_modified"
			FROM
				b2b_contact,
				jsonb_array_elements(b2b_contact.transaction_type) elems
			WHERE
				elems->>'microsite' = $1 
			AND 
				elems->>'type' = $2
			AND
				LOWER(b2b_contact.email) = LOWER($3)`

		tags[helper.TextQuery] = querySelect
		stmt, err := mq.db.Prepare(querySelect)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, querySelect)
			tags[helper.TextResponse] = err
			outputFindByEmail <- ResultQuery{Error: err}
			return
		}
		defer stmt.Close()

		// initialize needed variables
		var (
			contactMicrosite                                                                                        sharedModel.B2BContactData
			isNew                                                                                                   sql.NullBool
			birthDate                                                                                               pq.NullTime
			referenceID, navContactID, accountID, phoneNumber, avatar, status, token, micrositeName, password, salt sql.NullString
			createdBy, modifiedBy                                                                                   sql.NullInt64
			transactionTypeDB                                                                                       []byte
		)

		err = stmt.QueryRow(memberType, transactionType, email).Scan(
			&contactMicrosite.ID, &referenceID, &contactMicrosite.FirstName, &contactMicrosite.LastName, &contactMicrosite.Salutation,
			&contactMicrosite.JobTitle, &contactMicrosite.Email, &navContactID, &contactMicrosite.IsPrimary, &birthDate,
			&contactMicrosite.Note, &contactMicrosite.CreatedAt, &contactMicrosite.ModifiedAt, &createdBy, &modifiedBy,
			&accountID, &contactMicrosite.IsDisabled, &password, &token, &avatar,
			&status, &isNew, &phoneNumber,
			&transactionTypeDB, &micrositeName, &contactMicrosite.IsSync,
			&salt, &contactMicrosite.LastPasswordModified,
		)

		if err != nil {
			jsonM := fmt.Sprintf(`{"email":"%s", "transactionType":"%s", "memberType":"%s"}`, email, transactionType, memberType)
			helper.SendNotification("Microsite Client Login", jsonM, ctx, err)
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, querySelect)
			tags[helper.TextResponse] = err
			outputFindByEmail <- ResultQuery{Error: err}
			return
		}
		contactMicrosite.Password = helper.ValidateSQLNullString(password)
		contactMicrosite.Salt = helper.ValidateSQLNullString(salt)
		contactMicrosite.IsNew = isNew.Bool
		contactMicrosite.Status = helper.ValidateSQLNullString(status)
		contactMicrosite.PhoneNumber = helper.ValidateSQLNullString(phoneNumber)
		contactMicrosite.Avatar = helper.ValidateSQLNullString(avatar)

		outputFindByEmail <- ResultQuery{Result: contactMicrosite}

	})

	return outputFindByEmail
}

// FindContactByEmail function for getting detail contact by email
func (mq *ContactQueryPostgres) FindContactByEmail(ctxReq context.Context, email string) <-chan ResultQuery {
	ctx := "ContactQuery-Microsite-FindContactByEmail"

	outputFindByEmail := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(outputFindByEmail)
		querySelect := `
			SELECT b2b_contact.id, b2b_contact."reference_id", b2b_contact."first_name", b2b_contact."last_name", 
				b2b_contact."salutation", b2b_contact."job_title", b2b_contact."email", b2b_contact."nav_contact_id", 
				b2b_contact."is_primary", b2b_contact."birth_date", b2b_contact."note", b2b_contact."created_at", 
				b2b_contact."modified_at",b2b_contact."created_by", b2b_contact."modified_by", 
				b2b_contact."account_id",  b2b_contact."is_disabled", b2b_contact."password",b2b_contact."token", 
				b2b_contact."avatar", b2b_contact."status", b2b_contact."is_new", b2b_contact."phone_number", 
				b2b_contact.transaction_type, b2b_contact."is_sync", 
				b2b_contact."salt", b2b_contact."last_password_modified"
			FROM
				b2b_contact
			WHERE 
				LOWER(b2b_contact.email) = LOWER($1)`

		tags[helper.TextQuery] = querySelect
		stmt, err := mq.db.Prepare(querySelect)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, querySelect)
			tags[helper.TextResponse] = err
			outputFindByEmail <- ResultQuery{Error: err}
			return
		}
		defer stmt.Close()

		// initialize needed variables
		var (
			contactMicrosite                                                                         sharedModel.B2BContactData
			isNew                                                                                    sql.NullBool
			birthDate                                                                                pq.NullTime
			referenceID, navContactID, accountID, phoneNumber, avatar, status, token, password, salt sql.NullString
			createdBy, modifiedBy                                                                    sql.NullInt64
			transactionTypeDB                                                                        []byte
		)

		err = stmt.QueryRow(email).Scan(
			&contactMicrosite.ID, &referenceID, &contactMicrosite.FirstName, &contactMicrosite.LastName, &contactMicrosite.Salutation,
			&contactMicrosite.JobTitle, &contactMicrosite.Email, &navContactID, &contactMicrosite.IsPrimary, &birthDate,
			&contactMicrosite.Note, &contactMicrosite.CreatedAt, &contactMicrosite.ModifiedAt, &createdBy, &modifiedBy,
			&accountID, &contactMicrosite.IsDisabled, &password, &token, &avatar,
			&status, &isNew, &phoneNumber, &transactionTypeDB, &contactMicrosite.IsSync,
			&salt, &contactMicrosite.LastPasswordModified,
		)

		// convert transactionType jsonb to struct
		var transactionType []sharedModel.TransactionType
		json.Unmarshal(transactionTypeDB, &transactionType)
		contactMicrosite.TransactionType = transactionType

		if err != nil {
			jsonM := fmt.Sprintf(`{"email":"%s"}`, email)
			helper.SendNotification("Microsite Client Login", jsonM, ctx, err)
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, querySelect)
			tags[helper.TextResponse] = err
			outputFindByEmail <- ResultQuery{Error: err}
			return
		}
		contactMicrosite.Password = helper.ValidateSQLNullString(password)
		contactMicrosite.Salt = helper.ValidateSQLNullString(salt)
		contactMicrosite.IsNew = isNew.Bool
		contactMicrosite.Status = helper.ValidateSQLNullString(status)
		contactMicrosite.PhoneNumber = helper.ValidateSQLNullString(phoneNumber)
		contactMicrosite.Avatar = helper.ValidateSQLNullString(avatar)

		outputFindByEmail <- ResultQuery{Result: contactMicrosite}

	})

	return outputFindByEmail
}

// FindAccountByMemberType function for getting detail account by memberType
func (mq *ContactQueryPostgres) FindAccountByMemberType(ctxReq context.Context, memberType string) <-chan ResultQuery {
	ctx := "ContactQuery-Microsite-FindAccountByMemberType"

	outputFindByMemberType := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(outputFindByMemberType)
		querySelect := `SELECT * FROM b2b_account WHERE "member_type"=$1 AND "is_delete"=false AND "is_disabled"=false AND "status"='ACTIVATED'`

		tags[helper.TextQuery] = querySelect
		stmt, err := mq.db.Prepare(querySelect)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, querySelect)
			tags[helper.TextResponse] = err
			outputFindByMemberType <- ResultQuery{Error: err}
			return
		}
		defer stmt.Close()

		// initialize needed variables
		var (
			accountData sharedModel.B2BAccountCDC
		)

		err = stmt.QueryRow(memberType).Scan(
			&accountData.ID, &accountData.ParentID, &accountData.NavID, &accountData.LegalEntity, &accountData.Name, &accountData.IndustryID,
			&accountData.NumberEmployee, &accountData.OfficeEmployee, &accountData.BusinessSize, &accountData.BusinessGroup,
			&accountData.EstablishedYear, &accountData.IsDelete, &accountData.TermOfPayment, &accountData.CustomerCategory, &accountData.UserID,
			&accountData.CreatedAt, &accountData.ModifiedAt, &accountData.CreatedBy, &accountData.ModifiedBy, &accountData.AccountGroupID,
			&accountData.IsDisabled, &accountData.Status, &accountData.PaymentMethodID, &accountData.PaymentMethodType,
			&accountData.SubPaymentMethodName, &accountData.IsCf, &accountData.Logo, &accountData.IsParent, &accountData.IsMicrosite,
			&accountData.Membertype, &accountData.ErpID,
		)

		if err != nil {
			jsonM := fmt.Sprintf(`{"memberType":"%s"}`, memberType)
			helper.SendNotification("Microsite Bela Client Login", jsonM, ctx, err)
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, querySelect)
			tags[helper.TextResponse] = err
			outputFindByMemberType <- ResultQuery{Error: err}
			return
		}

		outputFindByMemberType <- ResultQuery{Result: accountData}

	})

	return outputFindByMemberType
}

// FindContactCorporateByEmail function for getting detail contact by email
func (mq *ContactQueryPostgres) FindContactCorporateByEmail(ctxReq context.Context, email string) <-chan ResultQuery {
	ctx := "ContactQuery-Microsite-FindContactCorporateByEmail"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		q := `
		SELECT
			b2b_contact.id, b2b_contact."reference_id", b2b_contact."first_name", b2b_contact."last_name",
			b2b_contact."salutation", b2b_contact."job_title", b2b_contact."email", b2b_contact."nav_contact_id",
			b2b_contact."is_primary", b2b_contact."birth_date", b2b_contact."note", b2b_contact."created_at",
			b2b_contact."modified_at",b2b_contact."created_by", b2b_contact."modified_by",
			b2b_contact."account_id",  b2b_contact."is_disabled", b2b_contact."password",b2b_contact."token",
			b2b_contact."avatar", b2b_contact."status", b2b_contact."is_new", b2b_contact."phone_number",
			elems->>'type' as login_type, b2b_contact."is_sync", 
			b2b_contact."salt", b2b_contact."last_password_modified"
		FROM
			b2b_contact,
			jsonb_array_elements(
				case jsonb_typeof(b2b_contact.transaction_type)
					when 'array' then b2b_contact.transaction_type
				else '[]' end
		) as elems
		WHERE			
			b2b_contact.is_disabled = false
		AND elems->>'microsite' = $2
		AND
			LOWER(b2b_contact.email) = LOWER($1)`

		tags[helper.TextQuery] = q
		stmt, err := mq.db.Prepare(q)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, q)
			tags[helper.TextResponse] = err
			output <- ResultQuery{Error: err}
			return
		}
		defer stmt.Close()

		var (
			contactCorporat                                                                                     sharedModel.B2BContactData
			isNew                                                                                               sql.NullBool
			birthDate                                                                                           pq.NullTime
			referenceID, navContactID, accountID, phoneNumber, avatar, status, token, password, loginType, salt sql.NullString
			createdBy, modifiedBy                                                                               sql.NullInt64
		)

		if err = stmt.QueryRow(email, model.LoginTypeCorporate).Scan(
			&contactCorporat.ID, &referenceID, &contactCorporat.FirstName, &contactCorporat.LastName, &contactCorporat.Salutation,
			&contactCorporat.JobTitle, &contactCorporat.Email, &navContactID, &contactCorporat.IsPrimary, &birthDate,
			&contactCorporat.Note, &contactCorporat.CreatedAt, &contactCorporat.ModifiedAt, &createdBy, &modifiedBy,
			&accountID, &contactCorporat.IsDisabled, &password, &token, &avatar,
			&status, &isNew, &phoneNumber, &loginType, &contactCorporat.IsSync,
			&salt, &contactCorporat.LastPasswordModified,
		); err != nil {
			helper.SendNotification("Corporate Client Login", email, ctx, err)
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, q)
			tags[helper.TextResponse] = err
			output <- ResultQuery{Error: err}
			return
		}
		contactCorporat.Password = helper.ValidateSQLNullString(password)
		if createdBy.Valid {
			contactCorporat.CreatedBy = createdBy.Int64
		}

		contactCorporat.IsNew = isNew.Bool
		contactCorporat.Status = helper.ValidateSQLNullString(status)
		contactCorporat.PhoneNumber = helper.ValidateSQLNullString(phoneNumber)
		contactCorporat.Avatar = helper.ValidateSQLNullString(avatar)
		contactCorporat.LoginType = helper.ValidateSQLNullString(loginType)
		contactCorporat.Salt = helper.ValidateSQLNullString(salt)
		tags["result"] = contactCorporat

		output <- ResultQuery{Result: contactCorporat}

	})

	return output
}
func (mq *ContactQueryPostgres) GetTransactionType(ctxReq context.Context, email string) <-chan ResultQuery {
	ctx := "ContactQuery-GetTransactionType"
	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		sq := `SELECT elems->>'type' as tipe
		FROM
			b2b_contact,
			jsonb_array_elements(b2b_contact.transaction_type) elems
		WHERE
			elems->>'microsite' = $1
		and 
			email = $2`

		tags[helper.TextQuery] = sq
		stmt, err := mq.db.Prepare(sq)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextQuery, err, sq)
			tags[helper.TextResponse] = err
			output <- ResultQuery{Error: err}
			return
		}

		defer stmt.Close()

		var (
			transaction_type string
		)

		if err = stmt.QueryRow(model.LoginTypeCorporate, email).Scan(
			&transaction_type,
		); err != nil {
			helper.SendNotification("Corporate Client Login", email, ctx, err)
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, sq)
			tags[helper.TextResponse] = err
			output <- ResultQuery{Error: err}
			return
		}
		tags["result"] = transaction_type

		output <- ResultQuery{Result: transaction_type}
	})

	return output

}
