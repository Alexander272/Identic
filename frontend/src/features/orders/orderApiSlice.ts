import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IOrder, IOrderCreate, IOrderSocketMessage } from './types/order'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

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
		getOrdersByYear: builder.query<{ data: IOrder[] }, string>({
			query: year => ({
				url: API.orders.byYear(year),
				method: 'GET',
			}),
			providesTags: (_res, _err, arg) => [{ type: 'Orders', id: arg }],
			async onCacheEntryAdded(year, { updateCachedData, cacheDataLoaded, cacheEntryRemoved }) {
				// Инициализируем переменные сразу, чтобы избежать "used before assignment"
				let ws: WebSocket | null = null
				let timeoutId: ReturnType<typeof setTimeout> | null = null
				let isNamespaceClosed = false

				const connect = () => {
					if (isNamespaceClosed) return

					const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
					// Явно типизируем сокет
					const socket = new WebSocket(`${protocol}//${window.location.host}/api/ws`)
					ws = socket

					socket.onmessage = (event: MessageEvent) => {
						const message: IOrderSocketMessage = JSON.parse(event.data)
						const orderYear = new Date(message.createdAt).getFullYear().toString()

						if (orderYear === year) {
							updateCachedData(draft => {
								const index = draft.data.findIndex(o => o.id === message.id)

								if (index !== -1) {
									draft.data[index] = { ...draft.data[index], ...message }
								} else if (message.action === 'INSERT') {
									draft.data.unshift(message)
								}
							})
						}
					}

					socket.onclose = e => {
						ws = null
						// Если закрытие не инициировано нами, пробуем переподключиться
						if (!isNamespaceClosed && e.code !== 1000) {
							timeoutId = setTimeout(connect, 5000)
						}
					}

					socket.onerror = () => {
						socket.close()
					}
				}

				try {
					await cacheDataLoaded
					connect()
				} catch {
					// Ошибка загрузки данных
				}

				// Ожидаем удаления данных из кэша (уход со страницы)
				await cacheEntryRemoved

				// Ставим флаг, чтобы реконнект больше не срабатывал
				isNamespaceClosed = true

				if (timeoutId) {
					clearTimeout(timeoutId)
				}

				const socketToClose = ws as WebSocket | null
				if (socketToClose) {
					socketToClose.close(1000, 'Closed by RTK Query')
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
	useGetUniqueDataQuery,
	useLazyGetUniqueDataQuery,
	useCreateOrderMutation,
} = orderApiSlice
