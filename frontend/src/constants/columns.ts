export const Columns = Object.freeze([
	{
		field: 'client',
		filter: 'list' as const,
		label: 'Контрагент',
	},
	{
		field: 'date',
		filter: 'date' as const,
		label: 'Дата',
	},
	{
		field: 'manager',
		filter: 'list' as const,
		label: 'Менеджер / помощник',
	},
	{
		field: 'isBargaining',
		filter: 'bool' as const,
		label: 'Тендер',
	},
	{
		field: 'isBudget',
		filter: 'bool' as const,
		label: 'Бюджет',
	},
])
