import { Box } from '@mui/material'

import { Permissions } from '@/features/access/components/Permissions'

export default function PermissionsPage() {
	return (
		<Box flexGrow={1}>
			<Permissions />
		</Box>
	)
}
