import { Controller, useFormContext } from 'react-hook-form'

import type { IFilter } from '../../types/params'
import { Switch } from '@/components/Switch/Switch'
import type { FC } from 'react'

type Props = {
	index: number
}

export const SwitchFilter: FC<Props> = ({ index }) => {
	const { control, setValue } = useFormContext<{ filters: IFilter[] }>()

	const handleChange = (val: boolean) => {
		setValue(`filters.${index}.value`, val ? 'true' : 'false')
	}

	return (
		<>
			<Controller
				name={`filters.${index}.value`}
				control={control}
				defaultValue='false'
				render={({ field }) => (
					<Switch
						value={field.value === 'true'}
						onChange={handleChange}
						sx={{ height: 40, borderRadius: 3 }}
					/>
				)}
			/>
		</>
	)
}
