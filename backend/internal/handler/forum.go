package handler

import (
	"warmisle/internal/service"
)

type ForumHandler struct {
	svc *service.ForumService
}

func NewForumHandler() *ForumHandler {
	return &ForumHandler{svc: service.NewForumService()}
}
