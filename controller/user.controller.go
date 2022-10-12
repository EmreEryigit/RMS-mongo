package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"rms/helper"
	"rms/model"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

		recordPerPage, err := strconv.Atoi(c.QueryParam("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}

		page, err1 := strconv.Atoi(c.QueryParam("page"))
		if err1 != nil || page < 1 {
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.QueryParam("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		projectStage := bson.D{
			{"$project", bson.D{
				{"_id", 0},
				{"total_count", 1},
				{"user_items", bson.D{{"$slice", []interface{}{"$data", startIndex, recordPerPage}}}},
			}}}

		result, err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, "error occured while listing user items")
		}

		var allUsers []bson.M
		if err = result.All(ctx, &allUsers); err != nil {
			log.Fatal(err)
		}
		return c.JSON(http.StatusOK, allUsers[0])
	}
}

func GetUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		userId := c.QueryParam("user_id")
		var user model.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": userId}).Decode(&user)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "error occured while listing users")
		}
		defer cancel()
		return c.JSON(http.StatusOK, user)
	}
}

func Signup() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		// first initialize private user for password validation
		var userPrivate model.UserPrivate
		if err := c.Bind(&userPrivate); err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "invalid request")
		}
		// validate
		validationError := validate.Struct(userPrivate)
		if validationError != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, validationError.Error())
		}
		count, err := userCollection.CountDocuments(ctx, bson.M{"email": userPrivate.Email})
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "user does not exist")
		}
		if count > 0 {
			defer cancel()
			return c.JSON(http.StatusConflict, "email already taken")
		}
		userPrivate.HashPassword()
		userPrivate.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		userPrivate.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		userPrivate.ID = primitive.NewObjectID()
		userPrivate.User_id = userPrivate.ID.Hex()
		user := userPrivate.User
		_, err = userCollection.InsertOne(ctx, user)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "could not save user")
		}

		jwtToken, err := helper.GenerateJWT(userPrivate.User_id, *userPrivate.First_name, *userPrivate.Email)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "error while generating jwt token")
		}
		session, _ := Store.Get(c.Request(), "auth-session")
		session.Values["auth"] = jwtToken
		err = session.Save(c.Request(), c.Response())
		if err != nil {
			defer cancel()
			c.JSON(http.StatusInternalServerError, "error while generating jwt token")
			return err
		}
		defer cancel()
		return c.JSON(http.StatusOK, user)
	}
}

func Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		var user model.UserPrivate
		var foundUser model.User
		if err := c.Bind(&user); err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "user does not exist")
		}
		if foundUser.Email == nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "user not found")
		}
		isValid := foundUser.VerifyPassword(*user.Password)
		if !isValid {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "invalid email or password")
		}
		token, err := helper.GenerateJWT(fmt.Sprint(foundUser.ID), *foundUser.First_name, *foundUser.Email)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "error while generating token")
		}
		session, _ := Store.Get(c.Request(), "auth-session")
		session.Values["auth"] = token
		err1 := session.Save(c.Request(), c.Response())
		if err1 != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "could not save the cookie")
		}
		defer cancel()
		return c.JSON(http.StatusOK, foundUser)
	}
}

func Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		session, err := Store.Get(c.Request(), "auth-session")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "error sessions")
		}
		session.Options.MaxAge = -1
		err = session.Save(c.Request(), c.Response().Writer)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "error saving session")
		}
		return err
	}
}

func WhoAmI() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		claims := c.Get("current-user").(*helper.SignedDetails)
		var user model.User
		err := userCollection.FindOne(ctx, bson.M{"user_id": claims.UserID}).Decode(&user)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "not logged in")
		}
		defer cancel()
		return c.JSON(http.StatusOK, user)
	}
}
