CREATE TABLE
  products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    title TEXT NOT NULL,
    price INT NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now ()
  );

CREATE INDEX IF NOT EXISTS products_id_idx ON products (id);

CREATE INDEX IF NOT EXISTS products_title_idx ON products (title);

CREATE INDEX IF NOT EXISTS products_price_idx ON products (price);