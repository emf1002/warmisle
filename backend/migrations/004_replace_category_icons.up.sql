-- +goose Up
-- Replace preset category emoji icons with Lucide icon names

-- Expense categories
UPDATE categories SET icon = 'UtensilsCrossed' WHERE name = '餐饮' AND preset = 1;
UPDATE categories SET icon = 'Car' WHERE name = '交通' AND preset = 1;
UPDATE categories SET icon = 'ShoppingBag' WHERE name = '购物' AND preset = 1;
UPDATE categories SET icon = 'Home' WHERE name = '居住' AND preset = 1;
UPDATE categories SET icon = 'Smartphone' WHERE name = '通讯' AND preset = 1;
UPDATE categories SET icon = 'HeartPulse' WHERE name = '医疗' AND preset = 1;
UPDATE categories SET icon = 'BookOpen' WHERE name = '教育' AND preset = 1;
UPDATE categories SET icon = 'Gamepad2' WHERE name = '娱乐' AND preset = 1;
UPDATE categories SET icon = 'Baby' WHERE name = '亲子' AND preset = 1;
UPDATE categories SET icon = 'Gift' WHERE name = '人情' AND preset = 1;
UPDATE categories SET icon = 'PawPrint' WHERE name = '宠物' AND preset = 1;
UPDATE categories SET icon = 'Sparkles' WHERE name = '美容' AND preset = 1;
UPDATE categories SET icon = 'Dumbbell' WHERE name = '运动' AND preset = 1;
UPDATE categories SET icon = 'ShieldCheck' WHERE name = '保险' AND preset = 1;
UPDATE categories SET icon = 'Package' WHERE name = '其他支出' AND preset = 1;

-- Income categories
UPDATE categories SET icon = 'Banknote' WHERE name = '工资' AND preset = 1;
UPDATE categories SET icon = 'Briefcase' WHERE name = '兼职' AND preset = 1;
UPDATE categories SET icon = 'TrendingUp' WHERE name = '理财' AND preset = 1;
UPDATE categories SET icon = 'Mail' WHERE name = '红包' AND preset = 1;
UPDATE categories SET icon = 'Package' WHERE name = '其他收入' AND preset = 1;

-- +goose Down
-- Restore original emoji icons

-- Expense categories
UPDATE categories SET icon = '🍱' WHERE name = '餐饮' AND preset = 1;
UPDATE categories SET icon = '🚗' WHERE name = '交通' AND preset = 1;
UPDATE categories SET icon = '🛒' WHERE name = '购物' AND preset = 1;
UPDATE categories SET icon = '🏠' WHERE name = '居住' AND preset = 1;
UPDATE categories SET icon = '📱' WHERE name = '通讯' AND preset = 1;
UPDATE categories SET icon = '🏥' WHERE name = '医疗' AND preset = 1;
UPDATE categories SET icon = '📚' WHERE name = '教育' AND preset = 1;
UPDATE categories SET icon = '🎮' WHERE name = '娱乐' AND preset = 1;
UPDATE categories SET icon = '👶' WHERE name = '亲子' AND preset = 1;
UPDATE categories SET icon = '🎁' WHERE name = '人情' AND preset = 1;
UPDATE categories SET icon = '🐱' WHERE name = '宠物' AND preset = 1;
UPDATE categories SET icon = '💄' WHERE name = '美容' AND preset = 1;
UPDATE categories SET icon = '⚽' WHERE name = '运动' AND preset = 1;
UPDATE categories SET icon = '🛡️' WHERE name = '保险' AND preset = 1;
UPDATE categories SET icon = '📦' WHERE name = '其他支出' AND preset = 1;

-- Income categories
UPDATE categories SET icon = '💰' WHERE name = '工资' AND preset = 1;
UPDATE categories SET icon = '💼' WHERE name = '兼职' AND preset = 1;
UPDATE categories SET icon = '📈' WHERE name = '理财' AND preset = 1;
UPDATE categories SET icon = '🧧' WHERE name = '红包' AND preset = 1;
UPDATE categories SET icon = '📦' WHERE name = '其他收入' AND preset = 1;
