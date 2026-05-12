CREATE TABLE
    IF NOT EXISTS carts (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        created_at timestamptz NOT NULL DEFAULT now (),
        updated_at timestamptz NOT NULL DEFAULT now (),
        user_id UUID NOT NULL,
        FOREIGN KEY (user_id) REFERENCES users (id)
    );

CREATE TABLE
    IF NOT EXISTS cart_items (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        cart_id UUID NOT NULL,
        product_id UUID NOT NULL,
        quantity INTEGER NOT NULL DEFAULT 1,
        created_at timestamptz NOT NULL DEFAULT now (),
        updated_at timestamptz NOT NULL DEFAULT now (),
        FOREIGN KEY (cart_id) REFERENCES carts (id)
    );

CREATE INDEX IF NOT EXISTS cart_items_cart_id_idx ON cart_items (cart_id);

CREATE INDEX IF NOT EXISTS cart_items_product_id_idx ON cart_items (product_id);