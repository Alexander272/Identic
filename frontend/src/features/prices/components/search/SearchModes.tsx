/* eslint-disable react-refresh/only-export-components */

import { useCallback, type FC } from 'react'
import { Box } from '@mui/material'
import { useForm, FormProvider } from 'react-hook-form'

import { SearchByText } from './SearchByText'
import { SearchByCodes } from './SearchByCodes'

export const DEFAULT_SEARCH_FIELDS = ['currentName', 'newName', 'template']

export const SEARCH_FIELD_OPTIONS = [
	{ value: 'currentName', label: 'Наименование СИБУР' },
	{ value: 'newName', label: 'Наименование СИЛУР' },
	{ value: 'template', label: 'Шаблон' },
	{ value: 'price', label: 'Цена' },
] as const

export type FormValues = {
	singleQuery: string
	codes: string[]
	extraQueries: { value: string }[]
	searchFields: string[]
}

export type SearchHandlers = {
	onSearch: (e?: React.BaseSyntheticEvent) => Promise<void>
	onClear?: () => void
	onReset?: () => void
}

type SearchModesProps = {
	onSearch: (params: { queries?: string[]; fields?: string[]; codes?: string[] }) => void
	isLoading: boolean
	onResetSearch?: () => void
}

export const SearchModes: FC<SearchModesProps> = ({ onSearch, isLoading, onResetSearch }) => {
	const form = useForm<FormValues>({
		defaultValues: {
			singleQuery: '',
			codes: [],
			extraQueries: [],
			searchFields: DEFAULT_SEARCH_FIELDS,
		},
	})

	const onTextSearch = form.handleSubmit(async formData => {
		const queries = [formData.singleQuery, ...formData.extraQueries.map(q => q.value)]
			.map(q => q.trim())
			.filter(Boolean)
		if (!queries.length) return
		const selected = formData.searchFields
		const fields = selected.length > 0 ? selected : undefined
		onSearch({ queries, fields })
	})

	const onCodesSearch = form.handleSubmit(async formData => {
		const codes = formData.codes
		if (!codes.length) return
		onSearch({ codes })
	})

	const handleReset = useCallback(() => {
		form.setValue('singleQuery', '')
		form.setValue('extraQueries', [])
		onResetSearch?.()
	}, [form, onResetSearch])

	const handleClear = useCallback(() => {
		form.setValue('codes', [])
		onResetSearch?.()
	}, [form, onResetSearch])

	return (
		<FormProvider {...form}>
			<Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap', mb: 1 }}>
				<SearchByCodes onSearch={onCodesSearch} onClear={handleClear} isLoading={isLoading} />
				<SearchByText onSearch={onTextSearch} onReset={handleReset} isLoading={isLoading} />
			</Box>
		</FormProvider>
	)
}
