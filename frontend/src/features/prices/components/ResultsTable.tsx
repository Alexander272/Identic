import { type FC, type ReactNode, useCallback, useMemo, useState, useEffect } from 'react'
import { Link } from 'react-router'
import {
	Table,
	TableBody,
	TableCell,
	TableContainer,
	TableHead,
	TableRow,
	Typography,
	Alert,
	Box,
	Menu,
	MenuItem,
	FormControlLabel,
	Checkbox,
	Button,
	Badge,
	type Theme,
} from '@mui/material'

import type { Price } from '@/features/prices/types/types'
import { AppRoutes } from '@/pages/router/routes'
import { priceFormat } from '@/utils/format'
import { ExportButton } from '@/features/prices/components/ExportButton'
import { SettingIcon } from '@/components/Icons/SettingIcon'
import { EditBoxIcon } from '@/components/Icons/EditBoxIcon'

const STORAGE_KEY = 'resultsTable_visibleColumns'

type ColumnDef = { key: string; label: string; sx?: Record<string, unknown> }

const COLUMNS: ColumnDef[] = [
	{ key: 'code', label: 'Код', sx: { fontFamily: 'monospace', fontWeight: 600, whiteSpace: 'nowrap' as const } },
	{ key: 'current_name', label: 'Текущее наименование' },
	{ key: 'new_name', label: 'Наименование АСВНСИ' },
	{ key: 'price', label: 'Цена', sx: { whiteSpace: 'nowrap' as const } },
	{ key: 'template', label: 'Шаблон' },
	{ key: 'note', label: 'Примечание' },
	{ key: 'technique', label: 'Техника' },
	{ key: 'under_drawing', label: 'Под чертеж' },
]

const COLUMN_KEYS = COLUMNS.map(c => c.key)

const buttonSx = ({ palette }: Theme) => ({
	minWidth: 48,
	textTransform: 'inherit',
	background: '#fff',
	border: '1px solid #c3c3c4',
	borderRadius: '6px',
	padding: '4px 10px',
	':hover': { svg: { fill: palette.primary.main }, color: palette.primary.main },
	'&:disabled': { svg: { fill: palette.action.disabled } },
})

type ResultsTableProps = {
	results: Price[]
	queries: string[]
	isLoading: boolean
	error: unknown
	hasFilters: boolean
	lastParams: { queries?: string[]; codes?: string[] }
}

const normalize = (s: string) => s.toLowerCase().replace(/[х]/g, 'x')

const esc = (s: string) => s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')

const highlight = (text: string | null, searchQueries: string[]) => {
	if (!text || !searchQueries.length) return text || ''
	const terms = searchQueries.filter(Boolean).map(q => esc(q).replace(/[xх]/gi, '[xх]'))
	if (!terms.length) return text || ''
	const regex = new RegExp(`(${terms.join('|')})`, 'gi')
	const parts = text.split(regex)
	const normalizedQueries = searchQueries.map(normalize)
	return parts.map((part, i) =>
		part && normalizedQueries.some(q => normalize(part) === q) ? (
			<span key={i} style={{ background: '#ffeb3b', fontWeight: 700 }}>
				{part}
			</span>
		) : (
			part
		),
	)
}

type TableRow = { kind: 'data'; position: Price } | { kind: 'not-found'; code: string }

export const ResultsTable: FC<ResultsTableProps> = ({ results, queries, isLoading, error, hasFilters, lastParams }) => {
	const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)

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
		const codes = lastParams.codes
		if (codes && codes.length > 0) {
			const map = new Map<string, Price>()
			for (const p of results) map.set(p.code, p)
			return codes.map(code => {
				const position = map.get(code)
				return position ? { kind: 'data', position } : { kind: 'not-found', code }
			})
		}
		return results.map(p => ({ kind: 'data', position: p }))
	}, [results, lastParams.codes])

	const cellRenderers: Record<string, (row: Price) => ReactNode> = useMemo(
		() => ({
			code: (row: Price) => (
				<TableCell key='code' sx={{ fontFamily: 'monospace', fontWeight: 600, whiteSpace: 'nowrap' }}>
					{renderCell(row.code, 'code')}
				</TableCell>
			),
			current_name: (row: Price) => (
				<TableCell key='current_name'>{renderCell(row.current_name, 'current_name')}</TableCell>
			),
			new_name: (row: Price) => <TableCell key='new_name'>{renderCell(row.new_name, 'new_name')}</TableCell>,
			price: (row: Price) => (
				<TableCell key='price' sx={{ whiteSpace: 'nowrap' }}>
					{renderCell(priceFormat(row.price || 0), 'price')}
				</TableCell>
			),
			template: (row: Price) => <TableCell key='template'>{renderCell(row.template, 'template')}</TableCell>,
			note: (row: Price) => <TableCell key='note'>{row.note || ''}</TableCell>,
			technique: (row: Price) => <TableCell key='technique'>{row.technique || ''}</TableCell>,
			under_drawing: (row: Price) => <TableCell key='under_drawing'>{row.under_drawing || ''}</TableCell>,
		}),
		[renderCell],
	)

	const handleToggleColumn = (key: string) => {
		setVisibleColumns(prev => (prev.includes(key) ? prev.filter(k => k !== key) : [...prev, key]))
	}

	const open = Boolean(anchorEl)

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
			<Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1 }}>
				{results.length > 0 && (
					<Box sx={{ display: 'flex', gap: 3, alignItems: 'center' }}>
						<Button
							onClick={e => setAnchorEl(e.currentTarget)}
							size='small'
							variant='outlined'
							color='inherit'
							sx={buttonSx}
						>
							<Badge color='primary' variant='dot' invisible={visibleColumns.length === COLUMNS.length}>
								<SettingIcon sx={{ fontSize: 18 }} />
							</Badge>
							<Typography sx={{ ml: 1 }}>Настройка колонок</Typography>
						</Button>
						<Menu
							anchorEl={anchorEl}
							open={open}
							onClose={() => setAnchorEl(null)}
							anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
							transformOrigin={{ vertical: 'top', horizontal: 'right' }}
						>
							{COLUMNS.map(col => (
								<MenuItem key={col.key} dense disableRipple>
									<FormControlLabel
										control={
											<Checkbox
												checked={visibleColumns.includes(col.key)}
												onChange={() => handleToggleColumn(col.key)}
												size='small'
											/>
										}
										label={col.label}
										sx={{ m: 0 }}
									/>
								</MenuItem>
							))}
						</Menu>

						<Typography variant='body2' sx={{ color: 'text.secondary' }}>
							Найдено: {results.length}
						</Typography>
					</Box>
				)}

				<Box sx={{ display: 'flex', gap: 1, ml: 'auto' }}>
					<Button
						component={Link}
						to={AppRoutes.PriceEdit}
						size='small'
						variant='outlined'
						color='inherit'
						sx={buttonSx}
					>
						<EditBoxIcon sx={{ fontSize: 16, mr: 1 }} />
						<Typography>Редактировать</Typography>
					</Button>
					{results.length > 0 && (
						<ExportButton
							lastParams={lastParams}
							visibleColumns={visibleColumns}
							allColumnKeys={COLUMN_KEYS}
						/>
					)}
				</Box>
			</Box>

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
									{hasFilters ? 'Ничего не найдено' : 'Введите запрос для поиска'}
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
		</Box>
	)
}
