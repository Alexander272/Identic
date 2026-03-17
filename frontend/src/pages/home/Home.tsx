import { memo, useState, type FC, type SyntheticEvent } from 'react'
import { Box, Tab, Tabs, Typography } from '@mui/material'

import { useGetUniqueDataQuery } from '@/features/orders/orderApiSlice'
import { OrdersList } from '@/features/orders/components/Orders/OrdersList'
import { PageBox } from '@/components/PageBox/PageBox'
import { BoxFallback } from '@/components/Fallback/BoxFallback'

export default function Home() {
	// const [searchParams, setSearchParams] = useSearchParams()

	// const year = searchParams.get('year') || new Date().getFullYear().toString()
	const [year, setYear] = useState(new Date().getFullYear().toString())

	const { data: years, isFetching } = useGetUniqueDataQuery({ field: 'year', sort: 'DESC' }, { skip: !year })

	const tabHandler = (_event: SyntheticEvent, newValue: string) => {
		setYear(newValue)
	}

	// const tabHandler = (_event: SyntheticEvent, value: string) => {
	// 	setSearchParams({ year: value })
	// }

	return (
		<PageBox>
			<Box
				borderRadius={3}
				paddingX={2}
				paddingY={1}
				width={'100%'}
				border={'1px solid rgba(0, 0, 0, 0.12)'}
				// flexGrow={1}
				height={'fit-content'}
				minHeight={600}
				display={'flex'}
				flexDirection={'column'}
				sx={{ backgroundColor: '#fff', userSelect: 'none' }}
			>
				{isFetching && <BoxFallback />}

				<Typography variant='h6' align='center' mt={1} mb={1}>
					Заказы по годам
				</Typography>
				<YearTabs value={year} onChange={tabHandler} years={years?.data || []} />

				<OrdersList year={Number(year)} />
			</Box>
		</PageBox>
	)
}

type YearProps = {
	value: string
	onChange: (event: SyntheticEvent, newValue: string) => void
	years: string[]
}

const YearTabs: FC<YearProps> = memo(({ value, onChange, years }) => {
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
			{years.map(t => (
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
