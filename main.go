package main

import (
	"log"
	"os"
	"rms/database"
	"rms/route"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

var foodCollection *mongo.Collection = database.OpenCollection(database.Client, "food")

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
	port := os.Getenv("PORT")
	/* if port == "" {
		port = "3000"
	} */
	r := echo.New()

	// userG := r.Group("/users")
	foodG := r.Group("/foods")
	// invoiceG := r.Group("/invoices")
	menuG := r.Group("/menus")
	// orderG := r.Group("/orders")
	// tableG := r.Group("/tables")
	// orderItemG := r.Group("/orderItems")

	// route.UserRoute(userG)
	route.FoodRoute(foodG)
	// route.InvoiceRoute(invoiceG)
	route.MenuRoute(menuG)
	// route.OrderRoute(orderG)
	// route.TableRoute(tableG)
	// route.OrderItemRoute(orderItemG)
	r.Logger.Fatal(r.Start(":" + port))
}
