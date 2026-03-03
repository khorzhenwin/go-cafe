ALTER TABLE gocafe_cafe_listings
ADD COLUMN IF NOT EXISTS address TEXT DEFAULT '';
