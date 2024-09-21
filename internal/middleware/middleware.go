package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/Magic-Kot/kode/pkg/utils/jwt_token"

	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

var errAuthorizationUser = fmt.Sprint("an unauthorized user")

type Middleware struct {
	logger *zerolog.Logger
	token  *jwt_token.Manager
}

func NewMiddleware(logger *zerolog.Logger, token *jwt_token.Manager) *Middleware {
	return &Middleware{
		logger: logger,
		token:  token,
	}
}

// AuthorizationUser - user authorization
func (m *Middleware) AuthorizationUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		header := c.Request().Header.Get("Authorization")

		if header == "" {
			m.logger.Debug().Msgf("empty 'Authorization' header: %v", header)

			return c.JSON(http.StatusUnauthorized, errAuthorizationUser)
		}

		headerParts := strings.Split(header, " ")

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			m.logger.Debug().Msgf("invalid 'Authorization' header: %v", header)

			return c.JSON(http.StatusUnauthorized, errAuthorizationUser)
		}

		id, err := m.token.ParseToken(headerParts[1])
		if err != nil {
			m.logger.Debug().Msgf("invalid authorization token: %v", err)

			return c.JSON(http.StatusUnauthorized, errAuthorizationUser)
		}

		c.Set("id", id)

		err = next(c)
		if err != nil {
			m.logger.Warn().Msgf("next HandlerFunc: %v", err)

			return c.JSON(http.StatusBadRequest, err.Error())
		}

		return nil
	}
}
