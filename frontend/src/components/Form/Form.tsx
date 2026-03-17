import type { FC } from 'react'

import type { Field, Option } from './type'
import { TextField } from './TextField'
import { NumberField } from './NumberField'
import { DateField } from './DateField'
import { SelectField } from './SelectField'

type Props = {
	item: Field
	options?: Option[]
	// isLoading?: boolean
	// onFocus?: () => void
}

export const FormField: FC<Props> = ({ item, options }) => {
	switch (item.type) {
		case 'text':
			return <TextField key={item.name} data={item} />
		case 'number':
			return <NumberField key={item.name} data={item} />
		case 'date':
			return <DateField key={item.name} data={item} />
		// case 'checkbox':
		// 	return <CheckboxField key={item.id} data={item} />
		// case 'file':
		// 	return <FileField key={item.id} data={item} instrumentId={instrumentId} />
		case 'list':
			return <SelectField key={item.name} data={item} options={options || []} />
		// case 'autocomplete':
		// 	return (
		// 		<AutocompleteField
		// 			key={item.name}
		// 			data={item}
		// 			options={options || []}
		// 			isLoading={isLoading}
		// 			onFocus={onFocus}
		// 		/>
		// 	)
		default:
			return null
	}
}
