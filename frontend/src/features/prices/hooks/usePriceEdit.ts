import { useState, useCallback, useMemo } from 'react'
import { toast } from 'react-toastify'

import type { SearchPriceRequest } from '@/features/prices/types'
import { gridRowToUpdate, positionToGridRow, type GridRow } from '@/features/prices/utils/grid'
import { isApiError } from '@/features/prices/utils/error'
import { useBatchPriceSaveMutation, useLazySearchAllPricesQuery } from '@/features/prices/priceApiSlice'

type SearchParams = {
	queries?: string[]
	fields?: string[]
	codes?: string[]
}

export const usePriceEdit = () => {
	const [batchSave, { isLoading: isSaving }] = useBatchPriceSaveMutation()
	const [search, { isLoading, data }] = useLazySearchAllPricesQuery()

	const [localRows, setLocalRows] = useState<GridRow[] | null>(null)

	const serverRows = useMemo(() => {
		return data?.data.map(positionToGridRow) || []
	}, [data])

	const rows = useMemo(() => {
		return localRows !== null ? localRows : serverRows
	}, [localRows, serverRows])

	const hasChanges = useMemo(() => rows.some(r => r.status !== 'ORIGINAL'), [rows])
	const hasInvalid = useMemo(() => rows.some(r => r.status !== 'DELETED' && !r.code), [rows])

	const handleRowsChange = useCallback((nextRows: GridRow[] | ((prev: GridRow[]) => GridRow[])) => {
		setLocalRows(prev => {
			const current = prev !== null ? prev : serverRows
			return typeof nextRows === 'function' ? nextRows(current) : nextRows
		})
	}, [serverRows])

	const handleSearch = useCallback(
		async (params: SearchParams) => {
			const isText = !!params.queries?.length
			const isCodes = !!params.codes?.length
			if (!isText && !isCodes) return

			try {
				const body: SearchPriceRequest = isCodes
					? { codes: params.codes }
					: { queries: params.queries, fields: params.fields }

				setLocalRows(null)
				await search(body).unwrap()
			} catch (err) {
				toast.error(isApiError(err) ? err.data.message : 'Ошибка поиска')
				setLocalRows([])
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
			setLocalRows(prev => {
				if (!prev) return null
				return prev
					.filter(r => r.status !== 'DELETED')
					.map(r => ({ ...r, status: 'ORIGINAL' as const }))
			})
		} catch (err) {
			toast.error(isApiError(err) ? err.data.message : 'Ошибка сохранения', { autoClose: false })
		}
	}, [rows, batchSave, hasInvalid])

	return {
		rows,
		setRows: handleRowsChange,
		hasChanges,
		hasInvalid,
		handleSearch,
		handleSave,
		isLoading,
		isSaving,
	}
}
