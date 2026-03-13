import { useState } from 'react'
import { Button, Divider, Stack, TextField, Tooltip, Typography, useTheme } from '@mui/material'
import { DatePicker } from '@mui/x-date-pickers'
import { Controller, useForm } from 'react-hook-form'
import { DataSheetGrid, floatColumn, keyColumn, textColumn, type Column } from 'react-datasheet-grid'
import { toast } from 'react-toastify'
import dayjs from 'dayjs'

import type { IFetchError } from '@/app/types/error'
import type { IOrderCreate } from '../../types/order'
import type { IPositionCreate } from '../../types/positions'
import { useCreateOrderMutation } from '../../orderApiSlice'
import { extractQuantity } from '@/utils/extract'
import { DateTextField } from '@/components/DatePicker/DatePicker'
import { AddRow } from '@/components/DataSheet/AddRow'
import { ContextMenu } from '@/components/DataSheet/ContextMenu'
import { SaveIcon } from '@/components/Icons/SaveIcon'
import { RefreshIcon } from '@/components/Icons/RefreshIcon'
import { HyperlinkIcon } from '@/components/Icons/HyperlinkIcon'

const defaultValues: IOrderCreate = {
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

export const CreateOrderForm = () => {
	const { palette } = useTheme()

	const [link, setLink] = useState('')
	const [data, setData] = useState<IPositionCreate[]>(defaultValues.positions)
	const {
		control,
		handleSubmit,
		reset,
		formState: { isDirty },
	} = useForm<IOrderCreate>({ defaultValues })

	const [create, { isLoading }] = useCreateOrderMutation()

	const saveHandler = handleSubmit(async form => {
		if (!isDirty) {
			toast.error('Данные заказа не заполнены')
			return
		}
		if (form.consumer == '' && form.customer == '') {
			toast.error('Заполните хотя бы одного контрагента')
			return
		}

		if (data.some(item => item.quantity === null || item.name === '')) {
			toast.error('Заполните хотя бы одну позицию')
			return
		}

		data.forEach((element, idx) => {
			element.rowNumber = idx + 1
		})
		form.positions = data
		console.log(form)

		try {
			const payload = await create(form).unwrap()
			const url = new URL(`/orders/${payload.id}`, window.location.origin)
			setLink(url.toString())
			toast.success('Заказ успешно создан')
		} catch (error) {
			const fetchError = error as IFetchError
			toast.error(fetchError.data.message, { autoClose: false })
		}
	})

	const copyHandler = () => {
		navigator.clipboard.writeText(link)
		toast.success('Ссылка скопирована')
	}

	const clearHandler = () => {
		reset(defaultValues)
		setData(defaultValues.positions)
	}

	return (
		<Stack>
			<Typography align='center' variant='h5' mt={1} mb={3}>
				Создание заказа
			</Typography>

			<Stack direction={'row'} spacing={2} px={2} mb={2}>
				<Stack spacing={2} width={'50%'}>
					<Typography fontSize={'1.1rem'}>Контрагенты</Typography>

					<Controller
						control={control}
						name={`customer`}
						render={({ field }) => <TextField {...field} label='Заказчик' fullWidth />}
					/>
					<Controller
						control={control}
						name={`consumer`}
						render={({ field }) => <TextField {...field} label='Конечник' fullWidth />}
					/>
				</Stack>

				<Stack spacing={2} width={'50%'}>
					<Typography fontSize={'1.1rem'}>Детали заказа</Typography>

					<Controller
						control={control}
						name={`date`}
						rules={{ required: true }}
						render={({ field, fieldState: { error } }) => (
							<DatePicker
								{...field}
								value={field.value ? dayjs(field.value) : null}
								onChange={value => field.onChange(value?.startOf('d').toISOString())}
								label={'Дата'}
								showDaysOutsideCurrentMonth
								fixedWeekNumber={6}
								slots={{
									textField: DateTextField,
								}}
								slotProps={{
									textField: {
										error: Boolean(error),
									},
								}}
								minDate={dayjs('01-01-2010')}
								sx={{ width: '100%' }}
							/>
						)}
					/>

					<Controller
						control={control}
						name={`manager`}
						render={({ field }) => <TextField {...field} label='Менеджер / Помощник' fullWidth />}
					/>
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
				render={({ field }) => <TextField {...field} label='Примечание' multiline minRows={2} sx={{ mx: 2 }} />}
			/>

			<Divider sx={{ mt: 3, mb: 2 }} />

			<Typography align='center' mb={1}>
				Позиции
			</Typography>
			<Stack position={'relative'}>
				<DataSheetGrid
					value={data}
					onChange={setData}
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

					<Tooltip title='Скопировать ссылку'>
						<Button
							onClick={copyHandler}
							disabled={!link}
							sx={{
								minWidth: 48,
								background: '#fff',
								border: '1px solid #dcdcdc',
								borderRadius: '6px',
								padding: '4px 10px',
								':disabled': { svg: { fill: palette.action.disabled } },
								':hover': { svg: { fill: palette.secondary.main } },
							}}
						>
							<HyperlinkIcon fontSize={18} />
						</Button>
					</Tooltip>

					<Tooltip title='Очистить форму'>
						<Button
							onClick={clearHandler}
							disabled={!isDirty || isLoading}
							sx={{
								minWidth: 48,
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
					</Tooltip>
				</Stack>
			</Stack>
		</Stack>
	)
}
