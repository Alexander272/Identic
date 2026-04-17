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
	Box,
	Typography,
	Grid,
} from '@mui/material'
import { ExpandMore, ExpandLess, Edit, Delete } from '@mui/icons-material'

import type { ActivityLog } from '../types/activity'
import { ActionType } from '../types/activity'
import { getActionLabel, getActionColor, getEntityTypeLabel, formatDate } from './utils'

interface ActivityTableProps {
	data: ActivityLog[]
}

export const ActivityTable = ({ data }: ActivityTableProps) => {
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
							<TableCell>Действие</TableCell>
							<TableCell>Сущность</TableCell>
							<TableCell>ID сущности</TableCell>
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
									<TableCell>{log.changedByName}</TableCell>
									<TableCell>
										<Chip
											label={getActionLabel(log.action)}
											color={getActionColor(log.action)}
											size='small'
											icon={
												log.action === ActionType.Insert ? (
													<Edit />
												) : log.action === ActionType.Delete ? (
													<Delete />
												) : undefined
											}
										/>
									</TableCell>
									<TableCell>{getEntityTypeLabel(log.entityType)}</TableCell>
									<TableCell>{log.entityId}</TableCell>
								</TableRow>
								{expandedId === log.id && (
									<TableRow>
										<TableCell colSpan={6}>
											<Box sx={{ pl: 4, py: 2 }}>
												<Grid container spacing={2}>
													{log.entity && (
														<Grid size={{ xs: 12 }}>
															<Typography variant='body2' color='text.secondary'>
																Сущность: {log.entity}
															</Typography>
														</Grid>
													)}
													{log.parentId && (
														<Grid size={{ xs: 12 }}>
															<Typography variant='body2' color='text.secondary'>
																Родитель: {log.parentId}
															</Typography>
														</Grid>
													)}
													{log.oldValues && (
														<Grid size={{ xs: 12, md: 6 }}>
															<Typography variant='body2' color='text.secondary' sx={{ mb: 1 }}>
																Старые значения:
															</Typography>
															<Paper variant='outlined' sx={{ p: 2, bgcolor: '#fff3f3' }}>
																<Typography
																	variant='body2'
																	component='pre'
																	sx={{ fontSize: 12 }}
																>
																	{JSON.stringify(log.oldValues, null, 2)}
																</Typography>
															</Paper>
														</Grid>
													)}
													{log.newValues && (
														<Grid size={{ xs: 12, md: 6 }}>
															<Typography variant='body2' color='text.secondary' sx={{ mb: 1 }}>
																Новые значения:
															</Typography>
															<Paper variant='outlined' sx={{ p: 2, bgcolor: '#f3fff3' }}>
																<Typography
																	variant='body2'
																	component='pre'
																	sx={{ fontSize: 12 }}
																>
																	{JSON.stringify(log.newValues, null, 2)}
																</Typography>
															</Paper>
														</Grid>
													)}
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