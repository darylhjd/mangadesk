# mangadesk - Terminal client for MangaDex 📖

<p align="center">Read and download manga directly from your terminal!</p>

<img src="assets/demo.gif" alt="">

<p align="center">This client retrieves information straight from MangaDex's API. <br>As the API is still a WIP, some changes (probably breaking) might be expected.</p>

## Features

- Read chapters straight on your computer.
- Login to keep track of your followed manga.
- Download multiple chapters together.
- Searching!
- Responsive UI (kind of)

Works for Windows/Linux. Not yet tested on MacOS.

## Installation & Usage

Check out the releases page for relevant files.

To compile from source,

```
git clone https://github.com/darylhjd/mangadesk.git
cd mangadesk
go get -d ./...
go build
```

### Keybindings

- Ctrl + L : Login/Logout
- Ctrl + H : Help
- Ctrl + S : Search
- Ctrl + E : Select multiple chapters
- Esc      : Going back

## Contributing

Always welcome and appreciated :)
