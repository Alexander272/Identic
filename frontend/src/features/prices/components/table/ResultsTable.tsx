import { type FC, useCallback, useMemo, useState, useEffect } from 'react'
import {
	Table,
	TableBody,
	TableCell,
	TableContainer,
	TableHead,
	TableRow,
	Alert,
	Box,
} from '@mui/material'

import type { Price } from '@/features/prices/types/types'
import { COLUMNS, STORAGE_KEY, COLUMN_KEYS, createCellRenderers, highlight } from './cells'
import { ResultsTableToolbar } from './ResultsTableToolbar'
import { PaginationBar } from './PaginationBar'
import { useCheckPermission } from '@/features/user/hooks/check'
import { PermRules } from '@/features/access/constants/permissions'

type TableRow = { kind: 'data'; position: Price } | { kind: 'not-found'; code: string }

export type ResultsTableProps = {
	results: Price[]
	queries: string[]
	isLoading: boolean
	error: unknown
	hasFilters: boolean
	codes?: string[]
	mode: 'browse' | 'search'
	totalCount: number
	page: number
	rowsPerPage: number
	showPagination: boolean
	onPageChange: (page: number) => void
	onRowsPerPageChange: (rowsPerPage: number) => void
}

export const ResultsTable: FC<ResultsTableProps> = ({
	results,
	queries,
	isLoading,
	error,
	hasFilters,
	codes,
	mode,
	totalCount,
	page,
	rowsPerPage,
	showPagination,
	onPageChange,
	onRowsPerPageChange,
}) => {
	const canEdit = useCheckPermission(PermRules.Prices.Write)

	const [visibleColumns, setVisibleColumns] = useState<string[]>(() => {
		try {
			const stored = localStorage.getItem(STORAGE_KEY)
			if (stored) {
				const parsed = JSON.parse(stored)
				const valid = parsed.filter((k: string) => COLUMN_KEYS.includes(k))
				if (valid.length > 0) return valid
			}
		} catch {
			console.error('Failed to parse resultsTable_visibleColumns')
		}
		return [...COLUMN_KEYS]
	})

	useEffect(() => {
		localStorage.setItem(STORAGE_KEY, JSON.stringify(visibleColumns))
	}, [visibleColumns])

	const matchedFields = useMemo(() => (results.length > 0 ? (results[0].matched_fields ?? []) : []), [results])

	const renderCell = useCallback(
		(value: string | null, field: string) => {
			const text = value || ''
			if (!queries.length || !matchedFields.includes(field)) return text
			return highlight(text, queries)
		},
		[queries, matchedFields],
	)

	const visibleColumnDefs = useMemo(() => COLUMNS.filter(col => visibleColumns.includes(col.key)), [visibleColumns])

	const tableRows = useMemo<TableRow[]>(() => {
		if (codes && codes.length > 0) {
			const map = new Map<string, Price>()
			for (const p of results) map.set(p.code, p)
			return codes.map(code => {
				const position = map.get(code)
				return position ? { kind: 'data', position } : { kind: 'not-found', code }
			})
		}
		return results.map(p => ({ kind: 'data', position: p }))
	}, [results, codes])

	const cellRenderers = useMemo(() => createCellRenderers(renderCell), [renderCell])

	const handleToggleColumn = (key: string) => {
		setVisibleColumns(prev => (prev.includes(key) ? prev.filter(k => k !== key) : [...prev, key]))
	}

	if (error) {
		return (
			<Alert severity='error'>
				{'data' in (error as object)
					? (error as { data?: { message?: string } }).data?.message || 'Ошибка поиска'
					: 'Ошибка поиска'}
			</Alert>
		)
	}

	return (
		<Box
			sx={{
				borderRadius: { xs: 2, sm: 3 },
				paddingY: 2,
				px: 2,
				flex: 1,
				minWidth: 320,
				position: 'relative',
				overflow: 'hidden',
				bgcolor: 'rgba(255,255,255,0.85)',
				border: '1px solid rgba(0,0,0,0.08)',
				backdropFilter: 'blur(20px)',
				boxShadow: '0 4px 12px rgba(0,0,0,0.04)',
				mb: 2,
			}}
		>
			<ResultsTableToolbar
				results={results}
				mode={mode}
				totalCount={totalCount}
				lastParams={{ codes }}
				visibleColumns={visibleColumns}
				onToggleColumn={handleToggleColumn}
				canEdit={canEdit}
			/>

			<TableContainer sx={{ maxHeight: 'calc(100vh - 350px)', borderRadius: 2 }}>
				<Table stickyHeader size='small'>
					<TableHead>
						<TableRow>
							{visibleColumnDefs.map(col => (
								<TableCell
									key={col.key}
									sx={{ fontWeight: 700, backgroundColor: '#fafafa', ...col.sx }}
								>
									{col.label}
								</TableCell>
							))}
						</TableRow>
					</TableHead>
					<TableBody>
						{tableRows.length === 0 && !isLoading && (
							<TableRow>
								<TableCell
									colSpan={visibleColumnDefs.length}
									align='center'
									sx={{ py: 4, color: 'text.secondary' }}
								>
									{mode === 'browse'
										? 'Нет данных'
										: hasFilters
											? 'Ничего не найдено'
											: 'Введите запрос для поиска'}
								</TableCell>
							</TableRow>
						)}
						{tableRows.map(row =>
							row.kind === 'data' ? (
								<TableRow
									key={row.position.id}
									sx={{
										'&:nth-of-type(even)': { backgroundColor: '#fafafa' },
										'&:hover': { backgroundColor: '#f0f4f8 !important' },
										transition: 'background-color 0.2s ease-in-out',
									}}
								>
									{visibleColumnDefs.map(col => cellRenderers[col.key](row.position))}
								</TableRow>
							) : (
								<TableRow
									key={`nf-${row.code}`}
									sx={{
										backgroundColor: '#ffced5',
										transition: 'background-color 0.2s ease-in-out',
										'&:hover': { backgroundColor: '#fdb9c3 !important' },
									}}
								>
									<TableCell sx={{ fontFamily: 'monospace', fontWeight: 600, whiteSpace: 'nowrap' }}>
										{row.code}
									</TableCell>
									{visibleColumnDefs.length > 1 && (
										<TableCell
											colSpan={visibleColumnDefs.length - 1}
											sx={{
												color: '#F44336',
												fontWeight: 'bold',
												textAlign: 'center',
											}}
										>
											— Не найдено —
										</TableCell>
									)}
								</TableRow>
							),
						)}
					</TableBody>
				</Table>
			</TableContainer>

			{showPagination && totalCount > 0 && (
				<PaginationBar
					page={page}
					rowsPerPage={rowsPerPage}
					totalCount={totalCount}
					onPageChange={onPageChange}
					onRowsPerPageChange={onRowsPerPageChange}
				/>
			)}
		</Box>
	)
}
