package repo

import (
	"context"
	"database/sql"

	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/merchant/v2/model"
)

// Save function for saving  data
func (mr *MerchantRepoPostgres) Save(merchant model.B2CMerchant) error {
	ctx := "MerchantRepo-create"

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
					"bankId", "isClosed", "merchantEmail", "bankBranch", "picKtpFile", "npwpFile", "businessType",  "source"
				)
			VALUES
				(
					$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 
					$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
					$21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
					$31, $32, $33, $34, $35, $36, $37, $38, $39, $40,
					$41, $42, $43, $44, $45, $46, $47, $48, $49, $50,
					$51, $52, $53, $54, $55, $56, $57, $58, $59, $60, 
					$61, $62, $63, $64, $66, $66, $67, $68, $69

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
			"bankId"=$62, "isClosed"=$63, "merchantEmail"=$64, "bankBranch"=$65, "picKtpFile"=$66, "npwpFile"=$67, "businessType"=$68, "source"=$69`

	tr := tracer.StartTrace(context.Background(), ctx)
	ctxReq := tr.NewChildContext()
	tags := make(map[string]interface{})
	defer func() {
		tr.Finish(tags)
	}()

	stmt, err := mr.WriteDB.Prepare(query)

	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, merchant)
		tags[helper.TextResponse] = err
		return err
	}

	tags[helper.TextQuery] = query
	tags[helper.TextMerchantIDCamel] = merchant.ID
	tags["args"] = merchant

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
	)

	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
		tags[helper.TextResponse] = err
		return err
	}

	return nil
}

// Delete function for delete merchant data
func (mr *MerchantRepoPostgres) Delete(id string) error {
	ctx := "MerchantRepo-delete"

	queryDelete := `DELETE FROM b2c_merchant WHERE id=$1;`
	tr := tracer.StartTrace(context.Background(), ctx)
	ctxReq := tr.NewChildContext()
	tags := make(map[string]interface{})
	defer func() {
		tr.Finish(tags)
	}()

	stmt, err := mr.WriteDB.Prepare(queryDelete)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, id)
		tags[helper.TextResponse] = err
		return err
	}

	tags[helper.TextQuery] = queryDelete
	tags[helper.TextMerchantIDCamel] = id
	_, err = stmt.Exec(id)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, id)
		tags[helper.TextResponse] = err
		return err
	}

	return nil
}

// UpdateMerchantGWS function for update merchant data
func (mr *MerchantRepoPostgres) UpdateMerchantGWS(ctxReq context.Context, merchant model.B2CMerchant) error {
	ctx := "MerchantRepo-UpdateMerchantGWS"
	query := `UPDATE b2c_merchant SET "userId"=$2, "merchantName"=$3, "vanityURL"=$4, "merchantCategory"=$5, "companyName"=$6, 
	"pic"=$7, "picOccupation"=$8, "dailyOperationalStaff"=$9, "storeClosureDate"=$10, "storeReopenDate"=$11, 
	"storeActiveShippingDate"=$12, "merchantAddress"=$13, "merchantVillage"=$14, "merchantDistrict"=$15, 
	"merchantCity"=$16, "merchantProvince"=$17, "zipCode"=$18, "storeAddress"=$19, "storeVillage"=$20, "storeDistrict"=$21, 
	"storeCity"=$22, "storeProvince"=$23, "storeZipCode"=$24, "phoneNumber"=$25, "mobilePhoneNumber"=$26, 
	"additionalEmail"=$27, "merchantDescription"=$28, "merchantLogo"=$29, "accountHolderName"=$30,
	"bankName"=$31, "accountNumber"=$32, "isPKP"=$33, "npwp"=$34, "npwpHolderName"=$35, "richContent"=$36, "notificationPreferences"=$37, 
	"merchantRank"=$38, "acquisitor"=$39, "accountManager"=$40, "launchDev"=$41, "skuLive"=$42, "mouDate"=$43, "note"=$44, 
	"agreementDate"=$45, "isActive"=$46, "creatorId"=$47, "creatorIp"=$48, "editorId"=$49, "editorIp"=$50, "version"=$51, 
	"created"=$52, "lastModified"=$53, "merchantVillageId"=$54, "merchantDistrictId"=$55, "merchantCityId"=$56, 
	"merchantProvinceId"=$57, "storeVillageId"=$58, "storeDistrictId"=$59, "storeCityId"=$60, "storeProvinceId"=$61, 
	"bankId"=$62, "isClosed"=$63, "merchantEmail"=$64, "bankBranch"=$65, "picKtpFile"=$66, "npwpFile"=$67, "businessType"=$68, "source"=$69, "deletedAt"=$70,
	"merchantType"=$71, "genderPic"=$72, "merchantGroup"=$73 , "upgradeStatus"=$74, "productType"=$75, "legalEntity"=$76, "numberOfEmployee"=$77
	WHERE "id"=$1;`

	err := mr.SaveUpdateMerchantGWS(ctxReq, merchant, query)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
		return err
	}

	return nil
}

// SaveUpdateMerchantGWS function for saving  data
func (mr *MerchantRepoPostgres) SaveUpdateMerchantGWS(ctxReq context.Context, merchant model.B2CMerchant, query string) error {
	ctx := "MerchantRepo-SaveUpdateMerchantGWS"

	tr := tracer.StartTrace(ctxReq, ctx)
	tags := make(map[string]interface{})
	defer func() {
		tr.Finish(tags)
	}()

	stmt, err := mr.WriteDB.Prepare(query)

	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
		tags[helper.TextResponse] = err
		return err
	}

	tags[helper.TextQuery] = query
	tags["args"] = merchant

	var (
		deletedAtString, storeClosureDate, storeReopenDate                               sql.NullString
		storeActiveShippingDate, mouDate, agreementDate                                  sql.NullString
		skuLive, merchantCategory, merchantType, genderPic, merchantGroup, upgradeStatus sql.NullString
	)

	if len(*merchant.StoreClosureDate) > 0 {
		storeClosureDate.Valid = true
		storeClosureDate.String = *merchant.StoreClosureDate
	}

	if len(*merchant.StoreReopenDate) > 0 {
		storeReopenDate.Valid = true
		storeReopenDate.String = *merchant.StoreReopenDate
	}

	if len(*merchant.StoreActiveShippingDate) > 0 {
		storeActiveShippingDate.Valid = true
		storeActiveShippingDate.String = *merchant.StoreActiveShippingDate
	}

	if len(*merchant.MouDate) > 0 {
		mouDate.Valid = true
		mouDate.String = *merchant.MouDate
	}

	if len(*merchant.AgreementDate) > 0 {
		agreementDate.Valid = true
		agreementDate.String = *merchant.AgreementDate
	}

	if len(*merchant.DeletedAt) > 0 {
		deletedAtString.Valid = true
		deletedAtString.String = *merchant.DeletedAt
	}

	if merchant.SkuLive != nil && len(*merchant.SkuLive) > 0 {
		skuLive.Valid = true
		skuLive.String = *merchant.SkuLive
	}

	if merchant.MerchantCategory != nil && len(*merchant.MerchantCategory) > 0 {
		merchantCategory.Valid = true
		merchantCategory.String = *merchant.MerchantCategory
	}

	if len(*merchant.MerchantType) > 0 {
		merchantType.Valid = true
		merchantType.String = *merchant.MerchantType
	}

	if len(*merchant.GenderPic) > 0 {
		genderPic.Valid = true
		genderPic.String = *merchant.GenderPic
	}

	if len(*merchant.MerchantGroup) > 0 {
		merchantGroup.Valid = true
		merchantGroup.String = *merchant.MerchantGroup
	}

	if len(*merchant.UpgradeStatus) > 0 {
		upgradeStatus.Valid = true
		upgradeStatus.String = *merchant.UpgradeStatus
	}

	_, err = stmt.Exec(
		merchant.ID, merchant.UserID, merchant.MerchantName, merchant.VanityURL,
		merchantCategory, merchant.CompanyName, merchant.Pic, merchant.PicOccupation,
		merchant.DailyOperationalStaff, storeClosureDate, storeReopenDate, storeActiveShippingDate,
		merchant.MerchantAddress, merchant.MerchantVillage, merchant.MerchantDistrict,
		merchant.MerchantCity, merchant.MerchantProvince, merchant.ZipCode, merchant.StoreAddress,
		merchant.StoreVillage, merchant.StoreDistrict, merchant.StoreCity, merchant.StoreProvince,
		merchant.StoreZipCode, merchant.PhoneNumber, merchant.MobilePhoneNumber,
		merchant.AdditionalEmail, merchant.MerchantDescription, merchant.MerchantLogo,
		merchant.AccountHolderName, merchant.BankName, merchant.AccountNumber, merchant.IsPKP,
		merchant.Npwp, merchant.NpwpHolderName, merchant.RichContent, merchant.NotificationPreferences,
		merchant.MerchantRank, merchant.Acquisitor, merchant.AccountManager,
		merchant.LaunchDev, skuLive, mouDate, merchant.Note,
		agreementDate, merchant.IsActive, merchant.CreatorID,
		merchant.CreatorIP, merchant.EditorID, merchant.EditorIP, merchant.Version,
		merchant.Created, merchant.LastModified,
		merchant.MerchantVillageID, merchant.MerchantDistrictID, merchant.MerchantCityID,
		merchant.MerchantProvinceID, merchant.StoreVillageID,
		merchant.StoreDistrictID, merchant.StoreCityID, merchant.StoreProvinceID, merchant.BankID,
		merchant.IsClosed, merchant.MerchantEmail,
		merchant.BankBranch, merchant.PicKtpFile, merchant.NpwpFile, merchant.BusinessType, merchant.Source, deletedAtString,
		merchantType, genderPic, merchantGroup, upgradeStatus, merchant.ProductType, merchant.LegalEntity, merchant.NumberOfEmployee,
	)

	return err
}

// SaveMerchantGWS function for saving  data
func (mr *MerchantRepoPostgres) SaveMerchantGWS(ctxReq context.Context, merchant model.B2CMerchant) error {
	ctx := "MerchantRepo-SaveMerchantGWS"
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
					"bankId", "isClosed", "merchantEmail", "bankBranch", "picKtpFile", "npwpFile", "businessType",  
					"source", "deletedAt", "merchantType", "genderPic", "merchantGroup", "upgradeStatus", "productType", "legalEntity", "numberOfEmployee"
				)
			VALUES
				(
					$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, 
					$11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
					$21, $22, $23, $24, $25, $26, $27, $28, $29, $30,
					$31, $32, $33, $34, $35, $36, $37, $38, $39, $40,
					$41, $42, $43, $44, $45, $46, $47, $48, $49, $50,
					$51, $52, $53, $54, $55, $56, $57, $58, $59, $60, 
					$61, $62, $63, $64, $66, $66, $67, $68, $69, $70,
					$71, $72, $73, $74, $75, $76, $77
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
			"bankId"=$62, "isClosed"=$63, "merchantEmail"=$64, "bankBranch"=$65, "picKtpFile"=$66, "npwpFile"=$67, "businessType"=$68, 
			"source"=$69, "deletedAt"=$70, "merchantType"=$71, "genderPic"=$72, "merchantGroup"=$73, "upgradeStatus"=$74, "productType"=$75,
			"legalEntity"=$76, "numberOfEmployee"=$77`

	err := mr.SaveUpdateMerchantGWS(ctxReq, merchant, query)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, helper.TextExecQuery, err, query)
		return err
	}

	return nil
}
