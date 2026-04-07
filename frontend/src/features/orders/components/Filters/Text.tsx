import type { FC } from 'react'
import { FormControl, InputLabel, MenuItem, Select } from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'

import type { IFilter } from '../../types/params'
import { TextField } from '@/components/Form/TextField'

type Props = {
	index: number
}

export const TextFilter: FC<Props> = ({ index }) => {
	const { control, watch } = useFormContext<{ filters: IFilter[] }>()
	const type = watch(`filters.${index}.compareType`)

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
							<MenuItem key='empty' value='null'>
								Не заполнено
							</MenuItem>
						</Select>
					</FormControl>
				)}
			/>

			<TextField
				data={{ name: `filters.${index}.value`, type: 'text', label: 'Значение', isRequired: type != 'null' }}
				disabled={type == 'null'}
			/>
		</>
	)
}
