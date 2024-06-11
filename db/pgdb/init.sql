CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS
  users (
    id SERIAL PRIMARY KEY,
    user_uuid UUID UNIQUE DEFAULT uuid_generate_v4 (),
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    updated_at TIMESTAMP NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
  );

INSERT INTO
  users (
    first_name,
    last_name,
    password,
    email,
    created_at,
    updated_at
  )
VALUES
  (
    'FirstName1',
    'LastName1',
    'Password1',
    'user1@example.com',
    now() AT TIME ZONE 'UTC',
    now() AT TIME ZONE 'UTC'
  ),
  (
    'FirstName2',
    'LastName2',
    'Password2',
    'user2@example.com',
    now() AT TIME ZONE 'UTC',
    now() AT TIME ZONE 'UTC'
  ) ON CONFLICT (email)
DO NOTHING;