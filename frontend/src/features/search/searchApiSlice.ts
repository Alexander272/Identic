import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IOrderMatchResult, ISearch } from './types/search'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'
import { wsService } from '@/app/services/socket'

export const searchApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		// findOrders: builder.query<{ data: IOrderMatchResult[] }, ISearch>({
		// 	query: data => ({
		// 		url: API.search.base,
		// 		method: 'POST',
		// 		body: data,
		// 	}),
		// 	// providesTags: [{ type: 'Sections', id: 'ALL' }],
		// 	onQueryStarted: async (_arg, api) => {
		// 		try {
		// 			await api.queryFulfilled
		// 		} catch (error) {
		// 			console.log(error)
		// 			const fetchError = (error as IBaseFetchError).error
		// 			toast.error(fetchError.data.message, { autoClose: false })
		// 		}
		// 	},
		// }),

		findOrders: builder.query<{ data: IOrderMatchResult[]; isProcessing: boolean; id: string }, ISearch>({
			query: data => ({
				url: API.search.stream,
				method: 'POST',
				body: data,
			}),
			async onCacheEntryAdded(_arg, { updateCachedData, cacheDataLoaded, cacheEntryRemoved }) {
				try {
					// 1. Ждем подтверждения от HTTP (получаем searchId)
					const { data } = await cacheDataLoaded
					const currentSearchId = data.id
					const topic = `SEARCH_RESULTS_${currentSearchId}`

					// 2. Подписываемся на топик
					wsService.send('SUBSCRIBE', { topic })
					let isFirstBatch = true

					const unsubscribe = wsService.subscribe('SEARCH_RESULT_PART', message => {
						updateCachedData(draft => {
							draft.isProcessing = true
							if (isFirstBatch) {
								// Первая пачка: очищаем старые результаты (от POST запроса или прошлого поиска)
								draft.data = message.items
								isFirstBatch = false
							} else {
								// Последующие пачки: просто дополняем
								draft.data.push(...message.items)
							}

							if (message.isLast) {
								draft.isProcessing = false
							}
						})
					})

					const unSubError = wsService.subscribe('SEARCH_ERROR', payload => {
						if (payload.searchId !== currentSearchId) return

						updateCachedData(draft => {
							draft.isProcessing = false
						})
						toast.error(`Ошибка поиска: ${payload.message}`)
					})

					await cacheEntryRemoved
					wsService.send('UNSUBSCRIBE', { topic })

					// unSubConnect()
					unsubscribe()
					unSubError()
				} catch (e) {
					console.error('Streaming error:', e)
				}
			},
			transformResponse: (response: { id: string }) => {
				return {
					id: response.id,
					data: [],
					isProcessing: true, // Сразу ставим в true, как только получили ID
				}
			},
			async onQueryStarted(arg, { dispatch, queryFulfilled }) {
				dispatch(
					searchApiSlice.util.updateQueryData('findOrders', arg, draft => {
						draft.data = []
						draft.isProcessing = true
					}),
				)

				try {
					await queryFulfilled
				} catch (error) {
					// Если HTTP-запрос упал, выключаем лоадер
					dispatch(
						searchApiSlice.util.updateQueryData('findOrders', arg, draft => {
							draft.isProcessing = false
						}),
					)
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError?.data?.message || 'Ошибка соединения', { autoClose: false })
				}
			},
		}),

		findOrdersOld: builder.query<{ data: IOrderMatchResult[]; isProcessing: boolean }, ISearch>({
			// query: data => ({
			// 	url: API.search.stream,
			// 	method: 'POST',
			// 	body: data,
			// }),
			queryFn: () => ({ data: { data: [], isProcessing: true } }),
			async onCacheEntryAdded(arg, { updateCachedData, cacheDataLoaded, cacheEntryRemoved }) {
				try {
					await cacheDataLoaded

					const startSearch = () => {
						wsService.send('SEARCH_STREAM', arg)
					}

					// 1. Первый запуск при монтировании
					startSearch()
					let isFirstBatch = true

					// const unSubConnect = wsService.subscribe('SYSTEM_CONNECTED', () => {
					// 	console.log('🔄 Re-sending search request after reconnect')
					// 	startSearch()
					// })

					const unsubscribe = wsService.subscribe('SEARCH_RESULT_PART', message => {
						updateCachedData(draft => {
							if (isFirstBatch) {
								// Первая пачка: очищаем старые результаты (от POST запроса или прошлого поиска)
								draft.data = message.items
								isFirstBatch = false
							} else {
								// Последующие пачки: просто дополняем
								draft.data.push(...message.items)
							}

							if (message.isLast) {
								draft.isProcessing = false
							}
						})
					})
					// const unsubscribe = wsService.subscribe('SEARCH_RESULT', message => {
					// 	updateCachedData(draft => {
					// 		if (draft && message) {
					// 			draft.data = message
					// 		}
					// 	})
					// })

					await cacheEntryRemoved
					// unSubConnect()
					unsubscribe()
				} catch (e) {
					console.error('Streaming error:', e)
				}
			},
			// onQueryStarted: async (_arg, api) => {
			// 	try {
			// 		await api.queryFulfilled
			// 	} catch (error) {
			// 		const fetchError = (error as IBaseFetchError).error
			// 		toast.error(fetchError.data.message, { autoClose: false })
			// 	}
			// },
			async onQueryStarted(arg, { dispatch, queryFulfilled }) {
				// Мгновенно переводим кэш в состояние загрузки для данных аргументов
				dispatch(
					searchApiSlice.util.updateQueryData('findOrders', arg, draft => {
						draft.data = []
						draft.isProcessing = true
					}),
				)

				try {
					await queryFulfilled
				} catch (error) {
					// Если HTTP-запрос упал, выключаем лоадер
					dispatch(
						searchApiSlice.util.updateQueryData('findOrders', arg, draft => {
							draft.isProcessing = false
						}),
					)
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
	}),
})

export const { useFindOrdersQuery, useLazyFindOrdersQuery } = searchApiSlice
