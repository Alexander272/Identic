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
			<Box sx={{ flex: 1.6 }}>
				<Stack spacing={3} direction={'row'}>
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
								{order.consumer || '-'}
							</Typography>
						</Stack>
					</Box>

					<Box>
						<Typography
							variant='caption'
							color='text.secondary'
							sx={{ fontWeight: 700, display: 'block', mb: 0.5 }}
						>
							Заказчик
						</Typography>
						<Typography variant='body2' color='text.primary' fontWeight={'bold'}>
							{order.customer || '-'}
						</Typography>
					</Box>
				</Stack>
			</Box>

			<Box sx={{ flex: 2.4 }}>
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
						fontSize: '14px',
						color: '#303030',
						bgcolor: '#f9f9f9',
						py: 0.5,
						px: 1.5,
						borderRadius: 2,
						border: '1px solid #eee',
					}}
				>
					{order.notes || '-'}
				</Typography>
			</Box>

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
						{manager || assistant ? (
							<>
								<Typography sx={{ px: 2, py: 0.5, borderRadius: 2, bgcolor: '#e3f2fd' }}>
									{manager}
								</Typography>
								<Typography sx={{ px: 2, py: 0.5, borderRadius: 2, bgcolor: '#e6e6e6' }}>
									{assistant}
								</Typography>
							</>
						) : (
							<Typography fontSize={14}>Менеджер не указан</Typography>
						)}
					</Stack>
				</Stack>
			</Box>

			<Box sx={{ flex: 0.75, textAlign: { md: 'right' } }}>
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
