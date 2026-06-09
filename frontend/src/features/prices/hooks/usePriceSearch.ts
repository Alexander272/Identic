import { useState, useCallback, useEffect } from 'react'

import { useGetPricesQuery, useLazySearchPriceQuery } from '@/features/prices/priceApiSlice'
import { STORAGE_KEYS } from '@/constants/storage'

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

type SearchParams = {
	queries?: string[]
	fields?: string[]
	codes?: string[]
}

export const usePriceSearch = () => {
	const [page, setPage] = useState(1)
	const [rowsPerPage, setRowsPerPage] = useState(getInitialRowsPerPage)

	useEffect(() => {
		localStorage.setItem(STORAGE_KEYS.pricesRowsPerPage, String(rowsPerPage))
	}, [rowsPerPage])

	const [mode, setMode] = useState<'browse' | 'search'>('browse')
	const [searchQueries, setSearchQueries] = useState<string[]>([])
	const [searchFields, setSearchFields] = useState<string[]>([])
	const [searchCodes, setSearchCodes] = useState<string[]>([])

	const {
		data: pricesData,
		isFetching: isPricesLoading,
		error: pricesError,
	} = useGetPricesQuery({ page, limit: rowsPerPage }, { skip: mode !== 'browse' })

	const [search, { data: searchData, isFetching: isSearchFetching, error: searchError }] =
		useLazySearchPriceQuery()

	const displayResults = mode === 'browse' ? (pricesData?.data ?? []) : (searchData?.data ?? [])
	const displayTotalCount = mode === 'browse'
		? (pricesData?.total ?? 0)
		: (searchCodes.length > 0 ? (searchData?.data?.length ?? 0) : (searchData?.total ?? 0))
	const displayIsLoading = mode === 'browse' ? isPricesLoading : isSearchFetching
	const displayError = mode === 'browse' ? pricesError : searchError
	const showPagination = !displayIsLoading && (mode === 'browse' || searchCodes.length === 0)

	const handleSearch = useCallback(
		(params: SearchParams) => {
			const isText = !!params.queries?.length
			const isCodes = !!params.codes?.length
			if (!isText && !isCodes) return

			setMode('search')
			setSearchQueries(params.queries ?? [])
			setSearchFields(params.fields ?? [])
			setSearchCodes(params.codes ?? [])
			setPage(1)

			if (isCodes) {
				search({ codes: params.codes })
			} else {
				search({ queries: params.queries, fields: params.fields, page: 1, limit: rowsPerPage })
			}
		},
		[search, rowsPerPage],
	)

	const handlePageChange = useCallback(
		(newPage: number) => {
			setPage(newPage + 1)
			if (searchQueries.length > 0) {
				search({ queries: searchQueries, fields: searchFields, page: newPage + 1, limit: rowsPerPage })
			}
		},
		[search, searchQueries, searchFields, rowsPerPage],
	)

	const handleRowsPerPageChange = useCallback(
		(newRowsPerPage: number) => {
			setRowsPerPage(newRowsPerPage)
			setPage(1)
			if (searchQueries.length > 0) {
				search({ queries: searchQueries, fields: searchFields, page: 1, limit: newRowsPerPage })
			}
		},
		[search, searchQueries, searchFields],
	)

	const handleResetSearch = useCallback(() => {
		setMode('browse')
		setSearchQueries([])
		setSearchFields([])
		setSearchCodes([])
		setPage(1)
		setRowsPerPage(getInitialRowsPerPage())
	}, [])

	return {
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
	}
}
