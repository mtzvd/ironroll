package httpapi

import "github.com/mtzvd/ironroll/core/roll"

// apiResponse defines the public JSON shape returned by the HTTP API.
//
// This is intentionally explicit and mirrors the roll.Result structure,
// without exposing internal types.
type apiResponse struct {
	ActionDie     int    `json:"action_die"`
	Modifier      int    `json:"modifier"`
	ChallengeDice [2]int `json:"challenge_dice"`
	Total         int    `json:"total"`
	Outcome       string `json:"outcome"`
}

func formatResult(r roll.Result) apiResponse {
	return apiResponse{
		ActionDie:     r.ActionDie,
		Modifier:      r.Modifier,
		ChallengeDice: r.ChallengeDice,
		Total:         r.Total,
		Outcome:       string(r.Outcome),
	}
}
