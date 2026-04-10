export interface IUser {
	id: string
	name: string
	role: string
	permissions: string[]
	token: string
}

export interface IUserData {
	id: string
	ssoId: string
	username: string
	firstName: string
	lastName: string
	email: string
}
