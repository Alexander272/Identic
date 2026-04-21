import type { FC } from 'react'
import { Dialog, DialogTitle, IconButton, Typography } from '@mui/material'

import type { IUserData } from '@/features/user/types/user'
import { TimesIcon } from '@/components/Icons/TimesIcon'

type Props = {
	user: IUserData | null
	onClose: () => void
}

export const LoginsModal: FC<Props> = ({ user, onClose }) => {
	return (
		<Dialog
			open={Boolean(user)}
			onClose={onClose}
			fullWidth
			maxWidth='sm'
			slotProps={{
				paper: {
					sx: {
						borderRadius: '16px',
						p: 1,
					},
				},
			}}
		>
			{/* Header */}
			<DialogTitle sx={{ m: 0, p: 2, display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
				<Typography variant='h6' component='div' sx={{ fontWeight: 'bold' }}>
					Данные о входах
				</Typography>

				<IconButton onClick={onClose} sx={{ color: 'text.secondary' }}>
					<TimesIcon fontSize={16} />
				</IconButton>
			</DialogTitle>

			{/* {isLoading && <BoxFallback />} */}
		</Dialog>
	)
}
