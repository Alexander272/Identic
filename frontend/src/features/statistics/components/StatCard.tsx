import { Box, Paper, Typography } from '@mui/material'
import type { ReactNode } from 'react'

interface StatCardProps {
	icon: ReactNode
	title: string
	value: number | undefined
	label?: string
	color?: string
}

export const StatCard = ({ icon, title, value, color = '#2196f3' }: StatCardProps) => (
	<Paper
		elevation={0}
		sx={{
			p: 2,
			border: '1px solid #eee',
			borderRadius: 2,
			height: '100%',
			display: 'flex',
			justifyContent: 'space-between',
		}}
	>
		<Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
			<Box sx={{ display: 'flex', alignItems: 'center', gap: 1, color }}>
				{icon}
				<Typography>{title}</Typography>
			</Box>
		</Box>
		<Typography variant='h5' sx={{ fontWeight: 'bold', lineHeight: 1, color }}>
			{value || 0}
		</Typography>
	</Paper>
)
