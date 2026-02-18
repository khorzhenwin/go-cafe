ALTER TABLE gocafe_cafe_listings
DROP CONSTRAINT IF EXISTS chk_gocafe_cafe_listings_visit_status;

DROP INDEX IF EXISTS idx_gocafe_cafe_listings_visit_status;

ALTER TABLE gocafe_cafe_listings
DROP COLUMN IF EXISTS visit_status;
