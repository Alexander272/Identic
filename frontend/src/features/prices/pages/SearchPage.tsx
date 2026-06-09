import { type FC } from 'react'
import { Box } from '@mui/material'

import { usePriceSearch } from '@/features/prices/hooks/usePriceSearch'
import { SearchModes } from '@/features/prices/components/search/SearchModes'
import { ResultsTable } from '@/features/prices/components/table/ResultsTable'

export const SearchPage: FC = () => {
	const {
		mode,
		searchQueries,
		searchFields,
		searchCodes,
		page,
		rowsPerPage,
		displayResults,
		displayTotalCount,
		displayIsLoading,
		displayError,
		showPagination,
		handleSearch,
		handlePageChange,
		handleRowsPerPageChange,
		handleResetSearch,
	} = usePriceSearch()

	return (
		<Box sx={{ p: 3, maxWidth: 1800, width: '100%', mx: 'auto' }}>
			<SearchModes onSearch={handleSearch} isLoading={displayIsLoading} onResetSearch={handleResetSearch} />

			<ResultsTable
				results={displayResults}
				queries={searchQueries}
				fields={searchFields}
				isLoading={displayIsLoading}
				error={displayError}
				hasFilters={mode === 'search'}
				codes={searchCodes}
				mode={mode}
				totalCount={displayTotalCount}
				page={page}
				rowsPerPage={rowsPerPage}
				showPagination={showPagination}
				onPageChange={handlePageChange}
				onRowsPerPageChange={handleRowsPerPageChange}
			/>
		</Box>
	)
}
