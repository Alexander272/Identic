import { useState, useCallback, type FC, useRef } from 'react'
import { Button, Stack, Typography } from '@mui/material'
import { Link } from 'react-router'
import { toast } from 'react-toastify'

import type { IFetchError } from '@/app/types/error'
import type { SearchPriceRequest } from '@/features/prices/types/types'
import { AppRoutes } from '@/pages/router/routes'
import { gridRowToUpdate, positionToGridRow, type GridRow } from '@/features/prices/utils/grid'
import { useBatchPriceSaveMutation, useSearchAllPricesMutation } from '@/features/prices/priceApiSlice'
import { SearchModes } from '@/features/prices/components/search/SearchModes'
import { EditPositionsGrid } from '@/features/prices/components/EditPositionsGrid'
import { FileImportButton } from '@/features/prices/components/FileImportButton'
import { Fallback } from '@/components/Fallback/Fallback'
import { LeftArrowIcon } from '@/components/Icons/LeftArrowIcon'

export const EditPage: FC = () => {
	const [batchSave, { isLoading: isSaving }] = useBatchPriceSaveMutation()

	const [rows, setRows] = useState<GridRow[]>([])
	const textSearchParams = useRef<{ queries: string[]; fields?: string[] } | null>(null)

	const hasChanges = rows.some(r => r.status !== 'ORIGINAL')
	const hasInvalid = rows.some(r => r.status !== 'DELETED' && !r.code)

	const [search, { isLoading }] = useSearchAllPricesMutation()

	const handleSearch = useCallback(
		async (params: { queries?: string[]; fields?: string[]; codes?: string[] }) => {
			const isText = !!params.queries?.length
			const isCodes = !!params.codes?.length
			if (!isText && !isCodes) return

			textSearchParams.current = isText ? { queries: params.queries!, fields: params.fields } : null

			try {
				const body: SearchPriceRequest = isCodes
					? { codes: params.codes }
					: { queries: params.queries, fields: params.fields }
				const result = await search(body).unwrap()

				setRows(result.data.map(positionToGridRow))
			} catch (err) {
				const fetchErr = err as IFetchError
				toast.error(fetchErr.data.message)
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
				positions: toUpdate.map(gridRowToUpdate),
				delete_codes: toDelete.length > 0 ? toDelete : undefined,
			}).unwrap()
			toast.success('Изменения сохранены')
			setRows(prev => prev.filter(r => r.status !== 'DELETED').map(r => ({ ...r, status: 'ORIGINAL' as const })))
		} catch (err) {
			const fetchError = err as { data?: { message?: string } }
			toast.error(fetchError.data?.message ?? 'Ошибка сохранения', { autoClose: false })
		}
	}, [rows, batchSave, hasInvalid])

	return (
		<>
			<Stack direction={'row'} spacing={2} sx={{ mt: 1, mb: 3, alignItems: 'center' }}>
				<Button
					variant='outlined'
					color='inherit'
					component={Link}
					to={AppRoutes.Price}
					sx={{ textTransform: 'none', width: 100 }}
				>
					<LeftArrowIcon fontSize={12} mr={1} />
					Назад
				</Button>

				<Typography align='center' variant='h5' sx={{ width: 'calc(100% - 100px)' }}>
					Редактирование позиций
				</Typography>
			</Stack>

			<SearchModes onSearch={handleSearch} isLoading={isLoading} />

			<FileImportButton />

			{isLoading && <Fallback />}

			<EditPositionsGrid
				rows={rows}
				onRowsChange={setRows}
				onSave={handleSave}
				isSaving={isSaving}
				hasChanges={hasChanges}
				hasInvalid={hasInvalid}
			/>
		</>
	)
}
