import { AppBar, Box, Stack, styled, Toolbar, Tooltip, useTheme } from '@mui/material'
import { Link } from 'react-router'

import { AppRoutes } from '@/pages/router/routes'

import Logo from '@/assets/logo.webp'
import { AddFileIcon } from '../Icons/AddFileIcon'
import { SearchIcon } from '../Icons/SearchIcon'

export const LayoutHeader = () => {
	const { palette } = useTheme()

	return (
		<AppBar position='relative' sx={{ borderRadius: 0, alignItems: 'center' }}>
			<Toolbar sx={{ justifyContent: 'space-between', width: '100%', maxWidth: 'xl' }}>
				<Link to='/' aria-label='home page'>
					<Stack
						display={'flex'}
						height={50}
						overflow={'hidden'}
						alignItems={'center'}
						justifyContent={'center'}
						sx={{ img: { height: '100%', width: 'auto' } }}
					>
						<img src={Logo} alt='logo' />
					</Stack>
				</Link>

				<Stack ml={'auto'} direction={'row'} spacing={0.5}>
					<Link to={AppRoutes.CreateOrder} aria-label='roles page'>
						<Tooltip title='Добавить заказ' disableInteractive>
							<NavBox sx={{ ':hover': { svg: { fill: palette.primary.main } } }}>
								<AddFileIcon fill={'#000'} fontSize={26} transition={'0.3s all ease-in-out'} />
							</NavBox>
						</Tooltip>
					</Link>
					{/* <Link to={AppRoutes.OrdersList} aria-label='roles page'>
						<Tooltip title='Список всех позиций' disableInteractive>
							<NavBox sx={{ ':hover': { svg: { fill: palette.primary.main } } }}>
								<TextDocIcon fill={'#000'} fontSize={26} transition={'0.3s all ease-in-out'} />
							</NavBox>
						</Tooltip>
					</Link> */}
					<Link to={AppRoutes.Search} aria-label='roles page'>
						<Tooltip title='Поиск' disableInteractive>
							<NavBox sx={{ ':hover': { svg: { fill: palette.primary.main } } }}>
								<SearchIcon fill={'#000'} fontSize={24} transition={'0.3s all ease-in-out'} />
							</NavBox>
						</Tooltip>
					</Link>
				</Stack>
			</Toolbar>
		</AppBar>
	)
}

const NavBox = styled(Box)(() => ({
	width: 46,
	height: 46,
	display: 'flex',
	justifyContent: 'center',
	alignItems: 'center',
	cursor: 'pointer',
	borderRadius: 12,
	transition: '.3s all ease-in-out',

	':hover': {
		background: '#05287f0a',
	},
}))
