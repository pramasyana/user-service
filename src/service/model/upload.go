package model

// UploadData data structure
type UploadData struct {
	URL string `json:"url"`
}

// ResponseUploadService data structure
type ResponseUploadService struct {
	Code    int        `json:"code"`
	Message string     `json:"message"`
	Data    UploadData `json:"data"`
}

// UploadServiceResult data structure
type UploadServiceResult struct {
	Result interface{}
	Error  error
}
