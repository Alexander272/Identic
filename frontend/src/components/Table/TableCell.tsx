import { FC, PropsWithChildren } from 'react'

import { TableCellContainer } from './style'

type Props = {
	width?: number
	colSpan?: number
	rowSpan?: number
	isActive?: boolean
	onClick?: () => void
}

export const TableCell: FC<PropsWithChildren<Props>> = ({ children, width, isActive, onClick }) => {
	return (
		<TableCellContainer width={width} isActive={isActive} onClick={onClick}>
			{children}
		</TableCellContainer>
	)
}
