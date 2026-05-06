import { useMemo } from 'react'

import type { IOrderMatchResult, ISearchItem } from '../../types/search'
import type { IPosition } from '@/features/orders/types/positions'
import { normalize } from '../../utils/normalize'

export type MatchStatus = 'found' | 'partial' | 'not_found'

export type MatchRowData = {
	index: number
	item: ISearchItem
	status: MatchStatus
	foundItem?: IPosition | null
	mismatch?: {
		name: boolean
		qty: boolean
	}
	matchedFrom?: 'name' | 'notes' | null
}

export const useMatchTable = (request: ISearchItem[], result: IOrderMatchResult, foundPositions: IPosition[]) => {
	const foundMap = useMemo(() => {
		return new Map<string, IPosition>(foundPositions.map(p => [p.id, p]))
	}, [foundPositions])

	const matchMap = useMemo(() => {
		return new Map<string, IOrderMatchResult['positions'][number]>(result.positions.map(p => [p.reqId, p]))
	}, [result.positions])

	const rows: MatchRowData[] = useMemo(() => {
		return request.map((item, index) => {
			const match = matchMap.get(index.toString())

			if (!match) {
				return {
					index,
					item,
					status: 'not_found',
					foundItem: null,
				}
			}

			const foundItem = foundMap.get(match.id)

			if (!foundItem) {
				return {
					index,
					item,
					status: 'not_found',
					foundItem: null,
				}
			}

			const reqQty = item.quantity ?? 0
			const foundQty = foundItem.quantity ?? 0

			const foundName = normalize(foundItem.name)
			const foundNotes = normalize(foundItem.notes)
			const reqName = item.name ? normalize(item.name) : ''

			const isNameMatch = foundName === reqName
			const isNotesMatch = foundNotes === reqName && reqName != ''

			const matchedFrom: 'name' | 'notes' | null = isNameMatch ? 'name' : isNotesMatch ? 'notes' : null

			const mismatch = {
				name: !matchedFrom,
				qty: foundQty !== reqQty,
			}

			if (mismatch.name || mismatch.qty) {
				return {
					index,
					item,
					status: 'partial',
					foundItem,
					mismatch,
					matchedFrom,
				}
			}

			return {
				index,
				item,
				status: 'found',
				foundItem,
				matchedFrom,
			}
		})
	}, [request, matchMap, foundMap])

	return { rows }
}
