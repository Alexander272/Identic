import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IFullRole } from './types/role'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

export const rolesApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getRoles: builder.query<{ data: IFullRole[] }, null>({
			query: () => ({
				url: API.roles,
				method: 'GET',
			}),
			providesTags: [{ type: 'Roles', id: 'all' }],
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

export const { useGetRolesQuery } = rolesApiSlice
