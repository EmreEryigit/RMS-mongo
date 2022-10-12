package route

import (
	"rms/controller"
	"rms/middleware"

	"github.com/labstack/echo/v4"
)

func UserRoute(g *echo.Group) {
	g.POST("/signup", controller.Signup())
	g.POST("/login", controller.Login())
	g.Use(middleware.Authenticate)
	g.GET("/", controller.GetUsers())
	g.GET("/:user_id", controller.GetUser())
	g.POST("/logout", controller.Logout())
	g.GET("/whoami", controller.WhoAmI())
}
