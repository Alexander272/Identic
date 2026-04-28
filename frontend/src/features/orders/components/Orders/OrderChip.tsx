import type { FC } from 'react'
import { Box } from '@mui/material'

export const OrderChip: FC<{ type: 'bargaining' | 'budget' }> = ({ type }) => {
	return (
		<Box
			sx={{
				display: 'inline-flex',
				alignItems: 'center',
				gap: 0.75,
				px: 1.5,
				py: 0.5,
				borderRadius: 2,
				fontSize: 10,
				fontWeight: 600,
				background: type === 'bargaining' ? 'rgb(231 92 92 / 15%)' : 'rgba(0, 184, 148, 0.1)',
				color: type == 'bargaining' ? '#fb0c20' : '#00b894',
			}}
		>
			{type === 'bargaining' ? 'Т' : 'Б'}
		</Box>
	)
}
