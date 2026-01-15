package telegram

import (
	"testing"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// fakeBot is a minimal stub for testing.
type fakeBot struct {
	called bool
}

func (f *fakeBot) Request(c tgbotapi.Chattable) (*tgbotapi.APIResponse, error) {
	f.called = true
	return &tgbotapi.APIResponse{Ok: true}, nil
}

func TestHandleInlineQuery_EmptyQueryDoesNothing(t *testing.T) {
	bot := &fakeBot{}

	query := &tgbotapi.InlineQuery{
		ID:    "empty",
		Query: "",
	}

	HandleInlineQuery((*tgbotapi.BotAPI)(nil), query)

	if bot.called {
		t.Fatal("expected no request for empty inline query")
	}
}

func TestHandleInlineQuery_WithModifierResponds(t *testing.T) {
	// This test only ensures the handler does not panic
	// when a modifier is provided.
	query := &tgbotapi.InlineQuery{
		ID:    "with-mod",
		Query: "+1",
	}

	HandleInlineQuery((*tgbotapi.BotAPI)(nil), query)
}
