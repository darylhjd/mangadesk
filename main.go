package main

import (
	"log"

	"github.com/darylhjd/mangadesk/core"
)

// Start the program.
func main() {
	// Initialise the application.
	core.Initialise()

	// Run the app.
	if err := core.App.Run(); err != nil {
		log.Println(err)
	}

	// Shutdown the application.
	core.App.Shutdown()
}
