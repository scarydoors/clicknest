-- +goose Up
CREATE TABLE events
(
    timestamp DateTime CODEC(Delta(4), ZSTD),
    domain LowCardinality(String),
    type LowCardinality(String),
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(timestamp)
ORDER BY (domain, timestamp, type);

-- +goose Down
DROP TABLE events;
