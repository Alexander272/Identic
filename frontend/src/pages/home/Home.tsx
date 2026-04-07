import { Box } from '@mui/material'

import { OrdersList } from '@/features/orders/components/Orders/OrdersList'
import { PageBox } from '@/components/PageBox/PageBox'

export default function Home() {
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
				sx={{ backgroundColor: '#fff' }}
			>
				<OrdersList />
			</Box>
		</PageBox>
	)
}
