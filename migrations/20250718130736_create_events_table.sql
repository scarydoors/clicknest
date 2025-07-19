-- +goose Up
CREATE TABLE events
(
    timestamp DateTime CODEC(Delta(4), ZSTD),
    domain LowCardinality(String),
    kind LowCardinality(String),
    pathname String,
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(timestamp)
ORDER BY (domain, timestamp, kind);

-- +goose Down
DROP TABLE events;
