export interface IUser {
	id: string
	name: string
	role: string
	permissions: string[]
	token: string
}

export interface IUserShort {
	id: string
	ssoId: string
	firstName: string
	lastName: string
	email: string
}

export interface IUserData {
	id: string
	ssoId: string
	username: string
	firstName: string
	lastName: string
	email: string
	roleId: string
	role: string
	isActive: boolean
	createdAt: string
	lastVisit: string
}

export interface IUserDataDTO {
	id: string
	ssoId: string
	roleId: string
	username: string
	firstName: string
	lastName: string
	email: string
	isActive: boolean
}
