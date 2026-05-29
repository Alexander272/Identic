export type Price = {
	id: string
	code: string
	current_name: string
	new_name: string | null
	price: number | null
	template: string | null
	note: string | null
	technique: string | null
	under_drawing: string | null
	matched_fields: string[] | null
}

export type SearchPriceRequest = {
	queries?: string[]
	codes?: string[]
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
	technique?: string | null
	under_drawing?: string | null
}

export type BatchSaveRequest = {
	positions: UpdatePrice[]
	delete_codes?: string[]
}

export type BatchSaveResponse = {
	status: string
}
