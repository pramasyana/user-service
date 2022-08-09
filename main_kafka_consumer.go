package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/user-service/helper"
	memberModel "github.com/Bhinneka/user-service/src/member/v1/model"
	memberRepo "github.com/Bhinneka/user-service/src/member/v1/repo"
	"github.com/Bhinneka/user-service/src/service"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/Bhinneka/user-service/src/shared"
	cluster "github.com/bsm/sarama-cluster"
)

func consumeKafka(dolphinLogRepository memberRepo.DolphinLogRepository, dolphinService *service.DolphinService, brokers []string) {
	ctx := "consume_kafka"
	scopeTopic := "kafka_find_topic"
	ctxReq := context.Background()

	topicUserService := golib.GetEnvOrFail(ctx, scopeTopic, "KAFKA_USER_SERVICE_TOPIC")
	// cluster kafka construct with partitions mode
	config := cluster.NewConfig()
	config.ClientID = "user-service-kafka"
	config.Group.Mode = cluster.ConsumerModePartitions

	// init consumer
	topics := []string{topicUserService}
	consumer, err := cluster.NewConsumer(brokers, "consumer-group", topics, config)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, scopeTopic, err, config)
		os.Exit(1)
	}
	defer consumer.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case partition, ok := <-consumer.Partitions():
			if !ok {
				os.Exit(1)
			}

			go func(pc cluster.PartitionConsumer) {
				for msg := range pc.Messages() {
					switch msg.Topic {
					case topicUserService:
						err = processSturgeonService(ctxReq, msg.Key, msg.Value, dolphinLogRepository, dolphinService)
					}
					//mark message as processed
					consumer.MarkOffset(msg, "")

				}
			}(partition)

		case <-signals:
			os.Exit(1)
		}
	}
}

func processSturgeonService(ctxReq context.Context, keyMessage, input []byte, dolphinLogRepository memberRepo.DolphinLogRepository, dolphinService *service.DolphinService) error {
	// get message key, and convert to shared.MessageKey type
	key := shared.MessageKeyFromString(string(keyMessage))
	switch key {
	case shared.MemberRegistration:
		return processDolphinRegistration(ctxReq, input, key.String(), dolphinLogRepository, dolphinService)
	case shared.MemberUpdate:
		return processDolphinUpdate(ctxReq, input, key.String(), dolphinLogRepository, dolphinService)
	case shared.MemberActivation:
		return processDolphinActivation(ctxReq, input, key.String(), dolphinLogRepository, dolphinService)
	}
	return nil
}

func processDolphinRegistration(ctxReq context.Context, input []byte, payloadID string, dolphinLogRepository memberRepo.DolphinLogRepository, dolphinService *service.DolphinService) (err error) {
	var (
		pl  serviceModel.DolphinPayloadNSQ
		ctx = "consumeDolphinRegistration"
	)
	if err := json.Unmarshal(input, &pl); err != nil {
		helper.SendErrorLog(ctxReq, ctx, "error_unmarshal_payload", err, pl)
		return err
	}

	// send to dolphin to register member
	if err = dolphinService.RegisterMember(ctxReq, pl.Payload); err != nil {
		return err
	}

	if err = saveLogDolphin(ctxReq, pl, payloadID, dolphinLogRepository); err != nil {
		return err
	}

	return nil
}

func processDolphinUpdate(ctxReq context.Context, input []byte, payloadID string, dolphinLogRepository memberRepo.DolphinLogRepository, dolphinService *service.DolphinService) (err error) {
	var (
		pl  serviceModel.DolphinPayloadNSQ
		ctx = "consumeDolphineUpdate"
	)
	if err := json.Unmarshal(input, &pl); err != nil {
		helper.SendErrorLog(ctxReq, ctx, "error_unmarshal_payload_update", err, pl)
		return err
	}

	if err = dolphinService.UpdateMember(ctxReq, pl.Payload); err != nil {
		return err
	}

	if err = saveLogDolphin(ctxReq, pl, payloadID, dolphinLogRepository); err != nil {
		return err
	}

	return nil
}

func processDolphinActivation(ctxReq context.Context, input []byte, payloadID string, dolphinLogRepository memberRepo.DolphinLogRepository, dolphinService *service.DolphinService) (err error) {
	var (
		pl  serviceModel.DolphinPayloadNSQ
		ctx = "consumeDolphineActivate"
	)
	if err := json.Unmarshal(input, &pl); err != nil {
		helper.SendErrorLog(ctxReq, ctx, "error_unmarshal_payload_activate", err, pl)
		return nil
	}

	if err = dolphinService.ActivateMember(ctxReq, pl.Payload); err != nil {
		return err
	}

	if err = saveLogDolphin(ctxReq, pl, payloadID, dolphinLogRepository); err != nil {
		return err
	}
	return nil
}

// save dolphin log to database
func saveLogDolphin(ctxReq context.Context, pl serviceModel.DolphinPayloadNSQ, payloadID string, dolphinLogRepository memberRepo.DolphinLogRepository) error {
	ctx := "saveDolphinLog"
	var dolphinLog memberModel.DolphinLog
	dolphinLog.UserID = pl.Payload.ID
	dolphinLog.EventType = pl.EventType
	dolphinLog.LogData = &pl
	dolphinLog.Created = time.Now()

	if err := dolphinLogRepository.Save(ctxReq, &dolphinLog); err != nil {
		helper.SendErrorLog(ctxReq, ctx, "error_save_log", err, payloadID)
		return err
	}
	return nil
}
