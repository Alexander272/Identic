import { Box } from '@mui/material'
import { useParams, useSearchParams } from 'react-router'

import { Order } from '@/features/orders/components/Order/Order'
import { PageBox } from '@/components/PageBox/PageBox'

export default function OrderPage() {
	const [searchParams] = useSearchParams()

	const { id } = useParams()
	const positionIds = searchParams.get('positions')?.split(',') || []

	// console.log('id', id)
	console.log('positionIds', positionIds)

	return (
		<PageBox>
			<Box
				borderRadius={3}
				paddingX={2}
				paddingY={1}
				minWidth={'80%'}
				border={'1px solid rgba(0, 0, 0, 0.12)'}
				// flexGrow={1}
				height={'fit-content'}
				minHeight={600}
				// maxHeight={800}
				display={'flex'}
				flexDirection={'column'}
				sx={{ backgroundColor: '#fff' }}
			>
				<Order id={id || ''} positionIds={positionIds} />
			</Box>
		</PageBox>
	)
}
