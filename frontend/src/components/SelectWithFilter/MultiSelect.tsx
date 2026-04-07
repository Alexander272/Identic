// import { type FC, useMemo, useRef, useState, useCallback, forwardRef } from 'react'
// import {
// 	Box,
// 	ClickAwayListener,
// 	InputAdornment,
// 	Popper,
// 	Stack,
// 	type SxProps,
// 	TextField,
// 	type Theme,
// 	useTheme,
// 	useAutocomplete,
// 	type UseAutocompleteProps,
// 	Divider,
// 	CircularProgress,
// 	alpha,
// } from '@mui/material'
// import { Virtuoso } from 'react-virtuoso'

// // Типизация с поддержкой ID для стабильных ключей
// export type Option = {
// 	label: string
// 	id?: string | number
// }

// type Props<T extends Option> = {
// 	label?: string
// 	headerLabel?: string
// 	values: T[]
// 	options: T[]
// 	disabled?: boolean
// 	onChange: (values: T[]) => void
// 	isLoading?: boolean
// 	onFocus?: () => void
// 	filterOptions?: UseAutocompleteProps<T, true, true, false>['filterOptions']
// 	sx?: SxProps<Theme>
// }

// export const MultiSelect = <T extends Option>({
// 	label = 'Значение',
// 	headerLabel,
// 	values,
// 	options,
// 	onChange,
// 	disabled,
// 	isLoading,
// 	onFocus,
// 	filterOptions,
// 	sx,
// }: Props<T>) => {
// 	const theme = useTheme()
// 	const [anchor, setAnchor] = useState<HTMLInputElement | null>(null)
// 	const inputRef = useRef<HTMLInputElement>(null)

// 	const openHandler = (e: React.MouseEvent<HTMLInputElement>) => {
// 		if (disabled) return
// 		setAnchor(e.currentTarget)
// 		if (onFocus) onFocus()
// 	}
// 	const closeHandler = () => setAnchor(null)

// 	// Гарантируем уникальность опций в списке
// 	const sortedOptions = useMemo(() => {
// 		const selectedIds = new Set(values.map(v => v.id ?? v.label))
// 		return [...values, ...options.filter(o => !selectedIds.has(o.id ?? o.label))]
// 	}, [options, values])

// 	const { getRootProps, getInputProps, getListboxProps, getOptionProps, groupedOptions, inputValue, onInputChange } =
// 		useAutocomplete<T, true, true, false>({
// 			id: 'multi-select-virtuoso',
// 			multiple: true,
// 			value: values,
// 			options: sortedOptions,
// 			disabled,
// 			filterOptions,
// 			getOptionLabel: option => option.label,
// 			isOptionEqualToValue: (option, value) => (option.id ?? option.label) === (value.id ?? value.label),
// 			onChange: (_, newValue) => onChange(newValue as T[]),
// 			onClose: (_, reason) => {
// 				if (reason === 'escape') closeHandler()
// 			},
// 			onInputChange: (_, newValue) => onInputChange(newValue),
// 		})

// 	// Предотвращаем сброс поиска при клике на элемент (удобно для мультивыбора)
// 	const handleInputChange = useCallback(
// 		(e: React.ChangeEvent<HTMLInputElement>) => {
// 			onInputChange(e.target.value)
// 		},
// 		[onInputChange],
// 	)

// 	return (
// 		<Stack width='100%' sx={sx}>
// 			<TextField
// 				label={label}
// 				value={values.map(v => v.label).join(', ')}
// 				onClick={openHandler}
// 				disabled={disabled}
// 				slotProps={{
// 					htmlInput: {
// 						readOnly: true,
// 						sx: { cursor: 'pointer', textOverflow: 'ellipsis' },
// 					},
// 					inputLabel: { shrink: Boolean(anchor) || values.length > 0 },
// 				}}
// 			/>

// 			<Popper
// 				open={Boolean(anchor)}
// 				anchorEl={anchor}
// 				placement='bottom-start'
// 				sx={{
// 					width: anchor?.clientWidth ? anchor.clientWidth + 2 : 'auto',
// 					border: `1px solid ${theme.palette.divider}`,
// 					boxShadow: theme.shadows[8],
// 					backgroundColor: theme.palette.background.paper,
// 					borderRadius: 1,
// 					zIndex: theme.zIndex.modal + 10,
// 				}}
// 			>
// 				<ClickAwayListener onClickAway={closeHandler}>
// 					<div {...getRootProps()}>
// 						{headerLabel && (
// 							<Box
// 								sx={{
// 									p: 1.5,
// 									fontWeight: 600,
// 									borderBottom: `1px solid ${theme.palette.divider}`,
// 									fontSize: 13,
// 								}}
// 							>
// 								{headerLabel}
// 							</Box>
// 						)}

// 						<TextField
// 							inputRef={inputRef}
// 							autoFocus
// 							fullWidth
// 							placeholder='Поиск...'
// 							value={inputValue}
// 							onChange={handleInputChange}
// 							sx={{ p: 1 }}
// 							slotProps={{
// 								input: {
// 									startAdornment: (
// 										<InputAdornment position='start'>
// 											{/* Замените на ваш SearchIcon */}
// 											<Box sx={{ fontSize: 18, ml: 1 }}>🔍</Box>
// 										</InputAdornment>
// 									),
// 								},
// 								htmlInput: {
// 									...getInputProps(), // Важно: пропсы от useAutocomplete для работы клавиш
// 								},
// 							}}
// 						/>

// 						<ListboxComponent
// 							{...getListboxProps()}
// 							options={groupedOptions as T[]}
// 							getOptionProps={getOptionProps}
// 							loading={isLoading}
// 						/>
// 					</div>
// 				</ClickAwayListener>
// 			</Popper>
// 		</Stack>
// 	)
// }

// type ListboxProps<T> = {
// 	options: readonly T[]
// 	getOptionProps: (props: { option: T; index: number }) => any
// 	loading?: boolean
// }

// const ListboxComponent = forwardRef<HTMLDivElement, ListboxProps<any>>(
// 	({ options, getOptionProps, loading, ...other }, ref) => {
// 		const theme = useTheme()
// 		const itemSize = 45
// 		const height = Math.min(options.length * itemSize, 350) + (options.length > 0 ? 10 : 0)

// 		if (loading) {
// 			return (
// 				<Stack alignItems='center' py={3} spacing={1}>
// 					<CircularProgress size={20} />
// 				</Stack>
// 			)
// 		}

// 		if (options.length === 0) {
// 			return <Box sx={{ p: 2, fontSize: 13, color: 'text.secondary' }}>Ничего не найдено</Box>
// 		}

// 		return (
// 			<div ref={ref} {...other}>
// 				<Virtuoso
// 					style={{ height }}
// 					totalCount={options.length}
// 					itemContent={index => {
// 						const option = options[index]
// 						const { key, ...optionProps } = getOptionProps({ option, index })
// 						const selected = optionProps['aria-selected']

// 						return (
// 							<Box key={key} {...optionProps} component='li' sx={{ listStyle: 'none', p: 0.5 }}>
// 								<Stack
// 									direction='row'
// 									alignItems='center'
// 									sx={{
// 										p: 1,
// 										borderRadius: 1,
// 										cursor: 'pointer',
// 										transition: '0.2s',
// 										'&:hover': { bgcolor: theme.palette.action.hover },
// 										...(selected && { bgcolor: alpha(theme.palette.primary.main, 0.08) }),
// 									}}
// 								>
// 									<Box
// 										sx={{
// 											width: 18,
// 											height: 18,
// 											mr: 1.5,
// 											borderRadius: 0.5,
// 											border: `1px solid ${selected ? theme.palette.primary.main : theme.palette.divider}`,
// 											display: 'flex',
// 											alignItems: 'center',
// 											justifyContent: 'center',
// 											bgcolor: selected ? theme.palette.primary.main : 'transparent',
// 										}}
// 									>
// 										{selected && <Box sx={{ color: '#fff', fontSize: 12 }}>✓</Box>}
// 									</Box>
// 									<Box sx={{ flex: 1, fontSize: 14, overflow: 'hidden', textOverflow: 'ellipsis' }}>
// 										{option.label}
// 									</Box>
// 								</Stack>
// 								{index < options.length - 1 && <Divider sx={{ mx: 1, mt: 0.5 }} />}
// 							</Box>
// 						)
// 					}}
// 				/>
// 			</div>
// 		)
// 	},
// )
