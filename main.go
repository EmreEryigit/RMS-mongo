package main

import (
	"log"
	"os"
	"rms/middleware"
	"rms/route"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file")
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	r := echo.New()

	userG := r.Group("/users")
	foodG := r.Group("/foods")
	invoiceG := r.Group("/invoices")
	menuG := r.Group("/menus")
	orderG := r.Group("/orders")
	tableG := r.Group("/tables")
	orderItemG := r.Group("/orderItems")

	r.Use(middleware.CurrentUser)

	route.UserRoute(userG)

	r.Use(middleware.Authenticate)

	route.FoodRoute(foodG)
	route.InvoiceRoute(invoiceG)
	route.MenuRoute(menuG)
	route.OrderRoute(orderG)
	route.TableRoute(tableG)
	route.OrderItemRoute(orderItemG)
	r.Logger.Fatal(r.Start(":" + port))
}
