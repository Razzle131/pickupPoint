package receptionRepo

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

type ReceptionRepoPg struct {
	db *sqlx.DB
}

func NewPg(db *sqlx.DB) *ReceptionRepoPg {
	return &ReceptionRepoPg{
		db: db,
	}
}

func scanRows(rows *sql.Rows) []model.Reception {
	res := make([]model.Reception, 0, consts.SliceMinCap)
	for rows.Next() {
		var id uuid.UUID
		var date time.Time
		var pvzId uuid.UUID
		var status string

		err := rows.Scan(&id, &date, &status, &pvzId)
		if err != nil {
			slog.Error(fmt.Sprintf("reception rows scan error: %s", err.Error()))
			continue
		}

		reception := model.Reception{
			Id:     id,
			Date:   date,
			PvzId:  pvzId,
			Status: status,
		}

		res = append(res, reception)
	}

	return res
}

func (r *ReceptionRepoPg) AddReception(ctx context.Context, reception model.Reception) (model.Reception, error) {
	rows, err := sq.Insert("receptions").
		Columns("reception_id", "reception_date", "reception_status", "pvz_id").
		Values(reception.Id, reception.Date, reception.Status, reception.PvzId).
		Suffix("RETURNING reception_id, reception_date, reception_status, pvz_id").
		RunWith(r.db).PlaceholderFormat(sq.Dollar).QueryContext(ctx)

	if err != nil {
		slog.Error(fmt.Sprintf("add reception error: %s", err.Error()))
		return model.Reception{}, err
	}

	res := scanRows(rows)
	if len(res) != 1 {
		return model.Reception{}, errors.New("bad result len for adding reception")
	}

	return res[0], nil
}

func (r *ReceptionRepoPg) GetActiveReceptionByPvzId(ctx context.Context, pvzId uuid.UUID) (model.Reception, error) {
	rows, err := sq.Select("reception_id", "reception_date", "reception_status", "pvz_id").
		From("receptions").
		Where(sq.And{
			sq.Eq{"pvz_id": pvzId},
			sq.Eq{"reception_status": consts.ReceptionStatusActive},
		}).RunWith(r.db).PlaceholderFormat(sq.Dollar).QueryContext(ctx)

	if err != nil {
		slog.Error(fmt.Sprintf("get reception error: %s", err.Error()))
		return model.Reception{}, err
	}

	res := scanRows(rows)
	if len(res) != 1 {
		return model.Reception{}, errors.New("bad result len for getting active reception")
	}

	return res[0], nil
}

func (r *ReceptionRepoPg) GetReceptionsByPvzId(ctx context.Context, pvzId uuid.UUID) ([]model.Reception, error) {
	rows, err := sq.Select("reception_id", "reception_date", "reception_status", "pvz_id").
		From("receptions").
		Where(sq.Eq{"pvz_id": pvzId}).
		RunWith(r.db).PlaceholderFormat(sq.Dollar).QueryContext(ctx)

	if err != nil {
		slog.Error(fmt.Sprintf("get receptions error: %s", err.Error()))
		return []model.Reception{}, err
	}

	res := scanRows(rows)

	return res, nil
}

func (r *ReceptionRepoPg) GetReceptionById(ctx context.Context, id uuid.UUID) (model.Reception, error) {
	rows, err := sq.Select("reception_id", "reception_date", "reception_status", "pvz_id").
		From("receptions").
		Where(sq.Eq{"reception_id": id}).
		RunWith(r.db).PlaceholderFormat(sq.Dollar).QueryContext(ctx)

	if err != nil {
		slog.Error(fmt.Sprintf("get reception by id error: %s", err.Error()))
		return model.Reception{}, err
	}

	res := scanRows(rows)
	if len(res) != 1 {
		return model.Reception{}, errors.New("bad result len for getting reception by id")
	}

	return res[0], nil
}

func (r *ReceptionRepoPg) UpdateReceptionStatusById(ctx context.Context, receptionId uuid.UUID, newStatus string) (model.Reception, error) {
	rows, err := sq.Update("receptions").
		Set("reception_status", newStatus).
		Where(sq.Eq{"reception_id": receptionId}).
		Suffix("RETURNING reception_id, reception_date, reception_status, pvz_id").
		RunWith(r.db).PlaceholderFormat(sq.Dollar).QueryContext(ctx)

	if err != nil {
		slog.Error(fmt.Sprintf("update reception error: %s", err.Error()))
		return model.Reception{}, err
	}

	res := scanRows(rows)
	if len(res) != 1 {
		return model.Reception{}, errors.New("bad result len for updating reception")
	}

	return res[0], nil
}
