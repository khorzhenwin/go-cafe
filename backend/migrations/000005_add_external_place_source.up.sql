ALTER TABLE gocafe_cafe_listings
ADD COLUMN IF NOT EXISTS source_provider VARCHAR(64) DEFAULT '',
ADD COLUMN IF NOT EXISTS external_place_id TEXT DEFAULT '';

CREATE INDEX IF NOT EXISTS idx_gocafe_cafe_listings_source_provider ON gocafe_cafe_listings (source_provider);
CREATE INDEX IF NOT EXISTS idx_gocafe_cafe_listings_external_place_id ON gocafe_cafe_listings (external_place_id);
