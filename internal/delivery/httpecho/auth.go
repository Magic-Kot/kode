package httpecho

import (
	"github.com/Magic-Kot/code/internal/controllers"

	"github.com/labstack/echo/v4"
)

func SetAuthRoutes(e *echo.Echo, apiController *controllers.ApiAuthController) {
	auth := e.Group("/auth")
	{
		auth.POST("/sign-in", apiController.SignIn)
		auth.POST("/refresh", apiController.RefreshToken)
	}
}
