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
import { useNavigate } from 'react-router'
import dayjs from 'dayjs'

import type { IOrder } from '../../types/order'
import { AppRoutes } from '@/pages/router/routes'
import { ModifyIcon } from '@/components/Icons/ModifyIcon'
import { renderManagers } from './RenderManagers'

type Props = {
	orders: IOrder[]
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

export const OrdersTable: FC<Props> = ({ orders }) => {
	const { palette } = useTheme()
	const navigate = useNavigate()

	const editHandler = (id: string) => (e: React.MouseEvent) => {
		e.stopPropagation()
		e.preventDefault()

		navigate(AppRoutes.EditOrder.replace(':id', id))
	}

	return (
		<Stack sx={{ mr: -2 }}>
			<TableContainer>
				<TableVirtuoso
					style={{ height: 700, paddingRight: 16, width: '100%' }}
					data={orders}
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
								Заказчик / Перекуп
							</TableCell>
							<TableCell width={120} align='center' sx={{ fontWeight: 'bold' }}>
								Кол-во позиций
							</TableCell>
							<TableCell width={180} sx={{ fontWeight: 'bold' }}>
								Менеджер / помощник
							</TableCell>
							<TableCell width={120} align='center' sx={{ fontWeight: 'bold' }}>
								Счет в 1С
							</TableCell>
							<TableCell width={280} sx={{ fontWeight: 'bold' }}>
								Примечание
							</TableCell>
							<TableCell width={40} />
						</TableRow>
					)}
					itemContent={(idx, d) => (
						<>
							<TableCell sx={{ borderTopLeftRadius: 8, borderBottomLeftRadius: 8 }}>{idx + 1}</TableCell>
							<TableCell align='center'>{dayjs(d.date).format('DD.MM.YYYY')}</TableCell>
							<TableCell>{d.consumer || '—'}</TableCell>
							<TableCell>{d.customer || '—'}</TableCell>
							<TableCell align='center'>{d.positionCount}</TableCell>
							<TableCell>{renderManagers(d.manager)}</TableCell>
							<TableCell align='center'>{d.bill || '—'}</TableCell>
							<TableCell>{d.notes || '—'}</TableCell>
							<TableCell sx={{ padding: 0, borderTopRightRadius: 8, borderBottomRightRadius: 8 }}>
								<Tooltip title='Редактировать заказ'>
									<Button
										onClick={editHandler(d.id)}
										sx={{
											minWidth: 40,
											minHeight: 38,
											borderRadius: '6px',
											':hover': { svg: { fill: palette.secondary.main } },
										}}
									>
										<ModifyIcon sx={{ fontSize: 18 }} />
									</Button>
								</Tooltip>
							</TableCell>
						</>
					)}
				/>
			</TableContainer>
		</Stack>
	)
}
