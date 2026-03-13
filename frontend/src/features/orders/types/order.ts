import type { IPosition, IPositionCreate } from './positions'

export interface IOrder {
	id: string
	customer: string
	consumer: string
	manager: string
	bill: string
	date: string
	notes: string
	createdAt: string
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
