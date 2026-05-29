import { type FC, useCallback, useState } from 'react'
import { Button, ButtonGroup, Menu, MenuItem, Typography, CircularProgress, type Theme } from '@mui/material'

import { useExportPricesMutation } from '@/features/prices/priceApiSlice'
import { LeftArrowIcon } from '@/components/Icons/LeftArrowIcon'
import { ExcelIcon } from '@/components/Icons/ExcelIcon'

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

type ExportButtonProps = {
	lastParams: { queries?: string[]; codes?: string[] }
	visibleColumns: string[]
	allColumnKeys: string[]
}

export const ExportButton: FC<ExportButtonProps> = ({ lastParams, visibleColumns, allColumnKeys }) => {
	const [exportPositions, { isLoading: isExporting }] = useExportPricesMutation()
	const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
	const open = Boolean(anchorEl)

	const handleExport = useCallback(
		async (columns: string[]) => {
			await exportPositions({ ...lastParams, columns })
		},
		[lastParams, exportPositions],
	)

	const allVisible = visibleColumns.length === allColumnKeys.length

	return (
		<>
			<ButtonGroup size='small' color='inherit'>
				<Button
					variant='outlined'
					onClick={() => handleExport(visibleColumns)}
					disabled={isExporting}
					sx={buttonSx}
				>
					{isExporting ? <CircularProgress size={16} /> : <ExcelIcon sx={{ fontSize: 16 }} />}
					<Typography sx={{ ml: 1 }}>Экспорт</Typography>
				</Button>
				{!allVisible && (
					<Button variant='outlined' sx={buttonSx} onClick={e => setAnchorEl(e.currentTarget)}>
						<LeftArrowIcon sx={{ fontSize: 12, transform: 'rotate(-90deg)' }} />
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
				<MenuItem
					dense
					onClick={() => {
						handleExport(allColumnKeys)
						setAnchorEl(null)
					}}
				>
					Экспортировать все колонки
				</MenuItem>
			</Menu>
		</>
	)
}
