import React from 'react'
import { Box, Typography } from '@mui/material'
import { diffWords } from 'diff'

interface DiffTextProps {
	expected: string // То, что просил пользователь
	actual: string // То, что реально нашли в базе
}

export const DiffText: React.FC<DiffTextProps> = ({ expected, actual }) => {
	// Вычисляем разницу посимвольно
	const diff = diffWords(expected, actual)

	return (
		<Typography variant='body2'>
			{diff.map((part, index) => {
				// part.added — символы, которые есть в базе, но нет в запросе (лишние/другие)
				// part.removed — символы, которые были в запросе, но нет в базе

				// let color = 'inherit'
				let bgcolor = 'transparent'

				if (part.added) {
					// color = '#af0202' // Красный (ошибка/несоответствие)
					bgcolor = '#ffced5'
				} else if (part.removed) {
					// Обычно "удаленные" символы из запроса в итоговой строке
					// можно либо не показывать, либо показывать зачеркнутыми
					return null
				}

				return (
					<Box
						key={index}
						component='span'
						sx={{
							// color,
							backgroundColor: bgcolor,
							px: '1px',
							borderRadius: '4px',
							fontWeight: 'bold',
						}}
					>
						{part.value}
					</Box>
				)
			})}
		</Typography>
	)
}
