import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { SearchLogResponse, ActivityLogResponse } from './types'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

export const statisticsApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getSearchLogs: builder.query<SearchLogResponse, null>({
			query: () => ({
				url: API.statistics.search,
				method: 'GET',
			}),
			providesTags: [{ type: 'SearchLogs', id: 'All' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		getActivityLogs: builder.query<ActivityLogResponse, null>({
			query: () => ({
				url: API.statistics.activity,
				method: 'GET',
			}),
			providesTags: [{ type: 'ActivityLogs', id: 'All' }],
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

export const { useGetSearchLogsQuery, useGetActivityLogsQuery } = statisticsApiSlice