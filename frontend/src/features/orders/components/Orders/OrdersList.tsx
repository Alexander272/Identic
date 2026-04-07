import { memo, useState, type FC, type SyntheticEvent } from 'react'
import { Button, Stack, Tab, Tabs, Tooltip, Typography, useTheme } from '@mui/material'
import { Link } from 'react-router'

import type { IFilter } from '../../types/params'
import { useGetOrdersByYearQuery, useGetOrdersQuery, useGetUniqueDataQuery } from '../../orderApiSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { AppRoutes } from '@/pages/router/routes'
import { TextDocIcon } from '@/components/Icons/TextDocIcon'
import { Filters } from '../Filters/Filters'
import { OrdersTable } from './OrdersTable'

export const OrdersList: FC = () => {
	const { palette } = useTheme()
	const [year, setYear] = useState(new Date().getFullYear().toString())
	const [filters, setFilters] = useState<IFilter[]>([])

	const { data, isFetching } = useGetOrdersByYearQuery(year, { skip: !year || filters.length > 0 })
	const { data: orders, isFetching: isFetchingOrders } = useGetOrdersQuery(filters, {
		skip: !filters.length,
	})

	const tabHandler = (_event: SyntheticEvent, newValue: string) => {
		setYear(newValue)
	}

	return (
		<>
			{isFetching || isFetchingOrders ? <BoxFallback /> : null}

			<Stack direction={'row'} alignItems={'center'} mt={1} mb={1}>
				<Stack direction={'row'} alignItems={'center'} mx={'auto'}>
					<Typography variant='h6' mr={1}>
						Заказы
					</Typography>

					<Filters filters={filters} onChange={setFilters} />
				</Stack>

				<Link to={AppRoutes.OrdersList} aria-label='roles page'>
					<Tooltip title='Список всех позиций' disableInteractive>
						<Button sx={{ minWidth: 48, ':hover': { svg: { fill: palette.primary.main } } }}>
							<TextDocIcon fill={'#000'} fontSize={22} transition={'0.3s all ease-in-out'} />
						</Button>
					</Tooltip>
				</Link>
			</Stack>

			{filters.length == 0 && <YearTabs value={year} onChange={tabHandler} />}

			<OrdersTable orders={filters.length ? orders?.data || [] : data?.data || []} />
		</>
	)
}

type YearProps = {
	value: string
	onChange: (event: SyntheticEvent, newValue: string) => void
}

const YearTabs: FC<YearProps> = memo(({ value, onChange }) => {
	const { data: years } = useGetUniqueDataQuery({ field: 'year', sort: 'DESC' })

	return (
		<Tabs
			value={value}
			onChange={onChange}
			variant='scrollable'
			scrollButtons
			sx={{
				borderBottom: 0.5,
				borderColor: 'divider',
				mb: 2,
				'.MuiTabs-scrollButtons': { transition: 'all .2s ease-in-out' },
				'.MuiTabs-scrollButtons.Mui-disabled': {
					height: 0,
				},
			}}
		>
			{years?.data.map(t => (
				<Tab
					key={t}
					label={t}
					value={t}
					sx={{
						textTransform: 'inherit',
						borderRadius: 3,
						transition: 'all 0.3s ease-in-out',
						maxWidth: '100%',
						flexGrow: 1,
						fontWeight: 600,
						':hover': {
							backgroundColor: '#f5f5f5',
						},
					}}
				/>
			))}
		</Tabs>
	)
})
