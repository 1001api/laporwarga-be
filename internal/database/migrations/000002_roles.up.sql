CREATE TABLE IF NOT EXISTS roles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    name VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);