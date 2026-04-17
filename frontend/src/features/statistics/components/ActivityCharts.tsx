import { Paper, Typography, Grid } from '@mui/material'
import {
	BarChart,
	Bar,
	XAxis,
	YAxis,
	CartesianGrid,
	Tooltip as ChartTooltip,
	Legend,
	PieChart,
	Pie,
	Cell,
	ResponsiveContainer,
} from 'recharts'

import type { ActivityLog } from '../types/activity'
import { getActionLabel, getEntityTypeLabel } from './utils'

const COLORS = ['#2196f3', '#4caf50', '#ff9800', '#f44336', '#9c27b0']

interface ActivityChartsProps {
	data: ActivityLog[]
}

export const ActivityCharts = ({ data }: ActivityChartsProps) => {
	const activityByActionData = () => {
		const grouped = data.reduce(
			(acc, log) => {
				const action = log.action
				if (!acc[action]) {
					acc[action] = { action: getActionLabel(action), count: 0 }
				}
				acc[action].count++
				return acc
			},
			{} as Record<string, { action: string; count: number }>
		)
		return Object.values(grouped)
	}

	const activityByEntityData = () => {
		if (!data) return []
		const grouped = data.reduce(
			(acc, log) => {
				const entity = log.entityType
				if (!acc[entity]) {
					acc[entity] = { entity: getEntityTypeLabel(entity), count: 0 }
				}
				acc[entity].count++
				return acc
			},
			{} as Record<string, { entity: string; count: number }>
		)
		return Object.values(grouped)
	}

	const byAction = activityByActionData()
	const byEntity = activityByEntityData()

	return (
		<Grid container spacing={3}>
			<Grid size={{ xs: 12, lg: 6 }}>
				<Paper elevation={0} sx={{ p: 3, border: '1px solid #eee', borderRadius: 2 }}>
					<Typography variant='h6' sx={{ mb: 2 }}>
						Активность по действиям
					</Typography>
					<ResponsiveContainer width='100%' height={300}>
						<PieChart>
							<Pie
								data={byAction}
								dataKey='count'
								nameKey='action'
								cx='50%'
								cy='50%'
								outerRadius={100}
								label
							>
								{byAction.map((_, index) => (
									<Cell
										key={`cell-${index}`}
										fill={COLORS[index % COLORS.length]}
									/>
								))}
							</Pie>
							<ChartTooltip />
							<Legend />
						</PieChart>
					</ResponsiveContainer>
				</Paper>
			</Grid>

			<Grid size={{ xs: 12, lg: 6 }}>
				<Paper elevation={0} sx={{ p: 3, border: '1px solid #eee', borderRadius: 2 }}>
					<Typography variant='h6' sx={{ mb: 2 }}>
						Активность по сущностям
					</Typography>
					<ResponsiveContainer width='100%' height={300}>
						<BarChart data={byEntity}>
							<CartesianGrid strokeDasharray='3 3' />
							<XAxis dataKey='entity' />
							<YAxis />
							<ChartTooltip contentStyle={{ borderRadius: 8 }} />
							<Legend />
							<Bar
								dataKey='count'
								name='Количество'
								fill='#4caf50'
								radius={[4, 4, 0, 0]}
							/>
						</BarChart>
					</ResponsiveContainer>
				</Paper>
			</Grid>
		</Grid>
	)
}