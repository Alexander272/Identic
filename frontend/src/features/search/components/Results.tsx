import type { FC } from 'react'
import { Stack, Typography, Box } from '@mui/material'

import type { IOrderMatchResult, ISearchItem } from '../types/search'
import { ResultsTable } from './ResultsTable'

import NotFoundImage from '@/assets/not-found.png'

type Props = {
	data: IOrderMatchResult[]
	search: ISearchItem[]
	isLoading: boolean | undefined
}

export const Results: FC<Props> = ({ data, search, isLoading }) => {
	return (
		<Stack
			borderRadius={3}
			paddingX={2}
			paddingY={1}
			border={'1px solid rgba(0, 0, 0, 0.12)'}
			sx={{ background: '#fff', height: '100%', flexGrow: 1, maxHeight: 750, overflow: 'auto', pr: 2 }}
		>
			<Stack direction={'row'} justifyContent={'center'} mb={2} position={'relative'}>
				<Typography component='h2' variant='h5'>
					Результаты
				</Typography>
			</Stack>

			{data.length === 0 && !isLoading ? (
				<Stack alignItems={'center'} justifyContent={'center'} height={'100%'}>
					<Box component='img' src={NotFoundImage} alt='not found' width={192} height={192} mb={-2} />
					<Typography align='center' fontSize={'1.3rem'} fontWeight={'bold'}>
						Ничего не найдено
					</Typography>
				</Stack>
			) : null}
			{isLoading && (
				<Stack alignItems={'center'} justifyContent={'center'} height={'100%'}>
					<Typography fontSize={'1.3rem'}>Идет поиск...</Typography>
					<Typography fontSize={'1.1rem'} variant='caption' color='text.secondary'>
						Поиск может занять некоторое время
					</Typography>
				</Stack>
			)}

			{data.length > 0 && !isLoading ? <ResultsTable data={data} search={search} /> : null}
		</Stack>
	)
}
