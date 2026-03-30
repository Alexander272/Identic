import { Box, Breadcrumbs } from '@mui/material'

import { PageBox } from '@/components/PageBox/PageBox'
import { FlatOrders } from '@/features/orders/components/FlatOrders/FlatOrders'
import { Breadcrumb } from '@/components/Breadcrumb/Breadcrumb'
import { AppRoutes } from '../router/routes'

export default function OrdersList() {
	return (
		<PageBox>
			<Box
				borderRadius={3}
				paddingX={2}
				paddingY={1}
				width={'100%'}
				border={'1px solid rgba(0, 0, 0, 0.12)'}
				// flexGrow={1}
				height={'fit-content'}
				minHeight={600}
				display={'flex'}
				flexDirection={'column'}
				sx={{ backgroundColor: '#fff', userSelect: 'none' }}
			>
				<Breadcrumbs aria-label='breadcrumb' sx={{ mb: 1 }}>
					<Breadcrumb to={AppRoutes.Home}>Главная</Breadcrumb>
					<Breadcrumb to={AppRoutes.OrdersList} active>
						Список позиций
					</Breadcrumb>
				</Breadcrumbs>

				<FlatOrders />
			</Box>
		</PageBox>
	)
}
