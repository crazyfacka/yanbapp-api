package users

import (
	"math/rand"
	"time"

	"github.com/crazyfacka/yanbapp-api/domain/accounts"
	"github.com/crazyfacka/yanbapp-api/domain/budget"
	"golang.org/x/crypto/bcrypt"
)

// User represents a user
type User struct {
	ID            uint   `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	PasswordPlain string `json:"password"`
	SessionKey    string `json:"session_key,omitempty"`

	Accounts     []accounts.Account `json:"accounts,omitempty"`
	BudgetGroups []budget.Group     `json:"budget_groups,omitempty"`
}

const (
	// SessionDuration is the default duration for a user session
	SessionDuration = time.Hour * 24 * 5
	// SessionCookieName is the cookie name
	SessionCookieName = "yanbapp_session"

	hashSeed   = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	hashLength = 32
)

// ValidateUser compares user generated hash with the provided one
func (u *User) ValidateUser(hash []byte) error {
	return bcrypt.CompareHashAndPassword(hash, []byte(u.PasswordPlain))
}

// GenRandSessionString generates a random session string for this user
func (u *User) GenRandSessionString() {
	b := make([]byte, hashLength)
	for i := range b {
		b[i] = hashSeed[rand.Int63()%int64(len(hashSeed))]
	}

	u.SessionKey = string(b)
}
