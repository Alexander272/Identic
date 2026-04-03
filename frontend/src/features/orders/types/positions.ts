export interface IPosition {
	id: string
	orderId: string
	rowNumber: number
	name: string
	quantity: number
	notes: string
	createdAt: string
}

export interface IPositionCreate {
	rowNumber: number
	name: string | null
	quantity: number | null
	notes: string | null
}

export interface IPositionUpdate {
	id: string
	rowNumber: number
	name: string | null
	quantity: number | null
	notes: string | null
	status: 'CREATED' | 'UPDATED' | 'DELETED'
}
