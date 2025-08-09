package api

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"planner-pro/internal/auth"
	"planner-pro/internal/db"
	"planner-pro/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

var isDevMode bool

func init() {
	isDevMode = os.Getenv("APP_ENV") == "development"
}

type Handler struct {
	pool     *pgxpool.Pool
	clients  *db.ClientRepo
	services *db.ServiceRepo
	apps     *db.AppointmentRepo
}

func NewHandler(pool *pgxpool.Pool) *Handler {
	return &Handler{
		pool:     pool,
		clients:  db.NewClientRepo(pool),
		services: db.NewServiceRepo(pool),
		apps:     db.NewAppointmentRepo(pool),
	}
}

// CreateClient expects JSON { name, email? }
func (h *Handler) CreateClient(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name  string  `json:"name"`
		Email *string `json:"email,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	if in.Name == "" {
		http.Error(w, "name required", http.StatusBadRequest)
		return
	}
	var profID uuid.UUID
	if isDevMode {
		profID = uuid.Nil
	} else {
		pid, ok := auth.FromContextProfessionalID(r.Context())
		if !ok || pid == uuid.Nil {
			http.Error(w, "missing professional_id claim", http.StatusUnauthorized)
			return
		}
		profID = pid
	}
	c := &models.Client{ProfessionalID: profID, Name: in.Name, Email: in.Email}
	if err := h.clients.Create(r.Context(), c); err != nil {
		log.Printf("Errore CreateClient: %v", err)
		http.Error(w, "create error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(c)
}

func (h *Handler) ListClients(w http.ResponseWriter, r *http.Request) {
	var profID uuid.UUID
	if isDevMode {
		profID = uuid.Nil
	} else {
		pid, ok := auth.FromContextProfessionalID(r.Context())
		if !ok || pid == uuid.Nil {
			http.Error(w, "missing professional_id claim", http.StatusUnauthorized)
			return
		}
		profID = pid
	}
	clients, err := h.clients.ListByProfessional(r.Context(), profID)
	if err != nil {
		log.Printf("Errore ListClients: %v", err)
		http.Error(w, "err", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(clients)
}

func (h *Handler) CreateService(w http.ResponseWriter, r *http.Request) {
	var in struct {
		Name     string `json:"name"`
		Duration int    `json:"duration_minutes"`
		Price    int    `json:"price_cents"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	if in.Name == "" || in.Duration <= 0 {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}
	var profID uuid.UUID
	if isDevMode {
		profID = uuid.Nil
	} else {
		pid, ok := auth.FromContextProfessionalID(r.Context())
		if !ok || pid == uuid.Nil {
			http.Error(w, "missing professional_id claim", http.StatusUnauthorized)
			return
		}
		profID = pid
	}
	s := &models.Service{ProfessionalID: profID, Name: in.Name, DurationMinutes: in.Duration, PriceCents: in.Price}
	if err := h.services.Create(r.Context(), s); err != nil {
		log.Printf("Errore CreateService: %v", err)
		http.Error(w, "create error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(s)
}

func (h *Handler) ListServices(w http.ResponseWriter, r *http.Request) {
	var profID uuid.UUID
	if isDevMode {
		profID = uuid.Nil
	} else {
		pid, ok := auth.FromContextProfessionalID(r.Context())
		if !ok || pid == uuid.Nil {
			http.Error(w, "missing professional_id claim", http.StatusUnauthorized)
			return
		}
		profID = pid
	}
	list, err := h.services.ListByProfessional(r.Context(), profID)
	if err != nil {
		log.Printf("Errore ListServices: %v", err)
		http.Error(w, "err", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}

func (h *Handler) CreateAppointment(w http.ResponseWriter, r *http.Request) {
	var in struct {
		ClientID  string    `json:"client_id"`
		ServiceID string    `json:"service_id"`
		StartAt   time.Time `json:"start_at"`
	}
	if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
		http.Error(w, "bad json", http.StatusBadRequest)
		return
	}
	if in.ClientID == "" || in.ServiceID == "" || in.StartAt.IsZero() {
		http.Error(w, "invalid", http.StatusBadRequest)
		return
	}
	// Resolve duration from service (simplified)
	// In real app fetch service and compute end
	start := in.StartAt.UTC()
	end := start.Add(30 * time.Minute)
	var profID uuid.UUID
	if isDevMode {
		profID = uuid.Nil
	} else {
		pid, ok := auth.FromContextProfessionalID(r.Context())
		if !ok || pid == uuid.Nil {
			http.Error(w, "missing professional_id claim", http.StatusUnauthorized)
			return
		}
		profID = pid
	}
	clientID, _ := uuid.Parse(in.ClientID)
	serviceID, _ := uuid.Parse(in.ServiceID)
	a := &models.Appointment{ProfessionalID: profID, ClientID: clientID, ServiceID: serviceID, StartAt: start, EndAt: end, Status: "scheduled"}
	if err := h.apps.CreateWithOverlapCheck(r.Context(), a); err != nil {
		log.Printf("Errore CreateAppointment: %v", err)
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(a)
}

func (h *Handler) ListAppointments(w http.ResponseWriter, r *http.Request) {
	// parse query params from/to
	now := time.Now().UTC()
	from := now.Add(-24 * time.Hour)
	to := now.Add(24 * time.Hour)
	var profID uuid.UUID
	if isDevMode {
		profID = uuid.Nil
	} else {
		pid, ok := auth.FromContextProfessionalID(r.Context())
		if !ok || pid == uuid.Nil {
			http.Error(w, "missing professional_id claim", http.StatusUnauthorized)
			return
		}
		profID = pid
	}
	list, err := h.apps.ListByProfessionalRange(r.Context(), profID, from, to)
	if err != nil {
		log.Printf("Errore ListAppointments: %v", err)
		http.Error(w, "err", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(list)
}
