import { type FC } from 'react'
import { Button, Stack, Typography } from '@mui/material'
import { Link } from 'react-router'

import { AppRoutes } from '@/pages/router/routes'
import { usePriceEdit } from '@/features/prices/hooks/usePriceEdit'
import { SearchModes } from '@/features/prices/components/search/SearchModes'
import { EditPositionsGrid } from '@/features/prices/components/EditPositionsGrid'
import { FileImportButton } from '@/features/prices/components/FileImportButton'
import { Fallback } from '@/components/Fallback/Fallback'
import { LeftArrowIcon } from '@/components/Icons/LeftArrowIcon'

export const EditPage: FC = () => {
	const { rows, setRows, hasChanges, hasInvalid, handleSearch, handleSave, isLoading, isSaving } = usePriceEdit()

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
