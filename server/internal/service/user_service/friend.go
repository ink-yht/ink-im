package user_service

import (
	"context"
	"ink-im-server/internal/domain/user_domain"
	"ink-im-server/internal/repository/user_repo"
	"ink-im-server/pkg/logger"
)

type FriendService interface {
	Info(ctx context.Context, uid uint) ([]user_domain.FriendsInfo, error)
}

type fiendService struct {
	repo user_repo.FriendRepository
	l    logger.Logger
}

func NewFriendService(repo user_repo.FriendRepository, l logger.Logger) FriendService {
	return &fiendService{
		repo: repo,
		l:    l,
	}
}

func (svc *fiendService) Info(ctx context.Context, uid uint) ([]user_domain.FriendsInfo, error) {
	return svc.repo.FindById(ctx, uid)
}
