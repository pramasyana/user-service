package model

// StaticData data structure
type StaticData struct {
	ID              int    `json:"id"`
	Title           string `json:"title"`
	SubTitle        string `json:"subTitle"`
	MetaTitle       string `json:"metaTitle"`
	MetaDescription string `json:"metaDescription"`
	Placement       string `json:"placement"`
	ContentType     string `json:"contentType"`
	Content         string `json:"content"`
	ReviveContent   string `json:"reviveContent"`
	Slug            string `json:"slug"`
	IsActive        bool   `json:"isActive"`
	Created         string `json:"created"`
	LastModified    string `json:"lastModified"`
	ZoneID          []int  `json:"zoneId"`
}

// ResponseGWSStatic data structure
type ResponseGWSStatic struct {
	Code    int           `json:"code"`
	Message string        `json:"message"`
	Data    StaticDataGWS `json:"data"`
}

// StaticDataGWS data structure
type StaticDataGWS struct {
	GetStaticDetail StaticDetail `json:"getStaticPageById"`
}

// StaticDetail data structure
type StaticDetail struct {
	Code    int        `json:"code"`
	Success bool       `json:"success"`
	Result  StaticData `json:"result"`
}
