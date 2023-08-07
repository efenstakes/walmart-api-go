package accounts

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type Account struct {
	mgm.DefaultModel `bson:",inline"`

	// ID         string   `json:"id" bson:"_id"`
	Name     string `json:"name" validate:"required,min=5,max=20"`
	Password string `json:"password" validate:"required,min=8,max=20,containsany=!@#?*"`
	Email    string `json:"email" validate:"required,email"`

	// REGULAR || ADMIN
	Type string `json:"type" validate:"required"`

	JoinedOn time.Time `json:"joinedOn" bson:"joinedOn"`
}

func New() *Account {
	account := new(Account)

	account.Type = "ADMIN"
	account.JoinedOn = time.Now()

	return account
}
