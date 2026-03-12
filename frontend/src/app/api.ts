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
	},
}
