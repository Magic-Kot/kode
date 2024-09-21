package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Magic-Kot/kode/internal/models"
	"github.com/Magic-Kot/kode/internal/services/user"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type ApiController struct {
	UserService user.UserService
	logger      *zerolog.Logger
	validator   *validator.Validate
}

func NewApiController(userService *user.UserService, logger *zerolog.Logger, validator *validator.Validate) *ApiController {
	return &ApiController{
		UserService: *userService,
		logger:      logger,
		validator:   validator,
	}
}

// SignUp - registering a new user
func (ac *ApiController) SignUp(c echo.Context) error {
	ctx := c.Request().Context()
	ctx = ac.logger.WithContext(ctx)

	ac.logger.Debug().Msg("starting the handler 'SignUp'")

	req := new(models.UserLogin)
	if err := c.Bind(req); err != nil {
		ac.logger.Warn().Msgf("bind: invalid request: %v", err)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid request"))
	}

	err := ac.validator.Struct(req)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.StructField() == "Username" {
				switch err.Tag() {
				case "required":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("Enter your login"))
				case "min":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The minimum login length is 4 characters"))
				case "max":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The maximum login length is 20 characters"))
				}
			}

			if err.StructField() == "Password" {
				switch err.Tag() {
				case "required":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("Enter your password"))
				case "min":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The minimum password length is 6 characters"))
				case "max":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("The maximum password length is 20 characters"))
				}
			}
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	id, err := ac.UserService.SignUp(ctx, req.Username, req.Password)
	if err != nil {
		// нет обработки ошибки на уникальность логина
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("successfully created user, id: %d", id))
}

// GetUser - getting a user by id
func (ac *ApiController) GetUser(c echo.Context) error {
	ctx := c.Request().Context()
	ctx = ac.logger.WithContext(ctx)

	ac.logger.Debug().Msg("starting the handler 'GetUser'")

	req := new(models.User)
	id := c.Get("id")

	userID, ok := id.(string)
	if ok != true {
		ac.logger.Debug().Msgf("invalid id: %v", id)
		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	userIdInt, err := strconv.Atoi(userID)
	if err != nil {
		ac.logger.Debug().Msgf("invalid id: %v", id)
		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	req.ID = userIdInt

	result, err := ac.UserService.GetUser(ctx, req)
	if err != nil {
		ac.logger.Debug().Msgf("error receiving user data: %v", err)
		return c.JSON(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, result)
}

// UpdateUser - updating user data by ID
func (ac *ApiController) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()
	ctx = ac.logger.WithContext(ctx)

	ac.logger.Debug().Msgf("starting the handler 'UpdateUser'")

	//req := new(models.User)
	var req models.User
	if err := c.Bind(&req); err != nil {
		ac.logger.Warn().Msgf("bind: invalid request: %v", err)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid request"))
	}

	id := c.Get("id")
	userId, ok := id.(string)
	if ok != true {
		ac.logger.Debug().Msgf("updateUser: invalid id: %d", id)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}
	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		ac.logger.Debug().Msgf("updateUser: invalid id: %d", id)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	req.ID = userIdInt

	err = ac.validator.Struct(&req)

	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			if err.StructField() == "ID" && err.Value() != "" {
				return c.JSON(http.StatusBadRequest, fmt.Sprintf("incorrect id"))
			} else if err.StructField() == "Age" && err.Value() != "" {
				return c.JSON(http.StatusBadRequest, fmt.Sprintf("incorrect age"))
			} else if err.StructField() == "Username" && err.Value() != "" {
				switch err.Tag() {
				case "min":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("the minimum login length is 4 characters"))
				case "max":
					return c.JSON(http.StatusBadRequest, fmt.Sprintf("the maximum login length is 20 characters"))
				}
			} else if err.StructField() == "Name" && err.Value() != "" {
				return c.JSON(http.StatusBadRequest, fmt.Sprintf("incorrect name"))
			} else if err.StructField() == "Surname" && err.Value() != "" {
				return c.JSON(http.StatusBadRequest, fmt.Sprintf("incorrect surname"))
			} else if err.StructField() == "Email" && err.Value() != "" {
				return c.JSON(http.StatusBadRequest, fmt.Sprintf("incorrect email"))
			}
			return c.JSON(http.StatusBadRequest, err.Error())
		}
	}

	err = ac.UserService.UpdateUser(ctx, &req)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprint("successfully updated"))
}

// DeleteUser - deleting a user by id
func (ac *ApiController) DeleteUser(c echo.Context) error {
	ctx := c.Request().Context()
	ctx = ac.logger.WithContext(ctx)

	ac.logger.Debug().Msgf("starting the handler 'DeleteUser'")

	id := c.Get("id")
	userId, ok := id.(string)
	if ok != true {
		ac.logger.Debug().Msgf("DeleteUser: invalid id: %d", id)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	userIdInt, err := strconv.Atoi(userId)
	if err != nil {
		ac.logger.Debug().Msgf("updateUser: invalid id: %d", id)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	err = ac.UserService.DeleteUser(ctx, userIdInt)
	if err != nil {
		return c.JSON(http.StatusNotFound, err.Error())
	}

	return c.JSON(http.StatusOK, fmt.Sprintf("successfully deleted user: %d", userIdInt))
}
