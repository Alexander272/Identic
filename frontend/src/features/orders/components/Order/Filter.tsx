import { useState, type FC } from 'react'
import { Box, Checkbox, Divider, FormControlLabel, IconButton, Popover, TextField, Typography } from '@mui/material'

import type { IFilter } from '../../types/filter'
import { FilterIcon } from '@/components/Icons/FilterIcon'
import { useDebounceFunc } from '@/hooks/useDebounceFunc'

type Props = {
	filters: IFilter
	onChange: (name: string, value: unknown) => void
}

export const Filter: FC<Props> = ({ filters, onChange }) => {
	const [anchor, setAnchor] = useState<HTMLButtonElement | null>(null)
	const [search, setSearch] = useState(filters.search)

	const openHandler = (event: React.MouseEvent<HTMLButtonElement>) => {
		setAnchor(event.currentTarget)
	}

	const handleClose = () => {
		setAnchor(null)
	}

	const debounce = useDebounceFunc(value => {
		onChange('search', value)
	}, 300)

	const searchHandler = (event: React.ChangeEvent<HTMLInputElement>) => {
		debounce(event.target.value)
		setSearch(event.target.value)
	}

	const checkboxHandler = (event: React.ChangeEvent<HTMLInputElement>) => {
		onChange(event.target.name, event.target.checked)
	}

	return (
		<>
			<IconButton onClick={openHandler} sx={{ ml: 1 }}>
				<FilterIcon sx={{ fontSize: 16 }} />
			</IconButton>

			<Popover
				open={Boolean(anchor)}
				onClose={handleClose}
				anchorEl={anchor}
				anchorOrigin={{
					vertical: 'bottom',
					horizontal: 'center',
				}}
				transformOrigin={{
					vertical: 'top',
					horizontal: 'center',
				}}
				slotProps={{
					paper: {
						elevation: 0,
						sx: {
							overflow: 'visible',
							filter: 'drop-shadow(0px 2px 8px rgba(0,0,0,0.32))',
							mt: 1,
							paddingX: 2,
							pt: 1.5,
							paddingBottom: 2,
							maxWidth: 500,
							width: '100%',
							'&:before': {
								content: '""',
								display: 'block',
								position: 'absolute',
								top: 0,
								width: 10,
								height: 10,
								left: '50%',
								bgcolor: 'background.paper',
								transform: 'translate(-50%, -50%) rotate(45deg)',
								zIndex: 0,
							},
						},
					},
				}}
			>
				<Box>
					<Typography align='center' fontWeight={'bold'}>
						Фильтры
					</Typography>
					<Divider sx={{ mt: 1, mb: 2 }} />

					<TextField
						value={search}
						onChange={searchHandler}
						label='Наименование содержит'
						variant='outlined'
						size='small'
						fullWidth
						sx={{ mb: 1 }}
					/>

					<FormControlLabel
						control={<Checkbox name='found' checked={filters.found} onChange={checkboxHandler} />}
						label='Показать найденное'
						sx={{
							mb: 1,
							ml: 0,
							pl: 0.5,
							transition: 'background-color 0.2s ease-in-out',
							borderRadius: 2,
							width: '100%',
							':hover': { backgroundColor: '#eff8ff' },
						}}
					/>
				</Box>
			</Popover>
		</>
	)
}
