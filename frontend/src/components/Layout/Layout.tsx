import { Suspense } from 'react'
import { Outlet, useLocation } from 'react-router'
import { Box } from '@mui/material'

import { Fallback } from '@/components/Fallback/Fallback'
import { LayoutHeader } from './LayoutHeader'
import { PriceLayoutHeader } from './PriceLayoutHeader'

export const Layout = () => {
	const location = useLocation()
	const isPriceRoute = location.pathname.startsWith('/price')

	return (
		<Box minHeight='100vh' height='100vh' display='flex' flexDirection='column' pb={4}>
			{isPriceRoute ? <PriceLayoutHeader /> : <LayoutHeader />}
			<Suspense key={location.key} fallback={<Fallback />}>
				<Outlet />
			</Suspense>
		</Box>
	)
}

export default Layout
