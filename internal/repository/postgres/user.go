package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/Magic-Kot/code/internal/models"
	"github.com/Magic-Kot/code/pkg/client/postg"

	"github.com/rs/zerolog"
)

var (
	errUserNotFound = errors.New("user not found")
	errCreateUser   = errors.New("failed to create user")
	errGetUser      = errors.New("failed to get user")
	errUpdateUser   = errors.New("failed to update user")
	errDeleteUser   = errors.New("failed to delete user")
)

type UserRepository struct {
	client postg.Client
}

func NewUserRepository(client postg.Client) *UserRepository {
	return &UserRepository{
		client: client,
	}
}

// CreateUser - creating a new user
func (r *UserRepository) CreateUser(ctx context.Context, username string, passwordHash string) (int, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'CreateUser' method")

	q := `
		INSERT INTO users 
		    (username, password) 
		VALUES 
		       ($1, $2) 
		RETURNING id
	`

	var id int

	if err := r.client.QueryRowx(q, username, passwordHash).Scan(&id); err != nil {
		logger.Debug().Msgf("failed to create user. %s", err)
		return 0, errCreateUser
	}

	return id, nil
}

// GetUser - getting a user by id
func (r *UserRepository) GetUser(ctx context.Context, user *models.User) (*models.User, error) {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'GetUser' method")
	logger.Debug().Msgf("postgres: get user by id: %d", user.ID)

	q := `
		SELECT username, name, surname, age, email
		FROM users
		WHERE id = $1
	`

	err := r.client.QueryRowx(q, user.ID).Scan(&user.Username, &user.Name, &user.Surname, &user.Age, &user.Email)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, errUserNotFound
	} else if err != nil {
		return nil, errGetUser
	}

	return user, nil
}

// UpdateUser - updating the data in the specified table by the specified column
func (r *UserRepository) UpdateUser(ctx context.Context, table string, column string, value string, arg []interface{}) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'UpdateUser' method")
	logger.Debug().Msgf("postgres: update table by table: %s, column: %s, value: %s, arg: %v", table, column, value, arg)

	q := fmt.Sprintf(`UPDATE %s SET %s WHERE %s = $1`, table, value, column)

	commandTag, err := r.client.Exec(q, arg...)

	if err != nil {
		logger.Debug().Msgf("failed table updates: %s", err)
		return errUpdateUser
	}

	if str, _ := commandTag.RowsAffected(); str != 1 {
		logger.Debug().Msgf("user not found: %s", err)
		return errUserNotFound
	}

	return nil
}

// DeleteUser - deleting a user from the 'users' table by id
func (r *UserRepository) DeleteUser(ctx context.Context, id int) error {
	logger := zerolog.Ctx(ctx)
	logger.Debug().Msg("accessing Postgres using the 'DeleteUser' method")

	q := `
		DELETE FROM users
		WHERE id = $1
	`

	commandTag, err := r.client.Exec(q, id)

	if err != nil {
		return errDeleteUser
	}

	if str, _ := commandTag.RowsAffected(); str != 1 {
		return errUserNotFound
	}

	return nil
}
