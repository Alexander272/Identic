import { useState, type FC } from 'react'
import { createFilterOptions } from '@mui/material'

import type { Field } from '@/components/Form/type'
import { useGetUniqueDataQuery } from '../../orderApiSlice'
import { AutocompleteField } from '@/components/Form/AutocompleteField'

const filterOptions = createFilterOptions<string>({
	limit: 15,
})

export const AutocompleteInput: FC<{ field: Field }> = ({ field }) => {
	const [shouldFetch, setShouldFetch] = useState(false)

	const { data, isFetching } = useGetUniqueDataQuery({ field: field.name }, { skip: !shouldFetch })

	const onFocus = () => {
		setShouldFetch(true)
	}

	return (
		<AutocompleteField
			data={field}
			options={data?.data || []}
			filterOptions={filterOptions}
			isLoading={isFetching}
			onFocus={onFocus}
		/>
	)
}
