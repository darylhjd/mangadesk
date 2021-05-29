package globals

import (
	"github.com/darylhjd/mangodex"
	"github.com/rivo/tview"
)

var (
	App  = tview.NewApplication()
	Dex  = mangodex.NewDexClient()
	Conf UserConfig
)

const (
	OffsetRange = 100
)
