import { useState, type FC } from 'react'
import {
	Autocomplete,
	FormControl,
	InputLabel,
	MenuItem,
	Select,
	TextField as MuiTextField,
	Checkbox,
	Typography,
} from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'

import type { IFilter } from '../../types/params'
import { useAutocompleteData } from '@/components/Autocomplete/useAutocompleteData'
import { TextField } from '@/components/Form/TextField'

type Props = {
	index: number
}

export const AutocompleteFilter: FC<Props> = ({ index }) => {
	const [open, setOpen] = useState(false)

	const { control, watch } = useFormContext<{ filters: IFilter[] }>()
	const field = watch(`filters.${index}.field`)
	const type = watch(`filters.${index}.compareType`)

	const { options, isFetching, focusHandler, filterOptions } = useAutocompleteData(field)

	const onOpen = () => {
		setOpen(true)
	}
	const onClose = () => {
		setOpen(false)
	}

	return (
		<>
			<Controller
				name={`filters.${index}.compareType`}
				control={control}
				rules={{ required: true }}
				render={({ field, fieldState: { error } }) => (
					<FormControl fullWidth sx={{ maxWidth: 170 }}>
						<InputLabel id={`filters.${index}.compareType`}>Условие</InputLabel>

						<Select
							{...field}
							error={Boolean(error)}
							labelId={`filters.${index}.compareType`}
							label='Условие'
						>
							<MenuItem key='con' value='con'>
								Содержит
							</MenuItem>
							<MenuItem key='like' value='like'>
								Равен
							</MenuItem>
							<MenuItem key='start' value='start'>
								Начинается с
							</MenuItem>
							<MenuItem key='end' value='end'>
								Заканчивается на
							</MenuItem>
							<MenuItem key='in' value='in'>
								В списке
							</MenuItem>
							<MenuItem key='empty' value='null'>
								Не заполнено
							</MenuItem>
						</Select>
					</FormControl>
				)}
			/>

			{type != 'in' && (
				<TextField
					data={{
						name: `filters.${index}.value`,
						type: 'text',
						label: 'Значение',
						isRequired: type != 'null',
					}}
					disabled={type == 'null'}
				/>
			)}

			{type == 'in' && (
				<Controller
					name={`filters.${index}.value`}
					control={control}
					rules={{ required: true }}
					render={({ field: { onChange, value, ref } }) => {
						const selectedLabels = value ? value.split('|') : []
						const selectedOptions = options.filter(opt => selectedLabels.includes(opt.label))

						return (
							<Autocomplete
								multiple
								fullWidth
								value={selectedOptions}
								renderValue={selected =>
									!open ? (
										<Typography
											noWrap={true}
											sx={{
												pl: 1,
												overflow: 'hidden',
												textOverflow: 'ellipsis',
												pointerEvents: 'none',
												width: '100%',
											}}
											color='textPrimary'
										>
											{selected.map(opt => opt.label).join('; ')}
										</Typography>
									) : null
								}
								options={options}
								getOptionLabel={option => (typeof option === 'string' ? option : option.label || '')}
								// filterOptions={(options, params) => {
								// 	const filtered = filterOptions(options, params) // Ваша функция из хука
								// 	return filtered.slice(0, 50) // Рендерим только первые 50 совпадений
								// }}
								filterOptions={filterOptions}
								isOptionEqualToValue={(option, value) => {
									const optionLabel = typeof option === 'string' ? option : option.label
									const valueLabel = typeof value === 'string' ? value : value.label
									return optionLabel === valueLabel
								}}
								disableCloseOnSelect
								disableClearable
								loading={isFetching}
								loadingText='Поиск похожих значений...'
								noOptionsText='Ничего не найдено'
								onFocus={focusHandler}
								onOpen={onOpen}
								onClose={onClose}
								onChange={(_event, newValue) => {
									const labels = newValue.map(v => (typeof v === 'string' ? v : v.label))
									onChange(labels.join('|'))
								}}
								renderOption={(props, option, { selected }) => {
									const { key, ...optionProps } = props
									return (
										<li key={key} {...optionProps}>
											<Checkbox style={{ marginRight: 8, marginLeft: -8 }} checked={selected} />
											{option.label}
										</li>
									)
								}}
								renderInput={params => (
									<MuiTextField
										{...params}
										fullWidth
										label={'Значение'}
										placeholder='Поиск'
										inputRef={ref}
									/>
								)}
								sx={{
									maxWidth: 362,
									'.MuiOutlinedInput-root': {
										flexWrap: 'nowrap',
									},
									'.MuiAutocomplete-inputRoot .MuiAutocomplete-input': {
										minWidth: 0,
									},
								}}
							/>
						)
					}}
				/>
			)}
		</>
	)
}
