package handler

import (
	"net/http"
	"os"
	"warmisle/internal/pkg"

	"github.com/gin-gonic/gin"
)

// TestReset 清空所有业务数据并重建种子数据
// 仅在 HC_TEST_MODE=true 时可用
func TestReset(c *gin.Context) {
	if os.Getenv("HC_TEST_MODE") != "true" {
		pkg.Error(c, http.StatusNotFound, 404, "not found")
		return
	}

	db := pkg.DB

	// 按外键依赖顺序删除业务数据
	tables := []string{
		"likes",
		"vote_records",
		"vote_options",
		"votes",
		"comments",
		"posts",
		"topics",
		"wish_votes",
		"wishes",
		"todo_logs",
		"todos",
		"ledger_members",
		"ledgers",
		"tags",
		"categories",
		"members",
	}

	for _, table := range tables {
		db.Exec("DELETE FROM " + table)
	}

	// 重建预置分类（支出 15 + 收入 5）
	db.Exec(`INSERT INTO categories (type, name, icon, sort_order, preset, created_at, updated_at) VALUES
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
		('income', '其他收入', '📦', 20, 1, datetime('now'), datetime('now'))
	`)

	// 重建预置标签（10 个）
	db.Exec(`INSERT INTO tags (name, preset, created_at) VALUES
		('家务', 1, datetime('now')),
		('育儿', 1, datetime('now')),
		('出行', 1, datetime('now')),
		('饮食', 1, datetime('now')),
		('健康', 1, datetime('now')),
		('教育', 1, datetime('now')),
		('财务', 1, datetime('now')),
		('购物', 1, datetime('now')),
		('装修', 1, datetime('now')),
		('宠物', 1, datetime('now'))
	`)

	pkg.Success(c, nil)
}
