import { forwardRef, type FC } from 'react'
import {
	Button,
	Paper,
	Stack,
	Table,
	TableBody,
	TableCell,
	TableContainer,
	TableHead,
	TableRow,
	Tooltip,
	useTheme,
	type TableProps,
	type TableRowProps,
} from '@mui/material'
import { TableVirtuoso } from 'react-virtuoso'
import { Link } from 'react-router'
import dayjs from 'dayjs'

import type { IOrder } from '../../types/order'
import { useGetOrdersByYearQuery } from '../../orderApiSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { PopupLinkIcon } from '@/components/Icons/PopupLinkIcon'
import { ManagerChip } from './ManagerChip'

type Props = {
	year: number
}

const TableComponents = {
	Scroller: forwardRef<HTMLDivElement>((props, ref) => (
		<TableContainer component={Paper} {...props} ref={ref} sx={{ boxShadow: 'none' }} />
	)),
	Table: (props: TableProps) => (
		<Table {...props} stickyHeader size='small' style={{ borderCollapse: 'separate', tableLayout: 'fixed' }} />
	),
	TableHead: forwardRef<HTMLTableSectionElement>((props, ref) => <TableHead {...props} ref={ref} />),
	TableRow: (props: TableRowProps & { item: IOrder }) => {
		const { item, ...rest } = props
		return (
			<TableRow
				{...rest}
				onClick={() => {
					const url = `/orders/${item.id}`
					window.open(url, '_blank', 'noopener,noreferrer')
				}}
				sx={{
					'&:nth-of-type(even)': { backgroundColor: '#fafafa' },
					'&:hover': { backgroundColor: '#f0f4f8 !important' },
					cursor: 'pointer',
				}}
			/>
		)
	},
	TableBody: forwardRef<HTMLTableSectionElement>((props, ref) => <TableBody {...props} ref={ref} />),
}

export const OrdersList: FC<Props> = ({ year }) => {
	const { palette } = useTheme()

	const { data, isFetching } = useGetOrdersByYearQuery(year.toString(), { skip: !year })

	return (
		<Stack sx={{ mr: -2 }}>
			<TableContainer>
				{isFetching ? <BoxFallback /> : null}

				<TableVirtuoso
					style={{ height: 700, paddingRight: 16, width: '100%' }}
					data={data?.data || []}
					components={TableComponents}
					fixedHeaderContent={() => (
						<TableRow>
							<TableCell width={60} sx={{ fontWeight: 'bold' }}>
								№
							</TableCell>
							<TableCell width={130} align='center' sx={{ fontWeight: 'bold' }}>
								Дата
							</TableCell>
							<TableCell width={400} sx={{ fontWeight: 'bold' }}>
								Конечник
							</TableCell>
							<TableCell width={400} sx={{ fontWeight: 'bold' }}>
								Заказчик
							</TableCell>
							<TableCell width={140} align='right' sx={{ fontWeight: 'bold' }}>
								Кол-во позиций
							</TableCell>
							<TableCell width={160} sx={{ fontWeight: 'bold' }}>
								Менеджер / помощник
							</TableCell>
							<TableCell width={120} sx={{ fontWeight: 'bold' }}>
								Счет в 1С
							</TableCell>
							<TableCell width={280} sx={{ fontWeight: 'bold' }}>
								Примечание
							</TableCell>
							<TableCell width={60} />
						</TableRow>
					)}
					itemContent={(idx, d) => (
						<>
							<TableCell sx={{ borderTopLeftRadius: 8, borderBottomLeftRadius: 8 }}>{idx + 1}</TableCell>
							<TableCell align='center'>{dayjs(d.date).format('DD.MM.YYYY')}</TableCell>
							<TableCell>{d.consumer || '—'}</TableCell>
							<TableCell>{d.customer || '—'}</TableCell>
							<TableCell align='right'>{d.positionCount}</TableCell>
							<TableCell>{renderManagers(d.manager)}</TableCell>
							<TableCell>{d.bill || '—'}</TableCell>
							<TableCell>{d.notes || '—'}</TableCell>
							<TableCell sx={{ borderTopRightRadius: 8, borderBottomRightRadius: 8 }}>
								<Link to={`/orders/${d.id}`} target='_blank' rel='noopener noreferrer'>
									<Tooltip title='Перейти к заказу'>
										<Button
											sx={{
												minWidth: 48,
												borderRadius: '6px',
												':hover': { svg: { fill: palette.secondary.main } },
											}}
										>
											<PopupLinkIcon fontSize={14} />
										</Button>
									</Tooltip>
								</Link>
							</TableCell>
						</>
					)}
				/>
			</TableContainer>
		</Stack>
	)
}

const renderManagers = (fullString: string) => {
	if (!fullString) return '—'

	// Разделяем по слэшу, запятой или пробелу
	const names = fullString
		.split(/[/,]/)
		.map(n => n.trim())
		.sort()

	return (
		<div style={{ display: 'flex', flexWrap: 'wrap', gap: '4px' }}>
			{names.map(name => (
				<ManagerChip key={name} name={name} />
			))}
		</div>
	)
}
