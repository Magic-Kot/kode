package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Magic-Kot/code/internal/models"
	"github.com/Magic-Kot/code/pkg/client/postg"

	"github.com/rs/zerolog"
)

var (
	errTransaction = errors.New("transaction error")
	errGetAllNotes = errors.New("error getting all notes")
)

type NoteRepository struct {
	client postg.Client
}

func NewNoteRepository(client postg.Client) *NoteRepository {
	return &NoteRepository{
		client: client,
	}
}

// CreateNote - creating a new note
func (nr *NoteRepository) CreateNote(ctx context.Context, userId *string, texts string) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'CreateNote' method")
	logger.Debug().Msgf("postgres. userId: %s, texts: %v", *userId, texts)

	tx, err := nr.client.Begin()
	if err != nil {
		logger.Debug().Msgf("transaction creation error. err: %s", err)
		return errTransaction
	}

	var id int
	createNoteQuery := fmt.Sprint("INSERT INTO notes (texts) VALUES ($1) RETURNING id")
	row := tx.QueryRow(createNoteQuery, texts)
	if err = row.Scan(&id); err != nil {
		logger.Debug().Msgf("error writing to the 'notes' table. err: %s", err)

		tx.Rollback()
		return errTransaction
	}

	createUsersNotesQuery := fmt.Sprint("INSERT INTO users_notes (user_id, note_id) VALUES ($1, $2)")
	_, err = tx.Exec(createUsersNotesQuery, userId, id)
	if err != nil {
		logger.Debug().Msgf("error writing to the 'users_notes' table. err: %s", err)

		tx.Rollback()
		return errTransaction
	}

	return tx.Commit()
}

// GetAllNote - getting all the notes
func (nr *NoteRepository) GetAllNote(ctx context.Context, userId string) ([]models.NotesResponse, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'GetAllNote' method")
	logger.Debug().Msgf("postgres. userId: %s", userId)

	var notes []models.NotesResponse

	query := fmt.Sprint("SELECT n.id, n.texts FROM notes n INNER JOIN users_notes un on n.id = un.note_id WHERE un.user_id = $1")

	err := nr.client.Select(&notes, query, userId)
	if err != nil {
		logger.Debug().Msgf("error getting all notes. err: %s", err)
		return nil, errGetAllNotes
	}

	return notes, nil
}
