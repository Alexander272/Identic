import { Fragment, useState, useMemo } from 'react'
import { Paper, Table, TableBody, TableContainer, TableHead, TableRow, TableCell, Box } from '@mui/material'

import type { SearchLog } from '../../types/search'
import { SearchTableRow } from './Row'
import { SearchTableExpanded } from './Expanded'
import { Pagination } from '@/components/Pagination/Pagination'

interface SearchTableProps {
	data: SearchLog[]
}

const ROWS_PER_PAGE = 10

export const SearchTable = ({ data }: SearchTableProps) => {
	const [expandedId, setExpandedId] = useState<string | null>(null)
	const [page, setPage] = useState(1)

	const toggleExpand = (id: string) => {
		setExpandedId(expandedId === id ? null : id)
	}

	const totalPages = Math.ceil(data.length / ROWS_PER_PAGE)

	const paginatedData = useMemo(() => {
		const start = (page - 1) * ROWS_PER_PAGE
		return data.slice(start, start + ROWS_PER_PAGE)
	}, [data, page])

	return (
		<Paper
			elevation={0}
			sx={{
				border: '1px solid',
				borderColor: 'divider',
				borderRadius: 3,
				overflow: 'hidden',
			}}
		>
			<TableContainer>
				<Table>
					<TableHead>
						<TableRow sx={{ background: 'action.hover' }}>
							<TableCell width={50} />
							<TableCell sx={headerCellStyle}>Время</TableCell>
							<TableCell sx={headerCellStyle}>Пользователь</TableCell>
							<TableCell sx={headerCellStyle}>Тип поиска</TableCell>
							<TableCell sx={{ ...headerCellStyle, textAlign: 'center' }}>Позиций</TableCell>
							<TableCell sx={{ ...headerCellStyle, textAlign: 'center' }}>Результаты</TableCell>
							<TableCell sx={{ ...headerCellStyle, textAlign: 'center' }}>Время</TableCell>
							<TableCell sx={headerCellStyle}>Статус</TableCell>
						</TableRow>
					</TableHead>
					<TableBody>
						{!paginatedData.length && (
							<TableRow>
								<TableCell colSpan={8} sx={{ py: 2, textAlign: 'center', fontWeight: 'bold' }}>
									Поисков не было
								</TableCell>
							</TableRow>
						)}

						{paginatedData.map(log => (
							<Fragment key={log.id}>
								<SearchTableRow
									key={log.id}
									log={log}
									isExpanded={expandedId === log.id}
									onToggle={() => toggleExpand(log.id)}
								/>
								{expandedId === log.id && <SearchTableExpanded log={log} />}
							</Fragment>
						))}
					</TableBody>
				</Table>
			</TableContainer>
			{totalPages > 1 && (
				<Box sx={{ py: 2, display: 'flex', justifyContent: 'center' }}>
					<Pagination page={page} totalPages={totalPages} onClick={setPage} />
				</Box>
			)}
		</Paper>
	)
}

const headerCellStyle = {
	fontSize: '11px',
	fontWeight: 700,
	textTransform: 'uppercase',
	letterSpacing: '0.5px',
}
