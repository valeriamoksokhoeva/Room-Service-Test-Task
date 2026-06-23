-- +migrate Up
CREATE TABLE bookings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    slot_id UUID NOT NULL REFERENCES slots(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(50) NOT NULL CHECK (status IN ('active', 'cancelled')),
    conference_link TEXT,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX ON bookings(user_id);

CREATE UNIQUE INDEX ON bookings(slot_id) WHERE status = 'active';

-- +migrate Down
DROP TABLE IF EXISTS bookings;