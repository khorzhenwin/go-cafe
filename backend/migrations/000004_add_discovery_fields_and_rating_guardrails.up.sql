ALTER TABLE gocafe_cafe_listings
ADD COLUMN IF NOT EXISTS city VARCHAR(255) DEFAULT '',
ADD COLUMN IF NOT EXISTS neighborhood VARCHAR(255) DEFAULT '',
ADD COLUMN IF NOT EXISTS latitude DOUBLE PRECISION,
ADD COLUMN IF NOT EXISTS longitude DOUBLE PRECISION,
ADD COLUMN IF NOT EXISTS image_url TEXT DEFAULT '',
ADD COLUMN IF NOT EXISTS source_cafe_id BIGINT;

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'fk_gocafe_cafe_listings_source_cafe'
    ) THEN
        ALTER TABLE gocafe_cafe_listings
        ADD CONSTRAINT fk_gocafe_cafe_listings_source_cafe
        FOREIGN KEY (source_cafe_id) REFERENCES gocafe_cafe_listings (id) ON DELETE SET NULL;
    END IF;
END $$;

CREATE INDEX IF NOT EXISTS idx_gocafe_cafe_listings_city ON gocafe_cafe_listings (city);
CREATE INDEX IF NOT EXISTS idx_gocafe_cafe_listings_source_cafe_id ON gocafe_cafe_listings (source_cafe_id);

DO $$
BEGIN
    IF NOT EXISTS (
        SELECT 1
        FROM pg_constraint
        WHERE conname = 'chk_gocafe_ratings_rating_range'
    ) THEN
        ALTER TABLE gocafe_ratings
        ADD CONSTRAINT chk_gocafe_ratings_rating_range CHECK (rating BETWEEN 1 AND 5);
    END IF;
END $$;
