package controllers

import (
	"fmt"
	"os"
	"server/config"
	"server/models"
	"server/security"
	"server/util"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	JwtSecretKey     = []byte(os.Getenv("JWT_SECRET_KEY"))
	JwtSigningMethod = jwt.SigningMethodHS256.Name
)

func Login(c *fiber.Ctx) error {
	userCollection := config.MI.DB.Collection("users")

	// Create a creds object with the User structure
	creds := new(models.User)

	// Parse and obtain provided credentials
	err := c.BodyParser(&creds)

	// if error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse JSON",
			"error":   err.Error(),
		})
	}

	// Query user with provided email
	user := &models.User{}
	query := bson.D{{Key: "email", Value: creds.Email}}

	err = userCollection.FindOne(c.Context(), query).Decode(user)
	if (err != nil) || (user.Password != "" && !util.CheckPasswordHash(creds.Password, user.Password)) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Login failed",
		})
	}

	claims := jwt.StandardClaims{
		Id:        user.ID,
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(time.Hour * 48).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Encode token
	encodedToken, err := token.SignedString(JwtSecretKey)

	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user":  util.GetSafeUser(user),
		"token": encodedToken,
	})

}

func GetByToken(c *fiber.Ctx) error {
	userCollection := config.MI.DB.Collection("users")
	token := util.ExtractToken(c)

	fmt.Println("token: ", token)

	claims, err := security.ParseToken(token)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid token",
			"error":   err.Error(),
		})
	}

	fmt.Println("claims.Id: ", claims.Id)

	// convert claims.Id to objectId
	id, err := primitive.ObjectIDFromHex(claims.Id)

	// if error while parsing claims.Id
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Something went wrong",
			"error":   err,
		})
	}

	// find user and return
	user := &models.User{}
	query := bson.D{{Key: "_id", Value: id}}

	err = userCollection.FindOne(c.Context(), query).Decode(user)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": util.GetSafeUser(user),
	})
}

// GetAll : get all users
func GetAll(c *fiber.Ctx) error {
	userCollection := config.MI.DB.Collection("users")

	// Query to filter
	query := bson.D{{}}
	projection := options.Find().SetProjection(bson.D{{"password", 0}})

	cursor, err := userCollection.Find(c.Context(), query, projection)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
	}

	var users []models.User = make([]models.User, 0)

	// iterate the cursor and decode each item into a User
	err = cursor.All(c.Context(), &users)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"users": users,
	})
}

// GetById : get user by their id
// PARAM: id
func GetById(c *fiber.Ctx) error {
	userCollection := config.MI.DB.Collection("users")

	// get parameter value
	paramID := c.Params("id")

	// convert paramID to objectId
	id, err := primitive.ObjectIDFromHex(paramID)

	// if error while parsing paramID
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse id",
			"error":   err,
		})
	}

	// find user and return
	user := &models.User{}
	query := bson.D{{Key: "_id", Value: id}}

	err = userCollection.FindOne(c.Context(), query).Decode(user)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "User not found",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": util.GetSafeUser(user),
	})
}

// Create : Create a new user
func Create(c *fiber.Ctx) error {
	userCollection := config.MI.DB.Collection("users")

	data := new(models.User)

	// Validation
	err := c.BodyParser(&data)

	// if error
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse JSON",
			"error":   err.Error(),
		})
	}

	if len(data.Password) == 0 {
		data.Password = "defaultpassword"
	}
	hashedPassword, err := util.HashPassword(data.Password)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
	}

	data.Password = hashedPassword
	data.CreatedAt = time.Now()
	data.UpdatedAt = time.Now()

	result, err := userCollection.InsertOne(c.Context(), data)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot insert user",
			"error":   err.Error(),
		})
	}

	// get the inserted data
	user := &models.User{}
	query := bson.D{{Key: "_id", Value: result.InsertedID}}
	projection := options.FindOne().SetProjection(bson.D{{"password", 0}})

	userCollection.FindOne(c.Context(), query, projection).Decode(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Something went wrong",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"user": user,
	})
}

// UpdateById : Update user corresponding to provided id
// PARAM: id
func UpdateById(c *fiber.Ctx) error {
	userCollection := config.MI.DB.Collection("users")

	// find parameter
	paramID := c.Params("id")

	// convert paramID to objectId
	id, err := primitive.ObjectIDFromHex(paramID)

	// if parameter cannot parse
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse id",
			"error":   err.Error(),
		})
	}

	// var data Request
	data := new(models.User)
	err = c.BodyParser(&data)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse JSON",
			"error":   err.Error(),
		})
	}

	query := bson.D{{Key: "_id", Value: id}}

	// updateData
	var dataToUpdate bson.D

	if data.Name != "" {
		dataToUpdate = append(dataToUpdate, bson.E{Key: "name", Value: data.Name})
	}

	if data.Email != "" {
		dataToUpdate = append(dataToUpdate, bson.E{Key: "email", Value: data.Email})
	}

	if data.Title != "" {
		dataToUpdate = append(dataToUpdate, bson.E{Key: "title", Value: data.Title})
	}

	if !data.Birthdate.IsZero() {
		dataToUpdate = append(dataToUpdate, bson.E{Key: "birthdate", Value: data.Birthdate})
	}

	dataToUpdate = append(dataToUpdate, bson.E{Key: "updatedAt", Value: time.Now()})

	update := bson.D{
		{Key: "$set", Value: dataToUpdate},
	}

	// update
	err = userCollection.FindOneAndUpdate(c.Context(), query, update).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User not found",
				"error":   err,
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot update user",
			"error":   err,
		})
	}

	// get updated data
	user := &models.User{}
	userCollection.FindOne(c.Context(), query).Decode(user)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"user": util.GetSafeUser(user),
	})
}

// DeleteById : Delete user by their id
// PARAM: id
func DeleteById(c *fiber.Ctx) error {
	userCollection := config.MI.DB.Collection("users")

	// get param
	paramID := c.Params("id")

	// convert parameter to object id
	id, err := primitive.ObjectIDFromHex(paramID)

	// if parameter cannot parse
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot parse id",
			"error":   err.Error(),
		})
	}

	// find and delete todo
	query := bson.D{{Key: "_id", Value: id}}

	err = userCollection.FindOneAndDelete(c.Context(), query).Err()

	if err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "User Not found",
				"error":   err.Error(),
			})
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Cannot delete user",
			"error":   err.Error(),
		})
	}

	return c.SendStatus(fiber.StatusNoContent)
}
