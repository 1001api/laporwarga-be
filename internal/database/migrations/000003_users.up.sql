CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    email_hash TEXT UNIQUE NOT NULL,  -- must be SHA-256 hash
    email_enc BYTEA NOT NULL,         -- must be AES-encrypted
    fullname_hash TEXT,               -- must be SHA-256 hash
    fullname_enc BYTEA NOT NULL,      -- must be AES-encrypted
    username VARCHAR(50) UNIQUE NOT NULL,
    phone_hash TEXT,                  -- must be SHA-256 hash
    phone_enc BYTEA,                  -- must be AES-encrypted

    credibility_score SMALLINT DEFAULT 50,
    status VARCHAR(20) DEFAULT 'probation' CHECK (status IN ('probation', 'regular', 'suspended')),

    auth_provider VARCHAR(50),
    oauth_id VARCHAR(255),

    password_hash VARCHAR(255),
    password_changed_at TIMESTAMPTZ,

    is_email_verified BOOLEAN DEFAULT FALSE,
    is_phone_verified BOOLEAN DEFAULT FALSE,
    last_login_at TIMESTAMPTZ,
    failed_login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMPTZ,

    role_id UUID REFERENCES roles(id) ON DELETE SET NULL,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    last_updated_at TIMESTAMPTZ,
    last_updated_by UUID REFERENCES users(id),
    deleted_at TIMESTAMPTZ,
    deleted_by UUID REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_users_username ON users(username);
CREATE INDEX IF NOT EXISTS idx_users_email_hash ON users(email_hash);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_role_id ON users(role_id);

CREATE TABLE IF NOT EXISTS officials (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    department TEXT,
    job_title TEXT,

    created_at TIMESTAMPTZ DEFAULT NOW(),
    created_by UUID REFERENCES users(id),
    last_updated_at TIMESTAMPTZ,
    last_updated_by UUID REFERENCES users(id),
    deleted_at TIMESTAMPTZ,
    deleted_by UUID REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    entity_name TEXT NOT NULL,
    entity_id UUID NOT NULL,
    action TEXT NOT NULL,
    metadata JSONB,
    performed_by UUID REFERENCES users(id),

    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_logs_entity ON audit_logs(entity_name, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_logs_performed_by ON audit_logs(performed_by);