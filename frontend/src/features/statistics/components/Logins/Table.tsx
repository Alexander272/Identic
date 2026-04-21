import { useMemo, useState } from 'react'
import { Paper, Table, TableBody, TableContainer, TableHead, TableRow, TableCell, Box, Typography } from '@mui/material'
import dayjs from 'dayjs'
import { Pagination } from '@/components/Pagination/Pagination'

import type { SearchLog } from '../../types/search'
import type { ActivityLog, EntityType } from '../../types/activity'
import type { IUserLoginWithUser } from '../../types/userLogins'
import { stringToHSLA } from '@/utils/colors'
import { getSmartDate } from '@/utils/date'
import { getInitials } from '../utils'

interface UserActivityStats {
	id: string
	firstName: string
	lastName: string
	email: string
	searches: number
	ordersModified: number
	itemsModified: number
	lastActivityAt: string | null
	lastLoginAt: string | null
	isOnline: boolean
}

interface LoginsTableProps {
	searchData: SearchLog[]
	activityData: ActivityLog[]
	userLogins: IUserLoginWithUser[]
}

const headerCellStyle = {
	fontSize: '11px',
	fontWeight: 700,
	textTransform: 'uppercase',
	letterSpacing: '0.5px',
}

const isUserOnline = (lastActivityAt: string): boolean => {
	return dayjs().diff(dayjs(lastActivityAt), 'minute') < 15
}

const ROWS_PER_PAGE = 10

export const LoginsTable = ({ searchData, activityData, userLogins }: LoginsTableProps) => {
	const [page, setPage] = useState(1)

	const userStats = useMemo(() => {
		const statsMap = new Map<string, UserActivityStats>()

		userLogins.forEach(login => {
			const userId = login.userId
			statsMap.set(userId, {
				id: userId,
				firstName: login.user.firstName,
				lastName: login.user.lastName,
				email: login.user.email,
				searches: 0,
				ordersModified: 0,
				itemsModified: 0,
				lastActivityAt: login.lastActivityAt,
				lastLoginAt: login.loginAt,
				isOnline: isUserOnline(login.lastActivityAt),
			})
		})

		searchData.forEach(log => {
			const userId = log.actor.id
			const stats = statsMap.get(userId)

			if (stats) {
				stats.searches++
			}
		})

		activityData.forEach(log => {
			const userId = log.changedBy
			const stats = statsMap.get(userId)

			if (stats) {
				if (log.entityType === ('order' as EntityType)) {
					stats.ordersModified++
				} else if (log.entityType === ('order_item' as EntityType)) {
					stats.itemsModified++
				}

				if (stats.lastActivityAt === null || dayjs(log.createdAt).isAfter(dayjs(stats.lastActivityAt))) {
					stats.lastActivityAt = log.createdAt
					stats.isOnline = isUserOnline(log.createdAt)
				}
			}
		})

		return Array.from(statsMap.values()).sort((a, b) => {
			if (a.isOnline !== b.isOnline) return a.isOnline ? -1 : 1
			return 0
		})
	}, [searchData, activityData, userLogins])

	const paginatedUsers = useMemo(() => {
		const start = (page - 1) * ROWS_PER_PAGE
		return userStats.slice(start, start + ROWS_PER_PAGE)
	}, [userStats, page])

	return (
		<Paper
			elevation={0}
			sx={{
				border: '1px solid',
				borderColor: 'divider',
				borderRadius: 3,
				overflow: 'hidden',
			}}
		>
			<TableContainer>
				<Table>
					<TableHead>
						<TableRow sx={{ background: 'action.hover' }}>
							<TableCell sx={headerCellStyle}>Пользователь</TableCell>
							<TableCell sx={{ ...headerCellStyle, textAlign: 'center' }}>Статус</TableCell>
							<TableCell sx={{ ...headerCellStyle, textAlign: 'center' }}>Поисков</TableCell>
							<TableCell sx={{ ...headerCellStyle, textAlign: 'center' }}>Заказов изменено</TableCell>
							<TableCell sx={{ ...headerCellStyle, textAlign: 'center' }}>Позиций изменено</TableCell>
							<TableCell sx={headerCellStyle}>Последний вход</TableCell>
						</TableRow>
					</TableHead>
					<TableBody>
						{paginatedUsers.length === 0 && (
							<TableRow>
								<TableCell colSpan={6} sx={{ py: 2, textAlign: 'center', fontWeight: 'bold' }}>
									<Typography>Ничего не найдено</Typography>
								</TableCell>
							</TableRow>
						)}

						{paginatedUsers.map(user => {
							const fullName = `${user.lastName} ${user.firstName}`.trim()
							const colors = stringToHSLA(fullName || user.email)

							return (
								<TableRow key={user.id} hover>
									<TableCell>
										<Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
											<Box
												sx={{
													width: 32,
													height: 32,
													borderRadius: '50%',
													background: colors.bg,
													color: colors.text,
													border: `1px solid ${colors.border}`,
													display: 'flex',
													alignItems: 'center',
													justifyContent: 'center',
													fontSize: 12,
													fontWeight: 700,
												}}
											>
												{getInitials(fullName || user.email)}
											</Box>
											<Box>
												<Typography sx={{ fontWeight: 600 }}>
													{fullName || user.email}
												</Typography>
												<Typography sx={{ fontSize: 12, color: 'text.secondary' }}>
													{user.email}
												</Typography>
											</Box>
										</Box>
									</TableCell>
									<TableCell align='center'>
										<Box
											sx={{
												display: 'inline-flex',
												alignItems: 'center',
												gap: 0.75,
												px: 1.5,
												py: 0.5,
												borderRadius: 10,
												fontSize: 14,
												fontWeight: 600,
												background: user.isOnline
													? 'rgba(0, 184, 148, 0.1)'
													: 'rgba(178, 190, 195, 0.2)',
												color: user.isOnline ? '#00b894' : '#636e72',
											}}
										>
											<Box
												sx={{
													width: 6,
													height: 6,
													borderRadius: '50%',
													background: 'currentColor',
												}}
											/>
											{user.isOnline ? 'Онлайн' : 'Офлайн'}
										</Box>
									</TableCell>
									<TableCell align='center'>
										<Typography sx={{ fontWeight: 600 }}>{user.searches}</Typography>
									</TableCell>
									<TableCell align='center'>
										<Typography sx={{ fontWeight: 600, color: '#0984e3' }}>
											{user.ordersModified}
										</Typography>
									</TableCell>
									<TableCell align='center'>
										<Typography sx={{ fontWeight: 600, color: '#6c5ce7' }}>
											{user.itemsModified}
										</Typography>
									</TableCell>
									<TableCell>
										{user.lastLoginAt ? (
											<Typography sx={{ color: 'text.secondary' }}>
												{getSmartDate(user.lastLoginAt)}
											</Typography>
										) : (
											<Typography sx={{ color: 'text.disabled' }}>-</Typography>
										)}
									</TableCell>
								</TableRow>
							)
						})}
					</TableBody>
				</Table>
			</TableContainer>
			{userStats.length > ROWS_PER_PAGE && (
				<Box sx={{ py: 2, display: 'flex', justifyContent: 'center' }}>
					<Pagination
						page={page}
						totalPages={Math.ceil(userStats.length / ROWS_PER_PAGE)}
						onClick={setPage}
					/>
				</Box>
			)}
		</Paper>
	)
}
