package rating

import "github.com/kamva/mgm/v3"

type ProductRating struct {
	mgm.DefaultModel `bson:",inline"`

	Rating    float64 `bson:"rating" json:"rating" validate:"required,min=1"`
	UserId    string  `bson:"userId" json:"userId" validate:"required"`
	ProductID string  `bson:"productId" json:"productId" validate:"required"`
	Comment   string  `bson:"comment" json:"comment" validate:"required,max=900"`
}
