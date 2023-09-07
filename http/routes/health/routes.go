package health

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(router chi.Router, logger *slog.Logger) {

	logger.Debug("Registering '/health' route.")

	router.Group(func(r chi.Router) {

		r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {

			logger.Debug("Health called")

			data := status{
				Ok: true,
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(data)

		})
	})

}

type status struct {
	Ok bool `json:"ok"`
}
