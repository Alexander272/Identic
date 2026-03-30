import { useEffect, type FC } from 'react'
import { Box, Button, Divider, Stack, TextField, Tooltip, Typography, useTheme } from '@mui/material'
import { Controller, FormProvider, useForm, useFormContext, useWatch } from 'react-hook-form'
import { DataSheetGrid, floatColumn, keyColumn, textColumn, type Column } from 'react-datasheet-grid'
import { toast } from 'react-toastify'
import dayjs from 'dayjs'

import type { IFetchError } from '@/app/types/error'
import type { IOrderUpdate } from '../../types/order'
import type { IPositionCreate } from '../../types/positions'
import { useGetOrderByIdQuery, useUpdateOrderMutation } from '../../orderApiSlice'
import { extractQuantity } from '@/utils/extract'
import { AddRow } from '@/components/DataSheet/AddRow'
import { ContextMenu } from '@/components/DataSheet/ContextMenu'
import { DateField } from '@/components/Form/DateField'
import { SaveIcon } from '@/components/Icons/SaveIcon'
import { RefreshIcon } from '@/components/Icons/RefreshIcon'
import { AutocompleteInput } from '@/components/Autocomplete/AutocompleteInput'
import { BoxFallback } from '@/components/Fallback/BoxFallback'

const defaultValues: IOrderUpdate = {
	id: '',
	customer: '',
	consumer: '',
	manager: '',
	bill: '',
	date: dayjs().startOf('d').toISOString(),
	notes: '',
	positions: [{ rowNumber: 1, name: '', quantity: null, notes: '' }],
}

const columns: Column<IPositionCreate>[] = [
	{
		...keyColumn<IPositionCreate, 'name'>('name', textColumn),
		title: 'Наименование',
		pasteValue: ({ rowData, value }) => {
			// 1. Если в колонке "Количество" уже есть данные (например, вставили 2 колонки из Excel),
			// то не пытаемся парсить наименование, просто обновляем имя.
			if (rowData.quantity) {
				return { ...rowData, name: value }
			}

			// 2. Если количество пустое, запускаем наш парсер
			const { name, quantity } = extractQuantity(value)

			// 3. Возвращаем ОБНОВЛЕННЫЙ объект всей строки
			return {
				...rowData,
				name: name,
				quantity: quantity ?? rowData.quantity, // берем из парсера или оставляем старое
			}
		},
	},
	{ ...keyColumn<IPositionCreate, 'quantity'>('quantity', floatColumn), title: 'Количество', width: 0.25 },
	{ ...keyColumn<IPositionCreate, 'notes'>('notes', textColumn), title: 'Примечание', width: 0.75 },
]

type Props = {
	orderId: string
}

export const EditOrderForm: FC<Props> = ({ orderId }) => {
	const methods = useForm<IOrderUpdate>({ defaultValues })
	const { control, reset } = methods

	const { data: order, isFetching } = useGetOrderByIdQuery(orderId, { skip: !orderId })

	useEffect(() => {
		if (order?.data) reset(order.data)
	}, [order, reset])

	return (
		<Stack>
			<Typography align='center' variant='h5' mt={1} mb={3}>
				Редактирование заказа
			</Typography>

			{isFetching ? <BoxFallback /> : null}

			<FormProvider {...methods}>
				<Stack direction={'row'} spacing={2} px={2} mb={2}>
					<Stack spacing={2} width={'50%'}>
						<Typography fontSize={'1.1rem'}>Контрагенты</Typography>

						<AutocompleteInput field={{ name: 'consumer', label: 'Конечник', type: 'list' }} />
						<AutocompleteInput field={{ name: 'customer', label: 'Заказчик', type: 'list' }} />
					</Stack>

					<Stack spacing={2} width={'50%'}>
						<Typography fontSize={'1.1rem'}>Детали заказа</Typography>

						<DateField data={{ name: 'date', label: 'Дата', isRequired: true, type: 'date' }} />
						<AutocompleteInput field={{ name: 'manager', label: 'Менеджер / Помощник', type: 'list' }} />
						<Controller
							control={control}
							name={`bill`}
							render={({ field }) => <TextField {...field} label='Счет в 1С' fullWidth />}
						/>
					</Stack>
				</Stack>
				<Controller
					control={control}
					name={`notes`}
					render={({ field }) => (
						<TextField {...field} label='Примечание' multiline minRows={2} sx={{ mx: 2 }} />
					)}
				/>

				<Divider sx={{ mt: 3, mb: 2 }} />

				<Typography align='center' mb={1}>
					Позиции
				</Typography>

				<Grid />
			</FormProvider>
		</Stack>
	)
}

const Grid = () => {
	const { palette } = useTheme()

	const [update, { isLoading }] = useUpdateOrderMutation()

	const {
		control,
		setValue,
		handleSubmit,
		reset,
		formState: { isDirty },
	} = useFormContext<IOrderUpdate>()
	const positions = useWatch({ control, name: 'positions' }) || []

	const saveHandler = handleSubmit(async form => {
		if (!isDirty) {
			toast.error('Данные заказа не заполнены')
			return
		}
		if (form.consumer == '' && form.customer == '') {
			toast.error('Заполните хотя бы одного контрагента')
			return
		}

		if (form.positions.some(item => item.quantity === null || item.name === '')) {
			toast.error('Заполните хотя бы одну позицию')
			return
		}

		form.positions.forEach((element, idx) => {
			element.rowNumber = idx + 1
		})
		console.log(form)

		try {
			await update(form).unwrap()
			toast.success('Заказ успешно обновлен')
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	const clearHandler = () => {
		reset(defaultValues)
	}

	return (
		<Stack position={'relative'}>
			<DataSheetGrid
				value={positions}
				onChange={newValue => setValue('positions', newValue, { shouldDirty: true })}
				columns={columns}
				contextMenuComponent={props => <ContextMenu {...props} />}
				addRowsComponent={props => <AddRow {...props} />}
				autoAddRow
			/>
			<Stack direction={'row'} spacing={1} sx={{ position: 'absolute', right: 8, bottom: 6 }}>
				<Button
					onClick={saveHandler}
					color='inherit'
					disabled={!isDirty || isLoading}
					sx={{
						minWidth: 48,
						textTransform: 'inherit',
						background: '#fff',
						border: '1px solid #dcdcdc',
						borderRadius: '6px',
						padding: '4px 10px',
						':disabled': { svg: { fill: palette.action.disabled } },
						':hover': { svg: { fill: palette.primary.main }, color: palette.primary.main },
					}}
				>
					<SaveIcon fontSize={18} mr={1} />
					Сохранить
				</Button>

				<Tooltip title={!isDirty || isLoading ? 'Очистить форму' : ''}>
					<Box>
						<Button
							onClick={clearHandler}
							disabled={!isDirty || isLoading}
							sx={{
								minWidth: 48,
								height: '100%',
								background: '#fff',
								border: '1px solid #dcdcdc',
								borderRadius: '6px',
								padding: '4px 10px',
								':disabled': { svg: { fill: palette.action.disabled } },
								':hover': { svg: { fill: palette.secondary.main } },
							}}
						>
							<RefreshIcon fontSize={18} />
						</Button>
					</Box>
				</Tooltip>
			</Stack>
		</Stack>
	)
}
