package budget

// Group represents a Budget Group
type Group struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Items  []Item `json:"items"`
	Active bool   `json:"active"`
}
