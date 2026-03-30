import { createBrowserRouter, type RouteObject } from 'react-router'

import { AppRoutes } from './routes'
import { Layout } from '@/components/Layout/Layout'
import { NotFound } from '@/pages/notFound/NotFoundLazy'
import { Home } from '@/pages/home/HomeLazy'
import { Order } from '@/pages/order/OrderLazy'
import { CreateOrder } from '@/pages/createOrder/CreateOrderLazy'
import { EditOrder } from '@/pages/editOrder/EditOrderLazy'
import { OrdersList } from '@/pages/ordersList/OrdersListLazy'
import { Search } from '@/pages/search/SearchLazy'

const config: RouteObject[] = [
	{
		element: <Layout />,
		errorElement: <NotFound />,
		children: [
			{
				path: AppRoutes.Home,
				element: <Home />,
			},
			{
				path: AppRoutes.Order,
				element: <Order />,
			},
			{
				path: AppRoutes.CreateOrder,
				element: <CreateOrder />,
			},
			{
				path: AppRoutes.EditOrder,
				element: <EditOrder />,
			},
			{
				path: AppRoutes.OrdersList,
				element: <OrdersList />,
			},
			{
				path: AppRoutes.Search,
				element: <Search />,
			},
		],
	},
]

export const router = createBrowserRouter(config)
