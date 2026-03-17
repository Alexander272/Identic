export type Field = {
	name: string
	type: 'text' | 'date' | 'number' | 'bool' | 'list' | 'autocomplete'
	label: string
	isRequired?: boolean
}

export type Option = {
	id: string
	name: string
}
