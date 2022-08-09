package helper

import (
	"os"

	"github.com/Bhinneka/golib"
)

func SendNotification(title, body, ctx string, err error) {
	os.Setenv("SERVER_ENV", os.Getenv("ENV"))
	golib.SendNotification(title, body, ctx, err)
}
