package query

import (
	"context"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/Bhinneka/bhinneka-go-sdk"
	localConfig "github.com/Bhinneka/user-service/config"
	"github.com/Bhinneka/user-service/src/auth/v1/model"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var (
	errDefault      = "default error"
	defaultOauthURL = "https://oauth.bhinnekatesting.com"
)

var testDataAzure = []struct {
	name            string
	wantError       bool
	code            string
	redirectURI     string
	serviceResponse interface{}
	statusCode      int
}{
	{
		name:       "Get Azure #1",
		statusCode: http.StatusOK,
		wantError:  false,
	},
	{
		name:            "Get Azure #2",
		statusCode:      http.StatusBadRequest,
		wantError:       true,
		serviceResponse: `{"code":401}`,
	},
	{
		name:            "Get Azure #3",
		statusCode:      http.StatusOK,
		wantError:       true,
		serviceResponse: model.AuthAzureToken{Error: errDefault},
	},
}

var testDataGoogle = []struct {
	name            string
	wantError       bool
	code            string
	redirectURI     string
	serviceResponse interface{}
	statusCode      int
}{
	{
		name:       "Get Google #1",
		statusCode: http.StatusOK,
		wantError:  false,
	},
	{
		name:            "Get Google #2",
		statusCode:      http.StatusBadRequest,
		wantError:       true,
		serviceResponse: `{"code":401}`,
	},
	{
		name:            "Get Google #3",
		statusCode:      http.StatusOK,
		wantError:       true,
		serviceResponse: model.AuthGoogleToken{Error: errDefault},
	},
}

var testDataFacebook = []struct {
	name            string
	wantError       bool
	code            string
	redirectURI     string
	serviceResponse interface{}
	statusCode      int
}{
	{
		name:       "Get Facebook #1",
		statusCode: http.StatusOK,
		wantError:  false,
	},
	{
		name:            "Get Facebook #2",
		statusCode:      http.StatusBadRequest,
		wantError:       true,
		serviceResponse: `{"code":401}`,
	},
	{
		name:            "Get Facebook #3",
		statusCode:      http.StatusOK,
		wantError:       true,
		serviceResponse: model.AuthFacebookToken{Error: &model.FacebookError{Message: errDefault}},
	},
}

func initTest() localConfig.OAuthService {
	AzureLoginBaseURL, _ := url.Parse(defaultOauthURL)
	AzureTenanID := "tenan-id"
	GoogleBaseTokenURL, _ := url.Parse(defaultOauthURL)
	FacebookTokenBaseURL, _ := url.Parse(defaultOauthURL)
	AppleBaseURL, _ := url.Parse(defaultOauthURL)
	os.Setenv("APPLE_AUTH_URL", defaultOauthURL)

	return localConfig.OAuthService{
		AzureLoginBaseURL:    AzureLoginBaseURL,
		AzureADTenanID:       AzureTenanID,
		GoogleBaseTokenURL:   GoogleBaseTokenURL,
		FacebookTokenBaseURL: FacebookTokenBaseURL,
		AppleBaseURL:         AppleBaseURL,
		AppleTeamID:          "appleTeamID",
		AppleKeyID:           "appleKeyID",
		AzureBaseURL:         AzureLoginBaseURL,
		FacebookBaseURL:      FacebookTokenBaseURL,
		GoogleBaseURL:        GoogleBaseTokenURL,
	}
}

func TestQueryOauthAzure(t *testing.T) {
	cfg := initTest()
	ctx := context.Background()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("Azure Login", func(t *testing.T) {
		for _, tc := range testDataAzure {
			oa := NewAuthQueryOAuth(cfg)
			bhinneka.MockHTTP(http.MethodPost, "https://oauth.bhinnekatesting.com/tenan-id/oauth2/token", tc.statusCode, tc.serviceResponse)
			sr := <-oa.GetAzureToken(ctx, tc.code, tc.redirectURI)
			if tc.wantError {
				assert.Error(t, sr.Error)
			} else {
				assert.NoError(t, sr.Error)
			}
		}
	})
}

var testDataAzureMember = []struct {
	name            string
	wantError       bool
	code            string
	redirectURI     string
	serviceResponse interface{}
	statusCode      int
}{
	{
		name:       "Get Azure Detail #1",
		statusCode: http.StatusOK,
		wantError:  false,
	},
	{
		name:            "Get Azure Detail #2",
		statusCode:      http.StatusBadRequest,
		wantError:       true,
		serviceResponse: `{"code":401}`,
	},
}

func TestAzureDetail(t *testing.T) {
	cfg := initTest()
	ctx := context.Background()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("Azure Get Member Detail", func(t *testing.T) {
		for _, tc := range testDataAzureMember {
			oa := NewAuthQueryOAuth(cfg)
			bhinneka.MockHTTP(http.MethodGet, "https://oauth.bhinnekatesting.com/tenan-id/me?api-version=1.6", tc.statusCode, tc.serviceResponse)
			sr := <-oa.GetDetailAzureMember(ctx, tc.code)
			if tc.wantError {
				assert.Error(t, sr.Error)
			} else {
				assert.NoError(t, sr.Error)
			}
		}
	})
}

func TestQueryOauthGoogle(t *testing.T) {
	cfg := initTest()
	ctx := context.Background()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("Google Login", func(t *testing.T) {
		for _, tc := range testDataGoogle {
			oa := NewAuthQueryOAuth(cfg)
			bhinneka.MockHTTP(http.MethodPost, "https://oauth.bhinnekatesting.com/o/oauth2/token", tc.statusCode, tc.serviceResponse)
			sr := <-oa.GetGoogleToken(ctx, tc.code, tc.redirectURI)
			if tc.wantError {
				assert.Error(t, sr.Error)
			} else {
				assert.NoError(t, sr.Error)
			}
		}
	})
}

func TestQueryOauthFacebook(t *testing.T) {
	cfg := initTest()
	ctx := context.Background()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("Facebook Login", func(t *testing.T) {
		for _, tc := range testDataFacebook {
			oa := NewAuthQueryOAuth(cfg)
			bhinneka.MockHTTP(http.MethodPost, "https://oauth.bhinnekatesting.com/oauth/access_token", tc.statusCode, tc.serviceResponse)
			sr := <-oa.GetFacebookToken(ctx, tc.code, tc.redirectURI)
			if tc.wantError {
				assert.Error(t, sr.Error)
			} else {
				assert.NoError(t, sr.Error)
			}
		}
	})
}

var testDataFacebookMember = []struct {
	name            string
	wantError       bool
	code            string
	redirectURI     string
	serviceResponse interface{}
	statusCode      int
}{
	{
		name:       "Get Facebook Detail #1",
		statusCode: http.StatusOK,
		wantError:  false,
		code:       "token",
	},
	{
		name:            "Get Facebook Detail #2",
		statusCode:      http.StatusBadRequest,
		wantError:       true,
		serviceResponse: `{"code":401}`,
		code:            "token",
	},
}

func TestFacebookDetail(t *testing.T) {
	cfg := initTest()
	ctx := context.Background()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("Facebook Get Member Detail", func(t *testing.T) {
		for _, tc := range testDataFacebookMember {
			oa := NewAuthQueryOAuth(cfg)
			bhinneka.MockHTTP(http.MethodGet, "https://oauth.bhinnekatesting.com/me?fields=id,name,first_name,last_name,email,birthday,gender&access_token=token", tc.statusCode, tc.serviceResponse)
			sr := <-oa.GetDetailFacebookMember(ctx, tc.code)
			if tc.wantError {
				assert.Error(t, sr.Error)
			} else {
				assert.NoError(t, sr.Error)
			}
		}
	})
}

var testGoogleMember = []struct {
	name            string
	wantError       bool
	code            string
	redirectURI     string
	serviceResponse interface{}
	statusCode      int
}{
	{
		name:       "Get Google #1",
		statusCode: http.StatusOK,
		wantError:  false,
	},
	{
		name:            "Get Google #2",
		statusCode:      http.StatusBadRequest,
		wantError:       true,
		serviceResponse: `{"code":401}`,
	},
	{
		name:            "Get Google #3",
		statusCode:      http.StatusOK,
		wantError:       true,
		serviceResponse: model.GoogleOAuth2Response{Error: &model.GooglError{Errors: errDefault}},
	},
}

func TestGoogleDetail(t *testing.T) {
	cfg := initTest()
	ctx := context.Background()

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	t.Run("Google Member", func(t *testing.T) {
		for _, tc := range testGoogleMember {
			oa := NewAuthQueryOAuth(cfg)
			bhinneka.MockHTTP(http.MethodGet, "https://oauth.bhinnekatesting.com/oauth2/v2/userinfo", tc.statusCode, tc.serviceResponse)
			sr := <-oa.GetDetailGoogleMember(ctx, tc.code)
			if tc.wantError {
				assert.Error(t, sr.Error)
			} else {
				assert.NoError(t, sr.Error)
			}
		}
	})
}
