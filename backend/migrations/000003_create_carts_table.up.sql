CREATE TABLE IF NOT EXISTS carts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    user_id UUID NOT NULL UNIQUE REFERENCES users (id),
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now ()
);
CREATE TABLE IF NOT EXISTS cart_items (
    cart_id UUID NOT NULL REFERENCES carts (id),
    product_id UUID NOT NULL REFERENCES products (id),
    quantity INTEGER NOT NULL DEFAULT 1 CHECK (quantity > 0),
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now (),
    PRIMARY KEY (cart_id, product_id)
);
CREATE INDEX IF NOT EXISTS cart_items_cart_id_idx ON cart_items (cart_id);
CREATE INDEX IF NOT EXISTS cart_items_product_id_idx ON cart_items (product_id);