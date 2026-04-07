import { useEffect, type FC } from 'react'
import { Controller, useFormContext } from 'react-hook-form'

import type { IFilter } from '../../types/params'
import { useAutocompleteData } from '@/components/Autocomplete/useAutocompleteData'
import { SelectWithFilter } from '@/components/SelectWithFilter/SelectWithFilter'

type Props = {
	index: number
}

export const ListFilter: FC<Props> = ({ index }) => {
	const { control, watch } = useFormContext<{ filters: IFilter[] }>()
	const field = watch(`filters.${index}.field`)
	const type = watch(`filters.${index}.compareType`)

	const { options, isFetching, focusHandler, filterOptions } = useAutocompleteData(field)

	useEffect(() => {
		focusHandler()
	})

	return (
		<>
			{/* <Controller
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
			)} */}

			{type == 'in' && (
				<Controller
					name={`filters.${index}.value`}
					control={control}
					rules={{ required: true }}
					render={({ field: { onChange, value } }) => {
						const selectedLabels = value ? value.split('|') : []
						const selectedOptions = options.filter(opt => selectedLabels.includes(opt.label))

						return (
							<SelectWithFilter
								values={selectedOptions}
								options={options}
								onChange={newValue => {
									const labels = newValue.map(v => (typeof v === 'string' ? v : v.label))
									onChange(labels.join('|'))
								}}
								label={'Значение'}
								isLoading={isFetching}
								onFocus={focusHandler}
								filterOptions={filterOptions}
							/>
						)
					}}
				/>
			)}
		</>
	)
}
