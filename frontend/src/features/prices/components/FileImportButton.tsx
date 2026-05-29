import { useRef, useState, useEffect, useCallback, type FC, type DragEvent } from 'react'
import { Box, Button, Dialog, IconButton, Stack, Typography, useTheme } from '@mui/material'
import { toast } from 'react-toastify'

import { useImportPricesMutation } from '@/features/prices/priceApiSlice'
import { AddFileIcon } from '@/components/Icons/AddFileIcon'
import { UploadIcon } from '@/components/Icons/UploadIcon'
import { TimesIcon } from '@/components/Icons/TimesIcon'

export const FileImportButton: FC = () => {
	const { palette } = useTheme()
	const [importPositions, { isLoading }] = useImportPricesMutation()

	const [open, setOpen] = useState(false)
	const [selectedFile, setSelectedFile] = useState<File | null>(null)
	const [dragOver, setDragOver] = useState(false)
	const [globalDragOver, setGlobalDragOver] = useState(false)
	const fileInputRef = useRef<HTMLInputElement>(null)
	const dragCounterRef = useRef(0)

	useEffect(() => {
		const handleDragEnter = () => {
			dragCounterRef.current += 1
			setGlobalDragOver(true)
		}
		const handleDragOver = (e: DragEvent) => {
			e.preventDefault()
		}
		const handleDragLeave = () => {
			dragCounterRef.current -= 1
			if (dragCounterRef.current <= 0) {
				dragCounterRef.current = 0
				setGlobalDragOver(false)
			}
		}
		const handleDrop = (e: DragEvent) => {
			e.preventDefault()
			dragCounterRef.current = 0
			setGlobalDragOver(false)
			const file = e.dataTransfer.files?.[0]
			if (file) {
				setSelectedFile(file)
				setOpen(true)
			}
		}

		window.addEventListener('dragenter', handleDragEnter as unknown as EventListener)
		window.addEventListener('dragover', handleDragOver as unknown as EventListener)
		window.addEventListener('dragleave', handleDragLeave as unknown as EventListener)
		window.addEventListener('drop', handleDrop as unknown as EventListener)

		return () => {
			window.removeEventListener('dragenter', handleDragEnter as unknown as EventListener)
			window.removeEventListener('dragover', handleDragOver as unknown as EventListener)
			window.removeEventListener('dragleave', handleDragLeave as unknown as EventListener)
			window.removeEventListener('drop', handleDrop as unknown as EventListener)
		}
	}, [])

	const handleFileSelect = useCallback((file: File | null) => {
		setSelectedFile(file)
	}, [])

	const handleConfirm = useCallback(async () => {
		if (!selectedFile) return
		const formData = new FormData()
		formData.append('file', selectedFile)
		try {
			await importPositions(formData).unwrap()
			toast.success(`Импортировано ${selectedFile.name}`)
			setOpen(false)
			setSelectedFile(null)
		} catch (err) {
			const fetchError = err as { data?: { message?: string } }
			toast.error(fetchError.data?.message ?? 'Ошибка импорта', { autoClose: false })
		}
	}, [selectedFile, importPositions])

	const handleClose = useCallback(() => {
		if (!isLoading) {
			setOpen(false)
			setSelectedFile(null)
		}
	}, [isLoading])

	const handleDragOver = useCallback((e: DragEvent) => {
		e.preventDefault()
		e.stopPropagation()
		setDragOver(true)
	}, [])

	const handleDragLeaveLocal = useCallback((e: DragEvent) => {
		e.preventDefault()
		e.stopPropagation()
		setDragOver(false)
	}, [])

	const handleDropLocal = useCallback(
		(e: DragEvent) => {
			e.preventDefault()
			e.stopPropagation()
			setDragOver(false)
			const file = e.dataTransfer.files?.[0]
			if (file) handleFileSelect(file)
		},
		[handleFileSelect],
	)

	return (
		<>
			<Box
				sx={{
					display: 'flex',
					justifyContent: 'center',
					alignItems: 'center',
					gap: 1,
					bgcolor: 'rgba(255,255,255,0.85)',
					border: '1px solid rgba(0,0,0,0.08)',
					backdropFilter: 'blur(20px)',
					boxShadow: '0 4px 12px rgba(0,0,0,0.04)',
					borderRadius: { xs: 2, sm: 3 },
					paddingX: { xs: 1.5, sm: 3 },
					paddingY: 1.5,
					mb: 1,
				}}
			>
				<Button
					variant='outlined'
					color='inherit'
					onClick={() => setOpen(true)}
					sx={{
						minWidth: 48,
						textTransform: 'inherit',
						background: '#fff',
						border: '1px solid #c3c3c4',
						borderRadius: '6px',
						padding: '4px 10px',
						':hover': { svg: { fill: palette.primary.main }, color: palette.primary.main },
						'&:disabled': { svg: { fill: palette.action.disabled } },
					}}
				>
					<UploadIcon sx={{ fontSize: 14, mr: 1 }} />
					Загрузить из XLSX
				</Button>
			</Box>

			<Dialog
				open={open}
				onClose={handleClose}
				maxWidth='sm'
				fullWidth
				slotProps={{
					backdrop: {
						sx: {
							backdropFilter: 'blur(4px)',
							backgroundColor: 'rgba(0,0,0,0.4)',
						},
					},
				}}
			>
				<Box sx={{ p: 3, position: 'relative' }}>
					<IconButton
						onClick={handleClose}
						size='large'
						sx={{ position: 'absolute', top: 12, right: 12, color: 'text.secondary' }}
					>
						<TimesIcon sx={{ fontSize: 14 }} />
					</IconButton>

					<Typography variant='h6' sx={{ mb: 2, pr: 4, fontWeight: 600 }}>
						Импорт XLSX
					</Typography>

					<Box
						onDragOver={handleDragOver}
						onDragLeave={handleDragLeaveLocal}
						onDrop={handleDropLocal}
						onClick={() => fileInputRef.current?.click()}
						sx={{
							border: '2px dashed',
							borderColor: dragOver ? palette.primary.main : palette.divider,
							borderRadius: 2,
							py: 5,
							px: 3,
							display: 'flex',
							flexDirection: 'column',
							alignItems: 'center',
							justifyContent: 'center',
							gap: 1,
							cursor: 'pointer',
							transition: 'all 0.2s ease',
							bgcolor: dragOver ? 'action.hover' : 'transparent',
							'&:hover': {
								borderColor: palette.primary.main,
								bgcolor: 'action.hover',
							},
						}}
					>
						<input
							ref={fileInputRef}
							type='file'
							accept='.xlsx,.xls'
							hidden
							onChange={e => {
								const file = e.target.files?.[0]
								if (file) handleFileSelect(file)
								e.target.value = ''
							}}
						/>

						<AddFileIcon sx={{ fontSize: 40, color: palette.text.secondary }} />
						<Typography variant='body1' color='text.secondary'>
							Нажмите или перетащите файл XLSX
						</Typography>
					</Box>

					{selectedFile && (
						<Typography variant='body2' color='text.secondary' sx={{ mt: 1.5 }}>
							Выбран: <strong>{selectedFile.name}</strong>
						</Typography>
					)}

					<Stack direction='row' spacing={1.5} sx={{ mt: 2.5 }}>
						<Button
							variant='outlined'
							disabled={!selectedFile || isLoading}
							onClick={handleConfirm}
							sx={{ flex: 1 }}
						>
							{isLoading ? 'Обрабатывается...' : 'Загрузить'}
						</Button>
						<Button variant='outlined' color='inherit' onClick={handleClose} sx={{ flex: 1 }}>
							Отмена
						</Button>
					</Stack>
				</Box>
			</Dialog>

			{globalDragOver && (
				<Box
					onDragOver={(e: DragEvent) => {
						e.preventDefault()
						e.stopPropagation()
					}}
					onDragLeave={() => setGlobalDragOver(false)}
					onDrop={(e: DragEvent) => {
						e.preventDefault()
						e.stopPropagation()
						dragCounterRef.current = 0
						setGlobalDragOver(false)
						const file = e.dataTransfer.files?.[0]
						if (file) {
							setSelectedFile(file)
							setOpen(true)
						}
					}}
					sx={{
						position: 'fixed',
						inset: 0,
						zIndex: 9999,
						display: 'flex',
						alignItems: 'center',
						justifyContent: 'center',
						backdropFilter: 'blur(4px)',
						backgroundColor: 'rgba(0,0,0,0.5)',
					}}
				>
					<Box
						sx={{
							width: '90%',
							maxWidth: 500,
							minHeight: 280,
							border: '2px dashed',
							borderColor: palette.primary.main,
							borderRadius: 4,
							display: 'flex',
							flexDirection: 'column',
							alignItems: 'center',
							justifyContent: 'center',
							gap: 2,
							bgcolor: 'rgba(255,255,255,0.12)',
						}}
					>
						<AddFileIcon sx={{ fontSize: 48, color: palette.primary.main }} />
						<Typography variant='h6' color='white'>
							Отпустите файл для импорта
						</Typography>
					</Box>
				</Box>
			)}
		</>
	)
}
