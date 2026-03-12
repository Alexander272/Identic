export interface IOrder {
	id: string
	customer: string
	consumer: string
	manager: string
	bill: string
	date: Date
	notes: string
	createdAt: Date
	positions: IPosition[]
}

export interface IPosition {
	id: string
	orderId: string
	rowNumber: number
	name: string
	quantity: number
	notes: string
	createdAt: Date
}
