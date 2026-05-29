import type { Price, UpdatePrice } from '../types/types'

export type GridRow = {
	id: string
	code: string | null
	current_name: string | null
	new_name: string | null
	price: number | null
	template: string | null
	note: string | null
	technique: string | null
	under_drawing: string | null
	status: 'ORIGINAL' | 'CREATED' | 'UPDATED' | 'DELETED'
}

export function positionToGridRow(p: Price): GridRow {
	return {
		id: p.id,
		code: p.code,
		current_name: p.current_name ?? '',
		new_name: p.new_name,
		price: p.price,
		template: p.template,
		note: p.note,
		technique: p.technique,
		under_drawing: p.under_drawing,
		status: 'ORIGINAL',
	}
}

export function gridRowToUpdate(r: GridRow): UpdatePrice {
	return {
		code: r.code || '',
		current_name: r.current_name || undefined,
		new_name: r.new_name,
		price: r.price,
		template: r.template,
		note: r.note,
		technique: r.technique,
		under_drawing: r.under_drawing,
	}
}
