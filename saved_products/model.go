package savedproducts

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type SavedProduct struct {
	mgm.DefaultModel `bson:",inline"`

	UserId    string    `bson:"userId" json:"userId"`
	ProductID string    `bson:"productId" json:"productId"`
	SavedOn   time.Time `bson:"savedOn" json:"savedOn"`
}
