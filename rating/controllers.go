package rating

import (
	"context"
	"fmt"
	"net/http"

	"github.com/efenstakes/walmart-api-g/accounts"
	"github.com/efenstakes/walmart-api-g/products"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// add a rating
// TODO !! check if user has bought the product before
func Rate(c *fiber.Ctx) error {
	accountLocal := c.Locals("account")
	if accountLocal == nil {
		return fiber.NewError(http.StatusUnauthorized)
	}

	account := accountLocal.(*accounts.Account)

	// get input data
	inputData := new(ProductRating)
	if err := c.BodyParser(inputData); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}
	inputData.UserId = account.ID.Hex()

	// validate
	v := validator.New()
	validationResult := v.Struct(inputData)
	fmt.Println("validationResult ", validationResult.Error())
	validationErr, ok := validationResult.(validator.ValidationErrors)

	if !ok {
		return fiber.NewError(http.StatusBadRequest, "Error")
	}

	if len(validationErr) > 0 {
		fmt.Println(validationErr)

		errors := make(map[string]string)
		for _, vErr := range validationErr {
			fmt.Printf("'%s' has a value of '%v' which does not satisfy '%s'.\n", vErr.Field(), vErr.Value(), vErr.Tag())
			errors[vErr.Field()] = vErr.Tag()
		}

		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"errors": errors})
	}

	// get product object id
	id, _ := primitive.ObjectIDFromHex(inputData.ProductID)

	// get the product
	product := new(products.Product)
	if err := mgm.Coll(product).FindByID(id, product); err != nil {
		fmt.Println("error ", err.Error())
		return fiber.NewError(http.StatusNotFound)
	}

	// add a rating
	if err := mgm.Coll(inputData).Create(inputData); err != nil {
		return c.Status(400).JSON(fiber.Map{})
	}

	// get new product rating by
	// (allRatings * theRating) + newRatingNumber / ( allRatings + 1 )
	currentTotalRating := product.Rating.Rating * float64(product.Rating.NoOfRatings)
	fmt.Println("currentTotalRating ", currentTotalRating)
	newAddedRatingTotal := currentTotalRating + float64(inputData.Rating)
	fmt.Println("newAddedRatingTotal ", newAddedRatingTotal)
	newTotalRaters := float64(product.Rating.NoOfRatings + 1)
	fmt.Println("newTotalRaters ", newTotalRaters)
	newProductRating := newAddedRatingTotal / newTotalRaters
	fmt.Println("newProductRating ", newProductRating)
	// Prepare the update fields
	// updates := bson.D{{"$set", bson.D{{"ratingAnalytics.rating", newProductRating}, {"ratingAnalytics.noOfRatings", newProductRating}}}}
	updates := bson.D{{"$inc", bson.D{{"ratingAnalytics.noOfRatings", 1}}}, {"$set", bson.D{{"ratingAnalytics.rating", newProductRating}}}}

	// Create a filter to find the product by ID
	filter := bson.D{{"_id", id}}

	opts := options.Update().SetUpsert(true)
	result, err := mgm.Coll(product).UpdateOne(context.TODO(), filter, updates, opts)
	if err != nil {
		fmt.Println("Error ", err.Error())
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"updated": false})
	}
	fmt.Printf("Number of documents updated: %v\n", result.ModifiedCount)
	fmt.Printf("Number of documents upserted: %v\n", result.UpsertedCount)

	return c.Status(http.StatusCreated).JSON(fiber.Map{"saved": true, "rating": inputData})
}

// get product ratings

func GetAll(c *fiber.Ctx) error {
	ratings := []ProductRating{}
	productID := c.Params("id", "")

	if productID == "" {
		return fiber.NewError(http.StatusBadRequest, "Arguments Error")
	}

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
	filters["productId"] = productID

	cursor, err := mgm.Coll(&ProductRating{}).Find(mgm.Ctx(), filters, findOptions)

	if err := cursor.All(mgm.Ctx(), &ratings); err != nil {
		fmt.Println("An Error in Fetch Cursor")
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	if err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(ratings)
}

// delete rating
func Delete(c *fiber.Ctx) error {
	ID := c.Params("id", "")

	if ID == "" {
		return fiber.NewError(http.StatusBadRequest)
	}

	// get our account
	accountLocal := c.Locals("account")

	if accountLocal == nil {
		return fiber.NewError(http.StatusUnauthorized)
	}

	// convert account Interface{} to account type
	account := accountLocal.(*accounts.Account)

	ratingId, _ := primitive.ObjectIDFromHex(ID)
	fmt.Println("ID ", ID)
	fmt.Println("ratingId ", ratingId)
	fmt.Println("account.ID.Hex() ", account.ID.Hex())

	// delete
	deletedItem := new(ProductRating)
	err := mgm.Coll(deletedItem).FindOneAndDelete(mgm.Ctx(), bson.D{{"_id", ratingId}, {"userId", account.ID.Hex()}}).Decode(deletedItem)
	if err != nil {
		fmt.Println(err.Error())
		return c.Status(400).JSON(fiber.Map{"error": "Rating Not Found"})
	}

	/** recalculate the new product rating **/

	// get the product
	product := new(products.Product)

	if err := mgm.Coll(product).FindByID(deletedItem.ProductID, product); err != nil {
		return fiber.NewError(http.StatusNotFound, "Not Found")
	}

	// get current rating total
	currentRatingTotal := product.Rating.Rating * float64(product.Rating.NoOfRatings)
	newRatingTotal := (currentRatingTotal - deletedItem.Rating) / float64(product.Rating.NoOfRatings-1)
	updates := bson.D{{"$inc", bson.D{{"ratingAnalytics.noOfRatings", -1}}}, {"$set", bson.D{{"ratingAnalytics.rating", newRatingTotal}}}}

	// Create a filter to find the product by ID
	id, _ := primitive.ObjectIDFromHex(deletedItem.ProductID)
	filter := bson.D{{"_id", id}}

	opts := options.Update().SetUpsert(true)
	result, err := mgm.Coll(product).UpdateOne(context.TODO(), filter, updates, opts)
	if err != nil {
		fmt.Println("Error ", err.Error())
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"updated": false})
	}
	fmt.Printf("Number of documents updated: %v\n", result.ModifiedCount)
	fmt.Printf("Number of documents upserted: %v\n", result.UpsertedCount)

	return c.Status(http.StatusCreated).JSON(fiber.Map{"deleted": true})
}
