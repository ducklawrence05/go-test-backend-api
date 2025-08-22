package router

import "github.com/ducklawrence05/go-test-backend-api/internal/router/user"

type RouterGroup struct {
	User user.UserRouterGroup
}

var RouterGroupApp = new(RouterGroup)
