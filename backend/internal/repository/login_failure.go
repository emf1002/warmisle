package repository

import (
	"warmisle/internal/model"
	"warmisle/internal/pkg"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// LoginFailureRepo 登录失败记录仓库
type LoginFailureRepo struct{}

// FindByUsername 按用户名查找记录
func (r *LoginFailureRepo) FindByUsername(username string) (*model.LoginFailure, error) {
	var lf model.LoginFailure
	err := pkg.DB.Where("username = ?", username).First(&lf).Error
	if err != nil {
		return nil, err
	}
	return &lf, nil
}

// Save 原子 upsert：INSERT OR REPLACE，消除 TOCTOU 竞态条件
func (r *LoginFailureRepo) Save(lf *model.LoginFailure) error {
	return pkg.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "username"}},       // 主键冲突
		DoUpdates: clause.AssignmentColumns([]string{"failed_count", "locked_until", "updated_at"}), // 仅更新关键字段
	}).Create(lf).Error
}

// Delete 删除指定用户的记录
func (r *LoginFailureRepo) Delete(username string) error {
	return pkg.DB.Where("username = ?", username).Delete(&model.LoginFailure{}).Error
}

// UpsertLocked 原子操作：增加失败计数并可选锁定
// 返回更新后的记录
func (r *LoginFailureRepo) UpsertLocked(lf *model.LoginFailure) error {
	result := pkg.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "username"}},
		DoUpdates: clause.Assignments(map[string]interface{}{
			"failed_count": gorm.Expr("failed_count + 1"),
			"locked_until": lf.LockedUntil,
			"updated_at":   gorm.Expr("CURRENT_TIMESTAMP"),
		}),
	}).Create(lf)
	return result.Error
}
