import { Paper, Typography, Grid } from '@mui/material'
import {
	BarChart,
	Bar,
	XAxis,
	YAxis,
	CartesianGrid,
	Tooltip as ChartTooltip,
	Legend,
	ResponsiveContainer,
} from 'recharts'

import type { SearchLog } from '../types/search'
import { getSearchTypeLabel } from './utils'

interface SearchChartsProps {
	data: SearchLog[]
}

export const SearchCharts = ({ data }: SearchChartsProps) => {
	const searchByTypeData = () => {
		const grouped = data.reduce(
			(acc, log) => {
				const type = log.searchType
				if (!acc[type]) {
					acc[type] = { type: getSearchTypeLabel(type), count: 0, duration: 0 }
				}
				acc[type].count++
				acc[type].duration += log.durationMs
				return acc
			},
			{} as Record<string, { type: string; count: number; duration: number }>
		)
		return Object.values(grouped).map(item => ({
			...item,
			avgDuration: Math.round(item.duration / item.count),
		}))
	}

	const topSearchQueries = () => {
		const grouped = data.reduce(
			(acc, log) => {
				const query = typeof log.query === 'string' ? log.query : JSON.stringify(log.query)
				if (!acc[query]) {
					acc[query] = { query, count: 0, results: 0 }
				}
				acc[query].count++
				acc[query].results += log.resultsCount
				return acc
			},
			{} as Record<string, { query: string; count: number; results: number }>
		)
		return Object.values(grouped)
			.map(q => ({ ...q, avgResults: Math.round(q.results / q.count) }))
			.sort((a, b) => b.count - a.count)
			.slice(0, 10)
	}

	const byType = searchByTypeData()
	const topQueries = topSearchQueries()

	return (
		<Grid container spacing={3}>
			<Grid size={{ xs: 12, lg: 6 }}>
				<Paper elevation={0} sx={{ p: 3, border: '1px solid #eee', borderRadius: 2 }}>
					<Typography variant='h6' sx={{ mb: 2 }}>
						Поиски по типу
					</Typography>
					<ResponsiveContainer width='100%' height={300}>
						<BarChart data={byType}>
							<CartesianGrid strokeDasharray='3 3' />
							<XAxis dataKey='type' />
							<YAxis />
							<ChartTooltip contentStyle={{ borderRadius: 8 }} />
							<Legend />
							<Bar dataKey='count' name='Количество' fill='#2196f3' radius={[4, 4, 0, 0]} />
							<Bar
								dataKey='avgDuration'
								name='Среднее время (мс)'
								fill='#ff9800'
								radius={[4, 4, 0, 0]}
							/>
						</BarChart>
					</ResponsiveContainer>
				</Paper>
			</Grid>

			<Grid size={{ xs: 12, lg: 6 }}>
				<Paper elevation={0} sx={{ p: 3, border: '1px solid #eee', borderRadius: 2 }}>
					<Typography variant='h6' sx={{ mb: 2 }}>
						Топ поисковых запросов
					</Typography>
					<ResponsiveContainer width='100%' height={300}>
						<BarChart data={topQueries} layout='vertical'>
							<CartesianGrid strokeDasharray='3 3' />
							<XAxis type='number' />
							<YAxis
								type='category'
								dataKey='query'
								width={100}
								style={{ fontSize: 12 }}
							/>
							<ChartTooltip contentStyle={{ borderRadius: 8 }} />
							<Legend />
							<Bar dataKey='count' name='Количество' fill='#2196f3' radius={[0, 4, 4, 0]} />
						</BarChart>
					</ResponsiveContainer>
				</Paper>
			</Grid>
		</Grid>
	)
}