export interface IScrollbarParameters {
	scrollbarWidth: string
	scrollbarHeight: string
	scrollbarBgColor: string
	thumbColor: string
	thumbColorHover: string
	thumbColorActive: string
	borderRadius: string
}

export const generateScrollbarStyles = (parameters: IScrollbarParameters) => {
	const {
		borderRadius,
		scrollbarHeight,
		scrollbarWidth,
		scrollbarBgColor,
		thumbColor,
		thumbColorActive,
		thumbColorHover,
	} = parameters

	return `
  * {
    scrollbar-color: ${thumbColor} ${scrollbarBgColor};
    scrollbar-width: thin;
  }
	*::-webkit-scrollbar {
		background-color: ${scrollbarBgColor};
		width: ${scrollbarWidth};
		height: ${scrollbarHeight}; 
	}
	*::-webkit-scrollbar-corner {
		background-color: ${scrollbarBgColor};
	}
	*::-webkit-scrollbar-thumb {
		background-color: ${thumbColor};
		border-radius: ${borderRadius};
	}
	*::-webkit-scrollbar-thumb:hover {
		background-color: ${thumbColorHover};
	}
	*::-webkit-scrollbar-thumb:active {
		background-color: ${thumbColorActive};
	}
	`
}
