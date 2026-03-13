import { useState, type FC } from 'react'
import { Checkbox, FormControlLabel, InputAdornment, Stack, TextField } from '@mui/material'

import type { IFilter } from '../../types/filter'
import { useDebounceFunc } from '@/hooks/useDebounceFunc'
import { SearchIcon } from '@/components/Icons/SearchIcon'

type Props = {
	filters: IFilter
	onChange: (name: string, value: unknown) => void
	showFound?: boolean
}

export const Filter: FC<Props> = ({ filters, onChange, showFound }) => {
	const [search, setSearch] = useState(filters.search)

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
		<Stack direction={'row'} width={'100%'} mt={1}>
			<TextField
				value={search}
				onChange={searchHandler}
				label='Наименование содержит'
				variant='outlined'
				size='small'
				sx={{ mb: 1, width: '40%', ml: '30%', mr: 1 }}
				slotProps={{
					input: {
						endAdornment: (
							<InputAdornment position='end'>
								<SearchIcon fontSize={14} />
							</InputAdornment>
						),
					},
				}}
			/>

			{showFound && (
				<FormControlLabel
					control={<Checkbox name='found' checked={filters.found} onChange={checkboxHandler} />}
					label='Показать только найденное'
					sx={{
						mb: 1,
						ml: 0,
						pl: 0.5,
						transition: 'background-color 0.2s ease-in-out',
						borderRadius: 2,
						flexGrow: 1,
						':hover': { backgroundColor: '#eff8ff' },
					}}
				/>
			)}
		</Stack>
	)
}
