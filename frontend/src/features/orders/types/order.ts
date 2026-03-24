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
