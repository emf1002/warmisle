package pkg

import (
	"github.com/microcosm-cc/bluemonday"
)

// SanitizeHTML 使用 bluemonday UGCPolicy 过滤 HTML，去除危险标签/属性，防止 XSS
// UGCPolicy 允许安全的 HTML 标签（如 <b>, <i>, <a>, <code> 等），去除 <script>, onclick 等
func SanitizeHTML(input string) string {
	return bluemonday.UGCPolicy().Sanitize(input)
}
