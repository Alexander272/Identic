import type { FC } from 'react'
import { FormControl, InputLabel, MenuItem, Select } from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'

import type { IFilter } from '../../types/params'
import { DateField } from '@/components/Form/DateField'

type Props = {
	index: number
}

export const DateFilter: FC<Props> = ({ index }) => {
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
							<MenuItem key='d_eq' value='eq'>
								Равна
							</MenuItem>
							<MenuItem key='d_gte' value='gte'>
								Больше или равна
							</MenuItem>
							<MenuItem key='d_lte' value='lte'>
								Меньше или равна
							</MenuItem>
							<MenuItem key='d_empty' value='null'>
								Не заполнена
							</MenuItem>
						</Select>
					</FormControl>
				)}
			/>

			<DateField
				data={{ name: `filters.${index}.value`, type: 'date', label: 'Значение', isRequired: type != 'null' }}
				disabled={type == 'null'}
			/>
		</>
	)
}
