CREATE TABLE IF NOT EXISTS report_votes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    report_id UUID NOT NULL REFERENCES reports(id),
    user_id UUID NOT NULL REFERENCES users(id),
    vote_type VARCHAR(10) NOT NULL CHECK (vote_type IN ('upvote', 'downvote')),
    created_at TIMESTAMPTZ DEFAULT NOW(),

    UNIQUE (report_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_report_votes_user_id ON report_votes(user_id);
CREATE INDEX IF NOT EXISTS idx_report_votes_report_id ON report_votes(report_id);

