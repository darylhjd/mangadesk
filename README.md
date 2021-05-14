# mangadesk - Terminal client for MangaDex 📖

<p align="center">Download manga directly from your terminal to read!</p>

<img src="assets/demo.gif" alt="">

<p align="center">This client retrieves information straight from MangaDex's API. <br>As the API is still a WIP, some changes (probably breaking) might be expected.</p>

## Features ✨

- Read chapters straight on your computer after downloading them.
- Login to keep track of your followed manga.
- Download multiple chapters together.
- Searching!
- Responsive UI (kind of)
- Written in Golang :)

Works for Windows/Linux. 

Not yet tested on MacOS. If you have compiled this on Mac and it works for you, do tell me so I can update this.

## Usage

Simply choose the chapters you want to read to download.

All downloads are stored in a folder titled`downloads`, which you will find in the same directory as where you ran this application.

### Keybindings

- Ctrl + L : Login/Logout
- Ctrl + H : Help
- Ctrl + S : Search
- Ctrl + E : Select multiple chapters
- Esc      : Going back

## Installation

Check out the releases page for relevant files.

If you want, you may compile from source,

```
git clone https://github.com/darylhjd/mangadesk.git
cd mangadesk
go get -d ./...
go build
```

## Issues ☠

- Search table does not clear itself for a new search. For now, you may do Ctrl+S to re-open the search page.

## Planning... maybe?

- Config files for personal settings (download folder, language selection etc...)
- More download information (notify user when download finished, show downloaded chapters etc...)

## Contributing

Always welcome and appreciated :)
