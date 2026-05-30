package service

import (
	"errors"
	"fmt"
	"time"

	"warmisle/internal/model"
	"warmisle/internal/repository"
)

var (
	ErrTodoNotFound         = errors.New("todo not found")
	ErrTodoPermissionDenied = errors.New("permission denied")
	ErrTodoAlreadyAssigned  = errors.New("todo already assigned")
	ErrTodoTitleRequired    = errors.New("title is required")
	ErrTodoInvalidPriority  = errors.New("invalid priority value")
	ErrTodoAssigneeNotFound = errors.New("assignee not found")
)

type TodoService struct {
	repo       *repository.TodoRepo
	memberRepo *repository.MemberRepo
}

func NewTodoService() *TodoService {
	return &TodoService{
		repo:       &repository.TodoRepo{},
		memberRepo: &repository.MemberRepo{},
	}
}

func (s *TodoService) List(filter repository.TodoFilter) (*repository.TodoListResult, error) {
	return s.repo.List(filter)
}

func (s *TodoService) FindByID(id uint) (*repository.TodoWithAssoc, error) {
	result, err := s.repo.FindByID(id)
	if err != nil {
		return nil, wrapNotFound(err, ErrTodoNotFound)
	}
	return result, nil
}

func (s *TodoService) Create(title, description, priority string, assigneeID *uint, dueDate *time.Time, creatorID uint) (*repository.TodoWithAssoc, error) {
	if title == "" {
		return nil, ErrTodoTitleRequired
	}

	if priority == "" {
		priority = "normal"
	}
	if !validPriority(priority) {
		return nil, ErrTodoInvalidPriority
	}

	if assigneeID != nil {
		_, err := s.memberRepo.FindByID(*assigneeID)
		if err != nil {
			return nil, wrapNotFound(err, ErrTodoAssigneeNotFound)
		}
	}

	todo := &model.Todo{
		Title:       title,
		Description: description,
		Priority:    priority,
		Status:      "pending",
		AssigneeID:  assigneeID,
		CreatorID:   creatorID,
		DueDate:     model.FromTimePtr(dueDate),
	}

	if err := s.repo.Create(todo); err != nil {
		return nil, err
	}

	if assigneeID != nil {
		s.repo.CreateLog(&model.TodoLog{
			TodoID:     todo.ID,
			FieldName:  "assignee",
			OldValue:   "",
			NewValue:   fmt.Sprintf("%d", *assigneeID),
			OperatorID: creatorID,
		})
	}

	return s.repo.FindByID(todo.ID)
}

func (s *TodoService) Update(id uint, title, description, priority *string, assigneeID *uint, dueDate *time.Time, currentMemberID uint, currentRole string) (*repository.TodoWithAssoc, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, wrapNotFound(err, ErrTodoNotFound)
	}

	// Permission: creator, assignee, or admin
	if existing.CreatorID != currentMemberID && currentRole != "admin" {
		if existing.AssigneeID == nil || *existing.AssigneeID != currentMemberID {
			return nil, ErrTodoPermissionDenied
		}
	}

	if title != nil {
		if *title == "" {
			return nil, ErrTodoTitleRequired
		}
		if *title != existing.Title {
			s.repo.CreateLog(&model.TodoLog{
				TodoID: id, FieldName: "title",
				OldValue: existing.Title, NewValue: *title, OperatorID: currentMemberID,
			})
			existing.Title = *title
		}
	}

	if description != nil {
		if *description != existing.Description {
			s.repo.CreateLog(&model.TodoLog{
				TodoID: id, FieldName: "description",
				OldValue: existing.Description, NewValue: *description, OperatorID: currentMemberID,
			})
			existing.Description = *description
		}
	}

	if priority != nil {
		if !validPriority(*priority) {
			return nil, ErrTodoInvalidPriority
		}
		if *priority != existing.Priority {
			s.repo.CreateLog(&model.TodoLog{
				TodoID: id, FieldName: "priority",
				OldValue: existing.Priority, NewValue: *priority, OperatorID: currentMemberID,
			})
			existing.Priority = *priority
		}
	}

	if dueDate != nil {
		oldStr := ""
		newStr := ""
		if existing.DueDate != nil {
			oldStr = existing.DueDate.ToTime().Format("2006-01-02")
		}
		if dueDate != nil {
			newStr = dueDate.Format("2006-01-02")
		}
		if oldStr != newStr {
			s.repo.CreateLog(&model.TodoLog{
				TodoID: id, FieldName: "due_date",
				OldValue: oldStr, NewValue: newStr, OperatorID: currentMemberID,
			})
			existing.DueDate = model.FromTimePtr(dueDate)
		}
	}

	if assigneeID != nil {
		if *assigneeID != 0 {
			_, err := s.memberRepo.FindByID(*assigneeID)
			if err != nil {
				return nil, wrapNotFound(err, ErrTodoAssigneeNotFound)
			}
		}

		oldStr := ""
		newStr := ""
		if existing.AssigneeID != nil {
			oldStr = fmt.Sprintf("%d", *existing.AssigneeID)
		}
		if *assigneeID != 0 {
			newStr = fmt.Sprintf("%d", *assigneeID)
		}
		if oldStr != newStr {
			s.repo.CreateLog(&model.TodoLog{
				TodoID: id, FieldName: "assignee",
				OldValue: oldStr, NewValue: newStr, OperatorID: currentMemberID,
			})
			if *assigneeID != 0 {
				existing.AssigneeID = assigneeID
			} else {
				existing.AssigneeID = nil
			}
		}
	}

	if err := s.repo.Update(&existing.Todo); err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *TodoService) Delete(id uint, currentMemberID uint, currentRole string) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return wrapNotFound(err, ErrTodoNotFound)
	}

	if existing.CreatorID != currentMemberID && currentRole != "admin" {
		return ErrTodoPermissionDenied
	}

	return s.repo.Delete(id)
}

func (s *TodoService) Toggle(id uint, currentMemberID uint, currentRole string) (*repository.TodoWithAssoc, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, wrapNotFound(err, ErrTodoNotFound)
	}

	if existing.CreatorID != currentMemberID && currentRole != "admin" {
		if existing.AssigneeID == nil || *existing.AssigneeID != currentMemberID {
			return nil, ErrTodoPermissionDenied
		}
	}

	now := model.FromTime(time.Now())
	if existing.Status == "pending" {
		existing.Status = "completed"
		existing.CompletedAt = &now
	} else {
		existing.Status = "pending"
		existing.CompletedAt = nil
	}

	if err := s.repo.Update(&existing.Todo); err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *TodoService) Claim(id uint, currentMemberID uint) (*repository.TodoWithAssoc, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, wrapNotFound(err, ErrTodoNotFound)
	}

	if existing.AssigneeID != nil {
		return nil, ErrTodoAlreadyAssigned
	}

	existing.AssigneeID = &currentMemberID

	if err := s.repo.Update(&existing.Todo); err != nil {
		return nil, err
	}

	s.repo.CreateLog(&model.TodoLog{
		TodoID:     id,
		FieldName:  "assignee",
		OldValue:   "",
		NewValue:   fmt.Sprintf("%d", currentMemberID),
		OperatorID: currentMemberID,
	})

	return s.repo.FindByID(id)
}
