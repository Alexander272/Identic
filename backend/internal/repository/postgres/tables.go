package postgres

var Tables = struct {
	Orders          string
	Positions       string
	Roles           string
	RoleHierarchy   string
	Permissions     string
	RolePermissions string
	Users           string
}{
	Orders:          "orders",
	Positions:       "positions",
	Roles:           "roles",
	RoleHierarchy:   "role_hierarchy",
	Permissions:     "permissions",
	RolePermissions: "role_permissions",
	Users:           "users",
}
