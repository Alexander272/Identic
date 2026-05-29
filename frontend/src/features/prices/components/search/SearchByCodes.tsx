import { type FC } from 'react'
import { Box, Button } from '@mui/material'
import { useFormContext, Controller, useWatch } from 'react-hook-form'

import type { FormValues, SearchHandlers } from './SearchModes'
import { CodesInput } from './CodesInput'
import { ListIcon } from '@/components/Icons/ListIcon'
import { DocSearchIcon } from '@/components/Icons/DocSearchIcon'
import { RefreshIcon } from '@/components/Icons/RefreshIcon'
import { cardSx, buttonSx } from './searchStyles'

type Props = SearchHandlers & {
	isLoading: boolean
}

export const SearchByCodes: FC<Props> = ({ onSearch, onClear, isLoading }) => {
	const { control } = useFormContext<FormValues>()

	const codes = useWatch({ name: 'codes', control })

	return (
		<Box component='form' onSubmit={onSearch} sx={cardSx}>
			<Box
				component='span'
				sx={{ fontWeight: 700, fontSize: 16, mb: 2, display: 'flex', alignItems: 'center', gap: 1 }}
			>
				<ListIcon sx={{ fontSize: 14 }} />
				Поиск по списку кодов
			</Box>

			<Controller
				name='codes'
				control={control}
				render={({ field }) => (
					<CodesInput codes={field.value} onChange={field.onChange} />
				)}
			/>

			<Box sx={{ display: 'flex', gap: 1 }}>
				<Button
					type='submit'
					variant='outlined'
					disabled={isLoading}
					color='inherit'
					sx={buttonSx}
				>
					<DocSearchIcon sx={{ fontSize: 14, mr: 1 }} /> Найти по списку
				</Button>
				<Button
					type='button'
					variant='outlined'
					onClick={onClear}
					disabled={codes.length === 0}
					color='inherit'
					sx={buttonSx}
				>
					<RefreshIcon sx={{ fontSize: 14, mr: 1 }} />
					Очистить
				</Button>
			</Box>
		</Box>
	)
}
