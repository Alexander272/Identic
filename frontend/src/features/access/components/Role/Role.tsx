import { EditBoxIcon } from '@/components/Icons/EditBoxIcon'
import { PlusIcon } from '@/components/Icons/PlusIcon'
import {
	Box,
	Button,
	Chip,
	IconButton,
	Paper,
	Table,
	TableBody,
	TableCell,
	TableContainer,
	TableHead,
	TableRow,
	Typography,
	useTheme,
} from '@mui/material'
import { ArrowRightIcon } from '@mui/x-date-pickers'
import { Fragment } from 'react/jsx-runtime'

const roles = [
	{
		id: 1,
		name: 'Администратор',
		slug: 'admin',
		desc: 'Полный доступ ко всем ресурсам системы',
		icon: '👑',
		iconBg: '#eef2ff', // indigo-100
		iconColor: '#4f46e5', // indigo-600
		status: 'active',
		parents: null,
		perms: { total: 48, own: 48, inherited: 0 },
		userCount: 3,
	},
	{
		id: 2,
		name: 'Менеджер',
		slug: 'manager',
		desc: 'Управляет пользователями и процессами',
		icon: '⚙️',
		iconBg: '#fffbeb', // amber-100
		iconColor: '#d97706', // amber-600
		status: 'active',
		parents: ['Администратор'],
		perms: { total: 42, own: 10, inherited: 32 },
		userCount: 12,
	},
]

export const Role = () => {
	const { palette } = useTheme()

	const isFetching = false

	return (
		<Box sx={{ p: 3 }}>
			{/* Page Header */}
			<Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
				<Box>
					<Typography variant='h4' sx={{ fontWeight: 'bold' }}>
						Роли
					</Typography>
					<Typography variant='body1' color='text.secondary'>
						Управление ролями и их правами
					</Typography>
				</Box>
				<Button
					variant='outlined'
					sx={{ borderRadius: '8px', textTransform: 'none', background: '#fff' }}
					onClick={() => {
						/* openUserModal() */
					}}
				>
					<PlusIcon fill={palette.primary.main} fontSize={16} mr={1.5} />
					Добавить
				</Button>
			</Box>

			<TableContainer
				component={Paper}
				elevation={0}
				sx={{ borderRadius: '24px', border: '1px solid #f3f4f6', overflow: 'hidden' }}
			>
				<Table sx={{ minWidth: 800 }}>
					<TableHead>
						<TableRow sx={{ borderBottom: '1px solid #f3f4f6' }}>
							<TableCell sx={{ py: 2.5, px: 4, color: 'text.secondary', fontSize: '0.875rem' }}>
								Роль
							</TableCell>
							<TableCell sx={{ py: 2.5, px: 3, color: 'text.secondary', fontSize: '0.875rem' }}>
								Статус
							</TableCell>
							<TableCell sx={{ py: 2.5, px: 3, color: 'text.secondary', fontSize: '0.875rem' }}>
								Наследование
							</TableCell>
							<TableCell sx={{ py: 2.5, px: 3, color: 'text.secondary', fontSize: '0.875rem' }}>
								Разрешения
							</TableCell>
							<TableCell sx={{ py: 2.5, px: 3, color: 'text.secondary', fontSize: '0.875rem' }}>
								Пользователей
							</TableCell>
							<TableCell align='right' sx={{ py: 2.5, px: 3, width: 64 }}></TableCell>
						</TableRow>
					</TableHead>

					<TableBody sx={{ '& tr:not(:last-child)': { borderBottom: '1px solid #f3f4f6' } }}>
						{roles.map(role => (
							<TableRow key={role.id} hover sx={{ cursor: 'pointer', '&:hover': { bgcolor: '#fafafa' } }}>
								{/* Роль */}
								<TableCell sx={{ py: 3, px: 4 }}>
									<Box>
										<Typography sx={{ fontWeight: 600, color: '#111827' }}>{role.name}</Typography>
										<Typography sx={{ fontSize: '0.75rem', color: '#9ca3af' }}>
											{role.slug}
										</Typography>
										<Typography sx={{ fontSize: '0.875rem', color: '#6b7280', mt: 0.5 }}>
											{role.desc}
										</Typography>
									</Box>
								</TableCell>

								{/* Статус */}
								<TableCell sx={{ px: 3 }}>
									{role.status === 'active' ? (
										// Активный статус (Зеленый)
										<Box
											sx={{
												display: 'inline-flex',
												alignItems: 'center',
												gap: 1,
												bgcolor: '#ecfdf5',
												color: '#047857',
												px: 2,
												py: 0.75,
												borderRadius: '16px',
												fontSize: '0.875rem',
												fontWeight: 500,
											}}
										>
											<Box
												component='span'
												sx={{ width: 8, height: 8, bgcolor: '#10b981', borderRadius: '50%' }}
											/>
											Активна
										</Box>
									) : (
										// Неактивный статус (Красный/Серый)
										<Box
											sx={{
												display: 'inline-flex',
												alignItems: 'center',
												gap: 1,
												bgcolor: '#fef2f2',
												color: '#b91c1c',
												px: 2,
												py: 0.75,
												borderRadius: '16px',
												fontSize: '0.875rem',
												fontWeight: 500,
											}}
										>
											<Box
												component='span'
												sx={{ width: 8, height: 8, bgcolor: '#ef4444', borderRadius: '50%' }}
											/>
											Неактивна
										</Box>
									)}
								</TableCell>

								{/* Наследование */}
								<TableCell sx={{ px: 3 }}>
									{role.parents && role.parents.length > 0 ? (
										<Box sx={{ display: 'flex', alignItems: 'center', flexWrap: 'wrap', gap: 0.5 }}>
											{role.parents.map((parentName, index) => (
												<Fragment key={parentName}>
													<Typography
														variant='caption'
														sx={{
															bgcolor: '#f3f4f6',
															px: 1,
															py: 0.2,
															borderRadius: '4px',
															color: '#6b7280',
														}}
													>
														{parentName}
													</Typography>
													{/* Если это не последний родитель, можно поставить какой-то разделитель, но обычно в множественном наследовании они равноправны */}
													{index < role.parents.length - 1 && (
														<Typography variant='caption' sx={{ color: '#d1d5db' }}>
															+
														</Typography>
													)}
												</Fragment>
											))}

											<ArrowRightIcon sx={{ fontSize: 14, color: '#d1d5db', mx: 0.5 }} />

											<Typography
												variant='caption'
												sx={{
													fontWeight: 600,
													color: 'primary.main',
													bgcolor: 'rgba(79, 70, 229, 0.08)',
													px: 1,
													py: 0.2,
													borderRadius: '4px',
												}}
											>
												{role.name}
											</Typography>
										</Box>
									) : (
										<Chip
											label='Нет наследования'
											size='small'
											variant='outlined'
											sx={{ borderStyle: 'dashed', color: '#9ca3af', fontSize: '11px' }}
										/>
									)}
								</TableCell>

								{/* Разрешения */}
								<TableCell sx={{ px: 3 }}>
									<Typography sx={{ fontSize: '0.875rem', fontWeight: 500 }}>
										Всего:{' '}
										<Box component='span' sx={{ color: 'primary.main' }}>
											{role.perms.total}
										</Box>
									</Typography>
									<Box sx={{ gap: 1, fontSize: '0.75rem', mt: 0.2 }}>
										<Typography variant='inherit' sx={{ color: '#059669' }}>
											Собственные: {role.perms.own}
										</Typography>
										{role.perms.inherited > 0 && (
											<Typography variant='inherit' sx={{ color: '#2563eb' }}>
												Наследованные: {role.perms.inherited}
											</Typography>
										)}
									</Box>
								</TableCell>

								{/* Кол-во пользователей */}
								<TableCell sx={{ px: 3, fontWeight: 600, color: '#374151' }}>
									{role.userCount}
								</TableCell>

								{/* Действия */}
								<TableCell align='right' sx={{ px: 3 }}>
									<Box sx={{ display: 'flex', gap: 1, justifyContent: 'flex-end' }}>
										<IconButton sx={{ color: '#9ca3af', '&:hover': { color: 'primary.main' } }}>
											<EditBoxIcon sx={{ fontSize: 18 }} />
										</IconButton>
										{/* <IconButton
											size='small'
											sx={{ color: '#9ca3af', '&:hover': { color: 'error.main' } }}
										>
											<DeleteIcon fontSize='small' />
										</IconButton> */}
									</Box>
								</TableCell>
							</TableRow>
						))}
						{!roles.length && !isFetching ? (
							<TableRow>
								<TableCell colSpan={6} align='center' sx={{ py: 3, color: 'text.secondary' }}>
									Роли не найдены.
								</TableCell>
							</TableRow>
						) : null}
					</TableBody>
				</Table>
			</TableContainer>
		</Box>
	)
}
