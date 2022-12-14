package middleware

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/labstack/echo"
	log "github.com/sirupsen/logrus"
)

const (
	// ContextKeyClientIP to pass client IP address
	ContextKeyClientIP = "ClientIP"
)

// Logger function for writing all request log into console
func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		name := "user-service"

		log.SetFormatter(&log.JSONFormatter{})
		l := log.StandardLogger()

		start := time.Now()
		req := c.Request()
		res := c.Response()

		remoteAddr := req.RemoteAddr
		if ip := req.Header.Get(echo.HeaderXRealIP); ip != "" {
			remoteAddr = ip
		} else if ip = req.Header.Get(echo.HeaderXForwardedFor); ip != "" {
			remoteAddr = ip
		} else {
			remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
		}
		ctx := context.WithValue(req.Context(), "log-basic-auth", true)
		r := req.WithContext(ctx)
		c.SetRequest(r)

		entry := l.WithFields(log.Fields{
			"request": req.RequestURI,
			"method":  req.Method,
			"remote":  remoteAddr,
		})

		if reqID := req.Header.Get("X-Request-Id"); reqID != "" {
			entry = entry.WithField("request_id", reqID)
		}

		entry.Info("started handling request")

		if err := next(c); err != nil {
			c.Error(err)
		}

		latency := time.Since(start)

		entry.WithFields(log.Fields{
			"size":                           res.Size,
			"status":                         res.Status,
			"text_status":                    http.StatusText(res.Status),
			"took":                           latency,
			fmt.Sprintf("#%s.latency", name): latency.Nanoseconds(),
		}).Info("completed handling request")

		return nil
	}
}
