import styled from '@emotion/styled'
import { keyframes } from '@emotion/react'

export const Field = styled.div`
	position: relative;
	display: flex;
	align-items: center;
`

export const Svg = styled.svg`
	position: absolute;
	pointer-events: none;
	user-select: none;
`

export const Label = styled.label`
	user-select: none;
	cursor: pointer;
	padding: 6px 8px;
	border-radius: 6px;
	overflow: hidden;
	transition: all 0.2s ease;

	&:hover {
		background: rgba(0, 119, 255, 0.06);

		span:first-of-type {
			border-color: #05287f;
		}
	}

	span {
		vertical-align: middle;
		transform: translate3d(0, 0, 0);

		&:first-of-type {
			position: relative;
			width: 18px;
			height: 18px;
			border-radius: 4px;
			transform: scale(1);
			border: 1px solid #cccfdb;
			transition: all 0.2s ease;
			box-shadow: 0 1px 1px rgba(0, 16, 75, 0.05);
			display: inline-block;

			${Svg} {
				position: absolute;
				top: 3px;
				left: 2px;
				fill: none;
				stroke: #fff;
				stroke-width: 2;
				stroke-linecap: round;
				stroke-linejoin: round;
				stroke-dasharray: 16px;
				stroke-dashoffset: 16px;
				transition: all 0.3s ease;
				transition-delay: 0.1s;
				transform: translate3d(0, 0, 0);
			}
		}

		&:last-child {
			padding-left: 8px;
			line-height: 18px;
		}
	}
`

const wave = keyframes`
    50% {
        transform: scale(0.9);
    }
`

export const Input = styled.input`
	position: absolute;
	visibility: hidden;

	&:checked + ${Label} span:first-of-type {
		background: #05287f;
		border-color: #05287f;
		animation: ${wave} 0.4s ease;

		${Svg} {
			stroke-dashoffset: 0;
		}
	}

	&:disabled {
		& + ${Label} {
			opacity: 0.6;
		}
	}
`
