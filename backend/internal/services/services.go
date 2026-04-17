package services

import (
	"github.com/Alexander272/Identic/backend/internal/config"
	"github.com/Alexander272/Identic/backend/internal/events"
	"github.com/Alexander272/Identic/backend/internal/repository"
	"github.com/Alexander272/Identic/backend/pkg/auth"
	"github.com/Alexander272/Identic/backend/pkg/ws_hub"
)

type MessageBroadcaster interface {
	BroadcastMessage(topic string, data []byte)
}

type Services struct {
	Import
	Orders
	OrdersStream
	Positions
	Search
	SearchStream

	Permissions
	Roles
	Users
	Session
	AccessPolices

	AuditLogs
	Activity
	UserLogins
	Statistic
}

type Deps struct {
	Repo     *repository.Repository
	Keycloak *auth.KeycloakClient
	Conf     *config.Config
	Hub      *ws_hub.Hub
}

func NewServices(deps *Deps) *Services {
	transaction := NewTransactionManager(deps.Repo.Transaction)

	updatePolicyEvent := &events.PolicyEventManager{}

	search := NewSearchService(deps.Repo.Search, deps.Conf.Links.Orders, deps.Conf.Search.CacheTTL)
	searchLog := NewSearchLogService(deps.Repo.SearchLogs)
	searchStream := NewSearchStreamService(search, deps.Hub, searchLog)

	ordersStream := NewOrderStreamService(deps.Repo.OrdersEvents, deps.Hub)

	activity := NewActivityService(deps.Repo.Activity)

	positions := NewPositionsService(deps.Repo.Positions, transaction, activity)
	orders := NewOrdersService(deps.Repo.Orders, transaction, positions, search, activity)
	import_file := NewImportService(transaction, orders, positions)

	permissions := NewPermissionService(deps.Repo.Permissions, transaction, updatePolicyEvent)
	roleHierarchy := NewRoleHierarchyService(deps.Repo.RoleHierarchy)
	role := NewRoleService(&RoleDeps{
		Repo:        deps.Repo.Roles,
		Hierarchy:   roleHierarchy,
		Permissions: permissions,
		EventBus:    updatePolicyEvent,
	})
	user := NewUserService(&UsersDeps{
		Repo:      deps.Repo.Users,
		TxManager: transaction,
		Keycloak:  deps.Keycloak,
		Role:      role,
		EventBus:  updatePolicyEvent,
	})

	adapter := NewAdapter(&AdapterDeps{Permissions: permissions, Users: user, RoleHierarchy: roleHierarchy})
	policies := NewAccessPoliciesService(&PoliciesDeps{
		Conf:     deps.Conf.Casbin,
		Adapter:  adapter,
		EventBus: updatePolicyEvent,
	})

	userLogins := NewUserLoginService(deps.Repo.UserLogins, transaction)

	session := NewSessionService(deps.Keycloak, user, policies, userLogins)

	auditLogs := NewAuditLogService(deps.Repo.AuditLogs, transaction, updatePolicyEvent)
	stats := NewStatisticService(activity, searchLog)

	return &Services{
		Import:       import_file,
		Orders:       orders,
		OrdersStream: ordersStream,
		Positions:    positions,
		Search:       search,
		SearchStream: searchStream,

		Permissions: permissions,
		Roles:       role,
		Users:       user,
		Session:     session,

		AccessPolices: policies,

		AuditLogs:  auditLogs,
		Activity:   activity,
		UserLogins: userLogins,
		Statistic:  stats,
	}
}
