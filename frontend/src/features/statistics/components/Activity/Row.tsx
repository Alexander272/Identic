import { TableRow, TableCell, IconButton, Box, Typography, Tooltip, Button } from '@mui/material'
import { ExpandMore, ExpandLess, Add, Edit, Delete } from '@mui/icons-material'
import { Link } from 'react-router'

import type { ActivityLog } from '../../types/activity'
import { ActionType } from '../../types/activity'
import { getActionLabel, getActionColor, getActionTextColor, getEntityTypeLabel, getInitials } from '../utils'
import { getSmartDate } from '@/utils/date'
import { stringToHSLA } from '@/utils/colors'
import { PopupLinkIcon } from '@/components/Icons/PopupLinkIcon'

interface ActivityTableRowProps {
	log: ActivityLog
	isExpanded: boolean
	onToggle: () => void
}

const getActionIcon = (action: ActionType) => {
	switch (action) {
		case ActionType.Insert:
			return <Add sx={{ fontSize: 14 }} />
		case ActionType.Delete:
			return <Delete sx={{ fontSize: 14 }} />
		default:
			return <Edit sx={{ fontSize: 14 }} />
	}
}

export const ActivityTableRow = ({ log, isExpanded, onToggle }: ActivityTableRowProps) => {
	const username = `${log.actor.lastName} ${log.actor.firstName}`
	const colors = stringToHSLA(username)

	return (
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
			<TableCell sx={{ whiteSpace: 'nowrap', color: 'text.secondary' }}>{getSmartDate(log.createdAt)}</TableCell>
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
						<Typography sx={{ fontWeight: 600 }}>{username.trim() || log.changedByName}</Typography>
						{log.actor.email && (
							<Typography sx={{ fontSize: 12, color: 'text.secondary' }}>{log.actor.email}</Typography>
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
						background: getActionColor(log.action),
						color: getActionTextColor(log.action),
					}}
				>
					{getActionIcon(log.action)}
					{getActionLabel(log.action)}
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
						background: log.entityType === 'order' ? 'rgba(108, 92, 231, 0.15)' : 'rgba(0, 184, 148, 0.1)',
						color: log.entityType === 'order' ? '#6c5ce7' : '#00b894',
					}}
				>
					{getEntityTypeLabel(log.entityType)}
				</Box>
			</TableCell>
			<TableCell>
				<Typography sx={{ fontWeight: 600, fontSize: 14 }}>{log.entity || log.entityId}</Typography>
			</TableCell>
			<TableCell sx={{ p: 0, px: 1 }}>
				<Link to={`/orders/${log.parentId || log.entityId}`} target='_blank' rel='noopener noreferrer'>
					<Tooltip title='Перейти к заказу'>
						<Button
							sx={theme => ({
								minWidth: 48,
								minHeight: 48,
								borderRadius: '6px',
								':hover': { svg: { fill: theme.palette.secondary.main } },
							})}
						>
							<PopupLinkIcon fontSize={16} />
						</Button>
					</Tooltip>
				</Link>
			</TableCell>
		</TableRow>
	)
}
