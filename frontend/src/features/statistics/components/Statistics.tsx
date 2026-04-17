import { useState } from 'react'
import { Box, Typography, Grid, Tabs, Tab } from '@mui/material'
import { TrendingUp, Search } from '@mui/icons-material'

import { useGetSearchLogsQuery, useGetActivityLogsQuery } from '../statisticsApiSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { StatCard } from './StatCard'
import { SearchCharts } from './SearchCharts'
import { ActivityCharts } from './ActivityCharts'
import { SearchTable } from './SearchTable'
import { ActivityTable } from './ActivityTable'

interface TabPanelProps {
	children?: React.ReactNode
	index: number
	value: number
}

const TabPanel = ({ children, value, index }: TabPanelProps) => (
	<Box hidden={value !== index} sx={{ py: 3 }}>
		{value === index && children}
	</Box>
)

export const Statistics = () => {
	const [tabValue, setTabValue] = useState(0)

	const { data: searchData, isFetching: isFetchingSearch } = useGetSearchLogsQuery(null)
	const { data: activityData, isFetching: isFetchingActivity } = useGetActivityLogsQuery(null)

	const isLoading = isFetchingSearch || isFetchingActivity

	return (
		<Box sx={{ p: 3 }}>
			<Box sx={{ mb: 2 }}>
				<Typography variant='h6'>Статистика поиска и активности</Typography>
			</Box>

			{isLoading && <BoxFallback />}

			<Grid container spacing={3}>
				<Grid size={{ xs: 12, md: 6 }}>
					<StatCard icon={<Search />} title='Поисковые запросы' value={searchData?.total} color='#2196f3' />
				</Grid>

				<Grid size={{ xs: 12, md: 6 }}>
					<StatCard
						icon={<TrendingUp />}
						title='Создание/редактирование заказов и позиций'
						value={activityData?.total}
						color='#4caf50'
					/>
				</Grid>
			</Grid>

			<Tabs
				value={tabValue}
				onChange={(_event, newValue) => {
					console.log('Tab changed to:', newValue)
					setTabValue(newValue)
				}}
				sx={{ mt: 2, borderBottom: 1, borderColor: 'divider' }}
			>
				<Tab label='Графики' />
				<Tab label='Поисковые запросы' />
				<Tab label='Активность' />
			</Tabs>

			<TabPanel value={tabValue} index={0}>
				<Box sx={{ display: 'flex', flexDirection: 'column', gap: 3 }}>
					{searchData?.data && <SearchCharts data={searchData.data} />}
					{activityData?.data && <ActivityCharts data={activityData.data} />}
				</Box>
			</TabPanel>

			<TabPanel value={tabValue} index={1}>
				{searchData?.data && <SearchTable data={searchData.data} />}
			</TabPanel>

			<TabPanel value={tabValue} index={2}>
				{activityData?.data && <ActivityTable data={activityData.data} />}
			</TabPanel>
		</Box>
	)
}
