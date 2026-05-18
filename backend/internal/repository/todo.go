package repository

import (
	"warmisle/internal/model"
	"warmisle/internal/pkg"
)

type TodoRepo struct{}

type TodoFilter struct {
	Status     string
	AssigneeID *uint
	Page       int
	PageSize   int
}

type TodoWithAssoc struct {
	model.Todo
	Assignee *model.Member `json:"assignee"`
	Creator  model.Member  `json:"creator"`
}

type TodoListResult struct {
	List     []TodoWithAssoc `json:"list"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

func (r *TodoRepo) List(filter TodoFilter) (*TodoListResult, error) {
	query := pkg.DB.Model(&model.Todo{}).
		Preload("Assignee").
		Preload("Creator")

	if filter.Status != "" {
		query = query.Where("todos.status = ?", filter.Status)
	}
	if filter.AssigneeID != nil {
		query = query.Where("todos.assignee_id = ?", *filter.AssigneeID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var todos []model.Todo
	err := query.
		Order("CASE todos.priority WHEN 'urgent' THEN 1 WHEN 'important' THEN 2 WHEN 'normal' THEN 3 END").
		Order("todos.due_date ASC").
		Order("todos.created_at DESC").
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		Find(&todos).Error
	if err != nil {
		return nil, err
	}

	list := make([]TodoWithAssoc, 0, len(todos))
	for _, t := range todos {
		list = append(list, TodoWithAssoc{
			Todo:     t,
			Assignee: t.Assignee,
			Creator:  t.Creator,
		})
	}

	return &TodoListResult{
		List:     list,
		Total:    total,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}, nil
}

func (r *TodoRepo) FindByID(id uint) (*TodoWithAssoc, error) {
	var todo model.Todo
	err := pkg.DB.
		Preload("Assignee").
		Preload("Creator").
		First(&todo, id).Error
	if err != nil {
		return nil, err
	}
	return &TodoWithAssoc{
		Todo:     todo,
		Assignee: todo.Assignee,
		Creator:  todo.Creator,
	}, nil
}

func (r *TodoRepo) Create(todo *model.Todo) error {
	return pkg.DB.Create(todo).Error
}

func (r *TodoRepo) Update(todo *model.Todo) error {
	return pkg.DB.Save(todo).Error
}

func (r *TodoRepo) Delete(id uint) error {
	return pkg.DB.Delete(&model.Todo{}, id).Error
}

func (r *TodoRepo) CreateLog(log *model.TodoLog) error {
	return pkg.DB.Create(log).Error
}
