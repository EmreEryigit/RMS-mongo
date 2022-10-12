package route

import (
	"rms/controller"

	"github.com/labstack/echo/v4"
)

func OrderItemRoute(g *echo.Group) {
	g.GET("/", controller.GetOrderItems())

	g.GET("/order/:order_id", controller.GetOrderItemsByOrder())

	g.GET("/:order_item_id", controller.GetOrderItem())
	g.POST("/", controller.CreateOrderItem())
	g.PATCH("/:order_item_id", controller.UpdateOrderItem())
}
