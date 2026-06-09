import { type FC } from 'react'
import { Box, Button, IconButton, Stack, TextField, Typography } from '@mui/material'
import { useFormContext, Controller, useFieldArray, useWatch } from 'react-hook-form'

import type { FormValues, SearchHandlers } from './SearchModes'
import { SEARCH_FIELD_OPTIONS } from './SearchModes'
import { SearchIcon } from '@/components/Icons/SearchIcon'
import { RefreshIcon } from '@/components/Icons/RefreshIcon'
import { TimesIcon } from '@/components/Icons/TimesIcon'
import { CheckIcon } from '@/components/Icons/CheckSimpleIcon'
import { buttonSx } from '../buttonStyles'
import { cardSx } from './searchStyles'

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
	const searchFields = useWatch({ name: 'searchFields', control })
	const noFieldsSelected = !searchFields?.length

	return (
		<Box component='form' onSubmit={onSearch} sx={cardSx}>
			<Stack direction='row' alignItems='center' justifyContent='space-between' sx={{ mb: 1.5 }}>
				<Typography sx={{ fontWeight: 700, fontSize: 16, display: 'flex', alignItems: 'center', gap: 1 }}>
					<SearchIcon sx={{ fontSize: 18 }} />
					Поиск
				</Typography>

				<Controller
					name='searchFields'
					control={control}
					render={({ field }) => (
						<Box sx={{ display: 'flex', gap: 0.5 }}>
							<Button
								size='small'
								color='inherit'
								sx={buttonSx}
								onClick={() => field.onChange(SEARCH_FIELD_OPTIONS.map(o => o.value))}
							>
								<CheckIcon sx={{ fontSize: 14, mr: 1 }} /> Выбрать все
							</Button>
							<Button size='small' color='inherit' sx={buttonSx} onClick={() => field.onChange([])}>
								<BoxIcon /> Снять все
							</Button>
						</Box>
					)}
				/>
			</Stack>

			<Controller
				name='searchFields'
				control={control}
				render={({ field }) => (
					<Box sx={{ display: 'flex', flexWrap: 'wrap', gap: 1, mb: 1 }}>
						{SEARCH_FIELD_OPTIONS.map(opt => {
							const selected = field.value.includes(opt.value)
							return (
								<Button
									key={opt.value}
									size='small'
									variant='outlined'
									color={selected ? 'primary' : 'inherit'}
									onClick={() =>
										field.onChange(
											selected
												? field.value.filter((v: string) => v !== opt.value)
												: [...field.value, opt.value],
										)
									}
									sx={{
										textTransform: 'none',
										minWidth: 48,
										borderRadius: '6px',
										py: 0.5,
										border: '1px solid #c3c3c4',
									}}
								>
									{selected ? (
										<CheckIcon
											sx={theme => ({
												fontSize: 14,
												mr: 1,
												fill: theme.palette.primary.main,
											})}
										/>
									) : (
										<BoxIcon />
									)}
									{opt.label}
								</Button>
							)
						})}
					</Box>
				)}
			/>

			{noFieldsSelected && (
				<Typography variant='caption' color='error' sx={{ mb: 1, display: 'block' }}>
					Выберите хотя бы одно поле для поиска
				</Typography>
			)}

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
						disabled={isLoading || !singleQuery.trim() || noFieldsSelected}
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

				<Button variant='text' onClick={() => append({ value: '' })} color='inherit' sx={buttonSx}>
					+ Добавить поле
				</Button>
			</Stack>
		</Box>
	)
}

const BoxIcon = () => (
	<Box
		component='span'
		sx={{
			width: 14,
			height: 14,
			mr: 1,
			border: '1.5px solid',
			borderColor: 'action.disabled',
			borderRadius: '2px',
			display: 'inline-flex',
			flexShrink: 0,
		}}
	/>
)
