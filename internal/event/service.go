package event

import "log/slog"

const EventTypePageview = "pageview"

type Event struct {
	Domain string `json:"domain"`
	Type   string `json:"type"`
	Url    string `json:"url"`
}

func (e Event) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("domain", e.Domain),
		slog.String("type", e.Type),
		slog.String("url", e.Url),
	)
}
