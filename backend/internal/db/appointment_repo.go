package db

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"planner-pro/internal/models"
)

type AppointmentRepo struct{ pool *pgxpool.Pool }

func NewAppointmentRepo(p *pgxpool.Pool) *AppointmentRepo { return &AppointmentRepo{pool: p} }

// CreateWithOverlapCheck: transactionally check overlaps and insert
func (r *AppointmentRepo) CreateWithOverlapCheck(ctx context.Context, ap *models.Appointment) error {
	if ap.ID == uuid.Nil { ap.ID = uuid.New() }
	if ap.CreatedAt.IsZero() { ap.CreatedAt = time.Now().UTC() }

	tx, err := r.pool.Begin(ctx)
	if err != nil { return err }
	defer func(){ _ = tx.Rollback(ctx) }()

	// check overlaps for same professional
	qCheck := `SELECT 1 FROM appointments WHERE professional_id=$1 AND NOT (end_at <= $2 OR start_at >= $3) LIMIT 1 FOR UPDATE`
	var dummy int
	row := tx.QueryRow(ctx, qCheck, ap.ProfessionalID, ap.StartAt, ap.EndAt)
	if err := row.Scan(&dummy); err == nil {
		return errors.New("time slot unavailable")
	}

	qIns := `INSERT INTO appointments (id, professional_id, client_id, service_id, location_id, start_at, end_at, status, notes, created_at) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`
	_, err = tx.Exec(ctx, qIns, ap.ID, ap.ProfessionalID, ap.ClientID, ap.ServiceID, ap.LocationID, ap.StartAt, ap.EndAt, ap.Status, ap.Notes, ap.CreatedAt)
	if err != nil { return fmt.Errorf("insert ap: %w", err) }

	if err := tx.Commit(ctx); err != nil { return err }
	return nil
}

func (r *AppointmentRepo) ListByProfessionalRange(ctx context.Context, professionalID uuid.UUID, from, to time.Time) ([]models.Appointment, error) {
	q := `SELECT id, professional_id, client_id, service_id, location_id, start_at, end_at, status, notes, created_at FROM appointments WHERE professional_id=$1 AND start_at >= $2 AND end_at <= $3 ORDER BY start_at`
	rows, err := r.pool.Query(ctx, q, professionalID, from, to)
	if err != nil { return nil, err }
	defer rows.Close()
	var out []models.Appointment
	for rows.Next() {
		var a models.Appointment
		if err := rows.Scan(&a.ID, &a.ProfessionalID, &a.ClientID, &a.ServiceID, &a.LocationID, &a.StartAt, &a.EndAt, &a.Status, &a.Notes, &a.CreatedAt); err != nil { return nil, err }
		out = append(out, a)
	}
	return out, nil
}
