import { SearchType } from '../types/search'
import { ActionType, EntityType } from '../types/activity'

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
			return 'success' as const
		case ActionType.Update:
			return 'warning' as const
		case ActionType.Delete:
			return 'error' as const
		default:
			return 'default' as const
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

export const formatDate = (dateStr: string) => {
	const date = new Date(dateStr)
	return date.toLocaleString('ru-RU')
}

export const formatDuration = (ms: number) => {
	if (ms < 1000) return `${ms}мс`
	return `${(ms / 1000).toFixed(1)}с`
}