
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_tables WHERE schemaname = 'public' AND tablename = 'users') THEN
        DELETE FROM users
        WHERE
            id IN (
                '20eebc99-9c0b-4ef8-bb6d-6bb9bd380b01'::uuid,
                '30eebc99-9c0b-4ef8-bb6d-6bb9bd380b01'::uuid,
                '40eebc99-9c0b-4ef8-bb6d-6bb9bd380b01'::uuid
            );
    END IF;
END $$;


-- Delete all order items referencing seeded products first (child table)
-- This also cleans up any test orders that use the seeded products
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_tables WHERE schemaname = 'public' AND tablename = 'order_items') THEN
        DELETE FROM order_items
        WHERE
            product_id IN (
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a03'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a04'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a05'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a06'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a07'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a08'::uuid
            );
    END IF;
END $$;


-- Check if the order items table exists before deleting
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_tables WHERE schemaname = 'public' AND tablename = 'order_items') THEN
        DELETE FROM order_items
        WHERE
            product_id IN (
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a03'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a04'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a05'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a06'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a07'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a08'::uuid
            );
    END IF;
END $$;

-- Check if the orders table exists before deleting
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_tables WHERE schemaname = 'public' AND tablename = 'orders') THEN
        DELETE FROM orders
        WHERE
            id IN (
                'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b01'::uuid,
                'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b02'::uuid,
                'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b03'::uuid,
                'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b04'::uuid,
                'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b05'::uuid,
                'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b06'::uuid,
                'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b07'::uuid,
                'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b08'::uuid
            );
    END IF;
END $$;

-- Check if the products table exists before deleting
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM pg_tables WHERE schemaname = 'public' AND tablename = 'products') THEN
        DELETE FROM products
        WHERE
            id IN (
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a03'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a04'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a05'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a06'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a07'::uuid,
                'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a08'::uuid
            );
    END IF;
END $$;