import type { FC } from 'react'
import {
	Stack,
	Table,
	Typography,
	TableBody,
	TableCell,
	TableHead,
	TableRow,
	Box,
	Button,
	CircularProgress,
	circularProgressClasses,
	Tooltip,
	useTheme,
} from '@mui/material'
import { Link } from 'react-router'

import type { IOrderMatchResult, ISearchItem } from '../types/search'
import { PopupLinkIcon } from '@/components/Icons/PopupLinkIcon'
import { Info } from './Info'

type Props = {
	data: IOrderMatchResult[]
	search: ISearchItem[]
}

export const ResultsTable: FC<Props> = ({ data, search }) => {
	const { palette } = useTheme()

	return (
		<Table size='small'>
			<TableHead>
				{/* <TableRow sx={{ bgcolor: '#f7f7f7', borderTopLeftRadius: 2, borderTopRightRadius: 2 }}>
						<TableCell rowSpan={2} align='center' sx={{ borderTopLeftRadius: 8 }}>
							Год
						</TableCell>
						<TableCell rowSpan={2}>Контрагент</TableCell>
						{/* <TableCell>Заказчик</TableCell>
						<TableCell>Конечник</TableCell>
						<TableCell rowSpan={2} align='center'>
							% совпадения
						</TableCell>
						<TableCell colSpan={2} align='center'>
							Совпало
						</TableCell>
						{/* <TableCell>Ссылка</TableCell>
						<TableCell rowSpan={2} sx={{ borderTopRightRadius: 8 }} />
					</TableRow> */}

				{/* Первый ряд шапки */}
				<TableRow sx={{ bgcolor: '#f7f7f7', '& th': { fontWeight: 700 } }}>
					<TableCell rowSpan={2} sx={{ borderTopLeftRadius: 8 }}>
						Год
					</TableCell>
					<TableCell rowSpan={2}>Контрагент</TableCell>
					<TableCell rowSpan={2} align='center'>
						% совпадения
					</TableCell>
					<TableCell colSpan={2} align='center' sx={{ pb: 0 }}>
						Совпало по
					</TableCell>
					<TableCell rowSpan={2} sx={{ borderTopRightRadius: 8 }} />
				</TableRow>

				{/* Второй ряд (под-шапка) */}
				<TableRow
					sx={{
						bgcolor: '#f7f7f7',
						'& th': { fontWeight: 600, color: 'text.secondary', fontSize: '0.75rem' },
					}}
				>
					<TableCell align='center' sx={{ pt: 0.5 }}>
						позициям
					</TableCell>
					<TableCell align='center' sx={{ pt: 0.5 }}>
						кол-ву
					</TableCell>
				</TableRow>
			</TableHead>
			<TableBody>
				{data.map(order => (
					<TableRow key={order.orderId}>
						<TableCell align='center' sx={{ fontWeight: 'bold' }}>
							{order.year}
						</TableCell>
						<TableCell>
							<Typography>Конечник: {order.consumer || '-'}</Typography>
							<Typography variant='body2' color='text.secondary'>
								Заказчик: {order.customer || '-'}
							</Typography>
						</TableCell>
						{/* <TableCell>{order.customer}</TableCell> */}
						{/* <TableCell>{order.consumer}</TableCell> */}
						<TableCell align='center'>
							<Box
								sx={{
									display: 'grid',
									gridTemplateColumns: '1fr 24px',
									alignItems: 'center',
									gap: 1,
								}}
							>
								<Typography fontWeight={'bold'}>{order.score}%</Typography>
								<CircularProgress
									variant='determinate'
									value={order.score}
									size={22}
									thickness={6}
									disableShrink
									enableTrackSlot
									sx={theme => ({
										color: () => {
											if (order.score >= 100) return theme.palette.success.main
											if (order.score >= 70) return '#8bc34a' // Кастомный светло-зеленый
											if (order.score >= 50) return theme.palette.info.main
											if (order.score >= 25) return theme.palette.warning.main
											return theme.palette.error.main
										},
										// Стилизация подложки (трека), если включен enableTrackSlot
										'& .MuiCircularProgress-circle': {
											strokeLinecap: 'round', // Делает концы прогресса закругленными
										},
										[`& .${circularProgressClasses.track}`]: {
											opacity: 1,
											stroke: (theme.vars || theme).palette.grey[300],
										},
									})}
								/>
							</Box>
						</TableCell>
						<TableCell align='center'>
							{order.matchedPos}/{order.totalCount}
						</TableCell>
						<TableCell align='center'>
							{order.matchedQuant}/{order.totalCount}
						</TableCell>
						<TableCell width={96} align='right' sx={{ padding: 0 }}>
							<Stack direction={'row'} justifyContent={'flex-end'} sx={{ height: '100%' }}>
								<Info data={order} search={search} />

								<Link to={order.link} target='_blank' rel='noopener noreferrer'>
									{/* <Button color='inherit' sx={{ textTransform: 'inherit', color: 'black' }}>
										Подробнее <DoubleRightIcon fontSize={10} ml={1} />
									</Button> */}
									<Tooltip title='Перейти к заказу'>
										<Button
											sx={{
												minWidth: 48,
												borderRadius: '6px',
												height: '100%',
												':hover': { svg: { fill: palette.secondary.main } },
											}}
										>
											<PopupLinkIcon fontSize={16} />
										</Button>
									</Tooltip>
								</Link>
							</Stack>
						</TableCell>
					</TableRow>
				))}
			</TableBody>
		</Table>
	)
}
