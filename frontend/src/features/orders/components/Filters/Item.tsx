import type { FC } from 'react'
import { Controller, useFormContext } from 'react-hook-form'
import { FormControl, IconButton, InputLabel, MenuItem, Select, Stack } from '@mui/material'
import dayjs from 'dayjs'

import type { ColumnTypes, CompareTypes, IFilter } from '../../types/params'
import { Columns } from '@/constants/columns'
import { TimesIcon } from '@/components/Icons/TimesIcon'
import { TextFilter } from './Text'
import { NumberFilter } from './Number'
import { DateFilter } from './Date'
import { AutocompleteFilter } from './Autocomplete'
import { ListFilter } from './List'
import { SwitchFilter } from './SwitchFilter'

interface FilterConfig {
	compareType: CompareTypes
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	defaultValue: any
}

const FILTER_SETTINGS: Record<ColumnTypes, FilterConfig> = {
	text: {
		compareType: 'con',
		defaultValue: '',
	},
	date: {
		compareType: 'eq',
		defaultValue: () => dayjs().toISOString(), // Используем функцию, чтобы дата генерировалась в момент вызова
	},
	number: {
		compareType: 'eq',
		defaultValue: '',
	},
	list: {
		compareType: 'in',
		defaultValue: '',
	},
	autocomplete: {
		compareType: 'in',
		defaultValue: '',
	},
	bool: {
		compareType: 'eq',
		defaultValue: 'false',
	},
}

const FILTER_COMPONENTS: Record<string, FC<{ index: number }>> = {
	text: TextFilter,
	number: NumberFilter,
	date: DateFilter,
	autocomplete: AutocompleteFilter,
	list: ListFilter,
	bool: SwitchFilter,
}

interface Props {
	index: number
	onRemove: (index: number) => void
	canRemove: boolean
}

export const FilterItem: FC<Props> = ({ index, onRemove, canRemove }) => {
	const { control, setValue, watch } = useFormContext<{ filters: IFilter[] }>()

	// Следим за типом, чтобы знать, какой компонент рендерить
	const fieldType = watch(`filters.${index}.fieldType`)

	const handleFieldChange = (selectedField: string, onChange: (val: string) => void) => {
		const column = Columns.find(c => c.field === selectedField)
		if (!column) return

		const settings = FILTER_SETTINGS[column.filter]
		if (column.filter !== fieldType) {
			setValue(`filters.${index}.fieldType`, column.filter)
			setValue(`filters.${index}.compareType`, settings.compareType)

			// Проверяем, функция ли это (для динамических дат) или просто значение
			const newValue =
				typeof settings.defaultValue === 'function' ? settings.defaultValue() : settings.defaultValue

			setValue(`filters.${index}.value`, newValue)
		}

		onChange(selectedField)
	}

	const SpecificFilter = FILTER_COMPONENTS[fieldType]

	return (
		<Stack direction='row' spacing={1} alignItems='center'>
			<FormControl fullWidth sx={{ maxWidth: 170 }}>
				<InputLabel>Колонка</InputLabel>
				<Controller
					control={control}
					name={`filters.${index}.field`}
					render={({ field, fieldState: { error } }) => (
						<Select
							{...field}
							label='Колонка'
							error={!!error}
							onChange={e => handleFieldChange(e.target.value, field.onChange)}
						>
							{Columns.map(c => (
								<MenuItem key={c.field} value={c.field}>
									{c.label}
								</MenuItem>
							))}
						</Select>
					)}
				/>
			</FormControl>

			{/* Динамический компонент фильтра */}
			{SpecificFilter ? <SpecificFilter index={index} /> : null}

			{canRemove && (
				<IconButton onClick={() => onRemove(index)} size='large' color='error'>
					<TimesIcon fontSize={14} />
				</IconButton>
			)}
		</Stack>
	)
}
