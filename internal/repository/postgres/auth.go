package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Magic-Kot/code/internal/models"
	"github.com/Magic-Kot/code/pkg/client/postg"

	"github.com/rs/zerolog"
)

var (
// errUserNotFound = errors.New("user not found")
// errGetUser      = errors.New("failed to get user")
)

type AuthPostgresRepository struct {
	client postg.Client
}

func NewAuthPostgresRepository(client postg.Client) *AuthPostgresRepository {
	return &AuthPostgresRepository{
		client: client,
	}
}

// SignIn - user authentication, getting the user's id and password
func (r *AuthPostgresRepository) SignIn(ctx context.Context, user *models.UserAuthorization) (*models.UserAuthorization, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'SignIn' method")
	logger.Debug().Msgf("postgres SignIn: by login: %s", user.Username)

	q := `
		SELECT id, password
		FROM users
		WHERE username = $1
	`

	err := r.client.QueryRowx(q, user.Username).Scan(&user.ID, &user.Password)

	if errors.Is(err, sql.ErrNoRows) {
		logger.Debug().Msgf("user not found: %s", err)
		return nil, errUserNotFound
	} else if err != nil {
		logger.Debug().Msgf("failed to get user: %s", err)
		return nil, errGetUser
	}

	return user, nil
}
