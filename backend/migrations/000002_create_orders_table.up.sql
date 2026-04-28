CREATE TYPE order_status AS ENUM (
    'pending',
    'confirmed',
    'paid',
    'processing',
    'shipped',
    'delivered',
    'cancelled',
    'refunded'
);

CREATE TABLE
    orders (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        status order_status NOT NULL DEFAULT 'pending',
        total_price INT NOT NULL,
        created_at timestamptz NOT NULL DEFAULT now (),
        updated_at timestamptz NOT NULL DEFAULT now (),
        deleted_at timestamptz NULL
    );

CREATE OR REPLACE FUNCTION set_updated_at ()
RETURNS trigger AS $$
BEGIN
    NEW.updated_at = now ();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_set_updated_at
BEFORE
UPDATE
    ON orders FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TABLE order_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    order_id UUID NOT NULL,
    product_id UUID NOT NULL,
    quantity INT NOT NULL,
    purchase_price INT NOT NULL,
    FOREIGN KEY (order_id) REFERENCES orders (id),
    FOREIGN KEY (product_id) REFERENCES products (id)
);