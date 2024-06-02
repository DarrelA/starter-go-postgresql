CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS
  users (
    id SERIAL PRIMARY KEY,
    user_uuid UUID UNIQUE DEFAULT uuid_generate_v4 (),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL
  );

INSERT INTO
  users (first_name, last_name, password, email)
VALUES
  (
    'FirstName1',
    'LastName1',
    'Password1',
    'user1@example.com'
  ),
  (
    'FirstName2',
    'LastName2',
    'Password2',
    'user2@example.com'
  ) ON CONFLICT (email)
DO NOTHING;