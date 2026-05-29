import { useState, type FC, useCallback } from 'react'
import { Box } from '@mui/material'

import type { Price } from '@/features/prices/types/types'
import { SearchModes } from '@/features/prices/components/search/SearchModes'
import { ResultsTable } from '@/features/prices/components/ResultsTable'

export const SearchPage: FC = () => {
	const [results, setResults] = useState<Price[]>([])
	const [queries, setQueries] = useState<string[]>([])
	const [isLoading, setIsLoading] = useState(false)
	const [error, setError] = useState<unknown>(null)
	const [hasSearched, setHasSearched] = useState(false)
	const [lastParams, setLastParams] = useState<{ queries?: string[]; codes?: string[] }>({})

	const handleSearchResults = useCallback((data: Price[], params: { queries?: string[]; codes?: string[] }) => {
		setResults(data)
		setQueries(params.queries ?? [])
		setLastParams(params)
		setHasSearched(true)
	}, [])

	const handleError = useCallback((err: unknown) => {
		setError(err)
		setHasSearched(true)
	}, [])

	return (
		<Box sx={{ p: 3, maxWidth: 1800, width: '100%', mx: 'auto' }}>
			<SearchModes onSearchResults={handleSearchResults} onError={handleError} onLoadingChange={setIsLoading} />

			<ResultsTable
				results={results}
				queries={queries}
				isLoading={isLoading}
				error={error}
				hasFilters={hasSearched}
				lastParams={lastParams}
			/>
		</Box>
	)
}
