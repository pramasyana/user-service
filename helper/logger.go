package helper

import (
	"os"

	"github.com/Bhinneka/golib"
	"github.com/getsentry/raven-go"
	log "github.com/sirupsen/logrus"
)

// ParamLog log param struct
type ParamLog struct {
	Error   error
	Level   log.Level
	Message string
	Context string
	Scope   string
	TraceID string
}

const (
	// TOPIC for setting topic of log
	TOPIC = "user-service-log"
	// LogTag default log tag
	LogTag = "user-service"
)

// IgnoredErrorMessage these error message should not send to sentry
var IgnoredErrorMessage = []string{
	"Invalid username or password",
	"password contains at least 1 capital letter, 1 lowercase, 1 numeric and 1 special character",
	"invalid token",
	"failed to login",
	"result is not login attempt",
	"no rows in result set",
	"Masih terdapat permintaan kirim ulang email aktifasi yang sebelumnya",
	"this token has been expired",
	"invalid old token",
	"cannot register your account, your email is hidden",
	"Akun Anda belum aktif. Periksa email subject Konfirmasi Alamat Email Anda dari Bhinneka untuk aktifasi akun. Atau kirimkan ulang email aktifasi",
	"failed to get token",
	"Email atau password yang Anda masukkan salah",
	"Maaf, Email anda belum terdaftar, silahkan registrasi terlebih dahulu",
	"refresh token is invalid",
	"Akun Anda telah diblokir. Silakan coba kembali setelah 5 menit",
	"value too long for type character varying(20)",
	"sql: no rows in result set",
	"client id is not found",
	"This authorization code has been used.",
	"This authorization code has expired.",
	"redis: nil",
	"invalid_grant",
}

// LogContext function for logging the context of echo
// c string context
// s string scope
func LogContext(c string, s string) *log.Entry {
	return log.WithFields(log.Fields{
		"topic":   TOPIC,
		"context": c,
		"scope":   s,
	})
}

// Log function for returning entry type
// level log.Level
// message string message of log
// context string context of log
// scope string scope of log
func Log(level log.Level, message string, context string, scope string) {
	log.SetFormatter(&log.JSONFormatter{})

	entry := LogContext(context, scope)
	switch level {
	case log.DebugLevel:
		entry.Debug(message)
	case log.InfoLevel:
		entry.Info(message)
	case log.WarnLevel:
		entry.Warn(message)
	case log.ErrorLevel:
		entry.Error(message)
	case log.FatalLevel:
		entry.Fatal(message)
	case log.PanicLevel:
		entry.Panic(message)
	}
}

// LogV2 function for returning entry type
// err error
// level log.Level
// message string message of log
// context string context of log
// scope string scope of log
// DEPRECATED user SendE
func LogV2(params *ParamLog) {
	log.SetFormatter(&log.JSONFormatter{})

	entry := LogContext(params.Context, params.Scope)
	message := params.Message
	switch params.Level {
	case log.DebugLevel:
		entry.Debug(message)
	case log.InfoLevel:
		entry.Info(message)
	case log.WarnLevel:
		entry.Warn(message)
	case log.ErrorLevel:
		entry.Error(message)
	case log.FatalLevel:
		entry.Fatal(message)
	case log.PanicLevel:
		entry.Panic(message)
	}

	if os.Getenv("SENTRY") == "true" && params.Error != nil && !IgnoredError(params.Error) {
		raven.CaptureError(params.Error, map[string]string{params.Context: params.Message, "trace_id": params.TraceID})
	}
}

// IgnoredError ignored error message which will be sent to sentry
func IgnoredError(err error) bool {
	return golib.StringInSlice(err.Error(), IgnoredErrorMessage)
}
