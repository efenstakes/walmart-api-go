package cart

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type Cart struct {
	mgm.DefaultModel `bson:",inline"`

	UserId    string    `bson:"userId" json:"userId" validate:"required"`
	Price     float64   `bson:"price" json:"price" validate:"required,min=1"`
	Quantity  int       `bson:"quantity" json:"quantity" validate:"required,min=0"`
	ProductID string    `bson:"productId" json:"productId" validate:"required"`
	SavedOn   time.Time `bson:"savedOn" json:"savedOn"`
}
