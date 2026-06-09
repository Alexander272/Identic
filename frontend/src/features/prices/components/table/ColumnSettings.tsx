import { useState, type FC } from 'react'
import { Badge, Button, Checkbox, FormControlLabel, Menu, MenuItem, Typography } from '@mui/material'

import { buttonSx } from '../buttonStyles'
import { COLUMNS } from './cells'
import { SettingIcon } from '@/components/Icons/SettingIcon'

type Props = {
	visibleColumns: string[]
	onToggleColumn: (key: string) => void
}

export const ColumnSettings: FC<Props> = ({ visibleColumns, onToggleColumn }) => {
	const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null)
	const open = Boolean(anchorEl)

	return (
		<>
			<Button
				onClick={e => setAnchorEl(e.currentTarget)}
				size='small'
				variant='outlined'
				color='inherit'
				sx={buttonSx}
			>
				<Badge color='primary' variant='dot' invisible={visibleColumns.length === COLUMNS.length}>
					<SettingIcon sx={{ fontSize: 18 }} />
				</Badge>
				<Typography sx={{ ml: 1 }}>Настройка колонок</Typography>
			</Button>
			<Menu
				anchorEl={anchorEl}
				open={open}
				onClose={() => setAnchorEl(null)}
				anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
				transformOrigin={{ vertical: 'top', horizontal: 'right' }}
			>
				{COLUMNS.map(col => (
					<MenuItem key={col.key} dense disableRipple>
						<FormControlLabel
							control={
								<Checkbox
									checked={visibleColumns.includes(col.key)}
									onChange={() => onToggleColumn(col.key)}
									size='small'
								/>
							}
							label={col.label}
							sx={{ m: 0 }}
						/>
					</MenuItem>
				))}
			</Menu>
		</>
	)
}
