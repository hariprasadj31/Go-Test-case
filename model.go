package main

type AssociateRequest struct {
	BusinessId string `json:"BusinessId"`
	PayeeRef   string `json:"PayeeRef"`
}

type AssociateWrapper struct {
	AssociateRequests []AssociateRequest `json:"AssociateRequests"`
}

type ErrorResponse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Message string `json:"message"`
}

type BusinessRequest struct {
	BusinessNm string `json:"BusinessNm"`
	FirstNm    string `json:"FirstNm"`
	MiddleNm   string `json:"MiddleNm"`
	LastNm     string `json:"LastNm"`
	Suffix     string `json:"Suffix"`
	PayerRef   string `json:"PayerRef"`
	TradeNm    string `json:"TradeNm"`
	IsEIN      bool   `json:"IsEIN"`
	EINorSSN   string `json:"EINorSSN"`
	Email      string `json:"Email"`
	ContactNm  string `json:"ContactNm"`
}
