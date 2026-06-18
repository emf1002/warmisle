package main

import (
	"fmt"
	"os"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cli <command> [args]")
		fmt.Println("Commands:")
		fmt.Println("  reset-password <username>  Reset user password to default")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "reset-password":
		if len(os.Args) < 3 {
			fmt.Println("Usage: cli reset-password <username>")
			os.Exit(1)
		}
		resetPassword(os.Args[2])
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}

func resetPassword(username string) {
	// 初始化数据库
	dbPath := getEnv("HC_DB_PATH", "./data/warmisle.db")
	if err := pkg.InitDatabase(dbPath); err != nil {
		fmt.Printf("failed to init database: %v\n", err)
		os.Exit(1)
	}

	// 查找用户
	var member model.Member
	if err := pkg.DB.Where("username = ?", username).First(&member).Error; err != nil {
		fmt.Printf("user not found: %s\n", username)
		os.Exit(1)
	}

	// 重置密码
	hash, err := pkg.HashPassword(pkg.DefaultPassword)
	if err != nil {
		fmt.Printf("failed to hash password: %v\n", err)
		os.Exit(1)
	}

	if err := pkg.DB.Model(&member).Update("password", hash).Error; err != nil {
		fmt.Printf("failed to reset password: %v\n", err)
		os.Exit(1)
	}

	// 同步清除登录失败记录
	if err := pkg.DB.Where("username = ?", username).Delete(&model.LoginFailure{}).Error; err != nil {
		fmt.Printf("warning: failed to clear login failures: %v\n", err)
	}

	fmt.Printf("Password for user '%s' has been reset to default: %s\n", username, pkg.DefaultPassword)
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
