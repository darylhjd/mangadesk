## Configurations ⚙

You may change the appropriate settings in the `config.json` file, which is stored in the `mangadesk` folder within your
OS' default configuration folder (this folder [depends](https://pkg.go.dev/os#UserConfigDir) on your OS).

If no location can be found, it will instead be found in the default home directory (also
[depends](https://pkg.go.dev/os#UserHomeDir) on your OS!)

### Download Folder

- `downloadDir`

By default, all downloads will be stored in a folder titled `downloads` relative to where the executable is run.

You may use environment variables to specify the new path. Take special note when specifying relative or absolute paths!

### Languages

- `langauges`

By default, only English (`en`) translated chapters will be shown. Please use
comma-separated [ISO language codes](https://www.andiamo.co.uk/resources/iso-language-codes/).

### Download Quality

- `downloadQuality`

Valid options are `data` (high quality) or `data-saver` (lower quality). Any other empty/invalid option will default
to `data`.

### Explicit Content

- `explicitContent`

Valid options are `true` or `false`. It is `false` by default.

Set to `true` if you want to see explicit content in your feeds.

Note that search results follow search-specific settings. If you do not check the `Explicit Content` checkbox, you will
not see explicit results regardless of your configuration settings.

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

Valid options are `zip` or `cbz`. This is ignored if `asZip` is set to `false`. Any other empty/invalid option will
default to `zip`.