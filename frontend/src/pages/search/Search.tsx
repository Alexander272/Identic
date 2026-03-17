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
				// width={'80%'}
				alignSelf={'center'}
				border={'1px solid rgba(0, 0, 0, 0.12)'}
				flexGrow={1}
				height={'fit-content'}
				minHeight={600}
				display={'flex'}
				flexDirection={'column'}
				sx={{ backgroundColor: '#fff', userSelect: 'none' }}
			>
				<Breadcrumbs aria-label='breadcrumb' sx={{ mb: -1 }}>
					<Breadcrumb to={AppRoutes.Home}>Главная</Breadcrumb>
					<Breadcrumb to={AppRoutes.Search} active>
						Поиск
					</Breadcrumb>
				</Breadcrumbs>

				<Search />
			</Box>
		</PageBox>
	)
}
