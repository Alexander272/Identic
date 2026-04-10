import { getPermissions } from '@/features/user/userSlice'
import { useAppSelector } from '@/hooks/redux'

// export const useCheckPermission = (rule: string) => {
// 	const permissions = useAppSelector(getPermissions)
// 	if (!permissions.length) return false

// 	for (let i = 0; i < permissions.length; i++) {
// 		if (permissions[i] === rule) return true
// 	}
// 	return false
// }

export const useCheckPermission = (requiredRule: string): boolean => {
	const permissions = useAppSelector(getPermissions) // Например, ["order:*", "user:read"]

	if (!permissions || !permissions.length) return false

	return permissions.some(p => {
		// 1. Полное совпадение (например, "order:read" === "order:read")
		if (p === requiredRule) return true

		// 2. Если в правах есть полный доступ "admin" или "*:*"
		if (p === '*' || p === '*:*') return true

		// 3. Логика со звездочкой (например, "order:*" покроет "order:read")
		const [pObj, pAct] = p.split(':')
		const [reqObj, reqAct] = requiredRule.split(':')

		const objMatch = pObj === '*' || pObj === reqObj
		const actMatch = pAct === '*' || pAct === reqAct

		return objMatch && actMatch
	})
}
