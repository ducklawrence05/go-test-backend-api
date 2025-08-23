package mapper

import (
	"github.com/ducklawrence05/go-test-backend-api/internal/delivery/payload"
	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
)

func ToUserInfoResponse(user *entities.User) payload.UserInfoResponse {
	return payload.UserInfoResponse{
		ID:        user.ID,
		UserName:  user.UserName,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		CreatedAt: user.CreatedAt,
		Role: payload.RoleInfoResponse{
			Name:        user.Role.Name,
			Description: user.Role.Description,
		},
	}
}
