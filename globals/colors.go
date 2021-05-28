package globals

import "github.com/gdamore/tcell/v2"

const ( // Login page colours
	LoginPageTitleColor = tcell.ColorOrange

	LoginFormBorderColor = tcell.ColorGrey
	LoginFormLabelColor  = tcell.ColorWhite
)

const ( // Main page colors
	MainPageGridTitleColor   = tcell.ColorOrange
	MainPageGridBorderColor  = tcell.ColorGrey
	MainPageTableTitleColor  = tcell.ColorLightSkyBlue
	MainPageTableBorderColor = tcell.ColorGrey

	LoggedMainPageTitleColor     = tcell.ColorLightGoldenrodYellow
	LoggedMainPagePubStatusColor = tcell.ColorSandyBrown

	GuestMainPageTitleColor = tcell.ColorOrange
	GuestMainPageDescColor  = tcell.ColorLightGrey
	GuestMainPageTagColor   = tcell.ColorLightSteelBlue
)

const ( // Modal colors
	ModalColor = tcell.ColorDarkSlateGrey
)
