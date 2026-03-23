import { useCallback, useState, type FC } from 'react'
import { Button, Dialog, DialogContent, DialogTitle, Divider, Stack, Tooltip, useTheme } from '@mui/material'
import dayjs from 'dayjs'

import type { IOrderMatchResult, ISearchItem } from '../types/search'
import { useLazyGetOrderInfoQuery } from '@/features/orders/orderApiSlice'
import { Header } from '@/features/orders/components/Order/Header'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { IndexingPageIcon } from '@/components/Icons/IndexingPageIcon'
import { TimesIcon } from '@/components/Icons/TimesIcon'
import { MatchTable } from './MatchTable'

type Props = {
	data: IOrderMatchResult
	search: ISearchItem[]
}

export const Info: FC<Props> = ({ data, search }) => {
	const { palette } = useTheme()
	const [open, setOpen] = useState(false)

	const [getOrder, { data: order, isFetching }] = useLazyGetOrderInfoQuery()

	const toggleHandler = () => {
		const isOpen = open
		setOpen(prev => !prev)
		if (!isOpen && !order) fetchOrder()
	}

	const fetchOrder = useCallback(async () => {
		await getOrder({ id: data.orderId, positions: data.positions?.map(p => p.id) })
	}, [data.orderId, data.positions, getOrder])

	return (
		<>
			<Tooltip title='Подробнее'>
				<Button
					onClick={toggleHandler}
					sx={{
						minWidth: 48,
						minHeight: 48,
						borderRadius: '6px',
						height: '100%',
						':hover': { svg: { fill: palette.secondary.main } },
					}}
				>
					<IndexingPageIcon fontSize={22} />
				</Button>
			</Tooltip>

			<Dialog open={open} onClose={toggleHandler} fullScreen>
				<Stack direction={'row'} width={'100%'} justifyContent={'space-between'} alignItems={'center'}>
					<DialogTitle sx={{ fontWeight: 'bold' }}>
						Информация о заказе от{' '}
						{order?.data?.date ? dayjs(order?.data?.date).format('DD.MM.YYYY') : '...'}
					</DialogTitle>

					<Button
						onClick={toggleHandler}
						sx={{
							minWidth: 40,
							minHeight: 40,
							borderRadius: '12px',
							mr: 2,
							zIndex: 50,
							':hover': { svg: { fill: palette.secondary.main } },
						}}
					>
						<TimesIcon fontSize={14} />
					</Button>
				</Stack>

				<DialogContent sx={{ mt: -4, position: 'relative' }}>
					{/* <Typography variant='h6' fontWeight='bold'>
						Основная информация о заказе от{' '}
						<Typography component={'span'} fontWeight={'bold'}>
							{order?.data?.date ? dayjs(order?.data?.date).format('DD.MM.YYYY') : '...'}
						</Typography>
					</Typography> */}
					{order?.data && <Header order={order?.data} />}

					{isFetching && <BoxFallback />}

					<Divider sx={{ my: 2 }} />

					<MatchTable request={search} result={data} foundPositions={order?.data.positions || []} />
					{/* <TableContainer component={Paper} sx={{ boxShadow: 3, borderRadius: 2 }}>
						<Table sx={{ minWidth: 650 }}>
							<TableHead sx={{ bgcolor: '#f5f5f5' }}>
								<TableRow>
									<TableCell width={50}>#</TableCell>
									<TableCell>Ваш запрос</TableCell>
									<TableCell>Статус</TableCell>
									<TableCell>ID в системе</TableCell>
								</TableRow>
							</TableHead>
							<TableBody>
								{search.map((item, index) => {
									const status = getRowStatus(index)
									const isFound = status === 'found'
									const matchedPos = data.positions.find(p => p.reqId === index.toString())

									return (
										<TableRow
											key={index}
											sx={{
												backgroundColor: isFound
													? 'rgba(76, 175, 80, 0.04)'
													: 'rgba(244, 67, 54, 0.04)',
												borderLeft: `5px solid ${isFound ? '#4caf50' : '#f44336'}`,
											}}
										>
											<TableCell>{index}</TableCell>
											<TableCell>
												<Typography variant='body1' fontWeight='bold'>
													{item.name || 'Без названия'}
												</Typography>
												<Typography variant='caption' color='textSecondary'>
													Требуется: {item.quantity ?? 0} шт.
												</Typography>
											</TableCell>
											<TableCell>
												{isFound ? (
													<Chip
														// icon={<CheckCircleIcon />}
														label='Найдено'
														color='success'
														variant='outlined'
														size='small'
													/>
												) : (
													<Chip
														// icon={<ErrorIcon />}
														label='Отсутствует'
														color='error'
														variant='outlined'
														size='small'
													/>
												)}
											</TableCell>
											<TableCell>
												<Typography sx={{ fontFamily: 'monospace' }}>
													{matchedPos ? matchedPos.id : '—'}
												</Typography>
											</TableCell>
										</TableRow>
									)
								})}
							</TableBody>
						</Table>
					</TableContainer> */}
				</DialogContent>
			</Dialog>
		</>
	)
}
