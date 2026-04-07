import { type FC } from 'react'
import { Autocomplete, TextField } from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'

import type { Field } from '@/components/Form/type'
import { useAutocompleteData } from './useAutocompleteData'

export const AutocompleteInput: FC<{ field: Field }> = ({ field }) => {
	const { control } = useFormContext()
	const { options, isFetching, focusHandler, filterOptions } = useAutocompleteData(field.name)

	return (
		<Controller
			name={field.name}
			control={control}
			rules={{ required: field.isRequired }}
			render={({ field: { onChange, value, ref }, fieldState: { error } }) => (
				<Autocomplete
					value={options?.find(opt => opt.label === value) || value || ''}
					freeSolo
					disableClearable
					options={options || []}
					getOptionLabel={option => (typeof option === 'string' ? option : option.label || '')}
					filterOptions={filterOptions}
					// Обработка выбора из выпадающего списка
					onChange={(_event, newValue) => {
						const val = typeof newValue === 'string' ? newValue : newValue?.label
						onChange(val || '')
					}}
					onInputChange={(_event, newInputValue) => {
						onChange(newInputValue)
					}}
					renderInput={params => (
						<TextField
							{...params}
							label={field.label}
							error={Boolean(error)}
							helperText={error?.message}
							inputRef={ref}
						/>
					)}
					loading={isFetching}
					loadingText='Поиск похожих значений...'
					noOptionsText='Ничего не найдено'
					onFocus={focusHandler}
				/>
			)}
		/>
	)
}
