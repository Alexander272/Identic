import {
	TableRow,
	TableCell,
	Box,
	Typography,
	Table,
	TableBody,
	TableHead,
	TableRow as MuiTableRow,
	Collapse,
} from '@mui/material'

import type { ActivityLog } from '../../types/activity'

interface ActivityTableExpandedProps {
	log: ActivityLog
	open: boolean
}

const formatValue = (value: unknown): string => {
	if (value === null || value === undefined) {
		return '-'
	}
	if (typeof value === 'object') {
		return JSON.stringify(value)
	}
	return String(value)
}

const sortObjectKeys = (obj: unknown): Record<string, unknown> | null => {
	if (!obj || typeof obj !== 'object') {
		return null
	}
	const parsed = obj as Record<string, unknown>
	return Object.keys(parsed)
		.sort()
		.reduce(
			(result, key) => {
				result[key] = parsed[key]
				return result
			},
			{} as Record<string, unknown>,
		)
}

const fieldLabels: Record<string, string> = {
	name: 'Наименование',
	quantity: 'Количество',
	notes: 'Примечание',
	status: 'Статус',
	consumer: 'Конечник',
	customer: 'Заказчик / Перекуп',
	manager: 'Менеджер / Помощник',
	bill: 'Счет в 1С',
	date: 'Дата',
	isBargaining: 'Тендер',
	isBudget: 'Бюджет',
}

const itemFields = ['name', 'quantity', 'notes']

const getFieldLabel = (key: string): string => {
	return fieldLabels[key] || key
}

export const ActivityTableExpanded = ({ log, open }: ActivityTableExpandedProps) => {
	const oldValuesObj = sortObjectKeys(log.oldValues)
	const newValuesObj = sortObjectKeys(log.newValues)

	const allKeys = new Set([
		...(oldValuesObj ? Object.keys(oldValuesObj) : []),
		...(newValuesObj ? Object.keys(newValuesObj) : []),
	])
	const sortedKeys = Array.from(allKeys).sort()

	let changes = sortedKeys.map(key => {
		const oldVal = oldValuesObj?.[key]
		const newVal = newValuesObj?.[key]
		return { key, oldVal, newVal }
	})

	if (log.entityType === 'order_item') {
		changes = changes.filter(c => itemFields.includes(c.key))
	}

	const hasChanges = oldValuesObj || newValuesObj
	const showEntity = log.entity || log.parentId || log.order

	return (
		<TableRow>
			<TableCell
				colSpan={6}
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
						{showEntity && (
							<Box sx={{ mb: 2 }}>
								{log.entity && (
									<Typography variant='body2' color='text.secondary'>
										Объект: <strong>{log.entity}</strong>
									</Typography>
								)}
								{log.order || log.parentId ? (
									<Typography variant='body2' color='text.secondary'>
										Заявка от <strong>{log.order || log.parentId}</strong>
									</Typography>
								) : null}
							</Box>
						)}

						{hasChanges ? (
							<Table size='small'>
								<TableHead>
									<MuiTableRow>
										<TableCell sx={{ fontWeight: 700, width: '40%' }}>Поле</TableCell>
										<TableCell sx={{ fontWeight: 700, width: '30%' }}>Было</TableCell>
										<TableCell sx={{ fontWeight: 700, width: '30%' }}>Стало</TableCell>
									</MuiTableRow>
								</TableHead>
								<TableBody>
									{changes.length > 0 ? (
										changes.map(({ key, oldVal, newVal }) => {
											const isChanged = oldVal !== newVal
											return (
												<MuiTableRow key={key}>
													<TableCell>
														<Typography sx={{ fontWeight: 600 }}>
															{getFieldLabel(key)}
														</Typography>
													</TableCell>
													<TableCell
														sx={{
															bgcolor: isChanged
																? 'rgba(255, 107, 107, 0.1)'
																: 'transparent',
															fontSize: 13,
														}}
													>
														{formatValue(oldVal)}
													</TableCell>
													<TableCell
														sx={{
															bgcolor: isChanged
																? 'rgba(0, 184, 148, 0.1)'
																: 'transparent',
															fontSize: 13,
														}}
													>
														{formatValue(newVal)}
													</TableCell>
												</MuiTableRow>
											)
										})
									) : (
										<MuiTableRow>
											<TableCell colSpan={3}>Нет данных об изменениях</TableCell>
										</MuiTableRow>
									)}
								</TableBody>
							</Table>
						) : (
							<Typography variant='body2' color='text.secondary'>
								Нет данных об изменениях
							</Typography>
						)}
					</Box>
				</Collapse>
			</TableCell>
		</TableRow>
	)
}
