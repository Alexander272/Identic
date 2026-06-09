import { TableRow, TableCell, IconButton, Box, Typography, Tooltip } from '@mui/material'
import { ExpandMore, ExpandLess } from '@mui/icons-material'

import type { PriceSearchLog } from '../../types/priceSearch'
import { stringToHSLA } from '@/utils/colors'
import { getSmartDate } from '@/utils/date'
import { formatDuration, getInitials } from '../utils'

interface PriceSearchTableRowProps {
	log: PriceSearchLog
	isExpanded: boolean
	onToggle: () => void
}

const getStatus = (log: PriceSearchLog) => {
	if (log.resultsCount === 0) {
		return { label: 'Без результатов', className: 'warning' }
	}
	if (log.durationMs > 30000) {
		return { label: 'Медленный', className: 'warning' }
	}
	return { label: 'Успешно', className: 'success' }
}

const getSearchType = (log: PriceSearchLog): { label: string; color: string; bg: string } => {
	if ((log.codes ?? []).length > 0) {
		return { label: 'По коду', color: '#6c5ce7', bg: 'rgba(108, 92, 231, 0.15)' }
	}
	return { label: 'По тексту', color: '#00b894', bg: 'rgba(0, 184, 148, 0.1)' }
}

export const PriceSearchTableRow = ({ log, isExpanded, onToggle }: PriceSearchTableRowProps) => {
	const status = getStatus(log)
	const searchType = getSearchType(log)
	const username = `${log.actor.lastName} ${log.actor.firstName}`
	const colors = stringToHSLA(username)

	return (
		<>
			<TableRow
				hover
				sx={{
					cursor: 'pointer',
					'&:hover': { background: 'action.hover' },
				}}
				onClick={onToggle}
			>
				<TableCell>
					<IconButton size='small'>{isExpanded ? <ExpandLess /> : <ExpandMore />}</IconButton>
				</TableCell>
				<TableCell sx={{ whiteSpace: 'nowrap', color: 'text.secondary' }}>
					{getSmartDate(log.createdAt)}
				</TableCell>
				<TableCell>
					<Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
						<Box
							sx={{
								width: 32,
								height: 32,
								borderRadius: '50%',
								background: colors.bg,
								color: colors.text,
								border: `1px solid ${colors.border}`,
								display: 'flex',
								alignItems: 'center',
								justifyContent: 'center',
								fontSize: 12,
								fontWeight: 700,
							}}
						>
							{getInitials(username)}
						</Box>
						<Box>
							<Typography sx={{ fontWeight: 600 }}>{username.trim()}</Typography>
							{log.actor.email && (
								<Typography sx={{ fontSize: 12, color: 'text.secondary' }}>
									{log.actor.email}
								</Typography>
							)}
						</Box>
					</Box>
				</TableCell>
				<TableCell>
					<Box
						sx={{
							display: 'inline-flex',
							alignItems: 'center',
							gap: 0.75,
							px: 1.5,
							py: 0.5,
							borderRadius: 10,
							fontSize: 12,
							fontWeight: 600,
							background: searchType.bg,
							color: searchType.color,
						}}
					>
						{searchType.label}
					</Box>
				</TableCell>
				<TableCell align='center'>{(log.codes ?? []).length > 0 ? (log.codes ?? []).length : (log.queries ?? []).length}</TableCell>
				<TableCell align='center'>
					<Typography
						component={'span'}
						sx={{ fontWeight: 700, color: log.resultsCount > 0 ? 'primary.main' : 'text.primary' }}
					>
						{log.resultsCount}
					</Typography>{' '}
					позиций
				</TableCell>
				<TableCell align='center'>
					<Tooltip title={`${log.durationMs}мс`}>
						<span>{formatDuration(log.durationMs)}</span>
					</Tooltip>
				</TableCell>
				<TableCell>
					<Box
						sx={{
							display: 'inline-flex',
							alignItems: 'center',
							gap: 0.75,
							px: 1.5,
							py: 0.5,
							borderRadius: 10,
							fontSize: 12,
							fontWeight: 600,
							background:
								status.className === 'success' ? 'rgba(0, 184, 148, 0.1)' : 'rgba(253, 203, 110, 0.1)',
							color: status.className === 'success' ? '#00b894' : '#fdcb6e',
						}}
					>
						<Box sx={{ width: 6, height: 6, borderRadius: '50%', background: 'currentColor' }} />
						{status.label}
					</Box>
				</TableCell>
			</TableRow>
		</>
	)
}
