import { Suspense } from 'react'
import { Stack } from '@mui/material'
import { Outlet } from 'react-router'

import { PageBox } from '@/components/PageBox/PageBox'
import { Fallback } from '@/components/Fallback/Fallback'
import Sidebar from './Sidebar'

export default function Permissions() {
	return (
		<PageBox>
			<Stack direction={'row'}>
				<Sidebar />

				<Suspense fallback={<Fallback />}>
					<Outlet />
				</Suspense>
			</Stack>
		</PageBox>
	)
}
