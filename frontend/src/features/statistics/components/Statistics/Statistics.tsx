import { useCallback, useState, useEffect } from 'react'
import { Box, Typography, Tabs, Tab } from '@mui/material'

import {
	useLazyGetSearchLogsQuery,
	useLazyGetActivityLogsQuery,
	useLazyGetLastUserLoginsQuery,
} from '../../statisticsApiSlice'
import { getDateRange } from '../utils'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { PeriodPicker } from '@/components/Period/Period'
import { SearchTable } from '../Search/Table'
import { ActivityTable } from '../Activity/Table'
import { LoginsTable } from '../Logins/Table'
import { SearchCards } from './SearchCards'
import { OrderCards } from './OrderCards'
import { useCheckPermission } from '@/features/user/hooks/check'
import { PermRules } from '@/features/access/constants/permissions'

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
	// const [period, setPeriod] = useState<Period>('week')
	const [dateRange, setDateRange] = useState<{ startDate: string; endDate: string } | undefined>(getDateRange('week'))

	const [triggerSearch, { data: searchData, isFetching: isFetchingSearch }] = useLazyGetSearchLogsQuery()
	const [triggerActivity, { data: activityData, isFetching: isFetchingActivity }] = useLazyGetActivityLogsQuery()
	const [triggerUserLogins, { data: userLoginsData, isFetching: isFetchingUserLogins }] =
		useLazyGetLastUserLoginsQuery()

	const fetchSearchLogs = useCallback(() => {
		const params = dateRange ? { startDate: dateRange.startDate, endDate: dateRange.endDate } : {}
		triggerSearch(params)
		triggerActivity(params)
		triggerUserLogins(params)
	}, [dateRange, triggerSearch, triggerActivity, triggerUserLogins])

	const isLoading = isFetchingSearch || isFetchingActivity || isFetchingUserLogins

	const periodHandler = useCallback(
		(newRange?: { startDate: string; endDate: string }) => {
			setDateRange(newRange)

			triggerSearch(newRange || {})
			triggerActivity(newRange || {})
			triggerUserLogins(newRange || {})
		},
		[triggerSearch, triggerActivity, triggerUserLogins],
	)

	useEffect(() => {
		fetchSearchLogs()
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [])

	const canSeeSearch = useCheckPermission(PermRules.SearchLog.Read)
	const canSeeActivity = useCheckPermission(PermRules.ActivityLog.Read)
	const canSeeLogins = useCheckPermission(PermRules.Logins.Read)

	return (
		<Box sx={{ p: 3 }}>
			<Box
				sx={{
					mb: 0.5,
					display: 'flex',
					justifyContent: 'space-between',
					alignItems: 'center',
					flexWrap: 'wrap',
					gap: 2,
				}}
			>
				<Typography variant='h6'>Статистика поиска и активности</Typography>

				<PeriodPicker value={dateRange} onChange={periodHandler} />
			</Box>

			<Tabs
				value={tabValue}
				onChange={(_event, newValue) => {
					setTabValue(newValue)
				}}
				sx={{ borderBottom: 1, borderColor: 'divider' }}
			>
				{canSeeSearch && <Tab value={0} label='Поисковые запросы' sx={{ textTransform: 'none' }} />}
				{canSeeActivity && <Tab value={1} label='Заявки' sx={{ textTransform: 'none' }} />}
				{canSeeLogins && <Tab value={2} label='Пользователи' sx={{ textTransform: 'none' }} />}
			</Tabs>

			{canSeeSearch && (
				<TabPanel value={tabValue} index={0}>
					{isLoading && <BoxFallback />}

					<SearchCards data={searchData?.data || []} total={searchData?.total || 0} />

					{searchData?.data && <SearchTable data={searchData.data} />}
				</TabPanel>
			)}

			{canSeeActivity && (
				<TabPanel value={tabValue} index={1}>
					{isLoading && <BoxFallback />}

					<OrderCards data={activityData?.data || []} total={activityData?.total || 0} />

					{activityData?.data && <ActivityTable data={activityData.data} />}
				</TabPanel>
			)}

			{canSeeLogins && (
				<TabPanel value={tabValue} index={2}>
					{isLoading && <BoxFallback />}

					{searchData?.data && activityData?.data && userLoginsData?.data && (
						<LoginsTable
							searchData={searchData.data}
							activityData={activityData.data}
							userLogins={userLoginsData.data}
						/>
					)}
				</TabPanel>
			)}
		</Box>
	)
}
