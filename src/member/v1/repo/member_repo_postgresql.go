package repo

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
	"github.com/lib/pq"
)

const (
	// number of column inserted
	queryPattern = `(
		?,?,?,?,?,?,?,?,?,?,
		?,?,?,?,?,?,?,?,?,?,
		?,?,?,?,?,?,?,?,?,?,
		?,?,?,?,?,?,?,?,?,?
		)`
)

// MemberRepoPostgres data structure
type MemberRepoPostgres struct {
	*repository.Repository
}

// NewMemberRepoPostgres function for initializing member repo
func NewMemberRepoPostgres(repo *repository.Repository) *MemberRepoPostgres {
	return &MemberRepoPostgres{repo}
}

// BulkSave save multiple rows
func (mr *MemberRepoPostgres) BulkSave(ctxReq context.Context, members []model.Member) <-chan ResultRepository {
	ctx := "MemberRepo-BulkSave"
	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags["args"] = members

		tx, err := mr.WriteDB.Begin()
		if err != nil {
			output <- ResultRepository{Error: err}
			return
		}
		for _, member := range members {
			member.ID = helper.GenerateMemberIDv2()
			err := mr.exec(ctxReq, tx, member)
			if err != nil {
				tx.Rollback()
				output <- ResultRepository{Error: err}
				return
			}
		}
		tx.Commit()
		output <- ResultRepository{Error: nil}
	})

	return output
}

// Save function for saving member data
func (mr *MemberRepoPostgres) Save(ctxReq context.Context, member model.Member) <-chan ResultRepository {
	ctx := "MemberRepo-Save"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tx, err := mr.WriteDB.Begin()
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, member)
			output <- ResultRepository{Error: err}
			return
		}
		readStmt, err := tx.Prepare(`SELECT "version", email FROM member WHERE id = $1`)

		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, member)
			output <- ResultRepository{Error: err}
			return
		}
		defer readStmt.Close()

		var (
			version int
			email   string
		)
		err = readStmt.QueryRow(member.ID).Scan(&version, &email)
		if err != nil && err != sql.ErrNoRows {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, member)
			tx.Rollback()
			output <- ResultRepository{Error: err}
			return
		}

		// MVCC https://en.wikipedia.org/wiki/Multiversion_concurrency_control
		member.Version++

		// SET isStaff true if domain @bhinneka.com
		member.IsStaff = member.IsBhinnekaEmail()

		if version > member.Version {
			tx.Rollback()
			err := fmt.Errorf("there is conflict during save, unable to save model. You can discard the changes or just try again. Persistence Version=%d, Entity Version=%d, ID=%s",
				version, member.Version, member.ID)
			helper.SendErrorLog(ctxReq, ctx, helper.ScopeSaveMember, err, member)
			output <- ResultRepository{Error: err}
			return
		}

		if err != sql.ErrNoRows && email != member.Email {
			tx.Rollback()
			err := fmt.Errorf("there is conflict during save, unable to save model. You can discard the changes or just try again. Persistence Email=%s, Entity Email=%s, ID=%s",
				email, member.Email, member.ID)
			helper.SendErrorLog(ctxReq, ctx, helper.ScopeSaveMember, err, member)
			output <- ResultRepository{Error: err}
			return
		}
		tags[helper.TextEmail] = member.Email
		err = mr.insertMember(ctxReq, tx, member)
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, member)
			output <- ResultRepository{Error: err}
			return
		}
		// commit statement
		tx.Commit()

		output <- ResultRepository{Error: nil}
	})

	return output
}

func (mr *MemberRepoPostgres) exec(ctxReq context.Context, tx *sql.Tx, member model.Member) error {
	ctx := "MemberRepo-exec"
	readStmt, err := tx.Prepare(`SELECT "version", email FROM member WHERE id = $1`)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, member)
		return err
	}
	defer readStmt.Close()

	var (
		version int
		email   string
	)

	err = readStmt.QueryRow(member.ID).Scan(&version, &email)

	if err != nil && err != sql.ErrNoRows {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, member)
		return err
	}

	// MVCC https://en.wikipedia.org/wiki/Multiversion_concurrency_control
	member.Version++

	// SET isStaff true if domain @bhinneka.com
	member.IsStaff = member.IsBhinnekaEmail()

	if version > member.Version {
		err := fmt.Errorf("there is conflict during save, unable to save model. You can discard the changes or just try again. Persistence Version=%d, Entity Version=%d, ID=%s",
			version, member.Version, member.ID)
		helper.SendErrorLog(ctxReq, ctx, helper.ScopeSaveMember, err, member)
		return err
	}

	if err != sql.ErrNoRows && email != member.Email {
		err := fmt.Errorf("there is conflict during save, unable to save model. You can discard the changes or just try again. Persistence Email=%s, Entity Email=%s, ID=%s",
			email, member.Email, member.ID)
		helper.SendErrorLog(ctxReq, ctx, helper.ScopeSaveMember, err, member)
		return err
	}

	err = mr.insertMember(ctxReq, tx, member)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "insert_member", err, member)
		return err
	}
	return nil
}

func (mr *MemberRepoPostgres) insertMember(ctxReq context.Context, tx *sql.Tx, member model.Member) error {
	ctx := "MemberRepo-insertMember"

	queryInsert := `INSERT INTO member
		(
			id, "firstName", "lastName", email, gender, mobile, phone, ext,
			"birthDate", password, salt,
			province, "provinceId", city, "cityId", district, "districtId",
			"subDistrict", "subDistrictId", "zipCode", address,
			"jobTitle", department,
			"isAdmin", "isStaff", status,
			"facebookId", "googleId", "appleId", "azureId", "signUpFrom",
			created, "lastModified", version, token, "ldapId",  
			"lastPasswordModified", "facebookConnect", "googleConnect", "appleConnect", "isActive"
		)
	VALUES
		(
			$1, $2, $3, $4, $5, $6, $7, $8,
			$9, $10, $11,
			$12, $13, $14, $15, $16, $17,
			$18, $19, $20, $21,
			$22, $23,
			$24, $25, $26,
			$27, $28, $29, $30, $31, 
			$32, $33, $34, $35, $36, 
			$37, $38, $39, $40, $41
		)
	ON CONFLICT(id)
	DO UPDATE SET
		"firstName" = $2, "lastName" = $3, email = $4, gender = $5, mobile = $6, phone = $7, ext = $8,
		"birthDate" = $9, password = $10, salt = $11,
		province = $12, "provinceId" = $13, city = $14, "cityId" = $15, district = $16, "districtId" = $17,
		"subDistrict" = $18, "subDistrictId" = $19, "zipCode" = $20, address = $21,
		"jobTitle" = $22, department = $23,
		"isAdmin"=$24, "isStaff" = $25, status = $26, 
		"facebookId" = $27, "googleId" = $28, "appleId" = $29, "azureId" = $30, "signUpFrom" = $31,
		"lastModified" = $33, version = $34, token = $35, "ldapId"= $36, 
		"lastPasswordModified"= $37 , "facebookConnect"=$38, "googleConnect"=$39, "appleConnect"=$40, "isActive"=$41`

	stmt, err := tx.Prepare(queryInsert)
	tr := tracer.StartTrace(ctxReq, ctx)
	tags := make(map[string]interface{})
	defer func() {
		tr.Finish(tags)
	}()

	tags[helper.TextQuery] = queryInsert
	tags[helper.TextArgs] = member

	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "insert_member", err, member)
		return err
	}

	defer stmt.Close()

	// set null-able variables
	var (
		lastName, mobile, phone, ext, password, salt                                  sql.NullString
		province, city, district, subDistrict, zipCode, address                       sql.NullString
		department, jobTitle, token, gender                                           sql.NullString
		provinceID, cityID, districtID, subDistrictID                                 sql.NullString
		birthDate, lastPasswordModified, facebookConnect, googleConnect, appleConnect pq.NullTime
	)

	lastName = helper.ValidateStringToSQLNullString(member.LastName)
	mobile = helper.ValidateStringToSQLNullString(member.Mobile)
	phone = helper.ValidateStringToSQLNullString(member.Phone)
	ext = helper.ValidateStringToSQLNullString(member.Ext)
	password = helper.ValidateStringToSQLNullString(member.Password)
	salt = helper.ValidateStringToSQLNullString(member.Salt)
	province = helper.ValidateStringToSQLNullString(member.Address.Province)
	provinceID = helper.ValidateStringToSQLNullString(member.Address.ProvinceID)
	city = helper.ValidateStringToSQLNullString(member.Address.City)
	cityID = helper.ValidateStringToSQLNullString(member.Address.CityID)
	district = helper.ValidateStringToSQLNullString(member.Address.District)
	districtID = helper.ValidateStringToSQLNullString(member.Address.DistrictID)
	subDistrict = helper.ValidateStringToSQLNullString(member.Address.SubDistrict)
	subDistrictID = helper.ValidateStringToSQLNullString(member.Address.SubDistrictID)
	zipCode = helper.ValidateStringToSQLNullString(member.Address.ZipCode)
	address = helper.ValidateStringToSQLNullString(member.Address.Address)
	jobTitle = helper.ValidateStringToSQLNullString(member.JobTitle)
	department = helper.ValidateStringToSQLNullString(member.Department)
	token = helper.ValidateStringToSQLNullString(member.Token)

	birthDate.Valid = false
	if member.BirthDate.Year() > 0 {
		birthDate.Valid = true
		birthDate.Time = member.BirthDate
	}

	gender.Valid = false
	if len(member.Gender.String()) > 0 {
		gender.Valid = true
		gender.String = member.Gender.String()
	}

	lastPasswordModified.Valid = false
	if !member.LastPasswordModified.IsZero() {
		lastPasswordModified.Valid = true
		lastPasswordModified.Time = member.LastPasswordModified
	}

	facebookConnect.Valid = false
	if !member.SocialMedia.FacebookConnect.IsZero() {
		facebookConnect.Valid = true
		facebookConnect.Time = member.SocialMedia.FacebookConnect
	}

	googleConnect.Valid = false
	if !member.SocialMedia.GoogleConnect.IsZero() {
		googleConnect.Valid = true
		googleConnect.Time = member.SocialMedia.GoogleConnect
	}

	appleConnect.Valid = false
	if !member.SocialMedia.AppleConnect.IsZero() {
		appleConnect.Valid = true
		appleConnect.Time = member.SocialMedia.AppleConnect
	}

	_, err = stmt.Exec(
		member.ID, member.FirstName, lastName, member.Email, gender, mobile, phone, ext,
		birthDate, password, salt,
		province, provinceID, city, cityID, district, districtID,
		subDistrict, subDistrictID, zipCode, address,
		jobTitle, department,
		member.IsAdmin, member.IsStaff, member.Status.String(),
		member.SocialMedia.FacebookID, member.SocialMedia.GoogleID, member.SocialMedia.AppleID, member.SocialMedia.AzureID, member.SignUpFrom,
		time.Now(), time.Now(), member.Version, token, member.SocialMedia.LDAPID,
		lastPasswordModified, facebookConnect, googleConnect, appleConnect, member.IsActive,
	)

	if err != nil {
		return err
	}

	return nil
}

// Load function for loading member data based on user id
func (mr *MemberRepoPostgres) Load(ctxReq context.Context, uid string) <-chan ResultRepository {
	ctx := "MemberRepo-Load"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		q := `SELECT email, "firstName", "lastName",
				gender, mobile, phone, ext, "birthDate",
				password, salt, token,
				province, "provinceId", city, "cityId",
				district, "districtId", "subDistrict", "subDistrictId", "zipCode", address,
				"jobTitle", department,
				"isAdmin", "isStaff", status, "facebookId", "googleId", "appleId", "azureId", "ldapId", "signUpFrom",
				"lastLogin", "lastBlocked", "created", "lastModified", version, "profilePicture", "mfaEnabled", "mfaKey", 
				"lastPasswordModified", "facebookConnect", "googleConnect", "appleConnect", "isSync", "isActive" 
			FROM member WHERE id = $1`

		tags[helper.TextQuery] = q

		stmt, err := mr.ReadDB.Prepare(q)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, q)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
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
			birthDate, lastLogin, lastModified, lastBlocked                    pq.NullTime
			lastPasswordModified, facebookConnect, googleConnect, appleConnect pq.NullTime
			status                                                             string
		)

		err = stmt.QueryRow(uid).Scan(
			&member.Email, &member.FirstName, &lastName,
			&gender, &mobile, &phone, &ext, &birthDate,
			&password, &salt, &token,
			&province, &provinceID, &city, &cityID,
			&district, &districtID, &subDistrict, &subDistrictID, &zipCode, &address,
			&jobTitle, &department,
			&member.IsAdmin, &member.IsStaff, &status, &facebookID, &googleID, &appleID, &azureID, &ldapID, &signUpFrom,
			&lastLogin, &lastBlocked, &member.Created, &lastModified, &member.Version, &profilePicture,
			&member.MFAEnabled, &mfaKey, &lastPasswordModified, &facebookConnect, &googleConnect, &appleConnect, &member.IsSync, &member.IsActive,
		)

		if err != nil {
			if err == sql.ErrNoRows {
				err = fmt.Errorf(helper.ErrorDataNotFound, "user")
			}
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, q)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
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

		if gender.Valid {
			member.Gender = model.StringToGender(gender.String)
			member.GenderString = gender.String
		}

		if birthDate.Valid {
			member.BirthDate = birthDate.Time
			member.BirthDateString = birthDate.Time.Format("02/01/2006")
		}

		member = mr.adjustLastDateDataRepo(member, lastLogin, lastBlocked, lastModified, lastPasswordModified)

		member.Status = model.StringToStatus(status)
		member.ID = uid
		member.CreatedString = member.Created.Format(time.RFC3339)

		memberAddressLoad := model.Address{}
		memberAddressLoad.Province = helper.ValidateSQLNullString(province)
		memberAddressLoad.ProvinceID = helper.ValidateSQLNullString(provinceID)
		memberAddressLoad.City = helper.ValidateSQLNullString(city)
		memberAddressLoad.CityID = helper.ValidateSQLNullString(cityID)
		memberAddressLoad.District = helper.ValidateSQLNullString(district)
		memberAddressLoad.DistrictID = helper.ValidateSQLNullString(districtID)
		memberAddressLoad.SubDistrict = helper.ValidateSQLNullString(subDistrict)
		memberAddressLoad.SubDistrictID = helper.ValidateSQLNullString(subDistrictID)
		memberAddressLoad.ZipCode = helper.ValidateSQLNullString(zipCode)
		memberAddressLoad.Address = helper.ValidateSQLNullString(address)

		// parse get streets
		street1, street2 := helper.SplitStreetAddress(memberAddressLoad.Address)
		memberAddressLoad.Street1 = street1
		memberAddressLoad.Street2 = street2
		member.Address = memberAddressLoad

		ms := model.SocialMedia{
			FacebookID: facebookID.String,
			GoogleID:   googleID.String,
			AzureID:    azureID.String,
			LDAPID:     ldapID.String,
			AppleID:    appleID.String,
		}

		// adjustSocmedDataRepo restruct socmed data with null values
		member = mr.adjustSocmedDataRepo(member, ms, facebookConnect, googleConnect, appleConnect)

		tags[helper.TextResponse] = member
		output <- ResultRepository{Result: member}
	})

	return output
}

// UpdatePasswordMemberByEmail ...
func (mr *MemberRepoPostgres) UpdatePasswordMemberByEmail(ctxReq context.Context, member model.Member) <-chan ResultRepository {
	ctx := "MemberRepo-UpdatePasswordMemberByEmail"

	outputUpdate := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(outputUpdate)

		txUpdate, err := mr.WriteDB.Begin()
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, member)
			outputUpdate <- ResultRepository{Error: err}
			return
		}

		query := `UPDATE "member" SET "password"=$1, "salt"=$2, "lastPasswordModified"=$3 WHERE "email"=$4`
		stmtUpdate, err := txUpdate.Prepare(query)
		tags[helper.TextExecQuery] = query
		tags["password"] = member.Password
		tags["salt"] = member.Salt
		tags["lastPasswordModified"] = member.LastPasswordModified
		tags["email"] = member.Email

		if err != nil {
			txUpdate.Rollback()
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, member)
			outputUpdate <- ResultRepository{Error: err}
			return
		}
		defer stmtUpdate.Close()

		_, err = stmtUpdate.Exec(member.Password, member.Salt, member.LastPasswordModified, member.Email)
		if err != nil {
			txUpdate.Rollback()
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, member.Email)
			outputUpdate <- ResultRepository{Error: err}
			return
		}
		// commit statement
		txUpdate.Commit()

		outputUpdate <- ResultRepository{Error: nil}
	})

	return outputUpdate
}

func (mr *MemberRepoPostgres) adjustLastDateDataRepo(memberData model.Member, lastLogin, lastBlocked, lastModified, lastPasswordModified pq.NullTime) model.Member {
	if lastLogin.Valid {
		memberData.LastLogin = lastLogin.Time
		memberData.LastLoginString = memberData.LastLogin.Format(time.RFC3339)
	}

	if lastBlocked.Valid {
		memberData.LastBlocked = lastBlocked.Time
		memberData.LastBlockedString = memberData.LastBlocked.Format(time.RFC3339)
	}

	if lastModified.Valid {
		memberData.LastModified = lastModified.Time
		memberData.LastModifiedString = memberData.LastModified.Format(time.RFC3339)
	}

	if lastPasswordModified.Valid {
		memberData.LastPasswordModified = lastPasswordModified.Time
		memberData.LastPasswordModifiedString = memberData.LastPasswordModified.Format(time.RFC3339)
	}

	return memberData
}

// adjustSocmedDataRepo function for restruct socmed data with null values
func (mr *MemberRepoPostgres) adjustSocmedDataRepo(memberData model.Member, ms model.SocialMedia, facebookConnect, googleConnect, appleConnect pq.NullTime) model.Member {

	if facebookConnect.Valid {
		memberData.SocialMedia.FacebookConnect = facebookConnect.Time
		ms.FacebookConnect = memberData.SocialMedia.FacebookConnect
		memberData.SocialMedia.FacebookConnectString = memberData.SocialMedia.FacebookConnect.Format(time.RFC3339)
		ms.FacebookConnectString = memberData.SocialMedia.FacebookConnect.Format(time.RFC3339)
	}

	if googleConnect.Valid {
		memberData.SocialMedia.GoogleConnect = googleConnect.Time
		ms.GoogleConnect = memberData.SocialMedia.GoogleConnect
		memberData.SocialMedia.GoogleConnectString = memberData.SocialMedia.GoogleConnect.Format(time.RFC3339)
		ms.GoogleConnectString = memberData.SocialMedia.GoogleConnect.Format(time.RFC3339)
	}

	if appleConnect.Valid {
		memberData.SocialMedia.AppleConnect = appleConnect.Time
		ms.AppleConnect = memberData.SocialMedia.AppleConnect
		memberData.SocialMedia.AppleConnectString = memberData.SocialMedia.AppleConnect.Format(time.RFC3339)
		ms.AppleConnectString = memberData.SocialMedia.AppleConnect.Format(time.RFC3339)
	}

	memberData.SocialMedia = ms
	return memberData
}

// LoadMember function without goroutine
func (mr *MemberRepoPostgres) LoadMember(uid string) ResultRepository {
	ctx := "MemberRepo-LoadMember"
	tr := tracer.StartTrace(context.Background(), ctx)
	ctxReq := tr.NewChildContext()
	tags := make(map[string]interface{})
	defer func() {
		tr.Finish(tags)
	}()

	q := `SELECT email, "firstName", "lastName",
				gender, mobile, phone, ext, "birthDate",
				password, salt, token,
				province, "provinceId", city, "cityId",
				district, "districtId", "subDistrict", "subDistrictId", "zipCode", address,
				"jobTitle", department,
				"isAdmin", "isStaff", status, "facebookId", "googleId", "appleId", "azureId", "ldapId", "signUpFrom",
				"lastLogin", "lastBlocked", "created", "lastModified", version, "isSync"
			FROM member WHERE id = $1`

	stmt, err := mr.ReadDB.Prepare(q)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, uid)
		return ResultRepository{Error: err}
	}

	defer stmt.Close()

	// initialize needed variables
	var (
		member                                                      model.Member
		lastName, mobile, phone, ext, password, salt, gender, token sql.NullString
		province, city, district, subDistrict, zipCode, address     sql.NullString
		provinceID, cityID, districtID, subDistrictID               sql.NullString
		facebookID, googleID, appleID, azureID, ldapID, signUpFrom  sql.NullString
		jobTitle, department                                        sql.NullString
		birthDate, lastLogin, lastModified, lastBlocked             pq.NullTime
		status                                                      string
	)
	tags[helper.TextQuery] = q
	tags["args"] = member

	err = stmt.QueryRow(uid).Scan(
		&member.Email, &member.FirstName, &lastName,
		&gender, &mobile, &phone, &ext, &birthDate,
		&password, &salt, &token,
		&province, &provinceID, &city, &cityID,
		&district, &districtID, &subDistrict, &subDistrictID, &zipCode, &address,
		&jobTitle, &department,
		&member.IsAdmin, &member.IsStaff, &status, &facebookID, &googleID, &appleID, &azureID, &ldapID, &signUpFrom,
		&lastLogin, &lastBlocked, &member.Created, &lastModified, &member.Version, &member.IsSync,
	)

	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, uid)
		return ResultRepository{Error: err}
	}

	// assign the nullable field to object
	member.LastName = helper.ValidateSQLNullString(lastName)
	member.Mobile = helper.ValidateSQLNullString(mobile)
	member.Phone = helper.ValidateSQLNullString(phone)
	member.Ext = helper.ValidateSQLNullString(ext)
	member.Salt = helper.ValidateSQLNullString(salt)
	member.Token = helper.ValidateSQLNullString(token)
	member.JobTitle = helper.ValidateSQLNullString(jobTitle)
	member.Department = helper.ValidateSQLNullString(department)
	member.SignUpFrom = helper.ValidateSQLNullString(signUpFrom)

	gender.Valid = false
	if len(member.GenderString) > 0 {
		gender.Valid = true
		gender.String = member.GenderString
		member.Gender = model.StringToGender(gender.String)
	}
	if birthDate.Valid {
		member.BirthDate = birthDate.Time
		member.BirthDateString = birthDate.Time.Format("02/01/2006")
	}
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

	member.Status = model.StringToStatus(status)
	member.ID = uid
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

	street1, street2 := helper.SplitStreetAddress(ma.Address)
	ma.Street1 = street1
	ma.Street2 = street2
	member.Address = ma

	ms := model.SocialMedia{
		FacebookID: facebookID.String,
		GoogleID:   googleID.String,
		AppleID:    appleID.String,
		AzureID:    azureID.String,
		LDAPID:     ldapID.String,
	}

	member.SocialMedia = ms

	return ResultRepository{Result: member}
}

// FindMaxID function for getting member max id
func (mr *MemberRepoPostgres) FindMaxID(ctxReq context.Context) <-chan ResultRepository {
	ctx := "MemberRepoFindMaxID"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		q := `SELECT COALESCE(MAX(CAST(SUBSTRING(id, 8, 15) AS INTEGER)), 0) as max 
			FROM "member" WHERE id LIKE CONCAT('USR', TO_CHAR(NOW(), 'YYmm'), '%')`
		tags[helper.TextQuery] = q

		stmt, err := mr.ReadDB.Prepare(q)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, q)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		defer stmt.Close()

		var maxID int64
		err = stmt.QueryRow().Scan(&maxID)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, q)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Result: maxID}
	})

	return output
}

// UpdateProfilePicture function for update profile picture member data
func (mr *MemberRepoPostgres) UpdateProfilePicture(ctxReq context.Context, data model.ProfilePicture) <-chan ResultRepository {
	ctx := "MemberRepo-UpdateProfilePicture"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		query := `UPDATE "member" SET "profilePicture" = $2 WHERE "id" = $1  RETURNING "id", "profilePicture"`

		tags[helper.TextQuery] = query
		stmt, err := mr.WriteDB.Prepare(query)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}
		var dataProfilePicture model.ProfilePicture
		err = stmt.QueryRow(data.ID, data.ProfilePicture).Scan(&dataProfilePicture.ID, &dataProfilePicture.ProfilePicture)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, data)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}
		tags[helper.TextResponse] = dataProfilePicture

		output <- ResultRepository{Result: dataProfilePicture}
	})
	return output
}

// UpdateFlagIsSyncMember ...
func (mr *MemberRepoPostgres) UpdateFlagIsSyncMember(ctxReq context.Context, member model.Member) <-chan ResultRepository {
	ctx := "MemberRepo-UpdateFlagIsSyncMember"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tx, err := mr.WriteDB.Begin()
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, member)
			output <- ResultRepository{Error: err}
			return
		}
		stmt, err := tx.Prepare(`UPDATE "member" SET "isSync" = true WHERE "id" = $1`)

		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, member)
			output <- ResultRepository{Error: err}
			return
		}
		defer stmt.Close()

		_, err = stmt.Exec(member.ID)
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, member.ID)
			output <- ResultRepository{Error: err}
			return
		}
		// commit statement
		tx.Commit()

		output <- ResultRepository{Error: nil}
	})

	return output
}
