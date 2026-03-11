import { FC, PropsWithChildren } from 'react'

import { TableGroupContainer } from './style'

export const TableGroup: FC<PropsWithChildren> = ({ children }) => {
	return <TableGroupContainer>{children}</TableGroupContainer>
}
