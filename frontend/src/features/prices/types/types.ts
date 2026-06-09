export type PaginatedPriceResponse = {
	data: Price[]
	total: number
}

export type Price = {
	id: string
	code: string
	currentName: string
	newName: string | null
	price: number | null
	template: string | null
	note: string | null
	needSiburApproval: string | null
	matchedFields: string[] | null
}

export type SearchPriceRequest = {
	queries?: string[]
	codes?: string[]
	fields?: string[]
	page?: number
	limit?: number
}

export type ExportPriceRequest = {
	queries?: string[]
	codes?: string[]
	fields?: string[]
	columns?: string[]
}

export type UpdatePrice = {
	code: string
	currentName?: string
	newName?: string | null
	price?: number | null
	template?: string | null
	note?: string | null
	needSiburApproval?: string | null
}

export type BatchSaveRequest = {
	prices: UpdatePrice[]
	deleteCodes?: string[]
}

export type BatchSaveResponse = {
	status: string
}
