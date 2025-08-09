package db

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"planner-pro/internal/models"
)

type ServiceRepo struct{ pool *pgxpool.Pool }

func NewServiceRepo(p *pgxpool.Pool) *ServiceRepo { return &ServiceRepo{pool: p} }

func (r *ServiceRepo) Create(ctx context.Context, s *models.Service) error {
	if s.ID == uuid.Nil { s.ID = uuid.New() }
	if s.CreatedAt.IsZero() { s.CreatedAt = time.Now().UTC() }
	q := `INSERT INTO services (id, professional_id, code, name, duration_minutes, price_cents, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7)`
	_, err := r.pool.Exec(ctx, q, s.ID, s.ProfessionalID, s.Code, s.Name, s.DurationMinutes, s.PriceCents, s.CreatedAt)
	return err
}

func (r *ServiceRepo) ListByProfessional(ctx context.Context, professionalID uuid.UUID) ([]models.Service, error) {
	q := `SELECT id, professional_id, code, name, duration_minutes, price_cents, created_at FROM services WHERE professional_id=$1 ORDER BY name`
	rows, err := r.pool.Query(ctx, q, professionalID)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []models.Service
	for rows.Next() {
		var s models.Service
		if err := rows.Scan(&s.ID, &s.ProfessionalID, &s.Code, &s.Name, &s.DurationMinutes, &s.PriceCents, &s.CreatedAt); err != nil { return nil, err }
		out = append(out, s)
	}
	return out, nil
}
