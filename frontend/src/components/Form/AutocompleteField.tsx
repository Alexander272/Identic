import { type FC } from 'react'
import { Autocomplete, TextField, type FilterOptionsState } from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'

import type { Field } from './type'

type Props = {
	data: Field
	options: string[]
	isLoading?: boolean
	onFocus?: () => void
	filterOptions?: (options: string[], state: FilterOptionsState<string>) => string[]
}

export const AutocompleteField: FC<Props> = ({ data, options, isLoading, onFocus, filterOptions }) => {
	const { control } = useFormContext()

	const focusHandler = () => {
		if (onFocus) onFocus()
	}

	return (
		<Controller
			name={data.name}
			control={control}
			rules={{ required: data.isRequired }}
			render={({ field: { onChange, value, ref }, fieldState: { error } }) => (
				<Autocomplete
					value={value || ''}
					freeSolo
					disableClearable
					autoComplete
					options={options}
					loading={isLoading}
					loadingText='Поиск похожих значений...'
					noOptionsText='Ничего не найдено'
					filterOptions={filterOptions}
					onChange={(_event, value) => {
						onChange(value)
					}}
					onFocus={focusHandler}
					renderInput={params => (
						<TextField
							{...params}
							label={data.label}
							onChange={onChange}
							error={Boolean(error)}
							helperText={error?.message}
							inputRef={ref}
						/>
					)}
				/>
			)}
		/>
	)
}
