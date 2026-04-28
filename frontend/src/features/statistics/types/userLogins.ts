import type { IUserLogin, IUserShort } from '@/features/user/types/user'

export interface IUserLoginWithUser extends IUserLogin {
	user: IUserShort
}

export interface IUserLoginsResponse {
	total: number
	data: IUserLoginWithUser[]
}

export interface IUserLoginsRequest {
	startDate?: string
	endDate?: string
	limit?: number
	offset?: number
}
