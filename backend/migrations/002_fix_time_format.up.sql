-- +goose Up
-- Fix time format: strip Go's timezone suffix and normalize to UTC

-- ledgers.occurred_at: format "2026-05-23 11:11:05.039 +0000 UTC" -> "2026-05-23 11:11:05"
UPDATE ledgers
SET occurred_at = datetime(substr(occurred_at, 1, 19))
WHERE occurred_at IS NOT NULL AND occurred_at != '';

-- todos.due_date
UPDATE todos
SET due_date = datetime(substr(due_date, 1, 19))
WHERE due_date IS NOT NULL AND due_date != '';

-- todos.completed_at
UPDATE todos
SET completed_at = datetime(substr(completed_at, 1, 19))
WHERE completed_at IS NOT NULL AND completed_at != '';

-- votes.deadline
UPDATE votes
SET deadline = datetime(substr(deadline, 1, 19))
WHERE deadline IS NOT NULL AND deadline != '';

-- Auto-managed timestamps (created_at, updated_at) use +0800 CST format.
-- Strip suffix and convert to UTC by subtracting 8 hours.
-- members
UPDATE members SET created_at = datetime(substr(created_at, 1, 19), '-8 hours') WHERE created_at IS NOT NULL;
UPDATE members SET updated_at = datetime(substr(updated_at, 1, 19), '-8 hours') WHERE updated_at IS NOT NULL;
UPDATE members SET last_login = datetime(substr(last_login, 1, 19), '-8 hours') WHERE last_login IS NOT NULL;

-- categories
UPDATE categories SET created_at = datetime(substr(created_at, 1, 19), '-8 hours') WHERE created_at IS NOT NULL;
UPDATE categories SET updated_at = datetime(substr(updated_at, 1, 19), '-8 hours') WHERE updated_at IS NOT NULL;

-- ledgers (created_at, updated_at)
UPDATE ledgers SET created_at = datetime(substr(created_at, 1, 19), '-8 hours') WHERE created_at IS NOT NULL;
UPDATE ledgers SET updated_at = datetime(substr(updated_at, 1, 19), '-8 hours') WHERE updated_at IS NOT NULL;

-- todos (created_at, updated_at)
UPDATE todos SET created_at = datetime(substr(created_at, 1, 19), '-8 hours') WHERE created_at IS NOT NULL;
UPDATE todos SET updated_at = datetime(substr(updated_at, 1, 19), '-8 hours') WHERE updated_at IS NOT NULL;

-- wishes
UPDATE wishes SET created_at = datetime(substr(created_at, 1, 19), '-8 hours') WHERE created_at IS NOT NULL;
UPDATE wishes SET updated_at = datetime(substr(updated_at, 1, 19), '-8 hours') WHERE updated_at IS NOT NULL;

-- posts
UPDATE posts SET created_at = datetime(substr(created_at, 1, 19), '-8 hours') WHERE created_at IS NOT NULL;
UPDATE posts SET updated_at = datetime(substr(updated_at, 1, 19), '-8 hours') WHERE updated_at IS NOT NULL;

-- topics
UPDATE topics SET created_at = datetime(substr(created_at, 1, 19), '-8 hours') WHERE created_at IS NOT NULL;
UPDATE topics SET updated_at = datetime(substr(updated_at, 1, 19), '-8 hours') WHERE updated_at IS NOT NULL;

-- comments
UPDATE comments SET created_at = datetime(substr(created_at, 1, 19), '-8 hours') WHERE created_at IS NOT NULL;
UPDATE comments SET updated_at = datetime(substr(updated_at, 1, 19), '-8 hours') WHERE updated_at IS NOT NULL;

-- votes (created_at, updated_at)
UPDATE votes SET created_at = datetime(substr(created_at, 1, 19), '-8 hours') WHERE created_at IS NOT NULL;
UPDATE votes SET updated_at = datetime(substr(updated_at, 1, 19), '-8 hours') WHERE updated_at IS NOT NULL;

-- vote_options
UPDATE vote_options SET created_at = datetime(substr(created_at, 1, 19), '-8 hours') WHERE created_at IS NOT NULL;

-- vote_records
UPDATE vote_records SET created_at = datetime(substr(created_at, 1, 19), '-8 hours') WHERE created_at IS NOT NULL;

-- likes
UPDATE likes SET created_at = datetime(substr(created_at, 1, 19), '-8 hours') WHERE created_at IS NOT NULL;

-- tags (no updated_at column)
UPDATE tags SET created_at = datetime(substr(created_at, 1, 19), '-8 hours') WHERE created_at IS NOT NULL;

-- todo_logs
UPDATE todo_logs SET created_at = datetime(substr(created_at, 1, 19), '-8 hours') WHERE created_at IS NOT NULL;

-- wish_votes
UPDATE wish_votes SET created_at = datetime(substr(created_at, 1, 19), '-8 hours') WHERE created_at IS NOT NULL;

-- +goose Down
-- Data migration, cannot be reversed