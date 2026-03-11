import { AppBar, Box, Toolbar } from '@mui/material'

import { AppRoutes } from '@/pages/router/routes'

import logo from '@/assets/logo.webp'
import { Link } from 'react-router'

export const LayoutHeader = () => {
	return (
		<AppBar sx={{ borderRadius: 0 }}>
			<Toolbar sx={{ justifyContent: 'space-between', alignItems: 'inherit' }}>
				<Box alignSelf={'center'} display={'flex'} alignItems={'center'} component={Link} to={AppRoutes.Home}>
					<img height={46} width={157} src={logo} alt='logo' />
				</Box>
			</Toolbar>
		</AppBar>
	)
}
