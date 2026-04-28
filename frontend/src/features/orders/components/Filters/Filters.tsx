import { useState, type FC } from 'react'
import { Button, Popover, Stack, Tooltip, Typography, useTheme } from '@mui/material'
import { FormProvider, useFieldArray, useForm } from 'react-hook-form'

import type { IFilter } from '../../types/params'
import { Columns } from '@/constants/columns'
import { Badge } from '@/components/Badge/Badge'
import { FilterIcon } from '@/components/Icons/FilterIcon'
import { TimesIcon } from '@/components/Icons/TimesIcon'
import { PlusIcon } from '@/components/Icons/PlusIcon'
import { CheckIcon } from '@/components/Icons/CheckSimpleIcon'
import { FilterItem } from './Item'

const defaultValue: IFilter = {
	field: Columns[0].field,
	fieldType: Columns[0].filter,
	compareType: 'in',
	value: '',
}

type Props = {
	filters: IFilter[]
	onChange: (data: IFilter[]) => void
}

export const Filters: FC<Props> = ({ filters, onChange }) => {
	const [anchor, setAnchor] = useState<HTMLButtonElement | null>(null)

	const { palette } = useTheme()

	const methods = useForm<{ filters: IFilter[] }>({
		values: {
			filters: filters.length ? filters : [defaultValue],
		},
	})
	const { fields, append, remove } = useFieldArray({ control: methods.control, name: 'filters' })

	const openHandler = (event: React.MouseEvent<HTMLButtonElement>) => setAnchor(event.currentTarget)
	const closeHandler = () => setAnchor(null)

	const addNewHandler = () => {
		append(defaultValue)
	}
	const removeHandler = (index: number) => {
		remove(index)
	}

	const resetHandler = () => {
		onChange([])
		methods.reset({ filters: [defaultValue] }) // Сбрасываем и состояние формы
		closeHandler()
	}

	// Чистая обработка данных без мутаций
	const onSubmit = methods.handleSubmit(form => {
		// const processedFilters = form.filters.map(f => ({
		// 	...f,
		// 	field: f.field.split('@')[0],
		// }))

		console.log('form', form)

		onChange(form.filters)
		closeHandler()
	})

	return (
		<>
			<Button
				onClick={openHandler}
				sx={{ minWidth: 40, borderRadius: '6px', ':hover': { svg: { fill: palette.primary.main } } }}
			>
				<Badge
					color='primary'
					variant={filters.length < 2 ? 'dot' : 'standard'}
					badgeContent={filters.length}
					anchorOrigin={{ horizontal: 'left' }}
				>
					<FilterIcon sx={{ fontSize: 18 }} />
				</Badge>
			</Button>

			<Popover
				open={Boolean(anchor)}
				onClose={closeHandler}
				anchorEl={anchor}
				anchorOrigin={{
					vertical: 'bottom',
					horizontal: 'center',
				}}
				transformOrigin={{
					vertical: 'top',
					horizontal: 'center',
				}}
				slotProps={{
					paper: {
						elevation: 0,
						sx: {
							overflow: 'visible',
							filter: 'drop-shadow(0px 2px 8px rgba(0,0,0,0.32))',
							mt: 1,
							paddingX: 2,
							pt: 1.5,
							paddingBottom: 2,
							maxWidth: 750,
							width: '100%',
							'&:before': {
								content: '""',
								display: 'block',
								position: 'absolute',
								top: 0,
								right: 'calc(50% - 10px)',
								width: 10,
								height: 10,
								bgcolor: 'background.paper',
								transform: 'translate(-50%, -50%) rotate(45deg)',
								zIndex: 0,
							},
						},
					},
				}}
			>
				<Stack>
					<Stack
						direction={'row'}
						mb={2.5}
						justifyContent={'space-between'}
						alignItems={'center'}
						component={'form'}
						onSubmit={onSubmit}
					>
						<Typography fontSize={'1.1rem'}>Применить фильтр</Typography>

						<Stack direction={'row'} spacing={1} height={34}>
							<Button
								onClick={addNewHandler}
								variant='outlined'
								sx={{ minWidth: 40, padding: '5px 14px' }}
							>
								<PlusIcon fill={palette.primary.main} fontSize={14} />
							</Button>

							<Tooltip title='Сбросить фильтры' enterDelay={700}>
								<Button onClick={resetHandler} variant='outlined' color='inherit' sx={{ minWidth: 40 }}>
									<TimesIcon fill={'#212121'} fontSize={12} />
								</Button>
							</Tooltip>

							<Tooltip title='Применить фильтры' enterDelay={700}>
								<Button type='submit' variant='contained' sx={{ minWidth: 40, padding: '6px 12px' }}>
									<CheckIcon sx={{ fontSize: 20, fill: palette.common.white }} />
								</Button>
							</Tooltip>
						</Stack>
					</Stack>

					<Stack spacing={1.5}>
						<FormProvider {...methods}>
							{fields.map((f, i) => (
								<FilterItem
									key={f.id}
									index={i}
									onRemove={removeHandler}
									canRemove={fields.length > 1}
								/>
							))}
						</FormProvider>
					</Stack>
				</Stack>
			</Popover>
		</>
	)
}
