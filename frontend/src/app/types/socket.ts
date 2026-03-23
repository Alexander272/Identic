import type { IOrder, IOrderSocketMessage } from '@/features/orders/types/order'
import type { IOrderMatchResult, ISearchError, ISearchItem } from '@/features/search/types/search'

// Базовый интерфейс сообщения
export interface ISocketEnvelope<T = unknown> {
	type: string
	data: T
}

// Типы входящих действий (из твоего Go Hub)
export type TServerAction = 'INSERT' | 'UPDATE' | 'DELETE' | 'INSERT_MANY'

export type OrderAction = 'INSERT' | 'UPDATE' | 'DELETE' | 'INSERT_MANY'
export type SearchAction = 'SEARCH_STREAM' | 'SEARCH_RESULT' | 'SEARCH_RESULT_PART'

export interface BaseWSMessage {
	type: string
}

export interface OrderWSMessage extends BaseWSMessage {
	type: OrderAction
	data: IOrderSocketMessage
}

export interface SearchWSMessage extends BaseWSMessage {
	type: SearchAction
	data: IOrderMatchResult[]
}

export type WSMessage = OrderWSMessage | SearchWSMessage

export type WSEvent =
	| { type: 'SYSTEM_CONNECTED'; data: null }
	| { type: 'SYSTEM_DISCONNECTED'; data: null }
	| { type: 'SYSTEM_RECONNECTING'; data: null }
	| { type: 'ORDER_INSERTED'; data: IOrder }
	| { type: 'ORDER_UPDATED'; data: IOrder }
	| { type: 'ORDER_DELETED'; data: { id: string; createdAt: string } }
	| { type: 'ORDERS_BULK_INSERTED'; data: { years: number[] } }
	| { type: 'SEARCH_STREAM'; data: ISearchItem[] }
	| { type: 'SEARCH_RESULT'; data: IOrderMatchResult[] }
	| { type: 'SEARCH_RESULT_PART'; data: { items: IOrderMatchResult[]; isLast: boolean; total: number } }
	| { type: 'SEARCH_ERROR'; data: ISearchError }
	| { type: 'CANCEL_SEARCH'; data: null }
	| { type: 'SUBSCRIBE'; data: null }
	| { type: 'UNSUBSCRIBE'; data: null }

// Превращаем Union в Map для сервиса
export type WSEventMap = {
	[E in WSEvent as E['type']]: E['data']
}

export type Listener<T> = (data: T) => void
