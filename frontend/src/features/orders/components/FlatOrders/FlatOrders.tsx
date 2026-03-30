import { forwardRef, useCallback, useEffect, useRef, useState } from 'react'
import {
	Box,
	CircularProgress,
	Paper,
	Stack,
	Table,
	TableBody,
	TableCell,
	TableContainer,
	TableFooter,
	TableHead,
	TableRow,
	Typography,
	type TableProps,
	type TableRowProps,
} from '@mui/material'
import { TableVirtuoso } from 'react-virtuoso'
import dayjs from 'dayjs'

import type { IFlatOrder } from '../../types/order'
import { numberFormat } from '@/utils/format'
import { useLazyGetFlatOrderQuery } from '../../orderApiSlice'
import { renderManagers } from '../Orders/RenderManagers'

const TableComponents = {
	Scroller: forwardRef<HTMLDivElement>((props, ref) => (
		<TableContainer component={Paper} {...props} ref={ref} sx={{ boxShadow: 'none' }} />
	)),
	Table: (props: TableProps) => (
		<Table {...props} stickyHeader size='small' style={{ borderCollapse: 'separate', tableLayout: 'fixed' }} />
	),
	TableHead: forwardRef<HTMLTableSectionElement>((props, ref) => <TableHead {...props} ref={ref} />),
	TableRow: (props: TableRowProps & { item: IFlatOrder }) => {
		// eslint-disable-next-line @typescript-eslint/no-unused-vars
		const { item: _, ...rest } = props
		return (
			<TableRow
				{...rest}
				sx={{
					'&:nth-of-type(even)': { backgroundColor: '#fafafa' },
					'&:hover': { backgroundColor: '#f0f4f8 !important' },
					cursor: 'pointer',
				}}
			/>
		)
	},
	TableBody: forwardRef<HTMLTableSectionElement>((props, ref) => <TableBody {...props} ref={ref} />),
	TableFooter: forwardRef<HTMLTableSectionElement>((props, ref) => <TableFooter {...props} ref={ref} />),
}

export const FlatOrders = () => {
	const [orders, setOrders] = useState<IFlatOrder[]>([])
	const [cursor, setCursor] = useState<string | null>(null)
	const [hasMore, setHasMore] = useState(true)
	const loadingRef = useRef(false)

	const [fetchOrders, { isFetching }] = useLazyGetFlatOrderQuery()

	const loadOrders = useCallback(
		async (cursorValue: string | null) => {
			if (loadingRef.current) return
			loadingRef.current = true

			try {
				const res = await fetchOrders({ cursor: cursorValue, limit: 20 }).unwrap()
				console.log('Получено строк:', res.data.orders.length)

				setOrders(prev => {
					const existingIds = new Set(prev.map(t => t.id))
					const newItems = res.data.orders.filter(t => !existingIds.has(t.id))
					console.log('Новых строк:', newItems.length)
					return [...prev, ...newItems]
				})

				setCursor(res.data.cursor)
				setHasMore(res.data.hasMore)
			} finally {
				loadingRef.current = false
			}
		},
		[fetchOrders],
	)

	useEffect(() => {
		const init = async () => {
			await loadOrders(null)
		}
		init()
		// eslint-disable-next-line react-hooks/exhaustive-deps
	}, [])

	const handleLoadMore = () => {
		if (!isFetching && !loadingRef.current && hasMore) {
			loadOrders(cursor)
		}
	}

	return (
		<Stack sx={{ mr: -2 }}>
			<TableContainer>
				<TableVirtuoso
					style={{ height: 750, paddingRight: 16, width: '100%' }}
					data={orders}
					endReached={handleLoadMore}
					increaseViewportBy={100}
					// rangeChanged={handleLoadMore}
					components={TableComponents}
					fixedHeaderContent={() => (
						<TableRow>
							<TableCell width={30} sx={{ fontWeight: 'bold' }}>
								№
							</TableCell>
							<TableCell width={120} align='center' sx={{ fontWeight: 'bold' }}>
								Дата
							</TableCell>
							<TableCell width={260} sx={{ fontWeight: 'bold' }}>
								Конечник
							</TableCell>
							<TableCell width={260} sx={{ fontWeight: 'bold' }}>
								Заказчик
							</TableCell>
							<TableCell width={180} sx={{ fontWeight: 'bold' }}>
								Примечание заказа
							</TableCell>
							<TableCell width={380} sx={{ fontWeight: 'bold' }}>
								Название продукции
							</TableCell>
							<TableCell width={100} align='right' sx={{ fontWeight: 'bold' }}>
								Кол-во
							</TableCell>
							<TableCell width={160} sx={{ fontWeight: 'bold' }}>
								Менеджер / помощник
							</TableCell>
							<TableCell width={120} sx={{ fontWeight: 'bold' }}>
								Счет в 1С
							</TableCell>
							<TableCell width={200} sx={{ fontWeight: 'bold' }}>
								Примечание позиции
							</TableCell>
						</TableRow>
					)}
					itemContent={(idx, d) => (
						<>
							<TableCell sx={{ borderTopLeftRadius: 8, borderBottomLeftRadius: 8 }}>{idx + 1}</TableCell>
							<TableCell align='center'>{dayjs(d.date).format('DD.MM.YYYY')}</TableCell>
							<TableCell>{d.consumer || '—'}</TableCell>
							<TableCell>{d.customer || '—'}</TableCell>
							<TableCell>{d.notes || '—'}</TableCell>
							<TableCell>{d.name || '—'}</TableCell>
							<TableCell align='right'>{numberFormat(d.quantity) || '—'}</TableCell>
							<TableCell>{renderManagers(d.manager)}</TableCell>
							<TableCell>{d.bill || '—'}</TableCell>
							<TableCell>{d.positionNotes || '—'}</TableCell>
						</>
					)}
					fixedFooterContent={() => {
						if (!isFetching) return null
						return (
							<TableRow>
								<TableCell colSpan={10} align='center' sx={{ py: 3 }}>
									<Box
										sx={{ display: 'flex', justifyContent: 'center', alignItems: 'center', gap: 2 }}
									>
										<CircularProgress size={20} />
										<Typography variant='body2' color='text.secondary'>
											Загрузка новых строк...
										</Typography>
									</Box>
								</TableCell>
							</TableRow>
						)
					}}
				/>
			</TableContainer>
		</Stack>
	)
}
