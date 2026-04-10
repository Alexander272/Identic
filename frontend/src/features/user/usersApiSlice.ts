import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IUserData } from './types/user'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

export const usersApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getAllUsers: builder.query<{ data: IUserData[] }, null>({
			query: () => ({
				url: API.users.base,
				method: 'GET',
			}),
			providesTags: [{ type: 'Users', id: 'all' }],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),
		getUserByAccess: builder.query<{ data: IUserData[] }, null>({
			query: () => ({
				url: `${API.users.access}`,
				method: 'GET',
			}),
			providesTags: [
				{ type: 'Users', id: 'all' },
				{ type: 'Users', id: 'access' },
			],
			onQueryStarted: async (_arg, api) => {
				try {
					await api.queryFulfilled
				} catch (error) {
					const fetchError = (error as IBaseFetchError).error
					toast.error(fetchError.data.message, { autoClose: false })
				}
			},
		}),

		syncUsers: builder.mutation<null, null>({
			query: () => ({
				url: API.users.sync,
				method: 'POST',
			}),
			invalidatesTags: [{ type: 'Users', id: 'all' }],
		}),
	}),
})

export const { useGetAllUsersQuery, useGetUserByAccessQuery, useSyncUsersMutation } = usersApiSlice
