import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IResource } from './types/resource'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

export const permsApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getResources: builder.query<{ data: IResource[] }, null>({
			query: () => ({
				url: API.permissions.resources,
				method: 'GET',
			}),
			providesTags: [{ type: 'Perms', id: 'All' }],
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

export const { useGetResourcesQuery } = permsApiSlice
