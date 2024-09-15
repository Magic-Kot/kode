package note

import (
	"context"
	"strings"

	"github.com/Magic-Kot/code/internal/models"
	"github.com/Magic-Kot/code/pkg/speller"

	"github.com/rs/zerolog"
)

type NoteRepository interface {
	CreateNote(ctx context.Context, id *string, text string) error
	GetAllNote(ctx context.Context, userId string) ([]models.NotesResponse, error)
}

type NoteService struct {
	NoteRepository NoteRepository
	speller        *speller.Speller
}

func NewNoteService(noteRepository NoteRepository, speller *speller.Speller) *NoteService {
	return &NoteService{
		NoteRepository: noteRepository,
		speller:        speller,
	}
}

// CreateNote - creating a new note
func (ns *NoteService) CreateNote(ctx context.Context, id *string, req models.TextsInput) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the 'CreateNote' service")

	spellerResp, err := ns.speller.CheckTexts(ctx, req.Texts)
	if err != nil {
		return err
	}

	correctorTexts := make([]string, len(req.Texts))
	for i, text := range req.Texts {
		for _, misspell := range spellerResp[i] {
			if len(misspell.Suggestions) > 0 {
				text = strings.Replace(text, misspell.Word, misspell.Suggestions[0], -1)
			}
		}
		correctorTexts[i] = text
	}

	s := strings.Join(correctorTexts, " ")

	err = ns.NoteRepository.CreateNote(ctx, id, s)
	if err != nil {
		return err
	}

	return nil
}

// GetAllNote - getting all the notes
func (ns *NoteService) GetAllNote(ctx context.Context, userId string) ([]models.NotesResponse, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("starting the 'GetAllNote' service")

	res, err := ns.NoteRepository.GetAllNote(ctx, userId)
	if err != nil {
		return nil, err
	}

	return res, nil
}
