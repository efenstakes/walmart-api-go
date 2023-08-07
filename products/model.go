package products

import (
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductRatingAnalytics struct {
	Rating      float64 `bson:"rating" json:"rating" validate:"required,min=1"`
	NoOfRatings int     `bson:"noOfRatings" json:"noOfRatings" validate:"required"`
}

type Product struct {
	mgm.DefaultModel `bson:",inline"`

	ID                primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Name              string                 `bson:"name" json:"name" validate:"required"`
	Price             float64                `bson:"price" json:"price" validate:"required,min=1"`
	Category          string                 `bson:"category" json:"category" validate:"required"`
	SubCategory       string                 `bson:"subCategory" json:"subCategory"`
	Discount          float64                `bson:"discount" json:"discount" validate:"required,min=0"`
	DiscountEndDate   time.Time              `bson:"discountEndDate" json:"discountEndDate"`
	DiscountStartDate time.Time              `bson:"discountStartDate" json:"discountStartDate"`
	Quantity          int                    `bson:"quantity" json:"quantity" validate:"required,min=1"`
	Description       string                 `bson:"description" json:"description" validate:"required,min=20,max=900"`
	Images            []string               `bson:"images" json:"images"`
	Rating            ProductRatingAnalytics `bson:"ratingAnalytics" json:"ratingAnalytics"`
}
