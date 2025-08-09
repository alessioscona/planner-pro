package models

import (
	"time"

	"github.com/google/uuid"
)

type Professional struct {
	ID uuid.UUID `json:"id"`
	UserID string `json:"user_id"`
	Name string `json:"name"`
	Email string `json:"email"`
	DefaultLocationID *uuid.UUID `json:"default_location_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Client struct {
	ID uuid.UUID `json:"id"`
	ProfessionalID uuid.UUID `json:"professional_id"`
	Name string `json:"name"`
	Email *string `json:"email,omitempty"`
	Phone *string `json:"phone,omitempty"`
	Address map[string]interface{} `json:"address,omitempty"`
	TaxCode *string `json:"tax_code,omitempty"`
	VatNumber *string `json:"vat_number,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Service struct {
	ID uuid.UUID `json:"id"`
	ProfessionalID uuid.UUID `json:"professional_id"`
	Code *string `json:"code,omitempty"`
	Name string `json:"name"`
	DurationMinutes int `json:"duration_minutes"`
	PriceCents int `json:"price_cents"`
	CreatedAt time.Time `json:"created_at"`
}

type Appointment struct {
	ID uuid.UUID `json:"id"`
	ProfessionalID uuid.UUID `json:"professional_id"`
	ClientID uuid.UUID `json:"client_id"`
	ServiceID uuid.UUID `json:"service_id"`
	LocationID *uuid.UUID `json:"location_id,omitempty"`
	StartAt time.Time `json:"start_at"`
	EndAt time.Time `json:"end_at"`
	Status string `json:"status"`
	Notes *string `json:"notes,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
