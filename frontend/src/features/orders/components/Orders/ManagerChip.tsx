import { useMemo } from 'react'
import { Chip } from '@mui/material'

import { stringToHSLA } from '@/utils/colors'

export const ManagerChip = ({ name }: { name: string }) => {
	const colors = useMemo(() => stringToHSLA(name), [name])

	return (
		<Chip
			label={name}
			size='small'
			style={{
				backgroundColor: colors.bg,
				color: colors.text,
				border: `1px solid ${colors.border}`,
				fontWeight: 500,
				fontSize: '0.75rem',
				height: '20px',
				borderRadius: '6px',
			}}
		/>
	)
}
