<div align="center">

# mangadesk 📖

  [![Top Language](https://img.shields.io/github/languages/top/darylhjd/mangadesk?style=flat-square)](https://github.com/darylhjd/mangadesk/search?l=go)
  [![License](https://img.shields.io/github/license/darylhjd/mangadesk?style=flat-square&color=blue)](https://github.com/darylhjd/mangadesk/blob/main/LICENSE)
  [![Go Report](https://goreportcard.com/badge/github.com/darylhjd/mangadesk?style=flat-square)](https://goreportcard.com/report/github.com/darylhjd/mangadesk)
  [![Downloads](https://img.shields.io/github/downloads/darylhjd/mangadesk/total?style=flat-square&color=success)](https://github.com/darylhjd/mangadesk/releases)

  **Note: Maintenance mode. Only severe bugs _may_ be fixed.**
  <h3>The ultimate MangaDex terminal client!</h3>
  <img src=".github/assets/demo.gif" alt="demo.gif">

</div>

## Features ✨

- Download chapters straight to your computer.
- Login to keep track of your followed manga.
- Download multiple chapters together.
- Searching!
- (Yes, you can use this to scrape manga).
- Written in Golang :)

<b>Works for Windows/Linux/macOS.</b>

## Installation 🔧

This application runs as a standalone executable, and does not need to be installed.

Check out the [Releases Page](https://github.com/darylhjd/mangadesk/releases) for new releases. To update, just
download the latest release.

For bleeding edge 🗡 updates, you may also compile from source:

```cmd
$ git clone https://github.com/darylhjd/mangadesk.git
$ cd mangadesk
$ go get -d ./...
$ go build
```

### Arch Linux

Mangadesk is available through the [AUR](https://aur.archlinux.org/packages/mangadesk/) and may be installed as such (thanks @AmaanHUB!).
It may be installed manually or with your preferred AUR helper:

```cmd
$ paru -S mangadesk
```

### Ubuntu

Mangadesk is available through the [MPR](https://mpr.makedeb.org/packages/mangadesk) and can be installed manually or with your preferred MPR helper:

```cmd
$ una install mangadesk
```

### Void Linux

Mangadesk is available through the official Void repositories and can be installed via XBPS:

```cmd
$ xbps-install -S mangadesk
```

## Uninstall ❌

To uninstall, simply delete the executable and its related folders and files.

Your downloads will not be removed by deleting the executable.

### Arch Linux

Uninstall with an AUR helper or with pacman:

```cmd
$ pacman -R mangadesk
```

### Ubuntu

Uninstall with a MPR helper or with APT:

```cmd
$ una remove mangadesk
```

### Void Linux

Uninstall with XBPS:

```cmd
$ xbps-remove -R mangadesk
```

## Usage ✍

To run the application, navigate to the directory where you stored the executable, and run the following command:

```cmd
$ ./mangadesk
```

Steps may differ for different OSes. For example, in Windows, use a backslash `\` instead.

### Keybindings ⌨

| Operation                                                                                 | Binding                          |
|-------------------------------------------------------------------------------------------|----------------------------------|
| Login/Logout                                                                              | <kbd>Ctrl</kbd> + <kbd>L</kbd>   |
| Keybindings/Help                                                                          | <kbd>Ctrl</kbd> + <kbd>K</kbd>   |
| Search                                                                                    | <kbd>Ctrl</kbd> + <kbd>S</kbd>   |
| Next/Prev Page                                                                            | <kbd>Ctrl</kbd> + <kbd>F/B</kbd> |
| Escape                                                                                    | <kbd>Esc</kbd>                   |
| Select a chapter                                                                          | <kbd>Ctrl</kbd> + <kbd>E</kbd>   |
| Toggle select all chapters                                                                | <kbd>Ctrl</kbd> + <kbd>A</kbd>   |
| Toggle chapter(s) read status<br/><br/>*Note: You can select multiple chapters to toggle! | <kbd>Ctrl</kbd> + <kbd>R</kbd>   |
| Toggle manga following                                                                    | <kbd>Ctrl</kbd> + <kbd>Q</kbd>   |

## Settings ⚙

Refer to [this document](app/core/CONFIG.md) for configurable settings.

## Issues ☠

Check out the Issues page for current issues/feature requests.

## Contributing 🤝

Always welcome and appreciated :)

Please take some time to familiarise yourself with the [contributing guidelines](.github/CONTRIBUTING.md).

## Learning Points 🧠

- Creating TUIs with tview/tcell.
- Working with the filesystem in Golang.
- Goroutines & Context.
- Go project structure.
