import { Box } from '@mui/material'

import { PageBox } from '@/components/PageBox/PageBox'
import { CreateOrderForm } from '@/features/orders/components/CreateOrder/CreateForm'

export default function CreateOrder() {
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
				<CreateOrderForm />
			</Box>
		</PageBox>
	)
}
