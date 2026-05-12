CREATE TABLE
    IF NOT EXISTS inventories (
        product_id UUID PRIMARY KEY REFERENCES products (id),
        available_quantity INTEGER NOT NULL DEFAULT 0 CHECK (available_quantity >= 0),
        reserved_quantity INTEGER NOT NULL DEFAULT 0 CHECK (reserved_quantity >= 0),
        created_at TIMESTAMP NOT NULL DEFAULT NOW (),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW ()
    );

CREATE INDEX IF NOT EXISTS inventories_product_id_idx ON inventories (product_id);

CREATE INDEX IF NOT EXISTS inventories_available_quantity_idx ON inventories (available_quantity);

CREATE INDEX IF NOT EXISTS inventories_reserved_quantity_idx ON inventories (reserved_quantity);