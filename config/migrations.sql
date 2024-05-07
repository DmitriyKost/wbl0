-- To run multiple queries, it's better to separate them via 'query_separator', so the code may execute them one by one.
CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    order_data JSONB NOT NULL
);
