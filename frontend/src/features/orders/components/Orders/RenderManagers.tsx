import { ManagerChip } from './ManagerChip'

export const renderManagers = (fullString: string) => {
	if (!fullString) return '—'

	// Разделяем по слэшу, запятой или пробелу
	const names = fullString
		.split(/[/,]/)
		.map(n => n.trim())
		.sort()

	return (
		<div style={{ display: 'flex', flexWrap: 'wrap', gap: '4px' }}>
			{names.map(name => (
				<ManagerChip key={name} name={name} />
			))}
		</div>
	)
}
