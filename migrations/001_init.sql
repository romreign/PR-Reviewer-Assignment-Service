CREATE TABLE IF NOT EXISTS teams (
  team_name TEXT PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS users (
  user_id TEXT PRIMARY KEY,
  username TEXT NOT NULL,
  team_name TEXT REFERENCES teams(team_name),
  is_active BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE IF NOT EXISTS pull_requests (
  pull_request_id TEXT PRIMARY KEY,
  pull_request_name TEXT NOT NULL,
  author_id TEXT NOT NULL REFERENCES users(user_id),
  status TEXT NOT NULL CHECK (status IN ('OPEN', 'MERGED')) DEFAULT 'OPEN',
  created_at TIMESTAMP,
  merged_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS pr_reviewers (
  pull_request_id TEXT NOT NULL REFERENCES pull_requests(pull_request_id) ON DELETE CASCADE,
  user_id TEXT NOT NULL REFERENCES users(user_id),
  PRIMARY KEY(pull_request_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_users_team_name ON users(team_name);
CREATE INDEX IF NOT EXISTS idx_pull_requests_author_id ON pull_requests(author_id);
CREATE INDEX IF NOT EXISTS idx_pull_requests_status ON pull_requests(status);
CREATE INDEX IF NOT EXISTS idx_pr_reviewers_user_id ON pr_reviewers(user_id);