package model

import "time"

// Application data structure
type Application struct {
	ID           string    `json:"id" form:"id" fieldname:"id"`
	Name         string    `json:"name" form:"name"`
	Logo         string    `json:"logo" form:"logo"`
	URL          string    `json:"url" form:"url"`
	Created      time.Time `json:"created"`
	LastModified time.Time `json:"lastModified"`
}

// ApplicationList data structure
type ApplicationList struct {
	Data      []Application `json:"data"`
	TotalData int           `json:"totalData"`
}

// ParametersApplication data structure
type ParametersApplication struct {
	StrPage  string `json:"strPage" form:"strPage" query:"strPage" validate:"omitempty,numeric" fieldname:"strPage" url:"strPage"`
	Page     int    `json:"page" form:"page" query:"page" validate:"omitempty,numeric" fieldname:"page" url:"page"`
	StrLimit string `json:"strLimit" form:"strLimit" query:"strLimit" validate:"omitempty" fieldname:"strLimit" url:"strLimit"`
	Limit    int    `json:"limit" form:"limit" query:"limit" validate:"omitempty,numeric" fieldname:"limit" url:"limit"`
	Offset   int    `json:"offset" form:"offset" query:"offset" validate:"omitempty,numeric" fieldname:"offset" url:"offset"`
}

// ListApplication data structure
type ListApplication struct {
	Application []*Application `jsonapi:"relation,application" json:"application"`
	TotalData   int            `json:"totalData"`
}
