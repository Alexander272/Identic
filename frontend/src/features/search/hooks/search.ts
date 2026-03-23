import { useEffect } from 'react'

export const useSearchHotkeys = (onSearchTrigger: (mode: 'exact' | 'fuzzy') => void) => {
	useEffect(() => {
		const handleKeyDown = (e: KeyboardEvent) => {
			// Проверяем, нажат ли Ctrl (или Command на Mac)
			const isCtrl = e.ctrlKey || e.metaKey

			// 1. Неточный поиск: Ctrl + Shift + F
			if (isCtrl && e.shiftKey && e.code === 'KeyF') {
				e.preventDefault()
				onSearchTrigger('fuzzy')
				return
			}

			// 2. Обычный поиск: Ctrl + S
			if (isCtrl && e.code === 'KeyS') {
				e.preventDefault() // КРИТИЧНО: отменяет сохранение страницы браузером
				onSearchTrigger('exact')
				return
			}
		}

		window.addEventListener('keydown', handleKeyDown, true)
		return () => window.removeEventListener('keydown', handleKeyDown, true)
	}, [onSearchTrigger])
}
