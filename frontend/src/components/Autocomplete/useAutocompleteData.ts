import { useCallback, useMemo, useState } from 'react'

import { useGetUniqueDataQuery } from '@/features/orders/orderApiSlice'

const orgTypesRegex = new RegExp(
	// Ищем ОПФ в начале ИЛИ в конце строки, окружённые пробелами/началом/концом
	'(^|\\s+)(ООО|АО|ПАО|НАО|ЗАО|ОАО|ИП|ГК|МФО|НКО|АНО|ФГУП|МУП|ТСЖ|СНТ|ОП|б/н|\\d+)(\\s+|$)|' + '[«»"\'„"""]', // кавычки
	'gi',
)

export const cleanString = (str: string) => {
	return str
		.replace(orgTypesRegex, '') // Удаляем типы организаций и кавычки
		.replace(/\s+/g, ' ') // Убираем двойные пробелы
		.trim()
		.toLowerCase()
}

export const useAutocompleteData = (fieldName: string) => {
	const [shouldFetch, setShouldFetch] = useState(false)

	const { data, isFetching } = useGetUniqueDataQuery({ field: fieldName }, { skip: !shouldFetch })

	const optionsWithSearch = useMemo(() => {
		return (
			data?.data.map(label => ({
				label,
				searchKey: cleanString(label), // Ваша функция очистки
			})) || []
		)
	}, [data?.data])

	const focusHandler = () => setShouldFetch(true)

	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	const filterOptions = useCallback((options: any[], { inputValue }: any) => {
		const query = cleanString(inputValue || '')

		if (!query) return options

		return options.filter(opt => {
			return opt.searchKey.includes(query) || query.includes(opt.searchKey)
		})
	}, [])

	return {
		options: optionsWithSearch,
		isFetching,
		focusHandler,
		filterOptions,
	}
}
