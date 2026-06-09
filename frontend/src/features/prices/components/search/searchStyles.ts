export const cardSx = {
	borderRadius: { xs: 2, sm: 3 },
	paddingX: { xs: 1.5, sm: 3 },
	paddingY: 2,
	flex: 1,
	minWidth: 320,
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
}


