package service

import (
	"errors"
	"sync"
	"time"

	"home-center/internal/pkg"
	"home-center/internal/repository"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrAccountLocked      = errors.New("account locked")
)

type lockInfo struct {
	failedCount int
	lockedUntil time.Time
}

type AuthService struct {
	repo     *repository.AuthRepo
	mu       sync.Mutex
	attempts map[string]*lockInfo
}

func NewAuthService() *AuthService {
	return &AuthService{
		repo:     &repository.AuthRepo{},
		attempts: make(map[string]*lockInfo),
	}
}

func (s *AuthService) isLocked(username string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	info, exists := s.attempts[username]
	if !exists {
		return false
	}
	if time.Now().Before(info.lockedUntil) {
		return true
	}
	// 锁定期已过，清除记录
	delete(s.attempts, username)
	return false
}

func (s *AuthService) recordFailed(username string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	info, exists := s.attempts[username]
	if !exists {
		info = &lockInfo{}
		s.attempts[username] = info
	}
	info.failedCount++
	if info.failedCount >= 5 {
		info.lockedUntil = time.Now().Add(15 * time.Minute)
	}
}

func (s *AuthService) clearAttempts(username string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.attempts, username)
}

func (s *AuthService) Login(username, password string) (string, error) {
	// 检查是否被锁定
	if s.isLocked(username) {
		return "", ErrAccountLocked
	}

	member, err := s.repo.FindByUsername(username)
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
	count, err := s.repo.Count()
	if err != nil {
		return false, err
	}
	return count == 0, nil
}
