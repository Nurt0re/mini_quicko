
CREATE TABLE IF NOT EXISTS price_history (
    id SERIAL PRIMARY KEY,
    product_id VARCHAR(255) NOT NULL,
    product_name TEXT NOT NULL,
    timestamp TIMESTAMP NOT NULL DEFAULT NOW(),
    min_price DECIMAL(10, 2) NOT NULL,
    max_price DECIMAL(10, 2) NOT NULL,
    avg_price DECIMAL(10, 2) NOT NULL,
    offer_count INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS offers_cache (
    id SERIAL PRIMARY KEY,
    product_id VARCHAR(255) NOT NULL,
    city_id VARCHAR(255) NOT NULL,
    merchant_id VARCHAR(255) NOT NULL,
    merchant_name TEXT NOT NULL,
    merchant_rating DECIMAL(3, 2),
    merchant_reviews_count INTEGER,
    title TEXT NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    total_offers INTEGER NOT NULL,
    offers_count INTEGER NOT NULL,
    fetched_at TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_price_history_product_id ON price_history(product_id);
CREATE INDEX IF NOT EXISTS idx_price_history_timestamp ON price_history(timestamp DESC);
CREATE INDEX IF NOT EXISTS idx_price_history_product_timestamp ON price_history(product_id, timestamp DESC);

CREATE INDEX IF NOT EXISTS idx_offers_cache_product_city ON offers_cache(product_id, city_id);
CREATE INDEX IF NOT EXISTS idx_offers_cache_fetched_at ON offers_cache(fetched_at DESC);
