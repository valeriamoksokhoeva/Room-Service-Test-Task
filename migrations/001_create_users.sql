-- +migrate Up
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    role VARCHAR(255) NOT NULL CHECK (role IN ('admin', 'user')),
    created_at TIMESTAMP DEFAULT NOW()
);

-- +migrate Down
DROP TABLE IF EXISTS users;