package route

import (
	"rms/controller"

	"github.com/labstack/echo/v4"
)

func MenuRoute(g *echo.Group) {
	g.GET("/", controller.GetMenus())
	g.GET("/:menu_id", controller.GetMenu())
	g.POST("/", controller.CreateMenu())
	g.PATCH("/:menu_id", controller.UpdateMenu())

}
