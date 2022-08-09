package helper

import (
	"context"
	"crypto/md5"
	"crypto/sha1"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/golib/jsonschema"
	validation "github.com/Bhinneka/golib/string"
	"github.com/Bhinneka/golib/tracer"
	"github.com/labstack/echo"
	"github.com/xeipuuv/gojsonschema"
)

const (
	// ErrorDataNotFound error message when data doesn't exist
	ErrorDataNotFound = "data %s not found"
	// ErrorParameterInvalid error message for parameter is invalid
	ErrorParameterInvalid = "%s parameter is invalid"
	// ErrorParameterRequired error message for parameter is missing
	ErrorParameterRequired = "%s parameter is required"
	// ErrorParameterLength error message for parameter length is invalid
	ErrorParameterLength = "length of %s parameter exceeds the limit %d"
	// ErrorUnauthorized error message for unauthorized user
	ErrorUnauthorized = "you are not authorized"
	// ErrorOccured error message for internal server error occurred
	ErrorOccured = "error occurred"
	// ErrorResultNotProper error message for response not proper
	ErrorResultNotProper = "result is not proper response"
	//ErrorStatusCode text
	ErrorStatusCode = "error status code"
	// ErrorRedis redis error
	ErrorRedis = "redis: nil"
	// ErrorLoginAttempt message for reaching login attempt limit
	ErrorLoginAttempt = "already reach login attempt limit"
	//ErrorDecode error message decode
	ErrorDecode = "failed decode response"
	// ErrorPayload error message payload
	ErrorPayload = "malformed payload"

	// SuccessMessage message for success process
	SuccessMessage = "succeed to process data"

	// FormatDateDB for inserting date to database
	FormatDateDB = "2006-01-02"
	//FormatDOB for parsing dob from form value
	FormatDOB = "02/01/2006"
	// DefaultDOB default value
	DefaultDOB = "01/01/0001"
	//HTTPRequest text
	HTTPRequest = "http_request"

	//EventProduceCreateMerchant text
	EventProduceCreateMerchant = "merchantCreated"
	//EventProduceUpdateMerchant text
	EventProduceUpdateMerchant = "merchantUpdated"
	//EventProduceDeleteMerchant text
	EventProduceDeleteMerchant = "merchantDeleted"

	// ScopeParseResponse scope parameter
	ScopeParseResponse = "parse_response"
	//ScopeSaveMember scope parameter
	ScopeSaveMember = "save_member"
	//ScopeJSONApiGenerate scope parameter
	ScopeJSONApiGenerate = "jsonapi_generate_message"

	//TextHeader text
	TextHeader = "headers"
	//TextAuthorization text string
	TextAuthorization = "Authorization"
	//TextQuery text
	TextQuery = "query"
	//TextResponse text
	TextResponse = "response"
	//TextStmtError text
	TextStmtError = "stmt_err"
	//TextExecQuery text
	TextExecQuery = "exec_query"
	//TextExecUsecase text
	TextExecUsecase = "exec_usecase"
	//TextPrepareDatabase text
	TextPrepareDatabase = "prepare_database"
	//TextQueryDatabase text
	TextQueryDatabase = "query_database"
	// TextDBBegin for transactional
	TextDBBegin = "db_begin"
	//TextParameter text
	TextParameter = "parameters"
	//TextConnection text
	TextConnection = "connection"
	//TextMemberID text
	TextMemberID = "memberID"
	//TextMerchantID text
	TextMerchantID = "merchantID"
	//TextEmail text
	TextEmail = "email"
	//TextMemberIDCamel text
	TextMemberIDCamel = "memberId"
	//TextMerchantIDCamel text
	TextMerchantIDCamel = "merchantId"
	//TextMerchantVanity
	TextMerchantVanity = "vanityUrl"
	//TextPhoneArea text
	TextPhoneArea = "phone_area"
	//TextFindServerKafkaConfig text
	TextFindServerKafkaConfig = "find_kafka_server_config"
	//TextCheckEnvirontmentVar text
	TextCheckEnvirontmentVar = "check_environment_var"
	//TextArgs text
	TextArgs = "args"
	//TextUpdate text
	TextUpdate = "update"
	TextCreate = "create"
	TextDelete = "delete"
	// TextAdd text
	TextAdd = "add"
	//TextDeleteUpper text
	TextDeleteUpper = "DELETE"
	//TextUpdateUpper text
	TextUpdateUpper = "UPDATE"
	//TextInsertUpper text
	TextInsertUpper = "INSERT"
	//TextMe text
	TextMe = "me"
	//TextURL text
	TextURL = "url"
	//TextToken text
	TextToken = "token"
	//TextHTTPS text
	TextHTTPS = "https://"
	//TextNpwp text
	TextNpwp = "NPWP"
	//TextKTP text
	TextKTP = "KTP"
	//TextTestCase1 text
	TextTestCase1 = "Testcase #1"
	//TextTestCase2 text
	TextTestCase2 = "Testcase #2"
	//TextTestCase3 text
	TextTestCase3 = "Testcase #3"
	//TextTestCase4 text
	TextTestCase4 = "Testcase #4"
	//TextTestCase5 text
	TextTestCase5 = "Testcase #5"
	// TextBearer header
	TextBearer = "Bearer"

	//TagError for error reporting
	TagError = "error"
	// TagCtx tag for error erporting
	TagCtx = "ctx"

	//HeaderClientID header for client login
	HeaderClientID = "X-Client-ID"
	//HeaderClientSecret client secret
	HeaderClientSecret = "X-Client-Secret"
	// TextActive mark value as active
	TextActive = "ACTIVE"
	// TextInactive mark value as inactive
	TextInactive = "INACTIVE"
	// TextRevoke mark value as Revoke
	TextRevoked = "REVOKED"
	// TextInvited mark value as Invited
	TextInvited = "INVITED"
	// TextNew mark value as new
	TextNew = "NEW"
	// Version1 first release
	Version1 = "v1"
	// Version2 second release
	Version2 = "v2"
	// Version3 third release
	Version3 = "v3"
	// Version4 fourth release
	Version4 = "v4"
	// EnvProd production
	EnvProd = "PROD"
	// EnvStaging staging
	EnvStaging = "STAGING"
	// EnvDev development
	EnvDev = "DEV"
	// DefaultFirstName used as default firstname
	DefaultFirstName = "Bhinneka"
	// DefaultLastName used as default lastname
	DefaultLastName = "User"
	// TagKey specific for redis
	TagKey = "key"
	// TextAccount account
	TextAccount = "account"
	// TextNarwhal narwhal
	TextNarwhal = "narwhal"
	// DefaultLKPPName specific for LKPP
	DefaultLKPPName = "Bhinneka LKPP"
	// status
	TextStatus = "status"

	ErrMsgEmptyContent      = "empty email template content"
	ErrMsgFailedGetTemplate = "failed to get email template"
	ErrMsgFailedSendEmail   = "failed to send email"

	FormatYmdhisz = "060102150405.000"
	FormatYmdhis  = "060102150405"
)

// GetHTTPNewRequest function for getting json API
func GetHTTPNewRequest(ctxReq context.Context, method, url string, body io.Reader, target interface{}, headers ...map[string]string) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	ctxReq, cancel := context.WithTimeout(ctxReq, time.Second*10)
	defer cancel()

	var respStatus string
	trace := tracer.StartTrace(ctxReq, fmt.Sprintf("%s %s%s", method, req.URL.Host, req.URL.Path))
	trace.InjectHTTPHeader(req)

	defer func() {
		tags := map[string]interface{}{
			"http.headers":    req.Header,
			"http.method":     req.Method,
			"http.url":        req.URL.String(),
			"response.status": respStatus,
			"response.body":   target,
		}
		trace.Finish(tags)
	}()

	// iterate optional data of headers
	for _, header := range headers {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Do(req)

	if err != nil {
		tracer.Log(ctxReq, HTTPRequest, err)
		return err
	}

	respStatus = r.Status
	defer r.Body.Close()

	e := json.NewDecoder(r.Body).Decode(target)
	if e != nil {
		tracer.Log(ctxReq, HTTPRequest, err)
		return e
	}

	return nil
}

// GetHTTPNewRequestV2 function for getting API
func GetHTTPNewRequestV2(ctxReq context.Context, method, url string, body io.Reader, target interface{}, headers ...map[string]string) error {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}
	ctxReq, cancel := context.WithTimeout(ctxReq, time.Second*10)
	defer cancel()

	var respStatus string

	trace := tracer.StartTrace(ctxReq, fmt.Sprintf("%s %s%s", method, req.URL.Host, req.URL.Path))
	trace.InjectHTTPHeader(req)

	defer func() {
		tags := map[string]interface{}{
			"http.headers":    req.Header,
			"http.method":     req.Method,
			"http.url":        req.URL.String(),
			"response.status": respStatus,
			"response.body":   target,
		}
		trace.Finish(tags)
	}()

	// iterate optional data of headers
	for _, header := range headers {
		for key, value := range header {
			req.Header.Set(key, value)
		}
	}

	client := &http.Client{Timeout: 10 * time.Second}
	r, err := client.Do(req)

	if err != nil {
		tracer.Log(ctxReq, HTTPRequest, err)
		return err
	}

	respStatus = r.Status
	defer r.Body.Close()

	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		tracer.Log(ctxReq, HTTPRequest, err)
		return err
	}

	e := json.Unmarshal(b, &target)
	if e != nil {
		tracer.Log(ctxReq, HTTPRequest, err)
		return e
	}

	if r.StatusCode != 200 {
		err := fmt.Errorf("error code %d with message: %s", r.StatusCode, string(b))
		tracer.Log(ctxReq, HTTPRequest, err)
		return err
	}

	return nil
}

// GenerateRandomID function for generating shipping ID
func GenerateRandomID(length int, prefix ...string) string {
	var strPrefix string

	if len(prefix) > 0 {
		strPrefix = prefix[0]
	}

	yearNow, monthNow, _ := time.Now().Date()
	year := strconv.Itoa(yearNow)[2:len(strconv.Itoa(yearNow))]
	month := int(monthNow)
	RandomString := golib.RandomString(length)

	id := fmt.Sprintf("%s%s%d%s", strPrefix, year, month, RandomString)
	return id
}

// RandomStringBase64 function for random string and base64 encoded
func RandomStringBase64(length int) string {
	rb := make([]byte, length)
	_, err := rand.Read(rb)

	if err != nil {
		return ""
	}
	rs := base64.URLEncoding.EncodeToString(rb)

	reg, _ := regexp.Compile("[^A-Za-z0-9]+")

	return reg.ReplaceAllString(rs, "")
}

// ClearHTML function for validating HTML
func ClearHTML(src string) string {
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	return re.ReplaceAllString(src, "")
}

// StringInSlice function for checking whether string in slice
// str string searched string
// list []string slice
func StringInSlice(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}

// SetLastName function for setting last name
func SetLastName(names []string) string {
	if len(names) <= 1 {
		return ""
	}

	names = append(names[:0], names[1:]...)

	return strings.Join(names, " ")
}

// GenerateMemberID function for generating member id
// based on member max id
func GenerateMemberID(lastUserID string) string {
	preffix := "USR"
	var strNextID string

	intLastUserID, _ := strconv.Atoi(lastUserID)

	if intLastUserID < 1000 {
		strNextID = "0" + strconv.Itoa(intLastUserID+1)
	} else if intLastUserID == 99999999 {
		strNextID = "00000001"
	} else {
		strNextID = strconv.Itoa(intLastUserID + 1)
	}

	now := time.Now()
	dt := now.Format("06") + now.Format("01")

	return preffix + dt + strNextID
}

// GenerateTokenByString function for generate token by string
func GenerateTokenByString(str string) string {
	var token string

	// generate random string
	RandomString := golib.RandomString(30)

	mix := str + "-" + RandomString
	md := md5.New()
	io.WriteString(md, mix)
	sumMd5 := md.Sum(nil)
	hash := hex.EncodeToString(sumMd5[:])

	sh := sha1.New()
	io.WriteString(sh, hash)
	sumSha := sh.Sum(nil)
	token = hex.EncodeToString(sumSha[:])

	return token
}

var camel = regexp.MustCompile("(^[^A-Z]*|[A-Z]*)([A-Z][^A-Z]+|$)")

// CamelToLowerCase function for convert camelcase to lower
func CamelToLowerCase(s string) string {
	var a []string
	for _, sub := range camel.FindAllStringSubmatch(s, -1) {
		if sub[1] != "" {
			a = append(a, sub[1])
		}
		if sub[2] != "" {
			a = append(a, sub[2])
		}
	}
	return strings.ToLower(strings.Join(a, " "))
}

// FloatToString function for convert float to string type
func FloatToString(s string) string {
	f, _ := strconv.ParseFloat(s, 64)
	return strconv.Itoa(int(f))
}

// ConvertTimeToDate function for convert time to date type
func ConvertTimeToDate(sourceTime time.Time) (time.Time, error) {

	result, err := time.Parse(FormatDateDB, sourceTime.Format(FormatDateDB))
	if err != nil {
		return time.Time{}, err
	}
	return result, nil
}

// XNOR function
func XNOR(a, b bool) bool {
	return !((a || b) && (!a || !b))
}

// RandomPassword function for set password random
func RandomPassword(s string, length int) string {
	rand.Seed(time.Now().UnixNano())
	letter := []rune(s)
	b := make([]rune, length)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

// ValidateDocumentFileURL validate function for merchant document url
func ValidateDocumentFileURL(url string) bool {

	AWSMerchantDocumentURL, _ := os.LookupEnv("AWS_MERCHANT_DOCUMENT_URL")
	AWSMerchantDocumentURLSalmon, _ := os.LookupEnv("AWS_MERCHANT_DOCUMENT_URL_SALMON")
	StaticUrl, _ := os.LookupEnv("STATIC_PROFILE_URL")

	contains := strings.Replace(AWSMerchantDocumentURL, TextHTTPS, "", -1)
	containsSalmon := strings.Replace(AWSMerchantDocumentURLSalmon, TextHTTPS, "", -1)
	containsStatic := strings.Replace(StaticUrl, TextHTTPS, "", -1)
	if !strings.Contains(url, contains) && !strings.Contains(url, containsSalmon) && !strings.Contains(url, containsStatic) {
		return false
	}
	return true
}

// ValidateSQLNullString function for validate string null
func ValidateSQLNullString(field sql.NullString) string {
	if field.Valid {
		return field.String
	}
	return ""
}

// ValidateSQLNullInt64 function for validate int64 null
func ValidateSQLNullInt64(field sql.NullInt64) int64 {
	if field.Valid {
		return field.Int64
	}
	return 0
}

// ValidateStringNull function for validate string null
func ValidateStringNull(field *string) string {
	if field == nil {
		return ""
	}
	return *field
}

// ValidateStringToSQLNullString function for validate string to set sqlnull
func ValidateStringToSQLNullString(field string) sql.NullString {
	var nullString sql.NullString
	nullString.Valid = false
	if len(field) > 0 {
		nullString.Valid = true
		nullString.String = field
		return nullString
	}
	return nullString
}

// SplitStreetAddress function for split street1 and street2
func SplitStreetAddress(address string) (string, string) {
	streets := strings.Split(address, "\n")
	if len(streets) > 1 {
		street1 := streets[0]

		// delete index 0 and join the rest
		streets = append(streets[:0], streets[0+1:]...)
		street2 := strings.TrimSpace(strings.Join(streets, " "))

		return street1, street2
	}
	return address, ""
}

// ExtractClientCred get client credential
func ExtractClientCred(c echo.Context) (clientID, clientSecret string) {
	h := c.Request().Header
	clientID = h.Get(HeaderClientID)
	clientSecret = h.Get(HeaderClientSecret)
	return clientID, clientSecret
}

// SplitStringByN return slice
// split string by specific length
func SplitStringByN(str string, size int) []string {
	stringLength := len(str)
	splitedLength := int(math.Ceil(float64(stringLength) / float64(size)))
	splited := make([]string, splitedLength)
	var start, stop int
	for i := 0; i < splitedLength; i++ {
		start = i * size
		stop = start + size
		if stop > stringLength {
			stop = stringLength
		}
		splited[i] = str[start:stop]
	}
	return splited
}

func validateAppleName(firstName, lastName string) (string, string) {
	var (
		validFNChar, validFNLen, validLNChar, validLNLen bool
	)

	if firstName != "" {
		if golib.ValidateAlphabetWithSpace(firstName) {
			validFNChar = true
		}

		if err := golib.ValidateMaxInput(firstName, 25); err == nil {
			validFNLen = true
		}
	}

	if len(lastName) != 0 {
		if golib.ValidateAlphabetWithSpace(lastName) {
			validLNChar = true
		}

		if err := golib.ValidateMaxInput(lastName, 25); err == nil {
			validLNLen = true
		}
	}
	if !validFNChar || !validFNLen {
		firstName = DefaultFirstName
	}
	if !validLNChar || !validLNLen {
		lastName = DefaultLastName
	}
	return firstName, lastName
}

// ValidateTemp from Go data type for response single error
func ValidateTemp(schemaID string, input interface{}) error {

	schema, err := jsonschema.Get(schemaID)

	if err != nil {
		return err
	}

	document := gojsonschema.NewGoLoader(input)
	return validateTemp(schema, document)
}

// ValidateTemp from Go data type for response single error
func validateTemp(schema *gojsonschema.Schema, document gojsonschema.JSONLoader) error {

	result, err := schema.Validate(document)
	if err != nil {
		return err
	}

	if !result.Valid() {
		var message string
		for _, desc := range result.Errors() {
			field := golib.CamelToLowerCase(desc.Field())
			message = field + " " + golib.CamelToLowerCase(desc.Description())
		}
		return errors.New(message)
	}

	return nil
}

// ValidateLatLong ...
func ValidateLatLong(latitude float64, longitude float64) bool {
	if latitude != 0 && longitude != 0 {
		return true
	}
	return false
}

// ValidateMamberID ...
func ValidateMamberID(params string) bool {
	err := validation.ValidateAlphanumeric(params, false)
	if err && strings.HasPrefix(params, "USR") {
		return true
	}
	return false
}

func ValidationMerchantName(str string) bool {
	var uppercase, lowercase, num, allowed, anothersymbol int
	for _, r := range str {
		if IsUppercase(r) {
			uppercase = +1
		} else if IsLowercase(r) {
			lowercase = +1
		} else if IsNumeric(r) {
			num = +1
		} else if IsAllowedSymbol(r) {
			allowed = +1
		} else {
			anothersymbol = +1
		}
	}
	if anothersymbol > 0 {
		return false
	}

	return uppercase >= 1 || lowercase >= 1 || num >= 1 || allowed >= 0
}

func IsUppercase(r rune) bool {
	return int(r) >= 65 && int(r) <= 90
}

// IsLowercase reusable rune check if char is lowercase
func IsLowercase(r rune) bool {
	return int(r) >= 97 && int(r) <= 122
}

// IsNumeric reusable rune check if char is numeric
func IsNumeric(r rune) bool {
	return int(r) >= 48 && int(r) <= 57
}

// IsAllowedSymbol check if rune is any of
func IsAllowedSymbol(r rune) bool {
	m := int(r)
	return m == 32 || m == 45 || m == 38 || m == 46 || m == 44 || m == 40 || m == 41
}

// RawQuery ...
func RawQuery(query string, offset, limit int, order, sort string) string {
	if order != "" {
		query += ` ORDER BY ` + order + ` ` + sort
	}
	if limit > 0 {
		query += ` LIMIT ` + strconv.Itoa(limit) + ` OFFSET ` + strconv.Itoa(offset)
	}
	return query
}

func TrimSpace(text string) string {
	trimmedText := strings.TrimSpace(text)
	return trimmedText
}
