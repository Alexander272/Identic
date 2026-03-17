import type { FC } from 'react'
import { TextField as MTextField } from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'

import type { Field } from './type'

type Props = {
	data: Field
	disabled?: boolean
}

export const TextField: FC<Props> = ({ data, disabled }) => {
	const { control } = useFormContext()

	return (
		<Controller
			control={control}
			name={data.name}
			rules={{ required: data.isRequired }}
			render={({ field, fieldState: { error } }) => (
				<MTextField
					{...field}
					value={field.value || ''}
					label={data.label}
					fullWidth
					error={Boolean(error)}
					disabled={disabled}
					multiline
				/>
			)}
		/>
	)
}
