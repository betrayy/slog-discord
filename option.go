package slogdiscord

import (
	"log/slog"
	"time"
)

type Option interface {
	apply(*DiscordHandler)
}

// WithHandler sets DiscordHandler.Handler
func WithHandler(handler slog.Handler) Option {
	return optFunc(func(h *DiscordHandler) {
		h.Handler = handler
	})
}

// WithMinLevel sets DiscordHandler.MinLevel.
// This is the minimum slog.Level for Discord logs to be sent
func WithMinLevel(level slog.Level) Option {
	return optFunc(func(h *DiscordHandler) {
		h.MinLevel = level
	})
}

// WithSyncMode sets DiscordHandler.SyncMode.
// If set to true, then logs will be sent to Discord synchronously.
// Otherwise, asynchronously. Note that this means logs may be sent out of order.
func WithSyncMode(sync bool) Option {
	return optFunc(func(h *DiscordHandler) {
		h.SyncMode = sync
	})
}

// WithTimeout sets DiscordHandler.Timeout
func WithTimeout(timeout time.Duration) Option {
	return optFunc(func(h *DiscordHandler) {
		h.Timeout = timeout
	})
}

// WithColorMap sets DiscordHandler.ColorMap
func WithColorMap(colorMap ColorMap) Option {
	return optFunc(func(h *DiscordHandler) {
		h.ColorMap = colorMap
	})
}

// WithEmbedBuilder sets DiscordHandler.BuildEmbed
func WithEmbedBuilder(buildEmbed BuildEmbed) Option {
	return optFunc(func(h *DiscordHandler) {
		h.BuildEmbed = buildEmbed
	})
}

type optFunc func(h *DiscordHandler)

func (opt optFunc) apply(h *DiscordHandler) {
	opt(h)
}
