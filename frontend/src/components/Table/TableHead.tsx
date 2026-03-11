import { FC, PropsWithChildren } from 'react'

import { TableHeadContainer } from './style'

export const TableHead: FC<PropsWithChildren> = ({ children }) => {
	return <TableHeadContainer>{children}</TableHeadContainer>
}
