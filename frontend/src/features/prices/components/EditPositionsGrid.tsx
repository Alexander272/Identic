import type { FC } from 'react'
import { Box, Button, Stack, useTheme } from '@mui/material'
import { DataSheetGrid, textColumn, floatColumn, keyColumn, type Column } from 'react-datasheet-grid'
import type { Operation } from 'react-datasheet-grid/dist/types'
import './editPositionsGrid.css'

import { ContextMenu } from '@/components/DataSheet/ContextMenu'
import { AddRow } from '@/components/DataSheet/AddRow'
import { SaveIcon } from '@/components/Icons/SaveIcon'
import type { GridRow } from '../utils/grid'

const defaultRow = (): GridRow => ({
	id: '',
	code: '',
	current_name: '',
	new_name: null,
	price: null,
	template: null,
	note: null,
	technique: null,
	under_drawing: null,
	status: 'CREATED',
})

const columns: Column<GridRow>[] = [
	{ ...keyColumn<GridRow, 'code'>('code', textColumn), title: 'Код', width: 0.5 },
	{ ...keyColumn<GridRow, 'current_name'>('current_name', textColumn), title: 'Наименование', width: 1.5 },
	{ ...keyColumn<GridRow, 'new_name'>('new_name', textColumn), title: 'Новое наименование', width: 1.5 },
	{ ...keyColumn<GridRow, 'price'>('price', floatColumn), title: 'Цена', width: 0.5 },
	{ ...keyColumn<GridRow, 'template'>('template', textColumn), title: 'Шаблон', width: 1 },
	{ ...keyColumn<GridRow, 'note'>('note', textColumn), title: 'Примечание', width: 1 },
	{ ...keyColumn<GridRow, 'technique'>('technique', textColumn), title: 'Техника', width: 1 },
	{ ...keyColumn<GridRow, 'under_drawing'>('under_drawing', textColumn), title: 'Под чертеж', width: 1 },
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
		const updated = [...newRows]

		for (const op of operations) {
			if (op.type === 'DELETE') {
				for (let i = op.fromRowIndex; i < op.toRowIndex; i++) {
					const row = updated[i]
					if (row && row.status === 'CREATED') {
						updated.splice(i, 1)
						i--
						op.toRowIndex--
					} else if (row) {
						row.status = 'DELETED'
					}
				}
			}
			if (op.type === 'UPDATE') {
				for (let i = op.fromRowIndex; i < op.toRowIndex; i++) {
					const row = updated[i]
					if (row && row.status === 'ORIGINAL') {
						row.status = 'UPDATED'
					}
				}
			}
			if (op.type === 'CREATE') {
				for (let i = op.fromRowIndex; i < op.toRowIndex; i++) {
					if (updated[i]) {
						updated[i].status = 'CREATED'
					}
				}
			}
		}

		onRowsChange(updated)
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
