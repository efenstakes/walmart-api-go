package orders

import (
	"fmt"
	"net/http"
	"time"

	"github.com/efenstakes/walmart-api-g/accounts"
	"github.com/efenstakes/walmart-api-g/products"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// place an order
func MakeOrder(c *fiber.Ctx) error {
	// get the user making the order
	accountLocal := c.Locals("account")

	if accountLocal == nil {
		return fiber.NewError(http.StatusUnauthorized, "StatusUnauthorized")
	}

	account := accountLocal.(*accounts.Account)

	type InputOrderData struct {
		Products []struct {
			ID       string
			Quantity int
		}
	}

	inputOrder := new(InputOrderData)

	if err := c.BodyParser(inputOrder); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	// object ids of the products we are ordering
	productObjectIds := []primitive.ObjectID{}

	// create object ids from the product ids
	for _, inputItem := range inputOrder.Products {
		objID, _ := primitive.ObjectIDFromHex(inputItem.ID)
		productObjectIds = append(productObjectIds, objID)
	}

	// get products with the ids we got
	cursor, err := mgm.Coll(&products.Product{}).Find(mgm.Ctx(), bson.M{"_id": bson.M{"$in": productObjectIds}})
	// cursor, err := mgm.Coll(&products.Product{}).Find(mgm.Ctx(), bson.D{"productID": bson.D{"$in": productIds}})
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	// get the products from the cursor
	orderProducts := []products.Product{}
	if err := cursor.All(mgm.Ctx(), &orderProducts); err != nil {
		fmt.Println("An Error in Fetch Cursor")
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	fmt.Println("Your order has these products")
	for _, v := range orderProducts {
		fmt.Print(v.Name)
	}

	// return c.Status(http.StatusCreated).JSON(fiber.Map{"saved": true})

	// now we can create order products with the data we got
	// fill in order data
	orderProductsInfo := []OrderProduct{}

	for i := 0; i < len(inputOrder.Products); i++ {
		// id
		productId := inputOrder.Products[i].ID

		// quantity
		quantity := inputOrder.Products[i].Quantity

		newOrderProduct := OrderProduct{
			ID:       productId,
			Quantity: quantity,
		}

		// find the price for it
		for _, v := range orderProducts {
			if v.ID.Hex() == productId {
				newOrderProduct.Price = v.Price
			}
		}

		if newOrderProduct.Price == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "No Product"})
		}

		orderProductsInfo = append(orderProductsInfo, newOrderProduct)
	}

	// calculate the prices
	totalPrice := float64(0)
	for _, v := range orderProductsInfo {
		totalPrice += v.Price * float64(v.Quantity)
	}

	// build out the order
	orderData := &Order{
		Products:   orderProductsInfo,
		TotalPrice: totalPrice,
		UserId:     account.ID.Hex(),
		MadeOn:     time.Now(),
	}

	// saved order
	if err := mgm.Coll(&Order{}).Create(orderData); err != nil {
		return c.Status(400).JSON(fiber.Map{})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"saved": true, "order": orderData})
}

// get product ratings
func GetAll(c *fiber.Ctx) error {
	orders := []Order{}
	accountLocal := c.Locals("account")

	if accountLocal == nil {
		return fiber.NewError(http.StatusUnauthorized, "StatusUnauthorized")
	}

	account := accountLocal.(*accounts.Account)

	limit, limitErr := c.ParamsInt("limit", 20)
	offset, offsetErr := c.ParamsInt("offset", 0)

	if limitErr != nil || offsetErr != nil {
		return fiber.NewError(http.StatusBadRequest, "Arguments Error")
	}

	// Define options for the query
	findOptions := options.Find()

	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	filters := map[string]string{}
	filters["userId"] = account.ID.Hex()

	cursor, err := mgm.Coll(&Order{}).Find(mgm.Ctx(), filters, findOptions)

	if err := cursor.All(mgm.Ctx(), &orders); err != nil {
		fmt.Println("An Error in Fetch Cursor")
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(orders)
}
