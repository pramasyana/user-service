package service

import (
	"context"
	"net/url"
	"os"

	"github.com/Bhinneka/user-service/helper"
	"github.com/Bhinneka/user-service/src/auth/v1/token"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
)

// NotificationService data structure
type NotificationService struct {
	BaseURL           *url.URL
	AuthBasicUserName string
	AuthBasicPassword string
	DeviceID          string
	TokenGenerator    token.AccessTokenGenerator
}

// NewNotificationService function for creating general third party client of email
func NewNotificationService(tokenGenerator token.AccessTokenGenerator) *NotificationService {
	ctx := "NewNotificationService"
	ctxReq := context.Background()
	emailBaseURL, err := url.Parse(os.Getenv("EMAIL_NOTIF_HOST"))

	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "init_service", err, nil)
		return nil
	}

	return &NotificationService{
		BaseURL:           emailBaseURL,
		AuthBasicUserName: os.Getenv("EMAIL_NOTIF_USER"),
		AuthBasicPassword: os.Getenv("EMAIL_NOTIF_PASS"),
		DeviceID:          serviceModel.DeviceIDAuth,
		TokenGenerator:    tokenGenerator,
	}
}

// Auth function
func (em *NotificationService) auth(ctxReq context.Context) (string, error) {
	tokenResult := <-em.TokenGenerator.GenerateAnonymous(ctxReq)
	if tokenResult.Error != nil {
		return "", tokenResult.Error
	}
	return tokenResult.AccessToken.AccessToken, nil
}
