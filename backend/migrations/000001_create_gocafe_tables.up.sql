-- gocafe_users
CREATE TABLE IF NOT EXISTS gocafe_users (
    id            SERIAL PRIMARY KEY,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT now(),
    email         VARCHAR(255) NOT NULL UNIQUE,
    name          VARCHAR(255) DEFAULT '',
    password_hash TEXT DEFAULT ''
);

CREATE INDEX IF NOT EXISTS idx_gocafe_users_email ON gocafe_users (email);

-- gocafe_cafe_listings
CREATE TABLE IF NOT EXISTS gocafe_cafe_listings (
    id           SERIAL PRIMARY KEY,
    created_at   TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at   TIMESTAMP WITH TIME ZONE DEFAULT now(),
    user_id      BIGINT NOT NULL,
    name         VARCHAR(255) NOT NULL,
    address      TEXT DEFAULT '',
    description  TEXT DEFAULT '',
    CONSTRAINT fk_gocafe_cafe_listings_user FOREIGN KEY (user_id) REFERENCES gocafe_users (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_gocafe_cafe_listings_user_id ON gocafe_cafe_listings (user_id);

-- gocafe_ratings
CREATE TABLE IF NOT EXISTS gocafe_ratings (
    id              SERIAL PRIMARY KEY,
    created_at      TIMESTAMP WITH TIME ZONE DEFAULT now(),
    updated_at      TIMESTAMP WITH TIME ZONE DEFAULT now(),
    user_id         BIGINT NOT NULL,
    cafe_listing_id BIGINT NOT NULL,
    visited_at      TIMESTAMP WITH TIME ZONE NOT NULL,
    rating          INTEGER NOT NULL,
    review          TEXT DEFAULT '',
    CONSTRAINT fk_gocafe_ratings_user FOREIGN KEY (user_id) REFERENCES gocafe_users (id) ON DELETE CASCADE,
    CONSTRAINT fk_gocafe_ratings_cafe_listing FOREIGN KEY (cafe_listing_id) REFERENCES gocafe_cafe_listings (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_gocafe_ratings_user_id ON gocafe_ratings (user_id);
CREATE INDEX IF NOT EXISTS idx_gocafe_ratings_cafe_listing_id ON gocafe_ratings (cafe_listing_id);
