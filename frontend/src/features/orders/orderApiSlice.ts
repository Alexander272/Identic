import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IOrder, IOrderCreate } from './types/order'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'
import { wsService } from '@/app/services/socket'

export const orderApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getOrderById: builder.query<{ data: IOrder }, string>({
			query: id => ({
				url: `${API.orders.base}/${id}`,
				method: 'GET',
			}),
			providesTags: (_res, _err, arg) => [{ type: 'Orders', id: arg }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		getOrderInfo: builder.query<{ data: IOrder }, { id: string; positions?: string[] }>({
			query: req => ({
				url: API.orders.info(req.id),
				method: 'GET',
				params: { positions: req.positions?.join(',') },
			}),
			providesTags: (_res, _err, arg) => [{ type: 'Orders', id: arg.id }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		// getOrdersByYear: builder.query<{ data: IOrder[] }, string>({
		// 	query: year => ({
		// 		url: `${API.orders.byYear(year)}`,
		// 		method: 'GET',
		// 	}),
		// 	providesTags: (_res, _err, arg) => [{ type: 'Orders', id: arg }],
		// 	onQueryStarted: async (_arg, api) => {
		// 		try {
		// 			await api.queryFulfilled
		// 		} catch (error) {
		// 			const fetchError = (error as IBaseFetchError).error
		// 			toast.error(fetchError.data.message, { autoClose: false })
		// 		}
		// 	},
		// }),
		// getOrdersByYear: builder.query<{ data: IOrder[] }, string>({
		// 	query: year => ({
		// 		url: API.orders.byYear(year),
		// 		method: 'GET',
		// 	}),
		// 	providesTags: (_res, _err, arg) => [{ type: 'Orders', id: arg }],
		// 	async onCacheEntryAdded(year, { updateCachedData, cacheDataLoaded, cacheEntryRemoved, dispatch }) {
		// 		// Инициализируем переменные сразу, чтобы избежать "used before assignment"
		// 		let ws: WebSocket | null = null
		// 		let timeoutId: ReturnType<typeof setTimeout> | null = null
		// 		let isNamespaceClosed = false

		// 		const connect = () => {
		// 			if (isNamespaceClosed) return

		// 			const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
		// 			// Явно типизируем сокет
		// 			const socket = new WebSocket(`${protocol}//${window.location.host}/api/ws`)
		// 			ws = socket

		// 			socket.onmessage = (event: MessageEvent) => {
		// 				const message: IOrderSocketMessage = JSON.parse(event.data)
		// 				const orderYear = new Date(message.createdAt).getFullYear().toString()

		// 				if (message.action === 'INSERT_MANY') {
		// 					// Проверяем, есть ли текущий год вкладки в списке измененных годов
		// 					const isCurrentYearAffected = message.years?.some((y: number) => y.toString() === year)
		// 					if (isCurrentYearAffected) {
		// 						updateCachedData(() => {
		// 							dispatch(apiSlice.util.invalidateTags([{ type: 'Orders', id: year }]))
		// 						})
		// 					}
		// 					return
		// 				}

		// 				if (orderYear === year) {
		// 					updateCachedData(draft => {
		// 						if (message.action === 'DELETE') {
		// 							draft.data = draft.data.filter(o => o.id !== message.id)
		// 						} else if (message.action === 'UPDATE') {
		// 							const index = draft.data.findIndex(o => o.id === message.id)
		// 							if (index !== -1) {
		// 								draft.data[index] = { ...draft.data[index], ...message }
		// 							}
		// 						} else if (message.action === 'INSERT') {
		// 							draft.data.unshift(message)
		// 						}

		// 						// const index = draft.data.findIndex(o => o.id === message.id)
		// 						// if (index !== -1) {
		// 						// 	draft.data[index] = { ...draft.data[index], ...message }
		// 						// } else if (message.action === 'INSERT') {
		// 						// 	draft.data.unshift(message)
		// 						// }
		// 					})
		// 				}
		// 			}

		// 			socket.onclose = e => {
		// 				ws = null
		// 				// Если закрытие не инициировано нами, пробуем переподключиться
		// 				if (!isNamespaceClosed && e.code !== 1000) {
		// 					timeoutId = setTimeout(connect, 5000)
		// 				}
		// 			}

		// 			socket.onerror = () => {
		// 				socket.close()
		// 			}
		// 		}

		// 		try {
		// 			await cacheDataLoaded
		// 			connect()
		// 		} catch {
		// 			// Ошибка загрузки данных
		// 		}

		// 		// Ожидаем удаления данных из кэша (уход со страницы)
		// 		await cacheEntryRemoved

		// 		// Ставим флаг, чтобы реконнект больше не срабатывал
		// 		isNamespaceClosed = true

		// 		if (timeoutId) {
		// 			clearTimeout(timeoutId)
		// 		}

		// 		const socketToClose = ws as WebSocket | null
		// 		if (socketToClose) {
		// 			socketToClose.close(1000, 'Closed by RTK Query')
		// 		}
		// 	},
		// 	onQueryStarted: async (_arg, api) => {
		// 		try {
		// 			await api.queryFulfilled
		// 		} catch (error) {
		// 			const fetchError = error as IBaseFetchError
		// 			const message = fetchError?.error?.data?.message || 'Ошибка загрузки'
		// 			toast.error(message, { autoClose: false })
		// 		}
		// 	},
		// }),
		getOrdersByYear: builder.query<{ data: IOrder[] }, string>({
			query: year => ({
				url: API.orders.byYear(year),
				method: 'GET',
			}),
			providesTags: (_res, _err, arg) => [{ type: 'Orders', id: arg }],
			async onCacheEntryAdded(year, { updateCachedData, cacheDataLoaded, cacheEntryRemoved, dispatch }) {
				try {
					await cacheDataLoaded

					const unsubs: Array<() => void> = []

					// Подписываемся на топик
					wsService.send('SUBSCRIBE', { topic: 'orders' })

					// 1. Слушаем новые заказы
					unsubs.push(
						wsService.subscribe('ORDER_INSERTED', order => {
							console.log('insert order', order)

							if (order.year.toString() === year) {
								updateCachedData(draft => {
									draft.data.unshift(order)
								})
							}
						}),
					)
					// 2. Слушаем массовые вставки
					unsubs.push(
						wsService.subscribe('ORDERS_BULK_INSERTED', payload => {
							if (payload.years.some(y => y.toString() === year)) {
								dispatch(apiSlice.util.invalidateTags([{ type: 'Orders', id: year }]))
							}
						}),
					)
					// 3. Слушаем удаления
					unsubs.push(
						wsService.subscribe('ORDER_DELETED', order => {
							if (order.year.toString() === year) {
								updateCachedData(draft => {
									draft.data = draft.data.filter(o => o.id !== order.id)
								})
							}
						}),
					)
					// 4. Слушаем обновления
					unsubs.push(
						wsService.subscribe('ORDER_UPDATED', order => {
							if (order.year.toString() === year) {
								updateCachedData(draft => {
									const index = draft.data.findIndex(o => o.id === order.id)
									if (index !== -1) {
										draft.data[index] = { ...draft.data[index], ...order }
									}
								})
							}
						}),
					)

					await cacheEntryRemoved
					// Отписываемся
					wsService.send('UNSUBSCRIBE', { topic: 'orders' })
					unsubs.forEach(u => u())
				} catch (e) {
					console.error('Socket sync error:', e)
				}
			},
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = error as IBaseFetchError
					const message = fetchError?.error?.data?.message || 'Ошибка загрузки'
					toast.error(message, { autoClose: false })
				}
			},
		}),
		getUniqueData: builder.query<{ data: string[] }, { field: string; sort?: 'ASC' | 'DESC' }>({
			query: req => ({
				url: API.orders.unique(req.field),
				method: 'GET',
				params: { sort: req.sort },
			}),
			providesTags: [{ type: 'Orders', id: 'Unique' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),

		createOrder: builder.mutation<{ id: string }, IOrderCreate>({
			query: order => ({
				url: `${API.orders.base}`,
				method: 'POST',
				body: order,
			}),
			invalidatesTags: ['Orders'],
		}),
	}),
})

export const {
	useGetOrderByIdQuery,
	useGetOrdersByYearQuery,
	useLazyGetOrderInfoQuery,
	useGetUniqueDataQuery,
	useLazyGetUniqueDataQuery,
	useCreateOrderMutation,
} = orderApiSlice
