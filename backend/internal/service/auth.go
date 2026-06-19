package service

import (
	"errors"
	"time"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountLocked      = errors.New("account locked")
)

type AuthService struct {
	authRepo     *repository.AuthRepo
	failureRepo  *repository.LoginFailureRepo
}

func NewAuthService() *AuthService {
	return &AuthService{
		authRepo:    &repository.AuthRepo{},
		failureRepo: &repository.LoginFailureRepo{},
	}
}

// isLocked 检查用户是否被锁定（从数据库读取）
func (s *AuthService) isLocked(username string) bool {
	lf, err := s.failureRepo.FindByUsername(username)
	if err != nil {
		// 记录不存在，说明未被锁定
		return false
	}
	// 检查 locked_until 是否有效
	if lf.LockedUntil == nil {
		return false
	}
	if time.Now().Before(*lf.LockedUntil) {
		return true
	}
	// 锁定期已过，删除记录
	_ = s.failureRepo.Delete(username)
	return false
}

// recordFailed 记录一次登录失败（原子 upsert，避免 TOCTOU 竞态）
func (s *AuthService) recordFailed(username string) {
	lf, err := s.failureRepo.FindByUsername(username)
	if err != nil {
		// 记录不存在，创建新记录
		lf = &model.LoginFailure{
			Username:    username,
			FailedCount: 1,
			LockedUntil: nil,
		}
		_ = s.failureRepo.Save(lf)
		return
	}

	// 原子递增失败次数
	lf.FailedCount++

	// 失败次数 >= 5，锁定 15 分钟
	if lf.FailedCount >= 5 {
		t := time.Now().Add(15 * time.Minute)
		lf.LockedUntil = &t
	}

	// 使用原子 upsert 保存（基于主键冲突自动更新）
	_ = s.failureRepo.Save(lf)
}

// clearAttempts 清除指定用户的登录失败记录
func (s *AuthService) clearAttempts(username string) {
	_ = s.failureRepo.Delete(username)
}

func (s *AuthService) Login(username, password string) (string, error) {
	// 检查是否被锁定
	if s.isLocked(username) {
		return "", ErrAccountLocked
	}

	member, err := s.authRepo.FindByUsername(username)
	if err != nil {
		s.recordFailed(username)
		return "", ErrInvalidCredentials
	}

	if member.Status == "disabled" {
		return "", ErrInvalidCredentials
	}

	if !pkg.CheckPassword(member.Password, password) {
		s.recordFailed(username)
		return "", ErrInvalidCredentials
	}

	// 登录成功，清除失败记录
	s.clearAttempts(username)

	token, err := pkg.GenerateToken(member.ID, member.Username, member.Role)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) InitCheck() (bool, error) {
	count, err := s.authRepo.Count()
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
