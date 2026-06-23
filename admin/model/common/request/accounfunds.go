package request

type AccountFunds struct {
	Id     uint    `json:"id" form:"id" `
	Amount float64 `json:"amount" form:"amount" `
	Points float64 `json:"points" form:"points" `
}
