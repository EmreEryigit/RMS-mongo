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
	"go.mongodb.org/mongo-driver/mongo"
)

func GetOrderItems() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		result, err := orderItemCollection.Find(ctx, bson.M{})
		if err != nil {
			defer cancel()
			return c.JSON(echo.ErrInternalServerError.Code, "error occured while listing order-item items")
		}
		var allOrderItems []bson.M
		if err = result.All(ctx, &allOrderItems); err != nil {
			defer cancel()
			log.Fatal(err)
		}
		defer cancel()
		return c.JSON(http.StatusOK, allOrderItems)
	}
}

func GetOrderItemsByOrder() echo.HandlerFunc {
	return func(c echo.Context) error {
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		orderId := c.Param("order_id")

		allOrderItems, err := ItemsByOrder(orderId)
		if err != nil {
			defer cancel()
			return c.JSON(echo.ErrInternalServerError.Code, "error occured while listing order-item items")
		}
		return c.JSON(http.StatusOK, allOrderItems)
	}
}

// not a route
func ItemsByOrder(id string) (OrderItems []primitive.M, err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	matchStage := bson.D{{"$match", bson.D{{"order_id", id}}}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "food"}, {"localField", "food_id"}, {"foreignField", "food_id"}, {"as", "food"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$food"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupOrderStage := bson.D{{"$lookup", bson.D{{"from", "order"}, {"localField", "order_id"}, {"foreignField", "order_id"}, {"as", "order"}}}}
	unwindOrderStage := bson.D{{"$unwind", bson.D{{"path", "$order"}, {"preserveNullAndEmptyArrays", true}}}}

	lookupTableStage := bson.D{{"$lookup", bson.D{{"from", "table"}, {"localField", "order.table_id"}, {"foreignField", "table_id"}, {"as", "table"}}}}
	unwindTableStage := bson.D{{"$unwind", bson.D{{"path", "$table"}, {"preserveNullAndEmptyArrays", true}}}}

	projectStage := bson.D{
		{"$project", bson.D{
			{"id", 0},
			{"amount", "$food.price"},
			{"total_count", 1},
			{"food_name", "$food.name"},
			{"food_image", "$food.food_image"},
			{"table_number", "$table.table_number"},
			{"table_id", "$table.table_id"},
			{"order_id", "$order.order_id"},
			{"price", "$food.price"},
			{"quantity", 1},
		}}}

	groupStage := bson.D{{"$group", bson.D{{"_id", bson.D{{"order_id", "$order_id"}, {"table_id", "$table_id"}, {"table_number", "$table_number"}}}, {"payment_due", bson.D{{"$sum", "$amount"}}}, {"total_count", bson.D{{"$sum", 1}}}, {"order_items", bson.D{{"$push", "$$ROOT"}}}}}}

	projectStage2 := bson.D{
		{"$project", bson.D{

			{"id", 0},
			{"payment_due", 1},
			{"total_count", 1},
			{"table_number", "$_id.table_number"},
			{"order_items", 1},
		}}}

	result, err := orderItemCollection.Aggregate(ctx, mongo.Pipeline{
		matchStage,
		lookupStage,
		unwindStage,
		lookupOrderStage,
		unwindOrderStage,
		lookupTableStage,
		unwindTableStage,
		projectStage,
		groupStage,
		projectStage2})

	if err != nil {
		panic(err)
	}

	if err = result.All(ctx, &OrderItems); err != nil {
		panic(err)
	}

	defer cancel()

	return OrderItems, err

}

func GetOrderItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		orderItemId := c.Param("order_item_id")
		var orderItem model.OrderItem
		err := orderItemCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "error while fetching the order item")
		}
		defer cancel()
		return c.JSON(http.StatusOK, orderItem)
	}
}

func CreateOrderItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var orderItemPack model.OrderItemPack
		var order model.Order
		if err := c.Bind(&orderItemPack); err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		order.Order_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		order.Table_id = orderItemPack.Table_id
		order_id := OrderItemOrderCreator(order)
		var orderItemsToBeInserted []interface{}
		for _, orderItem := range orderItemPack.Order_items {
			orderItem.Order_id = order_id
			orderItem.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
			orderItem.ID = primitive.NewObjectID() // try
			orderItem.Order_item_id = orderItem.ID.Hex()
			num := toFixed(*orderItem.Unit_price, 2)
			orderItem.Unit_price = &num
			validationErr := validate.Struct(orderItem)
			if validationErr != nil {
				defer cancel()
				return c.JSON(echo.ErrBadRequest.Code, validationErr.Error())
			}
			orderItemsToBeInserted = append(orderItemsToBeInserted, &orderItem)
		}
		result, err := orderItemCollection.InsertMany(ctx, orderItemsToBeInserted)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, err.Error())
		}
		defer cancel()
		return c.JSON(http.StatusCreated, result)
	}
}

func UpdateOrderItem() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var orderItem model.OrderItem

		orderItemId := c.Param("order_item_id")
		err := orderItemCollection.FindOne(ctx, bson.M{"order_item_id": orderItemId}).Decode(&orderItem)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "this order item does not exist")
		}
		if err := c.Bind(&orderItem); err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		orderItem.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		validationError := validate.Struct(orderItem)
		if validationError != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, validationError.Error())
		}
		result, updateErr := orderItemCollection.UpdateOne(ctx, bson.M{"order_item_id": orderItemId}, bson.D{{"$set", &orderItem}})
		if updateErr != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, updateErr.Error())
		}
		defer cancel()
		return c.JSON(http.StatusAccepted, result)
	}
}
