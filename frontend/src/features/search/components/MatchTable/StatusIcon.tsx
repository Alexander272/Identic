import type { FC } from 'react'

import { VerifiedIcon } from '@/components/Icons/VerifiedIcon'
import { WarnIcon } from '@/components/Icons/WarnIcon'
import { CloseRoundIcon } from '@/components/Icons/CloseRoundIcon'
import { colors } from './colors'

export const StatusIcon: FC<{ status: 'found' | 'partial' | 'not_found' }> = ({ status }) => {
	if (status === 'found') {
		return (
			// <SvgIcon sx={{ color: colors.successBorder, fontSize: '1.2rem' }}>
			// 	<path d='M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41L9 16.17z' />
			// </SvgIcon>
			<VerifiedIcon fill={colors.successBorder} fontSize={'1.3rem'} overflow={'visible'} />
		) // ✅
	} else if (status === 'partial') {
		return (
			// <SvgIcon sx={{ color: colors.warningBorder, fontSize: '1.2rem' }}>
			// 	<path d='M1 21h22L12 2 1 21zm12-3h-2v-2h2v2zm0-4h-2v-4h2v4z' />
			// </SvgIcon>
			<WarnIcon fill={colors.warningBorder} fontSize={'1.3rem'} overflow={'visible'} />
		) // ⚠️
	} else {
		return (
			// <SvgIcon sx={{ color: colors.errorBorder, fontSize: '1.2rem' }}>
			// 	<path d='M19 6.41L17.59 5 12 10.59 6.41 5 5 6.41 10.59 12 5 17.59 6.41 19 12 13.41 17.59 19 19 17.59 13.41 12 19 6.41z' />
			// </SvgIcon>
			<CloseRoundIcon fill={colors.errorBorder} fontSize={'1.3rem'} overflow={'visible'} />
		) // ❌
	}
}
