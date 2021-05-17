package globals

import (
	"github.com/darylhjd/mangodex"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var (
	App  = tview.NewApplication()
	Dex  = mangodex.NewDexClient()
	Conf UserConfig
)

const (
	LoginPageID  = "login_page" // Main Pages and their IDs
	MainPageID   = "main_page"
	MangaPageID  = "manga_page"
	HelpPageID   = "help_page"
	SearchPageID = "search_page"

	LoginLogoutFailureModalID   = "login_failure_modal" // Modals and their IDs
	LoginLogoutCfmModalID       = "logout_modal"
	StoreCredentialErrorModalID = "store_cred_error_modal"
	DownloadChaptersModalID     = "download_chapters_modal"
	DownloadFinishedModalID     = "download_error_modal"
	GenericAPIErrorModalID      = "api_error_modal"
)

const (
	LoggedMainPageTitleColor  = tcell.ColorLightGoldenrodYellow
	LoggedMainPageStatusColor = tcell.ColorSaddleBrown

	GuestMainPageTitleColor = tcell.ColorOrange
	GuestMainPageDescColor  = tcell.ColorLightGrey
	GuestMainPageTagColor   = tcell.ColorLightSteelBlue
)
