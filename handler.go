package slogdiscord

import (
	"context"
	"errors"
	"fmt"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/webhook"
	slogcommon "github.com/samber/slog-common"
	"log/slog"
	"net/http"
	"time"
)

type ColorMap map[slog.Level]int

type BuildEmbed func(slog.Record, []slog.Attr, *discord.EmbedBuilder) discord.Embed

var DefaultColorMappings = ColorMap{
	slog.LevelDebug: 0x63C5DA,
	slog.LevelInfo:  0x63C5DA,
	slog.LevelWarn:  0xFFA500,
	slog.LevelError: 0xFF0000,
}

type DiscordHandler struct {
	Handler    slog.Handler   // Underlying slog handler that gets called in every slog.Handler method (default: noop handler that does nothing)
	WebhookURL string         // Discord webhook URL (required)
	MinLevel   slog.Level     // minimum slog level for Discord logs (default: slog.LevelDebug)
	SyncMode   bool           // If true, then send logs to Discord synchronously. Otherwise, asynchronously (default: false)
	Timeout    time.Duration  // API request timeout (default: 10s)
	ColorMap   ColorMap       // The color mappings for slog.Level for Discord embeds (default: DefaultColorMappings)
	BuildEmbed BuildEmbed     // The function to build the discord.Embed before logging to Discord (default: DiscordHandler.defaultBuildEmbed)
	Client     webhook.Client // Disgo webhook client (required). Automatically set if you call NewDiscordHandler
	attrs      []slog.Attr    // Collects slog attribute fields
	groups     []string       // Collects slog attribute groups
}

var _ slog.Handler = (*DiscordHandler)(nil)

func NewDiscordHandler(webhookURL string, opts ...Option) (slog.Handler, error) {
	h := defaultHandler()
	h.WebhookURL = webhookURL

	for _, opt := range opts {
		opt.apply(h)
	}

	if h.WebhookURL == "" {
		return nil, errors.New("missing webhook URL")
	}

	if err := h.initWebhookClient(); err != nil {
		return nil, err
	}

	return h, nil
}

func defaultHandler() *DiscordHandler {
	h := &DiscordHandler{}
	h.MinLevel = slog.LevelDebug
	h.SyncMode = false
	h.Timeout = 10 * time.Second
	h.ColorMap = DefaultColorMappings
	h.BuildEmbed = h.defaultBuildEmbed
	return h
}

func (h *DiscordHandler) initWebhookClient() error {
	var restOpts []rest.ConfigOpt
	restOpts = append(restOpts, rest.WithHTTPClient(&http.Client{Timeout: h.Timeout}))

	var opts []webhook.ConfigOpt
	opts = append(opts, webhook.WithRestClientConfigOpts(restOpts...))
	if h.Handler != nil {
		opts = append(opts, webhook.WithLogger(slog.New(h.Handler)))
	}

	client, err := webhook.NewWithURL(h.WebhookURL, opts...)
	if err != nil {
		return fmt.Errorf("could not init webhook client: %w", err)
	}

	h.Client = client
	return nil
}

func (h *DiscordHandler) Enabled(ctx context.Context, level slog.Level) bool {
	handlerEnabled := false
	if h.Handler != nil {
		handlerEnabled = h.Handler.Enabled(ctx, level)
	}
	return level >= h.MinLevel || handlerEnabled
}

func (h *DiscordHandler) Handle(ctx context.Context, record slog.Record) error {
	if record.Level >= h.MinLevel {
		if h.SyncMode {
			h.logDiscord(record)
		} else {
			go h.logDiscord(record)
		}
	}
	if h.Handler == nil {
		return nil
	}
	return h.Handler.Handle(ctx, record)
}

func (h *DiscordHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	nh := *h
	nh.attrs = slogcommon.AppendAttrsToGroup(h.groups, h.attrs, attrs...)
	if h.Handler != nil {
		nh.Handler = h.Handler.WithAttrs(attrs)
	}
	return &nh
}

func (h *DiscordHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	nh := *h
	nh.groups = append(nh.groups, name)
	if h.Handler != nil {
		nh.Handler = h.Handler.WithGroup(name)
	}
	return &nh
}

func (h *DiscordHandler) logDiscord(record slog.Record) {
	attrs := slogcommon.AppendRecordAttrsToAttrs(h.attrs, h.groups, &record)
	attrs = slogcommon.RemoveEmptyAttrs(attrs)

	eb := discord.NewEmbedBuilder()
	eb.SetColor(h.ColorMap[record.Level])
	embed := h.BuildEmbed(record, attrs, eb)

	h.Client.CreateEmbeds([]discord.Embed{embed})
}

func (h *DiscordHandler) defaultBuildEmbed(record slog.Record, attrs []slog.Attr, eb *discord.EmbedBuilder) discord.Embed {
	eb.SetTitle(record.Level.String()).
		SetDescription(record.Message).
		SetTimestamp(record.Time)
	h.populateEmbedFields("", attrs, eb)
	return eb.Build()
}

func (h *DiscordHandler) populateEmbedFields(base string, attrs []slog.Attr, eb *discord.EmbedBuilder) {
	for _, attr := range attrs {
		key := attr.Key
		val := attr.Value
		kind := attr.Value.Kind()

		if kind == slog.KindGroup {
			h.populateEmbedFields(base+key+".", val.Group(), eb)
		} else {
			eb.AddField(base+key, slogcommon.ValueToString(val), false)
		}
	}
}
