import { useState, type FC, type MouseEvent } from 'react'
import {
	Stack,
	Table,
	Typography,
	TableCell,
	TableRow,
	Box,
	Button,
	CircularProgress,
	circularProgressClasses,
	Tooltip,
	useTheme,
	TableHead,
} from '@mui/material'
import { Link } from 'react-router'
import dayjs from 'dayjs'
import { TableVirtuoso } from 'react-virtuoso'

import type { IOrderMatchResult, ISearchItem } from '../types/search'
import { OrderChip } from '@/features/orders/components/Orders/OrderChip'
import { PopupLinkIcon } from '@/components/Icons/PopupLinkIcon'
import { IndexingPageIcon } from '@/components/Icons/IndexingPageIcon'
import { Info } from './Info'

type Props = {
	data: IOrderMatchResult[]
	search: ISearchItem[]
	searchId: string
}

export const ResultsTable: FC<Props> = ({ data, search, searchId }) => {
	const [openOrderId, setOpenOrderId] = useState<string | null>(null)

	const closeHandler = (event: MouseEvent) => {
		event.stopPropagation()
		setOpenOrderId(null)
	}

	const handleRowClick = (orderId: string) => {
		setOpenOrderId(orderId)
	}

	const openOrder = data.find(order => order.orderId === openOrderId)

	return (
		<>
			<TableVirtuoso
				data={data}
				components={{
					Table: props => <Table {...props} size='small' />,
					TableHead: props => <TableHead {...props} />,
					TableRow: props => {
						const index = (props as unknown as Record<string, unknown>)['data-index'] as number | undefined
						return (
							<TableRow
								{...props}
								hover
								onClick={() => {
									if (typeof index === 'number') {
										handleRowClick(data[index]?.orderId)
									}
								}}
								sx={{ cursor: 'pointer', transition: '0.2s all ease-in-out' }}
							/>
						)
					},
				}}
				fixedHeaderContent={() => (
					<>
						<TableRow sx={{ bgcolor: '#f7f7f7', '& th': { fontWeight: 700 } }}>
							<TableCell rowSpan={2} width={80} align='center' sx={{ borderTopLeftRadius: 8 }}>
								Дата
							</TableCell>
							<TableCell rowSpan={2} width={30} sx={{ p: 0 }} />
							<TableCell rowSpan={2} width={480}>
								Контрагент
							</TableCell>
							<TableCell rowSpan={2} width={130} align='center'>
								% совпадения
							</TableCell>
							<TableCell
								colSpan={2}
								width={160}
								align='center'
								sx={{ pb: 0, borderBottom: '1px solid #e0e0e0' }}
							>
								Совпало по
							</TableCell>
							<TableCell rowSpan={2} width={96} sx={{ borderTopRightRadius: 8 }} />
						</TableRow>
						<TableRow
							sx={{
								bgcolor: '#f7f7f7',
								'& th': { fontWeight: 600, color: 'text.secondary', fontSize: '0.75rem' },
							}}
						>
							<TableCell width={85} align='center' sx={{ pt: 0.5 }}>
								позициям
							</TableCell>
							<TableCell width={75} align='center' sx={{ pt: 0.5 }}>
								кол-ву
							</TableCell>
						</TableRow>
					</>
				)}
				itemContent={(_, order) => <ResultRow order={order} onRowClick={() => handleRowClick(order.orderId)} />}
			/>
			{openOrder && (
				<Info data={openOrder} search={search} searchId={searchId} open={true} onClose={closeHandler} />
			)}
		</>
	)
}

const getScoreColor = (score: number) => {
	if (score >= 100) return 'success.main'
	if (score >= 70) return '#8bc34a'
	if (score >= 50) return 'info.main'
	if (score >= 25) return 'warning.main'
	return 'error.main'
}

const ResultRow: FC<{
	order: IOrderMatchResult
	onRowClick: () => void
}> = ({ order, onRowClick }) => {
	const { palette } = useTheme()
	const scoreColor = getScoreColor(order.score)

	return (
		<>
			<TableCell align='center'>
				<Typography>{dayjs(order.date).format('DD.MM')}</Typography>
				<Typography fontWeight={'bold'} fontSize={'0.9rem'}>
					{order.year}
				</Typography>
			</TableCell>
			<TableCell sx={{ p: 0 }}>
				<Box display={'flex'} gap={0.5} alignItems={'center'} justifyContent={'center'} flexWrap={'wrap'}>
					{order.isBargaining ? <OrderChip type='bargaining' /> : null}
					{order.isBudget ? <OrderChip type='budget' /> : null}
				</Box>
			</TableCell>
			<TableCell>
				<Typography>Конечник: {order.consumer || '-'}</Typography>
				<Typography variant='body2' color='text.secondary'>
					Заказчик / Перекуп: {order.customer || '-'}
				</Typography>
			</TableCell>
			<TableCell align='center'>
				<Box
					sx={{
						display: 'grid',
						gridTemplateColumns: '1fr 24px',
						alignItems: 'center',
						gap: 1,
					}}
				>
					<Typography fontWeight={'bold'}>{Math.round(order.score)}%</Typography>
					<CircularProgress
						variant='determinate'
						value={order.score}
						size={22}
						thickness={6}
						disableShrink
						enableTrackSlot
						sx={theme => ({
							color: scoreColor,
							'& .MuiCircularProgress-circle': {
								strokeLinecap: 'round',
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
					<Tooltip title='Подробнее'>
						<Button
							onClick={e => {
								e.stopPropagation()
								onRowClick()
							}}
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

					<Link
						to={order.link}
						target='_blank'
						rel='noopener noreferrer'
						onClick={(e: MouseEvent) => e.stopPropagation()}
					>
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
		</>
	)
}
