package responses

type BaseResponse[T any] struct {
	Code    int           `json:"code,omitempty"`
	Message string        `json:"message"`
	Data    *T            `json:"data,omitempty"`
	Paging  *PageMetadata `json:"paging,omitempty"`
	Errors  string        `json:"errors,omitempty"`
}

type PageMetadata struct {
	Page      int `json:"page"`
	Size      int `json:"size"`
	TotalItem int `json:"total_item"`
	TotalPage int `json:"total_page"`
}
