# <p align="center">mangadesk - Terminal client for MangaDex 📖</p>
<p align="center">
  <img alt="Top Langauge" src="https://img.shields.io/github/languages/top/darylhjd/mangadesk?style=flat-square">
  <img alt="Go Report" src="https://goreportcard.com/badge/github.com/darylhjd/mangadesk?style=flat-square">
  <img alt="License" src="https://img.shields.io/github/license/darylhjd/mangadesk?style=flat-square">
  <img alt="Downloads" src="https://img.shields.io/github/downloads/darylhjd/mangadesk/total?style=flat-square">
</p>

<p align="center">Download manga directly from your terminal to read!</p>

<img src="assets/demo.gif" alt="demo.gif">

<p align="center">
  This client retrieves information straight from MangaDex v5's API.<br>
  As the API is still a WIP, some changes (probably breaking) might be expected.
</p>

## Features ✨

- Download chapters straight to your computer.
- Login to keep track of your followed manga.
- Keep track of already downloaded manga.
- Download multiple chapters together.
- Searching!
- Responsive UI (kind of)
- Written in Golang :)

Works for Windows/Linux/macOS.

## Installation 🔧

Check out the releases page for relevant files.

For bleeding edge 🗡 updates, you may compile from source:

```cmd
git clone https://github.com/darylhjd/mangadesk.git
cd mangadesk
go get -d ./...
go build
```

## Usage ✍

Simply choose the chapters you want to read to download.

### Keybindings ⌨

| Operation                 | Binding                          | Page       |
|---------------------------|----------------------------------|------------|
| Login/Logout              | <kbd>Ctrl</kbd> + <kbd>L</kbd>   | All        |
| Keybindings/Help          | <kbd>Ctrl</kbd> + <kbd>K</kbd>   | All        |
| Search                    | <kbd>Ctrl</kbd> + <kbd>S</kbd>   | All        |
| Next/Prev Page            | <kbd>Ctrl</kbd> + <kbd>F/B</kbd> | Some       |
| Escape                    | <kbd>Esc</kbd>                   | Some       |
| Select multiple chapters  | <kbd>Ctrl</kbd> + <kbd>E</kbd>   | Manga Page |
| Toggle select all         | <kbd>Ctrl</kbd> + <kbd>A</kbd>   | Manga Page |

## Settings

### Download folder

By default, all downloads are stored in a folder titled `downloads`, relative to where you run the application.

However, you can change this by changing the `downloadDir` field in the `usr/usr_config.json` file (this file only
appears after running the application at least once!).

### Languages

By default, English (`en`) is the main language.

You may change your desired language(s) through the `languages` field in the `usr_config.json` file.

Please use comma-separated [ISO language codes](https://www.andiamo.co.uk/resources/iso-language-codes/)!

### Force Port 443

Valid options for this are `true` or `false`. It is `false` by default.

Set to `true` if you are having trouble downloading or are using networks that block traffic to non-standard ports 
(such as school/office networks). [More info](https://api.mangadex.org/docs.html#operation/get-at-home-server-chapterId).

## Issues ☠

Check out the Issues page for current issues/feature requests.

## Contributing 🤝

Always welcome and appreciated :)

Please take some time to familiarise yourself with the [contributing guidelines](.github/CONTRIBUTING.md).

## Learning points 🧠

- Creating TUIs with tview/tcell.
- Working with the filesystem in Golang.
- Goroutines & Context.
- Go Project structure.
