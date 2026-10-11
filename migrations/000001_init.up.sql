CREATE TABLE IF NOT EXISTS users
(
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  version BIGINT NOT NULL DEFAULT 0,
  login VARCHAR(30) UNIQUE NOT NULL,
  pass_hash VARCHAR(100) NOT NULL,
  full_name VARCHAR(100) NOT NULL,
  currency VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS accounts
(
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  version BIGINT NOT NULL DEFAULT 0,
  user_id BIGINT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
  title VARCHAR(50) NOT NULL,
  description VARCHAR(500),
  icon VARCHAR(50),
  color VARCHAR(7),
  balance NUMERIC(20,2) NOT NULL,

  UNIQUE(id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_accounts_user_id ON accounts (user_id);

CREATE TABLE IF NOT EXISTS operation_types
(
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  title VARCHAR(50) NOT NULL
);

CREATE TABLE IF NOT EXISTS operation_categories
(
  id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  version BIGINT NOT NULL DEFAULT 0,
  user_id BIGINT REFERENCES users(id) ON DELETE CASCADE NOT NULL,
  operation_type_id BIGINT REFERENCES operation_types(id) NOT NULL,
  title VARCHAR(50) NOT NULL,
  icon VARCHAR(50),
  color VARCHAR(7)
);

CREATE INDEX IF NOT EXISTS idx_operation_categories_user_type ON operation_categories (user_id, operation_type_id);

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
    ON DELETE CASCADE
);

CREATE INDEX idx_operations_account_id ON operations (account_id);
CREATE INDEX idx_operations_user_account_exec_at ON operations (user_id, account_id, exec_at ASC);

INSERT INTO operation_types (title) VALUES ('income');
INSERT INTO operation_types (title) VALUES ('expense');

CREATE OR REPLACE VIEW v_account_monthly_stats AS
SELECT
    a.id    AS account_id,
    a.user_id,
    date_trunc('month', o.exec_at) AS month,
    COALESCE(SUM(CASE WHEN ot.title = 'income'  THEN o.sum ELSE 0 END), 0) AS income,
    COALESCE(SUM(CASE WHEN ot.title = 'expense' THEN ABS(o.sum) ELSE 0 END), 0) AS expense,
    COALESCE(SUM(CASE WHEN ot.title = 'income'  THEN o.sum ELSE 0 END), 0)
      - COALESCE(SUM(CASE WHEN ot.title = 'expense' THEN ABS(o.sum) ELSE 0 END), 0) AS profit
FROM accounts a
JOIN operations o            ON o.account_id = a.id
JOIN operation_categories oc ON oc.id = o.category_id
JOIN operation_types ot      ON ot.id = oc.operation_type_id
GROUP BY a.id, a.user_id, a.title, date_trunc('month', o.exec_at);

CREATE OR REPLACE VIEW v_user_monthly_stats AS
SELECT
    u.id AS user_id,
    date_trunc('month', o.exec_at) AS month,
    COALESCE(SUM(CASE WHEN ot.title = 'income'  THEN o.sum ELSE 0 END), 0) AS income,
    COALESCE(SUM(CASE WHEN ot.title = 'expense' THEN ABS(o.sum) ELSE 0 END), 0) AS expense,
    COALESCE(SUM(CASE WHEN ot.title = 'income'  THEN o.sum ELSE 0 END), 0)
      - COALESCE(SUM(CASE WHEN ot.title = 'expense' THEN ABS(o.sum) ELSE 0 END), 0) AS profit
FROM users u
JOIN operations o            ON o.user_id = u.id
JOIN operation_categories oc ON oc.id = o.category_id
JOIN operation_types ot      ON ot.id = oc.operation_type_id
GROUP BY u.id, date_trunc('month', o.exec_at);