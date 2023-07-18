package accounts

import "github.com/kamva/mgm/v3"

type Account struct {
	mgm.DefaultModel `bson:",inline"`

	// ID         string   `json:"id" bson:"_id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
	JoinedOn string `json:"joinedOn" bson:"created_at"`
}
