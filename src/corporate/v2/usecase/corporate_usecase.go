package usecase

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/tracer"
	localConfig "github.com/Bhinneka/user-service/config"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/corporate/v2/model"
	"github.com/Bhinneka/user-service/src/corporate/v2/query"
	"github.com/Bhinneka/user-service/src/corporate/v2/repo"
	"github.com/Bhinneka/user-service/src/service"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	sharedModel "github.com/Bhinneka/user-service/src/shared/model"
)

const (
	msgErrorFindContact = "cannot find contact"
)

//CorporateUseCaseImpl data structure
type CorporateUseCaseImpl struct {
	ContactRepo  repo.ContactRepository
	ContactQuery query.ContactQuery
	QPublisher   service.QPublisher
}

// NewCorporateUseCase function for initialise contact use case implementation
func NewCorporateUseCase(
	contactRepo repo.ContactRepository,
	contactQuery query.ContactQuery,
	services localConfig.ServiceShared) CorporateUseCase {
	return &CorporateUseCaseImpl{
		ContactRepo:  contactRepo,
		ContactQuery: contactQuery,
		QPublisher:   services.QPublisher,
	}
}

// GetAllListContact function for getting list of contact
func (s *CorporateUseCaseImpl) GetAllListContact(ctxReq context.Context, params *model.ParametersContact) <-chan ResultUseCase {
	ctx := "CorporateUseCase-GetAllListContact"

	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		var err error
		// validate all parameters
		paging, err := helper.ValidatePagination(
			helper.PaginationParameters{
				Page:     1, // default
				StrPage:  params.StrPage,
				Limit:    20, // default
				StrLimit: params.StrLimit,
			})

		if err != nil {
			tags[helper.TextResponse] = err.Error()
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		params.Page = paging.Page
		params.Limit = paging.Limit
		params.Offset = paging.Offset
		tags[helper.TextParameter] = params

		contactResult := <-s.ContactQuery.GetListContact(ctxReq, params)
		if contactResult.Error != nil {
			httpStatus := http.StatusInternalServerError

			// when data is not found
			if contactResult.Error == sql.ErrNoRows {
				httpStatus = http.StatusNotFound
				contactResult.Error = fmt.Errorf(helper.ErrorDataNotFound, "contact_list")
			}

			output <- ResultUseCase{Error: contactResult.Error, HTTPStatus: httpStatus}
			return
		}

		contact := contactResult.Result.(sharedModel.ListContact)

		totalResult := <-s.ContactQuery.GetTotalContact(ctxReq, params)
		if totalResult.Error != nil {
			output <- ResultUseCase{Error: totalResult.Error, HTTPStatus: http.StatusBadRequest}
			return
		}

		total := totalResult.Result.(int)
		contact.TotalData = total

		output <- ResultUseCase{Result: contact}
	})

	return output
}

// GetDetailContact function for getting detail of contact
func (s *CorporateUseCaseImpl) GetDetailContact(ctxReq context.Context, contactID string) <-chan ResultUseCase {
	ctx := "CorporateUseCase-GetDetailContact"
	output := make(chan ResultUseCase)
	go tracer.WithTraceFunc(ctxReq, ctx, func(ctxReq context.Context, tags map[string]interface{}) {
		defer close(output)

		// find contact by id
		findResult := <-s.ContactQuery.FindByID(ctxReq, contactID)
		if findResult.Error != nil {
			err := errors.New(msgErrorFindContact)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusNotFound}
			return
		}

		result, ok := findResult.Result.(sharedModel.B2BContactData)
		if !ok {
			err := errors.New(msgErrorFindContact)
			tags[helper.TextResponse] = err
			output <- ResultUseCase{Error: err, HTTPStatus: http.StatusBadRequest}
			return
		}
		result.Password = ""
		result.Token = ""
		output <- ResultUseCase{Result: result}

	})
	return output
}

func (s *CorporateUseCaseImpl) ImportContact(ctxReq context.Context, content []byte) ([]*model.ContactPayload, error) {
	// ctx := "CorporateUseCase-ImportContact"

	contacts := make([]*model.ContactPayload, 0)
	csvr := csv.NewReader(bytes.NewBuffer(content))

	// content structure
	// email	firstName	lastName	phoneNumber	accountId	transactionType
	// s.QPublisher.PublishKafka(ctxReq context.Context, topic string, messageKey string, message []byte)
	var emails []string
	for {
		row, err := csvr.Read()

		if err != nil {
			if err == io.EOF {
				emails = nil
				err = nil
				return contacts, s.publishMessage(ctxReq, contacts, "import")
			}
			return contacts, err
		}

		if row[0] == "email" {
			continue
		}
		email := row[0]

		if err := golib.ValidateEmail(email); err != nil {
			return contacts, err
		}
		if golib.StringInSlice(email, emails, false) {
			emails = nil
			err = fmt.Errorf("%s found as a duplicate email address in file source", email)
			return contacts, err
		}

		contact := model.ContactPayload{}
		contact.Email = row[0]
		contact.FirstName = row[1]
		contact.LastName = row[2]
		contact.PhoneNumber = row[3]
		contact.AccountID = row[4]
		contact.TransactionType = row[5]
		contact.CreatedAt = time.Now()
		contacts = append(contacts, &contact)
		emails = append(emails, email)
	}
}

func (s *CorporateUseCaseImpl) publishMessage(ctxReq context.Context, content []*model.ContactPayload, jobType string) error {
	topic := golib.GetEnvOrFail("CorporateUseCaseImpl-publishMessage", "env_shark_import_publisher", "KAFKA_SHARK_IMPORT")
	messages := make([]serviceModel.Messages, 0)

	for _, message := range content {
		payloadM := serviceModel.QueuePayload{
			GeneralPayload: serviceModel.GeneralPayload{
				EventType: jobType,
				Payload:   message,
			},
		}
		byteMessage, err := json.Marshal(payloadM)
		if err != nil {
			return err
		}

		messages = append(messages, serviceModel.Messages{Key: message.AccountID, Content: byteMessage})
	}

	return s.QPublisher.BulkPublishKafka(ctxReq, topic, messages)
}
