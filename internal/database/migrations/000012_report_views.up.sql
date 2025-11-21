CREATE TABLE IF NOT EXISTS report_views (
    report_id UUID PRIMARY KEY NOT NULL REFERENCES reports(id),
    session_id PRIMARY KEY NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    viewed_at TIMESTAMPTZ DEFAULT NOW(),
);

CREATE INDEX IF NOT EXISTS idx_report_views_user_id ON report_views(user_id);