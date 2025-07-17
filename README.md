# wsl-notify-send

A cross-platform notification tool for Windows and WSL2 that provides a clean interface to send desktop notifications, alerts, and beeps using the [beeep](https://github.com/gen2brain/beeep) library.

## Features

- **Cross-platform**: Works on Windows natively and in WSL2
- **Multiple notification types**: 
  - Silent notifications
  - Alert notifications with sound
  - Beep-only mode
- **Icon support**: PNG, JPG, ICO, BMP files and stock icons
- **Customizable**: App name, sound frequency, and duration
- **Clean CLI**: Built with Cobra framework for intuitive usage

## Installation

### Prerequisites

- Go 1.19 or later

### Build from source

```bash
git clone https://github.com/alexferrari88/wsl-notify-send
cd wsl-notify-send
go build -o wsl-notify-send.exe .
```

### Install binary

Copy the `wsl-notify-send.exe` to a directory in your PATH for global access.

## Usage

### Basic Examples

```bash
# Send a simple notification
wsl-notify-send "Hello" "World"

# Send an alert with sound
wsl-notify-send --alert "Warning" "Something happened"

# Just beep
wsl-notify-send --beep

# Send notification with icon
wsl-notify-send --icon icon.png "Info" "With custom icon"

# Set custom app name
wsl-notify-send --app-name "MyApp" "Custom" "From MyApp"
```

### Advanced Options

```bash
# Custom beep frequency and duration
wsl-notify-send --beep --freq 800 --duration 1000

# Quiet mode (suppress error output)
wsl-notify-send --quiet "Title" "Message"

# Show version
wsl-notify-send --version
```

### Command-line Options

```
Usage:
  wsl-notify-send [flags] <title> [message]

Flags:
  -a, --alert             Send alert notification with sound
      --app-name string   Application name (default "wsl-notify-send")
  -b, --beep              Just beep (no notification)
      --duration int      Beep duration in milliseconds (default 500)
      --freq float        Beep frequency in Hz (default 587)
  -h, --help              help for wsl-notify-send
  -i, --icon string       Icon file path or stock icon name
  -q, --quiet             Suppress error output
      --version           Show version information
```

## Icon Support

The tool supports various icon formats:
- **File paths**: PNG, JPG, JPEG, ICO, BMP files
- **Stock icons**: Platform-specific stock icon names
- **Embedded data**: Icons can be embedded as byte data

Examples:
```bash
# Use a PNG file
wsl-notify-send --icon /path/to/icon.png "Title" "Message"

# Use a stock icon (platform-specific)
wsl-notify-send --icon "warning" "Alert" "Warning message"
```

## Exit Codes

- `0`: Success
- `1`: General error
- `2`: Invalid arguments or configuration
- `3`: Notification failed to send

## Dependencies

- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [beeep](https://github.com/gen2brain/beeep) - Cross-platform notifications

## Platform Support

- **Windows 10/11**: Uses Windows Runtime COM API, falls back to PowerShell
- **Windows 7**: Uses win32 API
- **WSL2**: Forwards notifications to Windows host

## Error Handling

The tool provides clear error messages and appropriate exit codes:

```bash
# Invalid arguments
wsl-notify-send
# Error: requires at least a title argument

# Invalid icon file
wsl-notify-send --icon nonexistent.png "Title" "Message"
# Error: invalid configuration: icon file does not exist: nonexistent.png
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

[Add your license information here]

## Acknowledgments

- [beeep](https://github.com/gen2brain/beeep) - For the cross-platform notification library
- [Cobra](https://github.com/spf13/cobra) - For the CLI framework