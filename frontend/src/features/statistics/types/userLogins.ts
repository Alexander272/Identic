import type { IUserShort } from '@/features/user/types/user'

export interface IUserLogin {
	id: string
	userId: string
	loginAt: string
	ipAddress: string | null
	userAgent: string | null
	metadata: JSON
	lastActivityAt: string
}

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
