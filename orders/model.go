package orders

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type OrderProduct struct {
	ID       string  `bson:"id" json:"id"`
	Quantity int     `bson:"quantity" json:"quantity" validate:"required,min=1"`
	Price    float64 `bson:"price" json:"price" validate:"required,min=1"`
}

type Order struct {
	mgm.DefaultModel `bson:",inline"`

	UserId     string         `bson:"userId" json:"userId" validate:"required"`
	TotalPrice float64        `bson:"totalPrice" json:"totalPrice" validate:"required,min=1"`
	Products   []OrderProduct `bson:"products" json:"products"`
	MadeOn     time.Time      `bson:"madeOn" json:"madeOn"`
}
