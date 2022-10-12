package controller

import (
	"os"
	"rms/database"

	"github.com/go-playground/validator"
	"github.com/gorilla/sessions"
)

var foodCollection = database.OpenCollection(database.Client, "food")
var menuCollection = database.OpenCollection(database.Client, "menu")
var orderCollection = database.OpenCollection(database.Client, "order")
var tableCollection = database.OpenCollection(database.Client, "table")
var invoiceCollection = database.OpenCollection(database.Client, "invoice")
var orderItemCollection = database.OpenCollection(database.Client, "orderItem")
var userCollection = database.OpenCollection(database.Client, "user")
var validate = validator.New()
var Store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
