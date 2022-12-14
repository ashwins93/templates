-- migrate:up
CREATE TABLE IF NOT EXISTS users (
  id TEXT NOT NULL PRIMARY KEY,
  name TEXT,
  email TEXT UNIQUE NOT NULL,
  password TEXT NOT NULL,
  created_at INTEGER NOT NULL DEFAULT (UNIXEPOCH()),
  updated_at INTEGER NOT NULL DEFAULT (UNIXEPOCH())
) WITHOUT ROWID;
-- migrate:down
DROP TABLE IF EXISTS users;