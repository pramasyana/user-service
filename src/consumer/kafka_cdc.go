package consumer

import (
	"context"
	"encoding/json"

	localConfig "github.com/Bhinneka/user-service/config"
	"github.com/Bhinneka/user-service/helper"
	corporateModel "github.com/Bhinneka/user-service/src/corporate/v2/model"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
)

func plCheck(msg []byte, pl interface{}, ctx string, scope string) error {
	err := json.Unmarshal(msg, &pl)
	if err != nil {
		helper.SendErrorLog(context.Background(), ctx, scope, err, pl)
		return err
	}

	return nil
}

func ProcessSharkAccount(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl sharedModel.AccountPayloadCDC
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_account"); err != nil {
		return err
	}

	switch pl.Payload.Op {
	case "c", "u":
		// save account to database
		accountModel := sharedModel.RestructCorporateAccount(pl.Payload.After)
		err = cfg.CorporateAccountRepository.Save(ctxReq, accountModel)
	case "d":
		// delete account to database
		var accountModel sharedModel.B2BAccount
		accountModel.ID = pl.Payload.Before.ID
		err = cfg.CorporateAccountRepository.Delete(ctxReq, accountModel)
	}
	return err
}

func ProcessSharkAccountTemp(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl serviceModel.AccountTemporaryPayloadCDC
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_account_temporary"); err != nil {
		return err
	}

	switch pl.Payload.Op {
	case "c", "u":
		// save account to database
		accountModel := corporateModel.RestructCorporateAccountTemporary(pl)
		err = cfg.CorporateAccountTempRepository.Save(ctxReq, accountModel)
	case "d":
		// delete account to database
		var accountModel corporateModel.B2BAccountTemporary
		accountModel.ID = pl.Payload.Before.ID
		err = cfg.CorporateAccountTempRepository.Delete(ctxReq, accountModel)
	}
	return err
}

func ProcessSharkAccountContact(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl sharedModel.AccountContactPayloadCDC
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_account_contact"); err != nil {
		return err
	}

	switch pl.Payload.Op {
	case "c", "u":
		// save account to database
		accountModel := sharedModel.RestructCorporateAccountContact(pl.Payload.After)
		err = cfg.CorporateAccountContactRepository.Save(ctxReq, accountModel)
	case "d":
		// delete account to database
		var accountModel sharedModel.B2BAccountContact
		accountModel.ID = pl.Payload.Before.ID
		err = cfg.CorporateAccountContactRepository.Delete(ctxReq, accountModel)
	}
	return err
}

func ProcessSharkContact(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl sharedModel.ContactPayloadCDC
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_contact"); err != nil {
		return err
	}

	switch pl.Payload.Op {
	case "c", "u":
		// save contact to database
		contactModel := sharedModel.RestructCorporateContact(pl.Payload.After)
		err = cfg.CorporateContactRepository.Save(ctxReq, contactModel)
	case "d":
		// delete contact to database
		var contactModel sharedModel.B2BContact
		contactModel.ID = pl.Payload.Before.ID
		err = cfg.CorporateContactRepository.Delete(ctxReq, contactModel)
	}
	return err
}

func ProcessSharkAddress(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl sharedModel.AddressPayloadCDC
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_address"); err != nil {
		return err
	}

	switch pl.Payload.Op {
	case "c", "u":
		// save address to databaseID
		addressModel := sharedModel.RestructCorporateAddress(pl.Payload.After)
		err = cfg.CorporateAddressRepository.Save(ctxReq, addressModel)
	case "d":
		// delete address to database
		var addressModel sharedModel.B2BAddress
		addressModel.ID = pl.Payload.Before.ID
		err = cfg.CorporateAddressRepository.Delete(ctxReq, addressModel)
	}
	return err
}

func ProcessSharkPhone(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl sharedModel.PhonePayloadCDC
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_phone"); err != nil {
		return err
	}

	switch pl.Payload.Op {
	case "c", "u":
		// save phone to databaseID
		phoneModel := sharedModel.RestructCorporatePhone(pl.Payload.After)
		err = cfg.CorporatePhoneRepository.Save(ctxReq, phoneModel)
	case "d":
		// delete phone to database
		var phoneModel sharedModel.B2BPhone
		phoneModel.ID = pl.Payload.Before.ID
		err = cfg.CorporatePhoneRepository.Delete(ctxReq, phoneModel)
	}
	return err
}

func ProcessSharkDocument(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl sharedModel.DocumentPayloadCDC
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_document"); err != nil {
		return err
	}

	switch pl.Payload.Op {
	case "c", "u":
		// save document to database
		documentModule := sharedModel.RestructCorporateDocument(pl.Payload.After)
		err = cfg.CorporateDocumentRepository.Save(ctxReq, documentModule)
	case "d":
		// delete document to database
		var documentModule sharedModel.B2BDocument
		documentModule.ID = pl.Payload.Before.ID
		err = cfg.CorporateDocumentRepository.Delete(ctxReq, documentModule)
	}
	return err
}

func ProcessSharkContactNPWP(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl serviceModel.ContactNpwpPayloadCDC
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_contact_npwp"); err != nil {
		return err
	}

	switch pl.Payload.Op {
	case "c", "u":
		// save contactNpwp to database
		contactNpwpModule := corporateModel.RestructCorporateContactNpwp(pl)
		err = cfg.CorporateContactNPWPRepository.Save(ctxReq, contactNpwpModule)
	case "d":
		// delete contactNpwp to database
		var contactNpwpModule corporateModel.B2BContactNpwp
		contactNpwpModule.ID = pl.Payload.Before.ID
		err = cfg.CorporateContactNPWPRepository.Delete(ctxReq, contactNpwpModule)
	}
	return err
}

func ProcessSharkContactAddress(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl sharedModel.ContactAddressPayloadCDC
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_contact_address"); err != nil {
		return err
	}

	switch pl.Payload.Op {
	case "c", "u":
		// save contact address to database
		contactAddressModel := sharedModel.RestructCorporateContactAddress(pl.Payload.After)
		err = cfg.CorporateContactAddressRepository.Save(ctxReq, contactAddressModel)
	case "d":
		// delete contact address to database
		var contactAddressModel sharedModel.B2BContactAddress
		contactAddressModel.ID = pl.Payload.Before.ID
		err = cfg.CorporateContactAddressRepository.Delete(ctxReq, contactAddressModel)
	}
	return err
}

func ProcessSharkContactTemp(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl serviceModel.ContactTempPayloadCDC
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_contact_temp"); err != nil {
		return err
	}

	switch pl.Payload.Op {
	case "c", "u":
		// save contact to database
		contactTempModel := corporateModel.RestructCorporateContactTemp(pl)
		err = cfg.CorporateContactTempRepository.Save(ctxReq, contactTempModel)
	case "d":
		// delete contact to database
		var contactTempModel corporateModel.B2BContactTemp
		contactTempModel.ID = pl.Payload.Before.ID
		err = cfg.CorporateContactTempRepository.Delete(ctxReq, contactTempModel)
	}
	return err
}

func ProcessSharkLeads(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl sharedModel.LeadsPayloadCDC
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_leads"); err != nil {
		return err
	}

	switch pl.Payload.Op {
	case "c", "u":
		// save leads to database
		leadsModel := sharedModel.RestructCorporateLeads(pl.Payload.After)
		err = cfg.CorporateLeadsRepository.Save(ctxReq, leadsModel)
	case "d":
		// delete contact to database
		var leadsModel sharedModel.B2BLeads
		leadsModel.ID = pl.Payload.Before.ID
		err = cfg.CorporateLeadsRepository.Delete(ctxReq, leadsModel)
	}
	return err
}
