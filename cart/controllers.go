package cart

import (
	"fmt"
	"net/http"
	"time"

	"github.com/efenstakes/walmart-api-g/accounts"
	"github.com/efenstakes/walmart-api-g/products"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// add item to cart
func Add(c *fiber.Ctx) error {
	// get our account
	accountLocal := c.Locals("account")

	if accountLocal == nil {
		return fiber.NewError(http.StatusUnauthorized)
	}

	// convert account Interface{} to account type
	account := accountLocal.(*accounts.Account)

	// get input data
	inputCartItem := new(Cart)
	if err := c.BodyParser(inputCartItem); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}
	inputCartItem.SavedOn = time.Now()
	inputCartItem.UserId = account.ID.Hex()
	inputCartItem.SavedOn = time.Now()

	// check if the product exists
	product := new(products.Product)
	err := mgm.Coll(&products.Product{}).FindByID(inputCartItem.ProductID, product)
	if product.Name == "" || err != nil {
	}

	// check if product has already been added
	item := new(Cart)
	filters := map[string]string{}
	filters["productId"] = inputCartItem.ProductID
	filters["userId"] = account.ID.Hex()

	error := mgm.Coll(&Cart{}).FindOne(mgm.Ctx(), filters).Decode(item)
	fmt.Println(item)
	// add quantity
	if error != nil {
		fmt.Println("err.Error()")
		fmt.Println(error.Error())
		// just update quantity
		updates := bson.D{{"$set", bson.D{{"quantity", inputCartItem.Quantity}}}}

		if inputCartItem.ID.Hex() == "" {
			return c.Status(http.StatusNotFound).JSON(fiber.Map{"message": "No Product Id Found"})
		}

		// Create a filter to find the product by ID
		filter := bson.D{{"_id", inputCartItem.ID.Hex()}}

		opts := options.Update().SetUpsert(true)
		result, err := mgm.Coll(product).UpdateOne(mgm.Ctx(), filter, updates, opts)
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"saved": false})
		}

		fmt.Println(result.ModifiedCount)
		fmt.Println(result.UpsertedCount)

		return c.Status(http.StatusCreated).JSON(fiber.Map{"saved": true, "item": inputCartItem})
	}

	//
	if item.ProductID != "" {
		return c.Status(http.StatusOK).JSON(fiber.Map{"saved": false, "message": "Already Saved"})
	}

	// add it
	inputCartItem.Price = product.Price
	if err := mgm.Coll(inputCartItem).Create(inputCartItem); err != nil {
		return c.Status(400).JSON(fiber.Map{})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"saved": true, "item": inputCartItem})
}

// get cart products
func GetAll(c *fiber.Ctx) error {
	products := []Cart{}

	// get our account
	accountLocal := c.Locals("account")

	if accountLocal == nil {
		return fiber.NewError(http.StatusUnauthorized)
	}

	account := accountLocal.(*accounts.Account)

	// get pagination info
	limit, limitErr := c.ParamsInt("limit", 20)
	offset, offsetErr := c.ParamsInt("offset", 0)

	// deal with pagination info errors
	if limitErr != nil || offsetErr != nil {
		return fiber.NewError(http.StatusBadRequest, "Arguments Error")
	}

	// Define options for the query
	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	// build out filters
	filters := map[string]string{}
	filters["userId"] = account.ID.Hex()

	// get data
	cursor, err := mgm.Coll(&Cart{}).Find(mgm.Ctx(), filters, findOptions)

	// get data fron cursor
	if err := cursor.All(mgm.Ctx(), &products); err != nil {
		fmt.Println("An Error in Fetch Cursor")
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	// deal with error
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(products)
}

// delete item from cart
func Delete(c *fiber.Ctx) error {
	productID := c.Params("id", "")

	if productID == "" {
		return fiber.NewError(http.StatusBadRequest)
	}

	// get our account
	accountLocal := c.Locals("account")

	if accountLocal == nil {
		return fiber.NewError(http.StatusUnauthorized)
	}

	// convert account Interface{} to account type
	account := accountLocal.(*accounts.Account)

	// delete
	deletedItem := new(Cart)
	err := mgm.Coll(deletedItem).FindOneAndDelete(mgm.Ctx(), bson.D{{"productId", productID}, {"userId", account.ID.Hex()}}).Decode(deletedItem)
	if err != nil {
		fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{"error": "Product Not Found"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"deleted": true})
}
