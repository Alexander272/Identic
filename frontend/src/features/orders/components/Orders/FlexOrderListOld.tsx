import { type FC } from 'react'
import { useMediaQuery, useTheme, Box, Typography, Paper, Button, Tooltip, Stack } from '@mui/material'
import { Virtuoso } from 'react-virtuoso'
import { Link } from 'react-router'
import dayjs from 'dayjs'
import type { IOrder } from '../../types/order'
import { renderManagers } from './RenderManagers'
import { PopupLinkIcon } from '@/components/Icons/PopupLinkIcon'
import { useGetOrdersByYearQuery } from '../../orderApiSlice'
import { BoxFallback } from '@/components/Fallback/BoxFallback'

const gridTemplateColumns = '40px 110px 1fr 1fr 140px 160px 120px 1fr 60px'

// ─────────────────────────────────────────────────────────────
// Компонент строки: десктоп = ряд сетки, мобильный = карточка
// ─────────────────────────────────────────────────────────────
const MobileRow: FC<{ item: IOrder; idx: number }> = ({ item, idx }) => {
	const { palette } = useTheme()

	// 📱 Мобильный: карточка с вертикальным расположением полей
	return (
		<Paper
			elevation={0}
			sx={{
				p: 2,
				m: 1,
				borderRadius: 2,
				backgroundColor: idx % 2 === 0 ? '#fafafa' : 'white',
				'&:hover': { backgroundColor: '#f0f4f8 !important' },
				border: '1px solid #eee',
			}}
		>
			<Stack spacing={1.5}>
				{/* Заголовок карточки */}
				<Box
					sx={{
						display: 'flex',
						justifyContent: 'space-between',
						alignItems: 'center',
						pb: 1,
						borderBottom: '1px solid #f0f0f0',
					}}
				>
					<Typography variant='subtitle2' fontWeight='bold'>
						Заказ №{idx + 1}
					</Typography>
					<Typography
						variant='caption'
						color='text.secondary'
						sx={{ bgcolor: '#e3f2fd', px: 1, py: 0.5, borderRadius: 1 }}
					>
						{dayjs(item.date).format('DD.MM.YYYY')}
					</Typography>
				</Box>

				{/* Поля карточки: лейбл + значение */}
				{[
					{ label: 'Заказчик', value: item.customer },
					{ label: 'Конечник', value: item.consumer },
					{ label: 'Кол-во позиций', value: item.positionCount?.toString() },
					{ label: 'Менеджер / помощник', value: renderManagers(item.manager) },
					{ label: 'Счет в 1С', value: item.bill },
					{ label: 'Примечание', value: item.notes, multiline: true },
				].map((field, i) => (
					<Box key={i} sx={{ display: 'flex', gap: 1 }}>
						<Typography
							variant='caption'
							color='text.secondary'
							sx={{ minWidth: 140, flexShrink: 0, fontWeight: 500 }}
						>
							{field.label}:
						</Typography>
						<Typography
							variant='body2'
							sx={{
								wordBreak: 'break-word',
								whiteSpace: field.multiline ? 'pre-wrap' : 'nowrap',
								overflow: 'hidden',
								textOverflow: 'ellipsis',
								flex: 1,
							}}
						>
							{field.value || '—'}
						</Typography>
					</Box>
				))}

				{/* Кнопка перехода */}
				<Box sx={{ pt: 1, display: 'flex', justifyContent: 'flex-end' }}>
					<Link
						to={`/orders/${item.id}`}
						target='_blank'
						rel='noopener noreferrer'
						style={{ textDecoration: 'none' }}
					>
						<Tooltip title='Перейти к заказу'>
							<Button
								size='small'
								variant='outlined'
								endIcon={<PopupLinkIcon fontSize={14} />}
								sx={{
									borderRadius: 2,
									textTransform: 'none',
									'&:hover': {
										svg: { fill: palette.secondary.main },
										borderColor: palette.secondary.main,
									},
								}}
							>
								Открыть
							</Button>
						</Tooltip>
					</Link>
				</Box>
			</Stack>
		</Paper>
	)
}

type Props = {
	year: number
}

// ─────────────────────────────────────────────────────────────
// Заголовок таблицы (только для десктопа)
// ─────────────────────────────────────────────────────────────
const DesktopHeader = () => (
	<Box
		sx={{
			display: 'grid',
			alignItems: 'center',
			gridTemplateColumns: gridTemplateColumns,
			gap: 2,
			py: 2,
			px: 2,
			fontWeight: 600,
			backgroundColor: 'background.paper',
			borderBottom: '2px solid #ddd',
			zIndex: 20,
			boxShadow: '0 2px 8px rgba(0,0,0,0.08)',
		}}
	>
		<Typography>№</Typography>
		<Typography textAlign='center'>Дата</Typography>
		<Typography>Конечник</Typography>
		<Typography>Заказчик</Typography>

		<Typography textAlign='right' sx={{ pr: 1 }}>
			Кол-во позиций
		</Typography>

		<Typography sx={{ pl: 1 }}>Менеджер / помощник</Typography>

		<Typography>Счет в 1С</Typography>
		<Typography>Примечание</Typography>
		<Box />
	</Box>
)

const DesktopRow: FC<{ item: IOrder; idx: number }> = ({ item, idx }) => {
	const { palette } = useTheme()

	return (
		<Box
			sx={{
				display: 'grid',
				gridTemplateColumns: gridTemplateColumns,
				gap: 2,
				py: 1.5,
				px: 2,
				alignItems: 'center',
				borderBottom: '1px solid #eee',
				backgroundColor: idx % 2 === 0 ? '#fafafa' : '#fff',
				'&:hover': { backgroundColor: '#f0f4f8' },
				cursor: 'pointer',
				// '&:first-of-type': { borderTopLeftRadius: 8, borderTopRightRadius: 8 },
				// '&:last-of-type': { borderBottomLeftRadius: 8, borderBottomRightRadius: 8 },
			}}
		>
			<Typography variant='body2'>{idx + 1}</Typography>
			<Typography variant='body2' textAlign='center'>
				{dayjs(item.date).format('DD.MM.YYYY')}
			</Typography>
			<Typography variant='body2' sx={{ minWidth: 0, overflow: 'hidden' }}>
				{item.consumer || '—'}
			</Typography>
			<Typography variant='body2' sx={{ minWidth: 0, overflow: 'hidden' }}>
				{item.customer || '—'}
			</Typography>

			{/* 🔥 Кол-во позиций: отступ справа + выравнивание */}
			<Typography variant='body2' textAlign='right' sx={{ pr: 1, fontWeight: 500 }}>
				{item.positionCount}
			</Typography>

			{/* 🔥 Менеджер: отступ слева */}
			<Box sx={{ pl: 1, minWidth: 0, overflow: 'hidden' }}>{renderManagers(item.manager)}</Box>

			<Typography variant='body2'>{item.bill || '—'}</Typography>
			<Typography variant='body2' sx={{ minWidth: 0, overflow: 'hidden' }}>
				{item.notes || '—'}
			</Typography>

			<Box sx={{ display: 'flex', justifyContent: 'flex-end' }}>
				<Link to={`/orders/${item.id}`} target='_blank' rel='noopener noreferrer'>
					<Tooltip title='Перейти к заказу'>
						<Button
							sx={{
								minWidth: 40,
								padding: '4px',
								borderRadius: '6px',
								'&:hover': { svg: { fill: palette.secondary.main } },
							}}
						>
							<PopupLinkIcon fontSize={14} />
						</Button>
					</Tooltip>
				</Link>
			</Box>
		</Box>
	)
}

const HEADER_HEIGHT = 56

// ─────────────────────────────────────────────────────────────
// OrdersList
// ─────────────────────────────────────────────────────────────
export const OrdersList: FC<Props> = ({ year }) => {
	const theme = useTheme()
	const isMobile = useMediaQuery(theme.breakpoints.down(1000))

	const { data, isFetching } = useGetOrdersByYearQuery(year.toString(), { skip: !year })

	if (isFetching) {
		return <BoxFallback />
	}

	const orders = data?.data || []

	return (
		// 🔥 Контейнер с overflow: hidden для работы position: sticky внутри Virtuoso
		<Box
			sx={{
				height: 700,
				width: '100%',
				bgcolor: 'background.default',
				borderRadius: 2,
				overflow: 'hidden',
				border: '1px solid #e0e0e0',
				display: 'flex',
				flexDirection: 'column',
				scrollbarGutter: 'stable',
			}}
		>
			{!isMobile && <DesktopHeader />}
			<Virtuoso
				data={orders}
				totalCount={orders.length}
				// components={{
				// 	Header: !isMobile ? DesktopHeader : undefined,
				// 	Footer: () => (isMobile ? <Box sx={{ height: 16 }} /> : null),
				// }}
				itemContent={(idx, item: IOrder) =>
					isMobile ? <MobileRow item={item} idx={idx} /> : <DesktopRow item={item} idx={idx} />
				}
				style={{ height: isMobile ? '100%' : `calc(100% - ${HEADER_HEIGHT}px)` }}
				// 🔥 Опционально: плавный скролл
				// scrollerRef={ref => {
				// 	if (ref) {
				// 		ref.style.scrollBehavior = 'smooth'
				// 	}
				// }}
			/>
		</Box>
	)
}
