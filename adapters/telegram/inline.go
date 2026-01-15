package telegram

import (
	"log/slog"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/mtzvd/ironroll/core/roll"
)

// Telegram Inline Behavior
//
// The bot responds only to explicit inline queries.
// Typing only @ironrollbot produces no results.
//
// A roll is performed only when a modifier is explicitly provided.
// To roll with no modifier, the user must specify 0 or +0.
//
// RANDOM INLINE RESULTS (IMPORTANT)
//
// This bot produces non-deterministic (random) inline results.
// To ensure Telegram does NOT reuse previous results, the following
// conditions MUST be met (same as rollrobot):
//
//  1. cache_time = 0        → disables server-side caching
//  2. IsPersonal = true     → marks results as user-specific
//  3. unique result_id      → prevents client-side reuse
//
// Failing to satisfy ALL THREE will cause Telegram clients to
// reinsert the same inline result repeatedly.
func HandleInlineQuery(bot *tgbotapi.BotAPI, query *tgbotapi.InlineQuery) {
	slog.Info(
		"telegram inline query",
		"query_id", query.ID,
		"query", query.Query,
	)

	// Do not respond to empty inline queries.
	if strings.TrimSpace(query.Query) == "" {
		return
	}

	modifier := parseModifier(query.Query)

	// Perform the roll during InlineQuery handling
	// (same model as rollrobot).
	result := roll.Roll(modifier)

	text := formatResult(result)

	// Generate a unique result ID for every response.
	// Time-based uniqueness is sufficient and avoids extra dependencies.
	resultID := strconv.FormatInt(time.Now().UnixNano(), 10)

	article := tgbotapi.NewInlineQueryResultArticle(
		resultID,
		"Ironsworn Roll",
		text,
	)

	cfg := tgbotapi.InlineConfig{
		InlineQueryID: query.ID,
		Results:       []interface{}{article},
		CacheTime:     0,    // disable server-side caching
		IsPersonal:    true, // CRITICAL for random inline bots
	}

	if _, err := bot.Request(cfg); err != nil {
		slog.Error("telegram inline request failed", "err", err)
	}
}

func parseModifier(raw string) int {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return 0
	}

	m, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}

	return m
}
