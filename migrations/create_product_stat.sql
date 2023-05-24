-- Clickhouse БД для сбора статистики

CREATE TABLE products_from_orders
(
    count        UInt32,
    product_uuid UUID,
    price        Float64,
    product_name String,
    product_type String,
    created_at   DateTime
)
    ENGINE = MergeTree
        PARTITION BY toYYYYMM(created_at)
        ORDER BY (created_at, price);

CREATE TABLE products_from_orders_queue
(
    count        UInt32,
    product_uuid UUID,
    price        Float64,
    product_name String,
    product_type String,
    created_at   DateTime
)
    ENGINE = Kafka
        SETTINGS kafka_broker_list = '192.168.0.109:9092',
            kafka_topic_list = 'prodStat',
            kafka_group_name = 'prodStatGroup1',
            kafka_format = 'JSONEachRow',
            kafka_max_block_size = 1048576;

CREATE MATERIALIZED VIEW products_from_orders_queue_mv TO products_from_orders AS
SELECT *
FROM products_from_orders_queue;