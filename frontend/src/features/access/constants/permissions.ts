export const PermRules = Object.freeze({
	Orders: {
		Read: 'orders:read',
		Write: 'orders:write',
	},
	Users: {
		Read: 'user:read',
		Write: 'user:write',
	},
	SearchLog: { Read: 'search_log:read' },
	ActivityLog: { Read: 'activity_log:read' },
	Permissions: {
		Read: 'permissions:read',
		Write: 'permissions:write',
	},
})
