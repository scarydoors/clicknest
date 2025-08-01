package analytics

import (
	"fmt"
	"log/slog"
	"net/url"
	"time"
)

const EventKindPageview = "pageview"

type Event struct {
	Timestamp time.Time
	Domain string
	Kind string
	SessionId uint64
	UserId uint64
	Pathname string
}

func NewEvent(timestamp time.Time, domain, kind, rawUrl string) (Event, error) {
	parsed, err := parseURL(rawUrl)
	if err != nil {
		return Event{}, err
	}

	return Event{
		Timestamp: timestamp,
		Domain: domain,
		Kind: kind,
		Pathname: parsed.pathname,
	}, nil
}

func (e Event) LogValue() slog.Value {
	return slog.GroupValue(
		slog.Time("timestamp", e.Timestamp),
		slog.String("domain", e.Domain),
		slog.String("kind", e.Kind),
		slog.Uint64("session_id", e.SessionId),
		slog.Uint64("user_id", e.UserId),
		slog.String("pathname", e.Pathname),
	)
}

type parsedURL struct {
	pathname string
}

func parseURL(rawUrl string) (parsedURL, error) {
	u, err := url.Parse(rawUrl)
	if err != nil {
		return parsedURL{}, fmt.Errorf("parse url: %w", err)
	}

	// truncate trailing '/' to avoid counting /path and /path/ as seperate paths
	pathname := u.Path
	if last := len(pathname) - 1; last > 0 && pathname[last] == '/' {
		pathname = pathname[:last]
	}

	return parsedURL{
		pathname: pathname,
	}, nil
}
