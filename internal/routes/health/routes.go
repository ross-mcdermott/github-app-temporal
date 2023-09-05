package health

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// HTTP middleware setting a value on the request context
func Routes(r chi.Router) {
	r.Get("/healthz", getHealth)

}

type status struct {
	Ok bool `json:"ok"`
}

func getHealth(w http.ResponseWriter, r *http.Request) {

	data := status{
		Ok: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
