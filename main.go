package main

import (
	"github.com/darylhjd/mangadesk/app/service"
)

// Initialise the program.
func main() {
	// Initialise the application.
	service.Start()
	defer service.Shutdown()
}
