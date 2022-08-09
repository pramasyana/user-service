package service

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/Bhinneka/user-service/helper"
	serviceModel "github.com/Bhinneka/user-service/src/service/model"
)

// DolphinService data structure for dolphin receiver
type DolphinService struct {
	BaseURL   *url.URL
	BasicAuth string
}

const (
	// DolphinDataNotFound message
	DolphinDataNotFound   = "Data Not Found"
	contextDolphinService = "NewDolphinService"
	urlSyncDolphin        = "/index.php?entryPoint=syncAccountEntryPoint"
	contentType           = "Content-Type"
	applicationForm       = "application/x-www-form-urlencoded"
	basicToken            = "Basic %s"
	authorization         = "Authorization"
	status                = "status"
	dolphin               = "Dolphin"
	destination           = "destination"
	scopeDolphin          = "request_to_dolphin"
	parseStatus           = "parse_status"
	scopeStatusSuccess    = "status_success"
)

// NewDolphinService function for initializing dolphin service
func NewDolphinService() (*DolphinService, error) {
	var (
		dolphin DolphinService
		err     error
		ok      bool
	)

	dolphin.BasicAuth, ok = os.LookupEnv("DOLPHIN_BASIC_AUTH")
	if !ok {
		return &dolphin, errors.New("you need to specify DOLPHIN_BASIC_AUTH in the environment variable")
	}

	dolphinURL, ok := os.LookupEnv("DOLPHIN_BASE_URL")
	if !ok {
		return &dolphin, errors.New("you need to specify DOLPHIN_BASE_URL in the environment variable")
	}

	dolphin.BaseURL, err = url.Parse(dolphinURL)
	if err != nil {
		return &dolphin, errors.New("url is invalid")
	}

	return &dolphin, nil
}

// RegisterMember function for registering new member
// whose sign up from data is not from dolphin
func (ds *DolphinService) RegisterMember(ctxReq context.Context, data serviceModel.MemberDolphin) error {
	ctx := "Service-RegisterMember"

	uri := ds.BaseURL.String() + urlSyncDolphin
	resp := serviceModel.Response{}

	headers := map[string]string{
		contentType:   applicationForm,
		authorization: fmt.Sprintf(basicToken, ds.BasicAuth),
	}

	formRM := url.Values{}
	formRM.Add("id", data.ID)
	formRM.Add("email", data.Email)
	formRM.Add("firstName", data.FirstName)
	formRM.Add("lastName", data.LastName)
	formRM.Add("gender", data.Gender)
	formRM.Add("dob", data.DOB)
	formRM.Add("mobile", data.Mobile)
	formRM.Add(status, data.Status)
	formRM.Add("created", data.Created)
	formRM.Add("facebookId", data.FacebookID)
	formRM.Add("googleId", data.GoogleID)
	formRM.Add("azureId", data.AzureID)
	formRM.Add(destination, dolphin)
	encodedRM := formRM.Encode()
	formDataRM := strings.NewReader(encodedRM)

	errRM := helper.GetHTTPNewRequest(context.Background(), "POST", uri, formDataRM, &resp, headers)
	if errRM != nil {
		helper.SendErrorLog(ctxReq, ctx, "scope_http_dolphin", errRM, encodedRM)
		return errRM
	}

	// check message and status
	if !resp.Data.Attributes.IsSuccess {
		errRM := errors.New(resp.Data.Attributes.Message)
		helper.SendErrorLog(ctxReq, ctx, "dolphin_message_attributes", errRM, resp)
		return errRM
	}

	return nil
}

// UpdateMember function for updating member to dolphin
func (ds *DolphinService) UpdateMember(ctxReq context.Context, data serviceModel.MemberDolphin) error {
	ctx := "Service-UpdateMember"

	uri := ds.BaseURL.String() + urlSyncDolphin
	resp := serviceModel.Response{}

	headers := map[string]string{
		contentType:   applicationForm,
		authorization: fmt.Sprintf(basicToken, ds.BasicAuth),
	}

	form := url.Values{}
	form.Add("id", data.ID)
	form.Add("firstName", data.FirstName)
	form.Add("lastName", data.LastName)
	form.Add("gender", data.Gender)
	form.Add("dob", data.DOB)
	form.Add("phone", data.Phone)
	form.Add("ext", data.Ext)
	form.Add("mobile", data.Mobile)
	form.Add("street1", data.Street1)
	form.Add("street2", data.Street2)
	form.Add("postalCode", data.PostalCode)
	form.Add("subDistrictId", data.SubDistrictID)
	form.Add("subDistrictName", data.SubDistrictName)
	form.Add("districtId", data.DistrictID)
	form.Add("districtName", data.DistrictName)
	form.Add("cityId", data.CityID)
	form.Add("cityName", data.CityName)
	form.Add("provinceId", data.ProvinceID)
	form.Add("provinceName", data.ProvinceName)
	form.Add(status, data.Status)
	form.Add("facebookId", data.FacebookID)
	form.Add("googleId", data.GoogleID)
	form.Add("azureId", data.AzureID)
	form.Add("created", data.Created)
	form.Add("lastModified", data.LastModified)
	form.Add(destination, dolphin)
	encoded := form.Encode()
	formData := strings.NewReader(encoded)

	err := helper.GetHTTPNewRequest(ctxReq, "PUT", uri, formData, &resp, headers)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "dolphin_http_req", err, encoded)
		return err
	}

	// check message and status
	if !resp.Data.Attributes.IsSuccess {
		err := errors.New(resp.Data.Attributes.Message)
		helper.SendErrorLog(ctxReq, ctx, "dolphin_check_attr", err, resp)
		return err
	}

	return nil
}

// ActivateMember function for updating member to dolphin
func (ds *DolphinService) ActivateMember(ctxReq context.Context, data serviceModel.MemberDolphin) error {
	ctx := "Service-UpdateMember"

	uri := ds.BaseURL.String() + urlSyncDolphin
	resp := serviceModel.Response{}

	headers := map[string]string{
		contentType:   applicationForm,
		authorization: fmt.Sprintf(basicToken, ds.BasicAuth),
	}

	form := url.Values{}
	form.Add("id", data.ID)
	form.Add(status, data.Status)
	form.Add("lastModified", data.LastModified)
	form.Add(destination, dolphin)
	encoded := form.Encode()
	formData := strings.NewReader(encoded)

	err := helper.GetHTTPNewRequest(context.Background(), "PUT", uri, formData, &resp, headers)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, "dolphin_http_activate", err, encoded)
		return err
	}

	// check message and status
	if !resp.Data.Attributes.IsSuccess {
		err := errors.New(resp.Data.Attributes.Message)
		helper.SendErrorLog(ctxReq, ctx, "dolphin_http_activate", err, resp)
		return err
	}

	return nil
}

// GetMember function for get member data from dolphin
func (ds *DolphinService) GetMember(ctxReq context.Context, id string) (*serviceModel.MemberResponse, error) {
	ctx := "Service-GetMember"

	uri := fmt.Sprintf("%s/index.php?entryPoint=syncAccountEntryPoint&id=%s", ds.BaseURL.String(), id)
	resp := serviceModel.MemberResponse{}

	headers := map[string]string{
		contentType:   applicationForm,
		authorization: fmt.Sprintf(basicToken, ds.BasicAuth),
	}

	err := helper.GetHTTPNewRequest(ctxReq, "GET", uri, nil, &resp, headers)
	if err != nil {
		helper.SendErrorLog(ctxReq, ctx, scopeDolphin, err, uri)
		return nil, err
	}

	// check message and status
	if resp.Data.Attributes.Message == DolphinDataNotFound {
		err := errors.New(resp.Data.Attributes.Message)
		return nil, err
	}

	return &resp, nil
}
