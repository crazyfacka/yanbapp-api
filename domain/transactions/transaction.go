package transactions

import (
	"time"

	"github.com/crazyfacka/yanbapp-api/domain/budget"
)

// Transaction represents a transaction
type Transaction struct {
	ID          uint        `json:"id"`
	AccountID   uint        `json:"account_id"`
	BudgetItem  budget.Item `json:"budget_item"`
	Expense     bool        `json:"expense"`
	Amount      float64     `json:"amount"`
	When        time.Time   `json:"when"`
	Repeat      int         `json:"repeat"`
	Description string      `json:"description"`
	Created     time.Time   `json:"created_at"`
	Updated     time.Time   `json:"updated_at"`
}
