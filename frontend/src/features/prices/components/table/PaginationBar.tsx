import { type FC } from 'react'
import { Box, MenuItem, Select, Stack, Typography } from '@mui/material'

import { Pagination } from '@/components/Pagination/Pagination'

type Props = {
	page: number
	rowsPerPage: number
	totalCount: number
	onPageChange: (page: number) => void
	onRowsPerPageChange: (rowsPerPage: number) => void
}

export const PaginationBar: FC<Props> = ({ page, rowsPerPage, totalCount, onPageChange, onRowsPerPageChange }) => {
	return (
		<Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between', pt: 2, px: 1 }}>
			<Box sx={{ ml: '50%', transform: 'translateX(-50%)' }}>
				{Math.ceil(totalCount / rowsPerPage) > 1 && (
					<Pagination
						page={page}
						totalPages={Math.ceil(totalCount / rowsPerPage)}
						onClick={p => onPageChange(p - 1)}
					/>
				)}
			</Box>
			<Stack direction={'row'} alignItems={'center'}>
				<Typography variant='body2'>Строк на странице:</Typography>
				<Select
					value={rowsPerPage}
					onChange={e => onRowsPerPageChange(+e.target.value)}
					size='small'
					sx={{
						ml: 1,
						minWidth: 60,
						boxShadow: 'none',
						'.MuiOutlinedInput-notchedOutline': { border: 0 },
						'&.MuiOutlinedInput-root:hover .MuiOutlinedInput-notchedOutline': { border: 0 },
						'&.MuiOutlinedInput-root.Mui-focused .MuiOutlinedInput-notchedOutline': { border: 0 },
						'.MuiOutlinedInput-input': { padding: '6.5px 10px' },
					}}
				>
					{[15, 30, 50, 100].map(l => (
						<MenuItem key={l} value={l}>
							{l}
						</MenuItem>
					))}
				</Select>
				<Typography variant='body2' sx={{ ml: 2 }}>
					{(page - 1) * rowsPerPage || 1}-{Math.min(page * rowsPerPage, totalCount)} из {totalCount}
				</Typography>
			</Stack>
		</Box>
	)
}
