CREATE TABLE IF NOT EXISTS item(
    id TEXT PRIMARY KEY DEFAULT gen_random_uuid (),
    sku TEXT NOT NULL references product(sku),
    expiration_DATE DATE,
    received_on TIMESTAMP NOT NULL,
    shelf_id TEXT NOT NULL references shelf(id)
);