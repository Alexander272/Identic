import { type FC, useCallback } from 'react'
import {
	Button,
	Dialog,
	DialogContent,
	DialogTitle,
	IconButton,
	Table,
	TableBody,
	TableCell,
	TableContainer,
	TableHead,
	TableRow,
	Stack,
	Typography,
	CircularProgress,
} from '@mui/material'

import type { Price } from '@/features/prices/types'
import { COLUMNS } from './cells'
import { useExportPricesMutation } from '@/features/prices/priceApiSlice'
import { TimesIcon } from '@/components/Icons/TimesIcon'
import { TrashBinIcon } from '@/components/Icons/TrashBinIcon'
import { ExcelIcon } from '@/components/Icons/ExcelIcon'
import { buttonSx } from '../buttonStyles'

type Props = {
	open: boolean
	selected: Price[]
	visibleColumns: string[]
	onClose: () => void
	onRemove: (code: string) => void
	onClear: () => void
}

export const SelectedPositionsDialog: FC<Props> = ({ open, selected, visibleColumns, onClose, onRemove, onClear }) => {
	const [exportPositions, { isLoading: isExporting }] = useExportPricesMutation()
	const visibleCols = COLUMNS.filter(col => visibleColumns.includes(col.key))

	const handleExport = useCallback(async () => {
		const codes = selected.map(p => p.code)
		await exportPositions({ codes, columns: visibleColumns })
	}, [selected, visibleColumns, exportPositions])

	return (
		<Dialog
			open={open}
			onClose={onClose}
			fullWidth
			maxWidth='xl'
			slotProps={{
				paper: {
					sx: { borderRadius: '16px', p: 1 },
				},
				backdrop: {
					sx: { backdropFilter: 'blur(4px)', backgroundColor: 'rgba(0,0,0,0.4)' },
				},
			}}
		>
			<DialogTitle sx={{ m: 0, p: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
				<Typography variant='h6' sx={{ fontWeight: 'bold' }}>
					Выбранные позиции ({selected.length})
				</Typography>

				<Stack direction='row' spacing={8} alignItems='center'>
					<Stack direction='row' spacing={1.5} alignItems='center'>
						{selected.length > 0 && (
							<Button
								size='small'
								variant='outlined'
								color='inherit'
								onClick={handleExport}
								disabled={isExporting}
								sx={buttonSx}
							>
								{isExporting ? <CircularProgress size={16} /> : <ExcelIcon sx={{ fontSize: 16 }} />}
								<Typography sx={{ ml: 1 }}>Экспорт</Typography>
							</Button>
						)}
						{selected.length > 0 && (
							<Button size='small' color='inherit' onClick={onClear} sx={buttonSx}>
								<TrashBinIcon sx={{ fontSize: 16 }} />
								<Typography sx={{ ml: 1 }}>Очистить</Typography>
							</Button>
						)}
					</Stack>

					<IconButton onClick={onClose} sx={{ color: 'text.secondary' }}>
						<TimesIcon fontSize={16} />
					</IconButton>
				</Stack>
			</DialogTitle>

			<DialogContent sx={{ p: 0, pt: '4px !important' }}>
				{selected.length > 0 ? (
					<TableContainer sx={{ maxHeight: 'calc(100vh - 300px)' }}>
						<Table stickyHeader size='small'>
							<TableHead>
								<TableRow>
									{visibleCols.map(col => (
										<TableCell
											key={col.key}
											sx={{ fontWeight: 700, backgroundColor: '#fafafa', ...col.sx }}
										>
											{col.label}
										</TableCell>
									))}
									<TableCell sx={{ fontWeight: 700, backgroundColor: '#fafafa', width: 48 }} />
								</TableRow>
							</TableHead>
							<TableBody>
								{selected.map(item => (
									<TableRow key={item.code} hover>
										{visibleCols.map(col => (
											<TableCell key={col.key} sx={col.sx}>
												{col.key === 'price'
													? (item.price ?? '').toLocaleString('ru-RU')
													: String(item[col.key as keyof Price] ?? '')}
											</TableCell>
										))}
										<TableCell sx={{ textAlign: 'center' }}>
											<IconButton
												size='large'
												onClick={() => onRemove(item.code)}
												sx={{ color: 'text.secondary', '&:hover': { color: 'error.main' } }}
											>
												<TimesIcon fontSize={14} />
											</IconButton>
										</TableCell>
									</TableRow>
								))}
							</TableBody>
						</Table>
					</TableContainer>
				) : (
					<Typography sx={{ color: 'text.secondary', textAlign: 'center', py: 4 }}>
						Нет выбранных позиций
					</Typography>
				)}
			</DialogContent>
		</Dialog>
	)
}
