CREATE TABLE IF NOT EXISTS products (
    product_id        SERIAL PRIMARY KEY,
    name              TEXT NOT NULL,
    category          TEXT NOT NULL,
    description       TEXT,
    price             DECIMAL(10, 2) NOT NULL,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
