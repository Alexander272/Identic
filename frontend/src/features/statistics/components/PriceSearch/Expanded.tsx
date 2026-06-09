import { TableRow, TableCell, Box, Typography, Collapse, Chip, Stack, Table, TableBody } from '@mui/material'

import { COLUMNS } from '@/features/prices/components/table/cells'
import type { PriceSearchLog } from '../../types/priceSearch'

interface PriceSearchTableExpandedProps {
	log: PriceSearchLog
	open: boolean
}

const FIELD_LABELS = Object.fromEntries(COLUMNS.map(c => [c.key, c.label]))

export const PriceSearchTableExpanded = ({ log, open }: PriceSearchTableExpandedProps) => {
	const isCodes = (log.codes ?? []).length > 0

	return (
		<TableRow>
			<TableCell
				colSpan={8}
				sx={{
					py: 0,
					borderTop: '1px solid',
					borderColor: 'divider',
					background: 'action.hover',
					borderBottom: open ? '1px solid #eee' : 'none',
				}}
			>
				<Collapse in={open} timeout='auto' unmountOnExit>
					<Box sx={{ px: 4, py: 1.5, display: 'flex', flexDirection: 'column', gap: 1 }}>
						{isCodes ? (
							<>
								<Typography variant='subtitle2' sx={{ fontWeight: 600 }}>
									Коды поиска
								</Typography>
								<Stack direction='row' spacing={1}>
									{(log.codes ?? []).map((c, i) => (
										<Chip
											key={`${c}-${i}`}
											label={c}
											size='small'
											sx={{
												fontFamily: 'monospace',
												height: '26px',
												minWidth: '26px',
												backgroundColor: '#f1f5f9',
												border: 'none',
												borderRadius: 2,
											}}
										/>
									))}
								</Stack>
							</>
						) : (
							<Table size='small'>
								<TableBody>
									{(log.queries ?? []).length > 0 && (
										<TableRow>
											<TableCell sx={{ fontWeight: 700 }}>Запрос:</TableCell>
											<TableCell>{(log.queries ?? [])[0]}</TableCell>
										</TableRow>
									)}
									{(log.queries ?? []).length > 1 && (
										<TableRow>
											<TableCell sx={{ fontWeight: 700 }}>Доп запрос:</TableCell>
											<TableCell>{(log.queries ?? []).slice(1).join(', ')}</TableCell>
										</TableRow>
									)}
									{(log.fields ?? []).length > 0 && (
										<TableRow>
											<TableCell sx={{ fontWeight: 700 }}>Поля:</TableCell>
											<TableCell>
												<Stack direction='row' spacing={1} flexWrap='wrap' useFlexGap>
													{(log.fields ?? []).map((f, i) => (
														<Chip
															key={`${f}-${i}`}
															label={FIELD_LABELS[f] ?? f}
															size='small'
															sx={{
																height: '26px',
																px: 1.5,
																backgroundColor: '#f1f5f9',
																border: 'none',
																borderRadius: 2,
															}}
														/>
													))}
												</Stack>
											</TableCell>
										</TableRow>
									)}
								</TableBody>
							</Table>
						)}
					</Box>
				</Collapse>
			</TableCell>
		</TableRow>
	)
}
