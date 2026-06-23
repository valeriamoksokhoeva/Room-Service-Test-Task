-- +migrate Up
CREATE TABLE slots (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    start_at TIMESTAMP NOT NULL,
    end_at TIMESTAMP NOT NULL
);

CREATE INDEX ON slots(room_id, start_at);

-- +migrate Down
DROP TABLE IF EXISTS slots;