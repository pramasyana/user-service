package query

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/lib/pq"
)

// MemberQueryPostgres data structure
type MemberQueryPostgres struct {
	db *sql.DB
}

// NewMemberQueryPostgres function for initializing member query
func NewMemberQueryPostgres(db *sql.DB) *MemberQueryPostgres {
	return &MemberQueryPostgres{db: db}
}

// FindByID function for getting detail member by uid
func (mq *MemberQueryPostgres) FindByID(ctxReq context.Context, uid string) <-chan ResultQuery {
	ctx := "MemberQuery-FindByID"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		var (
			filterByID  string
			queryValues []interface{}
		)

		queryValues = append(queryValues, uid)
		filterByID = `WHERE "id" = $1`

		tags[helper.TextParameter] = filterByID
		member := <-mq.FindMember(ctxReq, filterByID, queryValues)
		if member.Error != nil {
			output <- ResultQuery{Error: member.Error}
			return
		}

		memberResult := member.Result.(model.Member)

		output <- ResultQuery{Result: memberResult}

	})

	return output
}

// BulkFindByEmail return members
func (mq *MemberQueryPostgres) BulkFindByEmail(ctxReq context.Context, emails []string) <-chan ResultQuery {
	ctx := "MemberQuery-BulkFindByEmail"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		var (
			queryValues []interface{}
		)
		n := 0
		bindVar := []string{}
		for _, m := range emails {
			n++
			queryValues = append(queryValues, strings.ToLower(m))
			bindVar = append(bindVar, "$"+strconv.Itoa(n))
		}
		bulkQuery := `SELECT "email", "mobile" FROM member WHERE "email" IN (` + strings.Join(bindVar, ",") + `)`
		tags[helper.TextQuery] = bulkQuery
		tags[helper.TextArgs] = queryValues

		rows, err := mq.db.Query(bulkQuery, queryValues...)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, queryValues)
			output <- ResultQuery{Error: err}
			return
		}
		defer rows.Close()
		var members []model.Member
		for rows.Next() {
			var member model.Member
			if err = rows.Scan(&member.Email, &member.Mobile); err != nil {
				helper.SendErrorLog(ctxReq, "MemberQuery-BulkFindByEmail-ScanRow", helper.ScopeParseResponse, err, member)
				output <- ResultQuery{Error: err}
				return
			}
			members = append(members, member)
		}
		tags[helper.TextResponse] = members
		output <- ResultQuery{Result: members}
	})

	return output
}

// FindByEmail function for getting detail member by email
func (mq *MemberQueryPostgres) FindByEmail(ctxReq context.Context, email string) <-chan ResultQuery {
	ctx := "MemberQuery-FindByEmail"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {

		defer close(output)
		var (
			filter      string
			queryValues []interface{}
		)

		email = strings.ToLower(email)
		queryValues = append(queryValues, email)
		filter = `WHERE "email" = $1`

		tags[helper.TextParameter] = filter
		member := <-mq.FindMember(ctxReq, filter, queryValues)
		if member.Error != nil {
			output <- ResultQuery{Error: member.Error}
			return
		}

		memberResult := member.Result.(model.Member)

		output <- ResultQuery{Result: memberResult}
	})

	return output
}

// FindByMobile function for getting detail member by mobile
func (mq *MemberQueryPostgres) FindByMobile(ctxReq context.Context, mobile string) <-chan ResultQuery {
	ctx := "MemberQuery-FindByMobile"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		var (
			filterByMobile string
			queryValues    []interface{}
		)

		queryValues = append(queryValues, mobile)
		filterByMobile = `WHERE "mobile" = $1`

		tags[helper.TextParameter] = filterByMobile
		member := <-mq.FindMember(ctxReq, filterByMobile, queryValues)
		if member.Error != nil {
			output <- ResultQuery{Error: member.Error}
			return
		}

		memberResult := member.Result.(model.Member)

		output <- ResultQuery{Result: memberResult}
	})

	return output
}

// FindMaxID function for getting member max id
func (mq *MemberQueryPostgres) FindMaxID(ctxReq context.Context) <-chan ResultQuery {
	ctx := "MemberQuery-FindMaxID"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		// Get date now format YYmm and defined to id format USR + YYmm
		dateFormat := fmt.Sprintf(`USR%s`, time.Now().UTC().Format("0601"))
		whereDate := dateFormat + "%"

		q := fmt.Sprintf(`
			SELECT id FROM "member" 
			WHERE id LIKE '%s' 
			ORDER BY id desc limit 1`,
			whereDate)

		tags[helper.TextQuery] = q

		stmt, err := mq.db.Prepare(q)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, q)
			output <- ResultQuery{Error: err}
			return
		}

		defer stmt.Close()

		var maxID string
		err = stmt.QueryRow().Scan(&maxID)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, q)
			if err == sql.ErrNoRows {
				output <- ResultQuery{Result: "0"}
				return
			}
			output <- ResultQuery{Error: err}
			return
		}

		// split id
		getMaxID := maxID[7:]
		output <- ResultQuery{Result: getMaxID}
	})

	return output
}

// FindByToken function for getting detail member by token
func (mq *MemberQueryPostgres) FindByToken(ctxReq context.Context, token string) <-chan ResultQuery {
	ctx := "MemberQuery-FindByToken"
	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		var (
			filter      string
			queryValues []interface{}
		)

		queryValues = append(queryValues, token)
		filter = `WHERE "token" = $1`

		tags[helper.TextParameter] = filter
		member := <-mq.FindMember(ctxReq, filter, queryValues)
		if member.Error != nil {
			output <- ResultQuery{Error: member.Error}
			return
		}

		memberResult := member.Result.(model.Member)

		output <- ResultQuery{Result: memberResult}
	})
	return output
}

// FindMember function for getting detail member by token
func (mq *MemberQueryPostgres) FindMember(ctxReq context.Context, filter string, queryValues []interface{}) <-chan ResultQuery {
	ctx := "MemberQuery-FindMember"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		querySelect := fmt.Sprintf(`SELECT id, email, "firstName", "lastName",
			gender, mobile, phone, ext, "birthDate",
			password, salt, token,
			province, "provinceId", city, "cityId",
			district, "districtId", "subDistrict", "subDistrictId", "zipCode", address,
			"jobTitle", department,
			"isAdmin", "isStaff", status, "facebookId", "googleId", "appleId", "azureId", "ldapId", "signUpFrom",
			"lastLogin", "lastBlocked", "created", "lastModified", version, "lastTokenAttempt", "profilePicture", "mfaEnabled", "mfaKey", 
			"lastPasswordModified", "facebookConnect", "googleConnect", "appleConnect", "mfaAdminEnabled", "mfaAdminKey", "isSync"
		FROM member %s`, filter)
		tags["filter"] = filter

		stmt, err := mq.db.Prepare(querySelect)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, querySelect)
			output <- ResultQuery{Error: err}
			return
		}

		defer stmt.Close()

		// initialize needed variables
		var (
			member                                                             model.Member
			lastName, mobile, phone, ext, password, salt, gender, token        sql.NullString
			province, city, district, subDistrict, zipCode, address            sql.NullString
			provinceID, cityID, districtID, subDistrictID                      sql.NullString
			facebookID, googleID, appleID, azureID, ldapID, signUpFrom         sql.NullString
			jobTitle, department, profilePicture, mfaKey                       sql.NullString
			birthDate, lastLogin, lastModified, lastBlocked, lastTokenAttempt  pq.NullTime
			lastPasswordModified, facebookConnect, googleConnect, appleConnect pq.NullTime
			status                                                             string
			mfaAdminKey                                                        sql.NullString
		)

		err = stmt.QueryRow(queryValues...).Scan(
			&member.ID, &member.Email, &member.FirstName, &lastName,
			&gender, &mobile, &phone, &ext, &birthDate,
			&password, &salt, &token,
			&province, &provinceID, &city, &cityID,
			&district, &districtID, &subDistrict, &subDistrictID, &zipCode, &address,
			&jobTitle, &department,
			&member.IsAdmin, &member.IsStaff, &status, &facebookID, &googleID, &appleID, &azureID, &ldapID, &signUpFrom,
			&lastLogin, &lastBlocked, &member.Created, &lastModified, &member.Version, &lastTokenAttempt, &profilePicture,
			&member.MFAEnabled, &mfaKey, &lastPasswordModified, &facebookConnect, &googleConnect, &appleConnect,
			&member.AdminMFAEnabled, &mfaAdminKey, &member.IsSync,
		)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, querySelect)
			output <- ResultQuery{Error: err}
			return
		}

		// assign the nullable field to object
		member.LastName = helper.ValidateSQLNullString(lastName)
		member.Mobile = helper.ValidateSQLNullString(mobile)
		member.Phone = helper.ValidateSQLNullString(phone)
		member.Ext = helper.ValidateSQLNullString(ext)
		member.Password = helper.ValidateSQLNullString(password)
		member.Salt = helper.ValidateSQLNullString(salt)
		member.Token = helper.ValidateSQLNullString(token)
		member.JobTitle = helper.ValidateSQLNullString(jobTitle)
		member.Department = helper.ValidateSQLNullString(department)
		member.SignUpFrom = helper.ValidateSQLNullString(signUpFrom)
		member.ProfilePicture = helper.ValidateSQLNullString(profilePicture)
		member.MFAKey = helper.ValidateSQLNullString(mfaKey)
		member.MFAAdminKey = helper.ValidateSQLNullString(mfaAdminKey)

		if gender.Valid {
			member.Gender = model.StringToGender(gender.String)
			member.GenderString = gender.String
		}

		if birthDate.Valid {
			member.BirthDate = birthDate.Time
			member.BirthDateString = birthDate.Time.Format("02/01/2006")
		}

		member = mq.adjustLastDateData(member, lastLogin, lastBlocked, lastModified, lastPasswordModified, lastTokenAttempt)

		ma := model.Address{}
		ma.Province = helper.ValidateSQLNullString(province)
		ma.ProvinceID = helper.ValidateSQLNullString(provinceID)
		ma.City = helper.ValidateSQLNullString(city)
		ma.CityID = helper.ValidateSQLNullString(cityID)
		ma.District = helper.ValidateSQLNullString(district)
		ma.DistrictID = helper.ValidateSQLNullString(districtID)
		ma.SubDistrict = helper.ValidateSQLNullString(subDistrict)
		ma.SubDistrictID = helper.ValidateSQLNullString(subDistrictID)
		ma.ZipCode = helper.ValidateSQLNullString(zipCode)
		ma.Address = helper.ValidateSQLNullString(address)

		// parse get streets
		street1, street2 := helper.SplitStreetAddress(ma.Address)
		ma.Street1 = street1
		ma.Street2 = street2

		member.Address = ma

		ms := model.SocialMedia{
			FacebookID: facebookID.String,
			GoogleID:   googleID.String,
			AzureID:    azureID.String,
			LDAPID:     ldapID.String,
		}

		// adjustSocmedData restruct socmed data with null values
		member = mq.adjustSocmedData(member, ms, facebookConnect, googleConnect, appleConnect)

		member.Status = model.StringToStatus(status)
		member.StatusString = status
		member.CreatedString = member.Created.Format(time.RFC3339)

		output <- ResultQuery{Result: member}
	})

	return output
}

func (mq *MemberQueryPostgres) adjustLastDateData(member model.Member, lastLogin, lastBlocked, lastModified, lastPasswordModified, lastTokenAttempt pq.NullTime) model.Member {
	if lastLogin.Valid {
		member.LastLogin = lastLogin.Time
		member.LastLoginString = member.LastLogin.Format(time.RFC3339)
	}

	if lastBlocked.Valid {
		member.LastBlocked = lastBlocked.Time
		member.LastBlockedString = member.LastBlocked.Format(time.RFC3339)
	}

	if lastModified.Valid {
		member.LastModified = lastModified.Time
		member.LastModifiedString = member.LastModified.Format(time.RFC3339)
	}

	if lastPasswordModified.Valid {
		member.LastPasswordModified = lastPasswordModified.Time
		member.LastPasswordModifiedString = member.LastPasswordModified.Format(time.RFC3339)
	}

	if lastTokenAttempt.Valid {
		member.LastTokenAttempt = lastTokenAttempt.Time
		member.LastTokenAttemptString = member.LastTokenAttempt.Format(time.RFC3339)
	}

	return member
}

// adjustSocmedData function for restruct socmed data with null values
func (mq *MemberQueryPostgres) adjustSocmedData(member model.Member, ms model.SocialMedia, facebookConnect, googleConnect, appleConnect pq.NullTime) model.Member {

	if facebookConnect.Valid {
		member.SocialMedia.FacebookConnect = facebookConnect.Time
		ms.FacebookConnect = member.SocialMedia.FacebookConnect
		member.SocialMedia.FacebookConnectString = member.SocialMedia.FacebookConnect.Format(time.RFC3339)
		ms.FacebookConnectString = member.SocialMedia.FacebookConnect.Format(time.RFC3339)
	}

	if googleConnect.Valid {
		member.SocialMedia.GoogleConnect = googleConnect.Time
		ms.GoogleConnect = member.SocialMedia.GoogleConnect
		member.SocialMedia.GoogleConnectString = member.SocialMedia.GoogleConnect.Format(time.RFC3339)
		ms.GoogleConnectString = member.SocialMedia.GoogleConnect.Format(time.RFC3339)
	}

	if appleConnect.Valid {
		member.SocialMedia.AppleConnect = appleConnect.Time
		ms.AppleConnect = member.SocialMedia.AppleConnect
		member.SocialMedia.AppleConnectString = member.SocialMedia.AppleConnect.Format(time.RFC3339)
		ms.AppleConnectString = member.SocialMedia.AppleConnect.Format(time.RFC3339)
	}

	member.SocialMedia = ms
	return member
}

// UpdateBlockedMember function for updating status only
func (mq *MemberQueryPostgres) UpdateBlockedMember(email string) <-chan ResultQuery {
	ctx := "MemberQuery-UpdateBlockedMember"

	output := make(chan ResultQuery)
	go func() {
		defer close(output)

		sq := `UPDATE member SET status = 'BLOCKED', "lastBlocked" = NOW() WHERE email = $1`

		tx, err := mq.db.Begin()
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(context.Background(), ctx, "update_blocked_member", err, email)
			output <- ResultQuery{Error: err}
			return
		}
		stmt, err := tx.Prepare(sq)
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(context.Background(), ctx, "update_blocked_member", err, email)
			output <- ResultQuery{Error: err}
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(email)
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, email)
			output <- ResultQuery{Error: err}
			return
		}

		// commit statement
		tx.Commit()

		output <- ResultQuery{Error: nil}
	}()

	return output
}

// UnblockMember function for updating status from BLOCKED TO ACTIVE
func (mq *MemberQueryPostgres) UnblockMember(email string) <-chan ResultQuery {
	ctx := "MemberQuery-UnblockMember"

	output := make(chan ResultQuery)
	go func() {
		defer close(output)

		sq := `UPDATE member SET status = 'ACTIVE' WHERE email = $1`

		tx, err := mq.db.Begin()
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(context.Background(), ctx, helper.TextDBBegin, err, email)
			output <- ResultQuery{Error: err}
			return
		}

		stmt, err := tx.Prepare(sq)
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(context.Background(), ctx, helper.TextPrepareDatabase, err, email)
			output <- ResultQuery{Error: err}
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(email)
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(context.Background(), ctx, helper.TextExecQuery, err, email)
			output <- ResultQuery{Error: err}
			return
		}

		// commit statement
		tx.Commit()

		output <- ResultQuery{Error: nil}
	}()

	return output
}

// GetListMembers function for getting list of members
func (mq *MemberQueryPostgres) GetListMembers(ctxReq context.Context, params *model.Parameters) <-chan ResultQuery {
	ctx := "MemberQuery-GetListMembers"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		strQuery, queryValues := mq.generateQuery(params)

		if len(params.OrderBy) > 0 {
			params.OrderBy = fmt.Sprintf(`"%s"`, params.OrderBy)
		}
		tags["params"] = queryValues

		querySelect := fmt.Sprintf(`
			SELECT id, "firstName", "lastName", email, gender, mobile, phone, ext, "birthDate",
			password, salt, token,province, "provinceId", city, "cityId",
			district, "districtId", "subDistrict", "subDistrictId", "zipCode",
			address, "jobTitle", department, "isAdmin", "isStaff", status, 
			"facebookId", "googleId",  "appleId", "azureId", "signUpFrom",
			"lastLogin", "lastBlocked", "created", "lastModified", "lastPasswordModified", "lastTokenAttempt", 
			version, "mfaEnabled", "mfaKey"
			FROM member 
			%s
			ORDER BY %s %s
			LIMIT %d OFFSET %d`, strQuery, params.OrderBy, params.Sort, params.Limit, params.Offset)

		rows, err := mq.db.Query(querySelect, queryValues...)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, params)
			output <- ResultQuery{Error: nil}
			return
		}
		defer rows.Close()

		var members model.ListMembers
		for rows.Next() {
			var (
				member                                                      model.Member
				lastName, mobile, phone, ext, password, salt, gender, token sql.NullString
				province, city, district, subDistrict, zipCode, address     sql.NullString
				provinceID, cityID, districtID, subDistrictID               sql.NullString
				facebookID, googleID, appleID, azureID, signUpFrom          sql.NullString
				jobTitle, department, mfaKey                                sql.NullString
				birthDate, lastLogin, lastModified, lastBlocked             pq.NullTime
				lastPasswordModified, lastTokenAttempt                      pq.NullTime
				status                                                      string
			)

			err = rows.Scan(
				&member.ID, &member.FirstName, &lastName, &member.Email,
				&gender, &mobile, &phone, &ext, &birthDate,
				&password, &salt, &token,
				&province, &provinceID, &city, &cityID,
				&district, &districtID, &subDistrict, &subDistrictID, &zipCode, &address,
				&jobTitle, &department,
				&member.IsAdmin, &member.IsStaff, &status, &facebookID, &googleID, &appleID, &azureID, &signUpFrom,
				&lastLogin, &lastBlocked, &member.Created, &lastModified, &lastPasswordModified, &lastTokenAttempt,
				&member.Version, &member.MFAEnabled, &mfaKey,
			)

			if err != nil {
				helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, params)
				output <- ResultQuery{Error: err}
				return
			}

			// assign the nullable field to object
			member.LastName = helper.ValidateSQLNullString(lastName)
			member.Mobile = helper.ValidateSQLNullString(mobile)
			member.Phone = helper.ValidateSQLNullString(phone)
			member.Ext = helper.ValidateSQLNullString(ext)
			member.Password = ""
			member.Salt = ""
			member.Token = helper.ValidateSQLNullString(token)
			member.JobTitle = helper.ValidateSQLNullString(jobTitle)
			member.Department = helper.ValidateSQLNullString(department)
			member.SignUpFrom = helper.ValidateSQLNullString(signUpFrom)
			member.MFAKey = helper.ValidateSQLNullString(mfaKey)

			member.SetGender(gender)
			member.SetBirthDate(birthDate)
			member.SetHasPassword(password)

			member = mq.adjustLastDateData(member, lastLogin, lastBlocked, lastModified, lastPasswordModified, lastTokenAttempt)

			member.Status = model.StringToStatus(status)
			member.StatusString = status
			member.CreatedString = member.Created.Format(time.RFC3339)

			ma := model.Address{}

			ma.Province = helper.ValidateSQLNullString(province)
			ma.ProvinceID = helper.ValidateSQLNullString(provinceID)
			ma.City = helper.ValidateSQLNullString(city)
			ma.CityID = helper.ValidateSQLNullString(cityID)
			ma.District = helper.ValidateSQLNullString(district)
			ma.DistrictID = helper.ValidateSQLNullString(districtID)
			ma.SubDistrict = helper.ValidateSQLNullString(subDistrict)
			ma.SubDistrictID = helper.ValidateSQLNullString(subDistrictID)
			ma.ZipCode = helper.ValidateSQLNullString(zipCode)
			ma.Address = helper.ValidateSQLNullString(address)

			// parse get streets
			street1, street2 := helper.SplitStreetAddress(ma.Address)
			ma.Street1 = street1
			ma.Street2 = street2
			member.Address = ma

			// parse address
			ms := model.SocialMedia{
				FacebookID: facebookID.String,
				GoogleID:   googleID.String,
				AppleID:    appleID.String,
				AzureID:    azureID.String,
			}
			member.SocialMedia = ms

			members.Members = append(members.Members, &member)
		}

		output <- ResultQuery{Result: members}
	})

	return output
}

// GetTotalMembers function for getting total of members
func (mq *MemberQueryPostgres) GetTotalMembers(params *model.Parameters) <-chan ResultQuery {
	ctx := "MemberQuery-GetTotalMembers"

	output := make(chan ResultQuery)
	go func() {
		defer close(output)

		var totalData int

		strQuery, queryValues := mq.generateQuery(params)

		sq := fmt.Sprintf(`SELECT count(id) FROM member %s`, strQuery)
		err := mq.db.QueryRow(sq, queryValues...).Scan(&totalData)
		if err != nil {
			helper.SendErrorLog(context.Background(), ctx, "get_row_from_query", err, params)
			output <- ResultQuery{Error: err}
			return
		}

		output <- ResultQuery{Result: totalData}
	}()

	return output
}

// generateQuery function for generating query
func (mq *MemberQueryPostgres) generateQuery(params *model.Parameters) (string, []interface{}) {
	var (
		strQuery, idx string
		queryStrOR    []string
		queryStrAND   []string
		queryValues   []interface{}
		lq            int
	)

	if len(params.Query) > 0 {
		queries := strings.Split(params.Query, " ")
		lq = len(queries)
		if lq > 1 {
			queryValues = append(queryValues, "%"+params.Query+"%")
			queryStrOR = append(queryStrOR, `("firstName" || ' ' || "lastName"  || ' ' || "mobile"   || ' ' || "id"  ilike $`+strconv.Itoa(len(queryStrOR)+1)+`)`)

		} else {
			queryStrOR = append(queryStrOR, `"firstName" ilike $`+strconv.Itoa(len(queryStrOR)+1))
			queryValues = append(queryValues, "%"+params.Query+"%")
			queryStrOR = append(queryStrOR, `"lastName" ilike $`+strconv.Itoa(len(queryStrOR)+1))
			queryValues = append(queryValues, "%"+params.Query+"%")
			queryStrOR = append(queryStrOR, `"mobile" ilike $`+strconv.Itoa(len(queryStrOR)+1))
			queryValues = append(queryValues, "%"+params.Query+"%")
			queryStrOR = append(queryStrOR, `"id" ilike $`+strconv.Itoa(len(queryStrOR)+1))
			queryValues = append(queryValues, "%"+params.Query+"%")
		}

		idx = strconv.Itoa(len(queryStrOR) + 1)

		queryStrOR = append(queryStrOR, `"email" ilike $`+idx)
		queryValues = append(queryValues, "%"+params.Query+"%")
	}

	// generate additional query
	queryStrAND, queryValues = mq.generateAdditionalQuery(params, queryStrOR, queryStrAND, queryValues, idx)

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

	return strQuery, queryValues
}

func (mq *MemberQueryPostgres) generateAdditionalQuery(params *model.Parameters, queryStrOR, queryStrAND []string, queryValues []interface{}, idx string) ([]string, []interface{}) {
	if len(params.Status) > 0 {
		intLentOR := 0
		if len(queryStrOR) > 0 {
			intLentOR = len(queryStrOR)
		}

		idx := strconv.Itoa(intLentOR + 1)

		queryStrAND = append(queryStrAND, `status = $`+idx)
		queryValues = append(queryValues, params.Status)
	}
	if len(params.Email) > 0 {
		intLentAnd := 0
		if len(queryStrAND) > 0 {
			intLentAnd = len(queryStrAND) + len(queryStrOR)
		}

		idx := strconv.Itoa(intLentAnd + 1)
		queryStrAND = append(queryStrAND, `email = $`+idx)
		queryValues = append(queryValues, params.Email)
	}

	if len(params.IsStaff) > 0 {
		intLentAnd := 0
		if len(queryStrAND) > 0 {
			intLentAnd = len(queryStrAND) + len(queryStrOR)
		}

		idx := strconv.Itoa(intLentAnd + 1)
		queryStrAND = append(queryStrAND, `"isStaff" = $`+idx)
		queryValues = append(queryValues, params.IsStaff)
	}

	if len(params.IsAdmin) > 0 {
		intLentAnd := 0
		if len(queryStrAND) > 0 {
			intLentAnd = len(queryStrAND) + len(queryStrOR)
		}

		idx := strconv.Itoa(intLentAnd + 1)
		queryStrAND = append(queryStrAND, `"isAdmin" = $`+idx)
		queryValues = append(queryValues, params.IsAdmin)
	}
	if len(params.UserID) > 0 {
		intLentAnd := 0
		if len(queryStrAND) > 0 {
			intLentAnd = len(queryStrAND) + len(queryStrOR)
		}

		idx := strconv.Itoa(intLentAnd + 1)
		queryStrAND = append(queryStrAND, `"id" = $`+idx)
		queryValues = append(queryValues, params.UserID)
	}

	return queryStrAND, queryValues
}

// UpdateLastTokenAttempt function for updating last token attempt only
func (mq *MemberQueryPostgres) UpdateLastTokenAttempt(ctxReq context.Context, email string) <-chan ResultQuery {
	ctx := "MemberQuery-UpdateLastTokenAttempt"

	output := make(chan ResultQuery)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		sq := `UPDATE member SET "lastTokenAttempt" = NOW() WHERE email = $1`

		tags[helper.TextQuery] = sq
		tx, err := mq.db.Begin()
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(ctxReq, ctx, helper.TextDBBegin, err, email)
			output <- ResultQuery{Error: err}
			return
		}

		stmt, err := tx.Prepare(sq)
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, email)
			output <- ResultQuery{Error: err}
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(email)
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, email)
			output <- ResultQuery{Error: err}
			return
		}

		// commit statement
		tx.Commit()

		output <- ResultQuery{Error: nil}
	})

	return output
}
