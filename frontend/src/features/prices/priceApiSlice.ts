import { toast } from 'react-toastify'

import type { IBaseFetchError, IFetchError } from '@/app/types/error'
import { saveAs } from '@/features/prices/utils/saveAs'
import type {
	Price,
	PaginatedPriceResponse,
	SearchPriceRequest,
	ExportPriceRequest,
	BatchSaveRequest,
	BatchSaveResponse,
} from './types/types'
import { apiSlice } from '@/app/apiSlice'
import { API } from '@/app/api'

export const priceApiSlice = apiSlice.injectEndpoints({
	endpoints: builder => ({
		getPrices: builder.query<PaginatedPriceResponse, { page: number; limit: number }>({
			query: ({ page, limit }) => `${API.price.base}?page=${page}&limit=${limit}`,
			providesTags: ['Prices'],
		}),
		searchPrice: builder.mutation<{ data: Price[]; total?: number }, SearchPriceRequest>({
			query: body => ({
				url: API.price.search,
				method: 'POST',
				body,
			}),
		}),
		searchAllPrices: builder.mutation<{ data: Price[]; total?: number }, SearchPriceRequest>({
			query: body => ({
				url: API.price.searchAll,
				method: 'POST',
				body,
			}),
		}),
		exportPrices: builder.mutation<null, ExportPriceRequest>({
			queryFn: async (params, _api, _options, baseQuery) => {
				const filename = `Книга цен.xlsx`
				const result = await baseQuery({
					url: API.price.export,
					method: 'POST',
					body: params,
					cache: 'no-cache',
					responseHandler: async response => {
						if (!response.ok) return response.json()
						return response.blob()
					},
				})
				if (result.error) {
					const fetchError = result.error as IFetchError
					if (fetchError.status !== 401) {
						toast.error(fetchError.data.message, { autoClose: false })
					}
					return { data: null }
				}
				if (result.data instanceof Blob) {
					saveAs(result.data, filename)
				}
				return { data: null }
			},
			invalidatesTags: ['Prices'],
		}),
		importPrices: builder.mutation<{ imported: number }, FormData>({
			query: body => ({
				url: API.price.import,
				method: 'POST',
				body,
				formData: true,
			}),
			invalidatesTags: ['Prices'],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		batchPriceSave: builder.mutation<BatchSaveResponse, BatchSaveRequest>({
			query: body => ({
				url: API.price.batch,
				method: 'PUT',
				body,
			}),
			invalidatesTags: ['Prices'],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
	}),
})

export const {
	useGetPricesQuery,
	useSearchPriceMutation,
	useSearchAllPricesMutation,
	useExportPricesMutation,
	useImportPricesMutation,
	useBatchPriceSaveMutation,
} = priceApiSlice
