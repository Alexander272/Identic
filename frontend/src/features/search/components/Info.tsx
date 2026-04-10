import { useEffect, type FC, type MouseEvent } from 'react'
import { Button, Dialog, DialogContent, DialogTitle, Divider, Stack, useTheme } from '@mui/material'
import dayjs from 'dayjs'

import type { IOrderMatchResult, ISearchItem } from '../types/search'
import { useLazyGetOrderInfoQuery } from '@/features/orders/orderApiSlice'
import { Header } from '@/features/orders/components/Order/Header'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { TimesIcon } from '@/components/Icons/TimesIcon'
import { MatchTable } from './MatchTable/MatchTable'

type Props = {
	data: IOrderMatchResult
	search: ISearchItem[]
	searchId: string
	open: boolean
	onClose: (e: MouseEvent) => void
}

export const Info: FC<Props> = ({ data, search, searchId, open, onClose }) => {
	const { palette } = useTheme()
	// const [open, setOpen] = useState(false)

	const [getOrder, { data: order, isFetching }] = useLazyGetOrderInfoQuery()

	useEffect(() => {
		if (open && !order) {
			getOrder({ id: data.orderId, searchId: searchId })
		}
	}, [open, data.orderId, searchId, getOrder, order])

	const closeHandler = (e: React.MouseEvent, reason?: string) => {
		if (reason === 'backdropClick') return
		console.log('close')

		e.stopPropagation?.()
		onClose(e)
	}

	// const toggleHandler = (e: MouseEvent) => {
	// 	e.stopPropagation()

	// 	if (!open && !order) fetchOrder()
	// 	onClose(e)
	// }

	// const fetchOrder = useCallback(async () => {
	// 	await getOrder({ id: data.orderId, positions: data.positions?.map(p => p.id) })
	// }, [data.orderId, data.positions, getOrder])

	return (
		<>
			<Dialog open={open} fullScreen>
				<Stack direction={'row'} width={'100%'} justifyContent={'space-between'} alignItems={'center'}>
					<DialogTitle sx={{ fontWeight: 'bold' }}>
						Информация о заказе от{' '}
						{order?.data?.date ? dayjs(order?.data?.date).format('DD.MM.YYYY') : '...'}
					</DialogTitle>

					<Button
						onClick={closeHandler}
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
