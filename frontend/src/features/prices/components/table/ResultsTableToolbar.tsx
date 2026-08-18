import { type FC, useState } from 'react'
import { Link } from 'react-router'
import { Box, Button, Chip, Typography } from '@mui/material'

import type { Price } from '@/features/prices/types'
import { AppRoutes } from '@/pages/router/routes'
import { buttonSx } from '../buttonStyles'
import { COLUMN_KEYS } from './cells'
import { ColumnSettings } from './ColumnSettings'
import { ExportButton } from '@/features/prices/components/table/ExportButton'
import { SelectedPositionsDialog } from '@/features/prices/components/table/SelectedPositionsDialog'
import { EditBoxIcon } from '@/components/Icons/EditBoxIcon'

type Props = {
	results: Price[]
	mode: 'browse' | 'search'
	totalCount: number
	lastParams: { queries?: string[]; codes?: string[]; fields?: string[] }
	visibleColumns: string[]
	onToggleColumn: (key: string) => void
	canEdit: boolean
	selected: Price[]
	onRemoveSelected: (code: string) => void
	onClearSelection: () => void
}

export const ResultsTableToolbar: FC<Props> = ({
	results,
	mode,
	totalCount,
	lastParams,
	visibleColumns,
	onToggleColumn,
	canEdit,
	selected,
	onRemoveSelected,
	onClearSelection,
}) => {
	const [dialogOpen, setDialogOpen] = useState(false)

	return (
		<Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1, gap: 1 }}>
			{results.length > 0 && (
				<Box sx={{ display: 'flex', gap: 3, alignItems: 'center' }}>
					<ColumnSettings visibleColumns={visibleColumns} onToggleColumn={onToggleColumn} />
					<Typography variant='body2' sx={{ color: 'text.secondary' }}>
						{mode === 'browse' ? `Всего: ${totalCount}` : `Найдено: ${totalCount}`}
					</Typography>
					{selected.length > 0 && (
						<>
							<Chip
								label={`Выбрано: ${selected.length}`}
								size='small'
								color='primary'
								variant='outlined'
								onClick={() => setDialogOpen(true)}
								sx={{ cursor: 'pointer' }}
							/>
							<SelectedPositionsDialog
								open={dialogOpen}
								selected={selected}
								visibleColumns={visibleColumns}
								onClose={() => setDialogOpen(false)}
								onRemove={onRemoveSelected}
								onClear={() => {
									onClearSelection()
									setDialogOpen(false)
								}}
							/>
						</>
					)}
				</Box>
			)}

			<Box sx={{ display: 'flex', gap: 1, ml: 'auto' }}>
				{canEdit ? (
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
				) : null}

				{results.length > 0 && (
					<ExportButton
						lastParams={lastParams}
						visibleColumns={visibleColumns}
						allColumnKeys={COLUMN_KEYS}
						selectedCodes={selected.map(p => p.code)}
					/>
				)}
			</Box>
		</Box>
	)
}
