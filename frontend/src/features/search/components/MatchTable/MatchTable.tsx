import { useState } from 'react'
import { Box, useTheme, useMediaQuery } from '@mui/material'
import { Virtuoso } from 'react-virtuoso'

import type { IOrderMatchResult, ISearchItem } from '../../types/search'
import type { IPosition } from '@/features/orders/types/positions'

import { useMatchTable } from './useMatchTable'
import { MatchRow } from './MatchRow'
import { MatchTableHeader } from './MatchTableHeader'

type Props = {
	request: ISearchItem[]
	result: IOrderMatchResult
	foundPositions: IPosition[]
}

type FilterType = 'all' | 'found' | 'not_found'

export const MatchTable = ({ request, result, foundPositions }: Props) => {
	const theme = useTheme()
	const isMobile = useMediaQuery(theme.breakpoints.down('md'))

	const [filter, setFilter] = useState<FilterType>('all')

	const { rows } = useMatchTable(request, result, foundPositions)

	const filteredRows = rows.filter(r => {
		if (filter === 'all') return true
		if (filter === 'found') return r.status !== 'not_found'
		if (filter === 'not_found') return r.status === 'not_found'
		return true
	})

	return (
		<Box>
			<MatchTableHeader
				filter={filter}
				setFilter={setFilter}
				total={request.length}
				matched={result.matchedPos}
				link={result.link}
			/>

			<Virtuoso
				style={{ height: 700 }}
				data={filteredRows}
				itemContent={(_index, row) => <MatchRow row={row} isMobile={isMobile} />}
			/>
		</Box>
	)
}
