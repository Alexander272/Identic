import { useState, useCallback, useRef, useEffect, useMemo, type FC } from 'react'
import { Box, Chip, Popover } from '@mui/material'
import { toast } from 'react-toastify'

const chipSx = {
	fontFamily: 'monospace',
	height: '26px',
	minWidth: '26px',
	backgroundColor: '#f1f5f9',
	border: 'none',
	borderRadius: 2,
	svg: { fill: '#24254ec4' },
}

type CodesInputProps = {
	codes: string[]
	onChange: (codes: string[]) => void
}

export const CodesInput: FC<CodesInputProps> = ({ codes, onChange }) => {
	const [value, setValue] = useState('')
	const [overflowAnchorEl, setOverflowAnchorEl] = useState<null | HTMLElement>(null)
	const [containerWidth, setContainerWidth] = useState(0)
	const inputRef = useRef<HTMLInputElement>(null)
	const containerRef = useRef<HTMLDivElement>(null)

	useEffect(() => {
		const el = containerRef.current
		if (!el) return
		const ro = new ResizeObserver(entries => {
			setContainerWidth(entries[0].contentRect.width)
		})
		ro.observe(el)
		return () => ro.disconnect()
	}, [])

	const { hiddenChips, visibleChips } = useMemo(() => {
		const PADDING = 28
		const GAP = 4
		const INPUT_MIN = 80
		const CHIP_BASE = 36
		const CHIP_PER_CHAR = 10.5
		const PLUS_N_W = 52

		if (codes.length === 0 || containerWidth <= 0) {
			return { hiddenChips: [] as string[], visibleChips: [...codes] }
		}

		const tryFit = (reservePlusN: boolean) => {
			let available = containerWidth - PADDING - INPUT_MIN
			if (reservePlusN) available -= PLUS_N_W + GAP
			const visible: string[] = []
			for (let i = codes.length - 1; i >= 0; i--) {
				const chipW = codes[i].length * CHIP_PER_CHAR + CHIP_BASE + GAP
				if (available >= chipW) {
					visible.unshift(codes[i])
					available -= chipW
				} else {
					break
				}
			}
			return visible
		}

		const noPlus = tryFit(false)
		if (noPlus.length === codes.length) {
			return { hiddenChips: [], visibleChips: codes }
		}

		const withPlus = tryFit(true)
		return {
			hiddenChips: codes.slice(0, codes.length - withPlus.length),
			visibleChips: withPlus,
		}
	}, [codes, containerWidth])

	const commitValue = useCallback(() => {
		const trimmed = value.trim()
		if (trimmed && /^\d+$/.test(trimmed) && !codes.includes(trimmed)) {
			onChange([...codes, trimmed])
			setValue('')
		}
	}, [value, codes, onChange])

	const addCode = useCallback(
		(raw: string) => {
			const trimmed = raw.trim()
			if (!trimmed || !/^\d+$/.test(trimmed)) return false
			if (codes.includes(trimmed)) {
				toast.warn('Такой код уже есть')
				return false
			}
			onChange([...codes, trimmed])
			return true
		},
		[codes, onChange],
	)

	const handlePaste = useCallback(
		(e: React.ClipboardEvent) => {
			const text = e.clipboardData.getData('text')
			if (!text) return
			const parts = text.split(/[^0-9]+/).filter(Boolean)
			if (parts.length > 1) {
				e.preventDefault()
				const newCodes = parts.filter(p => !codes.includes(p))
				const duplicates = parts.filter(p => codes.includes(p))
				if (duplicates.length > 0) {
					toast.warn(`Пропущено ${duplicates.length} повторяющихся кодов`)
				}
				onChange([...codes, ...newCodes])
			}
		},
		[codes, onChange],
	)

	const handleKeyDown = useCallback(
		(e: React.KeyboardEvent) => {
			if (e.key === 'Enter') {
				if (value.trim()) {
					e.preventDefault()
					if (addCode(value)) setValue('')
				}
			} else if (e.key === ' ' || e.key === ',') {
				e.preventDefault()
				if (addCode(value)) setValue('')
			}
			if (e.key === 'Backspace' && !value && codes.length > 0) {
				onChange(codes.slice(0, -1))
			}
		},
		[value, codes, addCode, onChange],
	)

	return (
		<>
			<Box
				ref={containerRef}
				sx={{
					display: 'flex',
					flexWrap: 'nowrap',
					gap: 0.5,
					alignItems: 'center',
					p: '6.5px 14px',
					border: '1px solid rgba(0,0,0,0.23)',
					borderRadius: 2,
					cursor: 'text',
					minHeight: 40,
					mb: 1,
					overflow: 'hidden',
					bgcolor: 'transparent',
					transition: 'all 0.2s ease-in-out',
					'&:hover': { borderColor: 'rgba(0,0,0,0.87)' },
					'&:focus-within': {
						borderColor: '#0432a5',
						borderWidth: 1.5,
					},
				}}
				onClick={() => inputRef.current?.focus()}
			>
				{hiddenChips.length > 0 && (
					<Chip
						label={`+${hiddenChips.length}`}
						size='small'
						variant='outlined'
						sx={{ cursor: 'pointer', fontWeight: 600 }}
						onClick={e => {
							e.stopPropagation()
							setOverflowAnchorEl(e.currentTarget)
						}}
					/>
				)}
				{visibleChips.map(code => (
					<Chip
						key={code}
						label={code}
						size='small'
						onDelete={() => onChange(codes.filter(c => c !== code))}
						sx={chipSx}
					/>
				))}
				<input
					ref={inputRef}
					value={value}
					onChange={e => setValue(e.target.value)}
					onBlur={commitValue}
					onPaste={handlePaste}
					onKeyDown={handleKeyDown}
					placeholder={codes.length === 0 ? 'Вставьте или введите коды...' : ''}
					style={{
						border: 'none',
						outline: 'none',
						flex: 1,
						minWidth: 80,
						fontSize: 14,
						background: 'transparent',
					}}
				/>
			</Box>
			<Popover
				open={Boolean(overflowAnchorEl)}
				anchorEl={overflowAnchorEl}
				onClose={() => setOverflowAnchorEl(null)}
				anchorOrigin={{ vertical: 'bottom', horizontal: 'left' }}
				transformOrigin={{ vertical: 'top', horizontal: 'left' }}
				slotProps={{
					paper: { sx: { maxHeight: 320, maxWidth: 500, width: '100%', mt: 0.5 } },
				}}
			>
				<Box
					sx={{
						display: 'grid',
						gridTemplateColumns: 'repeat(auto-fill, minmax(85px, 1fr))',
						gap: 0.5,
						px: 1,
						py: 1,
						minWidth: 200,
					}}
				>
					{hiddenChips.map(code => (
						<Chip
							key={code}
							label={code}
							size='small'
							onDelete={() => onChange(codes.filter(c => c !== code))}
							sx={chipSx}
						/>
					))}
				</Box>
			</Popover>
		</>
	)
}
