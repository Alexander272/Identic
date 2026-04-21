import { SearchType } from '../types/search'
import { ActionType, EntityType } from '../types/activity'

export const getInitials = (name: string) => {
	return name
		.split(' ')
		.map(n => n[0])
		.join('')
		.toUpperCase()
		.slice(0, 2)
}

export const getSearchTypeLabel = (type: string) => {
	switch (type) {
		case SearchType.Exact:
			return 'Точный'
		case SearchType.Fuzzy:
			return 'Нечеткий'
		default:
			return type
	}
}

export const getSearchTypeColor = (type: string) => {
	switch (type) {
		case SearchType.Exact:
			return 'primary' as const
		case SearchType.Fuzzy:
			return 'secondary' as const
		default:
			return 'default' as const
	}
}

export const getActionLabel = (action: string) => {
	switch (action) {
		case ActionType.Insert:
			return 'Создание'
		case ActionType.Update:
			return 'Изменение'
		case ActionType.Delete:
			return 'Удаление'
		default:
			return action
	}
}

export const getActionColor = (action: string) => {
	switch (action) {
		case ActionType.Insert:
			return 'rgba(0, 184, 148, 0.15)'
		case ActionType.Update:
			return 'rgba(9, 132, 227, 0.15)'
		case ActionType.Delete:
			return 'rgba(255, 107, 107, 0.15)'
		default:
			return 'rgba(0, 0, 0, 0.08)'
	}
}

export const getActionTextColor = (action: string) => {
	switch (action) {
		case ActionType.Insert:
			return '#00b894'
		case ActionType.Update:
			return '#0984e3'
		case ActionType.Delete:
			return '#ff6b6b'
		default:
			return 'inherit'
	}
}

export const getEntityTypeLabel = (type: string) => {
	switch (type) {
		case EntityType.Order:
			return 'Заказ'
		case EntityType.OrderItem:
			return 'Позиция заказа'
		default:
			return type
	}
}

export const formatDuration = (ms: number) => {
	if (ms < 1000) return `${ms}мс`
	return `${(ms / 1000).toFixed(1)}с`
}

import dayjs from 'dayjs'

export const getDateRange = (period: string): { startDate: string; endDate: string } | undefined => {
	const end = dayjs().endOf('day')
	let start = dayjs().startOf('day')

	switch (period) {
		case 'today':
			// start уже равен текущему дню
			break
		case 'week':
			start = end.subtract(7, 'day')
			break
		case 'month':
			start = end.subtract(30, 'day')
			break
		case 'quarter':
			start = end.subtract(3, 'month')
			break
		case 'year':
			start = end.subtract(1, 'year')
			break
	}

	return {
		startDate: start.toISOString(),
		endDate: end.toISOString(),
	}
}
