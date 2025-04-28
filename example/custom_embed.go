package main

import (
	"github.com/betrayy/slog-discord"
	"github.com/disgoorg/disgo/discord"
	"log/slog"
)

func main() {
	webhookURL := "webhook_url"

	// define your custom embed builder function (ensure it matches the signatures of slogdiscord.BuildEmbed)
	customBuilder := func(record slog.Record, attrs []slog.Attr, builder *discord.EmbedBuilder) discord.Embed {
		builder.SetTitlef("Custom Title - Level %s", record.Level)
		builder.SetDescriptionf("Custom Discription\nMsg: %s", record.Message)
		builder.SetTimestamp(record.Time)
		if record.Level == slog.LevelError {
			builder.SetColor(0xE30B5C) // overrides the color set by the configured color map
		}
		for _, attr := range attrs {
			builder.AddField(attr.Key, attr.Value.String(), true)
		}
		return builder.Build()
	}

	opts := []slogdiscord.Option{
		slogdiscord.WithSyncMode(true),
		slogdiscord.WithEmbedBuilder(customBuilder),
	}

	discordHandler, err := slogdiscord.NewDiscordHandler(webhookURL, opts...)
	if err != nil {
		panic(err)
	}

	logger := slog.New(discordHandler)

	logger = logger.With("env", "local")
	logger = logger.WithGroup("request")
	logger.Error("incoming request",
		slog.String("payload", "some payload"),
		slog.String("user_id", "some_user_id"),
	)
}
