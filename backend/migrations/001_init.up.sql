-- +goose Up
-- 家庭数字中心 V1 — 初始化数据库迁移
-- 创建所有基础表、索引和预置数据

-- ============================================================
-- 1. members 成员表
-- ============================================================
CREATE TABLE IF NOT EXISTS members (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    name TEXT NOT NULL DEFAULT '',
    avatar TEXT NOT NULL DEFAULT '👨',
    role TEXT NOT NULL DEFAULT 'member' CHECK(role IN ('admin', 'member')),
    status TEXT NOT NULL DEFAULT 'active' CHECK(status IN ('active', 'disabled')),
    last_login DATETIME,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at DATETIME
);
CREATE INDEX IF NOT EXISTS idx_members_username_deleted ON members(username, deleted_at);

-- ============================================================
-- 2. categories 分类表
-- ============================================================
CREATE TABLE IF NOT EXISTS categories (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    type TEXT NOT NULL CHECK(type IN ('income', 'expense')),
    name TEXT NOT NULL,
    icon TEXT NOT NULL DEFAULT '',
    sort_order INTEGER NOT NULL DEFAULT 0,
    preset INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at DATETIME
);

-- ============================================================
-- 3. tags 标签表
-- ============================================================
CREATE TABLE IF NOT EXISTS tags (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    preset INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at DATETIME
);

-- ============================================================
-- 4. ledgers 记账记录表
-- ============================================================
CREATE TABLE IF NOT EXISTS ledgers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    amount INTEGER NOT NULL,
    note TEXT DEFAULT '',
    category_id INTEGER NOT NULL REFERENCES categories(id),
    creator_id INTEGER NOT NULL REFERENCES members(id),
    occurred_at DATETIME NOT NULL,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at DATETIME
);
CREATE INDEX IF NOT EXISTS idx_ledgers_deleted_occurred ON ledgers(deleted_at, occurred_at);
CREATE INDEX IF NOT EXISTS idx_ledgers_creator_deleted ON ledgers(creator_id, deleted_at);
CREATE INDEX IF NOT EXISTS idx_ledgers_category_deleted ON ledgers(category_id, deleted_at);

-- ============================================================
-- 5. ledger_members 记账关联成员中间表
-- ============================================================
CREATE TABLE IF NOT EXISTS ledger_members (
    ledger_id INTEGER NOT NULL REFERENCES ledgers(id),
    member_id INTEGER NOT NULL REFERENCES members(id),
    PRIMARY KEY (ledger_id, member_id)
);
CREATE INDEX IF NOT EXISTS idx_ledger_members_member ON ledger_members(member_id);

-- ============================================================
-- 6. todos 待办事项表
-- ============================================================
CREATE TABLE IF NOT EXISTS todos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    priority TEXT NOT NULL DEFAULT 'normal' CHECK(priority IN ('normal', 'important', 'urgent')),
    status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'completed')),
    assignee_id INTEGER REFERENCES members(id),
    creator_id INTEGER NOT NULL REFERENCES members(id),
    due_date DATE,
    completed_at DATETIME,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at DATETIME
);
CREATE INDEX IF NOT EXISTS idx_todos_deleted_status_priority_due ON todos(deleted_at, status, priority, due_date);
CREATE INDEX IF NOT EXISTS idx_todos_assignee_deleted ON todos(assignee_id, deleted_at);

-- ============================================================
-- 7. todo_logs 待办变更日志表
-- ============================================================
CREATE TABLE IF NOT EXISTS todo_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    todo_id INTEGER NOT NULL REFERENCES todos(id),
    field_name TEXT NOT NULL,
    old_value TEXT DEFAULT '',
    new_value TEXT DEFAULT '',
    operator_id INTEGER NOT NULL REFERENCES members(id),
    created_at DATETIME NOT NULL DEFAULT (datetime('now'))
);

-- ============================================================
-- 8. wishes 愿望清单表
-- ============================================================
CREATE TABLE IF NOT EXISTS wishes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    category TEXT NOT NULL DEFAULT 'other' CHECK(category IN ('item', 'travel', 'experience', 'other')),
    amount INTEGER,
    priority TEXT NOT NULL DEFAULT 'normal' CHECK(priority IN ('normal', 'important', 'urgent')),
    type TEXT NOT NULL DEFAULT 'personal' CHECK(type IN ('personal', 'family')),
    status TEXT NOT NULL DEFAULT 'pending' CHECK(status IN ('pending', 'agreed', 'achieved', 'abandoned')),
    creator_id INTEGER NOT NULL REFERENCES members(id),
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at DATETIME
);

-- ============================================================
-- 9. wish_votes 愿望投票表
-- ============================================================
CREATE TABLE IF NOT EXISTS wish_votes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    wish_id INTEGER NOT NULL REFERENCES wishes(id),
    member_id INTEGER NOT NULL REFERENCES members(id),
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    UNIQUE(wish_id, member_id)
);
CREATE INDEX IF NOT EXISTS idx_wish_votes_wish ON wish_votes(wish_id);

-- ============================================================
-- 10. posts 动态表
-- ============================================================
CREATE TABLE IF NOT EXISTS posts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    creator_id INTEGER NOT NULL REFERENCES members(id),
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at DATETIME
);
CREATE INDEX IF NOT EXISTS idx_posts_deleted_created ON posts(deleted_at, created_at);

-- ============================================================
-- 11. topics 话题表
-- ============================================================
CREATE TABLE IF NOT EXISTS topics (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    content TEXT DEFAULT '',
    tag_id INTEGER REFERENCES tags(id),
    creator_id INTEGER NOT NULL REFERENCES members(id),
    is_pinned INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at DATETIME
);
CREATE INDEX IF NOT EXISTS idx_topics_deleted_created ON topics(deleted_at, created_at);
CREATE INDEX IF NOT EXISTS idx_topics_pinned_deleted_created ON topics(is_pinned, deleted_at, created_at);

-- ============================================================
-- 12. votes 投票表
-- ============================================================
CREATE TABLE IF NOT EXISTS votes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    creator_id INTEGER NOT NULL REFERENCES members(id),
    is_multi INTEGER NOT NULL DEFAULT 0,
    deadline DATE,
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at DATETIME
);

-- ============================================================
-- 13. vote_options 投票选项表
-- ============================================================
CREATE TABLE IF NOT EXISTS vote_options (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    vote_id INTEGER NOT NULL REFERENCES votes(id),
    content TEXT NOT NULL,
    sort_order INTEGER NOT NULL DEFAULT 0,
    created_at DATETIME NOT NULL DEFAULT (datetime('now'))
);

-- ============================================================
-- 14. vote_records 投票记录表
-- ============================================================
CREATE TABLE IF NOT EXISTS vote_records (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    vote_id INTEGER NOT NULL REFERENCES votes(id),
    option_id INTEGER NOT NULL REFERENCES vote_options(id),
    member_id INTEGER NOT NULL REFERENCES members(id),
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    UNIQUE(vote_id, member_id)
);
CREATE INDEX IF NOT EXISTS idx_vote_records_vote_member ON vote_records(vote_id, member_id);

-- ============================================================
-- 15. comments 评论表
-- ============================================================
CREATE TABLE IF NOT EXISTS comments (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    target_type TEXT NOT NULL CHECK(target_type IN ('post', 'topic', 'wish')),
    target_id INTEGER NOT NULL,
    parent_id INTEGER REFERENCES comments(id),
    content TEXT NOT NULL,
    creator_id INTEGER NOT NULL REFERENCES members(id),
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now')),
    deleted_at DATETIME
);
CREATE INDEX IF NOT EXISTS idx_comments_target_deleted ON comments(target_type, target_id, deleted_at);

-- ============================================================
-- 16. likes 点赞表
-- ============================================================
CREATE TABLE IF NOT EXISTS likes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    target_type TEXT NOT NULL CHECK(target_type IN ('post', 'topic', 'comment')),
    target_id INTEGER NOT NULL,
    member_id INTEGER NOT NULL REFERENCES members(id),
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    UNIQUE(target_type, target_id, member_id)
);
CREATE INDEX IF NOT EXISTS idx_likes_target ON likes(target_type, target_id);

-- ============================================================
-- 预置分类（支出 15 个 + 收入 5 个）
-- ============================================================
INSERT OR IGNORE INTO categories (type, name, icon, sort_order, preset, created_at, updated_at) VALUES
('expense', '餐饮', '🍱', 1, 1, datetime('now'), datetime('now')),
('expense', '交通', '🚗', 2, 1, datetime('now'), datetime('now')),
('expense', '购物', '🛒', 3, 1, datetime('now'), datetime('now')),
('expense', '居住', '🏠', 4, 1, datetime('now'), datetime('now')),
('expense', '通讯', '📱', 5, 1, datetime('now'), datetime('now')),
('expense', '医疗', '🏥', 6, 1, datetime('now'), datetime('now')),
('expense', '教育', '📚', 7, 1, datetime('now'), datetime('now')),
('expense', '娱乐', '🎮', 8, 1, datetime('now'), datetime('now')),
('expense', '亲子', '👶', 9, 1, datetime('now'), datetime('now')),
('expense', '人情', '🎁', 10, 1, datetime('now'), datetime('now')),
('expense', '宠物', '🐱', 11, 1, datetime('now'), datetime('now')),
('expense', '美容', '💄', 12, 1, datetime('now'), datetime('now')),
('expense', '运动', '⚽', 13, 1, datetime('now'), datetime('now')),
('expense', '保险', '🛡️', 14, 1, datetime('now'), datetime('now')),
('expense', '其他支出', '📦', 15, 1, datetime('now'), datetime('now')),
('income', '工资', '💰', 16, 1, datetime('now'), datetime('now')),
('income', '兼职', '💼', 17, 1, datetime('now'), datetime('now')),
('income', '理财', '📈', 18, 1, datetime('now'), datetime('now')),
('income', '红包', '🧧', 19, 1, datetime('now'), datetime('now')),
('income', '其他收入', '📦', 20, 1, datetime('now'), datetime('now'));

-- ============================================================
-- 预置标签（10 个）
-- ============================================================
INSERT OR IGNORE INTO tags (name, preset, created_at) VALUES
('家务', 1, datetime('now')),
('育儿', 1, datetime('now')),
('出行', 1, datetime('now')),
('饮食', 1, datetime('now')),
('健康', 1, datetime('now')),
('教育', 1, datetime('now')),
('财务', 1, datetime('now')),
('购物', 1, datetime('now')),
('装修', 1, datetime('now')),
('宠物', 1, datetime('now'));

-- +goose Down
-- 按依赖顺序删除表（先删子表，再删父表）
DROP TABLE IF EXISTS likes;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS vote_records;
DROP TABLE IF EXISTS vote_options;
DROP TABLE IF EXISTS votes;
DROP TABLE IF EXISTS topics;
DROP TABLE IF EXISTS posts;
DROP TABLE IF EXISTS wish_votes;
DROP TABLE IF EXISTS wishes;
DROP TABLE IF EXISTS todo_logs;
DROP TABLE IF EXISTS todos;
DROP TABLE IF EXISTS ledger_members;
DROP TABLE IF EXISTS ledgers;
DROP TABLE IF EXISTS tags;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS members;
