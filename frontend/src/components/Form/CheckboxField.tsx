import type { FC } from 'react'
import { Checkbox, FormControl, FormControlLabel, useTheme } from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'

import type { Field } from './type'

type Props = {
	data: Field
}

export const CheckboxField: FC<Props> = ({ data }) => {
	const { palette } = useTheme()
	const { control } = useFormContext()

	return (
		<FormControl>
			<Controller
				control={control}
				name={data.name}
				rules={{ required: data.isRequired }}
				render={({ field }) => (
					<FormControlLabel
						label={data.label}
						control={<Checkbox checked={field.value || false} />}
						onChange={field.onChange}
						sx={{
							transition: 'all 0.3s ease-in-out',
							borderRadius: 3,
							userSelect: 'none',
							margin: 0,
							':hover': { backgroundColor: palette.action.hover },
						}}
					/>
				)}
			/>
		</FormControl>
	)
}
