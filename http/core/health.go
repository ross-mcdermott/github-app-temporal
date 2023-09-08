package health

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func NewHealthHandler(logger *slog.Logger) *Health {

	// Create new instance of the hooks
	handler := &Health{
		logger: logger,
	}

	return handler
}

type Health struct {
	logger *slog.Logger
}

func (a *Health) Register(router chi.Router, route string) {

	a.logger.Debug(fmt.Sprintf("Registering '%s' route.", route))

	router.Group(func(r chi.Router) {
		r.Get(route, a.process)
	})

}

func (a *Health) process(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	r = r.WithContext(ctx)

	a.logger.Debug("Health called")

	data := status{
		Ok: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)

}

// func RegisterRoutes(router chi.Router, logger *slog.Logger) {

// 	logger.Debug("Registering '/health' route.")

// 	router.Group(func(r chi.Router) {

// 		r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {

// 			logger.Debug("Health called")

// 			data := status{
// 				Ok: true,
// 			}

// 			w.Header().Set("Content-Type", "application/json")
// 			json.NewEncoder(w).Encode(data)

// 		})
// 	})

// }

type status struct {
	Ok bool `json:"ok"`
}
