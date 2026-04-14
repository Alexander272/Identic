export const AppRoutes = Object.freeze({
	Home: '/' as const,
	Auth: '/auth' as const,
	Order: '/orders/:id' as const,
	OrdersByYear: '/orders/by-year/:year' as const,
	CreateOrder: '/orders/new' as const,
	EditOrder: '/orders/edit/:id' as const,
	OrdersList: '/orders/list' as const,
	Search: '/search' as const,
	Accesses: '/accesses' as const,
	UserAccess: '/accesses/user' as const,
	RoleAccess: '/accesses/roles' as const,
	Permissions: '/accesses/permissions' as const,
})
