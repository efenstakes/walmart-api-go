package main

import (
	"fmt"
	"os"

	"github.com/efenstakes/walmart-api-g/products"
	"github.com/joho/godotenv"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// this is always called before main making it a great place to initialize
func inite() {
	if err := godotenv.Load(); err != nil {
		panic("Couldn't load variables from environment")
	}
	err := mgm.SetDefaultConfig(
		nil, "walmart", options.Client().ApplyURI(os.Getenv("DB_URI")),
	)
	if err != nil {
		panic("Could not connect to MongoDB")
	} else {
		fmt.Println("Connected to db")
	}
}

func maine() {

	product := products.Product{}
	err := mgm.Coll(&products.Product{}).FindByID("64b6f2999debdaa3cd5ef947", &product)
	if err != nil {
		fmt.Println("Product Not Found")
	}

	fmt.Println("product id ", product.ID.Hex())
	fmt.Println("product name ", product.Name)

	if product.Name == "" {
		fmt.Println("Product Not Found")
	}

	return

	productIds := []string{}
	productIds = append(productIds, "64b6f2999debdaa3cd5ef647")
	productIds = append(productIds, "64b6e57808b7e2e7ed102442")

	// filters := bson.D{{"id", bson.D{{"$in", bson.A{"64b6f2999debdaa3cd5ef647"}}}}}
	// cursor, err := mgm.Coll(&products.Product{}).Find(mgm.Ctx(), filters)
	// cursor, err := mgm.Coll(&products.Product{}).Find(mgm.Ctx(), bson.M{"id": bson.M{"$in": productIds}})

	filters := map[string]primitive.ObjectID{}

	productObjectIds := []primitive.ObjectID{}

	for _, productID := range productIds {
		objID, _ := primitive.ObjectIDFromHex(productID)
		filters["_id"] = objID
		productObjectIds = append(productObjectIds, objID)
	}

	cursor, err := mgm.Coll(&products.Product{}).Find(mgm.Ctx(), bson.M{"_id": bson.M{"$in": productObjectIds}})

	// filters["_id"] = objID // "64b6f2999debdaa3cd5ef647"
	// cursor, err := mgm.Coll(&products.Product{}).Find(mgm.Ctx(), filters)
	if err != nil {
		fmt.Println("error ", err.Error())
	}

	orderProducts := []products.Product{}
	if err := cursor.All(mgm.Ctx(), &orderProducts); err != nil {
		fmt.Println("An Error in Fetch Cursor")
		fmt.Println("error ", err.Error())
	}

	fmt.Println(orderProducts)

	fmt.Println("Your order has these products")
	for _, v := range orderProducts {
		fmt.Println(v.Name)
	}
}
