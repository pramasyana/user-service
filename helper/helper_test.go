package helper

import (
	"context"
	"database/sql"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strings"
	"testing"
	"time"

	bhinneka "github.com/Bhinneka/bhinneka-go-sdk"
	"github.com/Bhinneka/golib"
	"github.com/jarcoal/httpmock"
	"github.com/labstack/echo"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

const (
	testCase1 = "Testcase #1"
	testCase2 = "Testcase #2"
	testCase3 = "Testcase #3"
	urlValid1 = "https://s3.ap-southeast-1.amazonaws.com/static.bmdstatic.com/sf/merchant_images/photoqu-1574132451.jpeg"
	urlValid2 = "https://bmd-upload.s3.ap-southeast-1.amazonaws.com/static.bmdstatic.com/sf/merchant_images/photoqu-1574132451.jpeg"
	awsHost1  = "https://bmd-upload.s3.ap-southeast-1.amazonaws.com"
	awsHost2  = "https://s3.ap-southeast-1.amazonaws.com"
	awsHost3  = "https://static.bmdstatic.com"
)

// PlanktonVariant data structure
type PlanktonVariant struct {
	Data []struct {
		Type       string `json:"type"`
		ID         string `json:"id"`
		Attributes Attr   `json:"attributes"`
	} `json:"data"`
	Included []struct {
		OffersProducts
	} `json:"included"`
}

// Attr data structure
type Attr struct {
	Name          string    `json:"name"`
	Type          []string  `json:"type"`
	SKU           string    `json:"skuNo"`
	RegisterNo    *string   `json:"registerNo"`
	Status        string    `json:"status"`
	Note          *string   `json:"note"`
	Cogs          int       `json:"cogs"`
	VideoURL      *[]string `json:"videoUrl"`
	AttachmentURL *[]string `json:"attachmentUrl"`
	Image         []Images  `json:"images"`
	Stock         []Stocks  `json:"stock"`
	Creator       Creators  `json:"creator"`
	CreatedAt     string    `json:"createdAt"`
	Editor        Editors   `json:"editor"`
	ModifiedAt    string    `json:"modifiedAt"`
}

// Images data structure
type Images struct {
	Default bool        `json:"default"`
	Order   int         `json:"order"`
	Variety []Varieties `json:"variety"`
}

// Varieties data structure
type Varieties struct {
	Type string `json:"type"`
	Size string `json:"size"`
	URL  string `json:"url"`
}

// Stocks data structure
type Stocks struct {
	LocationCode string `json:"locationCode"`
	Name         string `json:"name"`
	OnHand       int    `json:"onHand"`
	OnReserve    int    `json:"onReserve"`
	Available    int    `json:"available"`
}

// Creators data structure
type Creators struct {
	Email    string `json:"email"`
	ID       string `json:"id"`
	IP       string `json:"ip"`
	JobTitle string `json:"jobTitle"`
	Name     string `json:"name"`
}

// Editors data structure
type Editors struct {
	Email    string `json:"email"`
	ID       string `json:"id"`
	IP       string `json:"ip"`
	JobTitle string `json:"jobTitle"`
	Name     string `json:"name"`
}

// OffersProducts data structure
type OffersProducts struct {
	Type       string `json:"type"`
	ID         string `json:"id"`
	Attributes Attrs  `json:"attributes"`
}

// Attrs data structure
type Attrs struct {
	SkuNo                 string               `json:"skuNo,omitempty"`
	VendorSku             *string              `json:"vendorSku,omitempty"`
	MerchantID            *string              `json:"merchantId,omitempty"`
	ConditionID           int                  `json:"conditionId,omitempty"`
	ConditionName         string               `json:"conditionName,omitempty"`
	ConditionNote         *string              `json:"conditionNote,omitempty"`
	ShippingLength        *string              `json:"shippingLength,omitempty"`
	ShippingWidth         *string              `json:"shippingWidth,omitempty"`
	ShippingHeight        *string              `json:"shippingHeight,omitempty"`
	ShippingWeight        float32              `json:"shippingWeight,omitempty"`
	VplPrice              int                  `json:"vplPrice,omitempty"`
	VplSuggestedPrice     int                  `json:"vplSuggestedPrice,omitempty"`
	NormalPrice           int                  `json:"normalPrice,omitempty"`
	SpecialPrice          int                  `json:"specialPrice,omitempty"`
	SpecialPriceStartDate string               `json:"specialPriceStartDate,omitempty"`
	SpecialPriceEndDate   string               `json:"specialPriceEndDate,omitempty"`
	WarrantyTypeName      *string              `json:"warrantyTypeName,omitempty"`
	WarrantyPeriodID      int                  `json:"warrantyPeriodId,omitempty"`
	WarrantyPeriodName    *string              `json:"warrantyPeriodName,omitempty"`
	WarrantyNote          *string              `json:"warrantyNote,omitempty"`
	OfferInfo             *string              `json:"offerInfo,omitempty"`
	OfferStatus           string               `json:"offerStatus,omitempty"`
	HandlingTime          int                  `json:"handlingTime,omitempty"`
	ShippingNote          string               `json:"shippingNote,omitempty"`
	Name                  string               `json:"name,omitempty"`
	Model                 string               `json:"model,omitempty"`
	Status                string               `json:"status,omitempty"`
	RegisterNoTypeID      *int                 `json:"registerNoTypeId,omitempty"`
	RegisterNoTypeName    *string              `json:"registerNoTypeName,omitempty"`
	ExemptionTypeID       int                  `json:"exemptionTypeId,omitempty"`
	ExemptionTypeName     *string              `json:"exemptionTypeName,omitempty"`
	CategoryID            string               `json:"categoryId,omitempty"`
	CategoryName          *string              `json:"categoryName,omitempty"`
	CategoryStructure     []CategoryStructures `json:"categoryStructure,omitempty"`
	BrandID               string               `json:"brandId,omitempty"`
	BrandName             string               `json:"brandName,omitempty"`
	Description           []string             `json:"description,omitempty"`
	InTheBox              []string             `json:"inTheBox,omitempty"`
	KeyFeatures           *string              `json:"keyFeatures,omitempty"`
	SearchTerms           *string              `json:"searchTerms,omitempty"`
	IntendedUse           *string              `json:"intendedUse,omitempty"`
	TargetAudience        *string              `json:"targetAudience,omitempty"`
	RealLength            *string              `json:"realLength,omitempty"`
	RealWidth             *string              `json:"realWidth,omitempty"`
	RealHeight            *string              `json:"realHeight,omitempty"`
	RealWeight            float64              `json:"realWeight,omitempty"`
	Specs                 []Specs              `json:"specs,omitempty"`
	VariantTypeID         int                  `json:"variantTypeId,omitempty"`
	VariantTypeName       *string              `json:"variantTypeName,omitempty"`
	VariantType           *string              `json:"variantType,omitempty"`
	RichContent           string               `json:"richContent,omitempty"`
	Creator               Creators             `json:"creator"`
	CreatedAt             string               `json:"createdAt"`
	Editor                Editors              `json:"editor"`
	ModifiedAt            string               `json:"modifiedAt"`
}

// CategoryStructures data structure
type CategoryStructures struct {
	ID           string `json:"id"`
	OldID        string `json:"oldId"`
	Name         string `json:"name"`
	CurrentLevel int    `json:"currentLevel"`
}

// Specs data structure
type Specs struct {
	ID    string `json:"id"`
	OldID int    `json:"oldId"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ProductRatingNReview data structure
type ProductRatingNReview struct {
	Status string `json:"status"`
	Data   struct {
		TotalReview int     `json:"totalReview"`
		TotalRating float32 `json:"totalRating"`
	} `json:"data"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func getPass() string {
	return "da1c25d8-37c8-41b1-afe2-42dd4825bfea"
}

var testDatas = []struct {
	name       string
	method     string
	body       io.Reader
	wantError  bool
	target     interface{}
	statusCode int
	response   interface{}
	url        string
}{
	{
		"#1",
		"GET",
		nil,
		false,
		&PlanktonVariant{},
		http.StatusOK,
		nil,
		"http://api.plankton.bhinneka.com/master/variants?include=offers,product&filter[skuNo]=sku00113580",
	},
	{
		"#2",
		"GET",
		nil,
		true,
		PlanktonVariant{},
		http.StatusBadRequest,
		nil,
		"http://api.plankton.bhinneka.com/master/variants?include=offers,product&filter[skuNo]=sku00113581",
	},
}

var headers = map[string]string{
	"Content-Type": "application/vnd.api+json",
	"Accept":       "application/vnd.plankton_api.v3+json",
	"x-api-key":    "74bfb4277e931989784fe54f04ea3951a8631a60",
}

// TestGetHTTPNewRequestJSON test function
func TestGetHTTPNewRequestJSON(t *testing.T) {
	ctx := context.Background()
	err := GetHTTPNewRequest(ctx, "POST", "", strings.NewReader("something"), nil, headers)
	assert.Error(t, err)

	errM := GetHTTPNewRequest(ctx, "\\", "http://somedomain.com", nil, &PlanktonVariant{}, headers)
	assert.Error(t, errM)

	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	for _, tt := range testDatas {
		bhinneka.MockHTTP(tt.method, tt.url, tt.statusCode, tt.target)
		err := GetHTTPNewRequest(ctx, tt.method, tt.url, tt.body, tt.target, headers)
		if tt.wantError {
			assert.Error(t, err)
		}
	}
}

var testData = []struct {
	name       string
	wantError  bool
	body       io.Reader
	method     string
	url        string
	response   interface{}
	target     interface{}
	statusCode int
}{
	{
		name:       "Test V2 #2",
		method:     http.MethodGet,
		wantError:  false,
		statusCode: http.StatusOK,
		url:        "http://some-url.com",
		target:     JSONSchemaTemplate{},
		response:   `{"code":200}`,
	},
	{
		name:       "Test V2 #3",
		method:     http.MethodGet,
		wantError:  true,
		statusCode: http.StatusBadRequest,
		url:        "http://bhinnekalocal.com",
		target:     JSONSchemaTemplate{},
		response:   `{"code":401, "message":"unauthorized"}`,
	},
}

func TestGetHTTPNewRequestV2(t *testing.T) {
	ctx := context.Background()
	err := GetHTTPNewRequestV2(ctx, http.MethodPost, "", strings.NewReader("some oayload"), nil, headers)
	assert.Error(t, err)

	for _, tc := range testData {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		bhinneka.MockHTTP(tc.method, tc.url, tc.statusCode, tc.response)
		err := GetHTTPNewRequestV2(ctx, tc.method, tc.url, tc.body, tc.target, nil)
		if tc.wantError {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}

var testDataAuthBasic = []struct {
	name       string
	url        string
	user       string
	pass       string
	method     string
	wanterror  bool
	statusCode int
}{
	{
		"#1",
		"http://squid.bhinnekalocal.com:9003/api/product-rating/sku00112569",
		"bhinneka",
		getPass(),
		"GET",
		false,
		200,
	},
	{
		"#2",
		"",
		"bhinneka",
		getPass(),
		"GET",
		false,
		200,
	},
}

// TestClearHTML test function
func TestClearHTML(t *testing.T) {
	str := "<strong>Let us make Bhinneka great again</strong>"

	clearStr := ClearHTML(str)

	assert.Equal(t, "Let us make Bhinneka great again", clearStr)
}

// TestStringInSlice test function
func TestStringInSlice(t *testing.T) {
	strSlice := []string{
		"willy",
		"bern",
		"wurry",
	}

	ok := StringInSlice("willy", strSlice)
	nok := StringInSlice("linda", strSlice)

	assert.True(t, ok)
	assert.False(t, nok)

}

// TestRandomString test function
func TestRandomString(t *testing.T) {
	length := 8
	randString := golib.RandomString(length)

	assert.EqualValues(t, length, len(randString))
}

// TestGenerateRandomID test function
func TestGenerateRandomID(t *testing.T) {
	length := 8
	prefix := "USR"

	randString := GenerateRandomID(length, prefix)

	assert.True(t, strings.Contains(randString, prefix), "random string is false")
}

func TestRandomStringBase64(t *testing.T) {
	length := 8

	randString := RandomStringBase64(length)

	reg := regexp.MustCompile("^[A-Za-z0-9]*$")
	b := reg.Match([]byte(randString))

	assert.True(t, b, "error happens")
}

func TestSetLastName(t *testing.T) {
	type args struct {
		names []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Testcase #1: Positive, return true last name",
			args: args{names: []string{"agung", "dp"}},
			want: "dp",
		},
		{
			name: "Testcase #2: Positive, empty args return empty string",
			args: args{names: []string{}},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SetLastName(tt.args.names); got != tt.want {
				t.Errorf("SetLastName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateMemberID(t *testing.T) {
	type args struct {
		lastUserID string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: testCase1,
			args: args{lastUserID: "1"},
			want: "02",
		},
		{
			name: testCase2,
			args: args{lastUserID: "99999999"},
			want: "01",
		},
		{
			name: testCase3,
			args: args{lastUserID: "10000"},
			want: "01",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateMemberID(tt.args.lastUserID)
			got = got[len(got)-2:]
			if got != tt.want {
				t.Errorf("GenerateMemberID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateTokenByString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: testCase1,
			args: args{str: "12345"},
			want: "5f954acdef99eb0d0d74a5ccd49cf13886d34655",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateTokenByString(tt.args.str)
			assert.Equal(t, len(tt.want), len(got))
		})
	}
}

func TestCamelToLowerCase(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: testCase1,
			args: args{str: "ABcde"},
			want: "a bcde",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CamelToLowerCase(tt.args.str)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFloatToString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: testCase1,
			args: args{str: "12.29"},
			want: "12",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FloatToString(tt.args.str)
			assert.Equal(t, tt.want, got)
		})
	}
}

// TestValidateDocumentFileURL test function
func TestValidateDocumentFileURL(t *testing.T) {
	os.Setenv("AWS_MERCHANT_DOCUMENT_URL", awsHost1)
	os.Setenv("AWS_MERCHANT_DOCUMENT_URL_SALMON", awsHost2)
	os.Setenv("STATIC_PROFILE_URL", awsHost3)
	urlInValid := "asda"
	assert.True(t, ValidateDocumentFileURL(urlValid1), "string is true")
	assert.True(t, ValidateDocumentFileURL(urlValid2), "string is true")
	assert.False(t, ValidateDocumentFileURL(urlInValid), "string is false")

}

func TestDateTime(t *testing.T) {
	now := time.Now()
	nowDate, err := ConvertTimeToDate(now)
	assert.NoError(t, err)
	assert.NotNil(t, nowDate)
}

func TestXNOR(t *testing.T) {
	err := XNOR(false, false)
	assert.NotNil(t, err)
}

func TestSendErrorLog(*testing.T) {
	SendErrorLog(context.Background(), TextParameter, TextParameter, sql.ErrNoRows, nil)
}

func TestSqlNullString(t *testing.T) {
	var field sql.NullString
	fieldString := ValidateSQLNullString(field)
	assert.Equal(t, "", fieldString)

	field.Valid = true
	fieldString = ValidateSQLNullString(field)
	assert.Equal(t, "", fieldString)
}

func TestSqlNullInt64(t *testing.T) {
	var field sql.NullInt64
	fieldString := ValidateSQLNullInt64(field)
	assert.Equal(t, cast.ToInt64("0"), fieldString)

	field.Valid = true
	fieldString = ValidateSQLNullInt64(field)
	assert.Equal(t, cast.ToInt64("0"), fieldString)
}

func TestSplitStreetAddress(t *testing.T) {
	address1, address2 := SplitStreetAddress("Testing Split \n Dua Address")
	assert.Equal(t, "Testing Split ", address1)
	assert.Equal(t, "Dua Address", address2)

	address1, address2 = SplitStreetAddress("Testing one")
	assert.Equal(t, "Testing one", address1)
	assert.Equal(t, "", address2)

}

func TestStringToSQLNullString(t *testing.T) {
	fieldString := "ada"
	fieldNull := ValidateStringToSQLNullString(fieldString)
	assert.Equal(t, "ada", fieldNull.String)

	var field sql.NullString
	fieldString = ""
	fieldNull = ValidateStringToSQLNullString(fieldString)
	assert.Equal(t, field, fieldNull)
}

func TestGetHeader(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set("X-Client-ID", "someToken")
	req.Header.Set("X-Client-Secret", "someSecret")
	rec := httptest.NewRecorder()
	ctx := e.NewContext(req, rec)

	cID, cSecret := ExtractClientCred(ctx)
	assert.NotEqual(t, "", cID)
	assert.NotEqual(t, "", cSecret)
}

func TestSplitString(t *testing.T) {
	t.Run("Should return slice", func(t *testing.T) {
		input := "MyStringLengthMoreThanTwenty"
		splitted := SplitStringByN(input, 20)
		assert.Equal(t, 2, len(splitted))

		second := "If less than 20"
		splitted = SplitStringByN(second, 20)
		assert.Equal(t, 1, len(splitted))
	})
}

func TestValidateApplePArameter(t *testing.T) {

	var testData = []struct {
		firstname     string
		lastname      string
		expectedFname string
		expectedLname string
	}{
		{
			firstname:     "",
			lastname:      "string",
			expectedFname: DefaultFirstName,
			expectedLname: "string",
		},
		{
			firstname:     "Awan",
			lastname:      "string-kosong",
			expectedFname: "Awan",
			expectedLname: DefaultLastName,
		},
		{
			firstname:     "Awan12",
			lastname:      "string'ss",
			expectedFname: DefaultFirstName,
			expectedLname: DefaultLastName,
		},
		{
			firstname:     "My Name",
			lastname:      "string",
			expectedFname: "My Name",
			expectedLname: "string",
		},
	}
	t.Run("Test validating", func(t *testing.T) {
		for _, tc := range testData {
			f, l := validateAppleName(tc.firstname, tc.lastname)
			assert.Equal(t, tc.expectedFname, f)
			assert.Equal(t, tc.expectedLname, l)
		}
	})
}

func TestValidateTemp(t *testing.T) {
	t.Run("Should return error", func(t *testing.T) {
		schemaId := "add_merchant_params_v4"
		var payload interface{}
		validSchema := ValidateTemp(schemaId, payload)
		assert.Equal(t, "schema 'add_merchant_params_v4' not found", validSchema.Error())
	})
}
