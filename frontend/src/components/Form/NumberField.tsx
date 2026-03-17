import type { FC } from 'react'
import { TextField } from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'

import type { Field } from './type'

type Props = {
	data: Field
	disabled?: boolean
}

export const NumberField: FC<Props> = ({ data, disabled }) => {
	const { control } = useFormContext()

	return (
		<Controller
			control={control}
			name={data.name}
			rules={{ required: data.isRequired }}
			render={({ field, fieldState: { error } }) => (
				<TextField
					{...field}
					value={field.value || ''}
					onChange={e => field.onChange(+(e.target.value || 0))}
					label={data.label}
					disabled={disabled}
					fullWidth
					error={Boolean(error)}
					slotProps={{
						htmlInput: {
							type: 'number',
							step: 1,
							min: 1,
							// max: 100
						},
					}}
				/>
			)}
		/>
	)
}
