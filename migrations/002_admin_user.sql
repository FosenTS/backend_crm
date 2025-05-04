-- Insert admin user with Director role (0)
INSERT INTO users (role, username, pass_hash)
VALUES (
    0, -- Director role
    'admin',
    '$2a$10$1rNzpZOgK6y47J.lgk9Vq.bGxcBvK5C6kt/Esz73rp/RlAhiDIkJi' -- bcrypt hash of 'admin123'
)
ON CONFLICT (username) DO NOTHING; 