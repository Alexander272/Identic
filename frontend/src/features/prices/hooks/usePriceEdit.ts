import { useState, useCallback } from 'react'
import { toast } from 'react-toastify'

import type { SearchPriceRequest } from '@/features/prices/types'
import { gridRowToUpdate, positionToGridRow, type GridRow } from '@/features/prices/utils/grid'
import { isApiError } from '@/features/prices/utils/error'
import { useBatchPriceSaveMutation, useSearchAllPricesMutation } from '@/features/prices/priceApiSlice'

type SearchParams = {
	queries?: string[]
	fields?: string[]
	codes?: string[]
}

export const usePriceEdit = () => {
	const [batchSave, { isLoading: isSaving }] = useBatchPriceSaveMutation()
	const [search, { isLoading }] = useSearchAllPricesMutation()

	const [rows, setRows] = useState<GridRow[]>([])

	const hasChanges = rows.some(r => r.status !== 'ORIGINAL')
	const hasInvalid = rows.some(r => r.status !== 'DELETED' && !r.code)

	const handleSearch = useCallback(
		async (params: SearchParams) => {
			const isText = !!params.queries?.length
			const isCodes = !!params.codes?.length
			if (!isText && !isCodes) return

			try {
				const body: SearchPriceRequest = isCodes
					? { codes: params.codes }
					: { queries: params.queries, fields: params.fields }
				const result = await search(body).unwrap()

				setRows(result.data.map(positionToGridRow))
			} catch (err) {
				toast.error(isApiError(err) ? err.data.message : 'Ошибка поиска')
				setRows([])
			}
		},
		[search],
	)

	const handleSave = useCallback(async () => {
		if (hasInvalid) {
			toast.error('Заполните код для всех новых строк')
			return
		}

		const toUpdate = rows.filter(r => r.status === 'CREATED' || r.status === 'UPDATED')
		const toDelete = rows.filter(r => r.status === 'DELETED').map(r => r.code || '')

		try {
			await batchSave({
				prices: toUpdate.map(gridRowToUpdate),
				deleteCodes: toDelete.length > 0 ? toDelete : undefined,
			}).unwrap()
			toast.success('Изменения сохранены')
			setRows(prev => prev.filter(r => r.status !== 'DELETED').map(r => ({ ...r, status: 'ORIGINAL' as const })))
		} catch (err) {
			toast.error(isApiError(err) ? err.data.message : 'Ошибка сохранения', { autoClose: false })
		}
	}, [rows, batchSave, hasInvalid])

	return {
		rows,
		setRows,
		hasChanges,
		hasInvalid,
		handleSearch,
		handleSave,
		isLoading,
		isSaving,
	}
}
