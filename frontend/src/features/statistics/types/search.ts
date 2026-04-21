import type { IUserShort } from '@/features/user/types/user'

export const SearchType = {
	Exact: 'exact',
	Fuzzy: 'fuzzy',
} as const

export type SearchType = (typeof SearchType)[keyof typeof SearchType]

export interface SearchLog {
	id: string
	searchId: string
	actor: IUserShort
	actorId: string
	actorName: string
	searchType: SearchType
	query: JSON
	durationMs: number
	resultsCount: number
	itemsCount: number
	createdAt: string
}

export interface SearchLogResponse {
	total: number
	data: SearchLog[]
}

export interface SearchLogRequest {
	actorId?: string
	startDate?: string
	endDate?: string
	limit?: number
	offset?: number
}
