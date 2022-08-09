package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/Bhinneka/user-service/helper"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
)

func main() {

	var limit string
	fmt.Println("input limit data-merchat that get from sturgeon :")
	fmt.Scanf("%s", &limit)
	fmt.Printf("limit : %s \n", limit)

	dataMerchant, err := getDataMerchantSturgeon(limit)

	if err != nil {
		fmt.Println(err)
	}

	for _, merchant := range dataMerchant.Data {
		responseSendbird, err := createUserMerchantSendbird(merchant)
		if err != nil {
			fmt.Println(err)
		}

		data, _ := json.Marshal(responseSendbird)
		fmt.Printf("%s\n", data)

	}

}

// Meta string
type Meta struct {
	Page         int `json:"page"`
	Limit        int `json:"limit"`
	TotalRecords int `json:"totalRecords"`
	TotalPages   int `json:"totalPages"`
}

// DataMerchant struct
type DataMerchant struct {
	ID           string `json:"id"`
	UserID       string `json:"userID"`
	MerchantName string `json:"merchantName"`
	MerchantLogo string `json:"merchantLogo"`
}

// ListMerchantResponse struct
type ListMerchantResponse struct {
	Success bool           `json:"success"`
	Code    int            `json:"code"`
	Meta    Meta           `json:"meta"`
	Data    []DataMerchant `json:"data"`
}

// get data merchant from sturgeon
func getDataMerchantSturgeon(limit string) (ListMerchantResponse, error) {

	var response ListMerchantResponse
	initialToken := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJhZG0iOnRydWUsImF1ZCI6Im5vdGlmaWNhdGlvbi1zZXJ2aWNlIiwiYXV0aG9yaXNlZCI6dHJ1ZSwiZGlkIjoiYzBiNGQxYjRjNDQ3NCIsImRsaSI6IldFQiIsImVtYWlsIjoicnVzZGkuc3lhaHJlbkBiaGlubmVrYS5jb20iLCJleHAiOjE2MzA2NjczNTAsImlhdCI6MTYzMDY2MDE1MCwiaXNzIjoiYmhpbm5la2EuY29tIiwianRpIjoiYTQ3N2YzNmI3YzZmNTZiNzVmOWJmMDZjZTc4MjE1MGNjZGVkNjZlNSIsIm1lbWJlclR5cGUiOiJwZXJzb25hbCIsInNpZ25VcEZyb20iOiIiLCJzdGFmZiI6dHJ1ZSwic3ViIjoiVVNSMjAwMzA0NzEifQ.RLByY10oiZ3dqAPBVFiXs8Pgvwp0xhgYtNKjSyQdQ2dWKTtYuXoLV_a4y3QUO5mWIfxs--AjxuNcmrMk0sduLh1zgxGowg9HpIhRczfyox67rZya9Em7A9j4n8sUjG2B_1GuAjHjjnMsjVLauxgtrBUqJVa_3NpU6-HvGFZDKAhgPImmeSmbeL2zjewUOgGQyOlcIxyV6V5POMbugYzbvVr985qUzMOkRJ94ohcYTL7IqEa1kjii_i41jfjvhDZEb1z6bUN-0Gx2iTiz6tAHNlyqXdHUpHEG8Nx3FD0huXg-hsQ88ZuRyGudb-zDI6Nv_L4vJF9n9dKK5sHH-UvCOQ"
	token := fmt.Sprintf("Bearer %s", initialToken)
	// generate headers
	headers := map[string]string{
		"Authorization": token,
	}
	baseURL := "https://dev.bhinneka.com/user-service"
	uri := fmt.Sprintf("%s/api/v2/merchant/list?page=1&limit=%s", baseURL, limit)

	contex := context.Background()
	err := helper.GetHTTPNewRequest(contex, http.MethodGet, uri, nil, &response, headers)
	if err != nil {
		return response, err
	}

	return response, nil

}

// create user/merchant to sendbird
func createUserMerchantSendbird(data DataMerchant) (serviceModel.SendbirdResponse, error) {

	var err error
	contex := context.Background()
	var responseSendbird serviceModel.SendbirdResponse
	token := "d387ab27531840a8363fd3efa1a5e1fec312b5ee"
	// generate headers
	headers := map[string]string{
		"Api-Token": token,
	}
	// chat-dev application -> sendbird development
	applicationID := "05ACB8A0-12F7-4505-A783-E57AF527B865"
	uri := fmt.Sprintf("https://api-%s.sendbird.com/v3/users", applicationID)

	// var responseUser serviceModel.SendbirdStringResponse
	var userBody serviceModel.User

	// bind data to Interface user
	userBody.UserID = data.UserID
	bodyMershal, _ := json.Marshal(userBody)
	payload := strings.NewReader(string(bodyMershal))

	err = helper.GetHTTPNewRequest(contex, http.MethodPost, uri, payload, &responseSendbird, headers)
	if err != nil {
		return responseSendbird, err
	}

	uriMetadata := fmt.Sprintf("https://api-%s.sendbird.com/v3/users/%s/metadata", applicationID, data.UserID)
	var bodyMetadata serviceModel.MetadataRequest
	var responseMetadata serviceModel.MetaDataResponse

	var merchant serviceModel.Merchant

	merchant.Reference = data.UserID

	//merchantData, _ := json.Marshal(merchant)

	// bind data to Interface Metadata Request
	bodyMetadata.Metadata.Reference = string(merchant.Reference)
	bodyMetadataMershal, _ := json.Marshal(bodyMetadata)
	payloadMetadata := strings.NewReader(string(bodyMetadataMershal))

	// request for update metadata
	err = helper.GetHTTPNewRequest(contex, http.MethodPost, uriMetadata, payloadMetadata, &responseMetadata, headers)
	if err != nil {
		return responseSendbird, err

	}

	if responseMetadata.Error {
		e := errors.New(responseMetadata.Message)
		return responseSendbird, e

	}

	return responseSendbird, nil
}
