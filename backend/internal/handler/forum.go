package handler

import (
	"warmisle/internal/service"
)

// ForumHandler handles forum-related HTTP requests including posts,
// comments, votes, likes, feeds, and topics.
type ForumHandler struct {
	svc *service.ForumService
}

// NewForumHandler creates a new ForumHandler.
func NewForumHandler() *ForumHandler {
	return &ForumHandler{svc: service.NewForumService()}
}
