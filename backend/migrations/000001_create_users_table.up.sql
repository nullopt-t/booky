CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email CITEXT NOT NULL UNIQUE,
    phone TEXT NULL UNIQUE,
    email_verified_at TIMESTAMPTZ NULL,
    phone_verified_at TIMESTAMPTZ NULL,
    role TEXT NOT NULL DEFAULT 'customer' CHECK(
        role IN ('customer', 'admin', 'vendor')
    ),
    status TEXT NOT NULL DEFAULT 'active' CHECK(
        status IN ('active', 'inactive', 'suspended', 'locked', 'deleted')
    ),
    suspended_until TIMESTAMPTZ NULL,
    locked_until TIMESTAMPTZ NULL,
    password_hash TEXT NOT NULL,
    password_changed_at TIMESTAMPTZ NULL,
    last_login_at TIMESTAMPTZ NULL,
    last_login_ip TEXT NULL,
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
