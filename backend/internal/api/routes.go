package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"planner-pro/internal/auth"
)

func RegisterRoutes(r *chi.Mux, pool *pgxpool.Pool, oidcMiddleware auth.Middleware) {
	r.Route("/api/v1", func(r chi.Router) {
		r.With(oidcMiddleware).Post("/clients", func(w http.ResponseWriter, r *http.Request){ NewHandler(pool).CreateClient(w,r) })
		r.With(oidcMiddleware).Get("/clients", func(w http.ResponseWriter, r *http.Request){ NewHandler(pool).ListClients(w,r) })
		r.With(oidcMiddleware).Post("/services", func(w http.ResponseWriter, r *http.Request){ NewHandler(pool).CreateService(w,r) })
		r.With(oidcMiddleware).Get("/services", func(w http.ResponseWriter, r *http.Request){ NewHandler(pool).ListServices(w,r) })
		r.With(oidcMiddleware).Post("/appointments", func(w http.ResponseWriter, r *http.Request){ NewHandler(pool).CreateAppointment(w,r) })
		r.With(oidcMiddleware).Get("/appointments", func(w http.ResponseWriter, r *http.Request){ NewHandler(pool).ListAppointments(w,r) })
	})
}
