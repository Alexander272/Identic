import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { ISearch, ISearchResults } from './types/search'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

export const searchApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		findOrders: builder.query<{ data: ISearchResults[] }, ISearch>({
			query: data => ({
				url: API.search.base,
				method: 'POST',
				body: data,
			}),
			// providesTags: [{ type: 'Sections', id: 'ALL' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					console.log(error)
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
	}),
})

export const { useFindOrdersQuery, useLazyFindOrdersQuery } = searchApiSlice
