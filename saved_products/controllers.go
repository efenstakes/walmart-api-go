package savedproducts

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

// add a saved product
func Add(c *fiber.Ctx) error {
	accountLocal := c.Locals("account")

	if accountLocal == nil {
		return fiber.NewError(http.StatusUnauthorized, "StatusUnauthorized")
	}

	account := accountLocal.(*accounts.Account)

	// get input data
	inputSavedProduct := new(SavedProduct)
	if err := c.BodyParser(inputSavedProduct); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}
	inputSavedProduct.UserId = account.ID.Hex()
	inputSavedProduct.SavedOn = time.Now()

	// check if the product exists
	product := new(products.Product)
	err := mgm.Coll(&products.Product{}).FindByID(inputSavedProduct.ProductID, product)
	if product.Name == "" || err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Product Not Found"})
	}

	// check if product has already been added
	item := new(SavedProduct)
	filters := map[string]string{}
	filters["productId"] = inputSavedProduct.ProductID
	filters["userId"] = account.ID.Hex()

	error := mgm.Coll(&SavedProduct{}).FindOne(mgm.Ctx(), filters).Decode(item)
	fmt.Println(item)
	if error != nil {
		fmt.Println("err.Error()")
		fmt.Println(error.Error())
	}

	if item.ProductID != "" {
		return c.Status(http.StatusOK).JSON(fiber.Map{"saved": false, "message": "Already Saved"})
	}

	// add it
	if err := mgm.Coll(inputSavedProduct).Create(inputSavedProduct); err != nil {
		return c.Status(400).JSON(fiber.Map{})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"saved": true, "item": inputSavedProduct})
}

// get saved product
func GetAll(c *fiber.Ctx) error {
	products := []SavedProduct{}
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

	cursor, err := mgm.Coll(&SavedProduct{}).Find(mgm.Ctx(), filters, findOptions)

	if err := cursor.All(mgm.Ctx(), &products); err != nil {
		fmt.Println("An Error in Fetch Cursor")
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(products)
}

// delete saved item
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
	deletedItem := new(SavedProduct)
	err := mgm.Coll(deletedItem).FindOneAndDelete(mgm.Ctx(), bson.D{{"productId", productID}, {"userId", account.ID.Hex()}}).Decode(deletedItem)
	if err != nil {
		fmt.Println(err)
		return c.Status(400).JSON(fiber.Map{"error": "Product Not Found"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"deleted": true})
}
