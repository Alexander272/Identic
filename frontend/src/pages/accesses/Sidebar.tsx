import { List, ListItem, ListItemButton, ListItemIcon, ListItemText, Typography, Box } from '@mui/material'
import { useLocation, useNavigate } from 'react-router'

import { AppRoutes } from '../router/routes'
import { DashboardIcon } from '@/components/Icons/DashboardIcon'
import { AccessHandleIcon } from '@/components/Icons/AccessHandleIcon'
import { ShieldLockIcon } from '@/components/Icons/ShieldLockIcon'
import { ShieldIcon } from '@/components/Icons/ShieldIcon'
import { UserIcon } from '@/components/Icons/UserIcon'

const menuItems = [
	{ path: AppRoutes.Accesses, label: 'Дашборд', icon: <DashboardIcon sx={{ fontSize: 18 }} /> },
	{
		path: AppRoutes.UserAccess,
		label: 'Пользователи',
		icon: <UserIcon sx={{ fill: '#000', fontSize: 22 }} />,
	},
	{ path: AppRoutes.RoleAccess, label: 'Роли', icon: <ShieldLockIcon sx={{ fontSize: 22 }} /> },
	{ path: AppRoutes.Permissions, label: 'Права доступа', icon: <AccessHandleIcon fontSize='small' /> },
]

const Sidebar = () => {
	const location = useLocation()
	const navigate = useNavigate()

	const handleSwitch = (path: string) => {
		navigate(path)
	}

	return (
		<Box
			component='aside'
			borderRadius={3}
			paddingX={2}
			paddingY={1}
			border={'1px solid rgba(0, 0, 0, 0.12)'}
			flexGrow={1}
			minHeight={600}
			display={'flex'}
			flexDirection={'column'}
			maxWidth={240}
			sx={{
				// Чтобы сайдбар не сжимался в flex-контейнере
				flexShrink: 0,
				backgroundColor: '#fff',
			}}
		>
			{/* Sidebar Logo */}
			<Box sx={{ p: 3, display: 'flex', alignItems: 'center', gap: 1.5, fontWeight: 'bold' }}>
				<ShieldIcon sx={{ color: 'primary.main' }} />
				<Typography variant='h6' sx={{ fontSize: '1.1rem', fontWeight: 700 }}>
					AccessControl
				</Typography>
			</Box>

			{/* Navigation Items */}
			<List sx={{ flexGrow: 1 }}>
				{menuItems.map(item => (
					<ListItem key={item.path} disablePadding sx={{ mb: 0.5 }}>
						<ListItemButton
							selected={location.pathname === item.path}
							onClick={() => handleSwitch(item.path)}
							sx={{
								borderRadius: '8px',
								'&.Mui-selected': {
									backgroundColor: 'rgba(25, 118, 210, 0.08)',
									color: 'primary.main',
									'& .MuiListItemIcon-root': { color: 'primary.main' },
								},
							}}
						>
							<ListItemIcon sx={{ minWidth: 40 }}>{item.icon}</ListItemIcon>
							<ListItemText
								primary={item.label}
								slotProps={{
									primary: { fontSize: '14px', fontWeight: 500 },
								}}
							/>
						</ListItemButton>
					</ListItem>
				))}
			</List>
		</Box>
	)
}

export default Sidebar
