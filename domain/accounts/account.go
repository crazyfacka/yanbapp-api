package accounts

// Account represents an account
type Account struct {
	ID     uint    `json:"id"`
	Type   uint    `json:"type"`
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Active bool    `json:"active"`
}
