-- Seed dummy products
INSERT INTO
    products (id, title, price, stock, created_at, updated_at)
VALUES
    (
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01',
        'The Go Programming Language',
        4500,
        25,
        NOW () - INTERVAL '30 days',
        NOW ()
    ),
    (
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02',
        'Clean Code',
        3200,
        15,
        NOW () - INTERVAL '25 days',
        NOW ()
    ),
    (
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a03',
        'Design Patterns',
        2800,
        8,
        NOW () - INTERVAL '20 days',
        NOW ()
    ),
    (
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a04',
        'PostgreSQL Up & Running',
        3900,
        12,
        NOW () - INTERVAL '15 days',
        NOW ()
    ),
    (
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a05',
        'Microservices Patterns',
        5200,
        5,
        NOW () - INTERVAL '10 days',
        NOW ()
    ),
    (
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a06',
        'Kubernetes in Action',
        4800,
        0,
        NOW () - INTERVAL '5 days',
        NOW ()
    ),
    (
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a07',
        'Domain-Driven Design',
        4200,
        20,
        NOW () - INTERVAL '3 days',
        NOW ()
    ),
    (
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a08',
        'System Design Interview',
        3500,
        50,
        NOW () - INTERVAL '1 day',
        NOW ()
    );

-- Seed dummy orders with various statuses
INSERT INTO
    orders (id, status, total_price, created_at, updated_at)
VALUES
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b01',
        'delivered',
        7700,
        NOW () - INTERVAL '20 days',
        NOW () - INTERVAL '15 days'
    ),
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b02',
        'shipped',
        4500,
        NOW () - INTERVAL '10 days',
        NOW () - INTERVAL '2 days'
    ),
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b03',
        'pending',
        9700,
        NOW () - INTERVAL '2 days',
        NOW ()
    ),
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b04',
        'confirmed',
        8400,
        NOW () - INTERVAL '5 days',
        NOW () - INTERVAL '1 day'
    ),
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b05',
        'cancelled',
        3200,
        NOW () - INTERVAL '8 days',
        NOW () - INTERVAL '7 days'
    ),
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b06',
        'processing',
        6700,
        NOW () - INTERVAL '3 days',
        NOW ()
    ),
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b07',
        'paid',
        9100,
        NOW () - INTERVAL '4 days',
        NOW () - INTERVAL '3 days'
    ),
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b08',
        'pending',
        2800,
        NOW () - INTERVAL '1 day',
        NOW ()
    );

-- Seed order items linking orders to products
INSERT INTO
    order_items (order_id, product_id, quantity, purchase_price)
VALUES
    -- Order 1: delivered (Clean Code + Design Patterns)
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b01',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02',
        1,
        3200
    ),
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b01',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a03',
        1,
        2800
    ),
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b01',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a03',
        1,
        2800
    ),
    -- Order 2: shipped (The Go Programming Language)
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b02',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01',
        1,
        4500
    ),
    -- Order 3: pending (Multiple items)
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b03',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a01',
        1,
        4500
    ),
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b03',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a04',
        1,
        3900
    ),
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b03',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a05',
        1,
        5200
    ),
    -- Order 4: confirmed (PostgreSQL + Clean Code)
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b04',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a04',
        1,
        3900
    ),
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b04',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02',
        1,
        3200
    ),
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b04',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02',
        1,
        3200
    ),
    -- Order 5: cancelled (Clean Code)
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b05',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a02',
        1,
        3200
    ),
    -- Order 6: processing (Kubernetes + Design Patterns)
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b06',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a06',
        1,
        4800
    ),
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b06',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a03',
        1,
        2800
    ),
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b06',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a04',
        1,
        3900
    ),
    -- Order 7: paid (Domain-Driven Design + System Design)
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b07',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a07',
        1,
        4200
    ),
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b07',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a08',
        1,
        3500
    ),
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b07',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a07',
        1,
        4200
    ),
    -- Order 8: pending (Design Patterns)
    (
        'b0eebc99-9c0b-4ef8-bb6d-6bb9bd380b08',
        'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a03',
        1,
        2800
    );