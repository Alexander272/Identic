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
import { renderManagers } from './RenderManagers'

type Props = {
	year: number
}

const TableComponents = {
	Scroller: forwardRef<HTMLDivElement>((props, ref) => (
		<TableContainer component={Paper} {...props} ref={ref} sx={{ boxShadow: 'none' }} />
	)),
	Table: (props: TableProps) => (
		<Table
			{...props}
			stickyHeader
			size='small'
			sx={{
				borderCollapse: 'separate',
				// На маленьких экранах убираем фиксированную верстку
				tableLayout: { xs: 'auto', md: 'fixed' },
			}}
		/>
	),
	TableHead: forwardRef<HTMLTableSectionElement>((props, ref) => (
		<TableHead
			{...props}
			ref={ref}
			// Скрываем шапку на мобилках
			sx={{ display: { xs: 'none', md: 'table-header-group' } }}
		/>
	)),
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
					display: { xs: 'flex', md: 'table-row' },
					flexDirection: 'column',
					mb: { xs: 2, md: 0 },
					border: { xs: '1px solid #eee', md: 'none' },
					borderRadius: { xs: '8px', md: 0 },
					padding: { xs: 1, md: 0 },
				}}
			/>
		)
	},
	TableBody: forwardRef<HTMLTableSectionElement>((props, ref) => <TableBody {...props} ref={ref} />),
}

const cellStyle = (label: string) => ({
	display: { xs: 'flex', md: 'table-cell' },
	justifyContent: 'space-between',
	alignItems: 'center',
	textAlign: { xs: 'right', md: 'left' },
	width: { xs: '100% !important', md: 'auto' },
	borderBottom: { xs: '1px solid #f0f0f0', md: 'unset' },
	padding: { xs: '8px 4px', md: '6px 16px' },
	// Добавляем текст заголовка перед контентом на мобилках
	'&::before': {
		content: { xs: `"${label}"`, md: 'none' },
		fontWeight: 'bold',
		marginRight: 2,
		textAlign: 'left',
		fontSize: '0.75rem',
		color: 'text.secondary',
	},
})

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
							<TableCell sx={{ fontWeight: 'bold' }}>№</TableCell>
							<TableCell align='center' sx={{ fontWeight: 'bold' }}>
								Дата
							</TableCell>
							<TableCell sx={{ fontWeight: 'bold' }}>Конечник</TableCell>
							<TableCell sx={{ fontWeight: 'bold' }}>Заказчик</TableCell>
							<TableCell align='right' sx={{ fontWeight: 'bold' }}>
								Кол-во позиций
							</TableCell>
							<TableCell sx={{ fontWeight: 'bold' }}>Менеджер / помощник</TableCell>
							<TableCell sx={{ fontWeight: 'bold' }}>Счет в 1С</TableCell>
							<TableCell sx={{ fontWeight: 'bold' }}>Примечание</TableCell>
							<TableCell />
						</TableRow>
					)}
					// onClick={}
					itemContent={(idx, d) => (
						<>
							<TableCell sx={{ ...cellStyle('№'), borderTopLeftRadius: 8 }}>{idx + 1}</TableCell>
							<TableCell sx={cellStyle('Дата')} align='center'>
								{dayjs(d.date).format('DD.MM.YYYY')}
							</TableCell>
							<TableCell sx={cellStyle('Конечник')}>{d.consumer || '—'}</TableCell>
							<TableCell sx={cellStyle('Заказчик')}>{d.customer || '—'}</TableCell>
							<TableCell sx={cellStyle('Кол-во позиций')} align='right'>
								{d.positionCount}
							</TableCell>
							<TableCell sx={cellStyle('Менеджер / помощник')}>{renderManagers(d.manager)}</TableCell>
							<TableCell sx={cellStyle('Счет в 1С')}>{d.bill || '—'}</TableCell>
							<TableCell sx={cellStyle('Примечание')}>{d.notes || '—'}</TableCell>
							<TableCell sx={{ ...cellStyle('Перейти к заказ'), borderTopRightRadius: 8 }}>
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
