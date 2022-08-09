package main

import (
	"fmt"
	"strings"

	memberQuery "github.com/Bhinneka/user-service/src/member/v1/query"
	"github.com/Bhinneka/user-service/src/shared"
)

func redisSubscribe(pubSubClient *shared.RedisPubSub, memberQueryWrite memberQuery.MemberQuery) {
	reply := make(chan []byte)
	key := "__keyevent@*__:expired"
	if err := pubSubClient.Subscribe(key, reply); err != nil {
		panic(err)
	}

	for {
		select {
		case msg := <-reply:
			if strings.Contains(string(msg), "ATTEMPT:") {
				message := string(msg)
				messages := strings.Split(message, ":")
				email := messages[1]

				fmt.Println("------------------")
				fmt.Println(message)
				fmt.Println(email)

				updateResult := <-memberQueryWrite.UnblockMember(email)
				if updateResult.Error != nil {
					fmt.Println(updateResult.Error)
				}
			}
		}
	}
}
