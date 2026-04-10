package access

const (
	ResourcePerm  ResourceSlug = "permission"
	ResourceRole  ResourceSlug = "role"
	ResourceUser  ResourceSlug = "user"
	ResourceOrder ResourceSlug = "order"

	ResourceAudit    ResourceSlug = "audit_log"
	ResourceActivity ResourceSlug = "activity_log"
)

var Reg = NewRegistry(
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
	Resource{
		Slug:           ResourceUser,
		Name:           "Пользователи",
		Group:          "Администрирование",
		Description:    "Управление пользователями",
		AllowedActions: actions(All),
	},
	Resource{
		Slug:           ResourceOrder,
		Name:           "Заявки",
		Group:          "Операции",
		Description:    "Управление заявками",
		AllowedActions: actions(All),
	},

	Resource{
		Slug:           ResourceAudit,
		Name:           "Журнал изменений",
		Group:          "Логи",
		Description:    "История изменений прав доступа и разрешений",
		AllowedActions: actions(Read),
	},
	Resource{
		Slug:           ResourceActivity,
		Name:           "Журнал активности",
		Group:          "Логи",
		Description:    "Системный журнал действий пользователей",
		AllowedActions: actions(Read),
	},
)
