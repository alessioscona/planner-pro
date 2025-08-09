package db

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"planner-pro/internal/models"
)

type ClientRepo struct{
	pool *pgxpool.Pool
}

func NewClientRepo(pool *pgxpool.Pool) *ClientRepo { return &ClientRepo{pool: pool} }

func (r *ClientRepo) Create(ctx context.Context, c *models.Client) error {
	q := `INSERT INTO clients (id, professional_id, name, email, phone, address, tax_code, vat_number, created_at)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`
	if c.ID == uuid.Nil { c.ID = uuid.New() }
	if c.CreatedAt.IsZero() { c.CreatedAt = time.Now().UTC() }
	_, err := r.pool.Exec(ctx, q, c.ID, c.ProfessionalID, c.Name, c.Email, c.Phone, c.Address, c.TaxCode, c.VatNumber, c.CreatedAt)
	if err != nil { return fmt.Errorf("create client: %w", err) }
	return nil
}

func (r *ClientRepo) ListByProfessional(ctx context.Context, professionalID uuid.UUID) ([]models.Client, error) {
	q := `SELECT id, professional_id, name, email, phone, address, tax_code, vat_number, created_at FROM clients WHERE professional_id=$1 ORDER BY name`
	rows, err := r.pool.Query(ctx, q, professionalID)
	if err != nil { return nil, err }
	defer rows.Close()
	var res []models.Client
	for rows.Next() {
		var c models.Client
		var address map[string]interface{}
		if err := rows.Scan(&c.ID, &c.ProfessionalID, &c.Name, &c.Email, &c.Phone, &address, &c.TaxCode, &c.VatNumber, &c.CreatedAt); err != nil {
			return nil, err
		}
		c.Address = address
		res = append(res, c)
	}
	return res, nil
}
