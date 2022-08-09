package repo

import (
	"context"

	"github.com/Bhinneka/user-service/src/corporate/v2/model"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
)

// ResultRepository data structure
type ResultRepository struct {
	Result interface{}
	Error  error
}

// AccountRepository interface abstraction
type AccountRepository interface {
	Save(context.Context, sharedModel.B2BAccount) error
	Update(context.Context, sharedModel.B2BAccount) error
	Delete(context.Context, sharedModel.B2BAccount) error
}

// AccountTemporaryRepository interface abstraction
type AccountTemporaryRepository interface {
	Save(context.Context, model.B2BAccountTemporary) error
	Update(context.Context, model.B2BAccountTemporary) error
	Delete(context.Context, model.B2BAccountTemporary) error
}

// AccountContactRepository interface abstraction
type AccountContactRepository interface {
	Save(context.Context, sharedModel.B2BAccountContact) error
	Update(context.Context, sharedModel.B2BAccountContact) error
	Delete(context.Context, sharedModel.B2BAccountContact) error
}

// ContactRepository interface abstraction
type ContactRepository interface {
	Save(context.Context, sharedModel.B2BContact) error
	Update(context.Context, sharedModel.B2BContact) error
	Delete(context.Context, sharedModel.B2BContact) error
}

// AddressRepository interface abstraction
type AddressRepository interface {
	Save(context.Context, sharedModel.B2BAddress) error
	Update(context.Context, sharedModel.B2BAddress) error
	Delete(context.Context, sharedModel.B2BAddress) error
}

// PhoneRepository interface abstraction
type PhoneRepository interface {
	Save(context.Context, sharedModel.B2BPhone) error
	Update(context.Context, sharedModel.B2BPhone) error
	Delete(context.Context, sharedModel.B2BPhone) error
}

// DocumentRepository interface abstraction
type DocumentRepository interface {
	Save(context.Context, sharedModel.B2BDocument) error
	Update(context.Context, sharedModel.B2BDocument) error
	Delete(context.Context, sharedModel.B2BDocument) error
}

// ContactNpwpRepository interface abstraction
type ContactNpwpRepository interface {
	Save(context.Context, model.B2BContactNpwp) error
	Update(context.Context, model.B2BContactNpwp) error
	Delete(context.Context, model.B2BContactNpwp) error
}

// ContactAddressRepository interface abstraction
type ContactAddressRepository interface {
	Save(context.Context, sharedModel.B2BContactAddress) error
	Update(context.Context, sharedModel.B2BContactAddress) error
	Delete(context.Context, sharedModel.B2BContactAddress) error
}

// ContactTempRepository interface abstraction
type ContactTempRepository interface {
	Save(context.Context, model.B2BContactTemp) error
	Update(context.Context, model.B2BContactTemp) error
	Delete(context.Context, model.B2BContactTemp) error
}

// LeadsRepository interface abstraction
type LeadsRepository interface {
	Save(context.Context, sharedModel.B2BLeads) error
	Update(context.Context, sharedModel.B2BLeads) error
	Delete(context.Context, sharedModel.B2BLeads) error
}

type ContactDocumentRepository interface {
	Save(context.Context, sharedModel.B2BContactDocument) error
	Delete(context.Context, sharedModel.B2BContactDocument) error
}
