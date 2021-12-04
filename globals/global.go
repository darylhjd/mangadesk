package globals

import (
	"github.com/darylhjd/mangodex"
	"github.com/rivo/tview"
)

// Global values for the application.
var (
	App       = tview.NewApplication()  // The tview application.
	DexClient = mangodex.NewDexClient() // The MangaDex client for interfacing with the API.
	Conf      UserConfig
)

const (
	OffsetRange = 100
)
