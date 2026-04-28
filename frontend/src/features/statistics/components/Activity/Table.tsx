import { Fragment, useState, useMemo } from 'react'
import { Paper, Table, TableBody, TableContainer, TableHead, TableRow, TableCell, Box } from '@mui/material'

import type { ActivityLog } from '../../types/activity'
import { ActivityTableRow } from './Row'
import { ActivityTableExpanded } from './Expanded'
import { Pagination } from '@/components/Pagination/Pagination'

interface ActivityTableProps {
	data: ActivityLog[]
}

const headerCellStyle = {
	fontSize: '11px',
	fontWeight: 700,
	textTransform: 'uppercase',
	letterSpacing: '0.5px',
}

const ROWS_PER_PAGE = 10

export const ActivityTable = ({ data }: ActivityTableProps) => {
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
							<TableCell width={170} sx={headerCellStyle}>
								Время
							</TableCell>
							<TableCell width={290} sx={headerCellStyle}>
								Пользователь
							</TableCell>
							<TableCell width={180} sx={headerCellStyle}>
								Действие
							</TableCell>
							<TableCell width={180} sx={headerCellStyle}>
								Тип объекта
							</TableCell>
							<TableCell sx={{ ...headerCellStyle, maxWidth: 920 }}>Объект</TableCell>
							<TableCell width={50} sx={{ p: 0, px: 1 }} />
						</TableRow>
					</TableHead>
					<TableBody>
						{!paginatedData.length && (
							<TableRow>
								<TableCell colSpan={7} sx={{ py: 2, textAlign: 'center', fontWeight: 'bold' }}>
									Заказы не менялись
								</TableCell>
							</TableRow>
						)}

						{paginatedData.map(log => (
							<Fragment key={log.id}>
								<ActivityTableRow
									key={log.id}
									log={log}
									isExpanded={expandedId === log.id}
									onToggle={() => toggleExpand(log.id)}
								/>
								<ActivityTableExpanded log={log} open={expandedId === log.id} />
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
