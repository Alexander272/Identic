import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IFlatOrder, IGetFlatOrders, IOrder, IOrderCreate, IOrderUpdate } from './types/order'
import type { IFilter } from './types/params'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'
import { wsService } from '@/app/services/socket'
import { buildFilterUrlParams } from './utils/buildFilters'

export const orderApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getOrders: builder.query<{ data: IOrder[] }, IFilter[]>({
			query: req => ({
				url: API.orders.base,
				method: 'GET',
				params: buildFilterUrlParams(req),
			}),
			providesTags: [{ type: 'Orders', id: 'ALL' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		getOrderById: builder.query<{ data: IOrder }, { id: string; searchId?: string }>({
			query: req => ({
				url: `${API.orders.base}/${req.id}`,
				method: 'GET',
				params: { search: req.searchId },
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
		getOrderInfo: builder.query<{ data: IOrder }, { id: string; searchId?: string }>({
			query: req => ({
				url: API.orders.info(req.id),
				method: 'GET',
				params: { search: req.searchId },
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
		getFlatOrder: builder.query<
			{ data: { orders: IFlatOrder[]; cursor: string; hasMore: boolean } },
			IGetFlatOrders
		>({
			query: req => ({
				url: API.orders.flat,
				method: 'GET',
				params: {
					// fields: req.fields.join(','),
					// value: req.value,
					// sort: `${req.sort.field}_${req.sort.order}`,
					cursor: req.cursor,
					limit: req.limit,
				},
			}),
			providesTags: [{ type: 'Orders', id: 'Flat' }],
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
		updateOrder: builder.mutation<{ id: string }, IOrderUpdate>({
			query: order => ({
				url: `${API.orders.base}/${order.id}`,
				method: 'PUT',
				body: order,
			}),
			invalidatesTags: ['Orders'],
		}),
	}),
})

export const {
	useGetOrdersQuery,
	useGetOrderByIdQuery,
	useGetOrdersByYearQuery,
	useLazyGetOrderInfoQuery,
	useGetUniqueDataQuery,
	useLazyGetUniqueDataQuery,
	useGetFlatOrderQuery,
	useLazyGetFlatOrderQuery,
	useCreateOrderMutation,
	useUpdateOrderMutation,
} = orderApiSlice
