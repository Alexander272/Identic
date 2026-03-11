import { FC, JSX, MouseEvent, useRef, useState } from 'react'
import { Box, Button, Popover, Stack, Typography, useTheme } from '@mui/material'
import { WarningIcon } from '../Icons/WarningIcon'

type Props = {
	onClick: () => void
	width?: string
	iconColor?: string
	buttonComponent?: JSX.Element
	confirmTitle?: string
	confirmText: string
	buttonColor?: 'error' | 'inherit' | 'primary' | 'secondary' | 'success' | 'info' | 'warning'
}

export const Confirm: FC<Props> = ({
	width,
	onClick,
	iconColor,
	buttonComponent,
	confirmTitle,
	confirmText,
	buttonColor = 'error',
}) => {
	const [open, setOpen] = useState(false)
	const anchor = useRef<HTMLButtonElement>(null)

	const { palette } = useTheme()

	const toggleHandler = (event: MouseEvent) => {
		event.stopPropagation()
		setOpen(prev => !prev)
	}

	const confirmHandler = (event: MouseEvent) => {
		toggleHandler(event)
		onClick()
	}

	return (
		<Box ref={anchor} onClick={toggleHandler} sx={{ width: width ? width : 'inherit', height: '100%' }}>
			{buttonComponent ? (
				buttonComponent
			) : (
				<Button variant='contained' color='error' fullWidth>
					Удалить
				</Button>
			)}

			<Popover
				open={open}
				anchorEl={anchor.current}
				onClose={toggleHandler}
				anchorOrigin={{
					vertical: 'center',
					horizontal: 'center',
				}}
				transformOrigin={{
					vertical: 'top',
					horizontal: 'center',
				}}
			>
				<Stack spacing={2} paddingX={2} paddingY={1.2}>
					<Box>
						<Stack spacing={1} direction={'row'} justifyContent={'center'} alignItems={'center'} mb={1}>
							<WarningIcon fill={iconColor || palette.error.main} />
							<Typography fontSize={'1.1rem'} fontWeight={'bold'} align='center'>
								{confirmTitle || 'Удаление'}
							</Typography>
						</Stack>

						<Typography maxWidth={260} align='center'>
							{confirmText}
						</Typography>
					</Box>

					<Stack direction='row' spacing={2}>
						<Button onClick={confirmHandler} variant='contained' color={buttonColor} fullWidth>
							Да
						</Button>
						<Button onClick={toggleHandler} variant='outlined' fullWidth>
							Отмена
						</Button>
					</Stack>
				</Stack>
			</Popover>
		</Box>
	)
}
