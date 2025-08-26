package initialization

import (
	"context"

	"github.com/ducklawrence05/go-test-backend-api/internal/constants/rolecache"
	"github.com/ducklawrence05/go-test-backend-api/internal/usecase/role"
	"github.com/ducklawrence05/go-test-backend-api/pkg/logger"
	"go.uber.org/zap"
)

func InitRolesCache(roleManager role.RoleManager, logger logger.Interface) {
	roles, err := roleManager.GetAll(context.Background())
	if err != nil {
		logger.Fatal("Roles initialization failed", zap.Error(err))
	}

	rolecache.InitCache(roles)
}
