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

func GetInvoices() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		result, err := invoiceCollection.Find(ctx, bson.M{})
		if err != nil {
			defer cancel()
			return c.JSON(echo.ErrInternalServerError.Code, "error occured while listing invoice items")
		}
		var allInvoices []bson.M
		if err = result.All(ctx, &allInvoices); err != nil {
			defer cancel()
			log.Fatal(err)
		}
		defer cancel()
		return c.JSON(http.StatusOK, allInvoices)
	}
}

func GetInvoice() echo.HandlerFunc {
	return func(c echo.Context) error {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		invoiceId := c.Param("invoice_id")

		var invoice model.Invoice

		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, "error occured while listing invoice item")
		}

		var invoiceView model.InvoiceViewFormat

		allOrderItems, err := ItemsByOrder(invoice.Order_id)
		invoiceView.Order_id = invoice.Order_id
		invoiceView.Payment_due_date = invoice.Payment_due_date

		invoiceView.Payment_method = "null"
		if invoice.Payment_method != nil {
			invoiceView.Payment_method = *invoice.Payment_method
		}

		invoiceView.Invoice_id = invoice.Invoice_id
		invoiceView.Payment_status = *&invoice.Payment_status
		invoiceView.Payment_due = allOrderItems[0]["payment_due"]
		invoiceView.Table_number = allOrderItems[0]["table_number"]
		invoiceView.Order_details = allOrderItems[0]["order_items"]

		return c.JSON(http.StatusOK, invoiceView)
	}
}

func CreateInvoice() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var invoice model.Invoice
		if err := c.Bind(&invoice); err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, err.Error())
		}
		var order model.Order
		err := orderCollection.FindOne(ctx, bson.M{"order_id": invoice.Order_id}).Decode(&order)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "this order does not exist")
		}
		status := "PENDING"
		if invoice.Payment_status == nil {
			invoice.Payment_status = &status
		}
		invoice.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		invoice.ID = primitive.NewObjectID() // try
		invoice.Invoice_id = invoice.ID.Hex()
		validationError := validate.Struct(invoice)
		if validationError != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, validationError.Error())
		}
		result, insertErr := invoiceCollection.InsertOne(ctx, &invoice)
		if insertErr != nil {
			defer cancel()
			return c.JSON(http.StatusInternalServerError, "invoice item was not created")
		}
		defer cancel()
		return c.JSON(http.StatusCreated, result)
	}
}

func UpdateInvoice() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var invoice model.Invoice
		invoiceId := c.Param("invoice_id")
		err := invoiceCollection.FindOne(ctx, bson.M{"invoice_id": invoiceId}).Decode(&invoice)
		if err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, "this invoice does not exist")
		}
		if err := c.Bind(&invoice); err != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		invoice.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		status := "PENDING"
		if invoice.Payment_status == nil {
			invoice.Payment_status = &status
		}
		validationError := validate.Struct(invoice)
		if validationError != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, validationError.Error())
		}
		result, updateErr := invoiceCollection.UpdateOne(ctx, bson.M{"invoice_id": invoiceId}, bson.D{{"$set", &invoice}})
		if updateErr != nil {
			defer cancel()
			return c.JSON(http.StatusBadRequest, updateErr.Error())
		}
		defer cancel()
		return c.JSON(http.StatusAccepted, result)
	}
}
