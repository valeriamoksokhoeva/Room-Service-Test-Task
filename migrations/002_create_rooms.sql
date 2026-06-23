-- +migrate Up
CREATE TABLE rooms (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255),
    capacity INTEGER,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +migrate Down
DROP TABLE IF EXISTS rooms;