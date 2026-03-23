import { useState, type FC } from 'react'
import {
	Table,
	TableBody,
	TableCell,
	TableContainer,
	TableRow,
	Paper,
	Typography,
	Box,
	Stack,
	ButtonGroup,
	Button,
} from '@mui/material'

import type { IOrderMatchResult, ISearchItem } from '../types/search'
import type { IPosition } from '@/features/orders/types/positions'
import { VerifiedIcon } from '@/components/Icons/VerifiedIcon'
import { WarnIcon } from '@/components/Icons/WarnIcon'
import { CloseRoundIcon } from '@/components/Icons/CloseRoundIcon'
import { DiffText } from './DiffText'

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
	const [filter, setFilter] = useState<'all' | 'found' | 'not_found'>('all')

	// Логика определения статуса и данных для строки
	const getRowData = (index: number) => {
		// 1. Берем то, что пользователь ввел в поиске
		const requestedItem = request[index]

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

		const diff = {
			name: foundItem.name.trim() !== requestedItem.name?.trim(),
			qty: foundQty !== reqQty,
		}

		if (diff.name || diff.qty) {
			return {
				status: 'partial',
				foundItem,
				mismatch: diff,
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
			<Stack direction='row' justifyContent='space-between' alignItems='center' mb={1}>
				<Typography fontSize={'1.2rem'} fontWeight='bold'>
					Детализация позиций поиска
				</Typography>
				<ButtonGroup variant='outlined' size='small'>
					<Button onClick={() => setFilter('all')} variant={filter === 'all' ? 'contained' : 'outlined'}>
						Все ({request.length})
					</Button>
					<Button onClick={() => setFilter('found')} variant={filter === 'found' ? 'contained' : 'outlined'}>
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

			<TableContainer
				component={Paper}
				sx={{
					boxShadow: 'none',
					border: 'none',
					// backgroundColor: 'transparent',
					maxHeight: 700,
					overflow: 'auto',
					pr: 0.5,
				}}
			>
				<Table sx={{ borderCollapse: 'separate', borderSpacing: '0 8px' }}>
					{/* Отступы между строками */}
					<TableBody>
						{filteredIndices.map(index => {
							const item = request[index]
							const { status, foundItem, mismatch } = getRowData(index as number)

							let rowStyle
							if (status === 'found') rowStyle = tableRowStyles(colors.success, colors.successBorder)
							else if (status === 'partial')
								rowStyle = tableRowStyles(colors.warning, colors.warningBorder)
							else rowStyle = tableRowStyles(colors.error, colors.errorBorder)

							return (
								<TableRow key={index} sx={rowStyle}>
									{/* # Ячейка с индексом */}
									<TableCell sx={{ border: 'none', width: '30px', padding: '0 8px' }}>
										<Typography variant='body2' color='textSecondary'>
											{index + 1}
										</Typography>
									</TableCell>

									<TableCell
										sx={{
											border: 'none',
											width: '30px',
											padding: '0',
											flexShrink: 0,
											display: 'flex',
											alignItems: 'center',
											justifyContent: 'center',
											alignSelf: 'stretch', // Ячейка займет всю высоту строки
										}}
									>
										<StatusIcon status={status as 'found' | 'partial' | 'not_found'} />
									</TableCell>

									{/* Ячейка с запросом */}
									<TableCell sx={{ border: 'none', width: 750, padding: '0 16px' }}>
										<Box>
											<Typography variant='body1' fontWeight='bold'>
												{item.name}
											</Typography>
											<Typography variant='body2' color='textSecondary'>
												{item.quantity ?? 0} шт.
											</Typography>
										</Box>
									</TableCell>

									<TableCell sx={{ border: 'none', width: '30px', padding: '0 8px' }}>
										{status === 'partial' ? (
											<Typography variant='body1' color='textSecondary'>
												→
											</Typography>
										) : null}
									</TableCell>

									{status === 'not_found' ? (
										<TableCell
											colSpan={2}
											sx={{ border: 'none', flex: 1, padding: '0 16px', textAlign: 'center' }}
										>
											<Typography variant='body1' color={colors.errorBorder} fontWeight='bold'>
												— Не найдено в заказе —
											</Typography>
										</TableCell>
									) : null}
									{status === 'found' ? (
										<TableCell
											colSpan={2}
											sx={{ border: 'none', flex: 1, padding: '0 16px', textAlign: 'center' }}
										>
											<Typography variant='body1' color={colors.successBorder} fontWeight='bold'>
												— Полное совпадение —
											</Typography>
										</TableCell>
									) : null}

									{/* Ячейка с результатом (или "Не найдено") */}
									{status === 'partial' ? (
										<>
											<TableCell
												sx={{
													border: 'none',
													width: 200,
													textAlign: 'center',
												}}
											>
												{/* <Typography variant='body1' fontWeight='bold'>
													Найдено:
												</Typography> */}

												{status === 'partial' && (
													<Typography
														variant='caption'
														color={colors.warningBorder}
														fontWeight='bold'
													>
														Частичное совпадение
													</Typography>
												)}
											</TableCell>

											<TableCell
												sx={{
													border: 'none',
													width: 750,
													padding: '0 16px',
													textAlign: 'right',
												}}
											>
												<Stack
													direction='row'
													alignItems='center'
													spacing={1}
													justifyContent='flex-end'
												>
													<Box textAlign='right'>
														{mismatch?.name ? (
															<DiffText
																expected={request[index].name || ''} // То, что ввел юзер в строку
																actual={foundItem.name} // То, что вернул findOrders
															/>
														) : (
															<Typography variant='body2'>{item.name}</Typography>
														)}
														<Typography variant='body2' color='textSecondary'>
															<Typography
																component={'span'}
																mr={0.5}
																sx={{
																	px: '1px',
																	borderRadius: '4px',
																	backgroundColor: mismatch?.qty
																		? colors.errorBg
																		: 'transparent',
																	color: mismatch?.qty ? '#000' : undefined,
																}}
															>
																{foundItem?.quantity ?? 0}
															</Typography>
															шт.
														</Typography>
														{/* <Typography variant='body1' fontWeight='bold'>
															Найдено: {item.name} ({foundQuant} шт.)
														</Typography> */}
													</Box>
												</Stack>
											</TableCell>
										</>
									) : null}
								</TableRow>
							)
						})}
					</TableBody>
				</Table>
			</TableContainer>
		</Box>
	)
}
