package userRepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"

	sq "github.com/Masterminds/squirrel"
	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type UserRepoPg struct {
	db *sqlx.DB
}

func NewPg(db *sqlx.DB) *UserRepoPg {
	return &UserRepoPg{
		db: db,
	}
}

func scanRows(rows *sql.Rows) []model.User {
	res := make([]model.User, 0, consts.SliceMinCap)
	for rows.Next() {
		var id uuid.UUID
		var email string
		var password string
		var role string

		err := rows.Scan(&id, &email, &password, &role)
		if err != nil {
			slog.Error(fmt.Sprintf("user rows scan error: %s", err.Error()))
			continue
		}

		user := model.User{
			Id:       id,
			Email:    email,
			Password: password,
			Role:     role,
		}

		res = append(res, user)
	}

	return res
}

func (r *UserRepoPg) AddUser(ctx context.Context, user model.User) (model.User, error) {
	rows, err := sq.Insert("users").
		Columns("user_id", "user_email", "user_passwd", "user_role").
		Values(user.Id, user.Email, user.Password, user.Role).
		Suffix("RETURNING user_id, user_email, user_passwd, user_role").
		RunWith(r.db).PlaceholderFormat(sq.Dollar).QueryContext(ctx)

	if err != nil {
		slog.Error(fmt.Sprintf("add user error: %s", err.Error()))
		return model.User{}, err
	}

	res := scanRows(rows)
	if len(res) != 1 {
		return model.User{}, errors.New("bad result len for add user")
	}

	return res[0], nil
}

func (r *UserRepoPg) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	rows, err := sq.Select("user_id", "user_email", "user_passwd", "user_role").
		From("users").Where(sq.Eq{"user_email": email}).
		RunWith(r.db).PlaceholderFormat(sq.Dollar).QueryContext(ctx)

	if err != nil {
		slog.Error(fmt.Sprintf("get user error: %s", err.Error()))
		return model.User{}, err
	}

	res := scanRows(rows)
	if len(res) != 1 {
		return model.User{}, errors.New("bad result len for get user")
	}

	return res[0], nil
}
