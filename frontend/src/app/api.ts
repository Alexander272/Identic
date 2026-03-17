export const API = {
	auth: {
		signIn: `auth/sign-in` as const,
		refresh: `auth/refresh` as const,
		signOut: `auth/sign-out` as const,
	},
	search: {
		base: `search` as const,
	},
	orders: {
		base: `orders` as const,
		byYear: (year: string) => `orders/by-year/${year}` as const,
		unique: (field: string) => `orders/unique/${field}` as const,
	},
}
