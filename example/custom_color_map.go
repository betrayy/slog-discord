package main

import (
	"github.com/betrayy/slog-discord"
	"log/slog"
)

func main() {
	webhookURL := "webhook_url"

	// for undefined mappings, the color will be set to 0x0 (black)
	colorMap := slogdiscord.ColorMap{
		slog.LevelInfo: 0xA020F0, // purple
	}

	opts := []slogdiscord.Option{
		slogdiscord.WithSyncMode(true),
		slogdiscord.WithColorMap(colorMap),
	}

	discordHandler, err := slogdiscord.NewDiscordHandler(webhookURL, opts...)
	if err != nil {
		panic(err)
	}

	logger := slog.New(discordHandler)

	logger = logger.With("env", "local")
	logger = logger.WithGroup("request")
	logger.Info("incoming request",
		slog.String("payload", "some payload"),
		slog.String("user_id", "some_user_id"),
	)
}
