package products

import (
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ProductRatingAnalytics struct {
	Rating      float64 `bson:"rating" json:"rating"`
	NoOfRatings int     `bson:"noOfRatings" json:"noOfRatings"`
}

type Product struct {
	mgm.DefaultModel `bson:",inline"`

	ID                primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	Name              string                 `bson:"name" json:"name"`
	Price             float64                `bson:"price" json:"price"`
	Category          string                 `bson:"category" json:"category"`
	SubCategory       string                 `bson:"subCategory" json:"subCategory"`
	Discount          float64                `bson:"discount" json:"discount"`
	DiscountEndDate   time.Time              `bson:"discountEndDate" json:"discountEndDate"`
	DiscountStartDate time.Time              `bson:"discountStartDate" json:"discountStartDate"`
	Quantity          int                    `bson:"quantity" json:"quantity"`
	Description       string                 `bson:"description" json:"description"`
	Images            []string               `bson:"images" json:"images"`
	Rating            ProductRatingAnalytics `bson:"ratingAnalytics" json:"ratingAnalytics"`
}
