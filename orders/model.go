package orders

import (
	"time"

	"github.com/kamva/mgm/v3"
)

type OrderProduct struct {
	ID       string  `bson:"id" json:"id"`
	Quantity int     `bson:"quantity" json:"quantity"`
	Price    float64 `bson:"price" json:"price"`
}

type Order struct {
	mgm.DefaultModel `bson:",inline"`

	UserId     string         `bson:"userId" json:"userId"`
	TotalPrice float64        `bson:"totalPrice" json:"totalPrice"`
	Products   []OrderProduct `bson:"products" json:"products"`
	MadeOn     time.Time      `bson:"madeOn" json:"madeOn"`
}
