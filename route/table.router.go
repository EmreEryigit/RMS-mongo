package route

import (
	"rms/controller"

	"github.com/labstack/echo/v4"
)

func TableRoute(g *echo.Group) {
	g.GET("/", controller.GetTables())
	g.GET("/:table_id", controller.GetTable())
	g.POST("/", controller.CreateTable())
	g.PATCH("/:table_id", controller.UpdateTable())

}
