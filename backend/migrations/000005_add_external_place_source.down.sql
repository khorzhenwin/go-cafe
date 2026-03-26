DROP INDEX IF EXISTS idx_gocafe_cafe_listings_external_place_id;
DROP INDEX IF EXISTS idx_gocafe_cafe_listings_source_provider;

ALTER TABLE gocafe_cafe_listings
DROP COLUMN IF EXISTS external_place_id,
DROP COLUMN IF EXISTS source_provider;
