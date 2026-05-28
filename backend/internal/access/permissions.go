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
	ResourcePrice    ResourceSlug = "price"
)

var OrderOfResources = map[ResourceSlug]int{
	ResourceOrder:    1,
	ResourcePrice:    2,
	ResourceSearch:   3,
	ResourceLogins:   4,
	ResourceActivity: 5,
	ResourceAudit:    6,
	ResourceUser:     7,
	ResourceRole:     8,
	ResourcePerm:     9,
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
		Slug:           ResourcePrice,
		Name:           "Книга Цен",
		Group:          "Операции",
		Description:    "Управление книгой цен",
		AllowedActions: actions(Read, Write),
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
