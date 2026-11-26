BEGIN;

-- Menjaga kolom updated_at selalu terisi otomatis.
CREATE OR REPLACE FUNCTION set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE IF NOT EXISTS users (
    id            BIGSERIAL PRIMARY KEY,
    name          VARCHAR(120) NOT NULL,
    email         VARCHAR(150) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT users_email_unique UNIQUE (email)
);

CREATE TABLE IF NOT EXISTS movies (
    id               BIGSERIAL PRIMARY KEY,
    title            VARCHAR(200) NOT NULL,
    description      TEXT,
    genre            VARCHAR(100),
    duration_minutes INTEGER NOT NULL DEFAULT 0,
    release_date     DATE,
    created_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT movies_title_unique UNIQUE (title),
    CONSTRAINT movies_duration_check CHECK (duration_minutes >= 0)
);

CREATE TABLE IF NOT EXISTS branches (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(150) NOT NULL,
    city       VARCHAR(100) NOT NULL,
    address    TEXT,
    phone      VARCHAR(30),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT branches_name_city_unique UNIQUE (name, city)
);

CREATE TABLE IF NOT EXISTS studios (
    id          BIGSERIAL PRIMARY KEY,
    branch_id   BIGINT NOT NULL REFERENCES branches (id) ON DELETE CASCADE,
    name        VARCHAR(100) NOT NULL,
    total_seats INTEGER NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT studios_branch_name_unique UNIQUE (branch_id, name),
    CONSTRAINT studios_total_seats_check CHECK (total_seats >= 0)
);

CREATE TABLE IF NOT EXISTS schedules (
    id           BIGSERIAL PRIMARY KEY,
    movie_id     BIGINT NOT NULL REFERENCES movies (id) ON DELETE RESTRICT,
    branch_id    BIGINT NOT NULL REFERENCES branches (id) ON DELETE RESTRICT,
    studio_id    BIGINT NOT NULL REFERENCES studios (id) ON DELETE RESTRICT,
    show_date    DATE NOT NULL,
    show_time    TIME NOT NULL,
    ticket_price BIGINT NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT schedules_ticket_price_check CHECK (ticket_price >= 0),
    CONSTRAINT schedules_studio_slot_unique UNIQUE (studio_id, show_date, show_time)
);

CREATE INDEX IF NOT EXISTS idx_schedules_movie_id ON schedules (movie_id);
CREATE INDEX IF NOT EXISTS idx_schedules_branch_id ON schedules (branch_id);
CREATE INDEX IF NOT EXISTS idx_schedules_show_date ON schedules (show_date);

DROP TRIGGER IF EXISTS trg_users_updated_at ON users;
CREATE TRIGGER trg_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_movies_updated_at ON movies;
CREATE TRIGGER trg_movies_updated_at BEFORE UPDATE ON movies
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_branches_updated_at ON branches;
CREATE TRIGGER trg_branches_updated_at BEFORE UPDATE ON branches
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_studios_updated_at ON studios;
CREATE TRIGGER trg_studios_updated_at BEFORE UPDATE ON studios
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_schedules_updated_at ON schedules;
CREATE TRIGGER trg_schedules_updated_at BEFORE UPDATE ON schedules
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMIT;
