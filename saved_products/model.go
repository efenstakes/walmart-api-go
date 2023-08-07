package savedproducts

import (
	"time"

	"github.com/kamva/mgm/v3"
)

// Name     string `json:"name" validate:"required,min=5,max=20"`
// Password string `json:"password" validate:"required,min=8,max=20,containsany=!@#?*"`
// Email    string `json:"email" validate:"required,email"`
type SavedProduct struct {
	mgm.DefaultModel `bson:",inline"`

	UserId    string    `bson:"userId" json:"userId" validate:"required"`
	ProductID string    `bson:"productId" json:"productId" validate:"required"`
	SavedOn   time.Time `bson:"savedOn" json:"savedOn"`
}
