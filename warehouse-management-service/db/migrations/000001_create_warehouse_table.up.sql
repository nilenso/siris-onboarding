CREATE TABLE IF NOT EXISTS warehouse (
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid (),
    name TEXT,
    geolocation POINT
);