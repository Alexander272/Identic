import { useState } from 'react'
import {
	Paper,
	Table,
	TableBody,
	TableCell,
	TableContainer,
	TableHead,
	TableRow,
	Chip,
	IconButton,
	Tooltip,
	Box,
	Typography,
	Grid,
} from '@mui/material'
import { ExpandMore, ExpandLess } from '@mui/icons-material'

import type { SearchLog } from '../types/search'
import { getSearchTypeLabel, getSearchTypeColor, formatDate, formatDuration } from './utils'

interface SearchTableProps {
	data: SearchLog[]
}

export const SearchTable = ({ data }: SearchTableProps) => {
	const [expandedId, setExpandedId] = useState<string | null>(null)

	const toggleExpand = (id: string) => {
		setExpandedId(expandedId === id ? null : id)
	}

	return (
		<Paper elevation={0} sx={{ border: '1px solid #eee', borderRadius: 2 }}>
			<TableContainer>
				<Table>
					<TableHead>
						<TableRow>
							<TableCell />
							<TableCell>Дата</TableCell>
							<TableCell>Пользователь</TableCell>
							<TableCell>Тип поиска</TableCell>
							<TableCell>Запрос</TableCell>
							<TableCell align='right'>Результатов</TableCell>
							<TableCell align='right'>Время</TableCell>
						</TableRow>
					</TableHead>
					<TableBody>
						{data.map(log => (
							<>
								<TableRow key={log.id} hover>
									<TableCell>
										<IconButton size='small' onClick={() => toggleExpand(log.id)}>
											{expandedId === log.id ? <ExpandLess /> : <ExpandMore />}
										</IconButton>
									</TableCell>
									<TableCell>{formatDate(log.createdAt)}</TableCell>
									<TableCell>{log.actorName}</TableCell>
									<TableCell>
										<Chip
											label={getSearchTypeLabel(log.searchType)}
											color={getSearchTypeColor(log.searchType)}
											size='small'
										/>
									</TableCell>
									<TableCell>
										<Typography
											sx={{
												maxWidth: 200,
												overflow: 'hidden',
												textOverflow: 'ellipsis',
												whiteSpace: 'nowrap',
											}}
										>
											{typeof log.query === 'string'
												? log.query
												: JSON.stringify(log.query)}
										</Typography>
									</TableCell>
									<TableCell align='right'>{log.resultsCount}</TableCell>
									<TableCell align='right'>
										<Tooltip title={`${log.durationMs}мс`}>
											<span>{formatDuration(log.durationMs)}</span>
										</Tooltip>
									</TableCell>
								</TableRow>
								{expandedId === log.id && (
									<TableRow>
										<TableCell colSpan={7}>
											<Box sx={{ pl: 4, py: 2 }}>
												<Grid container spacing={2}>
													<Grid size={{ xs: 12, md: 4 }}>
														<Typography variant='body2' color='text.secondary'>
															Найдено результатов: {log.resultsCount}
														</Typography>
													</Grid>
													<Grid size={{ xs: 12, md: 4 }}>
														<Typography variant='body2' color='text.secondary'>
															Элементов: {log.itemsCount}
														</Typography>
													</Grid>
													<Grid size={{ xs: 12, md: 4 }}>
														<Typography variant='body2' color='text.secondary'>
															ID поиска: {log.searchId}
														</Typography>
													</Grid>
												</Grid>
											</Box>
										</TableCell>
									</TableRow>
								)}
							</>
						))}
					</TableBody>
				</Table>
			</TableContainer>
		</Paper>
	)
}