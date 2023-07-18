package accounts

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type Account struct {
	mgm.DefaultModel `bson:",inline"`

	// ID         string   `json:"id" bson:"_id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`

	// REGULAR || ADMIN
	Type string `json:"type"`

	JoinedOn time.Time `json:"joinedOn" bson:"joinedOn"`
}

func New() *Account {
	account := new(Account)

	account.Type = "ADMIN"
	account.JoinedOn = time.Now()

	return account
}
