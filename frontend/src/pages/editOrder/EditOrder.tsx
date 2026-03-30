import { Box, Breadcrumbs } from '@mui/material'
import { useParams } from 'react-router'

import { PageBox } from '@/components/PageBox/PageBox'
import { AppRoutes } from '../router/routes'
import { EditOrderForm } from '@/features/orders/components/EditOrder/EditForm'
import { Breadcrumb } from '@/components/Breadcrumb/Breadcrumb'

export default function EditOrder() {
	const { id } = useParams()

	return (
		<PageBox>
			<Box
				borderRadius={3}
				paddingX={2}
				paddingY={1}
				width={'80%'}
				alignSelf={'center'}
				border={'1px solid rgba(0, 0, 0, 0.12)'}
				// flexGrow={1}
				height={'fit-content'}
				minHeight={600}
				// maxHeight={800}
				display={'flex'}
				flexDirection={'column'}
				sx={{ backgroundColor: '#fff' }}
			>
				<Breadcrumbs aria-label='breadcrumb' sx={{ mb: -1 }}>
					<Breadcrumb to={AppRoutes.Home}>Главная</Breadcrumb>
					<Breadcrumb to={AppRoutes.CreateOrder} active>
						Редактировать заказ
					</Breadcrumb>
				</Breadcrumbs>

				<EditOrderForm orderId={id || ''} />
			</Box>
		</PageBox>
	)
}
