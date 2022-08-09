package repo

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
	"github.com/Bhinneka/user-service/src/shared/repository"
)

const private = "private"

// MerchantRepoPostgres data structure
type MerchantRepoPostgres struct {
	*repository.Repository
}

// NewMerchantRepoPostgres function for initializing  repo
func NewMerchantRepoPostgres(repo *repository.Repository) MerchantRepository {
	return &MerchantRepoPostgres{repo}
}

// AddUpdateMerchant function for saving  data
func (mr *MerchantRepoPostgres) AddUpdateMerchant(ctxReq context.Context, merchant model.B2CMerchantDataV2) <-chan ResultRepository {
	ctx := "MerchantRepo-AddUpdateMerchant"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		var (
			stmt                                                               *sql.Stmt
			id                                                                 string
			err                                                                error
			merchantType, genderPic, merchantGroup, upgradeStatus, productType sql.NullString
		)

		mr.ReadDB.QueryRow(`SELECT id FROM b2c_merchant WHERE id=$1`, merchant.ID).Scan(&id)
		query := `INSERT INTO b2c_merchant
					(
						"id", "userId", "merchantName", "vanityURL", "merchantCategory", "companyName", 
						"pic", "picOccupation", "dailyOperationalStaff", "storeClosureDate", "storeReopenDate", 
						"storeActiveShippingDate", "merchantAddress", "merchantVillage", "merchantDistrict", 
						"merchantCity", "merchantProvince", "zipCode", "storeAddress", "storeVillage", "storeDistrict", 
						"storeCity", "storeProvince", "storeZipCode", "phoneNumber", "mobilePhoneNumber",
						"additionalEmail", "merchantDescription", "merchantLogo", "accountHolderName", "bankName", 
						"accountNumber", "isPKP", "npwp", "npwpHolderName", "richContent", "notificationPreferences", 
						"merchantRank", "acquisitor", "accountManager", "launchDev", "skuLive", "mouDate", "note", 
						"agreementDate", "isActive", "creatorId", "creatorIp", "editorId", "editorIp", "version", 
						"created", "lastModified", "merchantVillageId", "merchantDistrictId", "merchantCityId", 
						"merchantProvinceId", "storeVillageId", "storeDistrictId", "storeCityId", "storeProvinceId", 
						"bankId", "isClosed", "merchantEmail", "bankBranch", "picKtpFile", "npwpFile", "businessType", "source",
						"merchantType", "genderPic", "merchantGroup", "upgradeStatus", "productType", "legalEntity", "numberOfEmployee", "status", "countUpdateNameAvailable", "sellerOfficerName", "sellerOfficerEmail"
					)
				VALUES
					(
						$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 
						$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
						$21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
						$31, $32, $33, $34, $35, $36, $37, $38, $39, $40,
						$41, $42, $43, $44, $45, $46, $47, $48, $49, $50,
						$51, $52, $53, $54, $55, $56, $57, $58, $59, $60,
						$61, $62, $63, $64, $65, $66, $67, $68, $69,
						$70, $71, $72, $73, $74, $75, $76, $77, $78,$79, $80
					)
				ON CONFLICT(id)
				DO UPDATE SET
				"id"=$1, "userId"=$2, "merchantName"=$3, "vanityURL"=$4, "merchantCategory"=$5, "companyName"=$6, 
				"pic"=$7, "picOccupation"=$8, "dailyOperationalStaff"=$9, "storeClosureDate"=$10, "storeReopenDate"=$11, 
				"storeActiveShippingDate"=$12, "merchantAddress"=$13, "merchantVillage"=$14, "merchantDistrict"=$15, 
				"merchantCity"=$16, "merchantProvince"=$17, "zipCode"=$18, "storeAddress"=$19, "storeVillage"=$20, "storeDistrict"=$21, 
				"storeCity"=$22, "storeProvince"=$23, "storeZipCode"=$24, "phoneNumber"=$25, "mobilePhoneNumber"=$26,
				"additionalEmail"=$27, "merchantDescription"=$28, "merchantLogo"=$29, "accountHolderName"=$30, "bankName"=$31, 
				"accountNumber"=$32, "isPKP"=$33, "npwp"=$34, "npwpHolderName"=$35, "richContent"=$36, "notificationPreferences"=$37, 
				"merchantRank"=$38, "acquisitor"=$39, "accountManager"=$40, "launchDev"=$41, "skuLive"=$42, "mouDate"=$43, "note"=$44, 
				"agreementDate"=$45, "isActive"=$46, "creatorId"=$47, "creatorIp"=$48, "editorId"=$49, "editorIp"=$50, "version"=$51, 
				"created"=$52, "lastModified"=$53, "merchantVillageId"=$54, "merchantDistrictId"=$55, "merchantCityId"=$56, 
				"merchantProvinceId"=$57, "storeVillageId"=$58, "storeDistrictId"=$59, "storeCityId"=$60, "storeProvinceId"=$61, 
				"bankId"=$62, "isClosed"=$63, "merchantEmail"=$64, "bankBranch"=$65, "picKtpFile"=$66, "npwpFile"=$67, "businessType"=$68, "source"=$69,
				"merchantType"=$70, "genderPic"=$71, "merchantGroup"=$72, "upgradeStatus"=$73, "productType"=$74, "legalEntity"=$75, "numberOfEmployee"=$76, "status"=$77, "countUpdateNameAvailable"=$78, "sellerOfficerName"=$79, "sellerOfficerEmail"=$80`

		tags[helper.TextQuery] = query
		if mr.Tx != nil {
			stmt, err = mr.Tx.Prepare(query)
		} else {
			stmt, err = mr.WriteDB.Prepare(query)
		}

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}
		defer stmt.Close()

		mr.setMerchantData(merchant, &merchantType, &genderPic, &merchantGroup, &upgradeStatus, &productType)

		_, err = stmt.Exec(
			merchant.ID, merchant.UserID, merchant.MerchantName, merchant.VanityURL,
			merchant.MerchantCategory, merchant.CompanyName, merchant.Pic, merchant.PicOccupation,
			merchant.DailyOperationalStaff, merchant.StoreClosureDate, merchant.StoreReopenDate, merchant.StoreActiveShippingDate,
			merchant.MerchantAddress, merchant.MerchantVillage, merchant.MerchantDistrict,
			merchant.MerchantCity, merchant.MerchantProvince, merchant.ZipCode, merchant.StoreAddress,
			merchant.StoreVillage, merchant.StoreDistrict, merchant.StoreCity, merchant.StoreProvince,
			merchant.StoreZipCode, merchant.PhoneNumber, merchant.MobilePhoneNumber,
			merchant.AdditionalEmail, merchant.MerchantDescription, merchant.MerchantLogo,
			merchant.AccountHolderName, merchant.BankName, merchant.AccountNumber, merchant.IsPKP,
			merchant.Npwp, merchant.NpwpHolderName, merchant.RichContent, merchant.NotificationPreferences,
			merchant.MerchantRank, merchant.Acquisitor, merchant.AccountManager,
			merchant.LaunchDev, merchant.SkuLive, merchant.MouDate, merchant.Note,
			merchant.AgreementDate, merchant.IsActive, merchant.CreatorID,
			merchant.CreatorIP, merchant.EditorID, merchant.EditorIP, merchant.Version,
			merchant.Created, merchant.LastModified,
			merchant.MerchantVillageID, merchant.MerchantDistrictID, merchant.MerchantCityID,
			merchant.MerchantProvinceID, merchant.StoreVillageID,
			merchant.StoreDistrictID, merchant.StoreCityID, merchant.StoreProvinceID, merchant.BankID,
			merchant.IsClosed, merchant.MerchantEmail,
			merchant.BankBranch, merchant.PicKtpFile, merchant.NpwpFile, merchant.BusinessType, merchant.Source,
			merchantType, genderPic, merchantGroup, upgradeStatus, productType, merchant.LegalEntity, merchant.NumberOfEmployee, merchant.Status, merchant.CountUpdateNameAvailable, merchant.SellerOfficerName, merchant.SellerOfficerEmail,
		)

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		tags["args"] = merchant

		output <- ResultRepository{Error: nil}
	})
	return output
}

func (mr *MerchantRepoPostgres) setMerchantData(merchant model.B2CMerchantDataV2, merchantType, genderPic, merchantGroup, upgradeStatus, productType *sql.NullString) {
	if merchant.MerchantTypeString.Valid {
		merchantType.Valid = true
		merchantType.String = merchant.MerchantTypeString.String
	}

	if merchant.GenderPic.String() != "" {
		genderPic.Valid = true
		genderPic.String = merchant.GenderPic.String()
	}

	if merchant.MerchantGroup.String != "" {
		merchantGroup.Valid = true
		merchantGroup.String = merchant.MerchantGroup.String
	}

	if merchant.UpgradeStatus.String != "" {
		upgradeStatus.Valid = true
		upgradeStatus.String = merchant.UpgradeStatus.String
	}
	if merchant.ProductType.String != "" {
		productType.Valid = true
		productType.String = merchant.ProductType.String
	}
}

// LoadMerchant function for getting detail merchant by id
func (mr *MerchantRepoPostgres) LoadMerchant(ctxReq context.Context, uid string, privacy string) ResultRepository {
	ctx := "MerchantRepo-LoadMerchant"

	var (
		filter      string
		queryValues []interface{}
	)

	queryValues = append(queryValues, uid)
	filter = `WHERE "deletedAt" IS NULL AND b2c_merchant."id" = $1`
	merchant, err := mr.findMerchant(ctxReq, ctx, filter, queryValues, privacy)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, merchant)
		return ResultRepository{Error: err}
	}

	return ResultRepository{Result: merchant}
}

// LoadMerchant function for getting detail merchant by id
func (mr *MerchantRepoPostgres) LoadMerchantByVanityURL(ctxReq context.Context, uid string) ResultRepository {
	ctx := "MerchantRepo-LoadMerchantByVanityURL"

	var (
		filter      string
		queryValues []interface{}
		public      string
	)

	queryValues = append(queryValues, uid)
	filter = `WHERE "deletedAt" IS NULL AND "vanityURL" = $1`
	public = "public"
	merchantData, err := mr.findMerchant(ctxReq, ctx, filter, queryValues, public)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, merchantData)
		return ResultRepository{Error: err}
	}

	return ResultRepository{Result: merchantData}
}

// FindMerchantByEmail function for getting detail merchant by id
func (mr *MerchantRepoPostgres) FindMerchantByEmail(ctxReq context.Context, email string) ResultRepository {
	ctx := "MerchantRepo-FindMerchantByEmail"

	var (
		filter      string
		queryValues []interface{}
		private     string
	)

	email = strings.ToLower(email)
	queryValues = append(queryValues, email)
	filter = `WHERE "deletedAt" IS NULL AND LOWER("merchantEmail") = $1`
	merchant, err := mr.findMerchant(ctxReq, ctx, filter, queryValues, private)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, email)
		return ResultRepository{Error: err}
	}

	return ResultRepository{Result: merchant}
}

// FindMerchantByUser function for getting detail merchant by user id
func (mr *MerchantRepoPostgres) FindMerchantByUser(ctxReq context.Context, uid string) ResultRepository {
	ctx := "MerchantRepo-FindMerchantByUser"

	var (
		filter      string
		queryValues []interface{}
		private     string
	)

	queryValues = append(queryValues, uid)
	filter = `WHERE "deletedAt" IS NULL AND "userId" = $1`
	merchant, err := mr.findMerchant(ctxReq, ctx, filter, queryValues, private)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, merchant)
		return ResultRepository{Error: err}
	}

	return ResultRepository{Result: merchant}
}

// FindMerchantByName function for getting detail merchant by name
func (mr *MerchantRepoPostgres) FindMerchantByName(ctxReq context.Context, name string) ResultRepository {
	ctx := "MerchantRepo-FindMerchantByName"

	var (
		filter      string
		queryValues []interface{}
		private     string
	)
	name = strings.ToLower(name)
	queryValues = append(queryValues, name)
	filter = `WHERE "deletedAt" IS NULL AND LOWER("merchantName") = $1`
	merchant, err := mr.findMerchant(ctxReq, ctx, filter, queryValues, private)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, name)
		return ResultRepository{Error: err}
	}

	return ResultRepository{Result: merchant}
}

// FindMerchantBySlug function for getting detail merchant by slug
func (mr *MerchantRepoPostgres) FindMerchantBySlug(ctxReq context.Context, slug string) ResultRepository {
	ctx := "MerchantRepo-FindMerchantBySlug"

	var (
		filter      string
		queryValues []interface{}
		private     string
	)

	slug = strings.ToLower(slug)
	queryValues = append(queryValues, slug)
	filter = `WHERE "deletedAt" IS NULL AND LOWER("vanityURL") = $1`

	merchant, err := mr.findMerchant(ctxReq, ctx, filter, queryValues, private)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, slug)
		return ResultRepository{Error: err}
	}

	return ResultRepository{Result: merchant}
}

// FindMerchantByID function for getting detail merchant by id and merchantid
func (mr *MerchantRepoPostgres) FindMerchantByID(ctxReq context.Context, id, userID string) ResultRepository {
	ctx := "MerchantRepo-FindMerchantByID"

	var (
		filter      string
		queryValues []interface{}
		private     string
	)

	queryValues = append(queryValues, id)
	queryValues = append(queryValues, userID)
	filter = `WHERE "deletedAt" IS NULL AND b2c_merchant."id" = $1 AND "userId" = $2`
	merchant, err := mr.findMerchant(ctxReq, ctx, filter, queryValues, private)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, merchant)
		return ResultRepository{Error: err}
	}

	return ResultRepository{Result: merchant}
}

func (mr *MerchantRepoPostgres) getSelect() string {
	return `
	SELECT b2c_merchant."id", "userId", "merchantName", "vanityURL", "merchantCategory", "companyName", 
	"pic", "picOccupation", "dailyOperationalStaff", "storeClosureDate", "storeReopenDate", 
	"storeActiveShippingDate", "merchantAddress", "merchantVillage", "merchantDistrict", 
	"merchantCity", "merchantProvince", "zipCode", "storeAddress", "storeVillage", "storeDistrict", 
	"storeCity", "storeProvince", "storeZipCode", "phoneNumber", "mobilePhoneNumber",
	"additionalEmail", "merchantDescription", "merchantLogo", "accountHolderName", "bankName", 
	"accountNumber", "isPKP", "npwp", "npwpHolderName", "richContent", "notificationPreferences", 
	"merchantRank", "acquisitor", "accountManager", "launchDev", "skuLive", "mouDate", "note", 
	"agreementDate", "isActive","status", "creatorId", "creatorIp", "editorId", "editorIp", "version", 
	"created", "lastModified", "merchantVillageId", "merchantDistrictId", "merchantCityId", 
	"merchantProvinceId", "storeVillageId", "storeDistrictId", "storeCityId", "storeProvinceId", 
	"bankId", "isClosed", "merchantEmail", "bankBranch", "picKtpFile", 
	"npwpFile", "businessType", "source", "merchantType", "genderPic", "merchantGroup", "upgradeStatus", "productType",
	"legalEntity", legal_entities."name" as "legalEntityName",
	"numberOfEmployee", number_of_employees."name" as "numberOfEmployeeName","countUpdateNameAvailable","sellerOfficerName","sellerOfficerEmail","reason"`

}

func (mr *MerchantRepoPostgres) getSelectPublic() string {
	return `
	SELECT b2c_merchant."id", "userId", "merchantName", "vanityURL", "merchantCategory", "companyName", 
	"pic", "picOccupation", "dailyOperationalStaff", "storeClosureDate", "storeReopenDate", 
	"storeActiveShippingDate", "merchantAddress", "merchantVillage", "merchantDistrict", 
	"merchantCity", "merchantProvince", "zipCode", "storeAddress", "storeVillage", "storeDistrict", 
	"storeCity", "storeProvince", "storeZipCode",
	"merchantDescription", "merchantLogo",
	"isPKP", "richContent", "notificationPreferences", 
	"merchantRank", "acquisitor", "accountManager", "launchDev", "skuLive", "mouDate", "note", 
	"merchantVillageId", "merchantDistrictId", "merchantCityId", 
	"merchantProvinceId", "storeVillageId", "storeDistrictId", "storeCityId", "storeProvinceId", 
	"isClosed",
	"businessType", "merchantType", "genderPic", "merchantGroup", "productType","isActive","status"`
}

func (mr *MerchantRepoPostgres) additionalSelect() string {
	return `
	LEFT JOIN legal_entities  ON (b2c_merchant."legalEntity"=legal_entities."id")
	LEFT JOIN number_of_employees  ON (b2c_merchant."numberOfEmployee"=number_of_employees."id")
	`
}
func (mr *MerchantRepoPostgres) findMerchant(ctxReq context.Context, ctx, filter string, queryValues []interface{}, privacy string) (model.B2CMerchantDataV2, error) {
	var (
		merchant model.B2CMerchantDataV2
	)
	query := fmt.Sprintf(`%s FROM b2c_merchant %s %s`, mr.getSelect(), mr.additionalSelect(), filter)
	if privacy == "public" {
		query = fmt.Sprintf(`%s FROM b2c_merchant %s`, mr.getSelectPublic(), filter)
	}
	tr := tracer.StartTrace(ctxReq, ctx)
	tags := make(map[string]interface{})
	defer func() {
		tr.Finish(tags)
	}()

	tags[helper.TextQuery] = query
	stmt, err := mr.ReadDB.Prepare(query)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
		tags[helper.TextResponse] = err
		return merchant, err
	}

	defer stmt.Close()

	if privacy == "public" {
		err = stmt.QueryRow(queryValues...).Scan(
			&merchant.ID, &merchant.UserID, &merchant.MerchantName, &merchant.VanityURL, &merchant.MerchantCategory, &merchant.CompanyName,
			&merchant.Pic, &merchant.PicOccupation, &merchant.DailyOperationalStaff, &merchant.StoreClosureDate, &merchant.StoreReopenDate,
			&merchant.StoreActiveShippingDate, &merchant.MerchantAddress, &merchant.MerchantVillage, &merchant.MerchantDistrict,
			&merchant.MerchantCity, &merchant.MerchantProvince, &merchant.ZipCode, &merchant.StoreAddress, &merchant.StoreVillage, &merchant.StoreDistrict,
			&merchant.StoreCity, &merchant.StoreProvince, &merchant.StoreZipCode,
			&merchant.MerchantDescription, &merchant.MerchantLogo,
			&merchant.IsPKP, &merchant.RichContent, &merchant.NotificationPreferences,
			&merchant.MerchantRank, &merchant.Acquisitor, &merchant.AccountManager, &merchant.LaunchDev, &merchant.SkuLive, &merchant.MouDate, &merchant.Note,
			&merchant.MerchantVillageID, &merchant.MerchantDistrictID, &merchant.MerchantCityID,
			&merchant.MerchantProvinceID, &merchant.StoreVillageID, &merchant.StoreDistrictID, &merchant.StoreCityID, &merchant.StoreProvinceID,
			&merchant.IsClosed,
			&merchant.BusinessType, &merchant.MerchantTypeString, &merchant.GenderPicString,
			&merchant.MerchantGroup, &merchant.ProductType, &merchant.IsActive, &merchant.Status,
		)
	} else {
		err = stmt.QueryRow(queryValues...).Scan(
			&merchant.ID, &merchant.UserID, &merchant.MerchantName, &merchant.VanityURL, &merchant.MerchantCategory, &merchant.CompanyName,
			&merchant.Pic, &merchant.PicOccupation, &merchant.DailyOperationalStaff, &merchant.StoreClosureDate, &merchant.StoreReopenDate,
			&merchant.StoreActiveShippingDate, &merchant.MerchantAddress, &merchant.MerchantVillage, &merchant.MerchantDistrict,
			&merchant.MerchantCity, &merchant.MerchantProvince, &merchant.ZipCode, &merchant.StoreAddress, &merchant.StoreVillage, &merchant.StoreDistrict,
			&merchant.StoreCity, &merchant.StoreProvince, &merchant.StoreZipCode, &merchant.PhoneNumber, &merchant.MobilePhoneNumber,
			&merchant.AdditionalEmail, &merchant.MerchantDescription, &merchant.MerchantLogo, &merchant.AccountHolderName, &merchant.BankName,
			&merchant.AccountNumber, &merchant.IsPKP, &merchant.Npwp, &merchant.NpwpHolderName, &merchant.RichContent, &merchant.NotificationPreferences,
			&merchant.MerchantRank, &merchant.Acquisitor, &merchant.AccountManager, &merchant.LaunchDev, &merchant.SkuLive, &merchant.MouDate, &merchant.Note,
			&merchant.AgreementDate, &merchant.IsActive, &merchant.Status, &merchant.CreatorID, &merchant.CreatorIP, &merchant.EditorID, &merchant.EditorIP, &merchant.Version,
			&merchant.Created, &merchant.LastModified, &merchant.MerchantVillageID, &merchant.MerchantDistrictID, &merchant.MerchantCityID,
			&merchant.MerchantProvinceID, &merchant.StoreVillageID, &merchant.StoreDistrictID, &merchant.StoreCityID, &merchant.StoreProvinceID,
			&merchant.BankID, &merchant.IsClosed, &merchant.MerchantEmail, &merchant.BankBranch, &merchant.PicKtpFile,
			&merchant.NpwpFile, &merchant.BusinessType, &merchant.Source, &merchant.MerchantTypeString, &merchant.GenderPicString,
			&merchant.MerchantGroup, &merchant.UpgradeStatus, &merchant.ProductType,
			&merchant.LegalEntity, &merchant.LegalEntityName, &merchant.NumberOfEmployee, &merchant.NumberOfEmployeeName, &merchant.CountUpdateNameAvailable, &merchant.SellerOfficerName, &merchant.SellerOfficerEmail, &merchant.Reason,
		)
	}

	if _, err := json.Marshal(&merchant); err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, merchant)
		return merchant, err
	}

	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
		tags[helper.TextResponse] = err
		return merchant, err
	}

	tags[helper.TextResponse] = merchant
	return merchant, nil
}

// SoftDelete function for flagging delete
func (mr *MerchantRepoPostgres) SoftDelete(ctxReq context.Context, merchantID string) <-chan ResultRepository {
	ctx := "MerchantRepoPostgres-SoftDelete"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		queryDelete := `UPDATE b2c_merchant SET "deletedAt"=$1, "status"=$2 WHERE id=$3;`

		tags[helper.TextQuery] = queryDelete
		tags[helper.TextMerchantIDCamel] = merchantID

		stmt, err := mr.WriteDB.Prepare(queryDelete)
		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, merchantID)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		if _, err = stmt.Exec(time.Now(), model.DeletedString, merchantID); err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, merchantID)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}
	})

	return output
}

// GetTotalMerchant return total merchants
func (mr *MerchantRepoPostgres) GetTotalMerchant(ctxReq context.Context, params *model.QueryParameters) <-chan ResultRepository {
	ctx := "MerchantRepoPostgres-GetTotalMerchant"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(_ context.Context, tags map[string]interface{}) {
		defer close(output)
		var totalData int
		query, bindVar := mr.makeQuery(params)
		fullQuery := `SELECT count(id) FROM b2c_merchant ` + query

		tags[helper.TextQuery] = fullQuery

		if err := mr.ReadDB.QueryRow(fullQuery, bindVar...).Scan(&totalData); err != nil {
			output <- ResultRepository{Error: err}
			return
		}
		output <- ResultRepository{Result: totalData}
	})
	return output
}

// GetMerchants retrieve merchant data based on given parameters
func (mr *MerchantRepoPostgres) GetMerchants(ctxReq context.Context, params *model.QueryParameters) <-chan ResultRepository {
	ctx := "MerchantRepoPostgres-GetMerchants"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		orderBy := `"created"`
		sort := "DESC"

		query, bindVar := mr.makeQuery(params)
		if params.OrderBy != "" && params.SortBy != "" {
			orderBy = strconv.Quote(params.OrderBy)
			sort = params.SortBy
		}

		fullQuery := mr.getSelect() + ` FROM b2c_merchant ` + mr.additionalSelect() + query
		fullQuery += fmt.Sprintf(` ORDER BY %s %s`, orderBy, sort)
		fullQuery += fmt.Sprintf(` LIMIT %d OFFSET %d`, params.Limit, params.Offset)

		tags[helper.TextQuery] = fullQuery

		rows, err := mr.ReadDB.Query(fullQuery, bindVar...)
		if err != nil {
			output <- ResultRepository{Error: err}
			return
		}
		results := []model.B2CMerchantDataV2{}
		for rows.Next() {
			var (
				row model.B2CMerchantDataV2
			)
			err := rows.Scan(
				&row.ID, &row.UserID, &row.MerchantName, &row.VanityURL, &row.MerchantCategory, &row.CompanyName,
				&row.Pic, &row.PicOccupation, &row.DailyOperationalStaff, &row.StoreClosureDate, &row.StoreReopenDate,
				&row.StoreActiveShippingDate, &row.MerchantAddress, &row.MerchantVillage, &row.MerchantDistrict,
				&row.MerchantCity, &row.MerchantProvince, &row.ZipCode, &row.StoreAddress, &row.StoreVillage, &row.StoreDistrict,
				&row.StoreCity, &row.StoreProvince, &row.StoreZipCode, &row.PhoneNumber, &row.MobilePhoneNumber,
				&row.AdditionalEmail, &row.MerchantDescription, &row.MerchantLogo, &row.AccountHolderName, &row.BankName,
				&row.AccountNumber, &row.IsPKP, &row.Npwp, &row.NpwpHolderName, &row.RichContent, &row.NotificationPreferences,
				&row.MerchantRank, &row.Acquisitor, &row.AccountManager, &row.LaunchDev, &row.SkuLive, &row.MouDate, &row.Note,
				&row.AgreementDate, &row.IsActive, &row.Status, &row.CreatorID, &row.CreatorIP, &row.EditorID, &row.EditorIP, &row.Version,
				&row.Created, &row.LastModified, &row.MerchantVillageID, &row.MerchantDistrictID, &row.MerchantCityID,
				&row.MerchantProvinceID, &row.StoreVillageID, &row.StoreDistrictID, &row.StoreCityID, &row.StoreProvinceID,
				&row.BankID, &row.IsClosed, &row.MerchantEmail, &row.BankBranch, &row.PicKtpFile,
				&row.NpwpFile, &row.BusinessType, &row.Source, &row.MerchantTypeString, &row.GenderPicString, &row.MerchantGroup, &row.UpgradeStatus, &row.ProductType,
				&row.LegalEntity, &row.LegalEntityName, &row.NumberOfEmployee, &row.NumberOfEmployeeName, &row.CountUpdateNameAvailable, &row.SellerOfficerName, &row.SellerOfficerEmail, &row.Reason)

			if err != nil {
				helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, params)
				output <- ResultRepository{Error: err}
				return
			}
			if _, err := json.Marshal(&row); err != nil {
				helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, params)
				output <- ResultRepository{Error: err}
				return
			}

			results = append(results, row)
		}

		output <- ResultRepository{Result: results}
	})

	return output
}
func (mr *MerchantRepoPostgres) GetMerchantsPublic(ctxReq context.Context, params *model.QueryParameters) <-chan ResultRepository {
	ctx := "MerchantRepoPostgres-GetMerchantsPublic"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		orderBy := `"userId"`
		sort := "DESC"

		query, bindVar := mr.makeQuery(params)
		if params.OrderBy != "" {
			orderBy = strconv.Quote(params.OrderBy)

		}
		if params.SortBy != "" {
			sort = params.SortBy
		}

		fullQuery := mr.getSelectPublic() + ` FROM b2c_merchant ` + query
		fullQuery += fmt.Sprintf(` ORDER BY %s %s`, orderBy, sort)
		fullQuery += fmt.Sprintf(` LIMIT %d OFFSET %d`, params.Limit, params.Offset)

		tags[helper.TextQuery] = fullQuery

		rows, err := mr.ReadDB.Query(fullQuery, bindVar...)
		if err != nil {
			output <- ResultRepository{Error: err}
			return
		}
		results := []model.B2CMerchantDataPublic{}
		for rows.Next() {
			var (
				row model.B2CMerchantDataPublic
			)
			err := rows.Scan(
				&row.ID, &row.UserID, &row.MerchantName, &row.VanityURL, &row.MerchantCategory, &row.CompanyName,
				&row.Pic, &row.PicOccupation, &row.DailyOperationalStaff, &row.StoreClosureDate, &row.StoreReopenDate,
				&row.StoreActiveShippingDate, &row.MerchantAddress, &row.MerchantVillage, &row.MerchantDistrict,
				&row.MerchantCity, &row.MerchantProvince, &row.ZipCode, &row.StoreAddress, &row.StoreVillage, &row.StoreDistrict,
				&row.StoreCity, &row.StoreProvince, &row.StoreZipCode,
				&row.MerchantDescription, &row.MerchantLogo,
				&row.IsPKP, &row.RichContent, &row.NotificationPreferences,
				&row.MerchantRank, &row.Acquisitor, &row.AccountManager, &row.LaunchDev, &row.SkuLive, &row.MouDate, &row.Note,
				&row.MerchantVillageID, &row.MerchantDistrictID, &row.MerchantCityID,
				&row.MerchantProvinceID, &row.StoreVillageID, &row.StoreDistrictID, &row.StoreCityID, &row.StoreProvinceID,
				&row.IsClosed,
				&row.BusinessType, &row.MerchantTypeString, &row.GenderPicString, &row.MerchantGroup, &row.ProductType, &row.IsActive, &row.Status,
			)

			if err != nil {
				helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, params)
				output <- ResultRepository{Error: err}
				return
			}
			if _, err := json.Marshal(&row); err != nil {
				helper.SendErrorLog(ctxReq, ctx, helper.TextQueryDatabase, err, params)
				output <- ResultRepository{Error: err}
				return
			}

			results = append(results, row)
		}

		output <- ResultRepository{Result: results}
	})

	return output
}

func (mr *MerchantRepoPostgres) makeQuery(params *model.QueryParameters) (string, []interface{}) {
	qstring := `WHERE "deletedAt" IS NULL `
	qs := []string{}
	bindVar := []interface{}{}

	lenParams := 0
	if params.BusinessType != "" {
		lenParams++
		qs = append(qs, fmt.Sprintf(`"businessType" = $%d`, lenParams))
		bindVar = append(bindVar, params.BusinessType)
	}
	if params.IsPKP != "" {
		lenParams++
		qs = append(qs, fmt.Sprintf(`"isPKP" = $%d`, lenParams))
		isPkp, _ := strconv.ParseBool(params.IsPKP)
		bindVar = append(bindVar, isPkp)
	}
	if params.MerchantType != "" {
		merchantTypes := strings.Split(params.MerchantType, ",")
		q := `"merchantType" IN (`
		vw := []string{}
		for _, merchantType := range merchantTypes {
			lenParams++
			vw = append(vw, fmt.Sprintf(`$%d`, lenParams))
			nMerchantType := helper.TrimSpace(merchantType)
			bindVar = append(bindVar, strings.ToUpper(nMerchantType))
		}
		q += strings.Join(vw, ",")

		q += `)`
		qs = append(qs, q)
	}
	if params.Status != "" {
		lenParams++
		qs = append(qs, fmt.Sprintf(`"isActive" = $%d`, lenParams))
		status, _ := strconv.ParseBool(params.Status)
		bindVar = append(bindVar, status)
	}
	if params.UpgradeStatus != "" {
		upgradeStatuses := strings.Split(params.UpgradeStatus, ",")
		q := `"upgradeStatus" IN (`
		vw := []string{}
		for _, status := range upgradeStatuses {
			lenParams++
			vw = append(vw, fmt.Sprintf(`$%d`, lenParams))
			nStatus := helper.TrimSpace(status)
			bindVar = append(bindVar, nStatus)
		}
		q += strings.Join(vw, ",")

		q += `)`
		qs = append(qs, q)
	}
	qss, par, lenParams := mr.parseSearchParam(params, lenParams)
	if qss != "" {
		qs = append(qs, qss)
		bindVar = append(bindVar, par)
	}
	if params.Officer != "" {
		q := `"sellerOfficerName" IN (`
		qe := `"sellerOfficerEmail" IN (`
		var vw []string
		var ve []string
		bindVar, lenParams, vw, ve = checkingEmailOrNameOfficer(params, lenParams, bindVar)
		q += strings.Join(vw, ",")
		q += `)`
		qe += strings.Join(ve, ",")
		qe += `)`
		text := CheckValueEmailOrNameOfficerExist(q, qe)
		newText := "(" + text + ")"
		qs = append(qs, newText)
	}
	if params.MerchantIDS != "" {
		merchantIds := strings.Split(params.MerchantIDS, ",")
		q := `b2c_merchant."id" IN (`
		vw := []string{}
		for _, t := range merchantIds {
			lenParams++
			vw = append(vw, fmt.Sprintf(`$%d`, lenParams))
			merchantId := strings.TrimSpace(t)
			bindVar = append(bindVar, merchantId)
		}
		q += strings.Join(vw, ",")

		q += `)`
		qs = append(qs, q)
	}

	if len(qs) > 0 {
		qstring += ` AND ` + strings.Join(qs, " AND ")
	}

	return qstring, bindVar
}

func checkingEmailOrNameOfficer(params *model.QueryParameters, lenParams int, bindVar []interface{}) ([]interface{}, int, []string, []string) {
	officer := strings.Split(params.Officer, ",")
	vw := []string{}
	ve := []string{}
	n := 0
	for _, t := range officer {
		if strings.Contains(officer[n], "@") {
			lenParams++
			ve = append(ve, fmt.Sprintf(`$%d`, lenParams))
			nOfficer := helper.TrimSpace(t)
			bindVar = append(bindVar, nOfficer)
			n++
		} else {
			lenParams++
			vw = append(vw, fmt.Sprintf(`$%d`, lenParams))
			nOfficer := helper.TrimSpace(t)
			bindVar = append(bindVar, nOfficer)
			n++
		}
	}
	return bindVar, lenParams, vw, ve
}
func CheckValueEmailOrNameOfficerExist(q, qe string) string {
	var text string
	if q != `"sellerOfficerName" IN ()` {
		if qe != `"sellerOfficerEmail" IN ()` {
			text = q + " OR " + qe
		} else {
			text = q
		}
	} else {
		text = qe
	}
	return text
}

func (mr *MerchantRepoPostgres) parseSearchParam(params *model.QueryParameters, lenParams int) (string, interface{}, int) {
	var searchParam, queryString string
	if params.Search != "" {
		lenParams++
		queryString = `(`
		if strings.HasPrefix(params.Search, "MCH") {
			queryString += fmt.Sprintf(`b2c_merchant."id" = $%d`, lenParams)
			searchParam = params.Search
		} else if err := golib.ValidateEmail(params.Search); err == nil {
			queryString += fmt.Sprintf(`lower("merchantEmail") LIKE $%d`, lenParams)
			searchParam = "%" + strings.ToLower(params.Search) + "%"
		} else {
			queryString += fmt.Sprintf(`lower("merchantName") LIKE $%d`, lenParams)
			searchParam = "%" + strings.ToLower(params.Search) + "%"
		}

		queryString += `)`
		// qs = append(qs, q)
		// bindVar = append(bindVar, searchParam)
	}
	return queryString, searchParam, lenParams
}

// LoadLegalEntity validate legal entity
func (mr *MerchantRepoPostgres) LoadLegalEntity(ctxReq context.Context, id int) ResultRepository {
	ctx := "MerchantRepoPostgres-LoadLegalEntity"

	query := `SELECT "id", "name" FROM legal_entities WHERE id = $1`

	stmt, err := mr.Repository.Prepare(ctxReq, query)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, query)
		return ResultRepository{Error: err}
	}
	defer stmt.Close()
	var legalEntity model.LegalEntity

	if err = stmt.QueryRow(id).Scan(&legalEntity.ID, &legalEntity.Name); err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
		return ResultRepository{Error: err}
	}
	return ResultRepository{Result: legalEntity}
}

// LoadCompanySize validate company size
func (mr *MerchantRepoPostgres) LoadCompanySize(ctxReq context.Context, id int) ResultRepository {
	ctx := "MerchantRepoPostgres-LoadEmployeeNumber"

	query := `SELECT "id", "name" FROM number_of_employees WHERE id = $1`

	stmt, err := mr.Repository.Prepare(ctxReq, query)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextPrepareDatabase, err, query)
		return ResultRepository{Error: err}
	}
	defer stmt.Close()
	var companySize model.CompanySize

	if err = stmt.QueryRow(id).Scan(&companySize.ID, &companySize.Name); err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
		return ResultRepository{Error: err}
	}
	return ResultRepository{Result: companySize}
}

func (mr *MerchantRepoPostgres) RejectUpgrade(ctxReq context.Context, merchant model.B2CMerchantDataV2, reasonReject string) <-chan ResultRepository {
	ctx := "MerchantRepo-RejectUpgrade"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		var (
			stmtUpdate   *sql.Stmt
			err          error
			merchantType sql.NullString
		)

		if merchant.MerchantTypeString.Valid {
			merchantType.Valid = true
			merchantType.String = merchant.MerchantTypeString.String
		}

		queryReject := `UPDATE b2c_merchant SET "upgradeStatus"=$1, "merchantType"=$2, "lastModified"=$3, "editorId"=$4, "editorIp"=$5, "reason"=$6 WHERE id=$7;`

		tags[helper.TextQuery] = queryReject
		tags[helper.TextMerchantIDCamel] = merchant.ID

		if mr.Tx != nil {
			stmtUpdate, err = mr.Tx.Prepare(queryReject)
		} else {
			stmtUpdate, err = mr.WriteDB.Prepare(queryReject)
		}

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, merchant.ID)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		if _, err = stmtUpdate.Exec(merchant.UpgradeStatus, merchantType, time.Now(), merchant.EditorID, merchant.EditorIP, reasonReject, merchant.ID); err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, merchant.ID)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Error: nil}
	})
	return output
}

func (mr *MerchantRepoPostgres) ClearRejectUpgrade(ctxReq context.Context, merchant model.B2CMerchantDataV2) <-chan ResultRepository {
	ctx := "MerchantRepo-RejectUpgrade"

	output := make(chan ResultRepository)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)
		var (
			stmtUpdate *sql.Stmt
			err        error
		)

		queryReject := `UPDATE b2c_merchant SET "upgradeStatus"= NULL, "lastModified"=$1, "editorId"=$2, "editorIp"=$3, "reason"=NULL WHERE id=$4;`

		tags[helper.TextQuery] = queryReject
		tags[helper.TextMerchantIDCamel] = merchant.ID

		if mr.Tx != nil {
			stmtUpdate, err = mr.Tx.Prepare(queryReject)
		} else {
			stmtUpdate, err = mr.WriteDB.Prepare(queryReject)
		}

		if err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, merchant.ID)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		if _, err = stmtUpdate.Exec(time.Now(), merchant.EditorID, merchant.EditorIP, merchant.ID); err != nil {
			helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, merchant.ID)
			tags[helper.TextResponse] = err
			output <- ResultRepository{Error: err}
			return
		}

		output <- ResultRepository{Error: nil}
	})
	return output
}
