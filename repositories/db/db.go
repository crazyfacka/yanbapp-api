package db

import (
	"database/sql"
	"log"
	"strconv"
	"time"

	"github.com/crazyfacka/yanbapp-api/domain/accounts"
	"github.com/crazyfacka/yanbapp-api/domain/budget"
	"github.com/crazyfacka/yanbapp-api/domain/transactions"
	"github.com/crazyfacka/yanbapp-api/domain/users"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
)

// DB repository struct
type DB struct {
	db *sql.DB
}

/* === PRIVATE FNs === */

func (db *DB) newRepeatTransaction(u *users.User, t *transactions.Transaction) (int64, error) {
	var (
		tx     *sql.Tx
		stmt   *sql.Stmt
		res    sql.Result
		lastID int64
		err    error
	)

	tx, err = db.db.Begin()
	if err != nil {
		return 0, err
	}

	defer tx.Rollback()

	stmt, err = tx.Prepare(`INSERT INTO repeat_transactions
														(user_id, budget_item, amount, day_of_month, description, active)
														VALUES (?, ?, ?, ?, ?, ?)`)

	defer stmt.Close()

	res, err = stmt.Exec(u.ID, t.BudgetItem.ID, t.Amount, t.When.Day(), t.Description, 1)
	if err != nil {
		return 0, err
	}

	lastID, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	stmt, err = tx.Prepare(`UPDATE transactions
														SET repeat_transaction_id = ?
														WHERE id = ?`)

	_, err = stmt.Exec(lastID, t.ID)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}

	return lastID, err
}

func (db *DB) disableRepeatTransaction(transactionID uint) error {
	var (
		tx       *sql.Tx
		stmt     *sql.Stmt
		repeatID int64
		err      error
	)

	err = db.db.QueryRow(`SELECT repeat_transaction_id FROM transactions WHERE id = ?`, transactionID).Scan(&repeatID)
	if err == nil && repeatID > 0 {
		tx, err = db.db.Begin()
		if err != nil {
			return err
		}

		defer tx.Rollback()

		stmt, err = tx.Prepare(`UPDATE transactions
																SET repeat_transaction_id = NULL
																WHERE id = ?`)

		defer stmt.Close()

		if err != nil {
			return err
		}

		_, err = stmt.Exec(transactionID)
		if err != nil {
			return err
		}

		stmt, err = tx.Prepare(`UPDATE repeat_transactions
																SET active = 0
																WHERE id = ?`)

		if err != nil {
			return err
		}

		_, err = stmt.Exec(repeatID)
		if err != nil {
			return err
		}

		err = tx.Commit()
		if err != nil {
			return err
		}
	}

	return nil
}

/* =================== */

// GetBudgetInfo returns all the budget info
func (db *DB) GetBudgetInfo(u *users.User) error {
	var (
		groupID   uint
		groupName string
		itemID    uint
		itemName  string
		amount    sql.NullFloat64

		budgetGroups map[uint]*budget.Group
	)

	stmt, err := db.db.Prepare(`SELECT bg.id AS bg_id, bg.name AS ` + "`group`" + `, bi.id AS bi_id, bi.name AS item, ba.amount
																FROM budget_items bi
																JOIN budget_groups bg
																	ON bi.group_id = bg.id
																LEFT JOIN budget_amounts ba
																	ON bi.id = ba.item_id
																WHERE bi.active = 1
																	AND bg.active = 1
																	AND bg.user_id = ?`)

	if err != nil {
		return err
	}

	defer stmt.Close()
	rows, err := stmt.Query(u.ID)

	if err != nil {
		return err
	}

	budgetGroups = make(map[uint]*budget.Group)

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&groupID, &groupName, &itemID, &itemName, &amount)
		if err != nil {
			log.Printf("error parsing row: %s", err.Error())
		}

		if _, ok := budgetGroups[groupID]; !ok {
			budgetGroups[groupID] = &budget.Group{
				ID:   groupID,
				Name: groupName,
			}
		}

		item := budget.Item{
			ID:     itemID,
			Name:   itemName,
			Amount: amount.Float64,
		}

		budgetGroups[groupID].Items = append(budgetGroups[groupID].Items, item)
	}

	for _, v := range budgetGroups {
		u.BudgetGroups = append(u.BudgetGroups, *v)
	}

	return nil
}

// SaveBudgetItem updates the information of a single budget item
func (db *DB) SaveBudgetItem(u *users.User, bi *budget.Item) (int64, error) {
	// For now it supports only setting the budgeted amount
	var (
		stmt   *sql.Stmt
		res    sql.Result
		lastID int64
		err    error
	)

	stmt, err = db.db.Prepare(`UPDATE budget_amounts
																SET amount = ?
																WHERE id = ?`)

	if err != nil {
		return 0, err
	}

	defer stmt.Close()
	res, err = stmt.Exec(bi.Amount, bi.ID)
	if err != nil {
		return 0, err
	}

	lastID, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

// GetTransactions returns all transactions for a given account ID
func (db *DB) GetTransactions(u *users.User, accountID uint) ([]transactions.Transaction, error) {
	var (
		id             uint
		budgetItemID   sql.NullInt64
		budgetItemName sql.NullString
		expense        uint8
		amount         float64
		when           time.Time
		repeat         sql.NullInt64
		description    sql.NullString
		createdAt      time.Time
		updatedAt      time.Time

		ts []transactions.Transaction
	)

	stmt, err := db.db.Prepare(`SELECT t.id, bi.id AS budget_item_id, bi.name AS budget_item_name, t.expense, t.amount, ` + "t.`when`" + `, t.repeat_transaction_id, t.description, t.created_at, t.updated_at
																FROM transactions t
																LEFT JOIN budget_items bi
																	ON t.budget_item = bi.id
																WHERE user_id = ? AND account_id = ?`)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	rows, err := stmt.Query(u.ID, accountID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id,
			&budgetItemID,
			&budgetItemName,
			&expense,
			&amount,
			&when,
			&repeat,
			&description,
			&createdAt,
			&updatedAt)

		if err != nil {
			log.Printf("error parsing row: %s", err.Error())
		}

		boolExpense := false
		if expense == 1 {
			boolExpense = true
		}

		t := transactions.Transaction{
			ID:        id,
			AccountID: accountID,
			BudgetItem: budget.Item{
				ID:   uint(budgetItemID.Int64),
				Name: budgetItemName.String,
			},
			Expense:     boolExpense,
			Amount:      amount,
			When:        when,
			Repeat:      int(repeat.Int64),
			Description: description.String,
			Created:     createdAt,
			Updated:     updatedAt,
		}

		ts = append(ts, t)
	}

	return ts, nil
}

// SaveTransaction saves a new or updates an existing transaction
func (db *DB) SaveTransaction(u *users.User, t *transactions.Transaction) (int64, int64, error) {
	var (
		stmt     *sql.Stmt
		res      sql.Result
		lastID   int64
		repeatID int64
		err      error
	)

	if t.Repeat == -1 { // New repeat transaction
		repeatID, err = db.newRepeatTransaction(u, t)
		if err != nil {
			return 0, 0, err
		}
	} else if t.Repeat == 0 { // Check if this t.ID has an existing repeat one, and inactivate it
		err = db.disableRepeatTransaction(t.ID)
		if err != nil {
			return 0, 0, err
		}
	}

	if t.ID == 0 {
		stmt, err = db.db.Prepare(`INSERT INTO transactions
																	(budget_item, amount, expense, ` + "`when`" + `, description, user_id, account_id)
																	VALUES (?, ?, ?, ?, ?, ?, ?)`)
	} else {
		stmt, err = db.db.Prepare(`UPDATE transactions t
																	SET budget_item = ?,
																			amount = ?,
																			expense = ?,
																			` + "t.`when`" + ` = ?,
																			description = ?
																	WHERE user_id = ? AND account_id = ? AND id = ?`)
	}

	if err != nil {
		return 0, 0, err
	}

	budgetItemID := sql.NullInt64{}
	if t.Expense {
		budgetItemID.Int64 = int64(t.BudgetItem.ID)
		budgetItemID.Valid = true
	}

	args := []interface{}{budgetItemID, t.Amount, t.Expense, t.When, t.Description, u.ID, t.AccountID}
	if t.ID != 0 {
		args = append(args, t.ID)
	}

	defer stmt.Close()
	res, err = stmt.Exec(args...)
	if err != nil {
		return 0, 0, err
	}

	lastID, err = res.LastInsertId()
	if err != nil {
		return 0, 0, err
	}

	return lastID, repeatID, nil
}

// GetAccounts returns all available Accounts for a given user
func (db *DB) GetAccounts(u *users.User) error {
	var (
		id     uint
		name   string
		amount float64
		typ    uint
		active uint8
	)

	stmt, err := db.db.Prepare(`SELECT id, name, amount, type, active
																FROM accounts
																WHERE user_id = ?`)

	if err != nil {
		return err
	}

	defer stmt.Close()
	rows, err := stmt.Query(u.ID)

	if err != nil {
		return err
	}

	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&id, &name, &amount, &typ, &active)
		if err != nil {
			log.Printf("error parsing row: %s", err.Error())
		}

		boolActive := false
		if active == 1 {
			boolActive = true
		}

		a := accounts.Account{
			ID:     id,
			Type:   typ,
			Name:   name,
			Amount: amount,
			Active: boolActive,
		}

		u.Accounts = append(u.Accounts, a)
	}

	return nil
}

// SaveAccount saves a new or updates an existing account
func (db *DB) SaveAccount(u *users.User, a *accounts.Account) (int64, error) {
	var (
		stmt   *sql.Stmt
		lastID int64
		err    error
	)

	if a.ID == 0 {
		stmt, err = db.db.Prepare(`INSERT INTO accounts
																	(name, amount, type, active, user_id)
																	VALUES (?, ?, ?, ?, ?)`)
	} else {
		stmt, err = db.db.Prepare(`UPDATE accounts
																	SET name = ?,
																			amount = ?,
																			type = ?,
																			active = ?
																	WHERE user_id = ? AND id = ?`)
	}

	args := []interface{}{a.Name, a.Amount, a.Type, a.Active, u.ID}
	if a.ID != 0 {
		args = append(args, a.ID)
	}

	defer stmt.Close()
	res, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}

	lastID, err = res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastID, nil
}

// ValidateUser validates that this users exists and populates its domain
func (db *DB) ValidateUser(u *users.User) error {
	var (
		passwordHash []byte
	)

	err := db.db.QueryRow("SELECT id, name, password FROM users WHERE email = ?", u.Email).Scan(&u.ID, &u.Name, &passwordHash)
	if err != nil {
		return err
	}

	if err = u.ValidateUser(passwordHash); err != nil {
		return err
	}

	log.Printf("validated user ID %d", u.ID)

	return nil
}

// Close closes the DB connection
func (db *DB) Close() error {
	return db.db.Close()
}

// NewDB initializes a DB repository
func NewDB(user string, password string, host string, port int, schema string) *DB {
	connString := user + ":" + password + "@tcp(" + host + ":" + strconv.Itoa(port) + ")/" + schema + "?parseTime=true"

	db, err := sql.Open("mysql", connString)
	if err != nil {
		log.Fatalln("can't connect to the database")
	}

	err = db.Ping()
	if err != nil {
		log.Fatalln("error pinging db")
	}

	return &DB{
		db: db,
	}
}
