package request

// PageInfo Paging common input parameter structure
type PageInfo struct {
	Page           int    `json:"page" form:"page" `        // 页码
	Limit          int    `json:"limit" form:"limit" `      // 每页大小
	Keyword        string `json:"kw" form:"kw"`             //关键字
	Count          bool   `json:"count" form:"count"`       //关键字
	Offset         int    `json:"offset" form:"offset"`     //关键字
	PageSize       int    `json:"pageSize" form:"pageSize"` // 每页大小
	Account        string `json:"account" form:"account"`   //
	GroupId        int    `json:"group_id" form:"group_id"`
	Status         string `json:"status" form:"status"` //
	Id             int    `json:"id" form:"id" `
	StockTypes     int    `json:"stock_types" form:"stock_types" `
	IsBulk         int    `json:"is_bulk" form:"is_bulk" `
	MStatus        int    `json:"m_status" form:"m_status" `
	OrderTradeType int    `json:"order_trade_type" form:"order_trade_type" `
	OrderType      int    `json:"order_type" form:"order_type" `
	SearchField    string `json:"search_field" form:"search_field" `
	Type           int    `json:"type" form:"type" `
	IsTest         int    `json:"is_test" form:"is_test" `
	Active         int    `json:"active" form:"active" `
}

// GetById Find by id structure
type GetById struct {
	ID uint `form:"id" bind:"required" json:"id"` // 主键ID
}
type GetByUserId struct {
	Id int `form:"id" bind:"required" json:"id"` // 主键ID
}

func (r *GetByUserId) Uint() uint {
	return uint(r.Id)
}
func (r *GetById) Uint() uint {
	return uint(r.ID)
}
func (r *GetById) Uint32() uint32 {

	return uint32(r.ID)
}

type IdsReq struct {
	Ids    []int  `json:"ids[]" form:"ids[]"`
	Action string `json:"action" form:"action"`
}

// GetAuthorityId Get role by id structure
type GetAuthorityId struct {
	AuthorityId uint `json:"authorityId" form:"authorityId"` // 角色ID
}

type Empty struct{}
