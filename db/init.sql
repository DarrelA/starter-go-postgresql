-- init.sql
CREATE TABLE IF NOT EXISTS
  users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL,
    email VARCHAR(50) NOT NULL UNIQUE
  );

INSERT INTO
  users (username, email)
VALUES
  ('user1', 'user1@example.com'),
  ('user2', 'user2@example.com') ON CONFLICT (email)
DO NOTHING;