import { Box, Breadcrumbs } from '@mui/material'

import { PageBox } from '@/components/PageBox/PageBox'
import { Search } from '@/features/search/components/Search'
import { Breadcrumb } from '@/components/Breadcrumb/Breadcrumb'
import { AppRoutes } from '../router/routes'

export default function Home() {
	return (
		<PageBox>
			<Box
				borderRadius={3}
				paddingX={2}
				paddingY={1}
				width={'100%'}
				border={'1px solid rgba(0, 0, 0, 0.12)'}
				mb={1}
				sx={{ backgroundColor: '#fff' }}
			>
				<Breadcrumbs aria-label='breadcrumb'>
					<Breadcrumb to={AppRoutes.Home}>Главная</Breadcrumb>
					<Breadcrumb to={AppRoutes.Search} active>
						Поиск
					</Breadcrumb>
				</Breadcrumbs>
			</Box>

			<Search />
		</PageBox>
	)
}
