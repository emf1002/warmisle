-- +goose Up
-- 暖屿 V1 — 网盘备份模块数据表

-- ============================================================
-- 1. cloud_drive_configs 云盘配置表（单行配置）
-- ============================================================
CREATE TABLE IF NOT EXISTS cloud_drive_configs (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    provider         TEXT NOT NULL DEFAULT 'alipan',
    app_id           TEXT NOT NULL DEFAULT '',
    encrypted_secret TEXT NOT NULL DEFAULT '',
    redirect_uri     TEXT NOT NULL DEFAULT '',
    encrypted_token  TEXT NOT NULL DEFAULT '',
    token_expiry     DATETIME,
    status           TEXT NOT NULL DEFAULT 'unconfigured',
    backup_dir       TEXT NOT NULL DEFAULT '/warmisle-backups/',
    schedule_enabled INTEGER NOT NULL DEFAULT 0,
    schedule_time    TEXT NOT NULL DEFAULT '03:00',
    retention_days   INTEGER NOT NULL DEFAULT 30,
    created_at       DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at       DATETIME NOT NULL DEFAULT (datetime('now'))
);

-- ============================================================
-- 2. backup_records 备份记录表
-- ============================================================
CREATE TABLE IF NOT EXISTS backup_records (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    file_name      TEXT NOT NULL,
    cloud_file_id  TEXT NOT NULL DEFAULT '',
    file_size      INTEGER NOT NULL DEFAULT 0,
    backup_type    TEXT NOT NULL DEFAULT 'manual',
    upload_status  TEXT NOT NULL DEFAULT 'pending',
    integrity_ok   INTEGER NOT NULL DEFAULT 0,
    error_message  TEXT NOT NULL DEFAULT '',
    is_pre_restore INTEGER NOT NULL DEFAULT 0,
    created_at     DATETIME NOT NULL DEFAULT (datetime('now'))
);

-- ============================================================
-- 3. 初始插入一条空配置
-- ============================================================
INSERT INTO cloud_drive_configs (provider, status) VALUES ('alipan', 'unconfigured');
