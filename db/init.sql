CREATE TABLE IF NOT EXISTS urls (
    id UUID PRIMARY KEY,
    original_url TEXT NOT NULL,
    short_url VARCHAR(10) UNIQUE NOT NULL,
    custom_alias VARCHAR(30) UNIQUE,
    created_at TIMESTAMP NOT NULL,
    expires_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS redirect_logs (
    id UUID PRIMARY KEY,
    short_url VARCHAR(10) NOT NULL,
    accessed_at TIMESTAMP NOT NULL,
    referrer TEXT,
    FOREIGN KEY (short_url) REFERENCES urls(short_url)
);
