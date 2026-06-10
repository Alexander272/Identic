import type { FC } from 'react'
import { Box, Button, Stack, useTheme } from '@mui/material'
import { DataSheetGrid, textColumn, floatColumn, keyColumn, type Column } from 'react-datasheet-grid'
import type { Operation } from 'react-datasheet-grid/dist/types'
import './editPositionsGrid.css'

import type { GridRow } from '../utils/grid'
import { ContextMenu } from '@/components/DataSheet/ContextMenu'
import { AddRow } from '@/components/DataSheet/AddRow'
import { SaveIcon } from '@/components/Icons/SaveIcon'

const defaultRow = (): GridRow => ({
	id: new Date().getTime().toString() + Math.random().toString(36).slice(2, 9),
	code: '',
	currentName: '',
	newName: null,
	price: null,
	template: null,
	note: null,
	needSiburApproval: null,
	status: 'CREATED',
})

const columns: Column<GridRow>[] = [
	{ ...keyColumn<GridRow, 'code'>('code', textColumn), title: 'Код', width: 0.5 },
	{ ...keyColumn<GridRow, 'currentName'>('currentName', textColumn), title: 'Наименование', width: 1.5 },
	{ ...keyColumn<GridRow, 'newName'>('newName', textColumn), title: 'Новое наименование', width: 1.5 },
	{ ...keyColumn<GridRow, 'price'>('price', floatColumn), title: 'Цена', width: 0.5 },
	{ ...keyColumn<GridRow, 'template'>('template', textColumn), title: 'Шаблон', width: 1 },
	{ ...keyColumn<GridRow, 'note'>('note', textColumn), title: 'Примечание', width: 1 },
	{
		...keyColumn<GridRow, 'needSiburApproval'>('needSiburApproval', textColumn),
		title: 'Требуется доп.согл. с СИБУР',
		width: 1,
	},
]

type Props = {
	rows: GridRow[]
	onRowsChange: (rows: GridRow[]) => void
	onSave: () => void
	isSaving: boolean
	hasChanges: boolean
	hasInvalid: boolean
}

export const EditPositionsGrid: FC<Props> = ({ rows, onRowsChange, onSave, isSaving, hasChanges, hasInvalid }) => {
	const { palette } = useTheme()

	const changeHandler = (newRows: GridRow[], operations: Operation[]) => {
		// DELETE indices refer to positions in `rows` (old value),
		// because the grid removes rows before calling onChange.
		// CREATE and UPDATE indices refer to positions in `newRows`.
		const rowsToDelete = new Set<number>()

		for (const op of operations) {
			if (op.type !== 'DELETE') continue
			for (let i = op.fromRowIndex; i < op.toRowIndex; i++) rowsToDelete.add(i)
		}

		// Step 1: Apply deletions to old rows
		const result: GridRow[] = []
		for (let i = 0; i < rows.length; i++) {
			if (rowsToDelete.has(i)) {
				if (rows[i].status !== 'CREATED') result.push({ ...rows[i], status: 'DELETED' })
			} else {
				result.push({ ...rows[i] })
			}
		}

		// Step 2: Apply updates from grid to result
		for (const op of operations) {
			if (op.type !== 'UPDATE') continue
			for (let i = op.fromRowIndex; i < op.toRowIndex; i++) {
				const newRow = newRows[i]
				if (!newRow?.id) continue
				const pos = result.findIndex(r => r.id === newRow.id)
				if (pos !== -1) {
					const status = result[pos].status === 'ORIGINAL' ? 'UPDATED' : result[pos].status
					result[pos] = { ...newRow, status }
				}
			}
		}

		// Step 3: Add genuinely new rows (not in old `rows`)
		for (const op of operations) {
			if (op.type !== 'CREATE') continue
			for (let i = op.fromRowIndex; i < op.toRowIndex; i++) {
				const newRow = newRows[i]
				if (!newRow) continue
				if (!newRow.id || !result.some(r => r.id === newRow.id)) result.push({ ...newRow, status: 'CREATED' })
			}
		}

		onRowsChange(result)
	}

	const addClasses = ({ rowData }: { rowData: GridRow }) => {
		switch (rowData.status) {
			case 'DELETED':
				return 'row-deleted'
			case 'CREATED':
				return 'row-created'
			case 'UPDATED':
				return 'row-updated'
			default:
				return ''
		}
	}

	return (
		<Box sx={{ mb: 1, position: 'relative' }}>
			<DataSheetGrid
				value={rows}
				onChange={changeHandler}
				columns={columns}
				createRow={defaultRow}
				rowClassName={addClasses}
				autoAddRow
				height={500}
				contextMenuComponent={props => <ContextMenu {...props} />}
				addRowsComponent={props => <AddRow {...props} />}
			/>
			<Stack direction='row' spacing={1} sx={{ position: 'absolute', right: 8, bottom: 6 }}>
				<Button
					onClick={onSave}
					color='inherit'
					disabled={!hasChanges || isSaving || hasInvalid}
					sx={{
						minWidth: 48,
						textTransform: 'inherit',
						background: '#fff',
						border: '1px solid #dcdcdc',
						borderRadius: '6px',
						padding: '4px 10px',
						':hover': { svg: { fill: palette.primary.main }, color: palette.primary.main },
						'&:disabled': { svg: { fill: palette.action.disabled } },
					}}
				>
					<SaveIcon mr={1} fontSize={16} />
					Сохранить
				</Button>
			</Stack>
		</Box>
	)
}
