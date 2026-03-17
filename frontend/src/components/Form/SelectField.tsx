import type { FC } from 'react'
import { FormControl, InputLabel, MenuItem, Select } from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'

import type { Field, Option } from './type'

type Props = {
	data: Field
	options: Option[]
}

export const SelectField: FC<Props> = ({ data, options }) => {
	const { control } = useFormContext()

	return (
		<FormControl>
			<InputLabel id={data.name}>{data.label}</InputLabel>
			<Controller
				control={control}
				name={data.name}
				rules={{ required: data.isRequired }}
				render={({ field, fieldState: { error } }) => (
					<Select
						{...field}
						value={field.value || ''}
						labelId={data.name}
						label={data.label}
						error={Boolean(error)}
					>
						{options.map(option => (
							<MenuItem key={option.id} value={option.id}>
								{option.name}
							</MenuItem>
						))}
					</Select>
				)}
			/>
		</FormControl>
	)
}
