package accounts

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
)

var jwtSigningKey = []byte(os.Getenv("JWT_SIGNING_KEY"))

type JWTCustomClaims struct {
	Account string `json:"account"`
	jwt.RegisteredClaims
}

func generateJwt(account *Account) (string, error) {
	account.Password = ""
	accountJson, err := json.Marshal(account)
	if err != nil {
		return "", err
	}

	// Create claims while leaving out some of the optional fields
	jwtClaims := JWTCustomClaims{
		string(accountJson),
		jwt.RegisteredClaims{
			// Also fixed dates can be used for the NumericDate
			// ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 30)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Messenger",
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
	account := new(Account)

	if err := c.BodyParser(account); err != nil {
		return fiber.NewError(http.StatusBadRequest, err.Error())
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
		Name:     "MessengerToken",
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
		Name:     "MessengerToken",
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
		if err := mgm.Coll(account).FindOne(context.TODO(), bson.M{"slug": id}); err != nil {
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
