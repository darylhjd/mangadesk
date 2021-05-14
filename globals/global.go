package globals

import (
	"github.com/darylhjd/mangodex"
	"github.com/rivo/tview"
)

var (
	App = tview.NewApplication()
	Dex = mangodex.NewDexClient()
)

const (
	LoginPageID  = "login_page"
	MainPageID   = "main_page"
	MangaPageID  = "manga_page"
	HelpPageID   = "help_page"
	SearchPageID = "search_page"

	LoginLogoutFailureModalID   = "login_failure_modal"
	LoginLogoutCfmModalID       = "logout_modal"
	StoreCredentialErrorModalID = "store_cred_error_modal"
	DownloadChaptersModalID     = "download_chapters_modal"
	DownloadFinishedModalID     = "download_error_modal"
	GenericAPIErrorModalID      = "api_error_modal"
)

const (
	UsrDir       = "usr"
	CredFileName = "cred"
	DownloadDir  = "downloads"
)
