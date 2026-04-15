package postgres

var Tables = struct {
	Orders          string
	Positions       string
	Roles           string
	RoleHierarchy   string
	Permissions     string
	RolePermissions string
	Users           string
	UserLogins      string
	AuditLogs       string
	ActivityLogs    string
	SearchLogs      string
}{
	Orders:          "orders",
	Positions:       "positions",
	Roles:           "roles",
	RoleHierarchy:   "role_hierarchy",
	Permissions:     "permissions",
	RolePermissions: "role_permissions",
	Users:           "users",
	UserLogins:      "user_logins",
	AuditLogs:       "policy_audit_logs",
	ActivityLogs:    "activity_logs",
	SearchLogs:      "search_logs",
}
