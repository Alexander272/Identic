import { type FC, useCallback, useState } from 'react'
import { Button, ButtonGroup, Menu, MenuItem, Typography, CircularProgress } from '@mui/material'

import { useExportPricesMutation } from '@/features/prices/priceApiSlice'
import { LeftArrowIcon } from '@/components/Icons/LeftArrowIcon'
import { ExcelIcon } from '@/components/Icons/ExcelIcon'
import { buttonSx } from '../buttonStyles'

type ExportButtonProps = {
	lastParams: { queries?: string[]; codes?: string[]; fields?: string[] }
	visibleColumns: string[]
	allColumnKeys: string[]
	selectedCodes?: string[]
}

export const ExportButton: FC<ExportButtonProps> = ({
	lastParams,
	visibleColumns,
	allColumnKeys,
	selectedCodes = [],
}) => {
	const [exportPositions, { isLoading: isExporting }] = useExportPricesMutation()
	const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
	const open = Boolean(anchorEl)

	const hasSelection = selectedCodes.length > 0

	const handleExport = useCallback(
		async (columns: string[], codes?: string[]) => {
			const params = codes && codes.length > 0 ? { columns, codes } : { ...lastParams, columns }
			await exportPositions(params)
		},
		[lastParams, exportPositions],
	)

	const allVisible = visibleColumns.length === allColumnKeys.length

	return (
		<>
			<ButtonGroup size='small' color='inherit'>
				<Button
					variant='outlined'
					onClick={() => handleExport(allColumnKeys, selectedCodes)}
					disabled={isExporting}
					sx={buttonSx}
				>
					{isExporting ? <CircularProgress size={16} /> : <ExcelIcon sx={{ fontSize: 16 }} />}
					<Typography sx={{ ml: 1 }}>Экспорт</Typography>
				</Button>
				{(hasSelection || !allVisible) && (
					<Button variant='outlined' sx={buttonSx} onClick={e => setAnchorEl(e.currentTarget)}>
						<LeftArrowIcon fontSize={12} transform={'rotate(-90deg)'} />
					</Button>
				)}
			</ButtonGroup>
			<Menu
				anchorEl={anchorEl}
				open={open}
				onClose={() => setAnchorEl(null)}
				anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
				transformOrigin={{ vertical: 'top', horizontal: 'right' }}
			>
				{hasSelection && (
					<MenuItem
						dense
						onClick={() => {
							handleExport(allColumnKeys)
							setAnchorEl(null)
						}}
					>
						Экспортировать всё
					</MenuItem>
				)}
				{!allVisible && (
					<MenuItem
						dense
						onClick={() => {
							handleExport(visibleColumns, selectedCodes)
							setAnchorEl(null)
						}}
					>
						Экспортировать видимое
					</MenuItem>
				)}
			</Menu>
		</>
	)
}
