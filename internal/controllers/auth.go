package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Magic-Kot/kode/internal/models"
	"github.com/Magic-Kot/kode/internal/services/auth"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type ApiAuthController struct {
	UserService auth.AuthService
	logger      *zerolog.Logger
	validator   *validator.Validate
}

func NewApiAuthController(AuthService *auth.AuthService, logger *zerolog.Logger, validator *validator.Validate) *ApiAuthController {
	return &ApiAuthController{
		UserService: *AuthService,
		logger:      logger,
		validator:   validator,
	}
}

// SignIn - user authentication
func (an *ApiAuthController) SignIn(c echo.Context) error {
	ctx := c.Request().Context()
	ctx = an.logger.WithContext(ctx)

	an.logger.Debug().Msg("starting the handler 'SignIn'")

	req := new(models.UserAuthorization)

	// getting the client's IP address
	IPAddress := c.Request().Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = c.Request().Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = c.Request().RemoteAddr
	}
	req.IPAddress = IPAddress
	fmt.Printf("IP: %s\n", IPAddress)

	// getting the client's GUID
	req.GUID = c.QueryParam("GUID")

	if err := c.Bind(req); err != nil {
		an.logger.Debug().Msgf("bind: invalid request: %v", err)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid request"))
	}

	err := an.validator.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.StructField() == "Username" {
				switch err.Tag() {
				case "required":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("Enter your login"))
				case "min":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The minimum login length is 1 characters"))
				case "max":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The maximum login length is 20 characters"))
				}
			}

			if err.StructField() == "Password" {
				switch err.Tag() {
				case "required":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("Enter your password"))
				case "min":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The minimum password length is 1 characters"))
				case "max":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The maximum password length is 20 characters"))
				}
			}
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	tokens, err := an.UserService.SignIn(ctx, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	cookie := new(http.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = tokens.RefreshToken
	cookie.Path = "/auth"
	cookie.Expires = time.Now().Add(4 * time.Hour) // Hard code
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, tokens.AccessToken)
}

// RefreshToken - getting new refresh and access tokens
func (an *ApiAuthController) RefreshToken(c echo.Context) error {
	ctx := c.Request().Context()
	ctx = an.logger.WithContext(ctx)

	an.logger.Debug().Msg("starting the handler 'RefreshToken'")

	cookieRequest, err := c.Cookie("refreshToken")

	if err != nil || cookieRequest.Value == "" {
		return c.JSON(http.StatusUnauthorized, errors.New("invalid refresh token"))
	}

	tokens, err := an.UserService.RefreshToken(ctx, cookieRequest.Value)
	if err != nil {
		//return c.Redirect(http.StatusUnauthorized, "/sign-in")
		return c.JSON(http.StatusUnauthorized, err.Error())
	}

	cookie := new(http.Cookie)
	cookie.Name = "refreshToken"
	cookie.Value = tokens.RefreshToken
	cookie.Path = "/auth"
	cookie.Expires = time.Now().Add(4 * time.Hour)
	cookie.HttpOnly = true
	c.SetCookie(cookie)

	return c.JSON(http.StatusOK, tokens.AccessToken)
}
