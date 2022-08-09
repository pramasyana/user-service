package main

import (
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"strconv"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/tracer"
	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/token"
	"gopkg.in/guregu/null.v4/zero"

	merchantModel "github.com/Bhinneka/user-service/src/merchant/v2/model"
	merchantRepo "github.com/Bhinneka/user-service/src/merchant/v2/repo"
	"github.com/Bhinneka/user-service/src/service"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	sharedRepo "github.com/Bhinneka/user-service/src/shared/repository"
	cluster "github.com/bsm/sarama-cluster"
)

const (
	//TextNpwp text
	TextNpwp = "NPWP"
	//TextKTP text
	TextKTP = "KTP"
	//TextMerchantID text
	TextMerchantID = "merchantId"

	//TextMasterBankID text
	TextMasterBankID = "bankId"

	//TextModuleMerchantGWS text
	TextModuleMerchantGWS = "consumerMerchantGWS"
)

func consumeKafkaGWS(sRepository *sharedRepo.Repository,
	merchantRepository merchantRepo.MerchantRepository,
	merchantDocumentRepository merchantRepo.MerchantDocumentRepository,
	merchantBankRepository merchantRepo.MerchantBankRepository,
	accessTokenGenerator token.AccessTokenGenerator,
	kafkaMessaging service.QPublisher,
	brokers []string,
) {
	ctx := "consume_kafka_gws"
	scopeTopic := "kafka_find_gws_topic"
	ctxReq := context.Background()

	topicGwsMerchantBank := golib.GetEnvOrFail(ctx, scopeTopic, "KAFKA_GWS_MERCHANT_BANK")
	topicGwsMerchant := golib.GetEnvOrFail(ctx, scopeTopic, "KAFKA_GWS_MERCHANT")

	//activity service
	activityService := service.NewActivityService("v2")
	merchantService, _ := service.NewMerchantService(kafkaMessaging, activityService)

	// cluster kafka construct with partitions mode
	config := cluster.NewConfig()
	config.ClientID = "user-service-gws-kafka"
	config.Group.Mode = cluster.ConsumerModePartitions

	// init consumer
	topics := []string{topicGwsMerchantBank}

	syncFromGws := golib.GetEnvOrFail(ctx, scopeTopic, "SYNC_MERCHANT_FROM_GWS")
	shouldSync, _ := strconv.ParseBool(syncFromGws)
	if shouldSync {
		topics = append(topics, topicGwsMerchant)
	}

	consumer, err := cluster.NewConsumer(brokers, "consumer-group-gws", topics, config)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, scopeTopic, err, nil)
		os.Exit(1)
	}
	defer consumer.Close()

	signalsGws := make(chan os.Signal, 1)
	signal.Notify(signalsGws, os.Interrupt)

	for {
		select {
		case partitionGws, okGws := <-consumer.Partitions():
			if !okGws {
				os.Exit(1)
			}

			go func(pcGws cluster.PartitionConsumer) {
				for msgGws := range pcGws.Messages() {
					tracer.WithTraceFunc(context.Background(), "KafkaGWSConsumer", func(ctxReq context.Context, tags map[string]interface{}) {

						// get message topic
						topicGws := msgGws.Topic
						tags["topic"] = msgGws.Topic
						tags["partition"] = msgGws.Partition
						tags["offset"] = msgGws.Offset
						tags["message"] = string(msgGws.Value)

						switch topicGws {
						case topicGwsMerchant:
							// update to function
							tag := processGwsMerchant(ctxReq, msgGws.Value, sRepository, merchantRepository, merchantDocumentRepository, merchantService, accessTokenGenerator)
							tags = mergeMaps(tags, tag)
						case topicGwsMerchantBank:
							tag := processGwsMerchantBank(ctxReq, msgGws.Value, merchantBankRepository)
							tags = mergeMaps(tags, tag)
						}
						consumer.MarkOffset(msgGws, "")
					})

				}
			}(partitionGws)

		case <-signalsGws:
			os.Exit(1)
		}
	}
}

// mergeMaps
func mergeMaps(maps ...map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for _, m := range maps {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

func processGwsMerchantBank(ctxReq context.Context, input []byte, merchantBankRepository merchantRepo.MerchantBankRepository) map[string]interface{} {
	var (
		pl  serviceModel.GWSMasterMerchantBankPayloadMessage
		ctx = "processGwsMerchantBank"
	)
	tags := make(map[string]interface{})
	if err := json.Unmarshal(input, &pl); err != nil {
		tags[helper.TextResponse] = err
		helper.SendErrorLog(ctxReq, ctx, "error_unmarshal_payload", err, pl)
		return tags
	}

	data := service.RestructMasterMerchantBankGWS(pl.Data)
	tags[TextMasterBankID] = data.ID
	if err := merchantBankRepository.SaveMasterBankGWS(ctxReq, data); err != nil {
		return tags
	}

	return tags
}

func processGwsMerchant(
	ctxReq context.Context,
	input []byte,
	sRepository *sharedRepo.Repository,
	merchantRepo merchantRepo.MerchantRepository,
	merchantDocumentRepository merchantRepo.MerchantDocumentRepository,
	merchantService *service.MerchantService,
	accessTokenGenerator token.AccessTokenGenerator,
) map[string]interface{} {

	var (
		pl              serviceModel.GWSMerchantPayloadMessage
		oldMerchantData merchantModel.B2CMerchantDataV2
		actionLog       string
		ctx             = "processGwsMerchant"
		data            merchantModel.B2CMerchant
		err             error
		privacy         string
	)
	tags := make(map[string]interface{})
	if err := json.Unmarshal(input, &pl); err != nil {
		tags[helper.TextResponse] = err
		helper.SendErrorLog(ctxReq, ctx, "error_unmarshal_payload_merchant", err, pl)
		return tags
	}
	privacy = "private"
	sRepository.StartTransaction()

	data = service.RestructMerchantGWS(pl.Data)
	tags[TextMerchantID] = data.ID

	if pl.EventType == helper.EventProduceCreateMerchant {
		// save merchant
		actionLog = helper.TextInsertUpper
		err = merchantRepo.SaveMerchantGWS(ctxReq, data)
	} else if pl.EventType == helper.EventProduceUpdateMerchant || pl.EventType == helper.EventProduceDeleteMerchant {
		// update merchant for update or delete
		actionLog = helper.TextUpdateUpper
		merchantResult := merchantRepo.LoadMerchant(ctxReq, data.ID, privacy)
		if merchantResult.Result != nil {
			oldMerchantData = merchantResult.Result.(merchantModel.B2CMerchantDataV2)
		}
		err = merchantRepo.UpdateMerchantGWS(ctxReq, data)
	}

	if err != nil {
		tags[helper.TextResponse] = err
		sRepository.Rollback()
		return tags
	}

	merchantData, err := merchantDocumentProcess(ctxReq, pl, data, merchantRepo, merchantDocumentRepository)
	if err != nil {
		tags[helper.TextResponse] = err
		sRepository.Rollback()
		return tags
	}

	sRepository.Commit()

	if err == nil {
		// publish to kafka for po nav
		// sturgeon consume from gws-merchant and re-publish to user-service-merchant
		merchantService.PublishToKafkaUserMerchant(ctxReq, &merchantData, pl.EventType, "gws")

		// insert log to activity service
		tokenResult := <-accessTokenGenerator.GenerateAnonymous(ctxReq)
		newCtx := context.WithValue(ctxReq, helper.TextAuthorization, tokenResult.AccessToken.AccessToken)
		merchantService.InsertLogMerchant(newCtx, oldMerchantData, merchantData, actionLog, TextModuleMerchantGWS)
	}
	return tags

}
func merchantDocumentProcess(ctxReq context.Context,
	pl serviceModel.GWSMerchantPayloadMessage,
	data merchantModel.B2CMerchant,
	merchantRepository merchantRepo.MerchantRepository,
	merchantDocumentRepository merchantRepo.MerchantDocumentRepository) (merchantModel.B2CMerchantDataV2, error) {

	merchantData := merchantModel.B2CMerchantDataV2{}
	merchantData = service.RestructMerchantDataGWS(pl.Data)

	merchantDocuments := make([]merchantModel.B2CMerchantDocumentData, 0)

	for _, document := range pl.Data.Documents {
		m := document.Value

		if document.Type == TextNpwp {
			data.NpwpFile = &m
			merchantData.NpwpFile = zero.StringFrom(document.Value)
		} else if document.Type == TextKTP {
			data.PicKtpFile = &m
			merchantData.PicKtpFile = zero.StringFrom(document.Value)
		}

		if document.Type == TextKTP || document.Type == TextNpwp {
			// update ktp & npwp merchant
			err := merchantRepository.UpdateMerchantGWS(ctxReq, data)
			if err != nil {
				return merchantData, err
			}
			continue
		}

		// save merchant document
		merchantDocumentModel := service.RestructMerchantDocumentGWS(document, pl.Data)
		err := merchantDocumentRepository.SaveMerchantDocumentGWS(merchantDocumentModel)
		if err != nil {
			return merchantData, err
		}

		merchantDocument := service.RestructMerchantDocumentDataGWS(document, pl.Data)
		merchantDocuments = append(merchantDocuments, merchantDocument)
	}

	merchantData.Documents = merchantDocuments

	return merchantData, nil
}
