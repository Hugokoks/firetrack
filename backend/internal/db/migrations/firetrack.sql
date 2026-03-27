CREATE TABLE users (
id UUID PRIMARY KEY,
name TEXT NOT NULL,
email TEXT NOT NULL UNIQUE,
password_hash TEXT NOT NULL,
role TEXT NOT NULL CHECK (role IN ('admin', 'technician')),
phone TEXT,
is_active BOOLEAN NOT NULL DEFAULT TRUE,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_users_email ON users(email);

--admin123
SELECT * FROM users

INSERT INTO users (
    id,
    name,
    email,
    password_hash,
    role,
    is_active,
    created_at,
    updated_at
)
VALUES (
    gen_random_uuid(),
    'Admin',
    'admin@firetrack.local',
    '$2a$10$qMEMdNetE7WISfSmE3ledeuq/nFv68tdDNoZPiZTFvuJM5FroX39u',
    'admin',
    TRUE,
    NOW(),
    NOW()
);

CREATE TABLE jobs (
    id UUID PRIMARY KEY,
    job_number TEXT UNIQUE,
    title TEXT NOT NULL,
    customer_name TEXT,
    address TEXT NOT NULL,
    city TEXT,
    zip TEXT,
    country TEXT NOT NULL DEFAULT 'CZ',

    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,

    scheduled_start TIMESTAMP NOT NULL,
    scheduled_end TIMESTAMP,

    completed_at TIMESTAMP,

    status TEXT NOT NULL CHECK (status IN ('planned', 'in_progress', 'done', 'cancelled')) DEFAULT 'planned',
    priority TEXT NOT NULL CHECK (priority IN ('low', 'normal', 'high', 'urgent')) DEFAULT 'normal',

    assigned_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,

    description TEXT,
    google_event_id TEXT,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_jobs_status ON jobs(status);
CREATE INDEX idx_jobs_assigned_user_id ON jobs(assigned_user_id);
CREATE INDEX idx_jobs_scheduled_start ON jobs(scheduled_start);
CREATE INDEX idx_jobs_customer_name ON jobs(customer_name);


CREATE TABLE job_notes (
    id UUID PRIMARY KEY,
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    content TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_job_notes_job_id ON job_notes(job_id);




CREATE TABLE job_files  (
    id UUID PRIMARY KEY,
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    uploaded_by UUID REFERENCES users(id) ON DELETE SET NULL,
    file_name TEXT NOT NULL,
    stored_name TEXT NOT NULL,
    file_path TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    file_size BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_job_files_job_id ON job_files(job_id);


CREATE TABLE job_activity (
    id UUID PRIMARY KEY,
    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action_type TEXT NOT NULL,
    action_label TEXT NOT NULL,
    meta JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE INDEX idx_job_activity_job_id ON job_activity(job_id);




CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at);

select * from sessions



CREATE TABLE job_map_points (
    id UUID PRIMARY KEY,

    job_id UUID NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,

    label TEXT,
    note TEXT,

    latitude DOUBLE PRECISION NOT NULL,
    longitude DOUBLE PRECISION NOT NULL,

    color TEXT, -- např. 'red', 'green', '#ff0000'

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_job_map_points_job_id ON job_map_points(job_id);




CREATE ROLE firetrack_app LOGIN PASSWORD 'silneheslo';
GRANT CONNECT ON DATABASE railway TO firetrack_app;
GRANT USAGE ON SCHEMA public TO firetrack_app;
GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO firetrack_app;
GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO firetrack_app;