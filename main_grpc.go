package main

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/Bhinneka/user-service/helper"
	rpc "github.com/Bhinneka/user-service/src/grpc/servers"
	healthDelivery "github.com/Bhinneka/user-service/src/health/delivery"
	memberDelivery "github.com/Bhinneka/user-service/src/member/v1/delivery"
	"github.com/Bhinneka/user-service/src/service"
	log "github.com/sirupsen/logrus"
)

// GRPCDefaultPort default port for GRPC
const GRPCDefaultPort = 8082

// GRPCServerMain function for serving GRPC
func (s *AppService) GRPCServerMain(publicKey *rsa.PublicKey, kafkaMessaging service.QPublisher) {
	ctx := "GRPCServerMain"

	// set GRPC port
	port := GRPCDefaultPort
	portGRPC, ok := os.LookupEnv("PORT_GRPC")
	if ok {
		intPort, _ := strconv.Atoi(portGRPC)

		port = intPort
	}

	// topic, ok := os.LookupEnv("NSQ_TOPIC_NOTIFICATION")
	// if !ok {
	// 	err := errors.New("you need to specify NSQ_TOPIC_NOTIFICATION in the environment variable")
	// 	helper.Log(log.ErrorLevel, err.Error(), "start_consumer", "find_nsq_topic_config")

	// 	os.Exit(1)
	// }

	kafkaUserServiceTopic, ok := os.LookupEnv("KAFKA_USER_SERVICE_TOPIC")
	if !ok {
		err := errors.New("you need to specify KAFKA_USER_SERVICE_TOPIC in the environment variable")
		helper.Log(log.ErrorLevel, err.Error(), ctx, "find_kafka_server_config")

		os.Exit(1)
	}

	helper.Log(log.InfoLevel, fmt.Sprintf("GRPC server will be running on port %d", port), "grpc_main", "initiate_grpc")

	hGRPCHandler := healthDelivery.NewGRPCHandler(s.HealthUseCase)
	mGRPCHandler := memberDelivery.NewGRPCHandler(s.MemberUseCase, kafkaMessaging, kafkaUserServiceTopic, publicKey)

	grpcServer := rpc.NewGRPCServer(hGRPCHandler, mGRPCHandler)

	err := grpcServer.Serve(uint(port))

	if err != nil {
		err = fmt.Errorf("error in Startup: %s", err.Error())
		helper.Log(log.ErrorLevel, err.Error(), ctx, "serve_grpc")
		os.Exit(1)
	}

}
