package helper

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	log "github.com/sirupsen/logrus"
)

var testLogDatas = []struct {
	level     log.Level
	wantError bool
}{
	{
		log.DebugLevel,
		false,
	},
	{
		log.InfoLevel,
		false,
	},
	{
		log.WarnLevel,
		false,
	},
	{
		log.PanicLevel,
		true,
	},
}

var testLogDataV2 = []struct {
	param     *ParamLog
	wantError bool
}{
	{
		&ParamLog{Level: log.DebugLevel},
		false,
	},
	{
		&ParamLog{Level: log.InfoLevel},
		false,
	},
	{
		&ParamLog{Level: log.WarnLevel},
		false,
	},
	{
		&ParamLog{Level: log.ErrorLevel},
		false,
	},
	{
		&ParamLog{Level: log.PanicLevel},
		true,
	},
}

func TestLog(t *testing.T) {
	for _, tt := range testLogDatas {

		if tt.wantError {
			assert.Panicsf(t, func() {
				Log(tt.level, "none", "no-context", "no-scope")
			}, "none", nil)
		} else {
			assert.NotPanics(t, func() {
				Log(tt.level, "none", "no-context", "no-scope")
			}, "none", nil)
		}
	}
}

func TestLogV2(t *testing.T) {
	for _, tt := range testLogDataV2 {

		if tt.wantError {
			assert.Panicsf(t, func() {
				LogV2(tt.param)
			}, "nope", nil)
		} else {
			assert.NotPanics(t, func() {
				LogV2(tt.param)
			}, "nope", nil)
		}
	}
}

var testMessageData = []struct {
	message     string
	ignoreError bool
}{
	{
		"Invalid username or password",
		true,
	},
	{
		"password contains at least 1 capital letter, 1 lowercase, 1 numeric and 1 special character",
		true,
	},
	{
		"invalid token",
		true,
	},
	{
		"failed to login",
		true,
	},
	{
		"result is not login attempt",
		true,
	},
	{
		"no rows in result set",
		true,
	},
	{
		"Masih terdapat permintaan kirim ulang email aktifasi yang sebelumnya",
		true,
	},
	{
		"this token has been expired",
		true,
	},
	{
		"another error message",
		false,
	},
}

func TestMessageLogV2(t *testing.T) {
	for _, tt := range testMessageData {
		err := errors.New(tt.message)
		assert.Equal(t, tt.ignoreError, IgnoredError(err))
	}
}
