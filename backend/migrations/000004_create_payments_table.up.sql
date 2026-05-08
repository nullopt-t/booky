CREATE TABLE
    IF NOT EXISTS payments (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        order_id UUID NOT NULL,
        status VARCHAR(50) NOT NULL,
        amount INT NOT NULL,
        provider VARCHAR(50) NOT NULL,
        provider_ref VARCHAR(255) NOT NULL,
        UNIQUE (provider, provider_ref),
        created_at timestamptz NOT NULL DEFAULT now (),
        updated_at timestamptz NOT NULL DEFAULT now ()
    );