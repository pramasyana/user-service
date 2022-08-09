package consumer

import (
	"context"

	localConfig "github.com/Bhinneka/user-service/config"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/member/v1/model"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
)

func ExecSharkAccount(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl sharedModel.PayloadAccount
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_account"); err != nil {
		return err
	}

	switch pl.EventType {
	case helper.TextCreate, helper.TextUpdate:
		err = cfg.CorporateAccountRepository.Save(ctxReq, sharedModel.RestructCorporateAccount(pl.Payload))
	case helper.TextDelete:
		// delete account to database
		var accountModel sharedModel.B2BAccount
		err = cfg.CorporateAccountRepository.Delete(ctxReq, accountModel)
	}
	return err
}

func ExecSharkAccountContact(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl sharedModel.PayloadAccountContact
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_account_contact"); err != nil {
		return err
	}

	switch pl.EventType {
	case helper.TextCreate, helper.TextUpdate:
		// save account to database
		accountModel := sharedModel.RestructCorporateAccountContact(pl.Payload)
		err = cfg.CorporateAccountContactRepository.Save(ctxReq, accountModel)
	case helper.TextDelete:
		// delete account to database
		var accountModel sharedModel.B2BAccountContact
		accountModel.ID = pl.Payload.ID
		err = cfg.CorporateAccountContactRepository.Delete(ctxReq, accountModel)
	}
	return err
}

func ExecSharkContact(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl sharedModel.PayloadContact
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_contact"); err != nil {
		return err
	}

	switch pl.EventType {
	case helper.TextCreate, helper.TextUpdate:
		contactModel := sharedModel.RestructCorporateContact(pl.Payload)
		err = cfg.CorporateContactRepository.Save(ctxReq, contactModel)

		// update password member id isSync true
		if contactModel.IsSync {
			var member model.Member
			member.Email = helper.ValidateStringNull(contactModel.Email)
			member.Password = helper.ValidateStringNull(contactModel.Password)
			member.Salt = helper.ValidateStringNull(contactModel.Salt)
			member.LastPasswordModified = *contactModel.LastPasswordModified

			memberResult := <-cfg.MemberRepository.UpdatePasswordMemberByEmail(ctxReq, member)
			if memberResult.Error != nil {
				return memberResult.Error
			}
		}
	case helper.TextDelete:
		var contactModel sharedModel.B2BContact
		contactModel.ID = pl.Payload.ID
		err = cfg.CorporateContactRepository.Delete(ctxReq, contactModel)
	}
	return err
}

func ExecSharkAddress(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl sharedModel.PayloadAddress
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_address"); err != nil {
		return err
	}

	switch pl.EventType {
	case helper.TextCreate, helper.TextUpdate:
		addressModel := sharedModel.RestructCorporateAddress(pl.Payload)
		err = cfg.CorporateAddressRepository.Save(ctxReq, addressModel)
	case helper.TextDelete:
		// delete address to database
		var addressModel sharedModel.B2BAddress
		addressModel.ID = pl.Payload.ID
		err = cfg.CorporateAddressRepository.Delete(ctxReq, addressModel)
	}
	return err
}

func ExecSharkPhone(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl sharedModel.PayloadPhone
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_phone"); err != nil {
		return err
	}

	switch pl.EventType {
	case helper.TextCreate, helper.TextUpdate:
		phoneModel := sharedModel.RestructCorporatePhone(pl.Payload)
		err = cfg.CorporatePhoneRepository.Save(ctxReq, phoneModel)
	case helper.TextDelete:
		var phoneModel sharedModel.B2BPhone
		phoneModel.ID = pl.Payload.ID
		err = cfg.CorporatePhoneRepository.Delete(ctxReq, phoneModel)
	}
	return err
}

func ExecSharkDocument(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl sharedModel.PayloadDocument
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_document"); err != nil {
		return err
	}

	switch pl.EventType {
	case helper.TextCreate, helper.TextUpdate:
		documentModule := sharedModel.RestructCorporateDocument(pl.Payload)
		err = cfg.CorporateDocumentRepository.Save(ctxReq, documentModule)
	case helper.TextDelete:
		var documentModule sharedModel.B2BDocument
		documentModule.ID = pl.Payload.ID
		err = cfg.CorporateDocumentRepository.Delete(ctxReq, documentModule)
	}
	return err
}

func ExecSharkContactAddress(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl sharedModel.PayloadContactAddress
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_contact_address"); err != nil {
		return err
	}

	switch pl.EventType {
	case helper.TextCreate, helper.TextUpdate:
		contactAddressModel := sharedModel.RestructCorporateContactAddress(pl.Payload)
		err = cfg.CorporateContactAddressRepository.Save(ctxReq, contactAddressModel)
	case helper.TextDelete:
		var contactAddressModel sharedModel.B2BContactAddress
		contactAddressModel.ID = pl.Payload.ID
		err = cfg.CorporateContactAddressRepository.Delete(ctxReq, contactAddressModel)
	}
	return err
}

func ExecSharkLeads(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl sharedModel.PayloadLeads
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_leads"); err != nil {
		return err
	}

	switch pl.EventType {
	case helper.TextCreate, helper.TextUpdate:
		leadsModel := sharedModel.RestructCorporateLeads(pl.Payload)
		err = cfg.CorporateLeadsRepository.Save(ctxReq, leadsModel)
	case helper.TextDelete:
		// delete contact to database
		var leadsModel sharedModel.B2BLeads
		leadsModel.ID = pl.Payload.ID
		err = cfg.CorporateLeadsRepository.Delete(ctxReq, leadsModel)
	}
	return err
}

func ExecSharkContactDocument(ctxReq context.Context, cfg localConfig.ServiceRepository, ctx string, input []byte) (err error) {
	var pl sharedModel.PayloadContactDocument
	if err := plCheck(input, &pl, ctx, "unmarshal_payload_contact_document"); err != nil {
		return err
	}

	switch pl.EventType {
	case helper.TextCreate, helper.TextUpdate:
		err = cfg.CorporateContactDocumentRepository.Save(ctxReq, pl.Payload)
	case helper.TextDelete:
		// delete contact to database
		var leadsModel sharedModel.B2BLeads
		leadsModel.ID = pl.Payload.ID
		err = cfg.CorporateContactDocumentRepository.Delete(ctxReq, pl.Payload)
	}
	return err
}
