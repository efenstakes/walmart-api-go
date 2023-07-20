package cart

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type Cart struct {
	mgm.DefaultModel `bson:",inline"`

	UserId    string    `bson:"userId" json:"userId"`
	Price     float64   `bson:"price" json:"price"`
	Quantity  int       `bson:"quantity" json:"quantity"`
	ProductID string    `bson:"productId" json:"productId"`
	SavedOn   time.Time `bson:"savedOn" json:"savedOn"`
}
