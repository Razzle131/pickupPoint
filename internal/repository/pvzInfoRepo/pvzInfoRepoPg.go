package pvzInfoRepo

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"

	"github.com/Razzle131/pickupPoint/internal/dto"
	"github.com/Razzle131/pickupPoint/internal/model"
	"github.com/jmoiron/sqlx"
)

type PvzInfoRepoPg struct {
	db *sqlx.DB
}

func NewPg(db *sqlx.DB) *PvzInfoRepoPg {
	return &PvzInfoRepoPg{
		db: db,
	}
}

func scanRows(rows *sql.Rows, limit int) map[model.Pvz]map[model.Reception][]model.Product {
	res := make(map[model.Pvz](map[model.Reception][]model.Product), limit)
	for rows.Next() {
		var pvzId uuid.UUID
		var pvzCity string
		var pvzDate time.Time

		var receptionId uuid.UUID
		var receptionDate time.Time
		var receptionStatus string

		var productId uuid.NullUUID
		var productDate sql.NullTime
		var productType sql.NullString

		err := rows.Scan(&pvzId, &pvzCity, &pvzDate, &receptionId, &receptionDate, &receptionStatus, &productId, &productDate, &productType)
		if err != nil {
			slog.Error(fmt.Sprintf("pvz list rows scan error: %s", err.Error()))
			continue
		}

		pvz := model.Pvz{
			Id:      pvzId,
			City:    pvzCity,
			RegDate: pvzDate,
		}

		reception := model.Reception{
			Id:     receptionId,
			Date:   receptionDate,
			PvzId:  pvzId,
			Status: receptionStatus,
		}

		product := model.Product{
			Id:          productId.UUID,
			Date:        productDate.Time,
			Type:        productType.String,
			ReceptionId: receptionId,
		}

		if _, found := res[pvz]; !found {
			res[pvz] = make(map[model.Reception][]model.Product, limit)
		}

		if _, found := res[pvz][reception]; !found {
			res[pvz][reception] = make([]model.Product, 0, limit)
		}

		if productId.Valid && productDate.Valid && productType.Valid {
			res[pvz][reception] = append(res[pvz][reception], product)
		}
	}

	return res
}

func (r *PvzInfoRepoPg) ListPvzInfo(ctx context.Context, params dto.PvzInfoFilterDto) (map[model.Pvz]map[model.Reception][]model.Product, error) {
	builder := sq.Select(
		"pz.pvz_id",
		"pz.pvz_city",
		"pz.pvz_date",
		"r.reception_id",
		"r.reception_date",
		"r.reception_status",
		"p.product_id",
		"p.product_date",
		"p.product_type").
		From("pvz pz").
		Join("receptions r USING (pvz_id)")

	if params.StartDateGiven {
		builder = builder.Where(sq.GtOrEq{"r.reception_date": params.StartDate})
	}

	if params.EndDateGiven {
		builder = builder.Where(sq.LtOrEq{"r.reception_date": params.EndDate})
	}

	offset := (params.Page - 1) * params.Limit

	rows, err := builder.
		LeftJoin("products p USING (reception_id)").
		OrderBy(
			"pz.pvz_city",
			"pz.pvz_date",
			"r.reception_date",
			"r.reception_status",
			"p.product_date",
			"p.product_type").
		Limit(uint64(params.Limit)).Offset(uint64(offset)).
		RunWith(r.db).PlaceholderFormat(sq.Dollar).QueryContext(ctx)

	if err != nil {
		slog.Error(fmt.Sprintf("get pvz list error: %s", err.Error()))
		return map[model.Pvz]map[model.Reception][]model.Product{}, err
	}

	return scanRows(rows, params.Limit), nil
}
