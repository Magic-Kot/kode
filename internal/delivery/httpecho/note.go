package httpecho

import (
	"github.com/Magic-Kot/code/internal/controllers"
	"github.com/Magic-Kot/code/internal/middleware"

	"github.com/labstack/echo/v4"
)

func SetNoteRoutes(e *echo.Echo, apiController *controllers.ApiNoteController, middleware *middleware.Middleware) {
	user := e.Group("/note", middleware.AuthorizationUser)
	{
		user.POST("/add", apiController.CreateNote)
		user.GET("/get", apiController.GetAllNote)
	}
}
