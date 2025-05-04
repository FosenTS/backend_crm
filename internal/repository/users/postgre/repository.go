package postgre

import (
	"backend_crm/internal/model"
	"backend_crm/internal/repository/users"
	"context"
	"database/sql"
	"errors"
)

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) users.Repository {
	return &repository{db: db}
}

func (r *repository) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT user_id, role, username, pass_hash
		FROM users
		WHERE username = $1
	`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.UserId,
		&user.Role,
		&user.Username,
		&user.PassHash,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, users.ErrNotFoundUser
		}
		return nil, err
	}

	return &user, nil
}

func (r *repository) Save(ctx context.Context, register *model.Register) error {
	query := `
		INSERT INTO users (role, username, pass_hash)
		VALUES ($1, $2, $3)
	`

	_, err := r.db.ExecContext(ctx, query,
		register.RoleId,
		register.Username,
		register.Password,
	)

	return err
}
