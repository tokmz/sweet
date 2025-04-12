package common

type IDReq struct {
	ID int64 `json:"id" form:"id"`
}

type IdsReq struct {
	Ids []int64 `json:"ids" form:"ids"`
}

type SID struct {
	ID string `json:"id" form:"id"`
}

type SIds struct {
	Ids []string `json:"ids" form:"ids"`
}

type PageReq struct {
	Page int `json:"page" form:"page"`
	Size int `json:"size" form:"size"`
}

type StatusReq struct {
	IdsReq
	Status int64 `json:"status" form:"status"`
}
