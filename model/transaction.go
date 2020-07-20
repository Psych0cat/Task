package model

type Transaction struct {
	State  string  `json:"state" required:"true"`
	Amount float64 `json:"amount" required:"true"`
	Id     string  `json:"transactionId" required:"true"`
}

