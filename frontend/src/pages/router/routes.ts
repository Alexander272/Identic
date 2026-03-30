export const AppRoutes = Object.freeze({
	Home: '/' as const,
	// Auth: '/auth' as const,
	Order: '/orders/:id' as const,
	OrdersByYear: '/orders/by-year/:year' as const,
	CreateOrder: '/orders/new' as const,
	OrdersList: '/orders/list' as const,
	Search: '/search' as const,
})
