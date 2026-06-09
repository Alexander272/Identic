import { type Theme } from '@mui/material'

export const buttonSx = ({ palette }: Theme) => ({
	minWidth: 48,
	textTransform: 'inherit',
	background: '#fff',
	border: '1px solid #c3c3c4',
	borderRadius: '6px',
	padding: '4px 10px',
	':hover': { svg: { fill: palette.primary.main }, color: palette.primary.main },
	'&:disabled': { svg: { fill: palette.action.disabled } },
})
