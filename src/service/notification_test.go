package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Bhinneka/bhinneka-go-sdk"
	"github.com/Bhinneka/user-service/src/auth/v1/token"
	jwtMock "github.com/Bhinneka/user-service/src/auth/v1/token/mocks"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func generateUsecaseResult(data token.AccessTokenResponse) <-chan token.AccessTokenResponse {
	output := make(chan token.AccessTokenResponse, 1)
	go func() {
		defer close(output)
		output <- data
	}()
	return output
}

func TestAuthEmailService(t *testing.T) {
	testData := []struct {
		name           string
		expectError    bool
		expectedResult token.AccessTokenResponse
		emailHost      string
	}{
		{
			name:           "Test #1 ",
			expectError:    false,
			emailHost:      "https://someUrl.com",
			expectedResult: token.AccessTokenResponse{AccessToken: token.AccessToken{AccessToken: "something"}},
		},
		{
			name:           "Test #2 ",
			expectError:    true,
			emailHost:      "https://someOtherUrl.com",
			expectedResult: token.AccessTokenResponse{Error: errors.New("some")},
		},
	}
	os.Setenv("EMAIL_NOTIF_USER", "myname")
	os.Setenv("EMAIL_NOTIF_PASS", "mypass")

	for _, tc := range testData {
		mockJwt := jwtMock.AccessTokenGenerator{}
		os.Setenv("EMAIL_NOTIF_HOST", tc.emailHost)

		m := NewNotificationService(&mockJwt)
		mockJwt.On("GenerateAnonymous", mock.Anything).Return(generateUsecaseResult(tc.expectedResult))

		res, err := m.auth(context.Background())
		if tc.expectError {
			assert.Error(t, err)
			assert.Equal(t, res, "")
		} else {
			assert.NoError(t, err)
			assert.NotEqual(t, res, "")
		}
	}
}

func TestSendEmail(t *testing.T) {
	testData := []struct {
		name          string
		expectError   bool
		tokenResponse token.AccessTokenResponse
		emailHost     string
		emailResponse interface{}
	}{
		{
			name:          "Test SendEmail #1",
			expectError:   true,
			emailHost:     "https://someUrl.com",
			tokenResponse: token.AccessTokenResponse{Error: errors.New("something bad happenned")},
		},
		{
			name:          "Test SendEmail #2",
			expectError:   false,
			emailHost:     "https://someEmailUrl.com",
			tokenResponse: token.AccessTokenResponse{AccessToken: token.AccessToken{AccessToken: "something other"}},
			emailResponse: serviceModel.SuccessMessage{Data: serviceModel.APIResponse{Attributes: serviceModel.Attributes{Message: "success message"}}},
		},
		{
			name:          "Test SendEmail #3",
			expectError:   false,
			emailHost:     "https://someOtherUrl.com",
			tokenResponse: token.AccessTokenResponse{AccessToken: token.AccessToken{AccessToken: "something other"}},
			emailResponse: serviceModel.SuccessMessage{Data: serviceModel.APIResponse{Attributes: serviceModel.Attributes{Message: "success message"}}},
		},
	}
	for _, tc := range testData {
		mockJwt := jwtMock.AccessTokenGenerator{}
		os.Setenv("EMAIL_NOTIF_HOST", tc.emailHost)

		m := NewNotificationService(&mockJwt)
		mockJwt.On("GenerateAnonymous", mock.Anything).Return(generateUsecaseResult(tc.tokenResponse))

		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		bhinneka.MockHTTP(http.MethodPost, fmt.Sprintf("%s/email/send", tc.emailHost), http.StatusOK, tc.emailResponse)

		email := serviceModel.Email{
			To:   []string{"my.email@bhinneka.com"},
			From: "noreply@somedomain.com",
		}
		resp, err := m.SendEmail(context.Background(), email)
		if !tc.expectError {
			assert.NoError(t, err)
			assert.NotEqual(t, "", resp)
		} else {
			assert.Error(t, err)
			assert.Equal(t, "", resp)
		}
	}
}
