CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email CITEXT NOT NULL UNIQUE CHECK (
        email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'
    ),
    
    role TEXT NOT NULL DEFAULT 'customer' CHECK(role IN ('customer', 'admin', 'vendor')),
    
    password_hash TEXT NOT NULL,
    
    -- password resetting columns
    reset_token TEXT NULL,
    reset_token_expires_at TIMESTAMPTZ NULL,

    failed_reset_attempts INT DEFAULT 0,
    last_reset_request_at TIMESTAMPTZ NULL,

    is_inactive BOOLEAN DEFAULT FALSE,
    failed_login_attempts INT DEFAULT 0,
    locked_until TIMESTAMPTZ NULL,
    deleted_at TIMESTAMPTZ NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- CREATE INDEX users_active_idx ON users(id)
-- WHERE is_inactive = false
--     AND deleted_at IS NULL;
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER AS $$ BEGIN NEW.updated_at = NOW();
RETURN NEW;
END;
$$ LANGUAGE plpgsql;
CREATE TRIGGER users_updated_at BEFORE
UPDATE ON users FOR EACH ROW EXECUTE FUNCTION set_updated_at();
