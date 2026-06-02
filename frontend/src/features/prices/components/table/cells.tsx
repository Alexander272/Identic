/* eslint-disable react-refresh/only-export-components */

import type { ReactNode } from 'react'
import { TableCell, type Theme } from '@mui/material'

import type { Price } from '@/features/prices/types/types'
import { priceFormat } from '@/utils/format'

export type ColumnDef = { key: string; label: string; sx?: Record<string, unknown> }

export const COLUMNS: ColumnDef[] = [
	{ key: 'code', label: 'Код', sx: { fontFamily: 'monospace', fontWeight: 600, whiteSpace: 'nowrap' as const } },
	{ key: 'current_name', label: 'Наименование СИБУР' },
	{ key: 'new_name', label: 'Наименование СИЛУР (для проверки спецификации)' },
	{ key: 'price', label: 'Цена, руб', sx: { whiteSpace: 'nowrap' as const } },
	{ key: 'template', label: 'Шаблон' },
	{ key: 'note', label: 'Примечание для СИЛУР' },
	{ key: 'need_sibur_approval', label: 'Требуется доп.согл. с СИБУР' },
]

export const COLUMN_KEYS = COLUMNS.map(c => c.key)

import { STORAGE_KEYS } from '@/constants/storage'

export const STORAGE_KEY = STORAGE_KEYS.resultsTableVisibleColumns

export const buttonSx = ({ palette }: Theme) => ({
	minWidth: 48,
	textTransform: 'inherit',
	background: '#fff',
	border: '1px solid #c3c3c4',
	borderRadius: '6px',
	padding: '4px 10px',
	':hover': { svg: { fill: palette.primary.main }, color: palette.primary.main },
	'&:disabled': { svg: { fill: palette.action.disabled } },
})

const normalize = (s: string) => s.toLowerCase().replace(/[х]/g, 'x')

const esc = (s: string) => s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')

export const highlight = (text: string | null, searchQueries: string[]) => {
	if (!text || !searchQueries.length) return text || ''
	const terms = searchQueries.filter(Boolean).map(q => esc(q).replace(/[xх]/gi, '[xх]'))
	if (!terms.length) return text || ''
	const regex = new RegExp(`(${terms.join('|')})`, 'gi')
	const parts = text.split(regex)
	const normalizedQueries = searchQueries.map(normalize)
	return parts.map((part, i) =>
		part && normalizedQueries.some(q => normalize(part).replace(/\s+/g, '') === q) ? (
			<span key={i} style={{ background: '#ffeb3b', fontWeight: 700 }}>
				{part}
			</span>
		) : (
			part
		),
	)
}

export const createCellRenderers = (renderCell: (value: string | null, field: string) => ReactNode): Record<string, (row: Price) => ReactNode> => ({
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
	need_sibur_approval: (row: Price) => <TableCell key='need_sibur_approval'>{row.need_sibur_approval || ''}</TableCell>,
})
