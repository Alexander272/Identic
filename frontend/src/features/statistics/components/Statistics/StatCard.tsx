import { Box, Paper, Typography } from '@mui/material'
import type { ReactNode } from 'react'

interface StatCardProps {
	icon: ReactNode
	title: string
	value: number | undefined
	label?: string
	color?: string
}

const iconColors: Record<string, string> = {
	'#6c5ce7': 'rgba(108, 92, 231, 0.15)',
	'#00b894': 'rgba(0, 184, 148, 0.1)',
	'#fdcb6e': 'rgba(253, 203, 110, 0.1)',
	'#ff6b6b': 'rgba(255, 107, 107, 0.1)',
	'#74b9ff': 'rgba(116, 185, 255, 0.1)',
	'#a29bfe': 'rgba(162, 155, 254, 0.15)',
	'#2196f3': 'rgba(108, 92, 231, 0.15)',
	'#4caf50': 'rgba(0, 184, 148, 0.1)',
	'#ff9800': 'rgba(253, 203, 110, 0.1)',
	'#f44336': 'rgba(255, 107, 107, 0.1)',
	'#00bcd4': 'rgba(116, 185, 255, 0.1)',
	'#9c27b0': 'rgba(162, 155, 254, 0.15)',
}

const formatValue = (value: number | undefined, label?: string): string => {
	if (value === undefined || value === null) return '0'

	if (label === 'сек' && value >= 1) {
		return `${value.toFixed(2)}с`
	}

	if (value >= 1000) {
		return `${(value / 1000).toFixed(1)}к`
	}

	return String(value)
}

export const StatCard = ({ icon, title, value, label, color = '#6c5ce7' }: StatCardProps) => {
	const iconBg = iconColors[color] || `${color}20`

	return (
		<Paper
			elevation={0}
			sx={{
				px: 1.5,
				py: 1,
				border: '1px solid',
				borderColor: 'divider',
				borderRadius: 3,
				height: '100%',
				transition: 'all 0.3s',
				cursor: 'pointer',
				display: 'flex',
				justifyContent: 'center',
				alignItems: 'center',
				gap: 2,
				'&:hover': {
					borderColor: color,
					transform: 'translateY(-2px)',
					boxShadow: '0 4px 24px rgba(0,0,0,0.1)',
				},
			}}
		>
			<Box
				sx={{
					width: 44,
					height: 44,
					borderRadius: 2,
					display: 'flex',
					alignItems: 'center',
					justifyContent: 'center',
					background: iconBg,
					color: color,
					fontSize: 22,
				}}
			>
				{icon}
			</Box>

			<Box>
				<Box sx={{ display: 'flex', justifyContent: 'center', alignItems: 'baseline', gap: 0.5 }}>
					<Typography
						sx={{
							fontSize: 32,
							fontWeight: 800,
							lineHeight: 1.2,
							mb: 0.5,
							color: 'text.primary',
						}}
					>
						{formatValue(value, label)}
					</Typography>
					{label && (
						<Typography
							sx={{
								fontSize: 14,
								fontWeight: 600,
								color: 'text.secondary',
							}}
						>
							{label}
						</Typography>
					)}
				</Box>
				<Typography
					textAlign={'center'}
					sx={{
						fontSize: 13,
						color: 'text.secondary',
					}}
				>
					{title}
				</Typography>
			</Box>
		</Paper>
	)
}
