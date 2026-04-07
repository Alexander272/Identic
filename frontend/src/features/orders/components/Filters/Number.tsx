import type { FC } from 'react'
import { FormControl, InputLabel, MenuItem, Select } from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'

import type { IFilter } from '../../types/params'
import { NumberField } from '@/components/Form/NumberField'

type Props = {
	index: number
}

export const NumberFilter: FC<Props> = ({ index }) => {
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
							<MenuItem key='n_eq' value='eq'>
								Равно
							</MenuItem>
							<MenuItem key='n_gte' value='gte'>
								Больше или равно
							</MenuItem>
							<MenuItem key='n_lte' value='lte'>
								Меньше или равно
							</MenuItem>
							<MenuItem key='n_empty' value='null'>
								Не заполнено
							</MenuItem>
						</Select>
					</FormControl>
				)}
			/>

			<NumberField
				data={{ name: `filters.${index}.value`, type: 'number', label: 'Значение', isRequired: type != 'null' }}
				disabled={type == 'null'}
			/>
		</>
	)
}
