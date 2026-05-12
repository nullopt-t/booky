CREATE TABLE
  products (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
    title TEXT NOT NULL,
    price INT NOT NULL,
    stock INT NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now ()
  );

CREATE INDEX IF NOT EXISTS products_stock_idx ON products (stock)
WHERE
  stock < 10;