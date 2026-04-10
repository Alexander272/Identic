import { useState, type FC } from 'react'
import { Paper, Typography, Box, Stack, ButtonGroup, Button, useTheme, useMediaQuery } from '@mui/material'
import { Link } from 'react-router'

import type { IOrderMatchResult, ISearchItem } from '../types/search'
import type { IPosition } from '@/features/orders/types/positions'
import { normalize } from '../utils/normalize'
import { VerifiedIcon } from '@/components/Icons/VerifiedIcon'
import { WarnIcon } from '@/components/Icons/WarnIcon'
import { CloseRoundIcon } from '@/components/Icons/CloseRoundIcon'
import { PopupLinkIcon } from '@/components/Icons/PopupLinkIcon'
import { DiffText } from './MatchTable/DiffText'

// --- Константы цветов и стилей ---
const colors = {
	success: '#E8F5E9', // Светло-зеленый
	successBorder: '#4CAF50',
	warning: '#FFFDE7', // Светло-желтый
	warningBorder: '#d19d19',
	error: '#FFEBEE', // Светло-красный
	errorBorder: '#F44336',
	errorBg: '#ffced5',
}

const tableRowStyles = (color: string, borderColor: string) => ({
	backgroundColor: color,
	'&:not(:last-child)': { marginBottom: '8px' }, // Отступ между строками
	borderRadius: '8px', // Скругление углов строки
	borderLeft: `5px solid ${borderColor}`,
	boxShadow: '0 2px 4px rgba(0,0,0,0.05)', // Легкая тень
	display: 'flex', // Для гибкости ячеек
	alignItems: 'center',
	padding: '8px 16px',
	// '&:hover': { backgroundColor: color, opacity: 0.9 }, // Чтобы MUI hover не перекрывал цвет
})

// --- Иконки статусов (Упрощенные SvgIcon для примера) ---
const StatusIcon = ({ status }: { status: 'found' | 'partial' | 'not_found' }) => {
	if (status === 'found') {
		return (
			// <SvgIcon sx={{ color: colors.successBorder, fontSize: '1.2rem' }}>
			// 	<path d='M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41L9 16.17z' />
			// </SvgIcon>
			<VerifiedIcon fill={colors.successBorder} fontSize={'1.3rem'} overflow={'visible'} />
		) // ✅
	} else if (status === 'partial') {
		return (
			// <SvgIcon sx={{ color: colors.warningBorder, fontSize: '1.2rem' }}>
			// 	<path d='M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z' />
			// </SvgIcon>
			<WarnIcon fill={colors.warningBorder} fontSize={'1.3rem'} overflow={'visible'} />
		) // ⚠️
	} else {
		return (
			// <SvgIcon sx={{ color: colors.errorBorder, fontSize: '1.2rem' }}>
			// 	<path d='M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12 19 6.41z' />
			// </SvgIcon>
			<CloseRoundIcon fill={colors.errorBorder} fontSize={'1.3rem'} overflow={'visible'} />
		) // ❌
	}
}

type Props = {
	request: ISearchItem[]
	result: IOrderMatchResult
	foundPositions: IPosition[]
}

// --- Основной компонент ---
export const MatchTable: FC<Props> = ({ request, result, foundPositions }) => {
	const { palette, breakpoints } = useTheme()
	const isMobile = useMediaQuery(breakpoints.down('md'))

	const [filter, setFilter] = useState<'all' | 'found' | 'not_found'>('all')

	// Логика определения статуса и данных для строки
	const getRowData = (index: number) => {
		// 1. Берем то, что пользователь ввел в поиске
		const requestedItem = { ...request[index] }

		// 2. Ищем в результате поиска, есть ли связь для этого индекса
		const match = result.positions.find(p => p.reqId === index.toString())

		// 3. Если связи нет — сразу возвращаем "Не найдено"
		if (!match) {
			return {
				status: 'not_found',
				foundItem: null,
			}
		}

		// 4. Если связь есть, ищем детальные данные в массиве найденных позиций по ID
		const foundItem = foundPositions.find(p => p.id === match.id)

		// 5. Определяем статус на основе сравнения количества
		const reqQty = requestedItem.quantity ?? 0
		const foundQty = foundItem?.quantity ?? 0

		if (!foundItem) {
			// На случай, если в связях ID есть, а в самом массиве позиций почему-то нет
			return { status: 'not_found', foundItem: null }
		}

		const foundName = normalize(foundItem.name)
		const foundNotes = normalize(foundItem.notes)
		const reqName = requestedItem?.name ? normalize(requestedItem.name) : ''

		const isNameMatch = foundName === reqName
		const isNotesMatch = foundNotes === reqName

		const matchedFrom = isNameMatch ? 'name' : isNotesMatch ? 'notes' : null

		const diff = {
			name: !matchedFrom,
			qty: foundQty !== reqQty,
		}

		if (diff.name || diff.qty) {
			return {
				status: 'partial',
				foundItem,
				mismatch: diff,
				matchedFrom,
			}
		}
		// if (foundQty < reqQty || foundQty > reqQty) {
		// 	return { status: 'partial', foundItem }
		// }

		return { status: 'found', foundItem }
	}

	// Фильтрация данных
	const filteredIndices = request.reduce((acc, _, index) => {
		const { status } = getRowData(index)
		if (
			filter === 'all' ||
			(filter === 'found' && status !== 'not_found') ||
			(filter === 'not_found' && status === 'not_found')
		) {
			acc.push(index)
		}
		return acc
	}, [] as number[])

	return (
		<Box>
			<Stack
				direction={{ md: 'row' }}
				spacing={{ md: 0, sm: 1 }}
				justifyContent='space-between'
				alignItems='center'
				mb={1.5}
			>
				<Typography fontSize={'1.2rem'} fontWeight='bold'>
					Детализация позиций поиска
				</Typography>

				<Stack direction='row' alignItems='center' gap={2}>
					<Link to={result.link} target='_blank' rel='noopener noreferrer'>
						<Button
							variant='outlined'
							size='small'
							sx={{ px: 2, textTransform: 'inherit' }}
							endIcon={<PopupLinkIcon fontSize={'14px !important'} fill={palette.primary.main} />}
						>
							Перейти к заказу
						</Button>
					</Link>

					<ButtonGroup variant='outlined' size='small'>
						<Button onClick={() => setFilter('all')} variant={filter === 'all' ? 'contained' : 'outlined'}>
							Все ({request.length})
						</Button>
						<Button
							onClick={() => setFilter('found')}
							variant={filter === 'found' ? 'contained' : 'outlined'}
						>
							Найдено ({result.matchedPos})
						</Button>
						<Button
							onClick={() => setFilter('not_found')}
							variant={filter === 'not_found' ? 'contained' : 'outlined'}
						>
							Не найдено ({request.length - result.matchedPos})
						</Button>
					</ButtonGroup>
				</Stack>
			</Stack>

			<Box
				sx={{
					maxHeight: 700,
					overflow: 'auto',
					display: 'flex',
					flexDirection: 'column',
					pr: 0.5,
				}}
			>
				{filteredIndices.map(index => {
					const item = request[index]
					const { status, foundItem, mismatch, matchedFrom } = getRowData(index as number)

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
							key={index}
							sx={{
								...rowColorStyles,
								display: 'flex',
								flexDirection: isMobile ? 'column' : 'row',
								alignItems: isMobile ? 'flex-start' : 'center',
								p: isMobile ? 2 : '8px 16px',
								minHeight: isMobile ? 'auto' : '64px',
								gap: isMobile ? 1.5 : 0,
								boxShadow: 'none',
							}}
						>
							{/* --- БЛОК 1: СЛУЖЕБНЫЙ (Номер, Иконка, Статус-текст) --- */}
							<Stack
								direction='row'
								spacing={2}
								alignItems='center'
								sx={{ minWidth: isMobile ? '100%' : '70px' }}
							>
								<Typography variant='body2' color='textSecondary'>
									{index + 1}
								</Typography>

								<StatusIcon status={status as 'found' | 'partial' | 'not_found'} />

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
									borderTop:
										isMobile && status !== 'found' ? `1px dashed ${palette.divider}` : 'none',
									pt: isMobile && status !== 'found' ? 1 : 0,
								}}
							>
								{/* На десктопе выводим статусы здесь, как в оригинале */}
								{!isMobile && status === 'found' && (
									<Typography
										variant='body1'
										color={colors.successBorder}
										fontWeight='bold'
										textAlign={'center'}
									>
										— Полное совпадение —
									</Typography>
								)}
								{!isMobile && status === 'not_found' && (
									<Typography
										variant='body1'
										color={colors.errorBorder}
										fontWeight='bold'
										textAlign={'center'}
									>
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
												actual={matchedFrom === 'notes' ? foundItem.notes : foundItem.name}
											/>
										) : (
											<Typography>
												{matchedFrom === 'notes' ? foundItem?.notes : foundItem?.name}
											</Typography>
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
				})}
			</Box>
		</Box>
	)
}
