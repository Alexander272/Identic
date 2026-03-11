export interface ISearchItem {
	name: string | null
	quantity: number | null
}

export interface ISearch {
	items: ISearchItem[]
	isFuzzy: boolean
}

export interface IOrderMatchResult {
	orderId: string
	customer: string
	consumer: string
	year: number
	link: string
	score: number // Общий процент совпадения (0-100)
	matchedCount: number // Сколько позиций совпало
	totalCount: number // Сколько позиций в запросе
}

export interface ISearchResults {
	year: number
	count: number
	orders: IOrderMatchResult[]
}
