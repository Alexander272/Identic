import { TableRow, TableCell, Box, Typography, Table, TableBody, TableHead, Collapse } from '@mui/material'

import type { SearchLog } from '../../types/search'

interface SearchTableExpandedProps {
	log: SearchLog
	open: boolean
}

export const SearchTableExpanded = ({ log, open }: SearchTableExpandedProps) => {
	const parseQuery = (query: SearchLog['query']) => {
		if (Array.isArray(query)) {
			return query
		}
		if (typeof query === 'string') {
			try {
				return JSON.parse(query)
			} catch {
				return []
			}
		}
		return []
	}

	const queryItems = parseQuery(log.query)

	return (
		<TableRow>
			<TableCell
				colSpan={8}
				sx={{
					py: 0,
					borderTop: '1px solid',
					borderColor: 'divider',
					background: 'action.hover',
					borderBottom: open ? '1px solid #eee' : 'none',
				}}
			>
				<Collapse in={open} timeout='auto' unmountOnExit>
					<Box sx={{ px: 4, py: 1.5 }}>
						<Typography variant='subtitle2' sx={{ mb: 0.5, fontWeight: 600 }}>
							Поисковый запрос
						</Typography>
						<Table size='small'>
							<TableHead>
								<TableRow>
									<TableCell sx={{ fontWeight: 700 }}>№</TableCell>
									<TableCell sx={{ fontWeight: 700 }}>Наименование</TableCell>
									<TableCell sx={{ fontWeight: 700 }} align='right'>
										Количество
									</TableCell>
								</TableRow>
							</TableHead>
							<TableBody>
								{queryItems.length > 0 ? (
									queryItems.map(
										(item: { id: number; name: string; quantity?: number }, idx: number) => (
											<TableRow key={item.id || idx}>
												<TableCell>{idx + 1}</TableCell>
												<TableCell>
													<Typography sx={{ fontSize: 13 }}>{item.name}</Typography>
												</TableCell>
												<TableCell align='right'>{item.quantity ?? '-'}</TableCell>
											</TableRow>
										),
									)
								) : (
									<TableRow>
										<TableCell colSpan={3}>Нет данных</TableCell>
									</TableRow>
								)}
							</TableBody>
						</Table>
					</Box>
				</Collapse>
			</TableCell>
		</TableRow>
	)
}
