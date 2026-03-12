import { Box } from '@mui/material'
import { useLocation } from 'react-router'

import { Order } from '@/features/orders/components/Order/Order'
import { PageBox } from '@/components/PageBox/PageBox'

export default function OrderPage() {
	const location = useLocation()

	const id = location.pathname.split('/').pop() || ''
	const positionIds = location.search.split('positions=')[1]?.split(',') || []

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
				<Order id={id} positionIds={positionIds} />
			</Box>
		</PageBox>
	)
}
