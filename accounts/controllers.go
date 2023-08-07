package accounts

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var jwtSigningKey = []byte(os.Getenv("JWT_SIGNING_KEY"))

type JWTCustomClaims struct {
	// NOTE: ID would coincide with ("jti" (JWT ID) Claim) which exists in jwt claims already
	AccountID string `json:"account_id"`
	Type      string `json:"type"`
	jwt.RegisteredClaims
}

func generateJwt(account *Account) (string, error) {

	// Create claims while leaving out some of the optional fields
	jwtClaims := JWTCustomClaims{
		AccountID: account.ID.Hex(),
		Type:      string(account.Type),
		RegisteredClaims: jwt.RegisteredClaims{
			// Also fixed dates can be used for the NumericDate
			// ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Walmart",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)

	// fmt.Println("os.Getenv(JWT_SECRET) ", os.Getenv("JWT_SECRET"))

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(jwtSigningKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func Create(c *fiber.Ctx) error {
	account := New()

	if err := c.BodyParser(account); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	// validate
	v := validator.New()
	validationResult := v.Struct(account)
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

	hashed, err := bcrypt.GenerateFromPassword([]byte(account.Password), 10)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{})
	}

	fmt.Println("hashed password: ", hashed)

	account.Password = string(hashed)

	if err := mgm.Coll(account).Create(account); err != nil {
		return c.Status(400).JSON(fiber.Map{})
	}

	tokenString, err := generateJwt(account)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{})
	}

	// set cookie too
	c.Cookie(&fiber.Cookie{
		Name:     "WalmartToken",
		Value:    tokenString,
		Expires:  time.Now().Add(24 * time.Hour * 30),
		HTTPOnly: false, // for testing purposes
		SameSite: "lax",
	})

	return c.Status(http.StatusCreated).JSON(fiber.Map{"account": account, "token": tokenString})
}

func Login(c *fiber.Ctx) error {
	inputAccount := new(Account)
	dbAccount := new(Account)

	if err := c.BodyParser(inputAccount); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	result := mgm.Coll(dbAccount).FindOne(context.TODO(), bson.M{"email": inputAccount.Email})

	if result.Err() != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{})
	}

	if err := result.Decode(dbAccount); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{})
	}

	// compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(dbAccount.Password), []byte(inputAccount.Password)); err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{})
	}

	tokenString, err := generateJwt(dbAccount)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{})
	}

	// set cookie too
	c.Cookie(&fiber.Cookie{
		Name:     "WalmartToken",
		Value:    string(tokenString),
		Expires:  time.Now().Add(24 * time.Hour * 30),
		HTTPOnly: false, // for testing purposes
		SameSite: "lax",
	})

	return c.Status(http.StatusOK).JSON(fiber.Map{"account": dbAccount, "token": tokenString})
}

func Get(c *fiber.Ctx) error {
	id := c.Params("id")
	fmt.Println("Find account", id)
	account := new(Account)

	if err := mgm.Coll(account).FindByID(id, account); err != nil {
		if err := mgm.Coll(account).FindOne(context.TODO(), bson.M{"id": id}); err != nil {
			return fiber.NewError(http.StatusNotFound, "Not Found")
		}
	}

	return c.JSON(account)
}

func GetAll(c *fiber.Ctx) error {
	accountList := []Account{}

	if err := mgm.Coll(&Account{}).SimpleFind(&accountList, bson.M{}); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(accountList)
}

func AccountExists(slug string) bool {
	account := new(Account)

	count, err := mgm.Coll(account).CountDocuments(context.TODO(), bson.M{"slug": slug})
	if err != nil {
		return false
	}

	return count > 0
}

func DecodeJwt(tokenString string) (*Account, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return jwtSigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		// fmt.Println("INFO :: claims")
		// fmt.Println(claims)
		// fmt.Println(claims["account_id"])
		// fmt.Println(claims["type"])

		accountID, err := primitive.ObjectIDFromHex(claims["account_id"].(string))

		if err != nil {
			fmt.Println("Error creating object id")
			fmt.Println(err)
			return nil, err
		}

		account := new(Account)
		account.ID = accountID
		account.Type = claims["type"].(string)

		return account, nil
	} else {
		fmt.Println("ERROR :: could not get claims")
		return new(Account), errors.New("Error Getting Claims")
	}
}
