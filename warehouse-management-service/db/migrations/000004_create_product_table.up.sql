CREATE TABLE IF NOT EXISTS product(
    sku TEXT PRIMARY KEY DEFAULT gen_random_uuid (),
    name TEXT NOT NULL,
    mrp NUMERIC NOT NULL,
    variant TEXT,
    length_in_cm NUMERIC,
    width_in_cm NUMERIC,
    height_in_cm NUMERIC,
    weight_in_kg NUMERIC,
    perishable BOOLEAN NOT NULL
);