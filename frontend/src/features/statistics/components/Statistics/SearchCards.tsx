import { useMemo, type FC } from 'react'
import { Box } from '@mui/material'
import { Search, AccessTime, ErrorOutline, Inbox } from '@mui/icons-material'

import type { SearchLog } from '../../types'
import { StatCard } from './StatCard'

type Props = {
	data: SearchLog[]
	total: number
}

export const SearchCards: FC<Props> = ({ data, total }) => {
	const calculatedStats = useMemo(() => {
		if (!data) {
			return { avgTime: 0, noResults: 0, foundOrders: 0 }
		}

		const totalDuration = data.reduce((sum, log) => sum + log.durationMs, 0)
		const avgTime = data.length > 0 ? totalDuration / data.length / 1000 : 0

		const noResults = data.filter(log => log.resultsCount === 0).length

		const foundOrders = data.reduce((sum, log) => sum + log.resultsCount, 0)

		return { avgTime, noResults, foundOrders }
	}, [data])

	return (
		<Box
			sx={{
				display: 'grid',
				gridTemplateColumns: {
					xs: '1fr',
					sm: 'repeat(2, 1fr)',
					lg: 'repeat(4, 1fr)',
				},
				gap: 3,
				mb: 3,
			}}
		>
			<StatCard icon={<Search />} title='Поисковые запросы' value={total} color='#6c5ce7' />

			<StatCard icon={<Inbox />} title='Найдено заказов' value={calculatedStats.foundOrders} color='#00b894' />

			<StatCard
				icon={<AccessTime />}
				title='Среднее время поиска'
				value={Math.round(calculatedStats.avgTime * 100) / 100}
				label='сек'
				color='#74b9ff'
			/>

			<StatCard
				icon={<ErrorOutline />}
				title='Запросов без результатов'
				value={calculatedStats.noResults}
				color='#ff6b6b'
			/>
		</Box>
	)
}
