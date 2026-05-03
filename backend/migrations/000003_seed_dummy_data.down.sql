-- Delete seeded order items first (child table)
DELETE FROM order_items 
WHERE order_id IN (
    'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b01',
    'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b02',
    'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b03',
    'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b04',
    'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b05',
    'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b06',
    'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b07',
    'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b08'
);

-- Delete seeded orders
DELETE FROM orders 
WHERE id IN (
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
WHERE id IN (
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a03',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a04',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a05',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a06',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a07',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a08'
);
