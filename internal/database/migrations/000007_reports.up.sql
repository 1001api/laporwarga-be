CREATE TABLE IF NOT EXISTS reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    title VARCHAR(255) NOT NULL,
    description TEXT NOT NULL,
    address TEXT,

    location GEOMETRY(POINT, 4326) NOT NULL,

    area_id UUID REFERENCES areas(id),
    category_id UUID NOT NULL REFERENCES categories(id),
    user_id UUID NOT NULL REFERENCES users(id),

    status VARCHAR(15) NOT NULL DEFAULT 'under_review'
        CHECK (status IN (
            'under_review',   -- waiting (new users / low score)
            'open',           -- live, visible to everyone
            'resolved',       -- closed successfully
            'hidden'          -- spam / troll / rejected by community
        )),

    view_count BIGINT DEFAULT 0,
    upvote_count BIGINT DEFAULT 0,
    downvote_count BIGINT DEFAULT 0,

    resolved_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_reports_location ON reports USING GIST(location);
CREATE INDEX IF NOT EXISTS idx_reports_category_id ON reports(category_id);
CREATE INDEX IF NOT EXISTS idx_reports_area_id ON reports(area_id);
CREATE INDEX IF NOT EXISTS idx_reports_status ON reports(status);
CREATE INDEX IF NOT EXISTS idx_reports_created_at ON reports(created_at);
CREATE INDEX IF NOT EXISTS idx_reports_user_id ON reports(user_id);

