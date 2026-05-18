package service

import (
	"errors"
	"regexp"
	"unicode/utf8"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrMemberNotFound       = errors.New("member not found")
	ErrUsernameTaken        = errors.New("username already exists")
	ErrCannotDeleteLastAdmin = errors.New("cannot delete last admin")
	ErrHasActivityRecords   = errors.New("member has activity records, must disable instead")
	ErrCannotDisableSelf    = errors.New("cannot disable yourself")
	ErrInvalidUsername      = errors.New("invalid username format")
	ErrInvalidPassword      = errors.New("invalid password length")
	ErrInvalidName          = errors.New("invalid name length")
	ErrUsernameRequired     = errors.New("username is required")
	ErrIncorrectPassword    = errors.New("incorrect password")
)

type MemberService struct {
	repo *repository.MemberRepo
}

func NewMemberService() *MemberService {
	return &MemberService{repo: &repository.MemberRepo{}}
}

var usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`)

func validateUsername(username string) error {
	if username == "" {
		return ErrUsernameRequired
	}
	if !usernameRegex.MatchString(username) {
		return ErrInvalidUsername
	}
	return nil
}

func validatePassword(password string) error {
	length := utf8.RuneCountInString(password)
	if length < 6 || length > 32 {
		return ErrInvalidPassword
	}
	return nil
}

func validateName(name string) error {
	length := utf8.RuneCountInString(name)
	if length < 1 || length > 20 {
		return ErrInvalidName
	}
	return nil
}

func (s *MemberService) List() ([]model.Member, error) {
	return s.repo.List()
}

func (s *MemberService) Create(username, password, name, avatar, role string) (*model.Member, error) {
	if err := validateUsername(username); err != nil {
		return nil, err
	}
	if err := validatePassword(password); err != nil {
		return nil, err
	}

	// Check uniqueness
	existing, err := s.repo.FindByUsername(username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrUsernameTaken
	}

	hash, err := pkg.HashPassword(password)
	if err != nil {
		return nil, err
	}

	if name == "" {
		name = username
	}
	if avatar == "" {
		avatar = "👨"
	}
	if role == "" {
		role = "member"
	}

	member := &model.Member{
		Username: username,
		Password: hash,
		Name:     name,
		Avatar:   avatar,
		Role:     role,
		Status:   "active",
	}
	if err := s.repo.Create(member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *MemberService) Update(id uint, name, avatar, role string) (*model.Member, error) {
	if name != "" {
		if err := validateName(name); err != nil {
			return nil, err
		}
	}

	member, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	// If changing role from admin to member, check it's not the last admin
	if member.Role == "admin" && role != "" && role != "admin" {
		count, err := s.repo.CountAdmins()
		if err != nil {
			return nil, err
		}
		if count <= 1 {
			return nil, ErrCannotDeleteLastAdmin
		}
	}

	if name != "" {
		member.Name = name
	}
	if avatar != "" {
		member.Avatar = avatar
	}
	if role != "" {
		member.Role = role
	}

	if err := s.repo.Update(member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *MemberService) Delete(id uint, currentMemberID uint) error {
	member, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMemberNotFound
		}
		return err
	}

	// Cannot delete last admin
	if member.Role == "admin" {
		count, err := s.repo.CountAdmins()
		if err != nil {
			return err
		}
		if count <= 1 {
			return ErrCannotDeleteLastAdmin
		}
	}

	// Check activity records
	activityCount, err := s.repo.CountActivityRecords(id)
	if err != nil {
		return err
	}
	if activityCount > 0 {
		return ErrHasActivityRecords
	}

	return s.repo.SoftDelete(id)
}

func (s *MemberService) Disable(id uint, currentMemberID uint) error {
	if id == currentMemberID {
		return ErrCannotDisableSelf
	}

	member, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMemberNotFound
		}
		return err
	}

	// Cannot disable last admin
	if member.Role == "admin" {
		count, err := s.repo.CountAdmins()
		if err != nil {
			return err
		}
		if count <= 1 {
			return ErrCannotDeleteLastAdmin
		}
	}

	member.Status = "disabled"
	return s.repo.Update(member)
}

func (s *MemberService) Enable(id uint) (*model.Member, error) {
	member, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	member.Status = "active"
	if err := s.repo.Update(member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *MemberService) ResetPassword(id uint) error {
	hash, err := pkg.HashPassword(pkg.DefaultPassword)
	if err != nil {
		return err
	}

	member, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMemberNotFound
		}
		return err
	}
	member.Password = hash
	return s.repo.Update(member)
}

func (s *MemberService) UpdateProfile(id uint, name, avatar string) (*model.Member, error) {
	if name != "" {
		if err := validateName(name); err != nil {
			return nil, err
		}
	}

	member, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}

	if name != "" {
		member.Name = name
	}
	if avatar != "" {
		member.Avatar = avatar
	}

	if err := s.repo.Update(member); err != nil {
		return nil, err
	}
	return member, nil
}

func (s *MemberService) GetProfile(id uint) (*model.Member, error) {
	member, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMemberNotFound
		}
		return nil, err
	}
	return member, nil
}

func (s *MemberService) ChangePassword(id uint, oldPwd, newPwd string) error {
	member, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMemberNotFound
		}
		return err
	}

	if !pkg.CheckPassword(member.Password, oldPwd) {
		return ErrIncorrectPassword
	}

	if err := validatePassword(newPwd); err != nil {
		return err
	}

	hash, err := pkg.HashPassword(newPwd)
	if err != nil {
		return err
	}

	member.Password = hash
	return s.repo.Update(member)
}
