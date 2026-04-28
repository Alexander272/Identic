import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { SearchLogResponse, ActivityLogResponse, SearchLogRequest } from './types'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'
import type { IUserLoginsRequest, IUserLoginsResponse } from './types/userLogins'

export const statisticsApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getSearchLogs: builder.query<SearchLogResponse, SearchLogRequest | undefined>({
			query: params => ({
				url: API.statistics.search,
				method: 'GET',
				params,
			}),
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		getActivityLogs: builder.query<ActivityLogResponse, SearchLogRequest | undefined>({
			query: params => ({
				url: API.statistics.activity,
				method: 'GET',
				params,
			}),
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		getLastUserLogins: builder.query<IUserLoginsResponse, IUserLoginsRequest>({
			query: params => ({
				url: API.statistics.logins,
				method: 'GET',
				params,
			}),
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
	useGetSearchLogsQuery,
	useGetActivityLogsQuery,
	useGetLastUserLoginsQuery,
	useLazyGetActivityLogsQuery,
	useLazyGetSearchLogsQuery,
	useLazyGetLastUserLoginsQuery,
} = statisticsApiSlice
