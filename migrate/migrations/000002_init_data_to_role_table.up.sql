INSERT INTO roles (name, description)
VALUES ('admin', 'Administrator'), ('user', 'Normal user')
ON CONFLICT DO NOTHING;