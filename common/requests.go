package common

type IDReq struct {
	ID int64 `json:"id" form:"id" binding:"required"`
}

type IdsReq struct {
	Ids []int64 `json:"ids" form:"ids" binding:"required"`
}

type SIDReq struct {
	SID string `json:"sid" form:"sid" binding:"required"`
}

type SIdsReq struct {
	SIds []string `json:"sids" form:"sids" binding:"required"`
}

type PageReq struct {
	Page  int `form:"page" json:"page" binding:"required,min=1"`
	Limit int `form:"limit" json:"limit" binding:"required,min=10,max=100"`
}

type StatusReq struct {
	ID     int64 `json:"id" form:"id" binding:"required"`
	Status int64 `json:"status" form:"status" binding:"required"`
}
