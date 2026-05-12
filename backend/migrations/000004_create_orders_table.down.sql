DROP TABLE IF EXISTS order_items;

DROP TABLE IF EXISTS orders;

DROP TYPE IF EXISTS order_status;

DROP INDEX IF EXISTS order_items_order_id_idx;
DROP INDEX IF EXISTS order_items_product_id_idx;
DROP INDEX IF EXISTS orders_created_at_idx;
DROP INDEX IF EXISTS orders_status_idx;
 