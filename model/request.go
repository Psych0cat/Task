package model

type Request struct {
	State  string `json:"state" required:"true"`
	Amount string `json:"amount" required:"true"`
	Id     string `json:"transactionId" required:"true"`
}

