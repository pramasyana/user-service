package model

import (
	"testing"

	"github.com/Bhinneka/user-service/src/member/v1/model"
	"github.com/stretchr/testify/assert"
)

func TestStatus(t *testing.T) {
	randString := StringToStatus(model.InactiveString)
	assert.EqualValues(t, 0, randString)

	getString := randString.String()
	assert.EqualValues(t, inactive, getString)

	randString = StringToStatus(model.ActiveString)
	assert.EqualValues(t, 1, randString)

	getString = randString.String()
	assert.EqualValues(t, active, getString)

	randString = StringToStatus(model.BlockedString)
	assert.EqualValues(t, 2, randString)

	getString = randString.String()
	assert.EqualValues(t, blocked, getString)

	randString = StringToStatus("UNDIFINED")
	assert.EqualValues(t, 1, randString)

	randString = 1
	getString = randString.String()
	assert.EqualValues(t, active, getString)
}

func TestNewClientApp(t *testing.T) {
	randString := NewClientApp("ABC")
	assert.EqualValues(t, randString, randString)
}

func TestAuthenticate(t *testing.T) {
	clientApp := NewClientApp("ABC")
	randString := clientApp.Authenticate("ABC")
	assert.EqualValues(t, false, randString)
}
