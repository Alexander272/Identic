import { useMemo, useState, type FC } from 'react'
import {
	Button,
	Stack,
	Table,
	TableBody,
	TableCell,
	TableContainer,
	TableHead,
	TableRow,
	Tooltip,
	Typography,
	useTheme,
} from '@mui/material'
import { useNavigate } from 'react-router'
import dayjs from 'dayjs'

import type { IFilter } from '../../types/filter'
import { AppRoutes } from '@/pages/router/routes'
import { useGetOrderByIdQuery } from '../../orderApiSlice'
import { numberFormat } from '@/utils/format'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { ModifyIcon } from '@/components/Icons/ModifyIcon'
import { Filter } from './FlatFilters'
import { Header } from './Header'

type Props = {
	id: string
	searchId: string
	// positionIds: string[]
}

export const Order: FC<Props> = ({ id, searchId }) => {
	const navigate = useNavigate()
	const { palette } = useTheme()
	const [filters, setFilters] = useState<IFilter>({
		search: '',
		found: false,
	})

	const { data, isFetching } = useGetOrderByIdQuery({ id, searchId }, { skip: !id })

	const updateFilter = (name: string, value: unknown) => {
		setFilters(prev => ({
			...prev,
			[name]: value,
		}))
	}

	const editHandler = (e: React.MouseEvent) => {
		e.stopPropagation()
		e.preventDefault()

		// navigate(`/orders/edit/${id}`)
		navigate(AppRoutes.EditOrder.replace(':id', id))
	}

	const filteredPositions = useMemo(() => {
		const positions = data?.data?.positions || []

		return positions.filter(item => {
			// 1. Фильтр по поиску
			const matchesSearch = item.name.toLowerCase().includes(filters.search.toLowerCase())

			// 2. Фильтр по найденному
			const matchesFound = filters.found ? item.isFound || false : true

			// // 3. Фильтр по флагу (чекбоксу)
			// const matchesStock = filters.showOnlyAvailable ? item.inStock : true

			return matchesSearch && matchesFound
		})
	}, [data, filters])

	return (
		<Stack position={'relative'} sx={{ mr: -2 }}>
			{isFetching ? <BoxFallback /> : null}

			<Stack direction={'row'} spacing={0.5} width={'100%'} justifyContent={'center'} alignContent={'center'}>
				<Typography component='h2' variant='h6' align='center'>
					Заказ от{' '}
					<Typography component='span' fontWeight={'bold'} variant='h6'>
						{data?.data?.date ? dayjs(data?.data?.date).format('DD.MM.YYYY') : '...'}
					</Typography>
				</Typography>

				<Tooltip title='Редактировать'>
					<Button
						onClick={editHandler}
						// variant='outlined'
						sx={{
							minWidth: 48,
							textTransform: 'inherit',
							// background: '#fff',
							// border: '1px solid #707070',
							borderRadius: '12px',
							padding: '4px 10px',
							':hover': { svg: { fill: palette.primary.main }, color: palette.primary.main },
						}}
					>
						<ModifyIcon sx={{ fontSize: 18 }} />
					</Button>
				</Tooltip>
			</Stack>

			<Stack mb={2} pr={2}>
				{/* <Table size='small'>
					<TableHead>
						<TableRow>
							<TableCell width={'20%'} align='center'>
								Заказчик
							</TableCell>
							<TableCell width={'20%'} align='center'>
								Конечник
							</TableCell>
							<TableCell width={'20%'} align='center'>
								Менеджер / помощник
							</TableCell>
							<TableCell width={'20%'} align='center'>
								Счет в 1С
							</TableCell>
							<TableCell width={'20%'} align='center'>
								Примечание
							</TableCell>
						</TableRow>
					</TableHead>

					<TableBody>
						<TableRow>
							<TableCell align='center'>{data?.data.customer}</TableCell>
							<TableCell align='center'>{data?.data.consumer}</TableCell>
							<TableCell align='center'>{data?.data.manager}</TableCell>
							<TableCell align='center'>{data?.data.bill}</TableCell>
							<TableCell align='center'>{data?.data.notes}</TableCell>
						</TableRow>
					</TableBody>
				</Table> */}

				{data?.data && <Header order={data?.data} />}
			</Stack>

			<Stack direction={'row'} spacing={2} justifyContent={'center'} alignItems={'center'} mb={2} mt={1}>
				{/* <Typography component='h5' fontSize={'1.2rem'}>
					Позиции
				</Typography> */}

				<Filter filters={filters} onChange={updateFilter} showFound={data?.data.posWereFound} />
			</Stack>

			<TableContainer sx={{ height: 570, overflow: 'auto', pr: 2 }}>
				<Table size='small' stickyHeader>
					<TableHead>
						<TableRow>
							<TableCell>№</TableCell>
							<TableCell>Наименование</TableCell>
							<TableCell align='center'>Количество</TableCell>
							<TableCell>Примечание</TableCell>
						</TableRow>
					</TableHead>

					<TableBody>
						{filteredPositions.map((pos, idx) => {
							let bgcolor = idx % 2 === 1 ? '#fafafa' : '#fff'
							if (pos.isFound) bgcolor = palette.rowActive.main

							return (
								<TableRow key={pos.id} hover sx={{ bgcolor: bgcolor }}>
									<TableCell sx={{ borderTopLeftRadius: 8, borderBottomLeftRadius: 8 }}>
										{numberFormat(pos.rowNumber)}
									</TableCell>
									<TableCell>{pos.name}</TableCell>
									<TableCell align='center'>{numberFormat(pos.quantity)}</TableCell>
									<TableCell sx={{ borderTopRightRadius: 8, borderBottomRightRadius: 8 }}>
										{pos.notes}
									</TableCell>
								</TableRow>
							)
						})}
					</TableBody>
				</Table>
			</TableContainer>
		</Stack>
	)
}
