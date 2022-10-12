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

func GetOrders() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		result, err := orderCollection.Find(ctx, bson.M{})
		if err != nil {
			defer cancel()
			return c.JSON(echo.ErrInternalServerError.Code, "error occured while listing order items")
		}
		var allOrders []bson.M
		if err = result.All(ctx, &allOrders); err != nil {
			defer cancel()
			log.Fatal(err)
		}
		defer cancel()
		return c.JSON(http.StatusOK, allOrders)

	}
}

func GetOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		orderId := c.Param("order_id")
		var order model.Order
		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "error while fetching the order item")
		}
		defer cancel()
		return c.JSON(http.StatusOK, order)
	}
}

func CreateOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var table model.Table
		var order model.Order
		if err := c.Bind(&order); err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		err := tableCollection.FindOne(ctx, bson.M{"table_id": order.Table_id}).Decode(&table)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "table was not found")
		}
		order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.ID = primitive.NewObjectID() // try
		order.Order_id = order.ID.Hex()
		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		validationError := validate.Struct(order)
		if validationError != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, validationError.Error())
		}
		result, insertErr := orderCollection.InsertOne(ctx, &order)
		if insertErr != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "order item was not created")
		}
		defer cancel()
		return c.JSON(http.StatusCreated, result)
	}
}

func UpdateOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var order model.Order
		orderId := c.Param("order_id")
		err := orderCollection.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&order)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "this order does not exist")
		}
		if err := c.Bind(&order); err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		validationError := validate.Struct(order)
		if validationError != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, validationError.Error())
		}
		result, updateErr := orderCollection.UpdateOne(ctx, bson.M{"order_id": orderId}, bson.D{{"$set", &order}})
		if updateErr != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, updateErr.Error())
		}
		defer cancel()
		return c.JSON(http.StatusAccepted, result)

	}
}

func OrderItemOrderCreator(order model.Order) string {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	order.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	order.ID = primitive.NewObjectID() // try
	order.Order_id = order.ID.Hex()
	orderCollection.InsertOne(ctx, &order)
	defer cancel()
	return order.Order_id
}
