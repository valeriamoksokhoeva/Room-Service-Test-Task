-- +migrate Up
CREATE TABLE schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    room_id UUID NOT NULL UNIQUE REFERENCES rooms(id),
    days_of_week INT[],
    start_time VARCHAR(255) NOT NULL,
    end_time VARCHAR(255) NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS schedules;