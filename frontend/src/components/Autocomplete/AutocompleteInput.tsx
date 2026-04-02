import { useMemo, useState, type FC } from 'react'
import { Autocomplete, TextField } from '@mui/material'
import { Controller, useFormContext } from 'react-hook-form'

import type { Field } from '@/components/Form/type'
import { useGetUniqueDataQuery } from '../../features/orders/orderApiSlice'

// const filterOptions = createFilterOptions<string>({
// 	// limit: 15,
// 	matchFrom: 'any',
// })

// const orgTypesRegex = new RegExp(
// 	'\\b(ООО|АО|ПАО|НАО|ЗАО|ОАО|ИП|ГК|МФО|НКО|АНО|ФГУП|МУП|ТСЖ|СНТ|ОП|б/н|\\d+)\\b|' + '[«»"\'„“”]', // Кавычки всех видов
// 	'gi',
// )
const orgTypesRegex = new RegExp(
	// Ищем ОПФ в начале ИЛИ в конце строки, окружённые пробелами/началом/концом
	'(^|\\s+)(ООО|АО|ПАО|НАО|ЗАО|ОАО|ИП|ГК|МФО|НКО|АНО|ФГУП|МУП|ТСЖ|СНТ|ОП|б/н|\\d+)(\\s+|$)|' + '[«»"\'„"""]', // кавычки
	'gi',
)

const cleanString = (str: string) => {
	return str
		.replace(orgTypesRegex, '') // Удаляем типы организаций и кавычки
		.replace(/\s+/g, ' ') // Убираем двойные пробелы
		.trim()
		.toLowerCase()
}

const defOption = { label: '', searchKey: '' }

export const AutocompleteInput: FC<{ field: Field }> = ({ field }) => {
	const [shouldFetch, setShouldFetch] = useState(false)

	const { control } = useFormContext()

	const { data, isFetching } = useGetUniqueDataQuery({ field: field.name }, { skip: !shouldFetch })

	const focusHandler = () => {
		setShouldFetch(true)
	}

	const optionsWithSearch = useMemo(() => {
		return (
			data?.data.map(label => ({
				label, // Оригинальная строка (то, что увидит юзер)
				searchKey: cleanString(label), // Очищенная строка для быстрого поиска
			})) || [defOption]
		)
	}, [data?.data])

	return (
		<Controller
			name={field.name}
			control={control}
			rules={{ required: field.isRequired }}
			render={({ field: { onChange, value, ref }, fieldState: { error } }) => (
				<Autocomplete
					value={optionsWithSearch?.find(opt => opt.label === value) || value || ''}
					freeSolo
					disableClearable
					options={optionsWithSearch || []}
					getOptionLabel={option => (typeof option === 'string' ? option : option.label || '')}
					filterOptions={(options, { inputValue }) => {
						const query = cleanString(inputValue || '')

						if (!query) return options

						return options.filter(opt => {
							return opt.searchKey.includes(query) || query.includes(opt.searchKey)
						})
					}}
					// Обработка выбора из выпадающего списка
					onChange={(_event, newValue) => {
						const val = typeof newValue === 'string' ? newValue : newValue?.label
						onChange(val || '')
					}}
					onInputChange={(_event, newInputValue) => {
						onChange(newInputValue)
					}}
					renderInput={params => (
						<TextField
							{...params}
							label={field.label}
							error={Boolean(error)}
							helperText={error?.message}
							inputRef={ref}
						/>
					)}
					loading={isFetching}
					loadingText='Поиск похожих значений...'
					noOptionsText='Ничего не найдено'
					onFocus={focusHandler}
				/>
			)}
		/>
	)
}
