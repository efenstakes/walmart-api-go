package rating

import "github.com/kamva/mgm/v3"

type ProductRating struct {
	mgm.DefaultModel `bson:",inline"`

	Rating    float64 `bson:"rating" json:"rating"`
	UserId    string  `bson:"userId" json:"userId"`
	ProductID string  `bson:"productId" json:"productId"`
	Comment   string  `bson:"comment" json:"comment"`
}
