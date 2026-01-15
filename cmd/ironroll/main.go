package main

import (
	cryptorand "crypto/rand"
	"encoding/binary"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/mtzvd/ironroll/adapters/discord"
	"github.com/mtzvd/ironroll/adapters/httpapi"
	"github.com/mtzvd/ironroll/adapters/telegram"
	"github.com/mtzvd/ironroll/core/roll"
	"github.com/mtzvd/ironroll/ratelimit"
	"github.com/mtzvd/ironroll/util/env"
	"github.com/mtzvd/ironroll/util/logging"
)

func main() {
	// ---------------------------------------------------------------------
	// Logging
	// ---------------------------------------------------------------------

	logging.Setup()
	slog.Info("starting ironroll service")

	// ---------------------------------------------------------------------
	// RNG initialization
	//
	// core/roll uses an injectable RNG via SetRand / ResetRand.
	// We use crypto/rand for seeding to ensure true randomness,
	// even in containerized environments where time-based seeds
	// may produce identical values across restarts.
	// ---------------------------------------------------------------------

	var seedBytes [8]byte
	if _, err := cryptorand.Read(seedBytes[:]); err != nil {
		slog.Error("failed to seed RNG", "err", err)
		os.Exit(1)
	}
	seed := int64(binary.LittleEndian.Uint64(seedBytes[:]))
	slog.Info("RNG initialized", "seed", seed)
	roll.SetRand(rand.New(rand.NewSource(seed)))

	// ---------------------------------------------------------------------
	// Environment
	// ---------------------------------------------------------------------

	if err := env.Load(); err != nil {
		slog.Error("failed to load .env", "err", err)
	}

	telegramToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	discordToken := os.Getenv("DISCORD_BOT_TOKEN")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// ---------------------------------------------------------------------
	// HTTP API + Rate Limiting
	// ---------------------------------------------------------------------

	limiter := ratelimit.New(
		30,            // requests
		time.Minute,   // per window
		5*time.Minute, // temporary block
	)

	httpHandler := httpapi.RateLimitMiddleware(
		limiter,
		http.HandlerFunc(httpapi.RollHandler),
	)

	http.Handle("/roll", httpHandler)

	go func() {
		slog.Info("http api started", "port", port)
		if err := http.ListenAndServe("127.0.0.1:"+port, nil); err != nil {
			slog.Error("http server failed", "err", err)
			os.Exit(1)
		}
	}()

	// ---------------------------------------------------------------------
	// Telegram Inline Bot
	// ---------------------------------------------------------------------

	slog.Info("telegram token check", "present", telegramToken != "")

	if telegramToken != "" {
		slog.Info("initializing telegram bot")

		bot, err := tgbotapi.NewBotAPI(telegramToken)
		if err != nil {
			slog.Error("failed to create telegram bot", "err", err)
			os.Exit(1)
		}

		slog.Info(
			"telegram bot connected",
			"username", bot.Self.UserName,
		)

		u := tgbotapi.NewUpdate(0)
		u.Timeout = 60
		updates := bot.GetUpdatesChan(u)

		go func() {
			slog.Info("telegram bot started")
			for update := range updates {
				if update.InlineQuery != nil {
					telegram.HandleInlineQuery(bot, update.InlineQuery)
				}
			}
		}()
	} else {
		slog.Warn("telegram bot disabled (no TELEGRAM_BOT_TOKEN)")
	}

	// ---------------------------------------------------------------------
	// Discord Bot
	// ---------------------------------------------------------------------

	if discordToken != "" {
		dg, err := discordgo.New("Bot " + discordToken)
		if err != nil {
			slog.Error("failed to create discord session", "err", err)
			os.Exit(1)
		}

		dg.AddHandler(discord.HandleInteraction)

		if err := dg.Open(); err != nil {
			slog.Error("failed to open discord connection", "err", err)
			os.Exit(1)
		}

		_, err = dg.ApplicationCommandCreate(
			dg.State.User.ID,
			"", // global command
			discord.Command,
		)
		if err != nil {
			slog.Error("failed to register discord command", "err", err)
			os.Exit(1)
		}

		slog.Info("discord bot started")
	} else {
		slog.Warn("discord bot disabled (no DISCORD_BOT_TOKEN)")
	}

	// ---------------------------------------------------------------------
	// Block forever
	// ---------------------------------------------------------------------

	select {}
}
