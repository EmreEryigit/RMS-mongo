package route

import (
	"rms/controller"

	"github.com/labstack/echo/v4"
)

func FoodRoute(g *echo.Group) {
	g.GET("/", controller.GetFoods())
	g.GET("/:food_id", controller.GetFood())
	g.POST("/", controller.CreateFood())
	g.PATCH("/:food_id", controller.UpdateFood())
}
