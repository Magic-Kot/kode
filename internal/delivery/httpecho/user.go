package httpecho

import (
	"github.com/Magic-Kot/code/internal/controllers"
	"github.com/Magic-Kot/code/internal/middleware"

	"github.com/labstack/echo/v4"
)

func SetUserRoutes(e *echo.Echo, apiController *controllers.ApiController, middleware *middleware.Middleware) {
	e.POST("/sign-up", apiController.SignUp)

	user := e.Group("/user", middleware.AuthorizationUser)
	{
		user.GET("/get", apiController.GetUser)
		user.PUT("/update", apiController.UpdateUser)
		user.DELETE("/delete", apiController.DeleteUser)
	}
}
