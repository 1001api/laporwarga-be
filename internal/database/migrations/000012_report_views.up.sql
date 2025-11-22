CREATE TABLE IF NOT EXISTS report_views (
    report_id UUID NOT NULL REFERENCES reports(id),
    session_id UUID NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id),
    viewed_at TIMESTAMPTZ DEFAULT NOW(),

    PRIMARY KEY (report_id, session_id)
);

CREATE INDEX IF NOT EXISTS idx_report_views_user_id ON report_views(user_id);
CREATE INDEX IF NOT EXISTS idx_report_views_report_id ON report_views(report_id);