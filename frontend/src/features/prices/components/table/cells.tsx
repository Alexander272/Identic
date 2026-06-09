/* eslint-disable react-refresh/only-export-components */

import type { ReactNode } from 'react'
import { TableCell } from '@mui/material'

import type { Price } from '@/features/prices/types'
import { priceFormat } from '@/utils/format'

import { STORAGE_KEYS } from '@/constants/storage'

export type ColumnDef = { key: string; label: string; sx?: Record<string, unknown> }

export const COLUMNS: ColumnDef[] = [
	{ key: 'code', label: 'Код', sx: { fontFamily: 'monospace', fontWeight: 600, whiteSpace: 'nowrap' as const } },
	{ key: 'currentName', label: 'Наименование СИБУР' },
	{ key: 'newName', label: 'Наименование СИЛУР (для проверки спецификации)' },
	{ key: 'price', label: 'Цена, руб', sx: { whiteSpace: 'nowrap' as const } },
	{ key: 'template', label: 'Шаблон' },
	{ key: 'note', label: 'Примечание для СИЛУР' },
	{ key: 'needSiburApproval', label: 'Требуется доп.согл. с СИБУР' },
]

export const COLUMN_KEYS = COLUMNS.map(c => c.key)

export const STORAGE_KEY = STORAGE_KEYS.resultsTableVisibleColumns

const normalize = (s: string) => s.toLowerCase().replace(/[х]/g, 'x')

const esc = (s: string) => s.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')

export const highlight = (text: string | null, searchQueries: string[]) => {
	if (!text || !searchQueries.length) return text || ''
	const terms = searchQueries
		.filter(Boolean)
		.map(q => esc(q.trim()).replace(/\s+/g, '\\s+').replace(/[xх]/gi, '[xх]'))
	if (!terms.length) return text || ''
	const regex = new RegExp(`(${terms.join('|')})`, 'gi')
	const parts = text.split(regex)
	const normalizedSet = new Set(searchQueries.map(q => normalize(q).trim().replace(/\s+/g, '')))
	return parts.map((part, i) => {
		if (!part) return part
		return normalizedSet.has(normalize(part).replace(/\s+/g, '')) ? (
			<span key={i} style={{ background: '#ffeb3b', fontWeight: 700 }}>
				{part}
			</span>
		) : (
			part
		)
	})
}

export const createCellRenderers = (
	renderCell: (value: string | null, field: string) => ReactNode,
): Record<string, (row: Price) => ReactNode> => ({
	code: (row: Price) => (
		<TableCell key='code' sx={{ fontFamily: 'monospace', fontWeight: 600, whiteSpace: 'nowrap' }}>
			{renderCell(row.code, 'code')}
		</TableCell>
	),
	currentName: (row: Price) => <TableCell key='currentName'>{renderCell(row.currentName, 'currentName')}</TableCell>,
	newName: (row: Price) => <TableCell key='newName'>{renderCell(row.newName, 'newName')}</TableCell>,
	price: (row: Price) => (
		<TableCell key='price' sx={{ whiteSpace: 'nowrap' }}>
			{renderCell(priceFormat(row.price || 0), 'price')}
		</TableCell>
	),
	template: (row: Price) => <TableCell key='template'>{renderCell(row.template, 'template')}</TableCell>,
	note: (row: Price) => <TableCell key='note'>{row.note || ''}</TableCell>,
	needSiburApproval: (row: Price) => <TableCell key='needSiburApproval'>{row.needSiburApproval || ''}</TableCell>,
})
