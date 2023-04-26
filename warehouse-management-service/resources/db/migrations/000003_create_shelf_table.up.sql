CREATE TABLE IF NOT EXISTS shelf(
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid (),
    label TEXT,
    section TEXT,
    level TEXT,
    shelf_block TEXT references shelf_block(id)
);