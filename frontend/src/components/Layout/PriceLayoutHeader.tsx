import { AppBar, Divider, Stack, Toolbar, Tooltip, Typography, useTheme } from '@mui/material'
import { Link } from 'react-router'
import transportationCosts from '@/assets/Приложение 3.1 - ТР.pdf'
import Logo from '@/assets/logo.webp'

import { useAppSelector } from '@/hooks/redux'
import { useSignOutMutation } from '@/features/auth/authApiSlice'
import { getToken } from '@/features/user/userSlice'
import { LogoutIcon } from '@/components/Icons/LogoutIcon'
import { DeliveryIcon } from '@/components/Icons/DeliveryIcon'
import { NavBox } from './NavBox'

export const PriceLayoutHeader = () => {
	const { palette } = useTheme()
	const token = useAppSelector(getToken)

	const [signOut] = useSignOutMutation()

	const logoutHandler = () => {
		void signOut(null)
	}

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

				<Stack
					direction='row'
					alignItems='center'
					spacing={0.5}
					sx={{
						fontWeight: 700,
						fontSize: 18,
						cursor: 'pointer',
						px: 1.5,
						py: 0.5,
						mr: 4,
						borderRadius: 2,
						transition: '.3s all ease-in-out',
						'&:hover': { bgcolor: 'rgba(0,139,150,0.08)' },
					}}
				>
					<Typography sx={{ color: '#042e60', fontWeight: 700, fontSize: 24 }}>
						Книга цен СИБУР - 2026{' '}
						<Typography component='span' sx={{ fontSize: 14 }}>
							(прайс действует по 31.05.2027)
						</Typography>
					</Typography>
					{/* <BottomArrowIcon sx={{ fontSize: 14, fill: '#008b96' }} /> */}
				</Stack>

				{token ? (
					<Stack
						direction={'row'}
						alignItems={'center'}
						gap={2}
						divider={<Divider orientation='vertical' flexItem variant='middle' />}
					>
						<Link to={transportationCosts} target='_blank' rel='noopener noreferrer'>
							<Tooltip title='Транспортные расходы' disableInteractive>
								<NavBox sx={{ ':hover': { svg: { stroke: palette.primary.main } } }}>
									<DeliveryIcon sx={{ fontSize: 26, transition: '0.3s all ease-in-out' }} />
								</NavBox>
							</Tooltip>
						</Link>

						<NavBox onClick={logoutHandler} sx={{ ':hover': { svg: { fill: palette.primary.main } } }}>
							<LogoutIcon fill={'#000'} fontSize={24} transition={'0.3s all ease-in-out'} />
						</NavBox>
					</Stack>
				) : null}
			</Toolbar>
		</AppBar>
	)
}
