package services

import (
	"fmt"
	"log"

	"github.com/Alexander272/Identic/backend/internal/config"
	"github.com/Alexander272/Identic/backend/internal/events"
	"github.com/Alexander272/Identic/backend/pkg/logger"
	"github.com/casbin/casbin/v3"
)

type accessPolicesService struct {
	enforcer casbin.IEnforcer
	adapter  Adapter
	eventBus *events.PolicyEventManager
}

type PoliciesDeps struct {
	Conf     config.CasbinConfig
	Adapter  Adapter
	EventBus *events.PolicyEventManager
}

func NewAccessPoliciesService(deps *PoliciesDeps) *accessPolicesService {
	enforcer, err := casbin.NewEnforcer(deps.Conf.ModelPath, deps.Adapter)
	if err != nil {
		log.Fatalf("failed to initialize permission service. error: %s", err.Error())
	}

	if err = enforcer.LoadPolicy(); err != nil {
		log.Fatalf("failed to load policy from DB: %s", err.Error())
	}

	s := &accessPolicesService{
		enforcer: enforcer,
		adapter:  deps.Adapter,
	}

	go func() {
		updateChan := deps.EventBus.Subscribe()
		for range updateChan {
			logger.Info("Received policy update event, reloading...")
			s.enforcer.LoadPolicy()
		}
	}()

	return s
}

type AccessPolices interface {
	Enforce(sub, obj, act string) (bool, error)
	Reload() error
	GetPolicies(user string) (*Access, error)
}

func (s *accessPolicesService) Enforce(sub, obj, act string) (bool, error) {
	return s.enforcer.Enforce(sub, obj, act)
}

func (s *accessPolicesService) Reload() error {
	err := s.enforcer.LoadPolicy()
	if err != nil {
		return fmt.Errorf("failed to reload policies: %w", err)
	}
	return nil
}

type Access struct {
	Role  string
	Perms []string
}

func (s *accessPolicesService) GetPolicies(user string) (*Access, error) {
	allPermissions, err := s.enforcer.GetImplicitPermissionsForUser(user)
	if err != nil {
		// обработка ошибки
		return nil, fmt.Errorf("failed to get implicit permissions for user: %w", err)
	}

	var perms []string
	seen := make(map[string]bool)

	for _, p := range allPermissions {
		// В стандартном конфиге: p[1] - объект, p[2] - действие
		rule := fmt.Sprintf("%s:%s", p[1], p[2])

		// Убираем возможные дубликаты (если право пришло и от роли, и напрямую)
		if !seen[rule] {
			perms = append(perms, rule)
			seen[rule] = true
		}
	}

	role := ""
	if len(allPermissions) > 0 {
		role = allPermissions[len(allPermissions)-1][0]
	}

	result := &Access{
		Role:  role,
		Perms: perms,
	}
	return result, nil
}
