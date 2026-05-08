DROP TRIGGER IF EXISTS trg_set_updated_at ON orders;

DROP TRIGGER IF EXISTS trg_set_updated_at ON products;

DROP FUNCTION IF EXISTS set_updated_at ();

DROP TABLE IF EXISTS order_items;

DROP TABLE IF EXISTS orders;

DROP TYPE IF EXISTS order_status;

DROP TABLE IF EXISTS products;
DROP INDEX IF EXISTS products_stock_idx;

DROP INDEX IF EXISTS order_items_order_id_idx;
DROP INDEX IF EXISTS order_items_product_id_idx;
DROP INDEX IF EXISTS orders_created_at_idx;
DROP INDEX IF EXISTS orders_status_idx;
 