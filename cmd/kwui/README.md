# kwui - Kenwood Radio GUI

A graphical interface for monitoring and controlling Kenwood radios using SDL2.

## System Dependencies

kwui requires the following system libraries to be installed:

### Fedora/RHEL
```bash
sudo dnf install SDL2-devel SDL2_ttf-devel
```

### Debian/Ubuntu
```bash
sudo apt-get install libsdl2-dev libsdl2-ttf-dev
```

### Arch Linux
```bash
sudo pacman -S sdl2 sdl2_ttf
```

## Building

Once dependencies are installed:

```bash
make kwui
```

Or directly:

```bash
go build -o kwui ./cmd/kwui
```

## Running

### On X11/Wayland (windowed mode)
```bash
./kwui
```

### On Linux console (framebuffer)
SDL2 will automatically detect and use the framebuffer when run from a Linux console (not in a terminal emulator).

## Configuration

### kwctl Command

By default, kwui executes `kwctl` to fetch radio status. You can configure a different command line using:

**Command-line flag:**
```bash
./kwui --kwctl-command 'ssh radio bin/kwctl -d /dev/radio0 -b 57600'
# or short form:
./kwui -k 'ssh radio bin/kwctl -d /dev/radio0 -b 57600'
# or just a simple path:
./kwui -k /path/to/kwctl
```

**Environment variable:**
```bash
export KWCTL_PATH='ssh radio bin/kwctl -d /dev/radio0 -b 57600'
./kwui
```

The value is treated as a command line and properly parsed for shell quoting. The `status` subcommand is automatically appended.

The command-line flag takes precedence over the environment variable.

## Features

- Real-time radio status display (updates every second via `kwctl status`)
- Dual VFO display with large frequency readout
- Retro amber terminal aesthetic
- TX/RX and control VFO indicators
- Mode, channel, and power display
- Keyboard controls:
  - Q or ESC: Quit
  - A: Switch to VFO A
  - B: Switch to VFO B

## Architecture

kwui uses SDL2 for cross-platform graphics that work on both X11/Wayland and Linux framebuffer:

- `main.go` - SDL initialization and main loop
- `display.go` - Rendering logic and event handling
- `model.go` - Application state and radio status polling
- `colors.go` - Color scheme definitions (amber retro theme)

The application executes `kwctl status` as a subprocess to fetch radio state, maintaining separation from the radio communication logic.
