package model

type Tokens struct {
	Access  string `json:"access"`
	Refresh string `json:"refresh"`
	Valid   bool   `json:"valid"`
}
