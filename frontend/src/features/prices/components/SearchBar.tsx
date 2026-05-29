import { type FC } from 'react'
import { TextField, Button, CircularProgress, Box } from '@mui/material'
import { SearchIcon } from '@/components/Icons/SearchIcon'

type SearchBarProps = {
	query: string
	onQueryChange: (v: string) => void
	onSearch: () => void
	isLoading: boolean
}

export const SearchBar: FC<SearchBarProps> = ({ query, onQueryChange, onSearch, isLoading }) => (
	<Box
		sx={{
			borderRadius: { xs: 2, sm: 3 },
			paddingX: { xs: 1.5, sm: 3 },
			paddingY: 2,
			width: '100%',
			mb: 2,

			position: 'relative',
			overflow: 'hidden',
			bgcolor: 'rgba(255,255,255,0.85)',
			border: '1px solid rgba(0,0,0,0.08)',
			backdropFilter: 'blur(20px)',
			boxShadow: '0 4px 12px rgba(0,0,0,0.04)',
			':before': {
				content: '""',
				position: 'absolute',
				top: 0,
				left: 0,
				right: 0,
				height: '1px',
				background: 'linear-gradient(90deg, transparent, #2563eb, transparent)',
			},
		}}
	>
		<TextField
			value={query}
			onChange={e => onQueryChange(e.target.value)}
			onKeyDown={e => {
				if (e.key === 'Enter') onSearch()
			}}
			label='Поиск'
			placeholder='Введите название, код, цену...'
			size='small'
			sx={{ flex: 1 }}
		/>
		<Button variant='contained' onClick={onSearch} disabled={isLoading} startIcon={<SearchIcon />}>
			{isLoading ? <CircularProgress size={20} /> : 'Найти'}
		</Button>
	</Box>
)
