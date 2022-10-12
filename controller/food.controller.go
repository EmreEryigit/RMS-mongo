package controller

import (
	"context"
	"log"
	"math"
	"net/http"
	"rms/model"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetFoods() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		recordPerPage, err := strconv.Atoi(c.QueryParam("recordPerPage"))
		if err != nil || recordPerPage < 1 {
			recordPerPage = 10
		}
		page, err := strconv.Atoi(c.QueryParam("page"))
		if err != nil || page < 1 {
			page = 1
		}
		startIndex := (page - 1) * recordPerPage
		startIndex, _ = strconv.Atoi(c.QueryParam("startIndex"))

		matchStage := bson.D{{"$match", bson.D{{}}}}
		groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"_id", "null"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"data", bson.D{{"$push", "$$ROOT"}}}}}}
		projectStage := bson.D{
			{
				"$project", bson.D{
					{"$id", 0},
					{"total_count", 1},
					{"food_items", bson.D{{
						"$slice", []interface{}{"$data", startIndex, recordPerPage},
					}}},
				},
			},
		}

		result, err := foodCollection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage, projectStage})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, "error occured while listing food items")
		}
		var allFoods []bson.M
		if err = result.All(ctx, &allFoods); err != nil {
			log.Fatal(err)
		}
		return c.JSON(http.StatusOK, allFoods[0])
	}
}

func GetFood() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		foodId := c.Param("food_id")
		var food model.Food

		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		defer cancel()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, "error while fetching the food item")

		}
		return c.JSON(http.StatusOK, food)
	}
}

func CreateFood() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var menu model.Menu
		var food model.Food
		if err := c.Bind(&food); err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		validationError := validate.Struct(food)
		if validationError != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, validationError.Error())
		}
		err := menuCollection.FindOne(ctx, bson.M{"menu_id": food.Menu_id}).Decode(&menu)
		defer cancel()
		if err != nil {
			return c.JSON(http.StatusBadRequest, "menu was not found")
		}
		food.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		food.ID = primitive.NewObjectID() // try
		food.Food_id = food.ID.Hex()
		var num = toFixed(*food.Price, 2)
		food.Price = &num
		result, insertErr := foodCollection.InsertOne(ctx, &food)
		if insertErr != nil {
			return c.JSON(http.StatusInternalServerError, "food item was not created")
		}
		defer cancel()
		return c.JSON(http.StatusCreated, result)
	}
}

func UpdateFood() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var food model.Food

		foodId := c.Param("food_id")
		err := foodCollection.FindOne(ctx, bson.M{"food_id": foodId}).Decode(&food)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "this food does not exist")
		}
		if err := c.Bind(&food); err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		validationError := validate.Struct(food)
		if validationError != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, validationError.Error())
		}

		food.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		result, updateErr := foodCollection.UpdateOne(ctx, bson.M{"food_id": foodId}, bson.D{{"$set", &food}})
		if updateErr != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, updateErr.Error())
		}
		defer cancel()
		return c.JSON(http.StatusAccepted, result)
	}
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}
