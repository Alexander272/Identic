import type { FC } from 'react'
import { Box, Divider, Stack, Typography, useTheme } from '@mui/material'

import type { IOrder } from '../../types/order'
import { BusinessIcon } from '@/components/Icons/BusinessIcon'
import { FileIcon } from '@/components/Icons/FileIcon'

type Props = {
	order: IOrder
}

export const Header: FC<Props> = ({ order }) => {
	const { palette } = useTheme()

	const parts = order.manager.split('/')
	const manager = parts.length > 0 && parts[0]
	const assistant = parts.length > 1 && parts[1]

	return (
		<Stack
			direction={{ xs: 'column', md: 'row' }}
			divider={<Divider orientation='vertical' flexItem />}
			spacing={4}
			sx={{ mx: 4, mt: 2 }}
		>
			{/* БЛОК 1: КОНТРАГЕНТЫ */}
			<Box sx={{ flex: 1.25 }}>
				<Stack spacing={2}>
					{order.consumer && (
						<Box>
							<Typography
								variant='caption'
								color='text.secondary'
								sx={{ fontWeight: 700, display: 'block', mb: 0.5 }}
							>
								Конечник
							</Typography>
							<Stack direction='row' spacing={1} alignItems='center'>
								<BusinessIcon fontSize='16px' fill={palette.primary.main} />
								<Typography variant='body1' sx={{ fontWeight: 600 }}>
									{order.consumer}
								</Typography>
							</Stack>
						</Box>
					)}

					{order.customer && (
						<Box>
							<Typography
								variant='caption'
								color='text.secondary'
								sx={{ fontWeight: 700, display: 'block', mb: 0.5 }}
							>
								Заказчик
							</Typography>
							<Typography variant='body2' color='text.primary' fontWeight={'bold'}>
								{order.customer}
							</Typography>
						</Box>
					)}
				</Stack>
			</Box>

			{order.notes && (
				<Box sx={{ flex: 2.25 }}>
					<Typography
						variant='caption'
						color='text.secondary'
						align='center'
						sx={{ fontWeight: 700, display: 'block', mb: 0.5 }}
					>
						Примечание
					</Typography>
					<Typography
						sx={{
							fontSize: '13px',
							color: '#303030',
							bgcolor: '#f9f9f9',
							p: 1.5,
							borderRadius: 2,
							border: '1px solid #eee',
						}}
					>
						{order.notes}
					</Typography>
				</Box>
			)}

			{/* БЛОК 2: ОТВЕТСТВЕННЫЕ */}
			<Box sx={{ flex: 1.25 }}>
				<Typography
					variant='caption'
					color='text.secondary'
					align='center'
					sx={{ fontWeight: 700, display: 'block', mb: 0.5 }}
				>
					Сопровождение
				</Typography>
				<Stack spacing={1} alignItems={'center'}>
					<Stack direction='row' spacing={1} alignItems='center'>
						{/* <UserDataIcon fontSize={28} fill={palette.primary.main} /> */}
						<Typography sx={{ px: 2, py: 0.5, borderRadius: 2, bgcolor: '#e3f2fd' }}>{manager}</Typography>
						<Typography sx={{ px: 2, py: 0.5, borderRadius: 2, bgcolor: '#e6e6e6' }}>
							{assistant}
						</Typography>

						{/* <Typography variant='body2'>
							<strong>Менеджер:</strong> {manager || '—'}
						</Typography> */}
					</Stack>
					{/* {assistant && (
						<Typography variant='body2' sx={{ ml: 4, color: 'text.secondary' }}>
							Помощник: {assistant}
						</Typography>
					)} */}
				</Stack>
			</Box>

			{/* БЛОК 3: ДОКУМЕНТЫ */}
			<Box sx={{ flex: 1.25, textAlign: { md: 'right' } }}>
				<Typography
					variant='caption'
					color='text.secondary'
					sx={{ fontWeight: 700, display: 'block', mb: 0.5 }}
				>
					Счет в 1С
				</Typography>
				<Typography color='primary' display={'flex'} alignItems={'center'} justifyContent={'flex-end'}>
					<FileIcon fontSize={16} fill={palette.primary.main} mr={1} /> Счет № {order.bill || 'Н/Д'}
				</Typography>
			</Box>
		</Stack>
	)
}
