CREATE TABLE IF NOT EXISTS users
(
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  version BIGINT NOT NULL DEFAULT 0,
  login VARCHAR(30) UNIQUE NOT NULL,
  pass_hash VARCHAR(100) NOT NULL,
  full_name VARCHAR(100) NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts
(
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  version BIGINT NOT NULL DEFAULT 0,
  user_id BIGINT REFERENCES users(id) NOT NULL,
  title VARCHAR(50) NOT NULL,
  description VARCHAR(500),
  icon VARCHAR(50),
  color VARCHAR(7),
  balance NUMERIC(20,2) NOT NULL,

  UNIQUE(id, user_id)
);

CREATE TABLE IF NOT EXISTS operation_types
(
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  title VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS operation_categories
(
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  version BIGINT NOT NULL DEFAULT 0,
  user_id BIGINT REFERENCES users(id),
  operation_type_id BIGINT REFERENCES operation_types(id) NOT NULL,
  title VARCHAR(50) NOT NULL,
  icon VARCHAR(50),
  color VARCHAR(7)
);

CREATE TABLE IF NOT EXISTS operations
(
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  version BIGINT NOT NULL DEFAULT 0,
  user_id BIGINT NOT NULL,
  account_id BIGINT NOT NULL,
  category_id BIGINT REFERENCES operation_categories(id) NOT NULL,
  title VARCHAR(50) NOT NULL,
  description VARCHAR(500),
  exec_at TIMESTAMPTZ NOT NULL,
  sum NUMERIC(20,2) NOT NULL,

  CONSTRAINT operations_account_owner_fk
        FOREIGN KEY (account_id, user_id)
        REFERENCES accounts (id, user_id)
);

INSERT INTO operation_types (title) VALUES ('income');
INSERT INTO operation_types (title) VALUES ('expense');