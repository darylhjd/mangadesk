# <p align="center">mangadesk - Terminal client for MangaDex 📖</p>

<p align="center">
  <img alt="Top Language" src="https://img.shields.io/github/languages/top/darylhjd/mangadesk?style=flat-square">
  <img alt="License" src="https://img.shields.io/github/license/darylhjd/mangadesk?style=flat-square&color=blue">
  <img alt="Go Report" src="https://goreportcard.com/badge/github.com/darylhjd/mangadesk?style=flat-square">  
  <img alt="Downloads" src="https://img.shields.io/github/downloads/darylhjd/mangadesk/total?style=flat-square&color=success">
</p>

<p align="center">
  Download manga directly from your terminal to read!<br><br>
  <img src="assets/demo.gif" alt="demo.gif">
</p>

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

This application runs as a standalone executable, and does not install itself.

Check out the [releases page](https://github.com/darylhjd/mangadesk/releases) for relevant files. To update, just
download the latest release.

For bleeding edge 🗡 updates, you may compile from source:

```cmd
git clone https://github.com/darylhjd/mangadesk.git
cd mangadesk
go get -d ./...
go build
```

**NOTE**: The application will create a `usr` folder in the same directory as where you run the application to store
your credentials/configurations.

## Uninstall ❌

To uninstall, simply delete the executable and its related folders and files (namely, the `usr` folder).

Your downloads will not be removed by deleting the executable.

## Usage ✍

To run the application, navigate to the directory where you stored the executable, and run the following command:

```cmd
$ ./mangadesk 
```

Steps may differ for different OSes. For example, in Windows, use a backslash `\` instead.

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

## Settings ⚙

You may change the appropriate settings in the `usr_config.json` file, which is stored in the `usr` folder.

### Download Folder

- `downloadDir`

By default, all downloads will be stored in a folder titled `downloads`.

You can change this by changing the `downloadDir` field.

### Languages

- `langauges`

By default, only English (`en`) translated chapters will be shown.

You may change your desired language(s) through the `languages` field. Please use
comma-separated [ISO language codes](https://www.andiamo.co.uk/resources/iso-language-codes/).

### Download Quality

- `downloadQuality`

Valid options are `data` (high quality) or `data-saver` (lower quality).

Any other empty/invalid option will default to `data`.

### Force Port 443

- `forcePort443`

Valid options are `true` or `false`. It is `false` by default.

Set to `true` if you are having trouble downloading or are using networks that block traffic to non-standard ports
(such as school/office networks).
[More info](https://api.mangadex.org/docs.html#operation/get-at-home-server-chapterId).

### As Zip

- `asZip`

Valid options are `true` or `false`. It is `false` by default.

Set to `true` if you want your chapter downloads to be compressed into a zip folder.

### Zip Type

- `zipType`

Valid options are `zip` or `cbz`. This is ignored if `asZip` is set to `false`.

Any other empty/invalid option will default to `zip`.

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
