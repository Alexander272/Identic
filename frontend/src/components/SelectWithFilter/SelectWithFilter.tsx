import { type FC, useMemo, useState } from 'react'
import {
	Autocomplete,
	autocompleteClasses,
	type AutocompleteCloseReason,
	Box,
	ClickAwayListener,
	InputAdornment,
	Popper,
	Stack,
	styled,
	type SxProps,
	TextField,
	type Theme,
	Tooltip,
	useTheme,
} from '@mui/material'

import { SearchIcon } from '../Icons/SearchIcon'
import { ListboxComponent } from './Listbox'

export type Option = {
	// id: string
	label: string
}

type Props = {
	label?: string
	headerLabel?: string
	values: Option[]
	options: Option[]
	disabled?: boolean
	onChange: (values: Option[]) => void
	isLoading?: boolean
	onFocus?: () => void
	// eslint-disable-next-line @typescript-eslint/no-explicit-any
	filterOptions?: (options: any[], { inputValue }: any) => any[]
	sx?: SxProps<Theme>
}

export const SelectWithFilter: FC<Props> = ({
	label,
	headerLabel,
	values,
	options,
	onChange,
	disabled,
	isLoading,
	onFocus,
	filterOptions,
	sx,
}) => {
	const theme = useTheme()
	const [anchor, setAnchor] = useState<HTMLInputElement | null>(null)
	const [inputValue, setInputValue] = useState('')

	const openHandler = (e: React.MouseEvent<HTMLInputElement>) => {
		setAnchor(e.target as HTMLInputElement)
	}
	const closeHandler = () => {
		setAnchor(null)
	}

	const sortedOptions = useMemo(() => {
		const selectedSet = new Set(values.map(v => v.label))

		return [...values, ...options.filter(o => !selectedSet.has(o.label))]
	}, [options, values])

	const fullValue = useMemo(() => {
		return values.map(v => v.label).join('; ')
	}, [values])

	return (
		<Stack width={'100%'} sx={sx}>
			<Tooltip title={fullValue} placement='top-start' disableInteractive>
				<TextField
					label={label || 'Значение'}
					value={fullValue}
					onClick={openHandler}
					disabled={disabled}
					slotProps={{
						htmlInput: {
							readOnly: true,
							sx: {
								cursor: 'pointer',
								overflow: 'hidden',
								textOverflow: 'ellipsis',
								whiteSpace: 'nowrap',
							},
						},
						inputLabel: {
							shrink: Boolean(anchor) || Boolean(values.length),
						},
					}}
				/>
			</Tooltip>

			<Popper
				open={Boolean(anchor)}
				anchorEl={anchor}
				placement='bottom-start'
				sx={{
					width: anchor ? anchor.clientWidth + 2 : 'auto',
					border: '1px solid #e1e4e8',
					boxShadow: `0 8px 24px rgba(149, 157, 165, 0.2)`,
					color: '#24292e',
					backgroundColor: '#fff',
					borderRadius: 2,
					zIndex: theme.zIndex.modal + 10,
					fontSize: 13,
				}}
			>
				<ClickAwayListener onClickAway={closeHandler}>
					<div>
						{headerLabel && (
							<Box
								sx={t => ({
									borderBottom: '1px solid #30363d',
									padding: '8px 10px',
									fontWeight: 600,
									...t.applyStyles('light', {
										borderBottom: '1px solid #eaecef',
									}),
								})}
							>
								{headerLabel}
							</Box>
						)}
						<Autocomplete
							open
							multiple
							onClose={(_event: React.ChangeEvent<object>, reason: AutocompleteCloseReason) => {
								if (reason === 'escape') {
									closeHandler()
								}
							}}
							value={values}
							inputValue={inputValue}
							onChange={(event, newValue, reason) => {
								if (
									event.type === 'keydown' &&
									((event as React.KeyboardEvent).key === 'Backspace' ||
										(event as React.KeyboardEvent).key === 'Delete') &&
									reason === 'removeOption'
								) {
									return
								}
								onChange(newValue)
							}}
							onInputChange={(_event, newInputValue, reason) => {
								if (reason !== 'input') return
								setInputValue(newInputValue)
							}}
							onFocus={onFocus}
							filterOptions={filterOptions}
							disableCloseOnSelect
							renderValue={() => null}
							loading={isLoading}
							loadingText='Поиск похожих значений...'
							noOptionsText='Ничего не найдено'
							renderOption={(props, option, state) => [props, option, state] as React.ReactNode}
							options={sortedOptions}
							getOptionLabel={option => option.label}
							renderInput={params => (
								<TextField
									ref={params.InputProps.ref}
									autoFocus
									fullWidth
									placeholder='Поиск'
									sx={{ padding: '8px 16px', borderBottom: '1px solid #eaecef' }}
									slotProps={{
										htmlInput: {
											...params.inputProps,
										},
										input: {
											startAdornment: (
												<InputAdornment position='start'>
													<SearchIcon fontSize={16} ml={1} />
												</InputAdornment>
											),
										},
									}}
								/>
							)}
							slotProps={{ listbox: { component: ListboxComponent } }}
							slots={{
								popper: PopperComponent,
							}}
						/>
					</div>
				</ClickAwayListener>
			</Popper>
		</Stack>
	)
}

interface PopperComponentProps {
	anchorEl?: unknown
	disablePortal?: boolean
	open: boolean
}

function PopperComponent(props: PopperComponentProps) {
	// eslint-disable-next-line @typescript-eslint/no-unused-vars
	const { disablePortal, anchorEl, open, ...other } = props
	return <StyledAutocompletePopper {...other} />
}

const StyledAutocompletePopper = styled('div')(({ theme }) => ({
	[`& .${autocompleteClasses.paper}`]: {
		boxShadow: 'none',
		color: 'inherit',
		fontSize: 13,
	},
	[`& .${autocompleteClasses.listbox}`]: {
		padding: 0,
		backgroundColor: '#fff',

		[`& .${autocompleteClasses.option}`]: {
			// flexDirection: 'column',
			// alignItems: 'flex-start',
			padding: 8,
			paddingBottom: 0,
			margin: 6,
			borderRadius: 8,

			[`&.${autocompleteClasses.focused}, &.${autocompleteClasses.focused}[aria-selected="true"]`]: {
				backgroundColor: theme.palette.action.hover,
			},
		},
		'& ul': {
			padding: 0,
			margin: 0,
		},
	},
	[`&.${autocompleteClasses.popperDisablePortal}`]: {
		position: 'relative',
	},
}))
