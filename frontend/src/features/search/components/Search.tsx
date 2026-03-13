import { useState } from 'react'
import { Button, FormControlLabel, Stack, Tooltip, Typography, useTheme, Checkbox, Divider } from '@mui/material'
import { DataSheetGrid, floatColumn, keyColumn, textColumn, type Column } from 'react-datasheet-grid'
import { toast } from 'react-toastify'

import type { ISearchItem, ISearchResults } from '../types/search'
import { useLazyFindOrdersQuery } from '../searchApiSlice'
import { extractQuantity } from '@/utils/extract'
import { ContextMenu } from '@/components/DataSheet/ContextMenu'
import { AddRow } from '@/components/DataSheet/AddRow'
import { BoxFallback } from '@/components/Fallback/BoxFallback'
import { SearchIcon } from '@/components/Icons/SearchIcon'
import { RefreshIcon } from '@/components/Icons/RefreshIcon'
import { Results } from './Results'

const defaultData = [{ name: '', quantity: null }]

const columns: Column<ISearchItem>[] = [
	{
		...keyColumn<ISearchItem, 'name'>('name', textColumn),
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
	{ ...keyColumn<ISearchItem, 'quantity'>('quantity', floatColumn), title: 'Количество', width: 0.5 },
]

export const Search = () => {
	const { palette } = useTheme()
	const [useFuzzy, setUseFuzzy] = useState(false)
	const [data, setData] = useState<ISearchItem[]>(defaultData)
	const [results, setResults] = useState<ISearchResults[]>([])

	const [findOrders, { isFetching, isSuccess }] = useLazyFindOrdersQuery()

	const clearHandler = () => {
		setData(defaultData)
	}

	const findHandler = async () => {
		const items = data.filter(item => Boolean(item.name) && item.quantity !== null)

		if (items.length === 0) {
			toast.error('Заполните хотя бы одну строку')
			return
		}

		const payload = await findOrders({ items: items, isFuzzy: useFuzzy }).unwrap()
		setResults(payload.data)
	}

	return (
		<Stack direction={'row'} position={'relative'} height={'100%'}>
			{isFetching ? <BoxFallback /> : null}

			<Stack>
				<Typography align='center' variant='h5'>
					Множественный поиск
				</Typography>
				<Typography align='center' mb={2}>
					Заполните таблицу и нажмите «Найти»
				</Typography>

				<FormControlLabel
					control={<Checkbox checked={useFuzzy} onChange={e => setUseFuzzy(e.target.checked)} />}
					label='Искать похожие наименования'
					sx={{
						mb: 1,
						ml: 0.5,
						pl: 1,
						transition: 'background-color 0.2s ease-in-out',
						borderRadius: 2,
						':hover': { backgroundColor: '#eff8ff' },
					}}
				/>

				<Stack width={800} position={'relative'} alignSelf={'center'}>
					<DataSheetGrid
						value={data}
						onChange={setData}
						columns={columns}
						contextMenuComponent={props => <ContextMenu {...props} />}
						addRowsComponent={props => <AddRow {...props} />}
						autoAddRow
					/>
					<Stack direction={'row'} spacing={1} sx={{ position: 'absolute', right: 8, bottom: 6 }}>
						<Tooltip title='Найти'>
							<Button
								onClick={findHandler}
								color='inherit'
								disabled={isFetching}
								sx={{
									minWidth: 48,
									textTransform: 'inherit',
									background: '#fff',
									border: '1px solid #dcdcdc',
									borderRadius: '6px',
									padding: '4px 10px',
									':hover': { svg: { fill: palette.primary.main }, color: palette.primary.main },
								}}
							>
								<SearchIcon fontSize={18} mr={1} />
								Найти
							</Button>
						</Tooltip>

						<Tooltip title='Очистить'>
							<Button
								onClick={clearHandler}
								disabled={isFetching}
								sx={{
									minWidth: 48,
									background: '#fff',
									border: '1px solid #dcdcdc',
									borderRadius: '6px',
									padding: '4px 10px',
									':hover': { svg: { fill: palette.secondary.main } },
								}}
							>
								<RefreshIcon fontSize={18} />
							</Button>
						</Tooltip>
					</Stack>
				</Stack>
			</Stack>

			{isSuccess && (
				<>
					<Divider orientation='vertical' flexItem sx={{ mx: 2 }} />
					{/* {!results || results?.length == 0 ? (
						
					) : null} */}
					<Results data={results || []} />
				</>
			)}
		</Stack>
	)
}
