import { Paper, Typography, Box, Stack, useTheme } from '@mui/material'

import type { MatchRowData } from './useMatchTable'
import { StatusIcon } from './StatusIcon'
import { DiffText } from './DiffText'
import { colors } from './colors'

type Props = {
	row: MatchRowData
	isMobile: boolean
}

const tableRowStyles = (color: string, borderColor: string) => ({
	backgroundColor: color,
	// '&:not(:last-child)': { marginBottom: '8px' }, // Отступ между строками
	borderRadius: '8px', // Скругление углов строки
	borderLeft: `5px solid ${borderColor}`,
	boxShadow: '0 2px 4px rgba(0,0,0,0.05)', // Легкая тень
	display: 'flex', // Для гибкости ячеек
	alignItems: 'center',
	padding: '8px 16px',
	// '&:hover': { backgroundColor: color, opacity: 0.9 }, // Чтобы MUI hover не перекрывал цвет
})

export const MatchRow = ({ row, isMobile }: Props) => {
	const { palette } = useTheme()

	const { index, item, status, foundItem, mismatch, matchedFrom } = row

	let rowColorStyles
	let statusColor

	switch (status) {
		case 'found':
			rowColorStyles = tableRowStyles(colors.success, colors.successBorder)
			statusColor = colors.successBorder
			break
		case 'partial':
			rowColorStyles = tableRowStyles(colors.warning, colors.warningBorder)
			statusColor = colors.warningBorder
			break
		default:
			rowColorStyles = tableRowStyles(colors.error, colors.errorBorder)
			statusColor = colors.errorBorder
			break
	}

	return (
		<Paper
			sx={{
				...rowColorStyles,
				display: 'flex',
				flexDirection: isMobile ? 'column' : 'row',
				alignItems: isMobile ? 'flex-start' : 'center',
				p: isMobile ? 2 : '8px 16px',
				minHeight: isMobile ? 'auto' : '64px',
				gap: isMobile ? 1.5 : 0,
				boxShadow: 'none',
				mb: 1,
			}}
		>
			{/* --- БЛОК 1: СЛУЖЕБНЫЙ (Номер, Иконка, Статус-текст) --- */}
			<Stack direction='row' spacing={2} alignItems='center' sx={{ minWidth: isMobile ? '100%' : '70px' }}>
				<Typography variant='body2' color='textSecondary'>
					{index + 1}
				</Typography>
				<StatusIcon status={status} />

				{isMobile && (
					<Typography variant='caption' fontWeight='bold' color={statusColor}>
						{status === 'found' && '— Полное совпадение —'}
						{status === 'partial' && '— Частичное совпадение —'}
						{status === 'not_found' && '— Не найдено в заказе —'}
					</Typography>
				)}
			</Stack>

			{/* --- БЛОК 2: ЗАПРОС (Наименование/Кол-во пользователя) --- */}
			<Box sx={{ flex: 1, width: '100%' }}>
				{isMobile && (
					<Typography variant='caption' color='textSecondary'>
						Запрос:
					</Typography>
				)}
				<Typography>{item.name}</Typography>
				<Typography color='textSecondary'>{item.quantity ?? 0} шт.</Typography>
			</Box>

			{!isMobile && status === 'partial' && (
				<Typography
					variant='caption'
					color={colors.warningBorder}
					fontWeight='bold'
					px={'16px'}
					textAlign={'center'}
					sx={{ width: 100 }}
				>
					Частичное совпадение
				</Typography>
			)}

			{/* --- БЛОК 3: РЕЗУЛЬТАТ (Найденное в заказе) --- */}
			<Box
				sx={{
					flex: 1,
					width: '100%',
					// textAlign: isMobile ? 'left' : 'right',
					borderTop: isMobile && status !== 'found' ? `1px dashed ${palette.divider}` : 'none',
					pt: isMobile && status !== 'found' ? 1 : 0,
				}}
			>
				{/* На десктопе выводим статусы здесь, как в оригинале */}
				{!isMobile && status === 'found' && (
					<Typography variant='body1' color={colors.successBorder} fontWeight='bold' textAlign={'center'}>
						— Полное совпадение —
					</Typography>
				)}
				{!isMobile && status === 'not_found' && (
					<Typography variant='body1' color={colors.errorBorder} fontWeight='bold' textAlign={'center'}>
						— Не найдено в заказе —
					</Typography>
				)}

				{/* Содержимое для частичного совпадения */}
				{status === 'partial' && (
					<Stack>
						{isMobile && (
							<Typography variant='caption' color='textSecondary'>
								Найдено:
							</Typography>
						)}
						{mismatch?.name ? (
							<DiffText
								expected={item.name || ''}
								actual={matchedFrom === 'notes' ? foundItem?.notes || '' : foundItem?.name || ''}
							/>
						) : (
							<Typography>{matchedFrom === 'notes' ? foundItem?.notes : foundItem?.name}</Typography>
						)}
						<Typography color='textSecondary'>
							<Box
								component='span'
								sx={{
									p: '1px 4px',
									borderRadius: '4px',
									bgcolor: mismatch?.qty ? colors.errorBg : 'transparent',
									color: mismatch?.qty ? '#000' : 'inherit',
									mr: 0.5,
								}}
							>
								{foundItem?.quantity ?? 0}
							</Box>
							шт.
						</Typography>
					</Stack>
				)}
			</Box>
		</Paper>
	)
}
