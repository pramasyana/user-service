package model

// Parameters data structure
type Parameters struct {
	ID         string `query:"id"`
	Module     string `query:"module"`
	Action     string `query:"action"`
	Pack       string `query:"pack"`
	Page       int    `query:"page"`
	Limit      int    `query:"limit"`
	Sort       string `query:"sort"`
	OrderBy    string `query:"orderBy"`
	DateFrom   string `query:"dateFrom"`
	DateTo     string `query:"dateTo"`
	ViewType   string `query:"viewType"`
	Creator    string `query:"creator"`
	Target     string `query:"target"`
	ObjectType string `query:"objectType"`
}
