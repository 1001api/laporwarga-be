CREATE TABLE IF NOT EXISTS report_comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    report_id UUID NOT NULL REFERENCES reports(id),
    parent_id UUID REFERENCES report_comments(id),
    user_id UUID REFERENCES users(id),
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_report_comments_report_id ON report_comments(report_id);
CREATE INDEX IF NOT EXISTS idx_report_comments_user_id ON report_comments(user_id);
CREATE INDEX IF NOT EXISTS idx_report_comments_created_at ON report_comments(created_at);

