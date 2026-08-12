import { type FC } from 'react'
import { Box, Button, Divider, List, ListItem, ListItemText, Popover, Typography } from '@mui/material'

import type { Price } from '@/features/prices/types'

type Props = {
	anchorEl: HTMLElement | null
	selected: Price[]
	onClose: () => void
	onClear: () => void
}

export const SelectedPositionsPopover: FC<Props> = ({ anchorEl, selected, onClose, onClear }) => {
	const open = Boolean(anchorEl)

	return (
		<Popover
			anchorEl={anchorEl}
			open={open}
			onClose={onClose}
			anchorOrigin={{ vertical: 'bottom', horizontal: 'left' }}
			transformOrigin={{ vertical: 'top', horizontal: 'left' }}
			slotProps={{ paper: { sx: { minWidth: 320, maxWidth: 480, maxHeight: 400, display: 'flex', flexDirection: 'column' } } }}
		>
			<Box sx={{ px: 2, py: 1.5, display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
				<Typography variant='subtitle2'>Выбранные позиции ({selected.length})</Typography>
				<Button size='small' color='inherit' onClick={onClear}>
					Очистить
				</Button>
			</Box>
			<Divider />
			<List dense sx={{ overflow: 'auto', py: 0 }}>
				{selected.map(item => (
					<ListItem key={item.code} sx={{ px: 2 }}>
						<ListItemText
							primary={
								<Typography sx={{ fontFamily: 'monospace', fontWeight: 600, fontSize: '0.85rem' }}>
									{item.code}
								</Typography>
							}
							secondary={item.currentName || undefined}
							secondaryTypographyProps={{ fontSize: '0.8rem' }}
						/>
					</ListItem>
				))}
			</List>
		</Popover>
	)
}
