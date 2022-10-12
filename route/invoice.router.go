package route

import (
	"rms/controller"

	"github.com/labstack/echo/v4"
)

func InvoiceRoute(g *echo.Group) {
	g.GET("/", controller.GetInvoices())
	g.GET("/:invoice_id", controller.GetInvoice())
	g.POST("/", controller.CreateInvoice())
	g.PATCH("/:invoice_id", controller.UpdateInvoice())

}
