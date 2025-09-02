package rolecache

import (
	"sync"

	"github.com/ducklawrence05/go-test-backend-api/internal/entities"
)

var (
	roleCache = map[string]entities.Role{}
	cacheMux  sync.RWMutex
)

func NewCache(roles []entities.Role) {
	tmp := map[string]entities.Role{}
	for _, r := range roles {
		tmp[r.Name] = r
	}

	cacheMux.Lock()
	defer cacheMux.Unlock()
	roleCache = tmp
}

func Add(role entities.Role) {
	cacheMux.Lock()
	defer cacheMux.Unlock()
	roleCache[role.Name] = role
}

func Get(roleName string) (entities.Role, bool) {
	cacheMux.RLock()
	defer cacheMux.RUnlock()
	r, ok := roleCache[roleName]
	return r, ok
}
