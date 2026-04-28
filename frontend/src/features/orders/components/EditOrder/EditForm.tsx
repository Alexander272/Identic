import { useEffect, type FC } from 'react'
import {
	Box,
	Button,
	Checkbox,
	Divider,
	FormControlLabel,
	Stack,
	TextField,
	Tooltip,
	Typography,
	useTheme,
} from '@mui/material'
import { Controller, FormProvider, useFieldArray, useForm, useFormContext, useWatch } from 'react-hook-form'
import { DataSheetGrid, floatColumn, keyColumn, textColumn, type Column } from 'react-datasheet-grid'
import type { Operation } from 'react-datasheet-grid/dist/types'
import { toast } from 'react-toastify'
import dayjs from 'dayjs'

import './styles.css'

import type { IFetchError } from '@/app/types/error'
import type { IOrder, IOrderUpdate } from '../../types/order'
import type { IPositionUpdate } from '../../types/positions'
import { useGetOrderByIdQuery, useUpdateOrderMutation } from '../../orderApiSlice'
import { extractQuantity } from '@/utils/extract'
import { handleGlobalPaste } from '@/utils/globalPaste'
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
	isBargaining: false,
	isBudget: false,
	bill: '',
	date: dayjs().startOf('d').toISOString(),
	notes: '',
	positions: [
		{
			id: new Date().getTime().toString(),
			rowNumber: 1,
			name: '',
			quantity: null,
			notes: '',
			status: 'CREATED',
		},
	],
}

const columns: Column<IPositionUpdate>[] = [
	{
		...keyColumn<IPositionUpdate, 'name'>('name', textColumn),
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
	{ ...keyColumn<IPositionUpdate, 'quantity'>('quantity', floatColumn), title: 'Количество', width: 0.25 },
	{ ...keyColumn<IPositionUpdate, 'notes'>('notes', textColumn), title: 'Примечание', width: 0.75 },
]

type Props = {
	orderId: string
}

export const EditOrderForm: FC<Props> = ({ orderId }) => {
	const methods = useForm<IOrderUpdate>({ defaultValues })
	const { control, reset } = methods
	const consumer = useWatch({ control, name: 'consumer' })

	const { data: order, isFetching } = useGetOrderByIdQuery({ id: orderId }, { skip: !orderId })

	useEffect(() => {
		if (order?.data) reset(order.data)
	}, [order, reset])

	useEffect(() => {
		// true — использование фазы захвата (capture)
		window.addEventListener('paste', handleGlobalPaste, true)

		return () => {
			window.removeEventListener('paste', handleGlobalPaste, true)
		}
	}, [])

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
						<AutocompleteInput field={{ name: 'customer', label: 'Заказчик / Перекуп', type: 'list' }} />

						<Stack direction={'row'} spacing={1}>
							<Controller
								control={control}
								name='isBargaining'
								disabled={!consumer}
								render={({ field }) => (
									<FormControlLabel
										control={<Checkbox {...field} checked={field.value} />}
										label='Тендер'
										sx={{
											pr: 2,
											borderRadius: '8px',
											transition: 'all 0.3s ease-in-out',
											':hover': { bgcolor: 'action.hover' },
										}}
									/>
								)}
							/>
							<Controller
								control={control}
								name='isBudget'
								disabled={!consumer}
								render={({ field }) => (
									<FormControlLabel
										control={<Checkbox {...field} checked={field.value} />}
										label='Бюджет'
										sx={{
											pr: 2,
											borderRadius: '8px',
											transition: 'all 0.3s ease-in-out',
											':hover': { bgcolor: 'action.hover' },
										}}
									/>
								)}
							/>
						</Stack>
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

				<Grid orderId={orderId} data={order?.data} />
			</FormProvider>
		</Stack>
	)
}

const Grid: FC<{ orderId: string; data?: IOrder }> = ({ orderId, data }) => {
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

	const { fields } = useFieldArray({
		control,
		name: 'positions',
		keyName: '_rhf_id',
	})

	const genId = () => new Date().getTime().toString()

	const createHandler = () => ({ ...defaultValues.positions[0], id: genId() })
	const duplicateHandler = ({ rowData }: { rowData: IPositionUpdate }) => ({ ...rowData, id: genId() })

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
		if (form.id == '') form.id = orderId

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
		reset(data)
	}

	const addClasses = ({ rowData }: { rowData: IPositionUpdate }) => {
		switch (rowData.status) {
			case 'DELETED':
				return 'row-deleted'
			case 'CREATED':
				return 'row-created'
			case 'UPDATED':
				return 'row-updated'
		}
	}

	const changeHandler = (value: IPositionUpdate[], operations: Operation[]) => {
		const updatedValue = [...value]

		for (const operation of operations) {
			if (operation.type === 'DELETE') {
				const deletedRows = fields.slice(operation.fromRowIndex, operation.toRowIndex)

				deletedRows.forEach((row, index) => {
					if (row.status == 'CREATED') {
						updatedValue.splice(operation.fromRowIndex, 1)
					} else {
						const deletedRow = { ...row, status: 'DELETED' as const }
						updatedValue.splice(operation.fromRowIndex + index, 0, deletedRow)
					}
				})
			}

			if (operation.type === 'UPDATE') {
				for (let i = operation.fromRowIndex; i < operation.toRowIndex; i++) {
					const row = updatedValue[i]
					// Если строка уже помечена как созданная или удаленная, статус не меняем
					if (!row.status) {
						updatedValue[i] = { ...row, status: 'UPDATED' }
					}
				}
			}

			if (operation.type === 'CREATE') {
				for (let i = operation.fromRowIndex; i < operation.toRowIndex; i++) {
					updatedValue[i] = { ...updatedValue[i], status: 'CREATED' }
				}
			}
		}

		setValue('positions', updatedValue, { shouldDirty: true })
	}

	return (
		<Stack position={'relative'} mb={1}>
			<DataSheetGrid
				value={positions}
				// onChange={newValue => setValue('positions', newValue, { shouldDirty: true })}
				createRow={createHandler}
				duplicateRow={duplicateHandler}
				onChange={changeHandler}
				rowClassName={addClasses}
				columns={columns}
				contextMenuComponent={props => <ContextMenu {...props} />}
				addRowsComponent={props => <AddRow {...props} />}
				autoAddRow
				height={300}
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
