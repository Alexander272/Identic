import { type FC } from 'react'
import { Box, TextField, Button, Chip, IconButton } from '@mui/material'
import { RefreshIcon } from '@/components/Icons/RefreshIcon'

type FilterPanelProps = {
	codeInput: string
	onCodeInputChange: (v: string) => void
	onAddCode: () => void
	onRemoveCode: (code: string) => void
	onClearCodes: () => void
	codes: string[]
}

export const FilterPanel: FC<FilterPanelProps> = ({
	codeInput,
	onCodeInputChange,
	onAddCode,
	onRemoveCode,
	onClearCodes,
	codes,
}) => (
	<Box sx={{ mt: 2, display: 'flex', gap: 1, flexWrap: 'wrap', alignItems: 'center' }}>
		<TextField
			value={codeInput}
			onChange={e => onCodeInputChange(e.target.value)}
			onKeyDown={e => {
				if (e.key === 'Enter') {
					e.preventDefault()
					onAddCode()
				}
			}}
			label='Код позиции'
			size='small'
			sx={{ width: 200 }}
		/>
		<Button size='small' onClick={onAddCode}>
			Добавить
		</Button>
		<Box sx={{ display: 'flex', gap: 0.5, flexWrap: 'wrap', ml: 1 }}>
			{codes.map(code => (
				<Chip
					key={code}
					label={code}
					onDelete={() => onRemoveCode(code)}
					size='small'
					color='primary'
					variant='outlined'
				/>
			))}
			{codes.length > 0 && (
				<IconButton size='small' onClick={onClearCodes}>
					<RefreshIcon sx={{ fontSize: 16 }} />
				</IconButton>
			)}
		</Box>
	</Box>
)
