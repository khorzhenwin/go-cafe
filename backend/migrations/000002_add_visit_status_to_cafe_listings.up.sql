ALTER TABLE gocafe_cafe_listings
ADD COLUMN IF NOT EXISTS visit_status VARCHAR(32) NOT NULL DEFAULT 'to_visit';

UPDATE gocafe_cafe_listings
SET visit_status = 'to_visit'
WHERE visit_status IS NULL OR visit_status = '';

ALTER TABLE gocafe_cafe_listings
DROP CONSTRAINT IF EXISTS chk_gocafe_cafe_listings_visit_status;

ALTER TABLE gocafe_cafe_listings
ADD CONSTRAINT chk_gocafe_cafe_listings_visit_status
CHECK (visit_status IN ('to_visit', 'visited'));

CREATE INDEX IF NOT EXISTS idx_gocafe_cafe_listings_visit_status
ON gocafe_cafe_listings (visit_status);
