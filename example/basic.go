package main

import (
	"github.com/betrayy/slog-discord"
	"log/slog"
	"os"
)

func main() {
	webhookURL := "webhook_url"

	opts := []slogdiscord.Option{
		slogdiscord.WithMinLevel(slog.LevelWarn), // only log to Discord if slog level is warning or higher
		slogdiscord.WithSyncMode(true),           // send logs to discord synchronously
		// also log to console
		slogdiscord.WithHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn})),
	}

	// instantiate the discord handler
	// typically, an error is returned if the webhook url is invalid
	discordHandler, err := slogdiscord.NewDiscordHandler(webhookURL, opts...)
	if err != nil {
		panic(err)
	}

	// instantiate the slog logger
	logger := slog.New(discordHandler)

	// start logging!
	logger = logger.With("env", "local")
	logger.Error("an error occurred", "error", "error msg") // logged to discord and console
}
