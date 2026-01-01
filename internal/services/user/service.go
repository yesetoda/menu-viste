package user

import (
	"context"

	"menuvista/internal/services/email"

	"menuvista/internal/storage/persistence"
)

type UserService struct {
	queries      *persistence.Queries
	emailService *email.Service
	// smsService   *sms.Service
}

func NewUserService(queries *persistence.Queries, emailService *email.Service) *UserService {
	return &UserService{
		queries:      queries,
		emailService: emailService,
		// smsService:   smsService,
	}
}

func (s *UserService) UpdateUser(ctx context.Context, updateUser persistence.UpdateUserParams) (persistence.User, error) {
	return s.queries.UpdateUser(ctx, updateUser)
}
