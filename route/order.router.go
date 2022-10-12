package route

import (
	"rms/controller"

	"github.com/labstack/echo/v4"
)

func OrderRoute(g *echo.Group) {
	g.GET("/", controller.GetOrders())
	g.GET("/:order_id", controller.GetOrder())
	g.POST("/", controller.CreateOrder())
	g.PATCH("/:order_id", controller.UpdateOrder())

}
