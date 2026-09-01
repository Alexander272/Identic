import { type FC, useCallback, useMemo, useState, useEffect, useRef } from 'react'
import {
	Table,
	TableBody,
	TableCell,
	TableContainer,
	TableHead,
	TableRow,
	Alert,
	Box,
	Stack,
	Typography,
} from '@mui/material'
import { Check } from '@mui/icons-material'

import type { Price } from '@/features/prices/types'
import { PermRules } from '@/features/access/constants/permissions'
import { COLUMNS, STORAGE_KEY, COLUMN_KEYS, createCellRenderers, highlight } from './cells'
import { useCheckPermission } from '@/features/user/hooks/check'
import { isApiError } from '../../utils/error'
import { PaginationBar } from './PaginationBar'
import { ResultsTableToolbar } from './ResultsTableToolbar'

type TableRow = { kind: 'data'; position: Price } | { kind: 'not-found'; code: string }

export type SelectedItem = { uid: number; position: Price; rowKey: string }

export type ResultsTableProps = {
	results: Price[]
	queries: string[]
	fields?: string[]
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
	fields,
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

	const uidRef = useRef(0)

	const [selected, setSelected] = useState<SelectedItem[]>([])

	const selectedItems = useMemo(
		() => selected.map(({ uid, position }) => ({ uid, position })),
		[selected],
	)

	const selectedRowKeys = useMemo(() => new Set(selected.map(s => s.rowKey)), [selected])

	const toggleSelect = useCallback((position: Price, rowKey: string) => {
		setSelected(prev => {
			const existing = prev.find(s => s.rowKey === rowKey)
			if (existing) {
				return prev.filter(s => s.rowKey !== rowKey)
			}
			return [...prev, { uid: ++uidRef.current, position, rowKey }]
		})
	}, [])

	const handleClearSelection = useCallback(() => setSelected([]), [])

	const handleRemoveSelected = useCallback((uid: number) => {
		setSelected(prev => prev.filter(s => s.uid !== uid))
	}, [])

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
		const isDefault =
			visibleColumns.length === COLUMN_KEYS.length && COLUMN_KEYS.every(k => visibleColumns.includes(k))
		if (isDefault) {
			localStorage.removeItem(STORAGE_KEY)
		} else {
			localStorage.setItem(STORAGE_KEY, JSON.stringify(visibleColumns))
		}
	}, [visibleColumns])

	const matchedFields = useMemo(() => {
		const set = new Set<string>()
		for (const r of results) r.matchedFields?.forEach(f => set.add(f))
		return [...set]
	}, [results])

	const renderCell = useCallback(
		(value: string | null, field: string) => {
			const text = value || ''
			if (!queries.length || !matchedFields.includes(field)) return text

			if (field === 'price') {
				const raw = text.replace(/\s/g, '').replace(/,.*$/, '')
				if (queries.some(q => raw.toLowerCase().includes(q.toLowerCase()))) {
					return <span style={{ background: '#ffeb3b', fontWeight: 700 }}>{text}</span>
				}
				return text
			}

			return highlight(text, queries)
		},
		[queries, matchedFields],
	)

	const visibleColumnDefs = useMemo(() => COLUMNS.filter(col => visibleColumns.includes(col.key)), [visibleColumns])

	const tableRows = useMemo<TableRow[]>(() => {
		if (isLoading && mode === 'search') return []
		if (codes && codes.length > 0) {
			const map = new Map<string, Price>()
			for (const p of results) map.set(p.code, p)
			return codes.map(code => {
				const position = map.get(code)
				return position ? { kind: 'data', position } : { kind: 'not-found', code }
			})
		}
		return results.map(p => ({ kind: 'data', position: p }))
	}, [results, codes, isLoading, mode])

	const cellRenderers = useMemo(() => createCellRenderers(renderCell), [renderCell])

	const handleToggleColumn = useCallback((key: string) => {
		setVisibleColumns(prev => (prev.includes(key) ? prev.filter(k => k !== key) : [...prev, key]))
	}, [])

	if (error) {
		return <Alert severity='error'>{isApiError(error) ? error.data.message : 'Ошибка поиска'}</Alert>
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
				lastParams={{ queries, codes, fields }}
				visibleColumns={visibleColumns}
				onToggleColumn={handleToggleColumn}
				canEdit={canEdit}
				selected={selectedItems}
				onRemoveSelected={handleRemoveSelected}
				onClearSelection={handleClearSelection}
			/>

			<TableContainer sx={{ maxHeight: 'calc(100vh - 350px)', borderRadius: 2 }}>
				<Table stickyHeader size='small'>
					<TableHead>
						<TableRow>
							<TableCell sx={{ width: 36, padding: 0.5, backgroundColor: '#fafafa' }} />
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
						{isLoading && (
							<TableRow>
								<TableCell colSpan={visibleColumnDefs.length + 1} align='center' sx={{ py: 4 }}>
									<Stack alignItems='center' justifyContent='center'>
										<Typography fontSize='1.3rem'>Идет поиск...</Typography>
										<Typography fontSize='1.1rem' variant='caption' color='text.secondary'>
											Поиск может занять некоторое время
										</Typography>
									</Stack>
								</TableCell>
							</TableRow>
						)}
						{tableRows.length === 0 && !isLoading && (
							<TableRow>
								<TableCell
									colSpan={visibleColumnDefs.length + 1}
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
						{tableRows.map((row, i) =>
							row.kind === 'data' ? (() => {
								const rowKey = `${i}-${row.position.id}`
								const isSelected = selectedRowKeys.has(rowKey)
								return (
									<TableRow
										key={rowKey}
										onClick={() => toggleSelect(row.position, rowKey)}
										sx={{
											'&:nth-of-type(even)': {
												backgroundColor: isSelected ? '#e3f2fd' : '#fafafa',
											},
											backgroundColor: isSelected ? '#e3f2fd' : undefined,
											'&:hover': {
												backgroundColor: isSelected ? '#dbeafe !important' : '#f0f4f8 !important',
											},
											cursor: 'pointer',
											transition: 'background-color 0.2s ease-in-out',
										}}
									>
										<TableCell sx={{ padding: 0.5, textAlign: 'center' }}>
											{isSelected && <Check sx={{ fontSize: 16, color: 'primary.main' }} />}
										</TableCell>
										{visibleColumnDefs.map(col => cellRenderers[col.key](row.position))}
									</TableRow>
								)
							})() : (
								<TableRow
									key={`${i}-nf-${row.code}`}
									sx={{
										backgroundColor: '#ffced5',
										transition: 'background-color 0.2s ease-in-out',
										'&:hover': { backgroundColor: '#fdb9c3 !important' },
									}}
								>
									<TableCell sx={{ padding: 0.5 }} />
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
