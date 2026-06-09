import type { Price, UpdatePrice } from '../types/types'

export type GridRow = {
	id: string
	code: string | null
	currentName: string | null
	newName: string | null
	price: number | null
	template: string | null
	note: string | null
	needSiburApproval: string | null
	status: 'ORIGINAL' | 'CREATED' | 'UPDATED' | 'DELETED'
}

export function positionToGridRow(p: Price): GridRow {
	return {
		id: p.id,
		code: p.code,
		currentName: p.currentName ?? '',
		newName: p.newName,
		price: p.price,
		template: p.template,
		note: p.note,
		needSiburApproval: p.needSiburApproval,
		status: 'ORIGINAL',
	}
}

export function gridRowToUpdate(r: GridRow): UpdatePrice {
	return {
		code: r.code || '',
		currentName: r.currentName || undefined,
		newName: r.newName,
		price: r.price,
		template: r.template,
		note: r.note,
		needSiburApproval: r.needSiburApproval,
	}
}
