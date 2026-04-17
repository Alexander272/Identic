export const ActionType = {
	Insert: 'INSERT',
	Update: 'UPDATE',
	Delete: 'DELETE',
} as const

export type ActionType = (typeof ActionType)[keyof typeof ActionType]

export const EntityType = {
	Order: 'order',
	OrderItem: 'order_item',
} as const

export type EntityType = (typeof EntityType)[keyof typeof EntityType]

export interface ActivityLog {
	id: string
	action: ActionType
	changedBy: string
	changedByName: string
	entityType: EntityType
	entityId: string
	entity?: string | null
	parentId?: string | null
	oldValues?: JSON
	newValues?: JSON
	createdAt: string
}

export interface ActivityLogResponse {
	total: number
	data: ActivityLog[]
}
