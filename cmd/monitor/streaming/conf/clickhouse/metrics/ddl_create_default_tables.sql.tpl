CREATE TABLE IF NOT EXISTS <database>.metrics ON CLUSTER '{cluster}'
(
    `org_name`            LowCardinality(String),
    `tenant_id`           LowCardinality(String),
    `metric_group`        LowCardinality(String),
    `timestamp`           DateTime64(9,'Asia/Shanghai') CODEC (DoubleDelta),
    `number_field_keys`   Array(LowCardinality(String)),
    `number_field_values` Array(Float64),
    `string_field_keys`   Array(LowCardinality(String)),
    `string_field_values` Array(String),
    `tag_keys`            Array(LowCardinality(String)),
    `tag_values`          Array(LowCardinality(String))
)
ENGINE = ReplicatedMergeTree('/clickhouse/tables/{cluster}-{shard}/{database}/metrics', '{replica}')
PARTITION BY toYYYYMMDD(timestamp)
ORDER BY (org_name, tenant_id, metric_group, timestamp)
TTL toDateTime(timestamp) + INTERVAL <ttl_in_days> DAY;

// create distributed table
// notice: ddls to the <table> table should be synced to the <table>_all table
CREATE TABLE IF NOT EXISTS <database>.metrics_all ON CLUSTER '{cluster}' AS <database>.metrics
ENGINE = Distributed('{cluster}', <database>, metrics, rand());

CREATE TABLE IF NOT EXISTS <database>.metrics_meta ON CLUSTER '{cluster}'
(
    `org_name`            LowCardinality(String),
    `tenant_id`           LowCardinality(String),
    `metric_group`        LowCardinality(String),
    `timestamp`           DateTime64(9,'Asia/Shanghai') CODEC (DoubleDelta),
    `number_field_keys`   Array(LowCardinality(String)),
    `string_field_keys`   Array(LowCardinality(String)),
    `tag_keys`            Array(LowCardinality(String)),

    INDEX idx_timestamp(timestamp) TYPE minmax GRANULARITY 2
)
ENGINE = ReplicatedReplacingMergeTree('/clickhouse/tables/{cluster}-{shard}/{database}/metrics_meta', '{replica}')
ORDER BY (org_name, tenant_id, metric_group, number_field_keys, string_field_keys, tag_keys)
TTL toDateTime(timestamp) + INTERVAL <ttl_in_days> DAY;

// create distributed table
// notice: ddls to the <table> table should be synced to the <table>_all table
CREATE TABLE IF NOT EXISTS <database>.metrics_meta_all ON CLUSTER '{cluster}' AS <database>.metrics_meta
ENGINE = Distributed('{cluster}', <database>, metrics_meta, rand());
