export type PaginatedPriceResponse = {
	data: Price[]
	total: number
}

export type Price = {
	id: string
	code: string
	current_name: string
	new_name: string | null
	price: number | null
	template: string | null
	note: string | null
	need_sibur_approval: string | null
	matched_fields: string[] | null
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
	columns?: string[]
}

export type UpdatePrice = {
	code: string
	current_name?: string
	new_name?: string | null
	price?: number | null
	template?: string | null
	note?: string | null
	need_sibur_approval?: string | null
}

export type BatchSaveRequest = {
	positions: UpdatePrice[]
	delete_codes?: string[]
}

export type BatchSaveResponse = {
	status: string
}
