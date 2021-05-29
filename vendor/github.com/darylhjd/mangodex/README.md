# mangodex

[![Go Reference](https://pkg.go.dev/badge/github.com/darylhjd/mangodex.svg)](https://pkg.go.dev/github.com/darylhjd/mangodex)

Golang API wrapper for MangaDex v5's MVP API.

Full documentation is found [here](https://api.mangadex.org/docs.html).

This API is still in Open Beta, so testing may not be complete. However, basic authentication has been tested.

## Installation

To install, do `go get -u github.com/darylhjd/mangodex`.

## Usage

```golang
package main

import (
	"fmt"
	"github.com/darylhjd/mangodex"
)

func main() {
	// Create new client.
	// Without logging in, you may not be able to access 
	// all API functionality.
	c := mangodex.NewDexClient()

	// Login using your username and password.
	err := c.Login("user", "password")
	if err != nil {
		fmt.Println("Could not login!")
	}
}
```

## Contributing

Rapid changes expected. Any contributions are welcome.
