import type { FC } from 'react'
import { Box } from '@mui/material'
import { TrendingUp, Add, Edit, Delete } from '@mui/icons-material'

import type { ActivityLog, ActionType, EntityType } from '../../types'
import { StatCard } from './StatCard'

type Props = {
	data: ActivityLog[]
	total: number
}

const countByActionAndEntity = (data: ActivityLog[], action: ActionType, entity: EntityType) =>
	data.filter(log => log.action === action && log.entityType === entity).length

export const OrderCards: FC<Props> = ({ data, total }) => {
	const createdOrders = countByActionAndEntity(data, 'INSERT', 'order')
	const updatedOrders = countByActionAndEntity(data, 'UPDATE', 'order')
	const addedItems =
		countByActionAndEntity(data, 'INSERT', 'order_item') + countByActionAndEntity(data, 'UPDATE', 'order_item')
	const modifiedItems = countByActionAndEntity(data, 'DELETE', 'order_item')

	return (
		<Box
			sx={{
				display: 'grid',
				gridTemplateColumns: {
					xs: '1fr',
					sm: 'repeat(2, 1fr)',
					lg: 'repeat(5, 1fr)',
				},
				gap: 3,
				mb: 3,
			}}
		>
			<StatCard icon={<TrendingUp />} title='Всего изменений' value={total} color='#fdcb6e' />

			<StatCard icon={<Add />} title='Создано заявок' value={createdOrders} color='#00b894' />

			<StatCard icon={<Edit />} title='Редактирований заявок' value={updatedOrders} color='#0984e3' />

			<StatCard icon={<TrendingUp />} title='Добавлено/Изменено позиций' value={addedItems} color='#6c5ce7' />

			<StatCard icon={<Delete />} title='Удалено позиций' value={modifiedItems} color='#ff6b6b' />
		</Box>
	)
}
