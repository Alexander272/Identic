import { toast } from 'react-toastify'

import type { IBaseFetchError } from '@/app/types/error'
import type { IOrder, IOrderCreate } from './types/order'
import { API } from '@/app/api'
import { apiSlice } from '@/app/apiSlice'

export const orderApiSlice = apiSlice.injectEndpoints({
	overrideExisting: false,
	endpoints: builder => ({
		getOrderById: builder.query<{ data: IOrder }, string>({
			query: id => ({
				url: `${API.orders.base}/${id}`,
				method: 'GET',
			}),
			providesTags: (_res, _err, arg) => [{ type: 'Orders', id: arg }],
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
	}),
})

export const { useGetOrderByIdQuery, useCreateOrderMutation } = orderApiSlice
