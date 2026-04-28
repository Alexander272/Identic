import React, { forwardRef, useRef, useEffect, useId } from 'react'
import styled from '@emotion/styled'

export interface CheckboxProps extends Omit<React.InputHTMLAttributes<HTMLInputElement>, 'type'> {
	label?: React.ReactNode
	indeterminate?: boolean
}

const Container = styled.label<{ disabled?: boolean }>`
	display: inline-flex;
	align-items: center;
	gap: 8px;
	cursor: ${({ disabled }) => (disabled ? 'not-allowed' : 'pointer')};
	user-select: none;
	border-radius: 8px;
	padding: 4px;
	font-family:
		system-ui,
		-apple-system,
		BlinkMacSystemFont,
		'Segoe UI',
		Roboto,
		sans-serif;
	transition: background-color 0.15s;

	&:focus-within {
		outline: 2px solid #3b82f6;
		outline-offset: 2px;
	}

	&:active {
		transform: scale(0.98);
	}
`

const Input = styled.input`
	position: absolute;
	opacity: 0;
	width: 0;
	height: 0;
`

const Box = styled.span<{ checked: boolean; disabled?: boolean }>`
	width: 20px;
	height: 20px;
	border: 2px solid ${({ disabled, checked }) => (disabled ? '#e2e8f0' : checked ? '#3b82f6' : '#cbd5e1')};
	border-radius: 6px;
	display: flex;
	align-items: center;
	justify-content: center;
	background: ${({ checked }) => (checked ? '#3b82f6' : 'transparent')};
	transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
	position: relative;
	flex-shrink: 0;
	pointer-events: none;

	&:hover {
		border-color: ${({ disabled }) => (disabled ? '#e2e8f0' : '#94a3b8')};
	}
`

const Checkmark = styled.svg<{ visible: boolean }>`
	width: 12px;
	height: 12px;
	stroke: white;
	stroke-width: 2.5;
	fill: none;
	stroke-linecap: round;
	stroke-linejoin: round;
	opacity: ${({ visible }) => (visible ? 1 : 0)};
	transform: ${({ visible }) => (visible ? 'scale(1)' : 'scale(0.8)')};
	transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
`

const IndeterminateLine = styled.div<{ visible: boolean }>`
	position: absolute;
	width: 8px;
	height: 2px;
	background: white;
	border-radius: 1px;
	opacity: ${({ visible }) => (visible ? 1 : 0)};
	transition: opacity 0.2s;
`

const LabelText = styled.span<{ disabled?: boolean }>`
	font-size: 14px;
	line-height: 1.5;
	color: ${({ disabled }) => (disabled ? '#94a3b8' : '#0f172a')};
	transition: color 0.2s;
`

export const Checkbox = forwardRef<HTMLInputElement, CheckboxProps>(
	({ label, indeterminate, disabled, checked, onChange, id, ...rest }, ref) => {
		const inputRef = useRef<HTMLInputElement>(null)
		const mergedRef = ref || inputRef
		const fallbackId = useId()
		const checkboxId = id || fallbackId

		// Синхронизация indeterminate с нативным элементом
		useEffect(() => {
			const el = typeof mergedRef === 'function' ? null : mergedRef?.current
			if (el) {
				el.indeterminate = !!indeterminate
			}
		}, [indeterminate, mergedRef])

		return (
			<Container disabled={disabled} htmlFor={checkboxId}>
				<Input
					ref={mergedRef}
					type='checkbox'
					id={checkboxId}
					checked={!!checked}
					disabled={disabled}
					onChange={onChange}
					aria-checked={indeterminate ? 'mixed' : checked}
					{...rest}
				/>
				<Box checked={!!checked} disabled={disabled}>
					<Checkmark visible={!!checked && !indeterminate}>
						<path d='M2 6l3 3 5-5' />
					</Checkmark>
					<IndeterminateLine visible={!!indeterminate} />
				</Box>
				{label && <LabelText disabled={disabled}>{label}</LabelText>}
			</Container>
		)
	},
)

Checkbox.displayName = 'Checkbox'
