package controller

import (
	"context"
	"log"
	"net/http"
	"rms/model"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetTables() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		result, err := tableCollection.Find(ctx, bson.M{})
		if err != nil {
			defer cancel()
			return c.JSON(echo.ErrInternalServerError.Code, "error occured while listing tables")
		}
		var allTables []bson.M
		if err = result.All(ctx, &allTables); err != nil {
			defer cancel()
			log.Fatal(err)
		}
		defer cancel()
		return c.JSON(http.StatusOK, allTables)
	}
}

func GetTable() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		tableId := c.Param("table_id")
		var table model.Table
		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "error while fetching the table")
		}
		defer cancel()
		return c.JSON(http.StatusOK, table)
	}
}

func CreateTable() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var table model.Table
		if err := c.Bind(&table); err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		table.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		table.ID = primitive.NewObjectID() // try
		table.Table_id = table.ID.Hex()
		validationErr := validate.Struct(table)
		if validationErr != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, validationErr.Error())
		}
		result, insertErr := tableCollection.InsertOne(ctx, &table)
		if insertErr != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "table was not created")
		}
		defer cancel()
		return c.JSON(http.StatusCreated, result)

	}
}

func UpdateTable() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var table model.Table
		tableId := c.Param("table_id")
		err := tableCollection.FindOne(ctx, bson.M{"table_id": tableId}).Decode(&table)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "this table does not exist")
		}
		if err := c.Bind(&table); err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		table.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		validationError := validate.Struct(table)
		if validationError != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, validationError.Error())
		}
		result, updateErr := tableCollection.UpdateOne(ctx, bson.M{"table_id": tableId}, bson.D{{"$set", &table}})
		if updateErr != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, updateErr.Error())
		}
		defer cancel()
		return c.JSON(http.StatusAccepted, result)
	}
}
