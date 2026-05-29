import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { Price, SearchPriceRequest, ExportPriceRequest, BatchSaveRequest, BatchSaveResponse } from './types/types'
import { apiSlice } from '@/app/apiSlice'
import { API } from '@/app/api'

export const priceApiSlice = apiSlice.injectEndpoints({
	endpoints: builder => ({
		searchPrice: builder.mutation<{ data: Price[] }, SearchPriceRequest>({
			query: body => ({
				url: API.price.search,
				method: 'POST',
				body,
			}),
		}),
		exportPrices: builder.mutation<Blob, ExportPriceRequest>({
			query: body => ({
				url: API.price.export,
				method: 'POST',
				body,
				responseHandler: async response => {
					const blob = await response.blob()
					const disposition = response.headers.get('Content-Disposition') || ''
					const match = disposition.match(/filename="?(.+?)"?$/)
					const filename = match ? match[1] : 'книга цен.xlsx'
					const url = window.URL.createObjectURL(blob)
					const a = document.createElement('a')
					a.href = url
					a.download = filename
					a.click()
					window.URL.revokeObjectURL(url)
					return blob
				},
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

export const { useSearchPriceMutation, useExportPricesMutation, useImportPricesMutation, useBatchPriceSaveMutation } =
	priceApiSlice
