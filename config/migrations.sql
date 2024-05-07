-- To run multiple queries, it's best to separate them using a "query separator" (separate the words in quotes with an underscore).
-- So the code can run them one by one.
CREATE TABLE IF NOT EXISTS orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    order_data JSONB NOT NULL
);
