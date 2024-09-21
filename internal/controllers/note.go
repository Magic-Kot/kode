package controllers

import (
	"fmt"
	"net/http"

	"github.com/Magic-Kot/kode/internal/models"
	"github.com/Magic-Kot/kode/internal/services/note"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/rs/zerolog"
)

type ApiNoteController struct {
	NoteService note.NoteService
	logger      *zerolog.Logger
	validator   *validator.Validate
}

func NewApiNoteController(noteService *note.NoteService, logger *zerolog.Logger, validator *validator.Validate) *ApiNoteController {
	return &ApiNoteController{
		NoteService: *noteService,
		logger:      logger,
		validator:   validator,
	}
}

func (an ApiNoteController) CreateNote(c echo.Context) error {
	ctx := c.Request().Context()
	ctx = an.logger.WithContext(ctx)

	an.logger.Debug().Msg("starting the handler 'CreateNote'")

	id := c.Get("id")

	userId, ok := id.(string)
	if ok != true {
		an.logger.Debug().Msgf("invalid id: %v", id)
		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	var reqInp models.TextsInput
	if err := c.Bind(&reqInp); err != nil {
		an.logger.Debug().Msgf("bind: invalid request: %v", err)

		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid request"))
	}

	err := an.NoteService.CreateNote(ctx, &userId, reqInp)
	if err != nil {
		an.logger.Debug().Msgf("create note: %v", err)
		return c.JSON(http.StatusInternalServerError, fmt.Sprint(err))
	}

	return c.JSON(http.StatusCreated, fmt.Sprint("successfully created note"))
}

func (an ApiNoteController) GetAllNote(c echo.Context) error {
	ctx := c.Request().Context()
	ctx = an.logger.WithContext(ctx)

	an.logger.Debug().Msg("starting the handler 'GetNote'")

	id := c.Get("id")

	userId, ok := id.(string)
	if ok != true {
		an.logger.Debug().Msgf("invalid id: %v", id)
		return c.JSON(http.StatusBadRequest, fmt.Sprint("invalid id"))
	}

	res, err := an.NoteService.GetAllNote(ctx, userId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, fmt.Sprint(err))
	}

	return c.JSON(http.StatusOK, res)
}
