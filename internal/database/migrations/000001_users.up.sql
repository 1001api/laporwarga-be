CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "postgis";
CREATE EXTENSION IF NOT EXISTS "timescaledb";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,

    email BYTEA NOT NULL, -- must be encrypted
    full_name BYTEA NOT NULL, -- must be encrypted
    username VARCHAR(50) NOT NULL UNIQUE,
    phone_number BYTEA, -- must be encrypted

    credibility_score SMALLINT DEFAULT 50 CHECK (credibility_score BETWEEN 0 AND 100),

    status VARCHAR(20) DEFAULT 'probation' CHECK (status IN ('probation', 'regular', 'suspended')),
    role VARCHAR(50) DEFAULT 'citizen' CHECK (role IN ('citizen', 'admin', 'superadmin')),

    auth_provider VARCHAR(50),
    oauth_id VARCHAR(255),
    password_hash VARCHAR(255),

    is_email_verified BOOLEAN DEFAULT false,
    is_phone_verified BOOLEAN DEFAULT false,

    last_login_at TIMESTAMPTZ,
    password_changed_at TIMESTAMPTZ,
    failed_login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMPTZ,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_users_role_status ON users (role, status);
CREATE INDEX IF NOT EXISTS idx_users_email_enc ON users (email);

-- Function to update updated_at automatically
CREATE OR REPLACE FUNCTION trigger_set_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create trigger for updated_at
CREATE TRIGGER set_users_timestamp
BEFORE UPDATE ON users
FOR EACH ROW
EXECUTE PROCEDURE trigger_set_timestamp();