package helper

import (
	"math/rand"
	"strconv"
	"time"
)

func GenerateMemberIDv2() string {
	now := time.Now()
	digit := ((uint64(now.Unix()) + uint64(now.Nanosecond()/int(time.Millisecond)) + uint64(rand.Int63())) % 10000000000) + uint64(now.Nanosecond())
	return now.Format("USR0601") + strconv.FormatUint(digit, 10)
}

func GenerateDocumentID() string {
	now := time.Now()
	temp := ((uint64(now.Unix())+uint64(now.Nanosecond()/int(time.Millisecond))+uint64(rand.Int63()))%1000000000 + uint64(rand.Int63())) % 100000000000
	return now.Format("DOC0601") + strconv.FormatUint(temp, 10)
}
