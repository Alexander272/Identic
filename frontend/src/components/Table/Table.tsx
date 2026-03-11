import { CSSProperties, FC, PropsWithChildren } from 'react'

import { TableContainer } from './style'

type Props = {
	sx?: CSSProperties
}

export const Table: FC<PropsWithChildren<Props>> = ({ children, sx }) => {
	return <TableContainer styles={sx}>{children}</TableContainer>
}
