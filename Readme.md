# Mookie Launcher

A simple cli menu tool to launch a set of predefined commands.

![mkl gif](img/mkl.gif)

Mookie launcher allows you to specify a menu structure in JSON format and execute the selected command. It can be useful for quickly launching hard to remember oneliners, bash scripts, programs, etc.

## Installation

With Go installed:

```bash
go install github.com/make0x20/mkl@latest
```

Optionally create a menu.json inside `~/config/mookie-launcher`. Example menu:

```json
{
  "List directory": "ls -alh",
  "Code editor": "nvim",
  "System - submenu": {
    "Update system": "sudo pacman -Syu",
    "Disk usage": "df -h"
  },
  "Some bash script": "~/path/some-script.sh",
  "System monitor": "htop"
}
```

## Usage

```bash
# simple usage - will look for menu.json inside ~/config/mookie-launcher
mkl

# specify a menu file
mkl -m <path to menu JSON file>

# set custom menu name
mkl -n <menu name>

# print the selected command to stdout instead of executing it - useful for piping to other commands
mkl -p
```

Search through entries using / - courtesy of [huh](https://github.com/charmbracelet/huh).
