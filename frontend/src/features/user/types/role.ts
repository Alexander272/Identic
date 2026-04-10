export interface IRole {
	id: string
	name: string
}

export interface IFullRole {
	id: string
	name: string
	description: string
	level: number
	extends: string[]
	isShow: boolean
}
