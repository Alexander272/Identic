import { type FC } from 'react'
import { Link } from 'react-router'
import { Box, Button, Typography } from '@mui/material'

import type { Price } from '@/features/prices/types/types'
import { AppRoutes } from '@/pages/router/routes'
import { COLUMN_KEYS, buttonSx } from './cells'
import { ColumnSettings } from './ColumnSettings'
import { ExportButton } from '@/features/prices/components/ExportButton'
import { EditBoxIcon } from '@/components/Icons/EditBoxIcon'

type Props = {
	results: Price[]
	mode: 'browse' | 'search'
	totalCount: number
	lastParams: { queries?: string[]; codes?: string[] }
	visibleColumns: string[]
	onToggleColumn: (key: string) => void
	canEdit: boolean
}

export const ResultsTableToolbar: FC<Props> = ({
	results,
	mode,
	totalCount,
	lastParams,
	visibleColumns,
	onToggleColumn,
	canEdit,
}) => {
	return (
		<Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', mb: 1 }}>
			{results.length > 0 && (
				<Box sx={{ display: 'flex', gap: 3, alignItems: 'center' }}>
					<ColumnSettings visibleColumns={visibleColumns} onToggleColumn={onToggleColumn} />
					<Typography variant='body2' sx={{ color: 'text.secondary' }}>
						{mode === 'browse' ? `Всего: ${totalCount}` : `Найдено: ${totalCount}`}
					</Typography>
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
					<ExportButton lastParams={lastParams} visibleColumns={visibleColumns} allColumnKeys={COLUMN_KEYS} />
				)}
			</Box>
		</Box>
	)
}
