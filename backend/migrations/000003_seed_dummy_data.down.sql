-- Delete all order items referencing seeded products first (child table)
-- This also cleans up any test orders that use the seeded products
DELETE FROM order_items
WHERE
    product_id IN (
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a03',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a04',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a05',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a06',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a07',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a08'
    );

-- Delete seeded orders (orders with no items will cascade or fail safely)
DELETE FROM orders
WHERE
    id IN (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b01',
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b02',
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b03',
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b04',
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b05',
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b06',
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b07',
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b08'
    );

-- Delete seeded products
DELETE FROM products
WHERE
    id IN (
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a03',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a04',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a05',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a06',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a07',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a08'
    );