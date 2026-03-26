ALTER TABLE gocafe_ratings
DROP CONSTRAINT IF EXISTS chk_gocafe_ratings_rating_range;

DROP INDEX IF EXISTS idx_gocafe_cafe_listings_source_cafe_id;
DROP INDEX IF EXISTS idx_gocafe_cafe_listings_city;

ALTER TABLE gocafe_cafe_listings
DROP CONSTRAINT IF EXISTS fk_gocafe_cafe_listings_source_cafe;

ALTER TABLE gocafe_cafe_listings
DROP COLUMN IF EXISTS source_cafe_id,
DROP COLUMN IF EXISTS image_url,
DROP COLUMN IF EXISTS longitude,
DROP COLUMN IF EXISTS latitude,
DROP COLUMN IF EXISTS neighborhood,
DROP COLUMN IF EXISTS city;
