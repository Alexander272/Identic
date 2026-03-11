import { FC, MouseEvent, PropsWithChildren } from 'react'

import { TableBodyContainer } from './style'

type Props = {
	onContext?: (event: MouseEvent<HTMLDivElement>) => void
}

export const TableBody: FC<PropsWithChildren<Props>> = ({ children, onContext }) => {
	return <TableBodyContainer onContextMenu={onContext}>{children}</TableBodyContainer>
}
