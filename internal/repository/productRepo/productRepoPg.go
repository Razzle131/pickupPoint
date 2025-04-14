package productRepo

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

type ProductRepoPg struct {
	db *sqlx.DB
}

func NewPg(db *sqlx.DB) *ProductRepoPg {
	return &ProductRepoPg{
		db: db,
	}
}

func scanRows(rows *sql.Rows) []model.Product {
	res := make([]model.Product, 0, consts.SliceMinCap)
	for rows.Next() {
		var id uuid.UUID
		var date time.Time
		var pType string
		var receptionId uuid.UUID

		err := rows.Scan(&id, &date, &pType, &receptionId)
		if err != nil {
			slog.Error(fmt.Sprintf("product rows scan error: %s", err.Error()))
			continue
		}

		product := model.Product{
			Id:          id,
			Date:        date,
			Type:        pType,
			ReceptionId: receptionId,
		}

		res = append(res, product)
	}

	return res
}

func (r *ProductRepoPg) AddProduct(ctx context.Context, product model.Product) (model.Product, error) {
	rows, err := sq.Insert("products").
		Columns("product_id", "product_date", "product_type", "reception_id").
		Values(product.Id, product.Date, product.Type, product.ReceptionId).
		Suffix("RETURNING product_id, product_date, product_type, reception_id").
		RunWith(r.db).PlaceholderFormat(sq.Dollar).QueryContext(ctx)

	if err != nil {
		slog.Error(fmt.Sprintf("add product error: %s", err.Error()))
		return model.Product{}, err
	}

	res := scanRows(rows)
	if len(res) != 1 {
		return model.Product{}, errors.New("bad result len for add product")
	}

	return res[0], nil
}

func (r *ProductRepoPg) DeleteProductById(ctx context.Context, productId uuid.UUID) error {
	_, err := sq.Delete("products").
		Where(sq.Eq{"product_id": productId}).
		Suffix("RETURNING product_id, product_date, product_type, reception_id").
		RunWith(r.db).PlaceholderFormat(sq.Dollar).ExecContext(ctx)

	if err != nil {
		slog.Error(fmt.Sprintf("delete product error: %s", err.Error()))
		return err
	}

	return nil
}

func (r *ProductRepoPg) GetProductsByReceptionId(ctx context.Context, receptionId uuid.UUID) ([]model.Product, error) {
	rows, err := sq.Select("product_id", "product_date", "product_type", "reception_id").
		From("products").Where(sq.Eq{"reception_id": receptionId}).
		RunWith(r.db).PlaceholderFormat(sq.Dollar).QueryContext(ctx)

	if err != nil {
		slog.Error(fmt.Sprintf("get products error: %s", err.Error()))
		return []model.Product{}, err
	}

	return scanRows(rows), nil
}

func (r *ProductRepoPg) GetReceptionLastProduct(ctx context.Context, receptionId uuid.UUID) (model.Product, error) {
	rows, err := sq.Select("product_id", "product_date", "product_type", "reception_id").
		From("products").
		Where(sq.Eq{"reception_id": receptionId}).
		OrderBy("product_date DESC").
		Limit(1).
		RunWith(r.db).PlaceholderFormat(sq.Dollar).QueryContext(ctx)

	if err != nil {
		slog.Error(fmt.Sprintf("get last product error: %s", err.Error()))
		return model.Product{}, err
	}

	res := scanRows(rows)
	if len(res) != 1 {
		slog.Error("bad result len for getting last product")
		return model.Product{}, errors.New("bad result len for getting last product")
	}

	return res[0], nil
}
