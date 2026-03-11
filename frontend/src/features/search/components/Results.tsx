import { useState, type FC } from 'react'
import {
	Accordion,
	AccordionDetails,
	AccordionSummary,
	Box,
	Button,
	Chip,
	Stack,
	Table,
	TableBody,
	TableCell,
	TableHead,
	TableRow,
	Tooltip,
	Typography,
} from '@mui/material'
import { Link } from 'react-router'

import type { ISearchResults } from '../types/search'
import { BottomArrowIcon } from '@/components/Icons/BottomArrowIcon'
import { DoubleRightIcon } from '@/components/Icons/DoubleRightIcon'
import { CalendarIcon } from '@/components/Icons/CalendarIcon'
import { UnfoldMoreIcon } from '@/components/Icons/UnfoldMoreIcon'
import { UnfoldLessIcon } from '@/components/Icons/UnfoldLessIcon'

type Props = {
	data: ISearchResults[]
}

export const Results: FC<Props> = data => {
	const [expanded, setExpanded] = useState<{ [key: number]: boolean }>({})
	const expendedAll = Object.keys(expanded).length === data.data.length

	const accordionChangeHandler = (year: number) => {
		setExpanded(prevState => ({
			...prevState,
			[year]: !prevState[year],
		}))
	}
	const expendAllHandler = () => {
		if (expendedAll) setExpanded({})
		else setExpanded(data.data.reduce((acc, item) => ({ ...acc, [item.year]: true }), {}))
	}

	if (data?.data?.length === 0) return null
	return (
		<Stack sx={{ height: '100%', maxHeight: 750, overflow: 'auto', mr: -2, pr: 2 }}>
			<Stack direction={'row'} justifyContent={'center'} mb={2}>
				<Typography component='h2' variant='h5' ml={'auto'}>
					Результаты
				</Typography>

				<Tooltip title={expendedAll ? 'Свернуть все' : 'Развернуть все'}>
					<Button
						onClick={expendAllHandler}
						variant='outlined'
						size='small'
						color='inherit'
						sx={{ minWidth: 48, borderColor: '#3f3f3f', ml: 'auto', mr: 1 }}
					>
						{expendedAll ? <UnfoldLessIcon /> : <UnfoldMoreIcon />}
					</Button>
				</Tooltip>
			</Stack>

			{data.data.map((item, index) => (
				<Accordion
					key={item.year}
					defaultExpanded={index === 0}
					expanded={expanded[item.year] || false}
					onChange={() => accordionChangeHandler(item.year)}
					disableGutters
					sx={{
						mb: 1,
						borderRadius: 2,
						border: '1px solid rgba(0,0,0,0.08)',
						boxShadow: '0px 2px 4px rgba(0,0,0,0.05)',
						'&:before': { display: 'none' }, // Убираем стандартный разделитель MUI
						'&.Mui-expanded': {
							boxShadow: '0px 8px 16px rgba(0,0,0,0.1)',
							margin: '16px 0',
						},
					}}
				>
					{/* <AccordionSummary
						expandIcon={<BottomArrowIcon fontSize={16} fill={'#1976d2'} />}
						sx={{
							bgcolor: '#fff',
							borderRadius: 2,
							minHeight: 56,
							'& .MuiAccordionSummary-content': {
								justifyContent: 'space-between', // Распределяем контент
								alignItems: 'center',
							},
						}}
					>
						<Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
							<Box
								sx={{
									width: 40,
									height: 40,
									borderRadius: '50%',
									bgcolor: '#e3f2fd',
									color: '#1976d2',
									display: 'flex',
									alignItems: 'center',
									justifyContent: 'center',
									fontWeight: 'bold',
								}}
							>
								#1024
							</Box>
							<Box>
								<Typography variant='subtitle1' fontWeight='600'>
									Заказы {item.year} года
								</Typography>
								<Typography variant='caption' color='text.secondary'>
									(Найдено: {item.count})
								</Typography>
							</Box>
						</Box>

						{/* <Typography component='span'>
							{item.year} (Найдено: {item.count})
						</Typography> 
					</AccordionSummary> */}

					<AccordionSummary
						expandIcon={<BottomArrowIcon fontSize={16} fill={'#1976d2'} />}
						sx={{
							minHeight: 64,
							px: 3,
							'& .MuiAccordionSummary-content': {
								justifyContent: 'space-between',
								alignItems: 'center',
								width: '100%', // Важно для распределения
							},
						}}
					>
						{/* Левая часть: Название группы */}
						<Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
							<Box
								sx={{
									width: 40,
									height: 40,
									borderRadius: '50%',
									bgcolor: '#e3f2fd',
									display: 'flex',
									alignItems: 'center',
									justifyContent: 'center',
									fontWeight: 'bold',
								}}
							>
								<CalendarIcon fontSize={20} />
							</Box>

							<Box>
								<Typography variant='subtitle1' color='text.primary'>
									Заказы{' '}
									<Typography
										component={'span'}
										fontWeight={'bold'}
										// color='#1976d2'
										fontSize={'1.1rem'}
									>
										{item.year}
									</Typography>{' '}
									года
								</Typography>
							</Box>
							<Chip
								label={`${item.count} поз.`}
								size='small'
								sx={{
									bgcolor: '#f5f5f5',
									fontWeight: 600,
									height: 24,
								}}
							/>
						</Box>
					</AccordionSummary>

					<AccordionDetails sx={{ bgcolor: '#fff', borderRadius: '0 0 12px 12px' }}>
						<Table size='small'>
							<TableHead>
								<TableRow sx={{ bgcolor: '#f7f7f7', borderTopLeftRadius: 2, borderTopRightRadius: 2 }}>
									<TableCell sx={{ borderTopLeftRadius: 8 }}>Заказчик</TableCell>
									<TableCell>Конечник</TableCell>
									<TableCell align='center'>Процент совпадения</TableCell>
									<TableCell align='center'>Совпало позиций</TableCell>
									{/* <TableCell>Ссылка</TableCell> */}
									<TableCell sx={{ borderTopRightRadius: 8 }} />
								</TableRow>
							</TableHead>
							<TableBody>
								{item.orders.map(order => (
									<TableRow key={order.orderId}>
										<TableCell>{order.customer}</TableCell>
										<TableCell>{order.consumer}</TableCell>
										<TableCell align='center'>{order.score}%</TableCell>
										<TableCell align='center'>
											{order.matchedCount}/{order.totalCount}
										</TableCell>
										<TableCell width={120} align='right'>
											<Link to={order.link} target='__blank'>
												<Button
													color='inherit'
													sx={{ textTransform: 'inherit', color: 'black' }}
												>
													Подробнее <DoubleRightIcon fontSize={10} ml={1} />
												</Button>
											</Link>
										</TableCell>
									</TableRow>
								))}
							</TableBody>
						</Table>
					</AccordionDetails>
				</Accordion>
			))}
		</Stack>
	)
}
