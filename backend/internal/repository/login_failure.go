package repository

import (
	"warmisle/internal/model"
	"warmisle/internal/pkg"
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

// Save 插入或更新记录（Upsert）
func (r *LoginFailureRepo) Save(lf *model.LoginFailure) error {
	var existing model.LoginFailure
	err := pkg.DB.Where("username = ?", lf.Username).First(&existing).Error
	if err != nil {
		// 不存在，创建新记录
		return pkg.DB.Create(lf).Error
	}
	// 存在，更新字段
	existing.FailedCount = lf.FailedCount
	existing.LockedUntil = lf.LockedUntil
	return pkg.DB.Save(&existing).Error
}

// Delete 删除指定用户的记录
func (r *LoginFailureRepo) Delete(username string) error {
	return pkg.DB.Where("username = ?", username).Delete(&model.LoginFailure{}).Error
}
