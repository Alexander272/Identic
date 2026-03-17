import { type FC } from 'react'
import { DatePicker } from '@mui/x-date-pickers'
import { Controller, useFormContext } from 'react-hook-form'
import dayjs from 'dayjs'

import type { Field } from './type'
import { NullDate } from '@/constants/defaultValues'
import { DateTextField } from '@/components/DatePicker/DatePicker'

type Props = {
	data: Field
	disabled?: boolean
}

export const DateField: FC<Props> = ({ data, disabled }) => {
	const { control } = useFormContext()

	return (
		<Controller
			control={control}
			name={data.name}
			rules={{ required: data.isRequired }}
			render={({ field, fieldState: { error } }) => (
				<DatePicker
					{...field}
					value={field.value != NullDate ? dayjs(field.value) : null}
					onChange={value => {
						field.onChange(value?.startOf('d').toISOString())
					}}
					label={data.label}
					showDaysOutsideCurrentMonth
					fixedWeekNumber={6}
					minDate={dayjs('2015-01-01')}
					maxDate={dayjs('2045-01-01')}
					slots={{
						textField: DateTextField,
					}}
					disabled={disabled}
					slotProps={{
						textField: {
							error: Boolean(error),
							fullWidth: true,
						},
					}}
				/>
			)}
		/>
	)
}
