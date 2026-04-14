import { useMemo, useState, type FC } from 'react'
import {
	Box,
	Typography,
	Button,
	TextField,
	MenuItem,
	Select,
	Table,
	TableBody,
	TableCell,
	TableContainer,
	TableHead,
	TableRow,
	Paper,
	InputAdornment,
	type SelectChangeEvent,
	useTheme,
	Avatar,
	Chip,
} from '@mui/material'
import dayjs from 'dayjs'

import type { IUserData } from '@/features/user/types/user'
import { getAvatarColor, getInitials } from './utils'
import { stringToHSLA } from '@/utils/colors'
import { getSmartDate } from '@/utils/date'
import { SearchIcon } from '@/components/Icons/SearchIcon'
import { PlusIcon } from '@/components/Icons/PlusIcon'
import { useGetAllUsersQuery } from '@/features/user/usersApiSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'

export const Users = () => {
	const { palette } = useTheme()
	const [search, setSearch] = useState('')
	const [roleFilter, setRoleFilter] = useState<string[]>([''])
	const [statusFilter, setStatusFilter] = useState('')

	const { data, isFetching } = useGetAllUsersQuery(null)

	const roles = [
		{ id: '1', slug: 'admin', name: 'Администратор' },
		{ id: '2', slug: 'manager', name: 'Менеджер' },
		{ id: '3', slug: 'reader', name: 'Наблюдатель' },
	]

	// const filteredUsers = useMemo(() => {
	// 	// Здесь будет логика users.filter(...) на основе search, roleFilter и statusFilter
	// 	return []
	// }, [search, roleFilter, statusFilter])

	const roleHandler = (event: SelectChangeEvent<string[]>) => {
		const value = event.target.value
		let newValue = typeof value === 'string' ? value.split(',') : value

		// 1. Если список стал пустым — возвращаем ''
		if (newValue.length === 0) {
			newValue = ['']
		}
		// 2. Если в списке больше одного элемента
		else if (newValue.length > 1) {
			// Если только что добавили что-то к '', то убираем ''
			if (newValue.includes('')) {
				newValue = newValue.filter(v => v !== '')
			}

			// 3. Если выбраны все доступные опции (кроме ''),
			// тут можно добавить условие сравнения с длиной исходного массива ролей
			if (newValue.length === roles.length) {
				newValue = ['']
			}
		}

		setRoleFilter(newValue)
	}

	return (
		<Box sx={{ p: 3 }}>
			{/* Page Header */}
			<Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 4 }}>
				<Box>
					<Typography variant='h4' sx={{ fontWeight: 'bold' }}>
						Пользователи
					</Typography>
					<Typography variant='body1' color='text.secondary'>
						Управление учётными записями
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

			{isFetching && <BoxFallback />}

			{/* Toolbar */}
			<Box sx={{ display: 'flex', gap: 2, mb: 3, flexWrap: 'wrap' }}>
				<TextField
					placeholder='Поиск по имени или email…'
					size='small'
					value={search}
					onChange={e => setSearch(e.target.value)}
					slotProps={{
						input: {
							startAdornment: (
								<InputAdornment position='start'>
									<SearchIcon fontSize='small' />
								</InputAdornment>
							),
						},
					}}
					sx={{ flexGrow: 1, minWidth: '200px', background: '#fff' }}
				/>

				<Select
					size='small'
					displayEmpty
					multiple
					value={roleFilter}
					onChange={roleHandler}
					sx={{ width: '400px', background: '#fff' }}
				>
					<MenuItem value=''>Все роли</MenuItem>
					{roles.map(role => (
						<MenuItem key={role.id} value={role.slug}>
							{role.name}
						</MenuItem>
					))}
				</Select>

				<Select
					size='small'
					displayEmpty
					value={statusFilter}
					onChange={e => setStatusFilter(e.target.value)}
					sx={{ width: '300px', background: '#fff' }}
				>
					<MenuItem value=''>Все статусы</MenuItem>
					<MenuItem value='active'>Активные</MenuItem>
					<MenuItem value='inactive'>Неактивные</MenuItem>
				</Select>
			</Box>

			{/* Table Container */}
			<TableContainer component={Paper} elevation={0} sx={{ border: '1px solid #eee', borderRadius: 2 }}>
				<Table>
					<TableHead>
						<TableRow sx={{ borderBottom: '1px solid #f3f4f6' }}>
							<TableCell sx={{ py: 2.5, px: 4, color: 'text.secondary', fontSize: '0.875rem' }}>
								Пользователь
							</TableCell>
							<TableCell
								sx={{ py: 2.5, px: 4, color: 'text.secondary', fontSize: '0.875rem', width: 250 }}
							>
								Роль
							</TableCell>
							<TableCell
								sx={{ py: 2.5, px: 4, color: 'text.secondary', fontSize: '0.875rem', width: 200 }}
							>
								Статус
							</TableCell>
							<TableCell
								sx={{ py: 2.5, px: 4, color: 'text.secondary', fontSize: '0.875rem', width: 250 }}
							>
								Создан
							</TableCell>
							<TableCell
								sx={{ py: 2.5, px: 4, color: 'text.secondary', fontSize: '0.875rem', width: 250 }}
							>
								Последний вход
							</TableCell>
							<TableCell
								sx={{ py: 2.5, px: 4, color: 'text.secondary', fontSize: '0.875rem', width: 100 }}
							>
								Действия
							</TableCell>
						</TableRow>
					</TableHead>
					<TableBody>
						{data?.data.map(user => (
							<UserRow key={user.id} u={user} />
						))}
						{!data?.data.length && !isFetching ? (
							<TableRow>
								<TableCell colSpan={6} align='center' sx={{ py: 3, color: 'text.secondary' }}>
									Пользователи не найдены.
								</TableCell>
							</TableRow>
						) : null}
					</TableBody>
				</Table>
			</TableContainer>
		</Box>
	)
}

const UserRow: FC<{ u: IUserData }> = ({ u }) => {
	const colors = useMemo(() => stringToHSLA(u.role), [u.role])

	return (
		<TableRow key={u.id} hover>
			{/* Пользователь */}
			<TableCell>
				<Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
					<Avatar
						sx={{
							bgcolor: getAvatarColor(u.id), // ваша функция
							fontSize: '14px',
							width: 36,
							height: 36,
						}}
					>
						{getInitials(u)} {/* ваша функция */}
					</Avatar>
					<Box>
						<Typography variant='body2' sx={{ fontWeight: 500 }}>
							{u.firstName} {u.lastName}
						</Typography>
						<Typography variant='caption' color='text.secondary' sx={{ display: 'block' }}>
							{u.email}
						</Typography>
					</Box>
				</Box>
			</TableCell>

			{/* Роль */}
			<TableCell>
				<Chip
					label={u.role}
					size={'small'}
					style={{
						backgroundColor: colors.bg,
						color: colors.text,
						border: `1px solid ${colors.border}`,
						fontWeight: 500,
						fontSize: '0.75rem',
						height: '20px',
						borderRadius: '6px',
					}}
				/>
			</TableCell>

			{/* Статус */}
			<TableCell>
				<Chip
					label={u.isActive ? 'Активный' : 'Неактивный'}
					size='small'
					variant='outlined'
					color={u.isActive ? 'success' : 'default'}
					sx={{ borderRadius: '6px', fontWeight: 500 }}
				/>
			</TableCell>

			{/* Создан */}
			<TableCell sx={{ color: 'text.secondary', fontSize: '13px' }}>
				{dayjs(u.createdAt).format('ddd, DD MMM, YYYY HH:mm')}
			</TableCell>

			{/* Последний вход */}
			<TableCell sx={{ color: 'text.secondary', fontSize: '13px' }}>{getSmartDate(u.lastVisit)}</TableCell>

			{/* Действия */}
			<TableCell>
				<Box sx={{ display: 'flex', gap: 0.5 }}>
					{/* <IconButton size='small' onClick={() => editUser(u.id)} title='Редактировать'>
						<EditIcon fontSize='small' />
					</IconButton>
					<IconButton
						size='small'
						onClick={() => deleteUser(u.id)}
						title='Удалить'
						sx={{ color: 'error.main' }}
					>
						<DeleteIcon fontSize='small' />
					</IconButton> */}
				</Box>
			</TableCell>
		</TableRow>
	)
}
