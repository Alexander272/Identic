import styled from '@emotion/styled'

export const Wrapper = styled.label<{ disabled?: boolean }>`
	display: inline-flex;
	align-items: center;
	gap: 10px;
	cursor: ${({ disabled }) => (disabled ? 'not-allowed' : 'pointer')};
	opacity: ${({ disabled }) => (disabled ? 0.5 : 1)};
	user-select: none;
	font-family: system-ui, sans-serif;
`

export const HiddenInput = styled.input`
	display: none;

	&:focus-visible + div {
		box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.35);
	}
`

export const Box = styled.div<{ checked: boolean }>`
	width: 22px;
	height: 22px;
	border-radius: 6px;
	border: 2px solid #d1d5db;
	display: flex;
	align-items: center;
	justify-content: center;
	transition: all 0.25s ease;
	background: ${({ checked }) => (checked ? 'linear-gradient(135deg, #6366f1, #8b5cf6)' : '#fff')};
	border-color: ${({ checked }) => (checked ? 'transparent' : '#d1d5db')};
	box-shadow: ${({ checked }) => (checked ? '0 4px 12px rgba(99, 102, 241, 0.35)' : 'none')};

	&:hover {
		border-color: #6366f1;
	}
`

export const CheckIcon = styled.svg`
	width: 14px;
	height: 14px;
	stroke: white;
	fill: none;
	stroke-width: 3;
	stroke-dasharray: 24;
	stroke-dashoffset: 24;
	transition: stroke-dashoffset 0.3s ease;

	${Box}[checked="true"] & {
		stroke-dashoffset: 0;
	}
`

export const Label = styled.span`
	font-size: 14px;
	color: #111827;
`
