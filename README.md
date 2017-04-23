# Wallpaper switcher
This application can switch wallpaper using `feh` by timeout or shortcut.
Application is controlled by named pipe.

Usage example:

```bash
./wallpaper-switcher  /path/to/images  /path/to/named/pipe   switch_timeout(in seconds)
```

Write something to pipe to initiate wallpaper switching:

```bash
echo next > /path/to/named/pipe
```

## License
This software is licensed under the [MIT license](LICENSE).
