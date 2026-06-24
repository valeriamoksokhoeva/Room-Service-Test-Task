-- +migrate Up
INSERT INTO users (id, email, password, role, created_at)
VALUES 
    ('a0000000-0000-0000-0000-000000000001', 'admin@dummy.local', 'dummy', 'admin', NOW()),
    ('b0000000-0000-0000-0000-000000000002', 'user@dummy.local',  'dummy', 'user',  NOW())
ON CONFLICT (id) DO NOTHING;

-- +migrate Down
DELETE FROM users WHERE id IN (
    'a0000000-0000-0000-0000-000000000001',
    'b0000000-0000-0000-0000-000000000002'
);