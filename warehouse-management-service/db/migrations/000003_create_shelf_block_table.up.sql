CREATE TABLE IF NOT EXISTS shelf_block(
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid (),
    aisle TEXT,
    rack TEXT,
    storage_type TEXT
);