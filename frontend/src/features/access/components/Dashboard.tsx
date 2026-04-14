import { Box, Typography, Grid, Paper, Skeleton } from '@mui/material'

const statCards: { label: string; key: 'users' | 'active' | 'roles'; color: string }[] = [
	{ label: 'Пользователи', key: 'users', color: '#2196f3' },
	{ label: 'Активные', key: 'active', color: '#4caf50' },
	{ label: 'Роли', key: 'roles', color: '#ff9800' },
]

export const Dashboard = () => {
	const isLoading = false

	const stats = {
		users: 0,
		active: 0,
		roles: 0,
	}

	return (
		<Box id='section-dashboard' className='section active' sx={{ p: 3 }}>
			{/* Page Header */}
			<Box className='page-header' sx={{ mb: 4 }}>
				<Typography variant='h4' component='h1' gutterBottom sx={{ fontWeight: 'bold' }}>
					Дашборд
				</Typography>
				<Typography variant='body1' color='text.secondary'>
					Состояние системы доступов
				</Typography>
			</Box>

			{/* Stats Grid */}
			<Grid container spacing={3} sx={{ mb: 4 }}>
				{statCards.map(card => (
					<Grid key={card.key} size={{ xs: 12, sm: 4 }}>
						<Paper sx={{ p: 3, border: '1px solid #eee', borderRadius: 2 }} elevation={0}>
							<Typography variant='caption' color='text.secondary' sx={{ display: 'block', mb: 1 }}>
								{card.label}
							</Typography>

							{/* Если загрузка — показываем скелетон, если нет — значение */}
							{isLoading ? (
								<Skeleton variant='text' width='60%' height={40} />
							) : (
								<Typography variant='h4' sx={{ color: card.color, fontWeight: 'bold' }}>
									{stats?.[card.key] || 0}
								</Typography>
							)}
						</Paper>
					</Grid>
				))}
			</Grid>

			{/* Activity Card */}
			<Paper elevation={0} sx={{ border: '1px solid rgba(0, 0, 0, 0.12)', borderRadius: 2 }}>
				<Box sx={{ p: '20px' }}>
					<Typography variant='subtitle1' sx={{ fontSize: '15px', fontWeight: 'bold', mb: 2 }}>
						Последние действия
					</Typography>
					<Box id='activity-log' sx={{ fontSize: '13px', color: 'text.secondary' }}>
						<Typography variant='body2'>Действий пока нет.</Typography>
					</Box>

					{/* {isLoading ? (
						<Box>
							<Skeleton height={20} sx={{ mb: 1 }} />
							<Skeleton height={20} width='80%' />
						</Box>
					) : (
						<Box sx={{ fontSize: '13px', color: 'text.secondary' }}>
							{stats?.recentActivity?.length > 0 ? (
								stats.recentActivity.map((log, i) => <Typography key={i}>{log.text}</Typography>)
							) : (
								<Typography variant='body2'>Действий пока нет.</Typography>
							)}
						</Box>
					)} */}
				</Box>
			</Paper>
		</Box>
	)
}
