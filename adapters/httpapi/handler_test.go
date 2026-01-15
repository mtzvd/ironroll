package httpapi

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mtzvd/ironroll/core/roll"
)

func TestRollHandler_WithModifier(t *testing.T) {
	// Make roll deterministic
	roll.SetRand(rand.New(rand.NewSource(42)))
	defer roll.ResetRand()

	req := httptest.NewRequest(http.MethodGet, "/roll?m=2", nil)
	rw := httptest.NewRecorder()

	RollHandler(rw, req)

	res := rw.Result()
	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", res.StatusCode)
	}

	var api apiResponse
	if err := json.NewDecoder(res.Body).Decode(&api); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if api.Modifier != 2 {
		t.Fatalf("modifier mismatch: got %d want %d", api.Modifier, 2)
	}
	if api.ActionDie < 1 || api.ActionDie > 6 {
		t.Fatalf("action die out of range: %d", api.ActionDie)
	}
}

func TestRollHandler_DefaultAndInvalidModifier(t *testing.T) {
	// Default modifier (no m param)
	roll.SetRand(rand.New(rand.NewSource(7)))
	defer roll.ResetRand()

	req := httptest.NewRequest(http.MethodGet, "/roll", nil)
	rw := httptest.NewRecorder()
	RollHandler(rw, req)
	if rw.Result().StatusCode != http.StatusOK {
		t.Fatalf("expected 200 OK for default modifier, got %d", rw.Result().StatusCode)
	}

	// Invalid modifier should return 400
	req2 := httptest.NewRequest(http.MethodGet, "/roll?m=bad", nil)
	rw2 := httptest.NewRecorder()
	RollHandler(rw2, req2)
	if rw2.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 Bad Request for invalid modifier, got %d", rw2.Result().StatusCode)
	}
}
