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
	// {
	// 	field: 'bill',
	// 	filter: 'string',
	// 	label: 'No счета',
	// },
	{
		field: 'manager',
		filter: 'list' as const,
		label: 'Менеджер / помощник',
	},
	// {
	// 	field: 'notes',
	// 	filter: 'string',
	// 	label: 'Примечание',
	// },
])
