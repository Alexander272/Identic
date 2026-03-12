import { createBrowserRouter, type RouteObject } from 'react-router'

import { AppRoutes } from './routes'
import { Layout } from '@/components/Layout/Layout'
import { NotFound } from '@/pages/notFound/NotFoundLazy'
import { Home } from '@/pages/home/HomeLazy'
import { Order } from '@/pages/order/OrderLazy'

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
		],
	},
]

export const router = createBrowserRouter(config)
