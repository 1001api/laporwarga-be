CREATE TABLE IF NOT EXISTS report_spam_flags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id UUID NOT NULL REFERENCES reports(id),
    user_id UUID NOT NULL REFERENCES users(id),
    reason VARCHAR(50) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE (report_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_report_spam_flags_report_id ON report_spam_flags(report_id);
CREATE INDEX IF NOT EXISTS idx_report_spam_flags_user_id ON report_spam_flags(user_id);