import { type FC } from 'react'
import { Box, Button, IconButton, Stack, TextField } from '@mui/material'
import { useFormContext, Controller, useFieldArray, useWatch } from 'react-hook-form'

import type { FormValues, SearchHandlers } from './SearchModes'
import { SearchIcon } from '@/components/Icons/SearchIcon'
import { RefreshIcon } from '@/components/Icons/RefreshIcon'
import { TimesIcon } from '@/components/Icons/TimesIcon'
import { cardSx, buttonSx } from './searchStyles'

type Props = SearchHandlers & {
	isLoading: boolean
}

export const SearchByText: FC<Props> = ({ onSearch, onReset, isLoading }) => {
	const { control } = useFormContext<FormValues>()

	const { fields, append, remove } = useFieldArray({
		control,
		name: 'extraQueries',
	})

	const singleQuery = useWatch({ name: 'singleQuery', control })

	return (
		<Box component='form' onSubmit={onSearch} sx={cardSx}>
			<Box
				component='span'
				sx={{ fontWeight: 700, fontSize: 16, mb: 2, display: 'flex', alignItems: 'center', gap: 1 }}
			>
				<SearchIcon sx={{ fontSize: 18 }} />
				Поиск
			</Box>

			<Controller
				name='singleQuery'
				control={control}
				render={({ field }) => (
					<TextField
						{...field}
						placeholder='Введите текст запроса...'
						size='small'
						fullWidth
						sx={{ mb: 1, '& .MuiOutlinedInput-notchedOutline': { borderRadius: 1.5 } }}
					/>
				)}
			/>

			{fields.map((field, index) => (
				<Box key={field.id} sx={{ display: 'flex', gap: 0.5, mb: 1 }}>
					<Controller
						name={`extraQueries.${index}.value`}
						control={control}
						render={({ field: innerField }) => (
							<TextField
								{...innerField}
								placeholder='Дополнительный запрос...'
								size='small'
								fullWidth
								sx={{ '& .MuiOutlinedInput-notchedOutline': { borderRadius: 1.5 } }}
							/>
						)}
					/>
					<IconButton size='large' onClick={() => remove(index)} sx={{ alignSelf: 'center' }}>
						<TimesIcon fontSize={14} />
					</IconButton>
				</Box>
			))}

			<Stack direction={'row'} spacing={2} sx={{ justifyContent: 'space-between' }}>
				<Box sx={{ display: 'flex', gap: 1 }}>
					<Button
						type='submit'
						variant='outlined'
						disabled={isLoading || !singleQuery.trim()}
						color='inherit'
						sx={buttonSx}
					>
						<SearchIcon sx={{ fontSize: 14, mr: 1 }} /> Найти
					</Button>
					<Button
						type='button'
						variant='outlined'
						onClick={onReset}
						disabled={!singleQuery.trim() && fields.length === 0}
						color='inherit'
						sx={buttonSx}
					>
						<RefreshIcon sx={{ fontSize: 14, mr: 1 }} /> Очистить
					</Button>
				</Box>

				<Button variant='text' size='small' onClick={() => append({ value: '' })} color='inherit' sx={buttonSx}>
					+ Добавить поле
				</Button>
			</Stack>
		</Box>
	)
}
