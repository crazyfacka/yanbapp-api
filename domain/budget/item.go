package budget

// Item represents a Budget Item
type Item struct {
	ID     uint    `json:"id"`
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
	Active bool    `json:"active"`
}
