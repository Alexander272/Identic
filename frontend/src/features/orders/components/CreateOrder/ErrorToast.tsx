import type { FC } from 'react'
import { Box, Stack, Typography, useTheme } from '@mui/material'
import type { ToastContentProps } from 'react-toastify'

import { LeftArrowIcon } from '@/components/Icons/LeftArrowIcon'

// 1. Определяем интерфейс для ваших данных
interface OrderErrorProps {
	errMessage: string
	message: string
	orderUrl: string
}

// 2. Объединяем ваши пропсы с пропсами Toastify
// Partial нужен, потому что Toastify прокидывает их автоматически
type FullProps = Partial<ToastContentProps> & OrderErrorProps

export const OrderErrorToast: FC<FullProps> = ({ errMessage, message, orderUrl }) => {
	const { palette } = useTheme()

	return (
		<Box ml={1}>
			<Typography mb={0.5}>{errMessage}</Typography>
			<a href={orderUrl} style={{ textDecoration: 'none' }}>
				<Stack
					direction={'row'}
					alignItems='center'
					sx={{
						borderBottom: '1px solid transparent',
						transition: 'all 0.3s ease-in-out',
						':hover': { borderBottomColor: palette.primary.main },
					}}
				>
					<Typography color='primary' fontSize={'14px'}>
						{message}
					</Typography>
					<LeftArrowIcon fontSize={8} ml={1} transform={'rotate(180deg)'} fill={palette.primary.main} />
				</Stack>
			</a>
		</Box>
	)
}
