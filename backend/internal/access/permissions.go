package access

const (
	ResourcePerm  ResourceSlug = "permission"
	ResourceRole  ResourceSlug = "role"
	ResourceUser  ResourceSlug = "user"
	ResourceOrder ResourceSlug = "order"

	ResourceLogins   ResourceSlug = "logins"
	ResourceAudit    ResourceSlug = "audit_log"
	ResourceActivity ResourceSlug = "activity_log"
	ResourceSearch   ResourceSlug = "search_log"
)

var OrderOfResources = map[ResourceSlug]int{
	ResourceOrder:    1,
	ResourceSearch:   2,
	ResourceLogins:   3,
	ResourceActivity: 4,
	ResourceAudit:    5,
	ResourceUser:     6,
	ResourceRole:     7,
	ResourcePerm:     8,
}

// TODO возможно стоит сделать какую-нибудь сортировку
var Reg = NewRegistry(
	Resource{
		Slug:           ResourceOrder,
		Name:           "Заявки",
		Group:          "Операции",
		Description:    "Управление заявками",
		AllowedActions: actions(All),
	},

	Resource{
		Slug:           ResourceActivity,
		Name:           "Журнал активности",
		Group:          "Логи",
		Description:    "Системный журнал действий пользователей",
		AllowedActions: actions(Read),
	},
	Resource{
		Slug:           ResourceAudit,
		Name:           "Журнал изменений",
		Group:          "Логи",
		Description:    "История изменений прав доступа и разрешений",
		AllowedActions: actions(Read),
	},
	Resource{
		Slug:           ResourceSearch,
		Name:           "Журнал поисков",
		Group:          "Логи",
		Description:    "История поисков пользователей",
		AllowedActions: actions(Read),
	},
	Resource{
		Slug:           ResourceLogins,
		Name:           "Логи входа",
		Group:          "Логи",
		Description:    "История входов пользователей",
		AllowedActions: actions(Read),
	},

	Resource{
		Slug:           ResourceUser,
		Name:           "Пользователи",
		Group:          "Администрирование",
		Description:    "Управление пользователями",
		AllowedActions: actions(Read, Write),
	},
	Resource{
		Slug:           ResourceRole,
		Name:           "Роли",
		Group:          "Администрирование",
		Description:    "Управление ролями пользователей",
		AllowedActions: actions(All),
	},
	Resource{
		Slug:           ResourcePerm,
		Name:           "Права",
		Group:          "Администрирование",
		Description:    "Действия, которые доступны пользователю",
		AllowedActions: actions(All),
	},
)
