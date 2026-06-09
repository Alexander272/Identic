import type { IUserShort } from '@/features/user/types/user'

export interface PriceSearchLog {
	id: string
	queries: string[]
	codes: string[]
	fields?: string[]
	actor: IUserShort
	resultsCount: number
	durationMs: number
	createdAt: string
}

export interface PriceSearchLogResponse {
	total: number
	data: PriceSearchLog[]
}

export interface PriceSearchLogRequest {
	actorId?: string
	startDate?: string
	endDate?: string
	limit?: number
	offset?: number
}
