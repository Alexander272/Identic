import { useState, type FC, useCallback, useRef, useEffect } from 'react'
import { Box } from '@mui/material'

import type { Price, SearchPriceRequest } from '@/features/prices/types/types'
import { useGetPricesQuery, useSearchPriceMutation } from '@/features/prices/priceApiSlice'
import type { IFetchError } from '@/app/types/error'
import { STORAGE_KEYS } from '@/constants/storage'
import { SearchModes } from '@/features/prices/components/search/SearchModes'
import { ResultsTable } from '@/features/prices/components/table/ResultsTable'

const ROWS_PER_PAGE_OPTIONS = [15, 30, 50, 100]

const getInitialRowsPerPage = (): number => {
	try {
		const stored = localStorage.getItem(STORAGE_KEYS.pricesRowsPerPage)
		if (stored && ROWS_PER_PAGE_OPTIONS.includes(JSON.parse(stored))) return JSON.parse(stored)
	} catch {
		// ignore
	}
	return 15
}

export const SearchPage: FC = () => {
	const [page, setPage] = useState(1)
	const [rowsPerPage, setRowsPerPage] = useState(getInitialRowsPerPage)

	useEffect(() => {
		localStorage.setItem(STORAGE_KEYS.pricesRowsPerPage, String(rowsPerPage))
	}, [rowsPerPage])

	const [results, setResults] = useState<Price[]>([])
	const [totalCount, setTotalCount] = useState(0)
	const [isLoading, setIsLoading] = useState(false)
	const [error, setError] = useState<unknown>(null)

	const [mode, setMode] = useState<'browse' | 'search'>('browse')
	const [searchQueries, setSearchQueries] = useState<string[]>([])
	const [searchCodes, setSearchCodes] = useState<string[]>([])

	const textSearchParams = useRef<{ queries: string[]; fields?: string[] } | null>(null)

	const {
		data: pricesData,
		isFetching: isPricesLoading,
		error: pricesError,
	} = useGetPricesQuery({ page, limit: rowsPerPage }, { skip: mode !== 'browse' })

	const [search] = useSearchPriceMutation()

	const displayResults = mode === 'browse' ? (pricesData?.data ?? []) : results
	const displayTotalCount = mode === 'browse' ? (pricesData?.total ?? 0) : totalCount
	const displayIsLoading = mode === 'browse' ? isPricesLoading : isLoading
	const displayError = mode === 'browse' ? pricesError : error
	const showPagination = mode === 'browse' || searchCodes.length === 0

	const handleSearch = useCallback(
		async (params: { queries?: string[]; fields?: string[]; codes?: string[] }) => {
			const isText = !!params.queries?.length
			const isCodes = !!params.codes?.length
			if (!isText && !isCodes) return

			setMode('search')
			setSearchQueries(params.queries ?? [])
			setSearchCodes(params.codes ?? [])
			setPage(1)
			setIsLoading(true)
			setError(null)
			textSearchParams.current = isText ? { queries: params.queries!, fields: params.fields } : null

			try {
				const body: SearchPriceRequest = isCodes
					? { codes: params.codes }
					: { queries: params.queries, fields: params.fields, page: 1, limit: rowsPerPage }
				const result = await search(body).unwrap()
				setResults(result.data ?? [])
				setTotalCount(isCodes ? (result.data?.length ?? 0) : (result.total ?? 0))
			} catch (err) {
				const fetchErr = err as IFetchError
				setError(fetchErr.data ?? err)

				setResults([])
				setTotalCount(0)
			} finally {
				setIsLoading(false)
			}
		},
		[search, rowsPerPage],
	)

	const handlePageChange = useCallback(
		(newPage: number) => {
			setPage(newPage + 1)
			const params = textSearchParams.current
			if (params) {
				search({ queries: params.queries, fields: params.fields, page: newPage + 1, limit: rowsPerPage })
					.unwrap()
					.then(result => {
						setResults(result.data ?? [])
						setTotalCount(result.total ?? 0)
					})
					.catch(err => {
						const fetchErr = err as IFetchError
						setError(fetchErr.data ?? err)
					})
			}
		},
		[search, rowsPerPage],
	)

	const handleRowsPerPageChange = useCallback(
		(newRowsPerPage: number) => {
			setRowsPerPage(newRowsPerPage)
			setPage(1)
			const params = textSearchParams.current
			if (params) {
				search({ queries: params.queries, fields: params.fields, page: 1, limit: newRowsPerPage })
					.unwrap()
					.then(result => {
						setResults(result.data ?? [])
						setTotalCount(result.total ?? 0)
					})
					.catch(err => {
						const fetchErr = err as IFetchError
						setError(fetchErr.data ?? err)
					})
			}
		},
		[search],
	)

	const handleResetSearch = useCallback(() => {
		setMode('browse')
		setSearchQueries([])
		setSearchCodes([])
		setResults([])
		setTotalCount(0)
		setError(null)
		setIsLoading(false)
		setPage(1)
		setRowsPerPage(15)
		textSearchParams.current = null
	}, [])

	return (
		<Box sx={{ p: 3, maxWidth: 1800, width: '100%', mx: 'auto' }}>
			<SearchModes onSearch={handleSearch} isLoading={displayIsLoading} onResetSearch={handleResetSearch} />

			<ResultsTable
				results={displayResults}
				queries={searchQueries}
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
