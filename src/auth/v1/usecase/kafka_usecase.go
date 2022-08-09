package usecase

import (
	"context"
	"encoding/json"
	"os"
	"time"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/client/v1/model"
	corporateModel "github.com/Bhinneka/user-service/src/corporate/v2/model"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
)

// PublishToKafkaDolphin function for publish new member to kafka dolphin
func (au *AuthUseCaseImpl) PublishToKafkaDolphin(ctxReq context.Context, member serviceModel.MemberDolphin, eventType string) error {
	NSQPayload := serviceModel.DolphinPayloadNSQ{
		EventOrchestration:     "UpdateMember",
		TimestampOrchestration: time.Now().Format(time.RFC3339),
		EventType:              eventType,
		Counter:                0,
		Payload:                member,
	}

	payloadJSON, err := json.Marshal(NSQPayload)
	if err != nil {
		helper.SendErrorLog(ctxReq, "PublishToKafkaDolphin", "unmarshal_payload", err, NSQPayload)
		return err
	}

	if err := au.QPublisher.PublishKafka(ctxReq, au.Topic, member.ID, payloadJSON); err != nil {
		helper.SendErrorLog(ctxReq, "PublishKafka", "publish_to_nav", err, member.Email)
		return err
	}
	return nil
}

// PublishToKafkaDolphin function for publish new member to kafka dolphin
func (au *AuthUseCaseImpl) PublishToKafkaContact(ctxReq context.Context, payload corporateModel.ContactPayload, eventType string) error {
	contactPayload := model.ContactKafka{
		EventType: eventType,
		Payload:   payload,
	}

	payloadJSON, err := json.Marshal(contactPayload)
	if err != nil {
		helper.SendErrorLog(ctxReq, "PublishToKafkaContact", "unmarshal_payload", err, contactPayload)
		return err
	}

	if err := au.QPublisher.PublishKafka(ctxReq, os.Getenv("KAFKA_SHARK_IMPORT"), payload.AccountID, payloadJSON); err != nil {
		helper.SendErrorLog(ctxReq, "PublishKafka", "publish_to_contact", err, payload.Email)
		return err
	}
	return nil
}
