-- +goose Up
-- 移除记账记录关联成员表
DROP TABLE IF EXISTS ledger_members;
