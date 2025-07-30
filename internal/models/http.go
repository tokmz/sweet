package models

type IDReq struct {
	ID int64 `json:"id" form:"id" binding:"required,min=1"`
}

type IdsReq struct {
	Ids []int64 `json:"ids" form:"ids" binding:"required,min=1,dive,min=1"`
}

type PageReq struct {
	Page int `json:"page" form:"page" binding:"min=1" default:"1"`
	Size int `json:"size" form:"size" binding:"min=1,max=100" default:"10"`
}

type SortReq struct {
	Field string `json:"field" form:"field"`
	Order string `json:"order" form:"order" binding:"oneof=asc desc" default:"desc"`
}

type TimeRangeReq struct {
	StartTime int64 `json:"start_time" form:"start_time"`
	EndTime   int64 `json:"end_time" form:"end_time"`
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type PageRes struct {
	Total int64 `json:"total"`
	List  any   `json:"list"`
}

func NewResponse(code int, message string, data any) *Response {
	return &Response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}

func NewPageRes(total int64, list any) *PageRes {
	return &PageRes{
		Total: total,
		List:  list,
	}
}
