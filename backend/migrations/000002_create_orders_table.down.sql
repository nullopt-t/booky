DROP TRIGGER IF EXISTS trg_set_updated_at ON orders;

DROP TRIGGER IF EXISTS trg_set_updated_at ON products;

DROP FUNCTION IF EXISTS set_updated_at ();

DROP TABLE IF EXISTS orders;

DROP TYPE IF EXISTS order_status;