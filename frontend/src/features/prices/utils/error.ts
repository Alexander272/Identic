export function isApiError(err: unknown): err is { data: { message: string } } {
	return typeof err === 'object' && err !== null && 'data' in err
}
