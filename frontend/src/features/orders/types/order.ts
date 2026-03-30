import type { IPosition, IPositionCreate } from './positions'

export interface IOrderSocketMessage extends IOrder {
	years: number[]
}

export interface IOrder {
	id: string
	customer: string
	consumer: string
	manager: string
	bill: string
	date: string
	notes: string
	year: number
	createdAt: string
	positionCount: number
	positions: IPosition[]
}

export interface IOrderCreate {
	customer: string
	consumer: string
	manager: string
	bill: string
	date: string
	notes: string
	positions: IPositionCreate[]
}

export interface IFlatOrder {
	id: string
	customer: string
	consumer: string
	manager: string
	bill: string
	date: string
	notes: string
	rowNumber: number
	name: string
	quantity: number
	positionNotes: string
	createdAt: string
}

export interface IGetFlatOrders {
	search?: {
		fields: string[]
		value: string
	}
	sort?: {
		field: string
		order: 'ASC' | 'DESC'
	}
	cursor: string | null
	limit?: number
}
