export interface IRole {
	id: string
	slug: string
	name: string
	description: string
	level: number
	isActive: boolean
	isSystem: boolean
	isEditable: boolean
	createdAt: Date
	updatedAt: Date
}

export interface IRoleWithStats extends IRole {
	inherited: string[]
	perms: IPermsCount
	userCount: number
}

export interface IPermsCount {
	own: number
	inherited: number
	total: number
}

export interface IFullRole {
	id: string
	name: string
	description: string
	level: number
	extends: string[]
	isShow: boolean
}
