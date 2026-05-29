import { useCallback, useEffect, type FC } from 'react'
import { Box } from '@mui/material'
import { useForm, FormProvider } from 'react-hook-form'

import type { Price } from '@/features/prices/types/types'
import { useSearchPriceMutation } from '@/features/prices/priceApiSlice'
import type { IFetchError } from '@/app/types/error'
import { SearchByText } from './SearchByText'
import { SearchByCodes } from './SearchByCodes'

export type FormValues = {
	singleQuery: string
	codes: string[]
	extraQueries: { value: string }[]
}

export type SearchHandlers = {
	onSearch: (e?: React.BaseSyntheticEvent) => Promise<void>
	onClear?: () => void
	onReset?: () => void
}

type SearchModesProps = {
	onSearchResults: (results: Price[], params: { queries?: string[]; codes?: string[] }) => void
	onError: (error: unknown) => void
	onLoadingChange: (loading: boolean) => void
}

export const SearchModes: FC<SearchModesProps> = ({ onSearchResults, onError, onLoadingChange }) => {
	const form = useForm<FormValues>({
		defaultValues: { singleQuery: '', codes: [], extraQueries: [] },
	})

	const [search, { isLoading }] = useSearchPriceMutation()

	useEffect(() => {
		onLoadingChange(isLoading)
	}, [isLoading, onLoadingChange])

	const onSearch = form.handleSubmit(async formData => {
		const queries = [formData.singleQuery, ...formData.extraQueries.map(q => q.value)]
			.map(q => q.trim())
			.filter(Boolean)
		if (!queries.length) return
		try {
			const data = await search({ queries }).unwrap()
			onSearchResults(data.data ?? [], { queries })
		} catch (err) {
			const fetchErr = err as IFetchError
			if (fetchErr.data) onError(fetchErr.data)
			else onError(err)
		}
	})

	const onCodesSearch = form.handleSubmit(async formData => {
		const codes = formData.codes
		if (!codes.length) return
		try {
			const data = await search({ codes }).unwrap()
			onSearchResults(data.data ?? [], { codes })
		} catch (err) {
			onError(err)
		}
	})

	const handleReset = useCallback(() => {
		form.setValue('singleQuery', '')
		form.setValue('extraQueries', [])
	}, [form])

	const handleClear = useCallback(() => {
		form.setValue('codes', [])
	}, [form])

	return (
		<FormProvider {...form}>
			<Box sx={{ display: 'flex', gap: 1, flexWrap: 'wrap', mb: 1 }}>
				<SearchByText onSearch={onSearch} onReset={handleReset} isLoading={isLoading} />
				<SearchByCodes onSearch={onCodesSearch} onClear={handleClear} isLoading={isLoading} />
			</Box>
		</FormProvider>
	)
}
