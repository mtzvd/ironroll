package httpapi

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/mtzvd/ironroll/core/roll"
)

// RollHandler handles GET /roll requests.
//
// Query parameters:
//   - m: optional integer modifier (defaults to 0)
//
// Responses:
//   - 200 OK with JSON roll result
//   - 400 Bad Request if modifier is invalid
func RollHandler(w http.ResponseWriter, r *http.Request) {
	modifier := 0

	if raw := r.URL.Query().Get("m"); raw != "" {
		m, err := strconv.Atoi(raw)
		if err != nil {
			http.Error(w, "invalid modifier", http.StatusBadRequest)
			return
		}
		modifier = m
	}

	result := roll.Roll(modifier)

	response := formatResult(result)

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(response)
}
