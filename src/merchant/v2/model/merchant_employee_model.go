package model

import (
	"fmt"
	"strings"
	"time"

	"github.com/Bhinneka/golib"
	"github.com/Bhinneka/user-service/helper"
	"gopkg.in/guregu/null.v4/zero"
)

var (
	allowedAllStatusEmpolyee      = []string{helper.TextInvited, helper.TextActive, helper.TextInactive, helper.TextRevoked}
	allowedStatusEmpolyeeInvited  = []string{helper.TextActive, helper.TextRevoked}
	allowedStatusEmpolyeeActive   = []string{helper.TextInactive, helper.TextRevoked}
	allowedStatusEmpolyeeInactive = []string{helper.TextActive, helper.TextRevoked}
)

const ErrorStatus = "status must be one of %s"

// B2CMerchantEmployee data structure
type B2CMerchantEmployee struct {
	ID         string    `json:"id"`
	MerchantID string    `json:"merchantId"`
	MemberID   string    `json:"memberId"`
	CreatedAt  time.Time `json:"createdAt"`
	CreatedBy  string    `json:"createdBy"`
	ModifiedAt time.Time `json:"modifiedAt"`
	ModifiedBy *string   `json:"modifiedBy"`
	Status     string    `json:"status"`
}

// B2CMerchantEmployeeData data structure
type B2CMerchantEmployeeData struct {
	ID              string      `json:"id"`
	MerchantID      string      `json:"merchantId"`
	MemberID        string      `json:"memberId"`
	FirstName       string      `json:"firstName"`
	LastName        zero.String `json:"lastName"`
	Email           string      `json:"email"`
	Gender          zero.String `json:"gender"`
	Mobile          zero.String `json:"mobile"`
	Phone           zero.String `json:"phone"`
	BirthDate       time.Time   `json:"-"`
	BirthDateString zero.String `json:"birthDate"`
	CreatedAt       time.Time   `json:"createdAt"`
	CreatedBy       string      `json:"createdBy"`
	ModifiedAt      *time.Time  `json:"modifiedAt"`
	ModifiedBy      *string     `json:"modifiedBy"`
	Status          zero.String `json:"status"`
	ProfilePicture  zero.String `json:"profilePicture"`
	MerchantLogo    zero.String `json:"merchantLogo"`
	MerchantName    zero.String `json:"merchantName"`
	MerchantType    zero.String `json:"merchantType"`
	VanityURL       zero.String `json:"vanityURL"`
	IsActive        bool        `json:"isActive"`
	IsPKP           bool        `json:"isPKP"`
}

// ListMerchantEmployee data structure
type ListMerchantEmployee struct {
	MerchantEmployee []*B2CMerchantEmployeeData `jsonapi:"relation,merchantEmployee" json:"merchantEmployee"`
	TotalData        int                        `json:"totalData"`
}

// ParametersMerchantEmployee data structure
type ParametersMerchantEmployee struct {
	StrPage  string `json:"strPage" form:"strPage" query:"strPage" validate:"omitempty,numeric" fieldname:"strPage" url:"strPage"`
	Page     int    `json:"page" form:"page" query:"page" validate:"omitempty,numeric" fieldname:"page" url:"page"`
	StrLimit string `json:"strLimit" form:"strLimit" query:"strLimit" validate:"omitempty" fieldname:"strLimit" url:"strLimit"`
	Limit    int    `json:"limit" form:"limit" query:"limit" validate:"omitempty,numeric" fieldname:"limit" url:"limit"`
	Offset   int    `json:"offset" form:"offset" query:"offset" validate:"omitempty,numeric" fieldname:"offset" url:"offset"`
	Sort     string `json:"sort" form:"sort" query:"sort" validate:"omitempty,oneof=asc desc" fieldname:"sort" url:"sort"`
	OrderBy  string `json:"orderBy" form:"orderBy" query:"orderBy" validate:"omitempty,oneof=id merchantId memberId" fieldname:"orderBy" url:"orderBy"`
	Status   string `json:"status" query:"status" validate:"omitempty,is-bool" fieldname:"status bank" url:"status"`
}

// AllowedSortFieldsMerchantEmployee is allowed field name for sorting
var AllowedSortFieldsMerchantEmployee = []string{
	"id",
	"merchantId",
	"memberId",
}

// QueryMerchantEmployeeParameters for search
type QueryMerchantEmployeeParameters struct {
	Page       int
	Limit      int
	Offset     int
	StrPage    string `json:"page" query:"page" form:"page" param:"page"`
	StrLimit   string `json:"limit" query:"limit" form:"limit" param:"limit"`
	OrderBy    string `json:"orderBy" query:"orderBy" form:"orderBy" param:"orderBy"`
	SortBy     string `json:"sortBy" query:"sortBy" form:"sortBy" param:"sortBy"`
	Search     string `json:"search" query:"search" form:"search" param:"search"`
	Status     string `json:"status" query:"status" form:"status" param:"status"`
	MerchantID string
	MemberID   string
	Name       string `json:"name" query:"name" form:"name" param:"name"`
	Email      string `json:"email" query:"email" form:"email" param:"email"`
}

// QueryCmsMerchantEmployeeParameters for search
type QueryCmsMerchantEmployeeParameters struct {
	Page       int
	Limit      int
	Offset     int
	StrPage    string `json:"page" query:"page" form:"page" param:"page"`
	StrLimit   string `json:"limit" query:"limit" form:"limit" param:"limit"`
	OrderBy    string `json:"orderBy" query:"orderBy" form:"orderBy" param:"orderBy"`
	SortBy     string `json:"sortBy" query:"sortBy" form:"sortBy" param:"sortBy"`
	Search     string `json:"search" query:"search" form:"search" param:"search"`
	Status     string `json:"status" query:"status" form:"status" param:"status"`
	MerchantID string `json:"merchantId" query:"merchantId" form:"merchantId" param:"merchantId"`
	MemberID   string `json:"memberId" query:"memberId" form:"memberId" param:"memberId"`
	Name       string `json:"name" query:"name" form:"name" param:"name"`
	Email      string `json:"email" query:"email" form:"email" param:"email"`
}

// Build ...
func (me *QueryMerchantEmployeeParameters) Build() ([]string, []interface{}) {
	var queries []string
	var queryValues []interface{}
	lenParams := 0

	if me.Status != "" {
		statusEmployee := strings.Split(me.Status, ",")
		qs := `me."status" IN (`
		vw := []string{}
		for _, t := range statusEmployee {
			lenParams++
			vw = append(vw, fmt.Sprintf(`$%d`, lenParams))
			status := strings.TrimSpace(t)
			queryValues = append(queryValues, status)
		}
		qs += strings.Join(vw, ",")
		qs += `)`
		queries = append(queries, qs)
	}

	if me.Search != "" {
		if strings.HasPrefix(me.Search, "MCH") {
			lenParams++
			queries = append(queries, fmt.Sprintf(`me."merchantId" = $%d`, lenParams))
			queryValues = append(queryValues, me.Search)
		} else if strings.HasPrefix(me.Search, "USR") {
			lenParams++
			queries = append(queries, fmt.Sprintf(`me."memberId" = $%d`, lenParams))
			queryValues = append(queryValues, me.Search)
		} else if err := golib.ValidateEmail(me.Search); err == nil {
			lenParams++
			queries = append(queries, fmt.Sprintf(`m."email" = $%d`, lenParams))
			queryValues = append(queryValues, me.Search)
		} else {
			lenParams++
			queries = append(queries, fmt.Sprintf(`CONCAT(m."firstName",' ', m."lastName") ILIKE $%d`, lenParams))
			me.Search = "%" + me.Search + "%"
			queryValues = append(queryValues, me.Search)
		}
	}

	if me.MerchantID != "" {
		lenParams++
		queries = append(queries, fmt.Sprintf(`me."merchantId" = $%d`, lenParams))
		queryValues = append(queryValues, me.MerchantID)
	}

	if me.MemberID != "" {
		lenParams++
		queries = append(queries, fmt.Sprintf(`me."memberId" = $%d`, lenParams))
		queryValues = append(queryValues, me.MemberID)
	}

	if me.Email != "" {
		lenParams++
		queries = append(queries, fmt.Sprintf(`m."email" = $%d`, lenParams))
		queryValues = append(queryValues, me.Email)
	}

	if me.Name != "" {
		lenParams++
		queries = append(queries, fmt.Sprintf(`CONCAT(m."firstName",' ', m."lastName") ILIKE $%d`, lenParams))
		me.Name = "%" + me.Name + "%"
		queryValues = append(queryValues, me.Name)
	}
	return queries, queryValues
}

// Validate allowed input
func (q *QueryMerchantEmployeeParameters) Validate() error {
	if !golib.StringInSlice(q.OrderBy, allowedOrderMerchantEmployee, false) {
		return fmt.Errorf("orderBy must be one of %s", strings.Join(allowedOrderMerchantEmployee[:4], delimiter))
	}
	if !golib.StringInSlice(q.SortBy, allowedSort, false) {
		return fmt.Errorf("sortBy must be one of %s", strings.Join(allowedSort[:2], delimiter))
	}
	if q.Status != "" {
		statuses := strings.Split(q.Status, ",")
		for _, ses := range statuses {
			sess := helper.TrimSpace(ses)
			if !golib.StringInSlice(sess, allowedStatusEmpolyee, false) {
				return fmt.Errorf(ErrorStatus, strings.Join(allowedStatusEmpolyee[:4], delimiter))
			}
		}
	}

	return nil
}

// Validate allowed input
func (cms *QueryCmsMerchantEmployeeParameters) Validate() error {
	if !golib.StringInSlice(cms.OrderBy, allowedOrderMerchantEmployee, false) {
		return fmt.Errorf("orderBy must be one of %s", strings.Join(allowedOrderMerchantEmployee[:4], delimiter))
	}
	if !golib.StringInSlice(cms.SortBy, allowedSort, false) {
		return fmt.Errorf("sortBy must be one of %s", strings.Join(allowedSort[:2], delimiter))
	}
	if cms.Status != "" {
		statuses := strings.Split(cms.Status, ",")
		for _, es := range statuses {
			ess := helper.TrimSpace(es)
			if !golib.StringInSlice(ess, allowedStatusEmpolyee, false) {
				return fmt.Errorf(ErrorStatus, strings.Join(allowedStatusEmpolyee[:4], delimiter))
			}
		}

	}
	return nil
}

// Validate allowed input
func (q *QueryMerchantEmployeeParameters) ValidateStatus(oldStatus, newStatus string) error {
	if !golib.StringInSlice(newStatus, allowedAllStatusEmpolyee, false) {
		return fmt.Errorf(ErrorStatus, strings.Join(allowedAllStatusEmpolyee[:4], delimiter))
	}

	switch oldStatus {
	case helper.TextInvited:
		if !golib.StringInSlice(newStatus, allowedStatusEmpolyeeInvited, false) {
			return fmt.Errorf(ErrorStatus, strings.Join(allowedStatusEmpolyeeInvited[:2], delimiter))
		}
	case helper.TextActive:
		if !golib.StringInSlice(newStatus, allowedStatusEmpolyeeActive, false) {
			return fmt.Errorf(ErrorStatus, strings.Join(allowedStatusEmpolyeeActive[:2], delimiter))
		}
	case helper.TextInactive:
		if !golib.StringInSlice(newStatus, allowedStatusEmpolyeeInactive, false) {
			return fmt.Errorf(ErrorStatus, strings.Join(allowedStatusEmpolyeeInactive[:2], delimiter))
		}
	case helper.TextRevoked:
		return fmt.Errorf("merchant status has been revoked")
	}

	return nil
}
