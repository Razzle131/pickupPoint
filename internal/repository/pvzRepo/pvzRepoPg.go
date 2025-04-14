package pvzRepo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/Razzle131/pickupPoint/internal/consts"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type PvzRepoPg struct {
	db *sqlx.DB
}

func NewPg(db *sqlx.DB) *PvzRepoPg {
	return &PvzRepoPg{
		db: db,
	}
}

func scanRows(rows *sql.Rows) []model.Pvz {
	res := make([]model.Pvz, 0, consts.SliceMinCap)
	for rows.Next() {
		var id uuid.UUID
		var city string
		var date time.Time

		err := rows.Scan(&id, &city, &date)
		if err != nil {
			slog.Error(fmt.Sprintf("pvz rows scan error: %s", err.Error()))
			continue
		}

		pvz := model.Pvz{
			Id:      id,
			City:    city,
			RegDate: date,
		}

		res = append(res, pvz)
	}

	return res
}

func (r *PvzRepoPg) AddPvz(ctx context.Context, pvz model.Pvz) (model.Pvz, error) {
	rows, err := sq.Insert("pvz").
		Columns("pvz_id", "pvz_city", "pvz_date").
		Values(pvz.Id, pvz.City, pvz.RegDate).
		Suffix("RETURNING pvz_id, pvz_city, pvz_date").
		RunWith(r.db).PlaceholderFormat(sq.Dollar).QueryContext(ctx)

	if err != nil {
		slog.Error(fmt.Sprintf("add pvz error: %s", err.Error()))
		return model.Pvz{}, err
	}

	res := scanRows(rows)
	if len(res) != 1 {
		slog.Error("bad result len for add pvz")
		return model.Pvz{}, errors.New("bad result len for add pvz")
	}

	return res[0], nil
}

func (r *PvzRepoPg) GetPvzById(ctx context.Context, pvzId uuid.UUID) (model.Pvz, error) {
	rows, err := sq.Select("pvz_id", "pvz_city", "pvz_date").
		From("pvz").
		Where(sq.Eq{"pvz_id": pvzId}).
		RunWith(r.db).PlaceholderFormat(sq.Dollar).QueryContext(ctx)

	if err != nil {
		slog.Error(fmt.Sprintf("get pvz error: %s", err.Error()))
		return model.Pvz{}, err
	}

	res := scanRows(rows)

	if len(res) != 1 {
		return model.Pvz{}, errors.New("bad result len for get pvz")
	}

	return res[0], nil
}

func (r *PvzRepoPg) ListPvz(ctx context.Context) ([]model.Pvz, error) {
	rows, err := sq.Select("pvz_id", "pvz_city", "pvz_date").
		From("pvz").
		RunWith(r.db).PlaceholderFormat(sq.Dollar).QueryContext(ctx)

	if err != nil {
		slog.Error(fmt.Sprintf("list pvz error: %s", err.Error()))
		return []model.Pvz{}, err
	}

	return scanRows(rows), nil
}
