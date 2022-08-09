package helper

import (
	"context"
	"encoding/json"
	"os"
	"strconv"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/tracer"
	"github.com/getsentry/raven-go"
)

// SendErrorLog function for log and send error to sentry
func SendErrorLog(ctxReq context.Context, actionContext string, scope string, err error, payload interface{}) {
	if err == nil {
		return
	}
	golib.Log(golib.ErrorLevel, err.Error(), actionContext, scope)

	isSentryActive, _ := strconv.ParseBool(os.Getenv("SENTRY"))
	if err.Error() == ErrorRedis {
		tracer.SetError(ctxReq, err)
	}
	if isSentryActive && !IgnoredError(err) {
		js, _ := json.Marshal(payload)
		sentryPayload := map[string]string{
			"ctx":      actionContext,
			"trace_id": tracer.GetTraceID(ctxReq),
			"error":    err.Error(),
		}
		if js != nil {
			sentryPayload["payload"] = string(js)
		}
		if clientName, ok := ctxReq.Value("clientName").(string); ok {
			sentryPayload["client"] = clientName
		}
		tracer.SetError(ctxReq, err)
		raven.CaptureError(err, sentryPayload)
	}

}
