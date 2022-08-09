package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRefreshTokenModel(t *testing.T) {

	refreshToken := &RefreshToken{ID: "1", Token: "788888"}

	t.Run("Refresh Token Match", func(t *testing.T) {

		assert.True(t, refreshToken.Match("788888"))
		assert.False(t, refreshToken.Match("788881"))
	})
}

func TestValidBhinnekaEmail(t *testing.T) {
	t.Run("should return true from google login", func(t *testing.T) {

		googleResponse := GoogleOAuth2Response{
			Email: "wuriyanto@bhinneka.com",
			HD:    "bhinneka.com",
		}

		assert.True(t, googleResponse.IsBhinnekaEmail())
	})

	t.Run("should return false from google login", func(t *testing.T) {

		googleResponse := GoogleOAuth2Response{
			Email: "wuriyanto@gmail.com",
			HD:    "google.com",
		}

		assert.False(t, googleResponse.IsBhinnekaEmail())
	})

	t.Run("should return true from azure login", func(t *testing.T) {

		azureResponse := AzureResponse{
			Email:    "wuriyanto@bhinneka.com",
			JobTitle: "Backend Dev",
		}

		assert.True(t, azureResponse.IsBhinnekaEmail())
	})

	t.Run("should return false from azure login", func(t *testing.T) {

		azureResponse := AzureResponse{
			Email:    "wuriyanto@gmail.com",
			JobTitle: "Backend Dev",
		}

		assert.False(t, azureResponse.IsBhinnekaEmail())
	})

	t.Run("should return false from apple login", func(t *testing.T) {

		appleResponse := AppleProfile{
			Email: "test@gmail",
		}

		assert.False(t, appleResponse.IsBhinnekaEmail())
	})
}

func TestNewRefreshTokenModel(t *testing.T) {
	randString := NewRefreshToken("1", "788889", 1)
	refreshToken := &RefreshToken{ID: "1", Token: "788889", RefreshTokenAge: 1}
	assert.EqualValues(t, refreshToken, randString)
}
