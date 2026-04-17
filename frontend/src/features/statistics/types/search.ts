export const SearchType = {
	Exact: 'exact',
	Fuzzy: 'fuzzy',
} as const

export type SearchType = (typeof SearchType)[keyof typeof SearchType]

export interface SearchLog {
	id: string
	searchId: string
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
