package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSessionInfoModel(t *testing.T) {
	email := "yderana@gmail.com"
	clientType := "ios"
	sessionTableTestCases := []*SessionInfoResponse{
		{UserName: &email},
		{ClientType: &clientType},
	}

	t.Run("Should return true test valid username", func(t *testing.T) {
		m := sessionTableTestCases[0]
		assert.True(t, m.UserName == &email)
	})

	t.Run("should return false valid client type", func(t *testing.T) {
		m := sessionTableTestCases[1]

		assert.True(t, m.ClientType == &clientType)
	})
}
