import type { CSSProperties, FC, MouseEvent, PropsWithChildren } from 'react'

import { TableRowContainer } from './style'

type Props = {
	onClick?: () => void
	onContext?: (event: MouseEvent<HTMLDivElement>) => void
	width?: number
	height?: number
	hover?: boolean
	sx?: CSSProperties
}

export const TableRow: FC<PropsWithChildren<Props>> = ({ children, onClick, onContext, width, height, hover, sx }) => {
	return (
		<TableRowContainer
			onClick={onClick}
			onContextMenu={onContext}
			width={width}
			height={height}
			hover={hover}
			styles={sx}
		>
			{children}
		</TableRowContainer>
	)
}
