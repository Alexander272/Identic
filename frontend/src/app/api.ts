export const API = {
	auth: {
		signIn: `auth/sign-in` as const,
		refresh: `auth/refresh` as const,
		signOut: `auth/sign-out` as const,
	},
	search: {
		base: `search` as const,
		stream: `search/stream` as const,
	},
	orders: {
		base: `orders` as const,
		info: (id: string) => `orders/info/${id}` as const,
		byYear: (year: string) => `orders/by-year/${year}` as const,
		unique: (field: string) => `orders/unique/${field}` as const,
		flat: `orders/flat` as const,
	},
	users: {
		base: '/users' as const,
		sync: '/users/sync' as const,
		access: '/users/access' as const,
	},
	roles: '/roles' as const,
	permissions: {
		base: '/permissions' as const,
		resources: '/permissions/resources' as const,
	},
}
