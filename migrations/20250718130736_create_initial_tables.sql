-- +goose Up
CREATE TABLE events
(
    timestamp DateTime CODEC(Delta(4), ZSTD),
    domain LowCardinality(String),
    kind LowCardinality(String),
    session_id UInt64,
    user_id UInt64,
    pathname String,
)
ENGINE = MergeTree
PARTITION BY toYYYYMM(timestamp)
PRIMARY KEY (domain, toDate(timestamp), kind, user_id)
ORDER BY (domain, toDate(timestamp), kind, user_id, timestamp);

CREATE TABLE sessions
(
    start DateTime CODEC(Delta(4), ZSTD),
    end DateTime CODEC(Delta(4), ZSTD),
    domain LowCardinality(String),
    duration UInt32,
    session_id UInt64,
    user_id UInt64,
    sign Int8,
    INDEX minmax_end end TYPE minmax GRANULARITY 1,
)
ENGINE = VersionedCollapsingMergeTree(sign, end)
PARTITION BY toYYYYMM(start)
PRIMARY KEY (domain, toDate(start), user_id, session_id)
ORDER BY (domain, toDate(start), user_id, session_id, start);

-- +goose Down
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS sessions;
