package products

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Create(c *fiber.Ctx) error {
	accountLocal := c.Locals("account")

	if accountLocal == nil {
		return fiber.NewError(http.StatusUnauthorized)
	}

	// account := accountLocal.(accounts.Account)

	product := new(Product)

	if err := c.BodyParser(product); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	if err := mgm.Coll(product).Create(product); err != nil {
		return c.Status(400).JSON(fiber.Map{})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"product": product})
}

func Update(c *fiber.Ctx) error {
	productID := c.Params("id")
	fmt.Println("productID ", productID)
	accountLocal := c.Locals("account")

	if accountLocal == nil {
		return fiber.NewError(http.StatusUnauthorized)
	}

	inputProduct := new(Product)

	if err := c.BodyParser(inputProduct); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	// Prepare the update fields
	updates := bson.D{{"$set", bson.D{{"name", inputProduct.Name}, {"description", inputProduct.Description}, {"price", inputProduct.Price}, {"category", inputProduct.Category}, {"subCategory", inputProduct.SubCategory}, {"images", inputProduct.Images}}}}

	// build id for mongo
	id, _ := primitive.ObjectIDFromHex(productID)

	// Create a filter to find the product by ID
	filter := bson.D{{"_id", id}}

	opts := options.Update().SetUpsert(true)
	result, err := mgm.Coll(inputProduct).UpdateOne(context.TODO(), filter, updates, opts)
	if err != nil {
		fmt.Println("Error ", err.Error())
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"product": inputProduct})
	}
	fmt.Printf("Number of documents updated: %v\n", result.ModifiedCount)
	fmt.Printf("Number of documents upserted: %v\n", result.UpsertedCount)

	return c.Status(http.StatusCreated).JSON(fiber.Map{"product": inputProduct})
}

func SetDiscount(c *fiber.Ctx) error {
	productID := c.Params("id")
	fmt.Println("productID ", productID)
	accountLocal := c.Locals("account")

	if accountLocal == nil {
		return fiber.NewError(http.StatusUnauthorized)
	}

	type DiscountInputData struct {
		Discount          float64 `json:"discount"`
		DiscountEndDate   string  `json:"discountEndDate"`
		DiscountStartDate string  `json:"discountStartDate"`
	}
	inputProduct := new(DiscountInputData)

	if err := c.BodyParser(inputProduct); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	// set dates
	discountStartDate, err := time.Parse("2006-01-02", inputProduct.DiscountStartDate)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	discountEndDate, err := time.Parse("2006-01-02", inputProduct.DiscountEndDate)
	if err != nil {
		fmt.Println(err)
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	// Prepare the update fields
	updates := bson.D{{"$set", bson.D{{"discount", inputProduct.Discount}, {"discountStartDate", discountStartDate}, {"discountEndDate", discountEndDate}}}}

	// build id for mongo
	id, _ := primitive.ObjectIDFromHex(productID)

	// Create a filter to find the product by ID
	filter := bson.D{{"_id", id}}

	opts := options.Update().SetUpsert(true)
	result, err := mgm.Coll(&Product{}).UpdateOne(context.TODO(), filter, updates, opts)
	if err != nil {
		fmt.Println("Error ", err.Error())
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"product": inputProduct})
	}
	fmt.Printf("Number of documents updated: %v\n", result.ModifiedCount)
	fmt.Printf("Number of documents upserted: %v\n", result.UpsertedCount)

	return c.Status(http.StatusCreated).JSON(fiber.Map{"updated": true})
}

func Get(c *fiber.Ctx) error {
	id := c.Params("id")
	fmt.Println("Find product", id)
	product := new(Product)

	if err := mgm.Coll(product).FindByID(id, product); err != nil {
		return fiber.NewError(http.StatusNotFound, "Not Found")
	}

	return c.JSON(product)
}

// get all products
func GetAll(c *fiber.Ctx) error {
	products := []Product{}
	category := c.Params("category", "")
	subCategory := c.Params("subCategory", "")
	name := c.Params("name", "")
	limit, limitErr := c.ParamsInt("limit", 20)
	offset, offsetErr := c.ParamsInt("offset", 0)

	if limitErr != nil || offsetErr != nil {
		return fiber.NewError(http.StatusBadRequest, "Arguments Error")
	}

	// filters := bson.M{}
	filters := map[string]string{}

	if category != "" {
		filters["category"] = category
	}
	if subCategory != "" {
		filters["subCategory"] = subCategory
	}

	if name != "" {
		filters["name"] = name
	}

	// Define options for the query
	findOptions := options.Find()

	findOptions.SetLimit(int64(limit))
	findOptions.SetSkip(int64(offset))

	// if err := mgm.Coll(&Product{}).SimpleFind(&products, filters, findOptions); err != nil {
	// 	return fiber.NewError(http.StatusBadRequest, err.Error())
	// }

	cursor, err := mgm.Coll(&Product{}).Find(mgm.Ctx(), filters, findOptions)

	if err := cursor.All(mgm.Ctx(), &products); err != nil {
		fmt.Println("An Error in Fetch Cursor")
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(products)
}
