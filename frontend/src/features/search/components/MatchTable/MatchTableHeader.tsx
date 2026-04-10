import type { FC } from 'react'
import { Stack, Typography, Button, ButtonGroup, useTheme } from '@mui/material'
import { Link } from 'react-router'

import { PopupLinkIcon } from '@/components/Icons/PopupLinkIcon'

type Props = {
	filter: string
	setFilter: (filter: 'all' | 'found' | 'not_found') => void
	total: number
	matched: number
	link: string
}

export const MatchTableHeader: FC<Props> = ({ filter, setFilter, total, matched, link }) => {
	const { palette } = useTheme()

	return (
		<Stack
			direction={{ md: 'row' }}
			spacing={{ md: 0, sm: 1 }}
			justifyContent='space-between'
			alignItems='center'
			mb={1.5}
		>
			<Typography fontSize={'1.2rem'} fontWeight='bold'>
				Детализация позиций поиска
			</Typography>

			<Stack direction='row' alignItems='center' gap={2}>
				<Link to={link} target='_blank'>
					<Button
						variant='outlined'
						size='small'
						sx={{ px: 2, textTransform: 'inherit' }}
						endIcon={<PopupLinkIcon fontSize={'14px !important'} fill={palette.primary.main} />}
					>
						Перейти к заказу
					</Button>
				</Link>

				<ButtonGroup variant='outlined' size='small'>
					<Button onClick={() => setFilter('all')} variant={filter === 'all' ? 'contained' : 'outlined'}>
						Все ({total})
					</Button>

					<Button onClick={() => setFilter('found')} variant={filter === 'found' ? 'contained' : 'outlined'}>
						Найдено ({matched})
					</Button>

					<Button
						onClick={() => setFilter('not_found')}
						variant={filter === 'not_found' ? 'contained' : 'outlined'}
					>
						Не найдено ({total - matched})
					</Button>
				</ButtonGroup>
			</Stack>
		</Stack>
	)
}
