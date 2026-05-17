package testutil

import (
	"home-center/internal/model"
	"home-center/internal/pkg"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB creates an in-memory SQLite database, auto-migrates all models, and sets pkg.DB
func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect test database: " + err.Error())
	}

	// Auto-migrate all models
	err = db.AutoMigrate(
		&model.Member{},
		&model.Category{},
		&model.Ledger{},
		&model.LedgerMember{},
		&model.Todo{},
		&model.TodoLog{},
		&model.Wish{},
		&model.WishVote{},
		&model.Post{},
		&model.Topic{},
		&model.Vote{},
		&model.VoteOption{},
		&model.VoteRecord{},
		&model.Comment{},
		&model.Like{},
		&model.Tag{},
	)
	if err != nil {
		panic("failed to migrate test database: " + err.Error())
	}

	pkg.DB = db
	return db
}

// TeardownTestDB clears the global DB reference
func TeardownTestDB() {
	pkg.DB = nil
}

// SeedMembers creates test members and returns them with IDs populated
func SeedMembers(db *gorm.DB, members []model.Member) []model.Member {
	for i := range members {
		db.Create(&members[i])
	}
	return members
}

// SeedCategories creates test categories
func SeedCategories(db *gorm.DB, categories []model.Category) []model.Category {
	for i := range categories {
		db.Create(&categories[i])
	}
	return categories
}

// CreateTestMember quickly creates a single test member
func CreateTestMember(db *gorm.DB, username, name, role string) model.Member {
	m := model.Member{
		Username: username,
		Password: "$2a$10$dummyhashedpassword1234567890abcdef",
		Name:     name,
		Avatar:   "\U0001F468",
		Role:     role,
		Status:   "active",
	}
	db.Create(&m)
	return m
}

// CreateTestCategory quickly creates a single test category
func CreateTestCategory(db *gorm.DB, ctype, name, icon string, sortOrder int) model.Category {
	c := model.Category{
		Type:      ctype,
		Name:      name,
		Icon:      icon,
		SortOrder: sortOrder,
		Preset:    false,
	}
	db.Create(&c)
	return c
}
