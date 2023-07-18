package products

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductRating struct {
	Rating      float64 `bson:"rating" json:"rating"`
	NoOfRatings int     `bson:"noOfRatings" json:"noOfRatings"`
}

type Product struct {
	ID                primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name              string             `bson:"name" json:"name"`
	Price             float64            `bson:"price" json:"price"`
	Discount          float64            `bson:"discount" json:"discount"`
	LastDiscountDate  string             `bson:"lastDiscountDate" json:"lastDiscountDate"`
	DiscountStartDate string             `bson:"discountStartDate" json:"discountStartDate"`
	Quantity          int                `bson:"quantity" json:"quantity"`
	Description       string             `bson:"description" json:"description"`
	Images            []string           `bson:"images" json:"images"`
	Rating            ProductRating      `bson:"rating" json:"rating"`
}
