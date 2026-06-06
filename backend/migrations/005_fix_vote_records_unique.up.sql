-- +goose Up
-- Change unique constraint from (vote_id, member_id) to (vote_id, member_id, option_id)
-- to support multi-select polls where one member can vote on multiple options.
DROP INDEX IF EXISTS idx_vote_records_vote_member;
CREATE UNIQUE INDEX IF NOT EXISTS idx_vote_records_vote_member_option ON vote_records(vote_id, member_id, option_id);
