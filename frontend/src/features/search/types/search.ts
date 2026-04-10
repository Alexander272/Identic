export interface ISearchItem {
	id: number
	name: string | null
	quantity: number | null
}

export interface ISearch {
	items: ISearchItem[]
	isFuzzy: boolean
	sessionId: string
}

export interface IOrderMatchResult {
	orderId: string
	customer: string
	consumer: string
	year: number
	date: string
	link: string
	score: number // Общий процент совпадения (0-100)
	matchedPos: number // Сколько позиций совпало
	matchedQuant: number // Сколько позиций + количество совпало
	totalCount: number // Сколько позиций в запросе
	// positionIds: string[]
	positions: IMatchPosition[]
}
export interface IMatchPosition {
	id: string
	reqId: string
	quantEqual: boolean
}

export interface ISearchResults {
	year: number
	count: number
	orders: IOrderMatchResult[]
}

export interface ISearchError {
	searchId: string
	message: string
}
