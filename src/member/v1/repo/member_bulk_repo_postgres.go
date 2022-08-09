package repo

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

// BulkImportSave save multiple rows in a single batch
func (mr *MemberRepoPostgres) BulkImportSave(ctxReq context.Context, members []*model.Member) <-chan ResultRepository {
	ctx := "MemberRepo-BulkImportSave"
	output := make(chan ResultRepository)

	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		tags["args"] = members

		tx, err := mr.WriteDB.Begin()
		if err != nil {
			output <- ResultRepository{Error: err}
			return
		}

		args := []interface{}{}

		queryInsert := `
		INSERT INTO member (id, "firstName", "lastName", email, gender, mobile, phone, ext,"birthDate", password, salt,
			province, "provinceId", city, "cityId", district, "districtId","subDistrict", "subDistrictId", "zipCode", address,
			"jobTitle", department, "isAdmin", "isStaff", status,"facebookId", "googleId", "appleId", "azureId", "signUpFrom", created, 
			"lastModified", version, token, "ldapId", "lastPasswordModified", "facebookConnect", "googleConnect", "appleConnect")
			VALUES %s`

		args = append(args, mr.generateArgs(members)...)

		queryInsert = mr.bulkReplaceSQLString(queryInsert, queryPattern, len(members))

		stmt, err := tx.Prepare(queryInsert)
		if err != nil {
			tx.Rollback()
			helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, members)
			output <- ResultRepository{Error: err}
			return
		}
		defer stmt.Close()

		if _, err = stmt.Exec(args...); err != nil {
			tx.Rollback()
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, members)
			output <- ResultRepository{Error: err}
			return
		}

		if err = tx.Commit(); err != nil {
			tx.Rollback()
			output <- ResultRepository{Error: err}
			return
		}
		output <- ResultRepository{Error: nil, Result: members}
	})

	return output
}

func (mr *MemberRepoPostgres) generateArgs(members []*model.Member) []interface{} {
	args := []interface{}{}
	var (
		lastName, mobile, phone, ext, password, salt                                  sql.NullString
		province, city, district, subDistrict, zipCode, address                       sql.NullString
		department, jobTitle, token, gender                                           sql.NullString
		provinceID, cityID, districtID, subDistrictID                                 sql.NullString
		birthDate, lastPasswordModified, facebookConnect, googleConnect, appleConnect pq.NullTime
	)

	for i, memberArg := range members {
		memberID := helper.GenerateMemberIDv2()
		members[i].ID = memberID
		memberArg.ID = memberID

		lastName = helper.ValidateStringToSQLNullString(memberArg.LastName)
		mobile = helper.ValidateStringToSQLNullString(memberArg.Mobile)
		phone = helper.ValidateStringToSQLNullString(memberArg.Phone)
		ext = helper.ValidateStringToSQLNullString(memberArg.Ext)
		password = helper.ValidateStringToSQLNullString(memberArg.Password)
		salt = helper.ValidateStringToSQLNullString(memberArg.Salt)
		province = helper.ValidateStringToSQLNullString(memberArg.Address.Province)
		provinceID = helper.ValidateStringToSQLNullString(memberArg.Address.ProvinceID)
		city = helper.ValidateStringToSQLNullString(memberArg.Address.City)
		cityID = helper.ValidateStringToSQLNullString(memberArg.Address.CityID)
		district = helper.ValidateStringToSQLNullString(memberArg.Address.District)
		districtID = helper.ValidateStringToSQLNullString(memberArg.Address.DistrictID)
		subDistrict = helper.ValidateStringToSQLNullString(memberArg.Address.SubDistrict)
		subDistrictID = helper.ValidateStringToSQLNullString(memberArg.Address.SubDistrictID)
		zipCode = helper.ValidateStringToSQLNullString(memberArg.Address.ZipCode)
		address = helper.ValidateStringToSQLNullString(memberArg.Address.Address)
		jobTitle = helper.ValidateStringToSQLNullString(memberArg.JobTitle)
		department = helper.ValidateStringToSQLNullString(memberArg.Department)
		token = helper.ValidateStringToSQLNullString(memberArg.Token)

		birthDate.Valid = false
		if memberArg.BirthDate.Year() > 0 {
			birthDate.Valid = true
			birthDate.Time = memberArg.BirthDate
		}

		gender.Valid = false
		if len(memberArg.Gender.String()) > 0 {
			gender.Valid = true
			gender.String = memberArg.Gender.String()
		}

		lastPasswordModified.Valid = false
		if !memberArg.LastPasswordModified.IsZero() {
			lastPasswordModified.Valid = true
			lastPasswordModified.Time = memberArg.LastPasswordModified
		}

		facebookConnect.Valid = false
		if !memberArg.SocialMedia.FacebookConnect.IsZero() {
			facebookConnect.Valid = true
			facebookConnect.Time = memberArg.SocialMedia.FacebookConnect
		}

		googleConnect.Valid = false
		if !memberArg.SocialMedia.GoogleConnect.IsZero() {
			googleConnect.Valid = true
			googleConnect.Time = memberArg.SocialMedia.GoogleConnect
		}

		appleConnect.Valid = false
		if !memberArg.SocialMedia.AppleConnect.IsZero() {
			appleConnect.Valid = true
			appleConnect.Time = memberArg.SocialMedia.AppleConnect
		}
		args = append(args, memberArg.ID, memberArg.FirstName, lastName, memberArg.Email, gender, mobile, phone, ext,
			birthDate, password, salt,
			province, provinceID, city, cityID, district, districtID,
			subDistrict, subDistrictID, zipCode, address,
			jobTitle, department,
			memberArg.IsAdmin, memberArg.IsStaff, memberArg.Status.String(),
			memberArg.SocialMedia.FacebookID, memberArg.SocialMedia.GoogleID, memberArg.SocialMedia.AppleID, memberArg.SocialMedia.AzureID, memberArg.SignUpFrom,
			time.Now(), time.Now(), memberArg.Version, token, memberArg.SocialMedia.LDAPID,
			lastPasswordModified, facebookConnect, googleConnect, appleConnect)
	}
	return args
}

// BulkReplaceSQLString replace string
func (mr *MemberRepoPostgres) bulkReplaceSQLString(queryString, pattern string, len int) string {
	pattern += ","
	queryString = fmt.Sprintf(queryString, strings.Repeat(pattern, len))
	n := 0
	for strings.IndexByte(queryString, '?') != -1 {
		n++
		param := "$" + strconv.Itoa(n)
		queryString = strings.Replace(queryString, "?", param, 1)
	}
	return strings.TrimSuffix(queryString, ",")
}
