CREATE TABLE IF NOT EXISTS users
(
  id SERIAL PRIMARY KEY,
  version BIGINT NOT NULL DEFAULT 0,
  login VARCHAR(30) UNIQUE NOT NULL,
  pass_hash VARCHAR(100) NOT NULL,
  full_name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts
(
  id SERIAL PRIMARY KEY,
  version BIGINT NOT NULL DEFAULT 0,
  user_id INTEGER REFERENCES users(id),
  title VARCHAR(50) NOT NULL,
  description VARCHAR(500),
  icon VARCHAR(50),
  color VARCHAR(30),
  balance MONEY NOT NULL
);

CREATE TABLE IF NOT EXISTS operation_type
(
  id SERIAL PRIMARY KEY,
  title VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS operation_categories
(
  id SERIAL PRIMARY KEY,
  version BIGINT NOT NULL DEFAULT 0,
  user_id INTEGER REFERENCES users(id),
  operation_type_id INTEGER REFERENCES operation_type(id),
  title VARCHAR(50) NOT NULL,
  icon VARCHAR(50),
  color VARCHAR(3)
);

CREATE TABLE IF NOT EXISTS operations
(
  id SERIAL PRIMARY KEY,
  version BIGINT NOT NULL DEFAULT 0,
  account_id INTEGER REFERENCES accounts(id),
  category_id INTEGER REFERENCES operation_categories(id),
  title VARCHAR(50) NOT NULL,
  description VARCHAR(500),
  exec_at TIMESTAMPTZ NOT NULL,
  sum MONEY NOT NULL
);

INSERT INTO operation_type (title) VALUES ('income');
INSERT INTO operation_type (title) VALUES ('expense');