package service

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/Bhinneka/bhinneka-go-sdk"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

const (
	constUploadAccessor       = "UPLOAD_SERVICE_ACCESSOR_ID"
	constUploadURL            = "UPLOAD_SERVICE_URL"
	constUploadTimeout        = "UPLOAD_SERVICE_TIMEOUT"
	uploadURL                 = "http://upload.bhinnekatesting.com"
	accessorID                = "somenumber"
	goodToken                 = "bearer token"
	constAwsMerchantDocSalmon = "AWS_MERCHANT_DOCUMENT_URL_SALMON"
	constAwsMerchantDoc       = "AWS_MERCHANT_DOCUMENT_URL"
	salmonURL                 = "https://s3.ap-southeast-1.amazonaws.com"
	bmdUploadURL              = "https://bmd-upload.s3.ap-southeast-1.amazonaws.com"
)

func TestInitUpload(t *testing.T) {
	_, err := NewUploadService()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "specify UPLOAD_SERVICE_URL")

	os.Setenv(constUploadURL, uploadURL)
	_, err = NewUploadService()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "specify UPLOAD_SERVICE_URL")

	os.Setenv(constUploadAccessor, accessorID)
	_, errm := NewUploadService()
	assert.Error(t, errm)

	os.Setenv(constUploadTimeout, "1")
	_, errm = NewUploadService()
	assert.NoError(t, errm)

	os.Setenv(constUploadURL, badURL)
	_, errm = NewUploadService()
	assert.Error(t, errm)
}

func TestGetURLImage(t *testing.T) {
	os.Setenv(constUploadTimeout, "1")
	os.Setenv(constUploadURL, uploadURL)
	os.Setenv(constUploadAccessor, accessorID)

	var testData = []struct {
		name            string
		wantError       bool
		token           string
		statusCode      int
		serviceResponse interface{}
		keyFile         string
		isAttachment    string
	}{
		{
			name:       "Get Image #1",
			wantError:  true,
			statusCode: http.StatusBadRequest,
		},
		{
			name:         "Get Image #2",
			wantError:    false,
			statusCode:   http.StatusOK,
			token:        goodToken,
			keyFile:      "https://static.bmdstatic.com/file/image.png",
			isAttachment: "false",
		},
		{
			name:         "Get Image #3",
			wantError:    true,
			statusCode:   http.StatusBadRequest,
			token:        goodToken,
			keyFile:      "%%2\r\n",
			isAttachment: "false",
		},
		{
			name:            "Get Image #4",
			wantError:       false,
			statusCode:      http.StatusOK,
			token:           goodToken,
			keyFile:         "image.png",
			serviceResponse: nil,
			isAttachment:    "false",
		},
		{
			name:            "Get Image #5",
			wantError:       true,
			statusCode:      http.StatusBadRequest,
			token:           goodToken,
			keyFile:         "image.png",
			serviceResponse: []byte(``),
			isAttachment:    "false",
		},
	}

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	for _, tc := range testData {
		us, _ := NewUploadService()
		uri := fmt.Sprintf("%s%s%s%s%s%s%s", uploadURL, "/presigned?timeout=", "1", "&key=", tc.keyFile, "&isAttachment=", tc.isAttachment)
		bhinneka.MockHTTP(http.MethodGet, uri, tc.statusCode, tc.serviceResponse)
		ctx := context.WithValue(context.Background(), TokenContextKey, tc.token)

		sr := <-us.GetURLImage(ctx, tc.keyFile, tc.isAttachment)
		if tc.wantError {
			assert.Error(t, sr.Error)
		} else {
			assert.NoError(t, sr.Error)
		}
	}
}

func TestReplaceURLFile(t *testing.T) {
	os.Setenv(constUploadTimeout, "1")
	os.Setenv(constUploadURL, uploadURL)
	os.Setenv(constUploadAccessor, accessorID)

	us, _ := NewUploadService()
	_, valid := us.replaceURLFile("https://s3.ap-southeast-1.amazonaws.com/image.png")
	assert.Equal(t, true, valid)

	os.Setenv(constAwsMerchantDoc, uploadURL)
	us, _ = NewUploadService()
	_, valid = us.replaceURLFile("https://somedomain.com/ll.png")
	assert.Equal(t, true, valid)

	os.Setenv(constAwsMerchantDocSalmon, bmdUploadURL)
	us, _ = NewUploadService()
	_, valid = us.replaceURLFile("https://bmd-upload.s3.ap-southeast-1.amazonaws.com/another-image.png")
	assert.Equal(t, true, valid)
}
